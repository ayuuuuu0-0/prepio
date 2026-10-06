package service

import (
	"context"

	"github.com/prepio/prepio/services/question/internal/dto"
)

// ListTopics returns the topic catalog in display order.
func (s *SkillService) ListTopics(ctx context.Context) ([]dto.TopicResponse, error) {
	topics, err := s.skills.ListTopics(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]dto.TopicResponse, 0, len(topics))
	for _, t := range topics {
		resp = append(resp, dto.TopicResponse{Slug: t.Slug, Name: t.Name, Description: t.Description})
	}
	return resp, nil
}
