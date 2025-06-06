# Person API
REST API для управления данными людей с использованием PostgreSQL и внешних API.

## Требования
- Go 1.24.1+
- Docker
- PostgreSQL

## Установка
1. Склонируй репозиторий:
    ```bash
    git clone https://github.com/Ekmor73/person-api.git
    cd person-api
    ```
2. Установи зависимости:
    ```bash
    go mod tidy
    ```
3. Настрой `.env`:
    ```bash
    cp .env.example .env
    nano .env
    ```
    Укажи:
    ```plaintext
    DB_HOST=localhost
    DB_PORT=5434
    DB_USER=postgres
    DB_PASSWORD=123
    DB_NAME=person_api
    PORT=8080
    ```
4. Запусти PostgreSQL:
    ```bash
    docker run -d -p 5434:5432 --name postgres -e POSTGRES_PASSWORD=123 postgres:16.4
    psql -U postgres -h localhost -p 5434 -c "CREATE DATABASE person_api;" -W
    ```
5. Запусти приложение:
    ```bash
    go run main.go
    ```
6. Или используй Docker Compose:
    ```bash
    docker-compose up -d
 
 ## Эндпоинты
- `POST /people` — Создать человека
- `GET /people` — Получить список людей (с фильтрами)
- `GET /people/:id` — Получить человека по ID
- `PUT /people/:id` — Обновить человека
- `DELETE /people/:id` — Удалить человека

## Swagger
- Доступен по: `http://localhost:8080/swagger/index.html`

## Структура проекта
- `main.go` — Точка входа, настройка сервера
- `database/` — Подключение к базе и миграции
- `models/` — Модели данных
- `handlers/` — Обработчики HTTP-запросов
- `migrations/` — SQL-миграции   ```
