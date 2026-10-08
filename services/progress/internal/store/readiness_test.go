package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/prepio/prepio/services/progress/internal/store"
	"github.com/prepio/prepio/test/testdb"
	"github.com/stretchr/testify/require"
)

func TestReadinessStoreListTopicSkills(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	ctx := context.Background()
	readinessStore := store.NewReadinessStore(pool)

	var userID string
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash)
		VALUES ('r@test.com', 'ruser', 'hash') RETURNING id`).Scan(&userID))

	skillID := "b2000001-0000-4000-8000-000000000002"
	_, err := pool.Exec(ctx, `
		INSERT INTO user_skill_scores (user_id, skill_id, mastery, attempts, last_practiced_at, source)
		VALUES ($1, $2, 72, 3, $3, 'live')`, userID, skillID, time.Now().UTC())
	require.NoError(t, err)

	rows, err := readinessStore.ListTopicSkills(ctx, userID)
	require.NoError(t, err)
	require.NotEmpty(t, rows)

	var practiced *store.TopicSkillRow
	for i := range rows {
		if rows[i].SkillSlug == "arrays" {
			practiced = &rows[i]
		} else {
			require.Nil(t, rows[i].Mastery, rows[i].SkillSlug)
		}
	}
	require.NotNil(t, practiced, "arrays skill should belong to a topic")
	require.NotNil(t, practiced.Mastery)
	require.Equal(t, 72, *practiced.Mastery)
	require.Equal(t, 3, practiced.Attempts)
}
