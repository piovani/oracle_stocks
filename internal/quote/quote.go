package quote

import (
	"time"

	"github.com/allison-piovani/oracle_stocks/internal/provider/cotahist"
)

type Quote struct {
	ID            uint64     `gorm:"column:id;primaryKey"`
	Date          time.Time  `gorm:"column:date"`
	Ticker        string     `gorm:"column:ticker"`
	BDICode       string     `gorm:"column:bdi_code"`
	MarketType    string     `gorm:"column:market_type"`
	ShortName     *string    `gorm:"column:short_name"`
	Specification *string    `gorm:"column:specification"`
	Currency      *string    `gorm:"column:currency"`
	OpenPrice     *float64   `gorm:"column:open_price"`
	HighPrice     *float64   `gorm:"column:high_price"`
	LowPrice      *float64   `gorm:"column:low_price"`
	AveragePrice  *float64   `gorm:"column:average_price"`
	ClosePrice    *float64   `gorm:"column:close_price"`
	BestBid       *float64   `gorm:"column:best_bid"`
	BestAsk       *float64   `gorm:"column:best_ask"`
	Trades        int64      `gorm:"column:trades"`
	Quantity      int64      `gorm:"column:quantity"`
	Volume        *float64   `gorm:"column:volume"`
	StrikePrice   *float64   `gorm:"column:strike_price"`
	ExpiryDate    *time.Time `gorm:"column:expiry_date"`
	QuoteFactor   int64      `gorm:"column:quote_factor"`
	ISIN          *string    `gorm:"column:isin"`
	DistNumber    *string    `gorm:"column:dist_number"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
}

func (Quote) TableName() string { return "quotes" }

func FromCotahist(r cotahist.Record) Quote {
	return Quote{
		Date:          r.Date,
		Ticker:        r.Ticker,
		BDICode:       r.BDICode,
		MarketType:    r.MarketType,
		ShortName:     nilIfZero(r.ShortName),
		Specification: nilIfZero(r.Specification),
		Currency:      nilIfZero(r.Currency),
		OpenPrice:     nilIfZero(r.OpenPrice),
		HighPrice:     nilIfZero(r.HighPrice),
		LowPrice:      nilIfZero(r.LowPrice),
		AveragePrice:  nilIfZero(r.AveragePrice),
		ClosePrice:    nilIfZero(r.ClosePrice),
		BestBid:       nilIfZero(r.BestBid),
		BestAsk:       nilIfZero(r.BestAsk),
		Trades:        r.Trades,
		Quantity:      r.Quantity,
		Volume:        nilIfZero(r.Volume),
		StrikePrice:   nilIfZero(r.StrikePrice),
		ExpiryDate:    r.ExpiryDate,
		QuoteFactor:   r.QuoteFactor,
		ISIN:          nilIfZero(r.ISIN),
		DistNumber:    nilIfZero(r.DistNumber),
	}
}

// nilIfZero returns nil for a zero value so empty COTAHIST fields persist as NULL.
func nilIfZero[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}
