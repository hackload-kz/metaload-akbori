package broker

import (
	"biletter-service/internal/models"
	"context"
)

// Publisher интерфейс для публикации событий
type Publisher interface {
	Publish(ctx context.Context, topic string, key string, event *models.DomainEvent) error
	Close() error
}
