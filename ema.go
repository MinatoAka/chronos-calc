package chronoscalc

import "sort"

// CalculateEMA — Exponential Moving Average dari data closes.
// Return nil jika data kurang dari period.
// Menggunakan rata-rata sederhana (SMA) sebagai seed awal untuk akurasi.
func CalculateEMA(closes []float64, period int) *float64 {
	if len(closes) < period {
		return nil
	}

	// Multiplier (k)
	k := 2.0 / float64(period+1)

	// Seed dengan rata-rata (SMA) dari data awal
	var sum float64
	for i := 0; i < period; i++ {
		sum += closes[i]
	}
	ema := sum / float64(period)

	// Hitung EMA secara berurutan sampai data terakhir
	for i := period; i < len(closes); i++ {
		ema = closes[i]*k + ema*(1-k)
	}

	return &ema
}

// calculateEMASeries — return slice EMA sepanjang data (khusus buat deteksi cross).
func calculateEMASeries(closes []float64, period int) []float64 {
	if len(closes) < period {
		return nil
	}

	k := 2.0 / float64(period+1)
	result := make([]float64, len(closes))

	// Seed
	var sum float64
	for i := 0; i < period; i++ {
		sum += closes[i]
	}
	smaSeed := sum / float64(period)

	// Isi index < period dengan seed (biar slice utuh)
	for i := 0; i < period; i++ {
		result[i] = smaSeed
	}

	// Hitung EMA
	for i := period; i < len(closes); i++ {
		result[i] = closes[i]*k + result[i-1]*(1-k)
	}

	return result
}

// ComputeEMATrend — Mode 1: Bandingkan EMA20 & EMA50 terakhir.
func ComputeEMATrend(rows []Candle) EMATrendResult {
	if len(rows) < 50 {
		return EMATrendResult{Error: "Data tidak cukup (min 50 candle)"}
	}

	// Clone & sort ASC (mencegah side-effect)
	sortedRows := make([]Candle, len(rows))
	copy(sortedRows, rows)
	sort.Slice(sortedRows, func(i, j int) bool {
		return sortedRows[i].Date < sortedRows[j].Date
	})

	closes := make([]float64, len(sortedRows))
	for i, r := range sortedRows {
		closes[i] = r.Close
	}

	last := sortedRows[len(sortedRows)-1]

	ema20 := CalculateEMA(closes, 20)
	ema50 := CalculateEMA(closes, 50)

	trend := "NEUTRAL"
	if ema20 != nil && ema50 != nil {
		if *ema20 > *ema50 {
			trend = "BULLISH"
		} else if *ema20 < *ema50 {
			trend = "BEARISH"
		}
	}

	return EMATrendResult{
		Pair:      last.Pair,
		AsOf:      last.Date,
		Close:     last.Close,
		EMA20:     ema20,
		EMA50:     ema50,
		Trend:     trend,
		RowsCount: len(sortedRows),
	}
}

// ComputeEMACross — Mode 2: Deteksi harga (Close) nembus EMA20 (Hari ini vs Kemarin).
func ComputeEMACross(rows []Candle) EMACrossResult {
	// Butuh min 21 candle (20 buat hitung EMA pertama, +1 buat bandingin kemarin)
	if len(rows) < 21 {
		return EMACrossResult{Error: "Data tidak cukup (min 21 candle untuk deteksi cross)"}
	}

	// Clone & sort ASC
	sortedRows := make([]Candle, len(rows))
	copy(sortedRows, rows)
	sort.Slice(sortedRows, func(i, j int) bool {
		return sortedRows[i].Date < sortedRows[j].Date
	})

	closes := make([]float64, len(sortedRows))
	for i, r := range sortedRows {
		closes[i] = r.Close
	}

	// Hitung series EMA20
	ema20Series := calculateEMASeries(closes, 20)
	if ema20Series == nil {
		return EMACrossResult{Error: "Gagal menghitung EMA series"}
	}

	lastIdx := len(closes) - 1
	prevIdx := lastIdx - 1

	lastClose := closes[lastIdx]
	prevClose := closes[prevIdx]
	lastEMA := ema20Series[lastIdx]
	prevEMA := ema20Series[prevIdx]

	// Logika Cross:
	// Bullish: Kemarin Close < EMA, Hari ini Close > EMA
	// Bearish: Kemarin Close > EMA, Hari ini Close < EMA
	crossType := "NONE"
	if prevClose < prevEMA && lastClose > lastEMA {
		crossType = "BULLISH_CROSS"
	} else if prevClose > prevEMA && lastClose < lastEMA {
		crossType = "BEARISH_CROSS"
	}

	last := sortedRows[len(sortedRows)-1]
	return EMACrossResult{
		Pair:      last.Pair,
		AsOf:      last.Date,
		Close:     last.Close,
		EMA20:     &lastEMA,
		CrossType: crossType,
	}
}