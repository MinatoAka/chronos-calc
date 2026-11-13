package chronoscalc

import "sort"

// CalculateSMA — Simple Moving Average dari N data terakhir.
// Return nil jika data kurang dari period.
func CalculateSMA(closes []float64, period int) *float64 {
	if len(closes) < period {
		return nil
	}
	slice := closes[len(closes)-period:]
	var sum float64
	for _, v := range slice {
		sum += v
	}
	res := sum / float64(period)
	return &res
}

// ComputeTrend — hitung SMA 20/50/200 dari data candle.
func ComputeTrend(rows []Candle) TrendResult {
	if len(rows) < 2 {
		return TrendResult{Error: "Data tidak cukup"}
	}

	// FIX 1: Clone slice untuk mencegah side-effect (mengubah urutan data asli caller)
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

	return TrendResult{
		Pair:      last.Pair,
		AsOf:      last.Date,
		Close:     last.Close,
		SMA20:     CalculateSMA(closes, 20),
		SMA50:     CalculateSMA(closes, 50),
		SMA200:    CalculateSMA(closes, 200),
		RowsCount: len(sortedRows),
	}
}
