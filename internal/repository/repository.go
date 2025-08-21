package repository

import (
	"database/sql"
)

type Repository struct {
	Event       EventRepository
	Seat        SeatRepository
	Booking     BookingRepository
	BookingSeat BookingSeatRepository
	User        UserRepository
	TxManager   *TransactionManager
}

func New(db *sql.DB) *Repository {
	return &Repository{
		Event:       NewEventRepository(db),
		Seat:        NewSeatRepository(db),
		Booking:     NewBookingRepository(db),
		BookingSeat: NewBookingSeatRepository(db),
		User:        NewUserRepository(db),
		TxManager:   NewTransactionManager(db),
	}
}

// InitializeCache предзагружает кэши при старте приложения
func (r *Repository) InitializeCache() error {
	return r.User.PreloadCache()
}
