package testing

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/services/progress/internal/store"
)

// NewService wires the progress service for integration tests.
func NewService(pool *pgxpool.Pool, publisher service.EventPublisher) *service.ProgressService {
	return service.NewProgressService(
		store.NewProgressStore(pool),
		store.NewLedgerStore(pool),
		publisher,
	)
}
