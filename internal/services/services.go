package services

import (
	"biletter-service/internal/config"
	"biletter-service/internal/repository"
	"biletter-service/pkg/broker"
	"biletter-service/pkg/cache"

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
	Reset          ResetService
	Analytics      AnalyticsService
}

func New(repos *repository.Repository, cacheClient cache.Cache, eventPublisher broker.Publisher, cfg *config.Config, logger *zap.Logger) *Services {
	// Создаем EventProvider сервис
	eventProvider := NewEventProviderService(cfg.ExternalService, logger)

	// Создаем PaymentGateway сервис
	paymentGateway := NewPaymentGatewayService(cfg.Payment, cfg.App.URL, logger)

	// Создаем UserService
	userService := NewUserService(repos.User)

	// Создаем PaymentService с зависимостями
	paymentService := NewPaymentService(repos.Booking, cfg.Payment, paymentGateway, userService, repos.TxManager)

	return &Services{
		Event:          NewEventService(repos.Event, cacheClient),
		Booking:        NewBookingService(repos.Booking, repos.BookingSeat, repos.Seat, repos.Event, repos.TxManager, eventPublisher, cfg.Kafka.Topics.BookingEvents),
		Seat:           NewSeatService(repos.Seat, eventProvider),
		Payment:        paymentService,
		User:           userService,
		EventProvider:  eventProvider,
		PaymentGateway: paymentGateway,
		Reset:          NewResetService(repos.Booking, repos.Seat, repos.TxManager, logger),
		Analytics:      NewAnalyticsService(repos.Seat, repos.Booking, logger),
	}
}
