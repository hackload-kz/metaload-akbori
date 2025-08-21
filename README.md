# Biletter Service - Go Implementation

Высокопроизводительная реализация сервиса бронирования билетов на Go, заменяющая Java версию для улучшения производительности и снижения потребления ресурсов.

## Архитектура

- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL с raw SQL запросами
- **Миграции**: golang-migrate
- **Логирование**: Zap (structured logging)
- **Конфигурация**: Viper
- **Валидация**: go-playground/validator

## Структура проекта

```
biletter-service/
├── cmd/server/           # Точка входа приложения
├── internal/
│   ├── config/          # Конфигурация
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Модели данных и DTOs
│   ├── repository/      # Слой доступа к данным
│   └── services/        # Бизнес-логика
├── pkg/
│   ├── database/        # Подключение к БД
│   └── logger/          # Настройка логгера
├── migrations/          # SQL миграции
└── config.yaml         # Конфигурация
```

## API Endpoints

### События
- `GET /api/events` - Список событий с фильтрацией и пагинацией

### Места
- `GET /api/seats/:event_id` - Список мест для события
- `POST /api/seats/select` - Выбрать место
- `POST /api/seats/release` - Освободить место

### Бронирования
- `POST /api/bookings` - Создать бронирование
- `GET /api/bookings/user/:user_id` - Бронирования пользователя
- `POST /api/bookings/cancel` - Отменить бронирование

### Платежи
- `POST /api/payments/initiate` - Инициировать платеж

### Мониторинг
- `GET /health` - Health check

## Быстрый старт

### Локальная разработка

1. Установить зависимости:
```bash
go mod download
```

2. Запустить базу данных:
```bash
docker-compose up postgres -d
```

3. Выполнить миграции:
```bash
migrate -path migrations -database "postgres://biletter_user:biletter_pass@localhost:5432/biletter_db?sslmode=disable" up
```

4. Запустить приложение:
```bash
make run
# или
go run cmd/server/main.go
```

### Docker

1. Собрать и запустить Go версию:
```bash
docker-compose --profile app up --build
```

### Доступные команды

```bash
make build      # Собрать бинарник
make run        # Собрать и запустить
make test       # Запустить тесты
make dev        # Запуск в режиме разработки
make fmt        # Форматирование кода
make lint       # Линтинг (требует golangci-lint)
make docker-build  # Собрать Docker образ
```

## Конфигурация

Приложение использует файл `config.yaml` и переменные окружения:

```yaml
port: "8081"
log_level: "info"
database:
  host: "localhost"
  port: 5432
  user: "biletter_user"  # или DB_USERNAME
  password: "biletter_pass"  # или DB_PASSWORD
```

## Производительность

Преимущества Go версии по сравнению с Java:

- **Память**: ~50MB vs ~2GB (Java)
- **Время старта**: ~1 секунда vs ~30 секунд
- **CPU**: Меньшее потребление благодаря отсутствию GC пауз
- **Размер образа**: ~20MB vs ~200MB

## Миграции

Применить миграции:
```bash
migrate -path migrations -database $DATABASE_URL up
```

Откатить миграции:
```bash
migrate -path migrations -database $DATABASE_URL down
```

## Мониторинг

- Health check: `GET /health`
- Логи в JSON формате для удобной обработки
- Structured logging с контекстом запросов

## Разработка

1. Форк репозитория
2. Создать feature branch
3. Внести изменения
4. Запустить тесты и линтинг
5. Создать Pull Request