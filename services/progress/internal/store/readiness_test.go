package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/services/progress/internal/store"
	"github.com/prepio/prepio/test/testdb"
	"github.com/stretchr/testify/require"
)

func TestReadinessStoreUpsertAndList(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	ctx := context.Background()
	readinessStore := store.NewReadinessStore(pool)

	var userID string
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash)
		VALUES ('r@test.com', 'ruser', 'hash') RETURNING id`).Scan(&userID))

	skillID := "b2000001-0000-4000-8000-000000000002"
	practicedAt := time.Now().UTC()
	require.NoError(t, readinessStore.UpsertUserSkillScore(
		ctx, userID, skillID, 72, 3, practicedAt, config.ReadinessSourceLive,
	))

	scores, err := readinessStore.ListUserSkillScores(ctx, userID)
	require.NoError(t, err)
	require.Len(t, scores, 1)
	require.Equal(t, "arrays", scores[0].SkillSlug)
	require.Equal(t, 72, scores[0].Mastery)
	require.Equal(t, 3, scores[0].Attempts)
}

func TestReadinessStoreCompanyWeights(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	ctx := context.Background()
	readinessStore := store.NewReadinessStore(pool)

	weights, err := readinessStore.ListCompanySkillWeights(ctx, "google")
	require.NoError(t, err)
	require.Len(t, weights, 7)

	total := 0
	for _, weight := range weights {
		total += weight.Weight
		require.NotEmpty(t, weight.SkillSlug)
	}
	require.Equal(t, 100, total)
}

func TestReadinessStoreQuestionSkillContributions(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	ctx := context.Background()
	readinessStore := store.NewReadinessStore(pool)

	contributions, err := readinessStore.ListQuestionSkillContributions(
		ctx, "b0000000-0000-4000-8000-000000000001",
	)
	require.NoError(t, err)
	require.Len(t, contributions, 2)
}

func TestReadinessStoreBackfillFromHistory(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	ctx := context.Background()
	readinessStore := store.NewReadinessStore(pool)

	var userID string
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash)
		VALUES ('rb@test.com', 'rbuser', 'hash') RETURNING id`).Scan(&userID))

	_, err := pool.Exec(ctx, `
		INSERT INTO user_question_history (user_id, question_id, correct, score, submitted_at, session_id)
		VALUES ($1, 'b0000000-0000-4000-8000-000000000001', true, 85, now(), 'a0000000-0000-4000-8000-000000000099')`, userID)
	require.NoError(t, err)

	scores, err := readinessStore.ListUserSkillScores(ctx, userID)
	require.NoError(t, err)
	require.NotEmpty(t, scores)
}
