package integration_test

import (
	"context"
	"testing"

	"github.com/prepio/prepio/constants"
	progresstest "github.com/prepio/prepio/services/progress/testing"
	questiontest "github.com/prepio/prepio/services/question/testing"
	usertest "github.com/prepio/prepio/services/user/testing"
	"github.com/prepio/prepio/test/testdb"
	"github.com/stretchr/testify/require"
)

func TestTopicsOnboardingAndMastery(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)
	ctx := context.Background()

	userID := newUser(t, pool, "topical")
	onboarding := usertest.NewOnboardingService(pool)
	companion := constants.StarterCompanionIDs[0]

	t.Run("the catalog has the four topics in display order", func(t *testing.T) {
		topics, err := questiontest.NewSkillService(pool).ListTopics(ctx)
		require.NoError(t, err)
		var slugs []string
		for _, tp := range topics {
			slugs = append(slugs, tp.Slug)
			require.NotEmpty(t, tp.Name)
			require.NotEmpty(t, tp.Description)
		}
		require.Equal(t, []string{"system-design", "backend-production", "low-level-design", "dsa-refresher"}, slugs)
	})

	t.Run("onboarding needs one to three distinct, real topics", func(t *testing.T) {
		bad := map[string][]string{
			"none":      nil,
			"empty":     {},
			"four":      {"system-design", "backend-production", "low-level-design", "dsa-refresher"},
			"duplicate": {"system-design", "system-design"},
			"unknown":   {"system-design", "astrology"},
			"blank":     {""},
		}
		for name, topics := range bad {
			_, err := onboarding.Complete(ctx, userID, usertest.OnboardingRequest{
				ExperienceLevel: "mid", CompanionID: companion, FocusTopics: topics,
			})
			require.ErrorIs(t, err, usertest.ErrInvalidRequest, name)
		}

		var done bool
		require.NoError(t, pool.QueryRow(ctx, `SELECT onboarding_completed FROM users WHERE id = $1`, userID).Scan(&done))
		require.False(t, done, "a rejected onboarding must not complete the profile")
		var n int
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM user_focus_topics WHERE user_id = $1`, userID).Scan(&n))
		require.Zero(t, n, "a rejected onboarding must not store partial topics")
	})

	t.Run("valid topics are stored in priority order and returned on the profile", func(t *testing.T) {
		profile, err := onboarding.Complete(ctx, userID, usertest.OnboardingRequest{
			ExperienceLevel: "mid", CompanionID: companion, FocusTopics: []string{"dsa-refresher", "backend-production"},
		})
		require.NoError(t, err)
		require.True(t, profile.OnboardingCompleted)
		require.Equal(t, []string{"dsa-refresher", "backend-production"}, profile.FocusTopics)

		again, err := onboarding.GetProfile(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, profile.FocusTopics, again.FocusTopics)
	})

	t.Run("changing the topics replaces them, and an invalid change keeps the old ones", func(t *testing.T) {
		_, err := onboarding.Complete(ctx, userID, usertest.OnboardingRequest{
			ExperienceLevel: "mid", CompanionID: companion, FocusTopics: []string{"system-design", "nope"},
		})
		require.ErrorIs(t, err, usertest.ErrInvalidRequest)
		kept, err := onboarding.GetProfile(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, []string{"dsa-refresher", "backend-production"}, kept.FocusTopics)

		profile, err := onboarding.Complete(ctx, userID, usertest.OnboardingRequest{
			ExperienceLevel: "senior", CompanionID: companion, FocusTopics: []string{"system-design"},
		})
		require.NoError(t, err)
		require.Equal(t, []string{"system-design"}, profile.FocusTopics)
	})

	t.Run("a new user starts every topic at not started", func(t *testing.T) {
		topics, err := progresstest.NewReadinessService(pool).GetTopicMastery(ctx, userID)
		require.NoError(t, err)
		require.Len(t, topics, 4)
		total := map[string]int{}
		for _, tp := range topics {
			require.Nil(t, tp.Mastery, tp.Slug)
			require.Zero(t, tp.SkillsStarted, tp.Slug)
			require.NotZero(t, tp.SkillsTotal, tp.Slug)
			total[tp.Slug] = tp.SkillsTotal
		}
		require.Equal(t, 4, total["system-design"])
		require.Equal(t, 4, total["backend-production"])
		require.Equal(t, 3, total["low-level-design"])
		require.Equal(t, 14, total["dsa-refresher"], "9 data-structure skills and 5 algorithm skills group under one topic")
	})

	t.Run("topic mastery is the mean of the started skills", func(t *testing.T) {
		_, err := pool.Exec(ctx, `
			INSERT INTO user_skill_scores (user_id, skill_id, mastery, attempts, last_practiced_at)
			SELECT $1, id, m.mastery, 1, now()
			FROM skills s JOIN (VALUES ('system-design-scaling', 30), ('system-design-fundamentals', 51)) AS m(slug, mastery) ON m.slug = s.slug`, userID)
		require.NoError(t, err)

		topics, err := progresstest.NewReadinessService(pool).GetTopicMastery(ctx, userID)
		require.NoError(t, err)
		sd := topics[0]
		require.Equal(t, "system-design", sd.Slug)
		require.NotNil(t, sd.Mastery)
		require.Equal(t, 41, *sd.Mastery, "(30+51)/2 = 40.5 rounds to 41")
		require.Equal(t, 2, sd.SkillsStarted)
		require.Equal(t, 4, sd.SkillsTotal)
		require.Nil(t, topics[1].Mastery, "other topics are untouched")
	})
}
