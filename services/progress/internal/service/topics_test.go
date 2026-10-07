package service_test

import (
	"testing"

	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/services/progress/internal/store"
	"github.com/stretchr/testify/require"
)

func intp(v int) *int { return &v }

func row(topic, skill string, mastery *int, attempts int) store.TopicSkillRow {
	return store.TopicSkillRow{
		TopicSlug: topic, TopicName: topic + "-name", TopicDescription: topic + "-desc",
		SkillSlug: skill, SkillName: skill + "-name", Mastery: mastery, Attempts: attempts,
	}
}

func TestBuildTopicMastery(t *testing.T) {
	topics := service.BuildTopicMastery([]store.TopicSkillRow{
		row("sd", "caching", intp(40), 1),
		row("sd", "sharding", nil, 0),
		row("sd", "queues", intp(61), 2),
		row("be", "api", nil, 0),
		row("be", "db", nil, 0),
		row("dsa", "arrays", intp(0), 1), // practiced but still at zero: started, not "not started"
	})
	require.Len(t, topics, 3)

	sd := topics[0]
	require.Equal(t, "sd", sd.Slug)
	require.Equal(t, "sd-name", sd.Name)
	require.NotNil(t, sd.Mastery)
	require.Equal(t, 51, *sd.Mastery, "mean of the started skills only: (40+61)/2 rounds to 51")
	require.Equal(t, 2, sd.SkillsStarted)
	require.Equal(t, 3, sd.SkillsTotal)
	require.Nil(t, sd.Skills[1].Mastery, "an unstarted skill is not a failure")
	require.Equal(t, 40, *sd.Skills[0].Mastery)

	be := topics[1]
	require.Nil(t, be.Mastery, "nothing started means no number at all")
	require.Zero(t, be.SkillsStarted)
	require.Equal(t, 2, be.SkillsTotal)

	dsa := topics[2]
	require.NotNil(t, dsa.Mastery)
	require.Equal(t, 0, *dsa.Mastery)
	require.Equal(t, 1, dsa.SkillsStarted)
}

func TestBuildTopicMasteryKeepsCatalogOrderAndIsNeverNil(t *testing.T) {
	require.NotNil(t, service.BuildTopicMastery(nil))
	require.Empty(t, service.BuildTopicMastery(nil))

	topics := service.BuildTopicMastery([]store.TopicSkillRow{row("b", "x", nil, 0), row("a", "y", nil, 0)})
	require.Equal(t, []string{"b", "a"}, []string{topics[0].Slug, topics[1].Slug}, "the query decides order, not the builder")
}
