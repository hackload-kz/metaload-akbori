package broker

import (
	"biletter-service/internal/models"
	"context"
)

// Consumer интерфейс для потребления событий из брокера сообщений
type Consumer interface {
	Subscribe(ctx context.Context, topics []string, handler EventHandler) error
	Close() error
}

// EventHandler интерфейс для обработки событий
type EventHandler interface {
	Handle(ctx context.Context, event *models.DomainEvent) error
}

// EventHandlerFunc функциональный тип для обработчиков
type EventHandlerFunc func(ctx context.Context, event *models.DomainEvent) error

func (f EventHandlerFunc) Handle(ctx context.Context, event *models.DomainEvent) error {
	return f(ctx, event)
}
