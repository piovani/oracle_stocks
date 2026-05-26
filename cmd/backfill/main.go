package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/allison-piovani/oracle_stocks/internal/backfill"
	"github.com/allison-piovani/oracle_stocks/internal/config"
	"github.com/allison-piovani/oracle_stocks/internal/database"
	"github.com/allison-piovani/oracle_stocks/internal/provider/cotahist"
	"github.com/allison-piovani/oracle_stocks/internal/quote"
)

func main() {
	from := flag.Int("from", 0, "first year to backfill (required)")
	to := flag.Int("to", time.Now().Year(), "last year to backfill")
	bdi := flag.String("bdi", "02", `BDI code filter ("02" = cash equities, "" = all instruments)`)
	batchSize := flag.Int("batch", 1000, "rows per insert batch")
	flag.Parse()

	if *from == 0 {
		slog.Error("flag -from is required")
		os.Exit(1)
	}
	if *to < *from {
		slog.Error("flag -to must be >= -from", "from", *from, "to", *to)
		os.Exit(1)
	}

	_ = godotenv.Load()
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, err := database.New(cfg.Database.DSN())
	if err != nil {
		slog.Error("connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	svc := backfill.NewService(
		cotahist.New(),
		quote.NewRepository(db),
		backfill.WithBDI(*bdi),
		backfill.WithBatchSize(*batchSize),
	)

	bdiLabel := *bdi
	if bdiLabel == "" {
		bdiLabel = "all"
	}
	slog.Info("backfill started", "from", *from, "to", *to, "bdi", bdiLabel)

	if err := svc.Years(ctx, *from, *to); err != nil {
		slog.Error("backfill", "error", err)
		os.Exit(1)
	}

	slog.Info("backfill complete", "from", *from, "to", *to)
}
