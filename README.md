# chronos-calc

[![Go Reference](https://pkg.go.dev/badge/github.com/MinatoAka/chronos-calc.svg)](https://pkg.go.dev/github.com/MinatoAka/chronos-calc)
[![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8?logo=go)](https://go.dev/)
[![Saweria](https://img.shields.io/badge/Support%20Me-Saweria-ff69b4?logo=coffee&logoColor=white)](https://saweria.co/AKaAturDompet)

Technical analysis library for crypto market — SMA, seasonality, and weekly support/resistance.

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

### FEATURES
| Function | Purpose |
| --- | --- |
| `CalculateSMA(closes, period)` | Simple moving average |
| `ComputeTrend(candles)` | SMA 20/50/200 + metadata |
| `ComputeSeasonality(candles)` | Monthly returns (mean, median, winrate) |
| `ComputeSupportResistance(candles, tolerancePct) `| Weekly swing high/low clustering |

### Notes
- `ComputeSupportResistance requires tolerancePct` (decimal, e.g. `0.005` for 0.5%).
- `SeasonalityResult.Insight` is left empty — caller fills it if needed.
- All functions are stateless & pure.
