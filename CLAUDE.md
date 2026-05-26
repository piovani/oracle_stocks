# CLAUDE.md

Project context for AI agents working on this codebase.

## Project

Go REST API serving Brazilian stock market (B3) historical data, backed by PostgreSQL.
Data is sourced from B3's free **COTAHIST** files (fixed-width, history since 1986).

- Module: `github.com/allison-piovani/oracle_stocks`
- Go version: 1.26

## Architecture

```
cmd/
  api/            → HTTP server (fasthttp): wires config, DB, router
  migrate/        → applies SQL migrations (golang-migrate, up only)
  backfill/       → CLI: downloads COTAHIST years into the DB
internal/
  config/         → env vars → typed Config
  database/       → GORM connection wrapper (*DB)
    migrations/   → embedded *.sql + golang-migrate iofs source (embed.go)
  provider/
    cotahist/     → B3 COTAHIST downloader (client.go) + fixed-width parser + Record DTO (dto.go)
  quote/          → Quote GORM model + Repository + FromCotahist mapper
  backfill/       → Service orchestrating cotahist → quote
```

## Data pipeline

1. `cotahist.Client` downloads/parses annual/monthly/daily COTAHIST ZIPs into `cotahist.Record` (DTO). Use `WalkAnnual`/`Walk*` for streaming — annual files have millions of rows.
2. `quote.FromCotahist` maps a `Record` → `quote.Quote`, converting empty/zero fields to `nil` (nullable columns).
3. `backfill.Service` streams records, filters by BDI code, and calls `quote.Repository.UpsertBatch` in batches.
4. `quote.Repository`: `UpsertBatch` (idempotent — `ON CONFLICT (date,ticker,bdi_code) DO NOTHING`), `ListByTicker`, `LatestDate`.

BDI `02` = cash equities (round lots); empty filter = all instruments (options, BDRs, indices).

## Key dependencies

| Package | Role |
|---|---|
| `github.com/valyala/fasthttp` + `github.com/fasthttp/router` | HTTP server + router (not net/http) |
| `gorm.io/gorm` + `gorm.io/driver/postgres` | ORM + PostgreSQL driver (pgx under the hood) |
| `github.com/golang-migrate/migrate/v4` | SQL migrations (iofs source + postgres driver) |
| `github.com/jackc/pgx/v5/stdlib` | `database/sql` driver used by the migrate command |
| `github.com/joho/godotenv` | `.env` file loading |

## HTTP layer

Handlers use `*fasthttp.RequestCtx`, not `http.ResponseWriter`/`*http.Request`.

```go
r := router.New()
r.GET("/path", func(ctx *fasthttp.RequestCtx) {
    ctx.SetContentType("application/json")
    ctx.SetStatusCode(fasthttp.StatusOK)
    ctx.SetBodyString(`{"key":"value"}`)
})
```

Graceful shutdown: `srv.ShutdownWithContext(ctx)`.

## Database & migrations

- Connection lives in `internal/database/database.go` as `*DB` wrapping `*gorm.DB`.
- Schema is defined **only** in migration SQL under `internal/database/migrations/`, embedded into binaries via `embed.FS` (`embed.go`). Files follow golang-migrate's `NNN_name.up.sql` / `.down.sql` convention.
- Migrations are applied by `make migrate-up` (the `cmd/migrate` command). They are **not** auto-applied on container start — the docker-compose initdb mount was removed to avoid clashing with golang-migrate's version tracking.
- First-time setup flow: `make docker-up && make migrate-up`. After a schema change re-run `make migrate-up`.

## Configuration

All config comes from environment variables (`internal/config/config.go`). Defaults in `.env.example`:

```
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=oracle_stocks
DB_SSLMODE=disable
```

## Common commands

```bash
make run                          # go run ./cmd/api
make build                        # build the API binary
make migrate-up                   # apply DB migrations
make backfill FROM=2020 TO=2024   # ingest COTAHIST years (BDI=02 default)
make test                         # go test ./...
make tidy                         # go mod tidy
make docker-up / docker-down      # start / stop PostgreSQL
make docker-reset                 # stop + remove volumes (wipes DB)
```

`backfill` also accepts `BDI=` and the command exposes `-batch`; for all instruments use `go run ./cmd/backfill -from=2024 -bdi=""`.

## Conventions

- Packages are flat under `internal/` — sub-packages only when the domain warrants it (e.g. `provider/cotahist`).
- GORM struct tags carry **only the column name** (`gorm:"column:x"`); types, constraints, defaults, and indexes live exclusively in the migration SQL.
- Nullable DB columns map to **pointer fields** in Go; use the `nilIfZero` helper to turn zero/empty values into `NULL`.
- Functional options for constructors (`cotahist.New`, `backfill.NewService`).
- No global variables; dependencies passed explicitly. Command logic stays thin — orchestration lives in a service (e.g. `internal/backfill`).
- Errors wrapped with `fmt.Errorf("context: %w", err)`.
- Logging via `log/slog`.
- No comments unless the WHY is non-obvious (the COTAHIST column-mapping comments in `cotahist/dto.go` are the exception worth keeping).

## Current API endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Returns `{"status":"ok"}` |

(Quote-serving endpoints are not implemented yet — the repository exists but isn't wired into the router.)
