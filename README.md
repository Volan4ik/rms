# Restaurant Management System (Go + PostgreSQL)

Курсовой проект по предмету «Базы данных»: управление рестораном (зал, заказы, склад).

## Запуск через Docker Compose
1. Установите переменные окружения (минимум `DB_PASSWORD`). По умолчанию используются:
   - `DB_USER=postgres`
   - `DB_PASSWORD=postgres`
   - `DB_NAME=rms`
2. Соберите и запустите:  
   ```sh
   docker-compose up --build
   ```
   PostgreSQL прогружает миграции из `migrations/schema.sql` и тестовые данные `migrations/test_data.sql` автоматически.

## Ручной запуск без Docker
1. Создайте БД PostgreSQL и примените миграции:
   ```sh
   psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f migrations/schema.sql
   psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f migrations/test_data.sql
   ```
2. Экспортируйте переменные окружения (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `HTTP_PORT`).
3. Запустите сервер:
   ```sh
   go run ./cmd/server
   ```

## Swagger / API
- Swagger UI: `http://localhost:8080/swagger/index.html`
- Базовый health-check: `GET /health`
- Основные эндпоинты (JSON):
  - `GET/POST/PUT/DELETE /api/customers`
  - `GET/POST/PUT/DELETE /api/employees`
  - `GET/POST/PUT/DELETE /api/tables`
  - `GET/POST/PUT/DELETE /api/menu-categories`
  - `GET/POST/PUT/DELETE /api/dishes`
  - `GET/POST/PUT/DELETE /api/products`
  - `GET/POST/PUT/DELETE /api/reservations`, `PUT /api/reservations/{id}/status`
  - `GET/POST /api/orders`, `PUT /api/orders/{id}/status`, `GET/POST/DELETE /api/orders/{id}/items`
  - `POST/DELETE /api/payments`
  - Отчёты: `/api/reports/shift-revenue`, `/api/reports/waiters`, `/api/reports/dishes-availability`
  - Батч: `POST /api/batch-import/products` (JSON массив или CSV файл с колонками `name,unit,cost_price,is_available`)

Примеры curl:
```sh
curl -X POST http://localhost:8080/api/customers \
  -H "Content-Type: application/json" \
  -d '{"full_name":"John Doe","phone":"+70001234567","vip_level":1}'

curl -X POST http://localhost:8080/api/batch-import/products \
  -H "Content-Type: application/json" \
  -d '[{"name":"New product","unit":"pcs","cost_price":10.5,"is_available":true}]'
```

## Структура
- `migrations/schema.sql` — DDL, функции, представления, триггеры, аудит.
- `migrations/test_data.sql` — генерация 600+ продуктов, 550 блюд, 800 клиентов, 5000+ заказов.
- `cmd/server` — точка входа.
- `internal/config` — загрузка ENV.
- `internal/db` — подключение PostgreSQL.
- `internal/domain` — модели.
- `internal/repository` — SQL-слой (только подготовленные запросы).
- `internal/http` — роутер и обработчики.
- `api/docs` — заглушка Swagger (можно пересобрать командой `swag init` при наличии swaggo).
- `Dockerfile`, `docker-compose.yml` — контейнеризация.

## Безопасность и ограничения
- Секреты только через переменные окружения, в коде отсутствуют пароли/URI.
- Все SQL-запросы параметризованы (placeholders `$1..$n`) — защита от SQL-инъекций.
- Аудит CRUD-операций для ключевых таблиц хранится в `audit_log`.
