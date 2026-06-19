package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/prepio/prepio/services/question/internal/dto"
	"github.com/prepio/prepio/services/question/internal/store"
)

// ErrNodeNotFound is returned when a journey node does not exist.
var ErrNodeNotFound = errors.New("journey node not found")

// ContentService exposes read operations on question pools and node content bindings.
type ContentService struct {
	content *store.ContentStore
	journey *store.JourneyStore
}

// NewContentService creates a ContentService.
func NewContentService(content *store.ContentStore, journey *store.JourneyStore) *ContentService {
	return &ContentService{content: content, journey: journey}
}

// GetNodeContent returns skills and pools bound to a journey node.
func (s *ContentService) GetNodeContent(ctx context.Context, nodeID string) (*dto.NodeContentResponse, error) {
	if len(nodeID) == 0 {
		return nil, fmt.Errorf("node id is required")
	}

	node, err := s.journey.GetNodeByID(ctx, nodeID)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, ErrNodeNotFound
	}

	content, err := s.content.GetNodeContent(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	skills := make([]dto.NodeSkillResponse, 0, len(content.Skills))
	for _, binding := range content.Skills {
		skills = append(skills, dto.NodeSkillResponse{
			SkillSlug: binding.SkillSlug,
			SkillName: binding.SkillName,
			IsPrimary: binding.IsPrimary,
		})
	}

	pools := make([]dto.NodePoolResponse, 0, len(content.Pools))
	for _, binding := range content.Pools {
		questionIDs, err := s.content.ListPoolQuestionIDs(ctx, binding.PoolID)
		if err != nil {
			return nil, err
		}
		pools = append(pools, dto.NodePoolResponse{
			PoolSlug:          binding.PoolSlug,
			PoolName:          binding.PoolName,
			SkillSlug:         binding.SkillSlug,
			SelectionStrategy: binding.SelectionStrategy,
			QuestionsRequired: binding.QuestionsRequired,
			QuestionCount:     len(questionIDs),
		})
	}

	return &dto.NodeContentResponse{
		NodeID: node.ID,
		Slug:   node.Slug,
		Label:  node.Label,
		Skills: skills,
		Pools:  pools,
	}, nil
}

// EnrichJourneyNodes attaches skill and pool bindings to journey nodes.
func (s *ContentService) EnrichJourneyNodes(ctx context.Context, nodes []dto.JourneyNodeResponse) error {
	for i := range nodes {
		content, err := s.content.GetNodeContent(ctx, nodes[i].ID)
		if err != nil {
			return err
		}

		skills := make([]dto.NodeSkillResponse, 0, len(content.Skills))
		for _, binding := range content.Skills {
			skills = append(skills, dto.NodeSkillResponse{
				SkillSlug: binding.SkillSlug,
				SkillName: binding.SkillName,
				IsPrimary: binding.IsPrimary,
			})
		}

		pools := make([]dto.NodePoolResponse, 0, len(content.Pools))
		for _, binding := range content.Pools {
			questionIDs, err := s.content.ListPoolQuestionIDs(ctx, binding.PoolID)
			if err != nil {
				return err
			}
			pools = append(pools, dto.NodePoolResponse{
				PoolSlug:          binding.PoolSlug,
				PoolName:          binding.PoolName,
				SkillSlug:         binding.SkillSlug,
				SelectionStrategy: binding.SelectionStrategy,
				QuestionsRequired: binding.QuestionsRequired,
				QuestionCount:     len(questionIDs),
			})
		}

		nodes[i].Skills = skills
		nodes[i].Pools = pools
	}
	return nil
}
