package service_test

import (
	"context"
	"testing"

	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/prepio/prepio/services/question/internal/store"
	"github.com/prepio/prepio/test/testdb"
	"github.com/stretchr/testify/require"
)

func TestSkillServiceListSkillTree(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	skillService := service.NewSkillService(store.NewSkillStore(pool))
	ctx := context.Background()

	tree, err := skillService.ListSkillTree(ctx)
	require.NoError(t, err)
	require.Len(t, tree, 8)

	var foundArrays bool
	for _, category := range tree {
		for _, skill := range category.Skills {
			if skill.Slug == "arrays" {
				foundArrays = true
				require.NotEmpty(t, skill.Subskills)
			}
		}
	}
	require.True(t, foundArrays)
}

func TestSkillServiceGetSkillBySlugNotFound(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	skillService := service.NewSkillService(store.NewSkillStore(pool))
	ctx := context.Background()

	_, err := skillService.GetSkillBySlug(ctx, "nonexistent-skill")
	require.ErrorIs(t, err, service.ErrSkillNotFound)
}

func TestSkillServiceListQuestionSkills(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	skillService := service.NewSkillService(store.NewSkillStore(pool))
	ctx := context.Background()

	mappings, err := skillService.ListQuestionSkills(ctx, "b0000000-0000-4000-8000-000000000009")
	require.NoError(t, err)
	require.Len(t, mappings, 1)
	require.Equal(t, "behavioral-star", mappings[0].SkillSlug)
}
