package main

import (
	"biletter-service/internal/config"
	"biletter-service/internal/handlers"
	"biletter-service/internal/repository"
	"biletter-service/internal/services"
	"biletter-service/pkg/database"
	"biletter-service/pkg/logger"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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
	services := services.New(repos, cfg, zapLogger)
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
