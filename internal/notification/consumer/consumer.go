package consumer

import (
	"context"

	"github.com/prepio/prepio/shared/kafka"
)

// Run starts the Kafka consumer loop with the given handler.
func Run(ctx context.Context, c *kafka.Consumer, handler func(context.Context, []byte) error) error {
	return c.Run(ctx, func(ctx context.Context, _, value []byte) error {
		return handler(ctx, value)
	})
}
