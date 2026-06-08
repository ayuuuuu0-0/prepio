package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Producer publishes JSON-encoded messages to Kafka topics.
type Producer struct {
	writer *kafka.Writer
}

// ProducerConfig holds Kafka producer settings.
type ProducerConfig struct {
	Brokers []string
}

// NewProducer creates a Kafka producer connected to the given brokers.
func NewProducer(cfg ProducerConfig) (*Producer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("kafka brokers are required")
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}, nil
}

// Publish serializes payload as JSON and writes it to the topic.
func (p *Producer) Publish(ctx context.Context, topic string, key string, payload any) error {
	if len(topic) == 0 {
		return fmt.Errorf("topic is required")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: body,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	return nil
}

// Close shuts down the producer.
func (p *Producer) Close() error {
	if p.writer == nil {
		return nil
	}
	return p.writer.Close()
}
