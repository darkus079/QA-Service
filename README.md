# QA Service API

API-сервис для вопросов и ответов

## Функциональность

Сервис предоставляет REST API для управления вопросами и ответами:

- Создание, получение и удаление вопросов
- Добавление ответов к вопросам
- Получение всех ответов на конкретный вопрос
- Каскадное удаление ответов при удалении вопроса

## Технологии

- **Go** 1.23
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - база данных
- **Gorilla Mux** - HTTP роутер
- **Goose** - миграции базы данных
- **Docker & Docker Compose** - контейнеризация

## API Endpoints

### Вопросы (Questions)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/api/v1/questions/` | Получить список всех вопросов |
| POST | `/api/v1/questions/` | Создать новый вопрос |
| GET | `/api/v1/questions/{id}` | Получить вопрос с ответами |
| DELETE | `/api/v1/questions/{id}` | Удалить вопрос (и все ответы) |

### Ответы (Answers)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/api/v1/questions/{id}/answers/` | Добавить ответ к вопросу |
| GET | `/api/v1/answers/{id}` | Получить конкретный ответ |
| DELETE | `/api/v1/answers/{id}` | Удалить ответ |

### Системные

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/health` | Проверка работоспособности сервиса |

## Запуск с помощью Docker Compose

### Предварительные требования

- Docker
- Docker Compose

### Быстрый запуск

1. **Клонируйте репозиторий:**
   ```bash
   git clone https://github.com/darkus079/QA-Service.git
   cd qa-service
   ```

2. **Запустите сервисы:**
   ```bash
   docker-compose up --build
   ```

3. **Запустите миграции (в отдельном терминале):**
   ```bash
   docker-compose --profile tools run --rm migrator
   ```

Сервис будет доступен по адресу: `http://localhost:8080`

### Остановка сервисов

```bash
docker-compose down
```

### Остановка с удалением данных

```bash
docker-compose down -v
```

## Локальный запуск без Docker

### Предварительные требования

- Go 1.23+
- PostgreSQL
- Goose (для миграций)

### Установка зависимостей

```bash
go mod tidy
```

### Настройка базы данных

1. **Создайте базу данных PostgreSQL:**
   ```sql
   CREATE DATABASE qa_service;
   ```

2. **Запустите миграции:**
   ```bash
   goose -dir migrations postgres "user=postgres password=password dbname=qa_service sslmode=disable" up
   ```

### Переменные окружения

Создайте файл `.env` или установите переменные окружения:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=your-user
DB_PASSWORD=your-password
DB_NAME=your-db-name
DB_SSLMODE=disable
PORT=8080
```

### Запуск приложения

```bash
go run cmd/server/main.go
```

## Тестирование

### Запуск интеграционных тестов

1. **Создайте тестовую базу данных:**
   ```sql
   CREATE DATABASE qa_service_test;
   ```

2. **Установите переменную окружения для тестов:**
   ```bash
   export TEST_DATABASE_URL="postgres://postgres:password@localhost:5432/qa_service_test?sslmode=disable"
   ```

3. **Запустите тесты:**
   ```bash
   go test ./tests/...
   ```

## Структура проекта

```
.
├── cmd/server/           # Точка входа приложения
├── internal/
│   ├── database/         # Настройка подключения к БД
│   ├── handlers/         # HTTP обработчики
│   ├── models/           # Модели данных
│   ├── repository/       # Репозитории для работы с БД
│   ├── routes/           # Настройка маршрутов
│   └── services/         # Бизнес-логика
├── migrations/           # Миграции базы данных
├── docker/               # Docker файлы
├── scripts/              # Скрипты для миграций
├── tests/                # Интеграционные тесты
└── docker-compose.yml    # Конфигурация Docker Compose
```

## Разработка

### Добавление новых миграций

```bash
goose -dir migrations postgres "user=postgres password=password dbname=qa_service sslmode=disable" create add_new_field sql
```

### Запуск линтера

```bash
go vet ./...
```

### Форматирование кода

```bash
go fmt ./...
```

## Мониторинг

Сервис логирует все HTTP запросы и важные события. Логи выводятся в stdout/stderr.

Проверка здоровья сервиса:
```bash
curl http://localhost:8080/health
```

## Безопасность

- Валидация входных данных
- Защита от SQL-инъекций через GORM
- Graceful shutdown
- CORS headers (можно добавить при необходимости)

## Производительность

- Использование prepared statements через GORM
- Connection pooling в PostgreSQL
- Graceful shutdown для корректного завершения соединений