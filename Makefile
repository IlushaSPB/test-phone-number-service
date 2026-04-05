.PHONY: run stop build test clean migrate sqlc setup

run:
	docker compose up -d --build

stop:
	docker compose down

build:
	go build -o bin/api ./cmd/api

test:
	@if [ ! -d "internal/db" ]; then \
		echo "Генерация internal/db/ из SQL..."; \
		if command -v sqlc >/dev/null 2>&1; then \
			sqlc generate; \
		else \
			echo ""; \
			echo "❌ sqlc не найден в PATH"; \
			echo ""; \
			echo "Установите sqlc:"; \
			echo "  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"; \
			echo ""; \
			echo "Затем добавьте ~/go/bin в PATH (если ещё не добавлено):"; \
			echo "  echo 'export PATH=\"\$$HOME/go/bin:\$$PATH\"' >> ~/.zshrc"; \
			echo "  source ~/.zshrc"; \
			echo ""; \
			exit 1; \
		fi \
	fi
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
