package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type CreateBookingRequest struct {
	EventID int64 `json:"event_id" validate:"required"`
}

type CreateBookingResponse struct {
	ID int64 `json:"id"`
}

type ListEventsResponseItem struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Type          string    `json:"type"`
	DatetimeStart time.Time `json:"datetime_start"`
	Provider      string    `json:"provider"`
}

type ListSeatsResponseItem struct {
	ID         int64           `json:"id"`
	RowNumber  int             `json:"row_number"`
	SeatNumber int             `json:"seat_number"`
	Status     SeatStatus      `json:"status"`
	Price      decimal.Decimal `json:"price"`
}

type ListBookingsResponseItem struct {
	ID          int64                          `json:"id"`
	EventTitle  string                         `json:"event_title"`
	Status      BookingStatus                  `json:"status"`
	TotalAmount decimal.Decimal                `json:"total_amount"`
	PaymentID   *string                        `json:"payment_id"`
	OrderID     *string                        `json:"order_id"`
	CreatedAt   time.Time                      `json:"created_at"`
	Seats       []ListBookingsResponseItemSeat `json:"seats"`
}

type ListBookingsResponseItemSeat struct {
	SeatID     int64 `json:"seat_id"`
	RowNumber  int   `json:"row_number"`
	SeatNumber int   `json:"seat_number"`
}

type SelectSeatRequest struct {
	BookingID int64 `json:"booking_id" validate:"required"`
	SeatID    int64 `json:"seat_id" validate:"required"`
}

type ReleaseSeatRequest struct {
	SeatID int64 `json:"seat_id" validate:"required"`
}

type CancelBookingRequest struct {
	BookingID int64 `json:"booking_id" validate:"required"`
}

type InitiatePaymentRequest struct {
	BookingID int64 `json:"booking_id" validate:"required"`
}
