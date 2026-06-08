package testing

import (
	"github.com/prepio/prepio/services/question/internal/dto"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/prepio/prepio/services/question/internal/store"
	"github.com/redis/go-redis/v9"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewService wires the question service for integration tests.
func NewService(pool *pgxpool.Pool, redisClient *redis.Client, publisher service.EventPublisher) *service.QuestionService {
	return service.NewQuestionService(
		store.NewQuestionStore(pool),
		store.NewDailyPaperStore(pool),
		store.NewHistoryStore(pool),
		store.NewUserStore(pool),
		redisClient,
		publisher,
	)
}

// SubmitRequest re-exports the submit DTO for integration tests.
type SubmitRequest = dto.SubmitRequest
