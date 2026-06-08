package kafka

import (
	"context"
	"os"
	"strings"

	"github.com/prepio/prepio/shared/devsync"
)

// EventPublisher is implemented by Kafka producers and the dev-sync HTTP publisher.
type EventPublisher interface {
	Publish(ctx context.Context, topic, key string, payload any) error
	Close() error
}

// NewEventPublisher returns a Kafka producer or dev-sync HTTP publisher based on DEV_SYNC_EVENTS.
func NewEventPublisher(cfg ProducerConfig) (EventPublisher, error) {
	if devSyncEnabled() {
		return devsync.NewPublisherFromEnv(), nil
	}
	return NewProducer(cfg)
}

func devSyncEnabled() bool {
	return strings.EqualFold(os.Getenv("DEV_SYNC_EVENTS"), "true")
}
