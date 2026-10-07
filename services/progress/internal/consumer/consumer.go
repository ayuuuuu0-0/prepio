package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/kafka"
)

// HandleStreakUpdated processes streak.updated events.
func HandleStreakUpdated(progress *service.ProgressService) kafka.MessageHandler {
	return func(ctx context.Context, _, value []byte) error {
		var event events.StreakUpdated
		if err := json.Unmarshal(value, &event); err != nil {
			return fmt.Errorf("decode streak updated: %w", err)
		}
		return progress.ProcessStreakUpdated(ctx, event)
	}
}

// HandleLessonCompleted processes lesson.completed events.
func HandleLessonCompleted(lessons *service.LessonService) kafka.MessageHandler {
	return func(ctx context.Context, _, value []byte) error {
		var event events.LessonCompleted
		if err := json.Unmarshal(value, &event); err != nil {
			return fmt.Errorf("decode lesson completed: %w", err)
		}
		return lessons.ProcessLessonCompleted(ctx, event)
	}
}
