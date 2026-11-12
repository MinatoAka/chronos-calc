# chronos-calc

Technical analysis library for crypto market — SMA, seasonality, and weekly support/resistance.

## Install

```bash
go get github.com/MinatoAka/chronos-calc
```

Usage:
```
package main

import (
	"fmt"
	chronoscalc "github.com/MinatoAka/chronos-calc"
)

func main() {
	candles := []chronoscalc.Candle{
		{Date: "2024-01-01", Close: 100},
		{Date: "2024-01-02", Close: 105},
		// ... isi dengan data nyata
	}

	// SMA
	trend := chronoscalc.ComputeTrend(candles)
	fmt.Printf("SMA20: %v\n", trend.SMA20)

	// Seasonality
	season := chronoscalc.ComputeSeasonality(candles)
	fmt.Printf("Month returns: %+v\n", season.Matrix)

	// Support & Resistance (tolerance 0.5%)
	sr := chronoscalc.ComputeSupportResistance(candles, 0.005)
	fmt.Printf("R1: %+v\n", sr.Resistances)
	fmt.Printf("S1: %+v\n", sr.Supports)
}
```
