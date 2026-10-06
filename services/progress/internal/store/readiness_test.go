package store_test

import (
	"context"
	"testing"
	"time"

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
	_, err := pool.Exec(ctx, `
		INSERT INTO user_skill_scores (user_id, skill_id, mastery, attempts, last_practiced_at, source)
		VALUES ($1, $2, 72, 3, $3, 'live')`, userID, skillID, practicedAt)
	require.NoError(t, err)

	scores, err := readinessStore.ListUserSkillScores(ctx, userID)
	require.NoError(t, err)
	require.Len(t, scores, 1)
	require.Equal(t, "arrays", scores[0].SkillSlug)
	require.Equal(t, 72, scores[0].Mastery)
	require.Equal(t, 3, scores[0].Attempts)
}
