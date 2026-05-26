# oracle_stocks

Go REST API serving Brazilian stock market (B3) historical data, backed by PostgreSQL.
Quote history is sourced from B3's free [COTAHIST](https://www.b3.com.br/pt_br/market-data-e-indices/servicos-de-dados/market-data/historico/mercado-a-vista/series-historicas/) files (fixed-width, history since 1986).

## Stack

- **HTTP**: [fasthttp](https://github.com/valyala/fasthttp) + [fasthttp/router](https://github.com/fasthttp/router)
- **ORM**: [GORM](https://gorm.io) with PostgreSQL driver
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate) (SQL embedded in the binary)
- **Database**: PostgreSQL 18
- **Build**: Go 1.26, multi-stage Docker image

## Requirements

- Go 1.26+
- Docker & Docker Compose

## Getting started

```bash
cp .env.example .env
make docker-up                    # start PostgreSQL
make migrate-up                   # create the schema
make backfill FROM=2024 TO=2024   # ingest a year of quotes
make run                          # start the API
```

> Migrations are **not** auto-applied on container start — always run `make migrate-up` after `make docker-up`.

## Environment variables

| Variable      | Default        | Description            |
|---------------|----------------|------------------------|
| `SERVER_PORT` | `8080`         | HTTP listen port       |
| `DB_HOST`     | `localhost`    | PostgreSQL host        |
| `DB_PORT`     | `5432`         | PostgreSQL port        |
| `DB_USER`     | `postgres`     | PostgreSQL user        |
| `DB_PASSWORD` | `postgres`     | PostgreSQL password    |
| `DB_NAME`     | `oracle_stocks`| Database name          |
| `DB_SSLMODE`  | `disable`      | SSL mode               |

## Backfilling history

The `backfill` command downloads and parses COTAHIST files into the `quotes` table.

```bash
make backfill FROM=2020 TO=2024          # range of years (BDI 02 = cash equities)
make backfill FROM=2010                  # from 2010 to the current year
go run ./cmd/backfill -from=2024 -bdi=""  # all instruments (options, BDRs, indices)
```

Command flags: `-from` (required), `-to` (default: current year), `-bdi` (default `02`, `""` = all), `-batch` (default 1000).
Re-running is idempotent — existing rows are skipped via `ON CONFLICT DO NOTHING`.

## Make targets

| Target              | Description                              |
|---------------------|------------------------------------------|
| `make run`          | Run the API locally                      |
| `make build`        | Build binary to `bin/oracle_stocks`      |
| `make migrate-up`   | Apply database migrations                |
| `make backfill`     | Ingest COTAHIST years (`FROM=`, `TO=`)   |
| `make test`         | Run tests                                |
| `make tidy`         | Tidy Go modules                          |
| `make docker-up`    | Start PostgreSQL in Docker               |
| `make docker-down`  | Stop containers                          |
| `make docker-logs`  | Tail container logs                      |
| `make docker-reset` | Stop containers and remove volumes       |

## API

| Method | Path      | Description  |
|--------|-----------|--------------|
| GET    | `/health` | Health check |

## Docker

Build and run the API image:

```bash
docker build -t oracle_stocks .
docker run --env-file .env -p 8080:8080 oracle_stocks
```

## Project structure

```
.
├── cmd/
│   ├── api/            # HTTP server entrypoint
│   ├── migrate/        # migration runner
│   └── backfill/       # COTAHIST ingestion CLI
├── internal/
│   ├── config/         # environment config
│   ├── database/       # GORM connection
│   │   └── migrations/ # embedded SQL migrations
│   ├── provider/
│   │   └── cotahist/   # B3 COTAHIST downloader + parser
│   ├── quote/          # Quote model + repository
│   └── backfill/       # ingestion service
├── Dockerfile
├── docker-compose.yml
└── Makefile
```
