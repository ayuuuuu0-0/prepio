package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prepio/prepio/services/streak/internal/service"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/kafka"
)

// QuestionAnsweredConsumer processes question.answered events.
type QuestionAnsweredConsumer struct {
	streaks *service.StreakService
}

// NewQuestionAnsweredConsumer creates a QuestionAnsweredConsumer.
func NewQuestionAnsweredConsumer(streaks *service.StreakService) *QuestionAnsweredConsumer {
	return &QuestionAnsweredConsumer{streaks: streaks}
}

// Handle decodes and processes a Kafka message.
func (c *QuestionAnsweredConsumer) Handle(ctx context.Context, _, value []byte) error {
	var event events.QuestionAnswered
	if err := json.Unmarshal(value, &event); err != nil {
		return fmt.Errorf("decode question answered: %w", err)
	}
	return c.streaks.ProcessQuestionAnswered(ctx, event)
}

// Run starts the Kafka consumer loop.
func Run(ctx context.Context, consumer *kafka.Consumer, handler *QuestionAnsweredConsumer) error {
	return consumer.Run(ctx, handler.Handle)
}
