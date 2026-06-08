package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/notification/internal/store"
	"github.com/prepio/prepio/shared/events"
	"github.com/redis/go-redis/v9"
)

// NotificationService dispatches notifications with daily caps.
type NotificationService struct {
	store *store.NotificationStore
	redis *redis.Client
}

// NewNotificationService creates a NotificationService.
func NewNotificationService(store *store.NotificationStore, redisClient *redis.Client) *NotificationService {
	return &NotificationService{store: store, redis: redisClient}
}

// HandleProgressUpdated sends a level-up notification when applicable.
func (s *NotificationService) HandleProgressUpdated(ctx context.Context, event events.ProgressUpdated) error {
	if !event.LeveledUp {
		return nil
	}
	return s.dispatch(ctx, event.UserID, events.NotificationLevelUp, map[string]any{
		"level_after": event.LevelAfter,
	})
}

// HandleStreakUpdated sends streak broken notifications.
func (s *NotificationService) HandleStreakUpdated(ctx context.Context, event events.StreakUpdated) error {
	if !event.StreakBroken {
		return nil
	}
	return s.dispatch(ctx, event.UserID, events.NotificationStreakBroken, nil)
}

// HandleDispatch sends an explicit notification request.
func (s *NotificationService) HandleDispatch(ctx context.Context, event events.NotificationsDispatch) error {
	return s.dispatch(ctx, event.UserID, event.NotificationType, event.Metadata)
}

func (s *NotificationService) dispatch(ctx context.Context, userID, notificationType string, metadata map[string]any) error {
	if len(userID) == 0 {
		return fmt.Errorf("user id is required")
	}

	today := time.Now().UTC().Format("20060102")
	capKey := constants.NotifCapKey(userID, today)

	count, err := s.redis.Incr(ctx, capKey).Result()
	if err != nil {
		return fmt.Errorf("increment notif cap: %w", err)
	}
	if count == 1 {
		s.redis.Expire(ctx, capKey, 36*time.Hour)
	}
	if count > 3 {
		return nil
	}

	log.Printf("fcm stub: user=%s type=%s metadata=%v", userID, notificationType, metadata)
	return s.store.Insert(ctx, userID, notificationType, "fcm")
}
