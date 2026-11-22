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
	Pair      string   `json:"pair"`
	AsOf      string   `json:"as_of"`
	Close     float64  `json:"close"`
	SMA20     *float64 `json:"sma20"`
	SMA50     *float64 `json:"sma50"`
	SMA200    *float64 `json:"sma200"`
	RowsCount int      `json:"rows_count"`
	Error     string   `json:"error,omitempty"`
}

// SeasonalityRow — satu baris matriks musiman (per bulan).
type SeasonalityRow struct {
	Bulan   string   `json:"bulan"`
	Mean    *float64 `json:"mean"`
	Median  *float64 `json:"median"`
	Winrate *float64 `json:"winrate"`
	N       int      `json:"n"`
}

// SeasonalityResult — hasil analisis musiman.
type SeasonalityResult struct {
	Matrix      []SeasonalityRow `json:"matrix"`
	MonthsTotal int              `json:"months_total"`
	YearsSpan   string           `json:"years_span"`
	Insight     string           `json:"insight,omitempty"` // diisi oleh caller
	Error       string           `json:"error,omitempty"`
}

// SRLevel — satu level support atau resistance.
type SRLevel struct {
	Price float64 `json:"price"`
	Count int     `json:"count"`
}

// SRResult — hasil analisis support & resistance mingguan.
type SRResult struct {
	CurrentPrice float64   `json:"current_price"`
	Resistances  []SRLevel `json:"resistances"`
	Supports     []SRLevel `json:"supports"`
	Error        string    `json:"error,omitempty"`
}

// EMATrendResult — hasil perbandingan EMA20 & EMA50.
type EMATrendResult struct {
	Pair      string   `json:"pair"`
	AsOf      string   `json:"as_of"`
	Close     float64  `json:"close"`
	EMA20     *float64 `json:"ema20"`
	EMA50     *float64 `json:"ema50"`
	Trend     string   `json:"trend"` // "BULLISH" | "BEARISH" | "NEUTRAL"
	RowsCount int      `json:"rows_count"`
	Error     string   `json:"error,omitempty"`
}

// EMACrossResult — hasil deteksi harga (Close) nembus EMA20.
type EMACrossResult struct {
	Pair      string  `json:"pair"`
	AsOf      string  `json:"as_of"`
	Close     float64 `json:"close"`
	EMA20     *float64 `json:"ema20"`
	CrossType string  `json:"cross_type"` // "BULLISH_CROSS" | "BEARISH_CROSS" | "NONE"
	Error     string  `json:"error,omitempty"`
}