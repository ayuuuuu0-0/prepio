package service

import (
	"context"

	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/dto"
	"github.com/prepio/prepio/services/question/internal/store"
)

const foundationWorldSlug = constants.FoundationForestWorldSlug

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

	sessionAnswered := map[string]bool{}
	records, err := s.history.ListBySession(ctx, userID, paper.SessionID)
	if err != nil {
		return nil, err
	}
	for _, rec := range records {
		sessionAnswered[rec.QuestionID] = true
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

	if config.JourneyPoolSelectionEnabled() {
		return s.buildPoolJourney(ctx, userID, paper, nodes, sessionAnswered, resp)
	}
	return s.buildIndexJourney(ctx, userID, paper, nodes, sessionAnswered, resp)
}

func (s *QuestionService) buildIndexJourney(
	ctx context.Context,
	userID string,
	paper *dto.DailyPaperResponse,
	nodes []store.JourneyNode,
	sessionAnswered map[string]bool,
	resp *dto.JourneyResponse,
) (*dto.JourneyResponse, error) {
	currentSet := false
	for i, node := range nodes {
		status := "locked"
		questionID := ""
		if i < len(paper.Questions) {
			q := paper.Questions[i]
			questionID = q.ID
			if sessionAnswered[q.ID] {
				status = "done"
				_ = s.journey.UpsertProgress(ctx, userID, node.ID, "done")
			} else if !currentSet {
				status = "current"
				currentSet = true
			}
		} else if node.NodeType == "boss" && !currentSet && len(paper.Questions) > 0 {
			if bossUnlockedIndex(node, paper, sessionAnswered) {
				status = "current"
				currentSet = true
			}
		}

		resp.Nodes = append(resp.Nodes, dto.JourneyNodeResponse{
			ID:         node.ID,
			Slug:       node.Slug,
			Label:      node.Label,
			NodeType:   node.NodeType,
			Status:     status,
			QuestionID: questionID,
			SortOrder:  node.SortOrder,
		})
	}
	return resp, nil
}

func (s *QuestionService) buildPoolJourney(
	ctx context.Context,
	userID string,
	paper *dto.DailyPaperResponse,
	nodes []store.JourneyNode,
	sessionAnswered map[string]bool,
	resp *dto.JourneyResponse,
) (*dto.JourneyResponse, error) {
	seen, err := s.history.AnsweredQuestionSet(ctx, userID)
	if err != nil {
		return nil, err
	}

	assignments := make([]nodeAssignment, len(nodes))
	for i, node := range nodes {
		assignments[i] = s.resolveNodeAssignment(ctx, true, userID, paper.SessionID, node, i, paper, seen)
	}

	currentSet := false
	for i, node := range nodes {
		assignment := assignments[i]
		status := "locked"
		questionID := displayQuestionID(assignment.questionIDs, sessionAnswered)

		if allAnsweredInSession(assignment.questionIDs, sessionAnswered) {
			status = "done"
			_ = s.journey.UpsertProgress(ctx, userID, node.ID, "done")
		} else if !currentSet && nodeUnlocked(i, assignments, sessionAnswered) && len(assignment.questionIDs) > 0 {
			status = "current"
			currentSet = true
		}

		resp.Nodes = append(resp.Nodes, dto.JourneyNodeResponse{
			ID:         node.ID,
			Slug:       node.Slug,
			Label:      node.Label,
			NodeType:   node.NodeType,
			Status:     status,
			QuestionID: questionID,
			SortOrder:  node.SortOrder,
		})
	}
	return resp, nil
}
