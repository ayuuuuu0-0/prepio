package integration_test

import (
	"context"
	"os"
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
	catalog, err := sync.Catalog(ctx)
	require.NoError(t, err)
	require.Empty(t, questiontest.Validate(content, catalog))

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
	path, err := questiontest.NewLessonService(pool, &fakes.KafkaProducer{}).GetPath(ctx, userID, nil)
	require.NoError(t, err)
	require.NotEmpty(t, path.Worlds)
	require.Equal(t, "current", path.Worlds[0].Nodes[0].Status)

	// Every authored world names its topic, and focus topics put their worlds first without
	// locking anything: the focused world gets the current node, the others stay available.
	for _, w := range path.Worlds {
		require.NotEmpty(t, w.Topic, "world %s names its topic", w.Slug)
		require.False(t, w.Focused, "no focus was given")
	}
	focused, err := questiontest.NewLessonService(pool, &fakes.KafkaProducer{}).GetPath(ctx, userID, []string{"dsa-refresher", "low-level-design"})
	require.NoError(t, err)
	require.Equal(t, "dsa-refresher", focused.Worlds[0].Topic)
	require.Equal(t, "low-level-design", focused.Worlds[1].Topic)
	require.True(t, focused.Worlds[0].Focused && focused.Worlds[1].Focused && !focused.Worlds[2].Focused)
	require.Equal(t, "current", focused.Worlds[0].Nodes[0].Status)
	require.Equal(t, "available", focused.Worlds[2].Nodes[0].Status, "an unfocused world stays open")
	require.Len(t, focused.Worlds, len(path.Worlds), "focus reorders, never hides")
}

// TestContentSyncReplacesLessonOnNode proves a lesson can be replaced on its node:
// the removed lesson is deprecated (never deleted) and the new one serves the path.
// It also proves two lessons can swap nodes in one sync.
func TestContentSyncReplacesLessonOnNode(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	ctx := context.Background()

	sync := questiontest.NewContentSync(pool)
	lessons := questiontest.NewLessonService(pool, &fakes.KafkaProducer{})
	catalog, err := sync.Catalog(ctx)
	require.NoError(t, err)
	userID := newUser(t, pool, "replacer")

	content := testContent()
	_, err = sync.Sync(ctx, content)
	require.NoError(t, err)

	lessonStatus := func(slug string) (id, status string) {
		t.Helper()
		require.NoError(t, pool.QueryRow(ctx, `SELECT id, status FROM lessons WHERE slug = $1`, slug).Scan(&id, &status))
		return id, status
	}
	pathLessons := func() []string {
		t.Helper()
		path, err := lessons.GetPath(ctx, userID, nil)
		require.NoError(t, err)
		require.Len(t, path.Worlds, 1)
		ids := []string{}
		for _, n := range path.Worlds[0].Nodes {
			ids = append(ids, n.LessonID)
		}
		return ids
	}

	// Replace lesson-one by lesson-three on the same node.
	replaced := testContent()
	replaced.Lessons[0].Slug = "lesson-three"
	replaced.Lessons[0].Title = "Lesson Three"
	require.Empty(t, questiontest.Validate(replaced, catalog))
	report, err := sync.Sync(ctx, replaced)
	require.NoError(t, err)
	require.Equal(t, 1, report.LessonsCreated)
	require.Equal(t, 1, report.LessonsDeprecated)

	oneID, oneStatus := lessonStatus("lesson-one")
	threeID, threeStatus := lessonStatus("lesson-three")
	twoID, _ := lessonStatus("lesson-two")
	require.Equal(t, "deprecated", oneStatus)
	require.Equal(t, "published", threeStatus)
	require.Equal(t, []string{threeID, twoID}, pathLessons(), "the path serves the replacement, never both")
	require.NotContains(t, pathLessons(), oneID)

	// Swap the two live lessons between their nodes in one sync.
	swapped := testContent()
	swapped.Lessons[0].Slug, swapped.Lessons[0].Title = "lesson-three", "Lesson Three"
	swapped.Lessons[0].Node, swapped.Lessons[1].Node = "n-next", "n-start"
	require.Empty(t, questiontest.Validate(swapped, catalog))
	_, err = sync.Sync(ctx, swapped)
	require.NoError(t, err)
	require.Equal(t, []string{twoID, threeID}, pathLessons())

	// Two live lessons on one node are still rejected (checked at commit).
	_, err = pool.Exec(ctx, `UPDATE lessons SET status = 'published' WHERE id = $1`, oneID)
	require.Error(t, err, "a node binds at most one live lesson")
}

// TestLessonNodeUniqueMigrationRollsBack proves 000040 rolls back and re-applies.
func TestLessonNodeUniqueMigrationRollsBack(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	ctx := context.Background()

	run := func(suffix string) {
		t.Helper()
		raw, err := os.ReadFile(filepath.Join("..", "..", "migrations", "000040_lessons_node_unique_live."+suffix+".sql"))
		require.NoError(t, err)
		_, err = pool.Exec(ctx, string(raw))
		require.NoError(t, err, suffix)
	}
	constraint := func(name string) bool {
		var ok bool
		require.NoError(t, pool.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM pg_constraint WHERE conrelid = 'lessons'::regclass AND conname = $1)`, name).Scan(&ok))
		return ok
	}

	require.True(t, constraint("lessons_node_id_live_excl"))
	require.False(t, constraint("lessons_node_id_key"))
	run("down")
	require.False(t, constraint("lessons_node_id_live_excl"))
	require.True(t, constraint("lessons_node_id_key"))
	run("up")
	require.True(t, constraint("lessons_node_id_live_excl"))
}
