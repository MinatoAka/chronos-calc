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
// Data di-sort ASC dulu untuk memastikan perhitungan benar.
func ComputeTrend(rows []Candle) TrendResult {
	if len(rows) < 2 {
		return TrendResult{Error: "Data tidak cukup"}
	}

	// Defensive sorting
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Date < rows[j].Date
	})

	closes := make([]float64, len(rows))
	for i, r := range rows {
		closes[i] = r.Close
	}
	last := rows[len(rows)-1]

	return TrendResult{
		Pair:      last.Pair,
		AsOf:      last.Date,
		Close:     last.Close,
		SMA20:     CalculateSMA(closes, 20),
		SMA50:     CalculateSMA(closes, 50),
		SMA200:    CalculateSMA(closes, 200),
		RowsCount: len(rows),
	}
}
