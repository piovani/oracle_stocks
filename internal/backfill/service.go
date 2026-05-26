package backfill

import (
	"context"
	"log/slog"

	"github.com/allison-piovani/oracle_stocks/internal/provider/cotahist"
	"github.com/allison-piovani/oracle_stocks/internal/quote"
)

type Service struct {
	client    *cotahist.Client
	repo      *quote.Repository
	bdi       string
	batchSize int
}

type Option func(*Service)

// WithBDI filters records to a single BDI code ("" keeps all instruments).
func WithBDI(bdi string) Option {
	return func(s *Service) { s.bdi = bdi }
}

func WithBatchSize(n int) Option {
	return func(s *Service) { s.batchSize = n }
}

func NewService(client *cotahist.Client, repo *quote.Repository, opts ...Option) *Service {
	s := &Service{
		client:    client,
		repo:      repo,
		bdi:       "02",
		batchSize: 1000,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Years backfills an inclusive range of calendar years.
func (s *Service) Years(ctx context.Context, from, to int) error {
	for year := from; year <= to; year++ {
		if err := s.Year(ctx, year); err != nil {
			return err
		}
	}
	return nil
}

// Year downloads a single year's COTAHIST file and persists matching quotes.
func (s *Service) Year(ctx context.Context, year int) error {
	slog.Info("backfilling year", "year", year)

	batch := make([]quote.Quote, 0, s.batchSize)
	var total int

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		if err := s.repo.UpsertBatch(ctx, batch); err != nil {
			return err
		}
		total += len(batch)
		batch = batch[:0]
		return nil
	}

	err := s.client.WalkAnnual(ctx, year, func(r cotahist.Record) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if s.bdi != "" && r.BDICode != s.bdi {
			return nil
		}
		batch = append(batch, quote.FromCotahist(r))
		if len(batch) >= s.batchSize {
			return flush()
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := flush(); err != nil {
		return err
	}

	slog.Info("year done", "year", year, "rows", total)
	return nil
}
