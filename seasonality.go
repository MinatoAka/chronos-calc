package chronoscalc

import (
	"sort"
	"strconv"
)

// ComputeSeasonality — hitung return bulanan per bulan (Mean, Median, Winrate, N).
func ComputeSeasonality(rows []Candle) SeasonalityResult {
	if len(rows) < 365 {
		return SeasonalityResult{Error: "Data tidak cukup (min 365 hari)"}
	}

	// FIX 5: Clone & sort untuk konsistensi
	sortedRows := make([]Candle, len(rows))
	copy(sortedRows, rows)
	sort.Slice(sortedRows, func(i, j int) bool {
		return sortedRows[i].Date < sortedRows[j].Date
	})

	// Ambil close terakhir per bulan (YYYY-MM)
	monthClose := make(map[string]float64)
	for _, r := range sortedRows {
		// FIX 6: Validasi format tanggal minimal YYYY-MM (7 char) sebelum slicing
		if len(r.Date) >= 7 && r.Date[4] == '-' {
			monthClose[r.Date[:7]] = r.Close
		}
	}

	var sortedYM []string
	for ym := range monthClose {
		sortedYM = append(sortedYM, ym)
	}
	sort.Strings(sortedYM)

	// Return bulanan
	var monthlyReturns []struct {
		ym  string
		ret float64
	}
	for i := 1; i < len(sortedYM); i++ {
		prev, cur := monthClose[sortedYM[i-1]], monthClose[sortedYM[i]]
		if prev > 0 {
			monthlyReturns = append(monthlyReturns, struct {
				ym  string
				ret float64
			}{sortedYM[i], (cur - prev) / prev})
		}
	}

	// Kelompokkan return per bulan (1-12)
	byMonth := make(map[int][]float64)
	for _, m := range monthlyReturns {
		// FIX 7: Safe parsing bulan, fallback jika format aneh
		if len(m.ym) >= 7 {
			mIdx, err := strconv.Atoi(m.ym[5:7])
			if err == nil && mIdx >= 1 && mIdx <= 12 {
				byMonth[mIdx] = append(byMonth[mIdx], m.ret)
			}
		}
	}

	monthNames := []string{"Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Agu", "Sep", "Okt", "Nov", "Des"}
	var matrix []SeasonalityRow

	for m := 1; m <= 12; m++ {
		arr := byMonth[m]
		if len(arr) == 0 {
			matrix = append(matrix, SeasonalityRow{Bulan: monthNames[m-1], N: 0})
			continue
		}

		sort.Float64s(arr)
		var sum float64
		for _, v := range arr {
			sum += v
		}
		mean := sum / float64(len(arr))

		var median float64
		if len(arr)%2 == 1 {
			median = arr[(len(arr)-1)/2]
		} else {
			median = (arr[len(arr)/2-1] + arr[len(arr)/2]) / 2
		}

		win := 0
		for _, v := range arr {
			if v > 0 {
				win++
			}
		}
		winrate := float64(win) / float64(len(arr))

		matrix = append(matrix, SeasonalityRow{
			Bulan: monthNames[m-1], Mean: &mean, Median: &median, Winrate: &winrate, N: len(arr),
		})
	}

	yearsSpan := ""
	if len(sortedYM) > 0 {
		yearsSpan = sortedYM[0] + " s/d " + sortedYM[len(sortedYM)-1]
	}

	return SeasonalityResult{
		Matrix:      matrix,
		MonthsTotal: len(monthlyReturns),
		YearsSpan:   yearsSpan,
	}
}
