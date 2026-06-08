package fakes

import (
	"context"
	"encoding/json"
	"sync"
)

// KafkaProducer records published messages for tests.
type KafkaProducer struct {
	mu       sync.Mutex
	Messages []RecordedMessage
}

// RecordedMessage is a captured Kafka publish.
type RecordedMessage struct {
	Topic   string
	Key     string
	Payload []byte
}

// Publish appends the message to the in-memory log.
func (p *KafkaProducer) Publish(ctx context.Context, topic, key string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Messages = append(p.Messages, RecordedMessage{
		Topic:   topic,
		Key:     key,
		Payload: body,
	})
	return nil
}

// Last returns the most recent recorded message.
func (p *KafkaProducer) Last() *RecordedMessage {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.Messages) == 0 {
		return nil
	}
	msg := p.Messages[len(p.Messages)-1]
	return &msg
}
