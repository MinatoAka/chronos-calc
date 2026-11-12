package chronoscalc

import (
	"math"
	"sort"
	"time"
)

// getMonday — kembalikan tanggal Senin dari minggu tanggal t.
func getMonday(t time.Time) time.Time {
	day := int(t.Weekday())
	diff := 0
	if day == 0 {
		diff = -6
	} else {
		diff = 1 - day
	}
	return t.AddDate(0, 0, diff)
}

// AggregateToWeekly — gabungkan candle harian menjadi candle mingguan.
func AggregateToWeekly(rows []Candle) []Candle {
	var weekly []Candle
	var current *Candle
	var currentMonday time.Time

	for _, r := range rows {
		t, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			continue
		}
		monday := getMonday(t)

		if current == nil || !monday.Equal(currentMonday) {
			if current != nil {
				weekly = append(weekly, *current)
			}
			c := r
			current = &c
			currentMonday = monday
		} else {
			if r.High > current.High {
				current.High = r.High
			}
			if r.Low < current.Low {
				current.Low = r.Low
			}
			current.Close = r.Close
			current.Volume += r.Volume
		}
	}
	if current != nil {
		weekly = append(weekly, *current)
	}
	return weekly
}

// ComputeSupportResistance — cari level S&R dari swing high/low mingguan.
// tolerancePct: toleransi jarak (dalam desimal, misal 0.005 = 0.5%).
func ComputeSupportResistance(rows []Candle, tolerancePct float64) SRResult {
	// Defensive sorting
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Date < rows[j].Date
	})

	weekly := AggregateToWeekly(rows)
	if len(weekly) < 10 {
		return SRResult{Error: "Data mingguan tidak cukup"}
	}
	if len(rows) == 0 {
		return SRResult{Error: "Data harian kosong"}
	}
	currentPrice := rows[len(rows)-1].Close

	windowSize := 2
	var swingHighs, swingLows []float64

	for i := windowSize; i < len(weekly)-windowSize; i++ {
		isHigh, isLow := true, true
		for j := 1; j <= windowSize; j++ {
			if weekly[i].High <= weekly[i-j].High || weekly[i].High <= weekly[i+j].High {
				isHigh = false
			}
			if weekly[i].Low >= weekly[i-j].Low || weekly[i].Low >= weekly[i+j].Low {
				isLow = false
			}
		}
		if isHigh {
			swingHighs = append(swingHighs, weekly[i].High)
		}
		if isLow {
			swingLows = append(swingLows, weekly[i].Low)
		}
	}

	cluster := func(prices []float64) []SRLevel {
		if len(prices) == 0 {
			return nil
		}

		sort.Float64s(prices)
		var res []SRLevel
		var currentMin, currentMax float64
		var currentCount int

		for i, p := range prices {
			if i == 0 {
				currentMin, currentMax, currentCount = p, p, 1
			} else if p <= currentMin*1.02 { // anti-clustering creep
				currentMax = p
				currentCount++
			} else {
				res = append(res, SRLevel{Price: (currentMin + currentMax) / 2, Count: currentCount})
				currentMin, currentMax, currentCount = p, p, 1
			}
		}
		if currentCount > 0 {
			res = append(res, SRLevel{Price: (currentMin + currentMax) / 2, Count: currentCount})
		}
		return res
	}

	cH, cL := cluster(swingHighs), cluster(swingLows)
	tolerance := currentPrice * tolerancePct

	var resistances, supports []SRLevel
	for _, r := range cH {
		if r.Price > (currentPrice + tolerance) {
			resistances = append(resistances, r)
		}
	}
	for _, s := range cL {
		if s.Price < (currentPrice - tolerance) {
			supports = append(supports, s)
		}
	}

	sort.Slice(resistances, func(i, j int) bool { return resistances[i].Price < resistances[j].Price })
	sort.Slice(supports, func(i, j int) bool { return supports[i].Price > supports[j].Price })

	if len(resistances) > 2 {
		resistances = resistances[:2]
	}
	if len(supports) > 2 {
		supports = supports[:2]
	}

	return SRResult{
		CurrentPrice: currentPrice,
		Resistances:  resistances,
		Supports:     supports,
	}
}

// helper untuk mencegah unused import math (dipakai di cluster jika diperlukan)
var _ = math.Max
