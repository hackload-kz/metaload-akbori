package broker

import (
	"biletter-service/internal/config"
	"biletter-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

// KafkaConsumer реализация Consumer для Kafka
type KafkaConsumer struct {
	consumerGroup sarama.ConsumerGroup
	groupID       string
	ready         chan bool
	wg            sync.WaitGroup
}

// NewKafkaConsumer создает новый Kafka consumer
func NewKafkaConsumer(cfg config.Kafka, groupID string) (Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Session.Timeout = 10000
	config.Consumer.Group.Heartbeat.Interval = 3000

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &KafkaConsumer{
		consumerGroup: consumerGroup,
		groupID:       groupID,
		ready:         make(chan bool),
	}, nil
}

// Subscribe подписывается на топики и начинает обработку событий
func (c *KafkaConsumer) Subscribe(ctx context.Context, topics []string, handler EventHandler) error {
	go func() {
		defer c.wg.Done()
		c.wg.Add(1)

		consumer := &consumerGroupHandler{
			handler: handler,
			ready:   c.ready,
		}

		for {
			select {
			case <-ctx.Done():
				log.Println("Terminating consumer")
				return
			default:
				if err := c.consumerGroup.Consume(ctx, topics, consumer); err != nil {
					log.Printf("Error from consumer: %v", err)
					return
				}

				if ctx.Err() != nil {
					return
				}
				c.ready = make(chan bool)
			}
		}
	}()

	<-c.ready
	log.Printf("Kafka consumer up and running for group %s", c.groupID)
	return nil
}

// Close закрывает consumer
func (c *KafkaConsumer) Close() error {
	c.wg.Wait()
	return c.consumerGroup.Close()
}

// consumerGroupHandler реализует sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	handler EventHandler
	ready   chan bool
}

// Setup запускается в начале новой сессии
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

// Cleanup запускается в конце сессии
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim обрабатывает сообщения из партиции
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			// Десериализуем событие
			var event models.DomainEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				session.MarkMessage(message, "")
				continue
			}

			// Обрабатываем событие
			ctx := context.Background()
			if err := h.handler.Handle(ctx, &event); err != nil {
				log.Printf("Failed to handle event %s: %v", event.Type, err)
				// В продакшене здесь может быть retry logic или DLQ
			}

			// Подтверждаем обработку
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
