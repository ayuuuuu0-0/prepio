package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/prepio/prepio/services/question/internal/store"
	"github.com/prepio/prepio/test/fakes"
	"github.com/prepio/prepio/test/testdb"
	"github.com/prepio/prepio/test/testredis"
	"github.com/stretchr/testify/require"
)

func TestGetJourneyPoolSelection(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	redisClient, _ := testredis.New(t)
	ctx := context.Background()

	userID := seedJourneyUser(t, pool)

	journeyStore := store.NewJourneyStore(pool)
	contentStore := store.NewContentStore(pool)
	questionService := service.NewQuestionService(
		store.NewQuestionStore(pool),
		store.NewDailyPaperStore(pool),
		store.NewHistoryStore(pool),
		journeyStore,
		contentStore,
		store.NewUserStore(pool),
		redisClient,
		&fakes.KafkaProducer{},
	)

	t.Run("index mode when flag disabled", func(t *testing.T) {
		require.NoError(t, os.Unsetenv(constants.EnvJourneyPoolSelection))

		resp, err := questionService.GetJourney(ctx, userID, "UTC")
		require.NoError(t, err)
		require.Len(t, resp.Nodes, 5)
		require.Equal(t, "current", resp.Nodes[0].Status)
		require.NotEmpty(t, resp.Nodes[0].QuestionID)
	})

	t.Run("pool mode assigns foundation forest pool questions", func(t *testing.T) {
		require.NoError(t, os.Setenv(constants.EnvJourneyPoolSelection, "true"))
		t.Cleanup(func() { _ = os.Unsetenv(constants.EnvJourneyPoolSelection) })

		resp, err := questionService.GetJourney(ctx, userID, "UTC")
		require.NoError(t, err)
		require.Len(t, resp.Nodes, 5)

		first := resp.Nodes[0]
		require.Equal(t, "arrays-basics", first.Slug)
		require.Equal(t, "current", first.Status)
		require.Equal(t, "b0000000-0000-4000-8000-000000000001", first.QuestionID)

		boss := resp.Nodes[4]
		require.Equal(t, "forest-boss", boss.Slug)
		require.Equal(t, "locked", boss.Status)
	})
}

func seedJourneyUser(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	ctx := context.Background()
	var userID string
	err := pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash)
		VALUES ('journey@test.com', 'journeyuser', 'hash')
		RETURNING id`).Scan(&userID)
	require.NoError(t, err)
	return userID
}
