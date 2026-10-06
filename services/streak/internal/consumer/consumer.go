package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prepio/prepio/services/streak/internal/service"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/kafka"
)

// HandleLessonCompleted processes lesson.completed events.
func HandleLessonCompleted(streaks *service.StreakService) kafka.MessageHandler {
	return func(ctx context.Context, _, value []byte) error {
		var event events.LessonCompleted
		if err := json.Unmarshal(value, &event); err != nil {
			return fmt.Errorf("decode lesson completed: %w", err)
		}
		return streaks.ProcessLessonCompleted(ctx, event)
	}
}
