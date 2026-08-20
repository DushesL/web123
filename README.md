# web123

MVP веб-портал для бронирования приёмов (Go + net/http + PostgreSQL).

## Стек
- Go 1.2x
- PostgreSQL 15
- Docker Compose
- golang-migrate

## Локальная настройка
1. Поднять БД: `docker compose up -d`
2. Применить миграции: `migrate -path migrations -database "postgres://postgres:password@localhost:5432/web123_db?sslmode=disable" up`
3. Запустить приложение: `go run cmd/web/main.go`

## Переменные окружения
Скопируй `.env.example` в `.env` и заполни значения.