package services

import (
	"biletter-service/internal/config"
	"biletter-service/internal/repository"
)

type Services struct {
	Event   EventService
	Booking BookingService
	Seat    SeatService
	Payment PaymentService
	User    UserService
}

func New(repos *repository.Repository, cfg *config.Config) *Services {
	return &Services{
		Event:   NewEventService(repos.Event),
		Booking: NewBookingService(repos.Booking, repos.Seat, repos.Event),
		Seat:    NewSeatService(repos.Seat),
		Payment: NewPaymentService(repos.Booking, cfg.Payment),
		User:    NewUserService(repos.User),
	}
}
