package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type SeatStatus string

const (
	SeatStatusFree     SeatStatus = "FREE"
	SeatStatusReserved SeatStatus = "RESERVED"
	SeatStatusSold     SeatStatus = "SOLD"
)

type Seat struct {
	ID         int64           `json:"id" db:"id"`
	EventID    int64           `json:"event_id" db:"event_id"`
	RowNumber  int             `json:"row_number" db:"row_number"`
	SeatNumber int             `json:"seat_number" db:"seat_number"`
	PlaceId    string          `json:"place_id" db:"place_id"`
	Status     SeatStatus      `json:"status" db:"status"`
	Price      decimal.Decimal `json:"price" db:"price"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
	Version    int64           `json:"version" db:"version"`
}
