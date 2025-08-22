package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type BookingStatus string

const (
	BookingStatusPending        BookingStatus = "PENDING"
	BookingStatusPaymentPending BookingStatus = "PAYMENT_PENDING"
	BookingStatusConfirmed      BookingStatus = "CONFIRMED"
	BookingStatusCancelled      BookingStatus = "CANCELLED"
)

type Booking struct {
	ID          int64           `json:"id" db:"id"`
	EventID     int64           `json:"event_id" db:"event_id"`
	UserID      int             `json:"user_id" db:"user_id"`
	Status      BookingStatus   `json:"status" db:"status"`
	TotalAmount decimal.Decimal `json:"total_amount" db:"total_amount"`
	PaymentID   *string         `json:"payment_id" db:"payment_id"`
	OrderID     *string         `json:"order_id" db:"order_id"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}
