package quote

import (
	"time"

	"github.com/allison-piovani/oracle_stocks/internal/provider/cotahist"
)

type Quote struct {
	ID            uint64     `gorm:"primaryKey"`
	Date          time.Time  `gorm:"type:date;not null;uniqueIndex:uq_quotes_date_ticker_bdi,priority:1;index:idx_quotes_ticker_date,priority:2,sort:desc"`
	Ticker        string     `gorm:"type:varchar(12);not null;uniqueIndex:uq_quotes_date_ticker_bdi,priority:2;index:idx_quotes_ticker_date,priority:1"`
	BDICode       string     `gorm:"type:varchar(2);not null;uniqueIndex:uq_quotes_date_ticker_bdi,priority:3"`
	MarketType    string     `gorm:"type:varchar(3);not null"`
	ShortName     string     `gorm:"type:varchar(12)"`
	Specification string     `gorm:"type:varchar(10)"`
	Currency      string     `gorm:"type:varchar(4)"`
	OpenPrice     float64    `gorm:"type:numeric(18,4)"`
	HighPrice     float64    `gorm:"type:numeric(18,4)"`
	LowPrice      float64    `gorm:"type:numeric(18,4)"`
	AveragePrice  float64    `gorm:"type:numeric(18,4)"`
	ClosePrice    float64    `gorm:"type:numeric(18,4)"`
	BestBid       float64    `gorm:"type:numeric(18,4)"`
	BestAsk       float64    `gorm:"type:numeric(18,4)"`
	Trades        int64      `gorm:"not null;default:0"`
	Quantity      int64      `gorm:"not null;default:0"`
	Volume        float64    `gorm:"type:numeric(20,2)"`
	StrikePrice   float64    `gorm:"type:numeric(18,4)"`
	ExpiryDate    *time.Time `gorm:"type:date"`
	QuoteFactor   int64      `gorm:"not null;default:1"`
	ISIN          string     `gorm:"type:varchar(12);index:idx_quotes_isin"`
	DistNumber    string     `gorm:"type:varchar(3)"`
	CreatedAt     time.Time
}

func (Quote) TableName() string { return "quotes" }

func FromCotahist(r cotahist.Record) Quote {
	return Quote{
		Date:          r.Date,
		Ticker:        r.Ticker,
		BDICode:       r.BDICode,
		MarketType:    r.MarketType,
		ShortName:     r.ShortName,
		Specification: r.Specification,
		Currency:      r.Currency,
		OpenPrice:     r.OpenPrice,
		HighPrice:     r.HighPrice,
		LowPrice:      r.LowPrice,
		AveragePrice:  r.AveragePrice,
		ClosePrice:    r.ClosePrice,
		BestBid:       r.BestBid,
		BestAsk:       r.BestAsk,
		Trades:        r.Trades,
		Quantity:      r.Quantity,
		Volume:        r.Volume,
		StrikePrice:   r.StrikePrice,
		ExpiryDate:    r.ExpiryDate,
		QuoteFactor:   r.QuoteFactor,
		ISIN:          r.ISIN,
		DistNumber:    r.DistNumber,
	}
}
