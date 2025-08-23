package domain_events

import (
	"biletter-service/internal/models"
	"biletter-service/internal/repository"
	"biletter-service/pkg/broker"
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

// Handlers содержит обработчики для различных типов доменных событий
type Handlers struct {
	repos  *repository.Repository
	logger *zap.Logger
}

// NewHandlers создает новый Handlers
func NewHandlers(repos *repository.Repository, logger *zap.Logger) *Handlers {
	return &Handlers{
		repos:  repos,
		logger: logger,
	}
}

// GetMainHandler возвращает основной обработчик событий
func (h *Handlers) GetMainHandler() broker.EventHandler {
	return broker.EventHandlerFunc(func(ctx context.Context, event *models.DomainEvent) error {
		h.logger.Info("Processing domain event",
			zap.String("event_id", event.ID),
			zap.String("event_type", string(event.Type)),
			zap.String("aggregate_id", event.AggregateID))

		switch event.Type {
		case models.BookingCreatedEvent:
			return h.handleBookingCreated(ctx, event)
		case models.BookingCancelledEvent:
			return h.handleBookingCancelled(ctx, event)
		case models.SeatSelectedEvent:
			return h.handleSeatSelected(ctx, event)
		case models.SeatReleasedEvent:
			return h.handleSeatReleased(ctx, event)
		default:
			h.logger.Warn("Unknown event type", zap.String("event_type", string(event.Type)))
			return nil // Игнорируем неизвестные события
		}
	})
}

// handleBookingCreated обрабатывает событие создания брони
func (h *Handlers) handleBookingCreated(ctx context.Context, event *models.DomainEvent) error {
	var data models.BookingCreatedData
	if err := h.unmarshalEventData(event, &data); err != nil {
		return fmt.Errorf("failed to unmarshal BookingCreatedData: %w", err)
	}

	h.logger.Info("Processing booking created event",
		zap.Int64("booking_id", data.BookingID),
		zap.Int64("event_id", data.EventID),
		zap.Int("user_id", data.UserID))

	// Здесь можно добавить логику:
	// - Отправка email уведомлений
	// - Обновление статистики
	// - Логирование аналитики
	// - Обновление кеша

	// Пример: очистка кеша событий при создании брони
	// TODO: Implement cache invalidation logic

	return nil
}

// handleBookingCancelled обрабатывает событие отмены брони
func (h *Handlers) handleBookingCancelled(ctx context.Context, event *models.DomainEvent) error {
	var data models.BookingCancelledData
	if err := h.unmarshalEventData(event, &data); err != nil {
		return fmt.Errorf("failed to unmarshal BookingCancelledData: %w", err)
	}

	h.logger.Info("Processing booking cancelled event",
		zap.Int64("booking_id", data.BookingID),
		zap.Int("user_id", data.UserID),
		zap.String("reason", data.Reason))

	// Здесь можно добавить логику:
	// - Отправка уведомлений о отмене
	// - Возврат средств (если было оплачено)
	// - Обновление статистики
	// - Освобождение ресурсов

	return nil
}

// handleSeatSelected обрабатывает событие выбора места
func (h *Handlers) handleSeatSelected(ctx context.Context, event *models.DomainEvent) error {
	var data models.SeatSelectedData
	if err := h.unmarshalEventData(event, &data); err != nil {
		return fmt.Errorf("failed to unmarshal SeatSelectedData: %w", err)
	}

	h.logger.Info("Processing seat selected event",
		zap.Int64("booking_id", data.BookingID),
		zap.Int64("seat_id", data.SeatID),
		zap.Int("user_id", data.UserID))

	// Здесь можно добавить логику:
	// - Обновление кеша доступных мест
	// - Отправка уведомлений
	// - Обновление статистики по местам
	// - Временная блокировка места с таймаутом

	return nil
}

// handleSeatReleased обрабатывает событие освобождения места
func (h *Handlers) handleSeatReleased(ctx context.Context, event *models.DomainEvent) error {
	var data models.SeatReleasedData
	if err := h.unmarshalEventData(event, &data); err != nil {
		return fmt.Errorf("failed to unmarshal SeatReleasedData: %w", err)
	}

	h.logger.Info("Processing seat released event",
		zap.Int64("booking_id", data.BookingID),
		zap.Int64("seat_id", data.SeatID),
		zap.Int("user_id", data.UserID))

	// Здесь можно добавить логику:
	// - Обновление кеша доступных мест
	// - Уведомление ожидающих пользователей
	// - Обновление статистики
	// - Очистка временных блокировок

	return nil
}

// unmarshalEventData десериализует данные события
func (h *Handlers) unmarshalEventData(event *models.DomainEvent, target interface{}) error {
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	if err := json.Unmarshal(dataBytes, target); err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	return nil
}
