# API Testing Guide

## 🚀 Быстрый старт для тестирования

### 1. Запуск окружения
```bash
# Запуск PostgreSQL и Kafka
docker-compose up -d

# Запуск приложения
mvn spring-boot:run
```

### 2. Проверка работоспособности
```bash
# Health check
curl http://localhost:8081/api/actuator/health

# Список событий (должен быть пустым)
curl http://localhost:8081/api/events
```

## 📝 Тестовые сценарии

### Сценарий 1: Создание события и бронирования

```bash
# 1. Создание события
curl -X POST http://localhost:8081/api/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Концерт в парке",
    "external": false
  }'

# Ожидаемый ответ: {"id": 1}

# 2. Создание бронирования
curl -X POST http://localhost:8081/api/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "eventId": 1
  }'

# Ожидаемый ответ: {"id": 1}

# 3. Проверка списка бронирований
curl http://localhost:8081/api/bookings
```

### Сценарий 2: Работа с местами

```bash
# 1. Получение списка мест (пока пустой)
curl "http://localhost:8081/api/seats?event_id=1"

# 2. Создание мест через базу данных (для тестирования)
# Можно добавить через psql или создать API endpoint

# 3. Выбор места (после создания мест)
curl -X PATCH http://localhost:8081/api/seats/select \
  -H "Content-Type: application/json" \
  -d '{
    "bookingId": 1,
    "seatId": 1
  }'

# 4. Проверка статуса места
curl "http://localhost:8081/api/seats?event_id=1&status=RESERVED"
```

### Сценарий 3: Управление платежами

```bash
# 1. Инициация платежа
curl -X PATCH http://localhost:8081/api/bookings/initiatePayment \
  -H "Content-Type: application/json" \
  -d '{
    "bookingId": 1
  }'

# 2. Симуляция успешного платежа
curl "http://localhost:8081/api/payments/success?orderId=1"

# 3. Проверка статуса бронирования
curl http://localhost:8081/api/bookings
```

### Сценарий 4: Отмена бронирования

```bash
# 1. Отмена бронирования
curl -X PATCH http://localhost:8081/api/bookings/cancel \
  -H "Content-Type: application/json" \
  -d '{
    "bookingId": 1
  }'

# 2. Проверка, что места освобождены
curl "http://localhost:8081/api/seats?event_id=1&status=FREE"
```

## 🔧 Полезные команды для отладки

### Проверка базы данных
```bash
# Подключение к PostgreSQL
docker exec -it metaload-akbori-postgres-1 psql -U biletter_user -d biletter_db

# Просмотр таблиц
\dt

# Просмотр данных
SELECT * FROM events;
SELECT * FROM bookings;
SELECT * FROM seats;
SELECT * FROM booking_seats;
```

### Проверка Kafka
```bash
# Просмотр топиков
docker exec -it metaload-akbori-kafka-1 kafka-topics --bootstrap-server localhost:9092 --list

# Просмотр сообщений в топике
docker exec -it metaload-akbori-kafka-1 kafka-console-consumer --bootstrap-server localhost:9092 --topic booking-events --from-beginning
```

### Логи приложения
```bash
# Просмотр логов Spring Boot
tail -f logs/spring.log

# Или через Maven
mvn spring-boot:run -Dspring-boot.run.arguments="--logging.level.com.metaload.biletter=DEBUG"
```

## 🐛 Частые проблемы и решения

### 1. Ошибка подключения к базе данных
```bash
# Проверьте статус PostgreSQL
docker-compose ps postgres

# Перезапустите контейнер
docker-compose restart postgres
```

### 2. Ошибка подключения к Kafka
```bash
# Проверьте статус Kafka
docker-compose ps kafka

# Перезапустите контейнер
docker-compose restart kafka
```

### 3. Приложение не запускается
```bash
# Проверьте логи
mvn spring-boot:run

# Убедитесь, что порт 8081 свободен
lsof -i :8081
```

## 📊 Мониторинг

### Health Checks
- `http://localhost:8081/api/actuator/health` - Общее состояние
- `http://localhost:8081/api/actuator/health/db` - Состояние БД
- `http://localhost:8081/api/actuator/health/kafka` - Состояние Kafka

### Метрики
- `http://localhost:8081/api/actuator/metrics` - Доступные метрики
- `http://localhost:8081/api/actuator/prometheus` - Prometheus метрики

## 🎯 Автоматизированное тестирование

### Запуск тестов
```bash
# Все тесты
mvn test

# Только unit тесты
mvn test -Dtest=*UnitTest

# Только интеграционные тесты
mvn test -Dtest=*IntegrationTest
```

### Тестирование с реальными сервисами
```bash
# Профиль для интеграции
mvn spring-boot:run -Dspring.profiles.active=integration
```

## 📚 Дополнительные ресурсы

- [Spring Boot Testing Guide](https://spring.io/guides/gs/testing-web/)
- [REST API Testing Best Practices](https://restfulapi.net/testing-rest-api-manually/)
- [Postman Collection](https://www.postman.com/) - для удобного тестирования API
