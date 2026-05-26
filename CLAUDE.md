# CLAUDE.md

Project context for AI agents working on this codebase.

## Project

Go REST API for stock market data backed by PostgreSQL.

- Module: `github.com/allison-piovani/oracle_stocks`
- Go version: 1.26
- Entry point: `cmd/api/main.go`

## Architecture

```
cmd/api/        → main: wires config, DB, router, server
internal/
  config/       → loads env vars into typed Config struct
  database/     → GORM wrapper with connect + ping + close
```

No service/repository layers yet — the project is in early scaffolding stage.

## Key dependencies

| Package | Role |
|---|---|
| `github.com/valyala/fasthttp` | HTTP server |
| `github.com/fasthttp/router` | URL router (replaces net/http mux) |
| `gorm.io/gorm` + `gorm.io/driver/postgres` | ORM + PostgreSQL driver |
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

## Database

- ORM: GORM
- Driver: `pgx` (via `gorm.io/driver/postgres`)
- Connection lives in `internal/database/database.go` as `*DB` wrapping `*gorm.DB`
- Migrations directory (used by docker-compose init): `internal/database/migrations/`

## Configuration

All config comes from environment variables. See `internal/config/config.go`.
Defaults in `.env.example`:

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
make run          # go run ./cmd/api
make build        # go build -o bin/oracle_stocks ./cmd/api
make test         # go test ./...
make tidy         # go mod tidy
make docker-up    # start PostgreSQL container
make docker-down  # stop containers
make docker-reset # stop + remove volumes (wipes DB)
```

## Docker

- Multi-stage build: `golang:1.26-alpine` builder → `alpine:3.22` runtime
- `docker-compose.yml` runs PostgreSQL 18 with health check and auto-runs SQL files from `internal/database/migrations/`

## Conventions

- Packages are flat under `internal/` — no sub-packages unless the domain clearly warrants it.
- No global variables; dependencies are passed explicitly.
- Errors are wrapped with `fmt.Errorf("context: %w", err)`.
- Logging via `log/slog` (structured, standard library).
- No comments unless the WHY is non-obvious.

## Current API endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Returns `{"status":"ok"}` |
