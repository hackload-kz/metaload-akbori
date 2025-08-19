package models

import (
	"time"
)

type BookingSeat struct {
	ID        int64     `json:"id" db:"id"`
	BookingID int64     `json:"booking_id" db:"booking_id"`
	SeatID    int64     `json:"seat_id" db:"seat_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}