# Phone Number Service

REST API для управления номерами телефонов с поддержкой импорта, нормализации в E.164 и поиска.

## Быстрый старт

```bash
# Скопировать .env.example → .env (при необходимости изменить порты/пароли)
cp .env.example .env

# Поднять Postgres, прогнать миграции и запустить API
docker compose up -d --build

# Проверить health-check
curl http://localhost:8080/health
```

## Структура

- `cmd/api` — HTTP-сервер
- `db/migrations` — миграции goose
- `sql/queries` — SQL-запросы для sqlc
- `internal/db` — сгенерированный код sqlc (не коммитится)

## Разработка

```bash
# Установить sqlc (если нет)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Сгенерировать Go-код из SQL
sqlc generate

# Запустить локально (Postgres должен быть доступен на порту из .env)
go run ./cmd/api
```

## Миграции

Миграции накатываются автоматически при `docker compose up` через сервис `migrate`.
