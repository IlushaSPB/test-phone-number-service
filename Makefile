.PHONY: run stop build test clean migrate sqlc

run:
	docker compose up -d --build

stop:
	docker compose down

build:
	go build -o bin/api ./cmd/api

test:
	go test -v ./...

test-cover:
	go test -cover ./internal/service/

clean:
	rm -rf bin/ coverage.out

migrate:
	docker compose up migrate

sqlc:
	sqlc generate

logs:
	docker compose logs -f api

.DEFAULT_GOAL := run
