package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/kafka"
)

// Handler routes Kafka messages to progress logic.
type Handler struct {
	progress  *service.ProgressService
	readiness *service.ReadinessService
}

// NewHandler creates a Handler.
func NewHandler(progress *service.ProgressService, readiness *service.ReadinessService) *Handler {
	return &Handler{progress: progress, readiness: readiness}
}

// HandleQuestionAnswered processes question.answered events.
func (h *Handler) HandleQuestionAnswered(ctx context.Context, _, value []byte) error {
	var event events.QuestionAnswered
	if err := json.Unmarshal(value, &event); err != nil {
		return fmt.Errorf("decode question answered: %w", err)
	}
	if err := h.progress.ProcessQuestionAnswered(ctx, event); err != nil {
		return err
	}
	if h.readiness != nil {
		return h.readiness.ProcessQuestionAnswered(ctx, event)
	}
	return nil
}

// HandleStreakUpdated processes streak.updated events.
func (h *Handler) HandleStreakUpdated(ctx context.Context, _, value []byte) error {
	var event events.StreakUpdated
	if err := json.Unmarshal(value, &event); err != nil {
		return fmt.Errorf("decode streak updated: %w", err)
	}
	return h.progress.ProcessStreakUpdated(ctx, event)
}

// RunQuestionAnswered starts the question.answered consumer loop.
func RunQuestionAnswered(ctx context.Context, consumer *kafka.Consumer, handler *Handler) error {
	return consumer.Run(ctx, handler.HandleQuestionAnswered)
}

// RunStreakUpdated starts the streak.updated consumer loop.
func RunStreakUpdated(ctx context.Context, consumer *kafka.Consumer, handler *Handler) error {
	return consumer.Run(ctx, handler.HandleStreakUpdated)
}
