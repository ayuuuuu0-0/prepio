package lesson_test

import (
	"path/filepath"
	"testing"

	"github.com/prepio/prepio/services/question/internal/lesson"
	"github.com/stretchr/testify/require"
)

// TestRepoContentIsValid validates the authored content/ directory structurally.
// It needs no database, so it always runs in CI; skill references are checked
// against the migrated catalog in the integration tests.
func TestRepoContentIsValid(t *testing.T) {
	content, err := lesson.Load(filepath.Join("..", "..", "..", "..", "content"))
	require.NoError(t, err)
	require.NotEmpty(t, content.Worlds)
	require.NotEmpty(t, content.Lessons)
	require.Empty(t, lesson.Validate(content, lesson.Catalog{}))
}
