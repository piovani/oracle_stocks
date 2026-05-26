package main

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/allison-piovani/oracle_stocks/internal/config"
	"github.com/allison-piovani/oracle_stocks/internal/database/migrations"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	m, err := newMigrator(cfg.Database.DSN())
	if err != nil {
		slog.Error("setup migrator", "error", err)
		os.Exit(1)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			slog.Error("close migration source", "error", srcErr)
		}
		if dbErr != nil {
			slog.Error("close migration db", "error", dbErr)
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("migrate up", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations applied")
}

func newMigrator(dsn string) (*migrate.Migrate, error) {
	source, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	driver, err := migratepg.WithInstance(db, &migratepg.Config{})
	if err != nil {
		return nil, err
	}

	return migrate.NewWithInstance("iofs", source, "postgres", driver)
}
