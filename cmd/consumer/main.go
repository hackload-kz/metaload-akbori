package main

import (
	"biletter-service/internal/config"
	"biletter-service/internal/repository"
	"biletter-service/internal/services"
	"biletter-service/pkg/database"
	"biletter-service/pkg/logger"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

	zapLogger := logger.New(cfg.LogLevel)
	defer zapLogger.Sync()

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	repos := repository.New(db)

	// Создаем consumer service
	consumerService, err := services.NewConsumerService(
		cfg.Kafka,
		cfg.Kafka.ConsumerGroup,
		repos,
		zapLogger,
	)
	if err != nil {
		log.Fatal("Failed to create consumer service:", err)
	}

	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	// Запускаем consumer service
	if err := consumerService.Start(ctx); err != nil {
		log.Fatal("Failed to start consumer service:", err)
	}

	log.Println("Consumer started successfully")

	// Ожидаем сигналы завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down consumer...")

	// Отменяем контекст
	cancel()

	// Останавливаем consumer service
	if err := consumerService.Stop(); err != nil {
		log.Printf("Error stopping consumer service: %v", err)
	}

	log.Println("Consumer stopped")
}
