.PHONY: run build test lint docker-up docker-down tidy migrate-up backfill

run:
	go run ./cmd/api

build:
	go build -o bin/oracle_stocks ./cmd/api

migrate-up:
	go run ./cmd/migrate up

backfill:
	go run ./cmd/backfill -from=$(FROM)$(if $(TO), -to=$(TO))$(if $(BDI), -bdi=$(BDI))

test:
	go test ./...

tidy:
	go mod tidy

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-reset:
	docker compose down -v
