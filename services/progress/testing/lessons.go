package testing

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/services/progress/internal/dto"
	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/services/progress/internal/store"
)

// AttemptRewards re-exports the rewards DTO for integration tests.
type AttemptRewards = dto.AttemptRewardsResponse

// NewLessonService wires the progress lesson service for integration tests.
func NewLessonService(pool *pgxpool.Pool, publisher service.EventPublisher) *service.LessonService {
	return service.NewLessonService(store.NewLessonStore(pool), publisher)
}

// TopicMastery re-exports the topic readiness DTO for integration tests.
type TopicMastery = dto.TopicMasteryResponse

// NewReadinessService wires the readiness service for integration tests.
func NewReadinessService(pool *pgxpool.Pool) *service.ReadinessService {
	return service.NewReadinessService(store.NewReadinessStore(pool))
}
