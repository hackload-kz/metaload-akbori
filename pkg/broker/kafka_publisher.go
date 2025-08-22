package broker

import (
	"biletter-service/internal/config"
	"biletter-service/internal/models"
	"context"
	"fmt"

	"github.com/IBM/sarama"
)

// KafkaPublisher реализация Publisher для Kafka
type KafkaPublisher struct {
	producer sarama.SyncProducer
}

// NewKafkaPublisher создает новый Kafka publisher
func NewKafkaPublisher(cfg config.Kafka) (Publisher, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Partitioner = sarama.NewManualPartitioner

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaPublisher{
		producer: producer,
	}, nil
}

// Publish отправляет событие в Kafka
func (p *KafkaPublisher) Publish(ctx context.Context, topic string, key string, event *models.DomainEvent) error {
	value, err := event.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("event_type"),
				Value: []byte(event.Type),
			},
			{
				Key:   []byte("event_id"),
				Value: []byte(event.ID),
			},
		},
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	// Логируем успешную отправку (можно заменить на proper logging)
	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)

	return nil
}

// Close закрывает producer
func (p *KafkaPublisher) Close() error {
	return p.producer.Close()
}
