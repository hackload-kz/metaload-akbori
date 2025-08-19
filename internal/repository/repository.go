package repository

import (
	"database/sql"
)

type Repository struct {
	Event   EventRepository
	Booking BookingRepository
	Seat    SeatRepository
	User    UserRepository
}

func New(db *sql.DB) *Repository {
	return &Repository{
		Event:   NewEventRepository(db),
		Booking: NewBookingRepository(db),
		Seat:    NewSeatRepository(db),
		User:    NewUserRepository(db),
	}
}