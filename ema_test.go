package chronoscalc

import (
	"testing"
	"time"
)

// generateCandles adalah helper buat bikin data dummy biar test rapi.
func generateCandles(count int, basePrice float64) []Candle {
	candles := make([]Candle, count)
	baseTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < count; i++ {
		candles[i] = Candle{
			Date:   baseTime.AddDate(0, 0, i).Format("2006-01-02"),
			Pair:   "BTCIDR",
			Open:   basePrice,
			High:   basePrice + 10,
			Low:    basePrice - 10,
			Close:  basePrice,
			Volume: 1000,
		}
	}
	return candles
}

func TestCalculateEMA(t *testing.T) {
	// 1. Test data kurang dari period (harus return nil)
	if res := CalculateEMA([]float64{1, 2, 3}, 5); res != nil {
		t.Errorf("Expected nil untuk data < period, got %v", *res)
	}

	// 2. Test data cukup (harus return nilai)
	closes := make([]float64, 10)
	for i := range closes {
		closes[i] = float64(i + 1)
	}
	if res := CalculateEMA(closes, 5); res == nil {
		t.Error("Expected non-nil EMA untuk data >= period")
	}
}

func TestComputeEMATrend(t *testing.T) {
	// 1. Test Error (data < 50)
	resErr := ComputeEMATrend(generateCandles(40, 100))
	if resErr.Error == "" {
		t.Error("Expected error untuk data < 50 candle")
	}

	// 2. Test BULLISH (Uptrend: harga naik terus)
	uptrend := generateCandles(50, 100)
	for i := range uptrend {
		uptrend[i].Close = float64(100 + i*2) // Naik 2 poin tiap candle
	}
	resBull := ComputeEMATrend(uptrend)
	if resBull.Error != "" || resBull.Trend != "BULLISH" {
		t.Errorf("Expected BULLISH, got %s (err: %s)", resBull.Trend, resBull.Error)
	}

	// 3. Test BEARISH (Downtrend: harga turun terus)
	downtrend := generateCandles(50, 200)
	for i := range downtrend {
		downtrend[i].Close = float64(200 - i*2) // Turun 2 poin tiap candle
	}
	resBear := ComputeEMATrend(downtrend)
	if resBear.Error != "" || resBear.Trend != "BEARISH" {
		t.Errorf("Expected BEARISH, got %s (err: %s)", resBear.Trend, resBear.Error)
	}
}

func TestComputeEMACross(t *testing.T) {
	// 1. Test Error (data < 21)
	resErr := ComputeEMACross(generateCandles(20, 100))
	if resErr.Error == "" {
		t.Error("Expected error untuk data < 21 candle")
	}

	// 2. Test BULLISH_CROSS (Kemarin Close < EMA, Hari ini Close > EMA)
	// Kita bikin 25 candle. 20 candle pertama stabil di harga 100 (biar EMA20 ~100).
	candlesBull := generateCandles(25, 100)
	// Candle ke-24 (index 23 / prev): banting harga ke 80 (di bawah EMA)
	candlesBull[23].Close = 80
	// Candle ke-25 (index 24 / last): pump harga ke 120 (di atas EMA)
	candlesBull[24].Close = 120

	resBullCross := ComputeEMACross(candlesBull)
	if resBullCross.Error != "" || resBullCross.CrossType != "BULLISH_CROSS" {
		t.Errorf("Expected BULLISH_CROSS, got %s (err: %s)", resBullCross.CrossType, resBullCross.Error)
	}

	// 3. Test BEARISH_CROSS (Kemarin Close > EMA, Hari ini Close < EMA)
	candlesBear := generateCandles(25, 100)
	// Candle ke-24 (prev): pump harga ke 120 (di atas EMA)
	candlesBear[23].Close = 120
	// Candle ke-25 (last): banting harga ke 80 (di bawah EMA)
	candlesBear[24].Close = 80

	resBearCross := ComputeEMACross(candlesBear)
	if resBearCross.Error != "" || resBearCross.CrossType != "BEARISH_CROSS" {
		t.Errorf("Expected BEARISH_CROSS, got %s (err: %s)", resBearCross.CrossType, resBearCross.Error)
	}

	// 4. Test NONE (Tidak ada cross, harga stabil di sekitar EMA)
	candlesNone := generateCandles(25, 100)
	// Biarkan 2 candle terakhir tetap di 100 (sama dengan EMA)
	resNone := ComputeEMACross(candlesNone)
	if resNone.Error != "" || resNone.CrossType != "NONE" {
		t.Errorf("Expected NONE, got %s (err: %s)", resNone.CrossType, resNone.Error)
	}
}