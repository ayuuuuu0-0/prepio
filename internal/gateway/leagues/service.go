package leagues

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Service aggregates the weekly league: standings come from Progress (which owns XP and
// leagues), names and companions from User. If either fails the league fails; it never
// shows a leaderboard with invented or missing people.
type Service struct {
	progressURL string
	userURL     string
	client      *http.Client
}

// NewService creates a league aggregation service.
func NewService(progressURL, userURL string) *Service {
	return &Service{progressURL: progressURL, userURL: userURL, client: &http.Client{Timeout: 10 * time.Second}}
}

// Tier is one rung of the league ladder.
type Tier struct {
	Index int    `json:"index"`
	Slug  string `json:"slug"`
	Name  string `json:"name"`
}

// Standing is one leaderboard row with the learner's public card attached.
type Standing struct {
	Rank             int    `json:"rank"`
	UserID           string `json:"user_id"`
	WeeklyXP         int    `json:"weekly_xp"`
	Zone             string `json:"zone"`
	IsMe             bool   `json:"is_me"`
	Username         string `json:"username"`
	CompanionName    string `json:"companion_name"`
	CompanionSpecies string `json:"companion_species"`
}

// Result is how the learner's most recent finished week ended.
type Result struct {
	WeekStart  string `json:"week_start"`
	Rank       int    `json:"rank"`
	CohortSize int    `json:"cohort_size"`
	Outcome    string `json:"outcome"`
	FromTier   Tier   `json:"from_tier"`
	ToTier     Tier   `json:"to_tier"`
}

// League is returned by GET /api/v1/league.
type League struct {
	WeekStart string     `json:"week_start"`
	EndsAt    string     `json:"ends_at"`
	Joined    bool       `json:"joined"`
	Tier      Tier       `json:"tier"`
	Tiers     []Tier     `json:"tiers"`
	MyRank    int        `json:"my_rank"`
	Standings []Standing `json:"standings"`
	Last      *Result    `json:"last_result"`
}

type publicCard struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	CompanionName    string `json:"companion_name"`
	CompanionSpecies string `json:"companion_species"`
}

// GetLeague returns the authenticated learner's league for this week.
func (s *Service) GetLeague(ctx context.Context, token string) (*League, error) {
	var league League
	if err := s.call(ctx, http.MethodGet, s.progressURL+"/api/v1/progress/league", token, nil, &league); err != nil {
		return nil, err
	}
	if league.Standings == nil {
		league.Standings = []Standing{}
	}
	if len(league.Standings) == 0 {
		return &league, nil
	}

	ids := make([]string, len(league.Standings))
	for i, st := range league.Standings {
		ids[i] = st.UserID
	}
	var cards []publicCard
	if err := s.call(ctx, http.MethodPost, s.userURL+"/api/v1/users/public-cards", token,
		map[string][]string{"user_ids": ids}, &cards); err != nil {
		return nil, err
	}
	byID := make(map[string]publicCard, len(cards))
	for _, c := range cards {
		byID[c.ID] = c
	}
	for i := range league.Standings {
		card, ok := byID[league.Standings[i].UserID]
		if !ok {
			return nil, fmt.Errorf("league: no public card for user %s", league.Standings[i].UserID)
		}
		league.Standings[i].Username = card.Username
		league.Standings[i].CompanionName = card.CompanionName
		league.Standings[i].CompanionSpecies = card.CompanionSpecies
	}
	return &league, nil
}

func (s *Service) call(ctx context.Context, method, url, token string, in, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode >= 300 {
		return fmt.Errorf("upstream %s: status %d", url, res.StatusCode)
	}
	envelope := struct {
		Data any `json:"data"`
	}{Data: out}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return fmt.Errorf("decode %s: %w", url, err)
	}
	return nil
}
