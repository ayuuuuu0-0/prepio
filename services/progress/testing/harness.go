package testing

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/internal/progress/service"
	"github.com/prepio/prepio/internal/progress/store"
)

// NewService wires the progress service for integration tests.
func NewService(pool *pgxpool.Pool, publisher service.EventPublisher) *service.ProgressService {
	return service.NewProgressService(
		store.NewProgressStore(pool),
		store.NewAchievementStore(pool),
		publisher,
	)
}
