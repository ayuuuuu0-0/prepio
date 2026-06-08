package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// MessageHandler processes a single Kafka message value.
type MessageHandler func(ctx context.Context, key, value []byte) error

// Consumer reads messages from a topic within a consumer group.
type Consumer struct {
	reader *kafka.Reader
}

// ConsumerConfig holds Kafka consumer settings.
type ConsumerConfig struct {
	Brokers       []string
	Topic         string
	GroupID       string
	MinBytes      int
	MaxBytes      int
	MaxWait       int // milliseconds
	StartOffset   int64
}

// NewConsumer creates a Kafka consumer for the given topic and group.
func NewConsumer(cfg ConsumerConfig) (*Consumer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("kafka brokers are required")
	}
	if len(cfg.Topic) == 0 {
		return nil, fmt.Errorf("topic is required")
	}
	if len(cfg.GroupID) == 0 {
		return nil, fmt.Errorf("group id is required")
	}

	minBytes := cfg.MinBytes
	if minBytes == 0 {
		minBytes = 1
	}

	maxBytes := cfg.MaxBytes
	if maxBytes == 0 {
		maxBytes = 10e6
	}

	startOffset := cfg.StartOffset
	if startOffset == 0 {
		startOffset = kafka.FirstOffset
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Brokers,
		Topic:       cfg.Topic,
		GroupID:     cfg.GroupID,
		MinBytes:    minBytes,
		MaxBytes:    maxBytes,
		StartOffset: startOffset,
	})

	return &Consumer{reader: reader}, nil
}

// Run blocks and invokes handler for each message until ctx is cancelled.
func (c *Consumer) Run(ctx context.Context, handler MessageHandler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("fetch message: %w", err)
		}

		if err := handler(ctx, msg.Key, msg.Value); err != nil {
			return fmt.Errorf("handle message: %w", err)
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("commit message: %w", err)
		}
	}
}

// Close shuts down the consumer.
func (c *Consumer) Close() error {
	if c.reader == nil {
		return nil
	}
	return c.reader.Close()
}
