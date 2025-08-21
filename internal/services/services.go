package services

import (
	"biletter-service/internal/config"
	"biletter-service/internal/repository"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Services struct {
	Event          EventService
	Booking        BookingService
	Seat           SeatService
	Payment        PaymentService
	User           UserService
	EventProvider  EventProviderService
	PaymentGateway PaymentGatewayService
}

func New(repos *repository.Repository, cfg *config.Config, logger *zap.Logger) *Services {
	// Создаем Redis клиент
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Создаем EventProvider сервис
	eventProvider := NewEventProviderService(cfg.ExternalService, logger)

	// Создаем PaymentGateway сервис
	paymentGateway := NewPaymentGatewayService(cfg.Payment, cfg.App.URL, logger)

	// Создаем UserService
	userService := NewUserService(repos.User)

	// Создаем PaymentService с зависимостями
	paymentService := NewPaymentService(repos.Booking, cfg.Payment, paymentGateway, userService, repos.TxManager)

	return &Services{
		Event:          NewEventService(repos.Event, redisClient),
		Booking:        NewBookingService(repos.Booking, repos.Seat, repos.Event, repos.TxManager),
		Seat:           NewSeatService(repos.Seat, eventProvider),
		Payment:        paymentService,
		User:           userService,
		EventProvider:  eventProvider,
		PaymentGateway: paymentGateway,
	}
}
