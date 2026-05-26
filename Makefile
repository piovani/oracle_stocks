.PHONY: run build test lint docker-up docker-down tidy migrate-up

run:
	go run ./cmd/api

build:
	go build -o bin/oracle_stocks ./cmd/api

migrate-up:
	go run ./cmd/migrate up

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
