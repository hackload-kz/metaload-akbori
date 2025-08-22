package main

import (
	"biletter-service/internal/config"
	"biletter-service/internal/handlers"
	"biletter-service/internal/repository"
	"biletter-service/internal/services"
	"biletter-service/pkg/cache"
	"biletter-service/pkg/database"
	"biletter-service/pkg/logger"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	if err := runMigrations(cfg.Database); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Создаем cache клиент
	cacheClient := cache.NewRedisCache(cfg.Redis)

	repos := repository.New(db)
	services := services.New(repos, cacheClient, cfg, zapLogger)
	handlers := handlers.New(services, zapLogger)

	if err := repos.InitializeCache(); err != nil {
		log.Fatal("Failed to initialize cache:", err)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	handlers.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}

func runMigrations(cfg config.Database) error {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	m, err := migrate.New(
		"file://migrations",
		databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
