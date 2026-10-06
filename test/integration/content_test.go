package integration_test

import (
	"context"
	"path/filepath"
	"testing"

	questiontest "github.com/prepio/prepio/services/question/testing"
	"github.com/prepio/prepio/test/fakes"
	"github.com/prepio/prepio/test/testdb"
	"github.com/stretchr/testify/require"
)

// TestAuthoredContent validates the real content/ directory against the migrated
// skill catalog and proves content-sync is idempotent on it.
func TestAuthoredContent(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	ctx := context.Background()

	content, err := questiontest.LoadContent(filepath.Join("..", "..", "content"))
	require.NoError(t, err)
	require.NotEmpty(t, content.Lessons, "the path must ship at least one real lesson")

	sync := questiontest.NewContentSync(pool)
	skills, err := sync.SkillSlugs(ctx)
	require.NoError(t, err)
	require.Empty(t, questiontest.Validate(content, skills))

	first, err := sync.Sync(ctx, content)
	require.NoError(t, err)
	require.Equal(t, len(content.Lessons), first.LessonsCreated)

	second, err := sync.Sync(ctx, content)
	require.NoError(t, err)
	require.Zero(t, second.LessonsCreated)
	require.Zero(t, second.LessonsNewVersion)
	require.Equal(t, len(content.Lessons), second.LessonsUnchanged)

	// A published lesson is reachable from the path of a brand-new user.
	var userID string
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash) VALUES ('c@test.com', 'content', 'h') RETURNING id`).Scan(&userID))
	path, err := questiontest.NewLessonService(pool, &fakes.KafkaProducer{}).GetPath(ctx, userID)
	require.NoError(t, err)
	require.NotEmpty(t, path.Worlds)
	require.Equal(t, "current", path.Worlds[0].Nodes[0].Status)
}
