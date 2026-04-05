# Phone Number Service

REST API для управления номерами телефонов с поддержкой импорта, нормализации в E.164 и поиска.

## Запуск

```bash
cp .env.example .env
make run
```

API будет доступен на `http://localhost:8080`

## Примеры использования

### Импорт номеров

```bash
curl -X POST http://localhost:8080/api/numbers/import \
  -H "Content-Type: application/json" \
  -d '{
    "numbers": ["+79161234567", "89261234567", "9031234567"],
    "source": "telegram"
  }'
```

Ответ:
```json
{
  "accepted": 2,
  "skipped": 1,
  "errors": 0
}
```

### Поиск номеров

```bash
# Все номера
curl "http://localhost:8080/api/numbers/search"

# Фильтр по стране
curl "http://localhost:8080/api/numbers/search?country=Russia&limit=10"

# Частичное совпадение номера
curl "http://localhost:8080/api/numbers/search?number=916"

# Комбинация фильтров
curl "http://localhost:8080/api/numbers/search?region=Москва&provider=МТС"
```

## Тесты

```bash
# Установить sqlc (один раз)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

make test
```
