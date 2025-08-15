# Biletter Service

Сервис для управления событиями, бронированиями и местами с интеграцией внешних сервисов.

## 🚀 Быстрый запуск

```bash
# Запуск всего окружения (PostgreSQL + Kafka + Приложение)
./start.sh

# Или пошагово:
docker-compose up -d
mvn spring-boot:run
```

## 📋 Реализованные API эндпоинты

### Events (События)
- `POST /api/events` - Создание события
- `GET /api/events` - Список событий с фильтрацией

### Bookings (Бронирования)
- `POST /api/bookings` - Создание бронирования
- `GET /api/bookings` - Список бронирований
- `PATCH /api/bookings/initiatePayment` - Инициация платежа
- `PATCH /api/bookings/cancel` - Отмена бронирования

### Seats (Места)
- `GET /api/seats` - Список мест с фильтрацией
- `PATCH /api/seats/select` - Выбор места
- `PATCH /api/seats/release` - Освобождение места

### Payments (Платежи)
- `GET /api/payments/success` - Успешный платеж
- `GET /api/payments/fail` - Неудачный платеж
- `POST /api/payments/notifications` - Webhook уведомления

## 🏗️ Архитектура

- **Spring Boot 3.2.0** - Основной фреймворк
- **PostgreSQL** - Основная база данных
- **Kafka** - Очередь сообщений
- **Flyway** - Миграции базы данных
- **Spring Data JPA** - Работа с базой данных
- **Spring WebFlux** - HTTP клиент для внешних API

## 🔌 Внешние интеграции

### Hackload Ticketing Service
- Создание и управление заказами
- Управление местами
- API версия: v1

### Payment Gateway
- Инициация платежей
- Проверка статуса
- Подтверждение/отмена
- HMAC-SHA256 аутентификация

## 🗄️ База данных

### Основные таблицы
- `events` - События
- `seats` - Места
- `bookings` - Бронирования
- `booking_seats` - Связь бронирований и мест
- `payment_notifications` - Уведомления о платежах

### Миграции
- Автоматическое создание схемы при запуске
- Flyway для управления версиями

## 🐳 Docker

```bash
# Запуск инфраструктуры
docker-compose up -d

# Сборка приложения
docker build -t biletter-service .

# Запуск приложения
docker run -p 8081:8081 biletter-service
```

## 📊 Kafka Topics

- `booking-events` - События бронирований
- `payment-events` - События платежей
- `seat-selection-events` - События выбора мест

## 🔧 Конфигурация

Основные настройки в `application.yml`:
- Порт: 8081
- Контекст: `/api`
- База данных: PostgreSQL на localhost:5432
- Kafka: localhost:9092

## 🧪 Тестирование

```bash
# Компиляция
mvn clean compile

# Сборка JAR
mvn clean package

# Запуск тестов
mvn test

# Запуск приложения
mvn spring-boot:run
```

## 📝 Примеры API запросов

### Создание события
```bash
curl -X POST http://localhost:8081/api/events \
  -H "Content-Type: application/json" \
  -d '{"title": "Концерт", "external": false}'
```

### Создание бронирования
```bash
curl -X POST http://localhost:8081/api/bookings \
  -H "Content-Type: application/json" \
  -d '{"eventId": 1}'
```

### Выбор места
```bash
curl -X PATCH http://localhost:8081/api/seats/select \
  -H "Content-Type: application/json" \
  -d '{"bookingId": 1, "seatId": 1}'
```

## 📁 Структура проекта

```
src/main/java/com/metaload/biletter/
├── BiletterApplication.java          # Главный класс приложения
├── config/                           # Конфигурации
│   ├── ExternalServiceConfig.java    # Настройки внешних сервисов
│   └── KafkaConfig.java             # Конфигурация Kafka
├── controller/                       # REST контроллеры
│   ├── EventController.java         # API событий
│   ├── BookingController.java       # API бронирований
│   ├── SeatController.java          # API мест
│   └── PaymentController.java       # API платежей
├── dto/                             # Data Transfer Objects
│   ├── CreateEventRequest.java      # Запрос создания события
│   ├── CreateBookingRequest.java    # Запрос создания бронирования
│   ├── ListSeatsResponseItem.java   # Ответ со списком мест
│   └── ...                          # Другие DTO
├── model/                           # JPA сущности
│   ├── Event.java                   # Событие
│   ├── Seat.java                    # Место
│   ├── Booking.java                 # Бронирование
│   └── BookingSeat.java             # Связь бронирования и места
├── repository/                      # Репозитории
│   ├── EventRepository.java         # Репозиторий событий
│   ├── SeatRepository.java          # Репозиторий мест
│   ├── BookingRepository.java       # Репозиторий бронирований
│   └── BookingSeatRepository.java   # Репозиторий связей
└── service/                         # Бизнес-логика
    ├── EventService.java            # Сервис событий
    ├── BookingService.java          # Сервис бронирований
    ├── SeatService.java             # Сервис мест
    ├── PaymentService.java          # Сервис платежей
    ├── HackloadService.java         # Интеграция с Hackload
    └── PaymentGatewayService.java   # Интеграция с Payment Gateway
```

## 🚀 Развертывание

### Локальная разработка
1. Запустить `docker-compose up -d`
2. Запустить `mvn spring-boot:run`
3. API доступен на `http://localhost:8081/api`

### Продакшн
1. Собрать JAR: `mvn clean package`
2. Запустить: `java -jar target/biletter-service-1.0.0.jar`
3. Настроить переменные окружения для продакшн