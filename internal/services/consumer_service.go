package services

import (
	"biletter-service/internal/config"
	"biletter-service/internal/domain_events"
	"biletter-service/internal/repository"
	"biletter-service/pkg/broker"
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// ConsumerService управляет обработкой событий из брокера
type ConsumerService struct {
	consumer      broker.Consumer
	eventHandlers *domain_events.Handlers
	logger        *zap.Logger
	topics        []string
	wg            sync.WaitGroup
	cancelFunc    context.CancelFunc
}

// NewConsumerService создает новый ConsumerService
func NewConsumerService(cfg config.Kafka, groupID string, repos *repository.Repository, logger *zap.Logger) (*ConsumerService, error) {
	// Создаем consumer
	consumer, err := broker.NewKafkaConsumer(cfg, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	// Создаем обработчики событий
	eventHandlers := domain_events.NewHandlers(repos, logger)

	// Определяем топики для подписки
	topics := []string{cfg.Topics.BookingEvents}

	return &ConsumerService{
		consumer:      consumer,
		eventHandlers: eventHandlers,
		logger:        logger,
		topics:        topics,
	}, nil
}

// Start запускает обработку событий
func (s *ConsumerService) Start(ctx context.Context) error {
	// Создаем отменяемый контекст
	ctx, cancel := context.WithCancel(ctx)
	s.cancelFunc = cancel

	s.logger.Info("Starting consumer service",
		zap.Strings("topics", s.topics))

	// Запускаем consumer в отдельной горутине
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		if err := s.consumer.Subscribe(ctx, s.topics, s.eventHandlers.GetMainHandler()); err != nil {
			s.logger.Error("Consumer subscription failed", zap.Error(err))
		}
	}()

	return nil
}

// Stop останавливает обработку событий
func (s *ConsumerService) Stop() error {
	s.logger.Info("Stopping consumer service")

	// Отменяем контекст
	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	// Ждем завершения всех горутин
	s.wg.Wait()

	// Закрываем consumer
	if err := s.consumer.Close(); err != nil {
		s.logger.Error("Failed to close consumer", zap.Error(err))
		return err
	}

	s.logger.Info("Consumer service stopped")
	return nil
}

// Wait ожидает завершения работы consumer'а
func (s *ConsumerService) Wait() {
	s.wg.Wait()
}
