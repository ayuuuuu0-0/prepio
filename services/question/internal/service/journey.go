package service

import (
	"context"

	"github.com/prepio/prepio/services/question/internal/dto"
)

const foundationWorldSlug = "foundation-forest"

// GetJourney returns the user's current world path with live node status.
func (s *QuestionService) GetJourney(ctx context.Context, userID, timezone string) (*dto.JourneyResponse, error) {
	if len(userID) == 0 {
		return nil, ErrInvalidRequest
	}

	world, err := s.journey.GetWorldBySlug(ctx, foundationWorldSlug)
	if err != nil {
		return nil, err
	}

	nodes, err := s.journey.ListNodesByWorld(ctx, world.ID)
	if err != nil {
		return nil, err
	}

	paper, err := s.GetDailyPaper(ctx, userID, timezone)
	if err != nil {
		return nil, err
	}

	answered := map[string]bool{}
	records, err := s.history.ListBySession(ctx, userID, paper.SessionID)
	if err != nil {
		return nil, err
	}
	for _, rec := range records {
		answered[rec.QuestionID] = true
	}

	resp := &dto.JourneyResponse{
		World: dto.JourneyWorldResponse{
			ID:          world.ID,
			Slug:        world.Slug,
			Name:        world.Name,
			Description: world.Description,
			Theme:       world.Theme,
		},
		SessionID: paper.SessionID,
	}

	currentSet := false
	for i, node := range nodes {
		status := "locked"
		questionID := ""
		if i < len(paper.Questions) {
			q := paper.Questions[i]
			questionID = q.ID
			if answered[q.ID] {
				status = "done"
				_ = s.journey.UpsertProgress(ctx, userID, node.ID, "done")
			} else if !currentSet {
				status = "current"
				currentSet = true
			}
		} else if node.NodeType == "boss" && !currentSet && len(paper.Questions) > 0 {
			allDone := true
			for _, q := range paper.Questions {
				if !answered[q.ID] {
					allDone = false
					break
				}
			}
			if allDone {
				status = "current"
				currentSet = true
			}
		}

		resp.Nodes = append(resp.Nodes, dto.JourneyNodeResponse{
			ID:         node.ID,
			Label:      node.Label,
			NodeType:   node.NodeType,
			Status:     status,
			QuestionID: questionID,
			SortOrder:  node.SortOrder,
		})
	}

	return resp, nil
}
