package quote

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/clause"

	"github.com/allison-piovani/oracle_stocks/internal/database"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// UpsertBatch inserts quotes ignoring rows that collide on (date, ticker, bdi_code).
func (r *Repository) UpsertBatch(ctx context.Context, quotes []Quote) error {
	if len(quotes) == 0 {
		return nil
	}
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "date"}, {Name: "ticker"}, {Name: "bdi_code"}},
			DoNothing: true,
		}).
		CreateInBatches(quotes, 1000).Error
	if err != nil {
		return fmt.Errorf("upsert quotes: %w", err)
	}
	return nil
}

func (r *Repository) ListByTicker(ctx context.Context, ticker string, from, to time.Time) ([]Quote, error) {
	var quotes []Quote
	err := r.db.WithContext(ctx).
		Where("ticker = ? AND date BETWEEN ? AND ?", ticker, from, to).
		Order("date DESC").
		Find(&quotes).Error
	if err != nil {
		return nil, fmt.Errorf("list quotes for %s: %w", ticker, err)
	}
	return quotes, nil
}

func (r *Repository) LatestDate(ctx context.Context) (time.Time, error) {
	var latest time.Time
	err := r.db.WithContext(ctx).
		Model(&Quote{}).
		Select("COALESCE(MAX(date), '0001-01-01'::date)").
		Row().Scan(&latest)
	if err != nil {
		return time.Time{}, fmt.Errorf("latest quote date: %w", err)
	}
	return latest, nil
}
