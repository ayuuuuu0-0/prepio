package service

import (
	"context"
	"math"

	"github.com/prepio/prepio/services/progress/internal/dto"
	"github.com/prepio/prepio/services/progress/internal/store"
)

// BuildTopicMastery turns per-skill rows into per-topic readiness.
//
// Topic readiness is the mean mastery of the skills the user has started, with every
// started skill weighted equally. Skills not yet practiced are reported as "not started"
// (nil mastery) and do not count against the topic; SkillsStarted / SkillsTotal tells the
// learner how much of the topic they have covered. A topic with nothing started has nil
// mastery. It never uses XP, levels, or streaks.
func BuildTopicMastery(rows []store.TopicSkillRow) []dto.TopicMasteryResponse {
	topics := make([]dto.TopicMasteryResponse, 0)
	index := map[string]int{}
	sums := map[string]int{}

	for _, r := range rows {
		i, ok := index[r.TopicSlug]
		if !ok {
			i = len(topics)
			index[r.TopicSlug] = i
			topics = append(topics, dto.TopicMasteryResponse{
				Slug:        r.TopicSlug,
				Name:        r.TopicName,
				Description: r.TopicDescription,
				Skills:      []dto.TopicSkillResponse{},
			})
		}

		skill := dto.TopicSkillResponse{Slug: r.SkillSlug, Name: r.SkillName}
		topics[i].SkillsTotal++
		if r.Attempts > 0 && r.Mastery != nil {
			m := *r.Mastery
			skill.Mastery = &m
			topics[i].SkillsStarted++
			sums[r.TopicSlug] += m
		}
		topics[i].Skills = append(topics[i].Skills, skill)
	}

	for i := range topics {
		if started := topics[i].SkillsStarted; started > 0 {
			mean := int(math.Round(float64(sums[topics[i].Slug]) / float64(started)))
			topics[i].Mastery = &mean
		}
	}
	return topics
}

// GetTopicMastery returns the user's readiness in every topic.
func (s *ReadinessService) GetTopicMastery(ctx context.Context, userID string) ([]dto.TopicMasteryResponse, error) {
	if len(userID) == 0 {
		return nil, ErrInvalidRequest
	}
	rows, err := s.readiness.ListTopicSkills(ctx, userID)
	if err != nil {
		return nil, err
	}
	return BuildTopicMastery(rows), nil
}
