package service

import (
	"context"
	"fmt"

	"github.com/prepio/prepio/services/question/internal/dto"
	"github.com/prepio/prepio/services/question/internal/store"
)

// SkillService exposes read operations on the skill graph.
type SkillService struct {
	skills *store.SkillStore
}

// NewSkillService creates a SkillService.
func NewSkillService(skills *store.SkillStore) *SkillService {
	return &SkillService{skills: skills}
}

// ListSkillTree returns categories with nested skills and subskills.
func (s *SkillService) ListSkillTree(ctx context.Context) ([]dto.SkillCategoryResponse, error) {
	categories, err := s.skills.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.SkillCategoryResponse, 0, len(categories))
	for _, category := range categories {
		skills, err := s.skills.ListSkillsByCategory(ctx, category.ID)
		if err != nil {
			return nil, err
		}

		skillResponses := make([]dto.SkillResponse, 0, len(skills))
		for _, skill := range skills {
			subskills, err := s.skills.ListSubskillsBySkillID(ctx, skill.ID)
			if err != nil {
				return nil, err
			}

			subskillResponses := make([]dto.SubskillResponse, 0, len(subskills))
			for _, subskill := range subskills {
				subskillResponses = append(subskillResponses, dto.SubskillResponse{
					ID:   subskill.ID,
					Slug: subskill.Slug,
					Name: subskill.Name,
				})
			}

			skillResponses = append(skillResponses, dto.SkillResponse{
				ID:          skill.ID,
				Slug:        skill.Slug,
				Name:        skill.Name,
				Description: skill.Description,
				Subskills:   subskillResponses,
			})
		}

		result = append(result, dto.SkillCategoryResponse{
			ID:     category.ID,
			Slug:   category.Slug,
			Name:   category.Name,
			Skills: skillResponses,
		})
	}
	return result, nil
}

// GetSkillBySlug returns a skill with its subskills.
func (s *SkillService) GetSkillBySlug(ctx context.Context, slug string) (*dto.SkillResponse, error) {
	if len(slug) == 0 {
		return nil, fmt.Errorf("skill slug is required")
	}

	skill, err := s.skills.GetSkillBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if skill == nil {
		return nil, ErrSkillNotFound
	}

	subskills, err := s.skills.ListSubskillsBySkillID(ctx, skill.ID)
	if err != nil {
		return nil, err
	}

	subskillResponses := make([]dto.SubskillResponse, 0, len(subskills))
	for _, subskill := range subskills {
		subskillResponses = append(subskillResponses, dto.SubskillResponse{
			ID:   subskill.ID,
			Slug: subskill.Slug,
			Name: subskill.Name,
		})
	}

	return &dto.SkillResponse{
		ID:          skill.ID,
		Slug:        skill.Slug,
		Name:        skill.Name,
		Description: skill.Description,
		Subskills:   subskillResponses,
	}, nil
}
