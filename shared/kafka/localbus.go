package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// LocalBus is an in-process Kafka replacement for local development.
type LocalBus struct {
	mu       sync.RWMutex
	handlers map[string][]MessageHandler
}

// NewLocalBus creates a LocalBus.
func NewLocalBus() *LocalBus {
	return &LocalBus{handlers: make(map[string][]MessageHandler)}
}

// Subscribe registers a handler for a topic.
func (b *LocalBus) Subscribe(topic string, handler MessageHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], handler)
}

// Publish serializes payload and invokes all topic handlers synchronously.
func (b *LocalBus) Publish(ctx context.Context, topic, key string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	b.mu.RLock()
	handlers := append([]MessageHandler(nil), b.handlers[topic]...)
	b.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(ctx, []byte(key), body); err != nil {
			return fmt.Errorf("handle %s: %w", topic, err)
		}
	}
	return nil
}

// Close is a no-op for LocalBus.
func (b *LocalBus) Close() error { return nil }
