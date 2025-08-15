## Что у нас есть

✅ **Spring Boot 3.2** приложение с полной архитектурой  
✅ **PostgreSQL** база данных с миграциями  
✅ **Kafka** для асинхронной обработки  
✅ **Интеграция с Hackload** Ticketing Service  
✅ **Интеграция с Payment Gateway**  
✅ **Docker Compose** для локальной разработки  

## 🎯 Запуск за 3 минуты

### 1. Клонируйте и перейдите в проект
```bash
cd metaload-akbori
```

### 2. Запустите всё одной командой
```bash
./start.sh
```

### 3. Готово! 🎉
- **API**: http://localhost:8081
- **Kafka UI**: http://localhost:8080
- **База данных**: localhost:5432

## 📋 Что происходит при запуске

1. **Проверка зависимостей** (Docker, Maven, Java 17+)
2. **Запуск инфраструктуры** (PostgreSQL + Kafka)
3. **Сборка проекта** (Maven compile)
4. **Запуск приложения** (Spring Boot)

## 🔧 Ручной запуск

### Только инфраструктура
```bash
docker-compose up -d
```

### Только приложение
```bash
mvn spring-boot:run
```

## 📱 Тестирование API

### Создать событие
```bash
curl -X POST http://localhost:8081/api/events \
  -H "Content-Type: application/json" \
  -d '{"title": "Концерт в парке", "external": false}'
```

### Получить список событий
```bash
curl http://localhost:8081/api/events
```

## 🐛 Отладка

### Логи приложения
```bash
docker-compose logs -f
```

### Проверка здоровья
```bash
curl http://localhost:8081/actuator/health
```

### Kafka топики
```bash
docker exec -it biletter-kafka kafka-topics --bootstrap-server localhost:9092 --list
```

## 🆘 Проблемы

### Docker не запущен
```bash
# macOS
open -a Docker

# Linux
sudo systemctl start docker
```

### Порт занят
```bash
# Остановить всё
docker-compose down

# Запустить заново
./start.sh
```

### База не подключается
```bash
# Проверить статус
docker-compose ps postgres

# Посмотреть логи
docker-compose logs postgres
```

## 📚 Документация

- **Полный README**: [README.md](README.md)
- **API спецификация**: [Biletter-api.json](Biletter-api.json)
- **Hackload API**: [hackload ticketing service provider](hackload%20ticketing%20service%20provider)
- **Payment Gateway**: [payment gateway documentation](payment%20gateway%20documentation)

## 🎯 Следующие шаги

1. **Настройте переменные окружения** в `env.example`
2. **Добавьте свои эндпоинты** в контроллеры
3. **Создайте бизнес-логику** в сервисах
4. **Добавьте тесты** для вашего кода
5. **Настройте CI/CD** для автоматического деплоя