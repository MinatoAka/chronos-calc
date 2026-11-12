package chronoscalc

import "testing"

func TestCalculateSMANotEnoughData(t *testing.T) {
	closes := []float64{1, 2, 3}
	if CalculateSMA(closes, 20) != nil {
		t.Error("harusnya nil kalau data kurang")
	}
}

func TestCalculateSMA(t *testing.T) {
	closes := []float64{1, 2, 3, 4, 5}
	sma := CalculateSMA(closes, 5)
	if sma == nil || *sma != 3.0 {
		t.Errorf("SMA = %v, mau 3.0", sma)
	}
}

func TestComputeTrend(t *testing.T) {
	var rows []Candle
	for i := 1; i <= 250; i++ {
		rows = append(rows, Candle{
			Date:  dateAdd(i),
			Close: float64(100 + i),
		})
	}
	trend := ComputeTrend(rows)
	if trend.Error != "" {
		t.Fatalf("Error: %s", trend.Error)
	}
	if trend.SMA20 == nil || trend.SMA50 == nil || trend.SMA200 == nil {
		t.Fatal("SMA harus terisi untuk 250 data")
	}
}

// helper: generate tanggal "YYYY-MM-DD" dari index hari
func dateAdd(day int) string {
	return "2024-01-01" // placeholder — tidak dipakai untuk logic, cukup konsisten sorting
}
