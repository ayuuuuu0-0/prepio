package service

import (
	"context"
	"time"

	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/services/progress/internal/dto"
	"github.com/prepio/prepio/services/progress/internal/store"
)

// LeagueService serves weekly league standings. League XP is written by LessonStore in
// the same transaction that grants the XP.
type LeagueService struct {
	leagues *store.LeagueStore
	now     func() time.Time
}

// NewLeagueService creates a LeagueService.
func NewLeagueService(leagues *store.LeagueStore) *LeagueService {
	return &LeagueService{leagues: leagues, now: time.Now}
}

// GetLeague returns the learner's league for the current week.
func (s *LeagueService) GetLeague(ctx context.Context, userID string) (*dto.LeagueResponse, error) {
	if !validIDs(userID) {
		return nil, ErrInvalidRequest
	}
	view, err := s.leagues.Get(ctx, userID, s.now())
	if err != nil {
		return nil, err
	}
	return buildLeagueResponse(view, userID), nil
}

func buildLeagueResponse(view *store.LeagueView, userID string) *dto.LeagueResponse {
	resp := &dto.LeagueResponse{
		WeekStart: view.WeekStart.Format("2006-01-02"),
		EndsAt:    view.WeekStart.AddDate(0, 0, 7).Format(time.RFC3339),
		Joined:    view.Joined,
		Tier:      tierResponse(view.Tier),
		Tiers:     make([]dto.LeagueTierResponse, len(config.LeagueTiers)),
		Standings: make([]dto.LeagueStandingResponse, 0, len(view.Standings)),
	}
	for i := range config.LeagueTiers {
		resp.Tiers[i] = tierResponse(i)
	}

	size := len(view.Standings)
	for i, st := range view.Standings {
		rank := i + 1
		outcome, _ := config.LeagueOutcome(rank, size, view.Tier)
		row := dto.LeagueStandingResponse{
			Rank: rank, UserID: st.UserID, WeeklyXP: st.WeeklyXP,
			Zone: zoneFor(outcome), IsMe: st.UserID == userID,
		}
		if row.IsMe {
			resp.MyRank = rank
		}
		resp.Standings = append(resp.Standings, row)
	}

	if last := view.Last; last != nil {
		resp.Last = &dto.LeagueResultResponse{
			WeekStart:  last.WeekStart.Format("2006-01-02"),
			Rank:       last.Rank,
			CohortSize: last.CohortSize,
			Outcome:    last.Outcome,
			FromTier:   tierResponse(last.Tier),
			ToTier:     tierResponse(last.NextTier),
		}
	}
	return resp
}

func tierResponse(i int) dto.LeagueTierResponse {
	i = config.ClampLeagueTier(i)
	t := config.LeagueTiers[i]
	return dto.LeagueTierResponse{Index: i, Slug: t.Slug, Name: t.Name}
}

func zoneFor(outcome string) string {
	switch outcome {
	case config.LeaguePromoted:
		return "promote"
	case config.LeagueDemoted:
		return "demote"
	default:
		return "stay"
	}
}
