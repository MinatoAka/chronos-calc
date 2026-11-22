<div align="center">
# chronos-calc

Calculator technical analysis library for crypto market — SMA, EMA, seasonality, and weekly support/resistance.

[![Go Reference](https://pkg.go.dev/badge/github.com/MinatoAka/chronos-calc.svg)](https://pkg.go.dev/github.com/MinatoAka/chronos-calc)
[![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8?logo=go)](https://go.dev/)
[![Support Me](https://img.shields.io/badge/Support%20Me-Saweria-ff69b4?logo=coffee&logoColor=white)](https://saweria.co/)

</div>
---


## Features

* **Simple Moving Average (SMA)** — Calculate SMA 20/50/200 for trend analysis
* **Exponential Moving Average (EMA)** — Calculate EMA 20/50 and detect price cross (Bullish/Bearish Cross)
* **Seasonality Analysis** — Monthly return statistics (mean, median, winrate)
* **Support & Resistance** — Weekly swing high/low clustering with configurable tolerance
* **Pure Functions** — No side effects, input data remains unchanged
* **Zero Dependencies** — Standard library only

## Install

```bash
go get github.com/MinatoAka/chronos-calc
```

## Test

```bash
go test ./...
```

## Usage

### ComputeTrend — SMA Analysis

```go
package main

import (
	"encoding/json"
	"fmt"
	chronoscalc "github.com/MinatoAka/chronos-calc"
)

func main() {
	candles := []chronoscalc.Candle{
		{Date: "2024-01-01", Pair: "BTCIDR", Open: 650000000, High: 660000000, Low: 645000000, Close: 655000000, Volume: 100},
		{Date: "2024-01-02", Pair: "BTCIDR", Open: 655000000, High: 665000000, Low: 650000000, Close: 660000000, Volume: 120},
		// ... add more candles (minimum 2 for basic SMA, 200 for SMA200)
	}

	trend := chronoscalc.ComputeTrend(candles)
	jsonBytes, _ := json.MarshalIndent(trend, "", "  ")
	fmt.Println(string(jsonBytes))
}
```

**Output:**
```json
{
  "pair": "BTCIDR",
  "as_of": "2024-12-31",
  "close": 660000000,
  "sma20": 655200000,
  "sma50": 648500000,
  "sma200": 635800000,
  "rows_count": 1825,
  "error": ""
}
```

### ComputeEMATrend & ComputeEMACross — EMA Analysis

```go
// 1. Cek Tren (EMA20 vs EMA50)
trend := chronoscalc.ComputeEMATrend(candles)
jsonBytes, _ := json.MarshalIndent(trend, "", "  ")
fmt.Println(string(jsonBytes))

// 2. Cek Cross (Harga nembus EMA20)
cross := chronoscalc.ComputeEMACross(candles)
jsonBytes, _ = json.MarshalIndent(cross, "", "  ")
fmt.Println(string(jsonBytes))
```

**Output (Trend):**
```json
{
  "pair": "BTCIDR",
  "as_of": "2024-12-31",
  "close": 660000000,
  "ema20": 658000000,
  "ema50": 645000000,
  "trend": "BULLISH",
  "rows_count": 150,
  "error": ""
}
```

**Output (Cross):**
```json
{
  "pair": "BTCIDR",
  "as_of": "2024-12-31",
  "close": 660000000,
  "ema20": 658000000,
  "cross_type": "BULLISH_CROSS",
  "error": ""
}
```
*(CrossType values: `BULLISH_CROSS`, `BEARISH_CROSS`, `NONE`)*

### ComputeSeasonality — Monthly Return Statistics

```go
season := chronoscalc.ComputeSeasonality(candles)
jsonBytes, _ := json.MarshalIndent(season, "", "  ")
fmt.Println(string(jsonBytes))
```

**Output:**
```json
{
  "matrix": [
    {
      "bulan": "Jan",
      "mean": 0.025,
      "median": 0.018,
      "winrate": 0.6,
      "n": 5
    }
    // ... other months
  ],
  "months_total": 60,
  "years_span": "2020-01 s/d 2024-12",
  "insight": "",
  "error": ""
}
```

### ComputeSupportResistance — Weekly S&R Levels

```go
// tolerancePct: clustering tolerance (e.g., 0.005 = 0.5%)
sr := chronoscalc.ComputeSupportResistance(candles, 0.005)
jsonBytes, _ := json.MarshalIndent(sr, "", "  ")
fmt.Println(string(jsonBytes))
```

**Output:**
```json
{
  "current_price": 660000000,
  "resistances": [
    { "price": 670500000, "count": 3 },
    { "price": 685200000, "count": 2 }
  ],
  "supports": [
    { "price": 650300000, "count": 4 },
    { "price": 635800000, "count": 2 }
  ],
  "error": ""
}
```

## API Reference

| Function | Signature | Description |
| :--- | :--- | :--- |
| `CalculateSMA` | `CalculateSMA(closes []float64, period int) *float64` | Simple Moving Average of last N values |
| `CalculateEMA` | `CalculateEMA(closes []float64, period int) *float64` | Exponential Moving Average (SMA seeded) |
| `ComputeTrend` | `ComputeTrend(rows []Candle) TrendResult` | SMA 20/50/200 + metadata |
| `ComputeEMATrend`| `ComputeEMATrend(rows []Candle) EMATrendResult`| Compare EMA20 vs EMA50 (BULLISH/BEARISH) |
| `ComputeEMACross`| `ComputeEMACross(rows []Candle) EMACrossResult`| Detect price cross over EMA20 (BULLISH_CROSS/BEARISH_CROSS) |
| `ComputeSeasonality`| `ComputeSeasonality(rows []Candle) SeasonalityResult` | Monthly return statistics (min 365 days) |
| `ComputeSupportResistance`| `ComputeSupportResistance(rows []Candle, tolerancePct float64) SRResult` | Weekly swing high/low clustering |

## Requirements

* `ComputeTrend`: Minimum 2 candles (returns null for SMA if data < period)
* `ComputeEMATrend`: Minimum 50 candles
* `ComputeEMACross`: Minimum 21 candles
* `ComputeSeasonality`: Minimum 365 days of data
* `ComputeSupportResistance`: Minimum 10 weeks of data

## Thread Safety

All functions are safe for concurrent use. Input slices are cloned internally to prevent side effects.