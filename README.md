# oracle_stocks

Go REST API for stock market data, backed by PostgreSQL.

## Stack

- **HTTP**: [fasthttp](https://github.com/valyala/fasthttp) + [fasthttp/router](https://github.com/fasthttp/router)
- **ORM**: [GORM](https://gorm.io) with PostgreSQL driver
- **Database**: PostgreSQL 18
- **Build**: Go 1.26, multi-stage Docker image

## Requirements

- Go 1.26+
- Docker & Docker Compose

## Getting started

```bash
cp .env.example .env
make docker-up   # start PostgreSQL
make run         # start the API
```

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

## Make targets

| Target          | Description                          |
|-----------------|--------------------------------------|
| `make run`      | Run the API locally                  |
| `make build`    | Build binary to `bin/oracle_stocks`  |
| `make test`     | Run tests                            |
| `make tidy`     | Tidy Go modules                      |
| `make docker-up`| Start PostgreSQL in Docker           |
| `make docker-down` | Stop containers                   |
| `make docker-logs` | Tail container logs               |
| `make docker-reset`| Stop containers and remove volumes|

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
├── cmd/api/        # main entrypoint
├── internal/
│   ├── config/     # environment config
│   └── database/   # GORM connection
├── Dockerfile
├── docker-compose.yml
└── Makefile
```
