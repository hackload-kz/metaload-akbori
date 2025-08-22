package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventType определяет тип события
type EventType string

const (
	BookingCreatedEvent   EventType = "booking.created"
	BookingCancelledEvent EventType = "booking.cancelled"
	SeatSelectedEvent     EventType = "seat.selected"
	SeatReleasedEvent     EventType = "seat.released"
)

// DomainEvent базовая структура для всех доменных событий
type DomainEvent struct {
	ID          string    `json:"id"`
	Type        EventType `json:"type"`
	AggregateID string    `json:"aggregate_id"` // ID брони для всех событий
	Version     int       `json:"version"`
	Data        any       `json:"data"`
	Timestamp   time.Time `json:"timestamp"`
}

// BookingCreatedData данные события создания брони
type BookingCreatedData struct {
	BookingID   int64 `json:"booking_id"`
	EventID     int64 `json:"event_id"`
	UserID      int   `json:"user_id"`
	TotalAmount int64 `json:"total_amount"` // в копейках
}

// BookingCancelledData данные события отмены брони
type BookingCancelledData struct {
	BookingID int64  `json:"booking_id"`
	UserID    int    `json:"user_id"`
	Reason    string `json:"reason,omitempty"`
}

// SeatSelectedData данные события выбора места
type SeatSelectedData struct {
	BookingID int64 `json:"booking_id"`
	SeatID    int64 `json:"seat_id"`
	UserID    int   `json:"user_id"`
}

// SeatReleasedData данные события освобождения места
type SeatReleasedData struct {
	BookingID int64 `json:"booking_id"`
	SeatID    int64 `json:"seat_id"`
	UserID    int   `json:"user_id"`
}

// NewDomainEvent создает новое доменное событие
func NewDomainEvent(eventType EventType, aggregateID string, data any) *DomainEvent {
	return &DomainEvent{
		ID:          uuid.New().String(),
		Type:        eventType,
		AggregateID: aggregateID,
		Version:     1,
		Data:        data,
		Timestamp:   time.Now().UTC(),
	}
}

// ToJSON сериализует событие в JSON
func (e *DomainEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
