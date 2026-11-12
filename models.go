package chronoscalc

// Candle — satu candle OHLCV.
type Candle struct {
	Date   string  `json:"date"`
	Pair   string  `json:"pair"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

// TrendResult — hasil analisis tren (SMA 20/50/200).
type TrendResult struct {
	Pair      string
	AsOf      string
	Close     float64
	SMA20     *float64
	SMA50     *float64
	SMA200    *float64
	RowsCount int
	Error     string
}

// SeasonalityRow — satu baris matriks musiman (per bulan).
type SeasonalityRow struct {
	Bulan   string
	Mean    *float64
	Median  *float64
	Winrate *float64
	N       int
}

// SeasonalityResult — hasil analisis musiman.
type SeasonalityResult struct {
	Matrix      []SeasonalityRow
	MonthsTotal int
	YearsSpan   string
	Insight     string // diisi oleh caller (tidak dihitung di library)
	Error       string
}

// SRLevel — satu level support atau resistance.
type SRLevel struct {
	Price float64
	Count int
}

// SRResult — hasil analisis support & resistance mingguan.
type SRResult struct {
	CurrentPrice float64
	Resistances  []SRLevel
	Supports     []SRLevel
	Error        string
}
