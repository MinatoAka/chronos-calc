[![Go Reference](https://pkg.go.dev/badge/github.com/MinatoAka/chronos-calc.svg)](https://pkg.go.dev/github.com/MinatoAka/chronos-calc)
[![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8?logo=go)](https://go.dev/)
[![Saweria](https://img.shields.io/badge/Support%20Me-Saweria-ff69b4?logo=coffee&logoColor=white)](https://saweria.co/AKaAturDompet)

# chronos-calc

Technical analysis library for crypto market — SMA, seasonality, and weekly support/resistance.

## Features

- **Simple Moving Average (SMA)** — Calculate SMA 20/50/200 for trend analysis
- **Seasonality Analysis** — Monthly return statistics (mean, median, winrate)
- **Support & Resistance** — Weekly swing high/low clustering with configurable tolerance
- **Pure Functions** — No side effects, input data remains unchanged
- **Zero Dependencies** — Standard library only

## Install

```bash
go get github.com/MinatoAka/chronos-calc
```

### Usage:

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
  "Pair": "BTCIDR",
  "AsOf": "2024-12-31",
  "Close": 660000000,
  "SMA20": 655200000,
  "SMA50": 648500000,
  "SMA200": 635800000,
  "RowsCount": 1825,
  "Error": ""
}
```

**Case: Less than 200 days of data**
```json
{
  "Pair": "BTCIDR",
  "AsOf": "2024-12-31",
  "Close": 660000000,
  "SMA20": 655200000,
  "SMA50": 648500000,
  "SMA200": null,
  "RowsCount": 150,
  "Error": ""
}
```
---

### ComputeSeasonality — Monthly Return Statistics

```go
season := chronoscalc.ComputeSeasonality(candles)

jsonBytes, _ := json.MarshalIndent(season, "", "  ")
fmt.Println(string(jsonBytes))
```

**Output:**
```json
{
  "Matrix": [
    {
      "Bulan": "Jan",
      "Mean": 0.025,
      "Median": 0.018,
      "Winrate": 0.6,
      "N": 5
    },
    {
      "Bulan": "Feb",
      "Mean": 0.032,
      "Median": 0.029,
      "Winrate": 0.8,
      "N": 5
    },
    {
      "Bulan": "Mar",
      "Mean": -0.015,
      "Median": -0.012,
      "Winrate": 0.2,
      "N": 5
    },
    {
      "Bulan": "Apr",
      "Mean": 0.058,
      "Median": 0.061,
      "Winrate": 1.0,
      "N": 5
    },
    {
      "Bulan": "Mei",
      "Mean": -0.023,
      "Median": -0.021,
      "Winrate": 0.0,
      "N": 5
    },
    {
      "Bulan": "Jun",
      "Mean": 0.012,
      "Median": 0.009,
      "Winrate": 0.6,
      "N": 5
    },
    {


      "Bulan": "Jul",
      "Mean": 0.045,
      "Median": 0.042,
      "Winrate": 0.8,
      "N": 5
    },
    {
      "Bulan": "Agu",
      "Mean": -0.008,
      "Median": -0.005,
      "Winrate": 0.4,
      "N": 5
    },
    {
      "Bulan": "Sep",
      "Mean": 0.021,
      "Median": 0.019,
      "Winrate": 0.6,
      "N": 5
    },
    {
      "Bulan": "Okt",
      "Mean": 0.063,
      "Median": 0.065,
      "Winrate": 1.0,
      "N": 5
    },
    {
      "Bulan": "Nov",
      "Mean": 0.038,
      "Median": 0.035,
      "Winrate": 0.8,
      "N": 5
    },
    {
      "Bulan": "Des",
      "Mean": -0.011,
      "Median": -0.009,
      "Winrate": 0.2,
      "N": 5
    }
  ],
  "MonthsTotal": 60,
  "YearsSpan": "2020-01 s/d 2024-12",
  "Insight": "",
  "Error": ""
}
```

**Field Descriptions:**
- `Mean` — Average monthly return (e.g., 0.025 = 2.5%)
- `Median` — Median monthly return
- `Winrate` — Percentage of positive months (e.g., 0.6 = 60%)
- `N` — Number of years with data for that month

**Case: Month without data (N=0)**
```json
{
  "Bulan": "Feb",
  "Mean": null,
  "Median": null,
  "Winrate": null,
  "N": 0
}
```
---

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
  "CurrentPrice": 660000000,
  "Resistances": [
    {
      "Price": 670500000,
      "Count": 3
    },
    {
      "Price": 685200000,
      "Count": 2
    }
  ],
  "Supports": [
    {
      "Price": 650300000,
      "Count": 4
    },
    {
      "Price": 635800000,
      "Count": 2
    }
  ],
  "Error": ""
}
```

**Field Descriptions:**
- `CurrentPrice` — Latest close price
- `Resistances` — Top 2 resistance levels above current price (sorted ascending)
- `Supports` — Top 2 support levels below current price (sorted descending)
- `Count` — Number of weekly touches at that level

**Case: Insufficient data**
```json
{
  "CurrentPrice": 0,
  "Resistances": null,
  "Supports": null,
  "Error": "Data mingguan tidak cukup"
}
```
---

## API Reference

| Function | Signature | Description |
|---|---|---|
| `CalculateSMA` | `CalculateSMA(closes []float64, period int) *float64` | Simple Moving Average of last N values |
| `ComputeTrend` | `ComputeTrend(rows []Candle) TrendResult` | SMA 20/50/200 + metadata |
| `ComputeSeasonality` | `ComputeSeasonality(rows []Candle) SeasonalityResult` | Monthly return statistics (min 365 days) |
| `ComputeSupportResistance` | `ComputeSupportResistance(rows []Candle, tolerancePct float64) SRResult` | Weekly swing high/low clustering |

## Requirements

- **ComputeTrend**: Minimum 2 candles (returns nil for SMA if data < period)
- **ComputeSeasonality**: Minimum 365 days of data
- **ComputeSupportResistance**: Minimum 10 weeks of data

## Thread Safety

All functions are safe for concurrent use. Input slices are cloned internally to prevent side effects.

## Test

```bash
go test ./...
```
---
