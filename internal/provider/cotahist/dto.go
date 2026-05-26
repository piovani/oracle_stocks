package cotahist

import "time"

type Record struct {
	Date          time.Time  // DATA
	BDICode       string     // CODBDI — market segment code
	Ticker        string     // CODNEG — trading code
	MarketType    string     // TPMERC — 010=cash, 020=fractional, 070=options call, 080=options put, ...
	ShortName     string     // NOMRES — short company name
	Specification string     // ESPECI — share class (ON, PN, UNT, ...)
	Currency      string     // MODREF — reference currency
	OpenPrice     float64    // PREABE
	HighPrice     float64    // PREMAX
	LowPrice      float64    // PREMIN
	AveragePrice  float64    // PREMED
	ClosePrice    float64    // PREULT
	BestBid       float64    // PREOFC
	BestAsk       float64    // PREOFV
	Trades        int64      // TOTNEG — number of trades
	Quantity      int64      // QUATOT — total shares traded
	Volume        float64    // VOLTOT — total volume in BRL
	StrikePrice   float64    // PREEXE — option strike
	ExpiryDate    *time.Time // DATVEN — option expiry (nil for non-derivatives)
	QuoteFactor   int64      // FATCOT — quote factor (usually 1 or 1000)
	ISIN          string     // CODISI
	DistNumber    string     // DISMES — distribution number
}
