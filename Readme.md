# Order Service WB

Микросервис для обработки заказов с интеграцией Kafka и PostgreSQL. Поддерживает приём заказов через Kafka, хранение в PostgreSQL, кэширование и REST-доступ к данным.

## ⚙️ Технологии
- Go 1.24
- Kafka (franz-go)
- PostgreSQL (sqlx)
- Goose (миграции)
- Docker / Docker Compose
- Gin (REST API)
- go-playground/validator (валидация)
- golangci-lint (линтинг)

## 🔧 Функциональность
- ✅ Приём заказов через Kafka
- ✅ Валидация данных при получении
- ✅ Хранение заказов, доставок, оплат, товаров в PostgreSQL
- ✅ Кэширование заказов в памяти
- ✅ Восстановление кэша при старте из базы данных
- ✅ API: получение заказа по `order_uid`

## 🏑 Запуск через Docker
```bash
docker compose up -d
```

## ⚖️ Миграции
Перед запуском нужно применить миграции:
```bash
make goose-install
make migrate-up
```

## ⚡️ Makefile-команды

```bash
# Установка goose
make goose-install

# Применить все миграции
make migrate-up

# Откатить миграцию
make migrate-down

# Статус миграций
make migrate-status

# Создать новую миграцию
make migrate-create name=create_orders_table

# Форсировать версию
make migrate-force version=1

# Запуск приложения
make run

# Генератор заказов в Kafka
make generator

# Запуск линтера
make lint

# Юнит-тесты
make test


# Запуск docker-compose
make docker-up

# Остановка docker-compose
make docker-down
```

## 🚀 Запуск Go-приложения вручную
```bash
go run ./cmd/app
```

## 🔎 Пример API-запроса
```bash
curl http://localhost:8081/order/b563feb7b2b84b6test
```