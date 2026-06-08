package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/prepio/prepio/config"
)

// Service aggregates dashboard data from upstream microservices.
type Service struct {
	userURL     string
	progressURL string
	streakURL   string
	client      *http.Client
}

// NewService creates a dashboard aggregation service.
func NewService(userURL, progressURL, streakURL string) *Service {
	return &Service{
		userURL:     userURL,
		progressURL: progressURL,
		streakURL:   streakURL,
		client:      &http.Client{},
	}
}

// HomeResponse is returned by GET /api/v1/dashboard/home.
type HomeResponse struct {
	Streak            StreakCard            `json:"streak"`
	Progress          ProgressCard          `json:"progress"`
	Companion         CompanionCard         `json:"companion"`
	Readiness         []ReadinessCard       `json:"readiness"`
	League            LeagueCard            `json:"league"`
	DailyQuests       []DailyQuestCard      `json:"daily_quests"`
	CompanionMessage  string                `json:"companion_message"`
	OnboardingNeeded  bool                  `json:"onboarding_needed"`
}

// StreakCard summarizes streak state.
type StreakCard struct {
	CurrentStreak     int    `json:"current_streak"`
	LongestStreak     int    `json:"longest_streak"`
	FreezeCount       int    `json:"freeze_count"`
	StreakActiveToday bool   `json:"streak_active_today"`
}

// ProgressCard summarizes XP and gems.
type ProgressCard struct {
	TotalXP       int `json:"total_xp"`
	CurrentLevel  int `json:"current_level"`
	GemBalance    int `json:"gem_balance"`
	XPToNextLevel int `json:"xp_to_next_level"`
}

// CompanionCard summarizes the active companion.
type CompanionCard struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Species string `json:"species"`
}

// ReadinessCard shows company-specific readiness.
type ReadinessCard struct {
	Company string `json:"company"`
	Score   int    `json:"score"`
}

// LeagueCard is a placeholder until Phase 9.
type LeagueCard struct {
	Tier  string `json:"tier"`
	Rank  int    `json:"rank"`
	Label string `json:"label"`
}

// DailyQuestCard is a placeholder until Phase 8.
type DailyQuestCard struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Progress    int    `json:"progress"`
	Target      int    `json:"target"`
	Completed   bool   `json:"completed"`
	RewardXP    int    `json:"reward_xp"`
	RewardGems  int    `json:"reward_gems"`
}

// GetHome aggregates dashboard data for the authenticated user.
func (s *Service) GetHome(ctx context.Context, token string) (*HomeResponse, error) {
	profile, err := s.fetchProfile(ctx, token)
	if err != nil {
		return nil, err
	}

	progress, err := s.fetchProgress(ctx, token)
	if err != nil {
		return nil, err
	}

	streak, err := s.fetchStreak(ctx, token)
	if err != nil {
		return nil, err
	}

	resp := &HomeResponse{
		Progress:         progress,
		Streak:           streak,
		OnboardingNeeded: !profile.OnboardingCompleted,
		League: LeagueCard{
			Tier:  "bronze",
			Rank:  1,
			Label: "Bronze League",
		},
		DailyQuests: defaultDailyQuests(streak.CurrentStreak > 0),
	}

	if profile.Companion != nil {
		resp.Companion = CompanionCard{
			ID:      profile.Companion.ID,
			Name:    profile.Companion.Name,
			Species: profile.Companion.Species,
		}
	}

	resp.Readiness = computeReadiness(profile.TargetCompanies, progress)
	resp.CompanionMessage = companionMessage(resp.Companion.Name, progress)

	return resp, nil
}

type profilePayload struct {
	OnboardingCompleted bool `json:"onboarding_completed"`
	TargetCompanies     []string `json:"target_companies"`
	Companion           *CompanionCard `json:"companion"`
}

func (s *Service) fetchProfile(ctx context.Context, token string) (*profilePayload, error) {
	body, err := s.get(ctx, s.userURL+"/api/v1/users/profile", token)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data profilePayload `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("decode profile: %w", err)
	}
	return &envelope.Data, nil
}

func (s *Service) fetchProgress(ctx context.Context, token string) (ProgressCard, error) {
	body, err := s.get(ctx, s.progressURL+"/api/v1/progress/me", token)
	if err != nil {
		return ProgressCard{}, err
	}
	var envelope struct {
		Data ProgressCard `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return ProgressCard{}, fmt.Errorf("decode progress: %w", err)
	}
	return envelope.Data, nil
}

func (s *Service) fetchStreak(ctx context.Context, token string) (StreakCard, error) {
	body, err := s.get(ctx, s.streakURL+"/api/v1/streaks/me", token)
	if err != nil {
		return StreakCard{}, err
	}
	var envelope struct {
		Data StreakCard `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return StreakCard{}, fmt.Errorf("decode streak: %w", err)
	}
	return envelope.Data, nil
}

func (s *Service) get(ctx context.Context, url, token string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("upstream %s: status %d", url, res.StatusCode)
	}
	return body, nil
}

func computeReadiness(targets []string, progress ProgressCard) []ReadinessCard {
	base := 10 + progress.CurrentLevel*4
	if progress.TotalXP > 0 {
		base += progress.TotalXP / 40
	}
	if base > 95 {
		base = 95
	}

	if len(targets) == 0 {
		return []ReadinessCard{}
	}

	cards := make([]ReadinessCard, 0, len(targets))
	for i, company := range targets {
		score := base - i*2
		if score < 5 {
			score = 5
		}
		cards = append(cards, ReadinessCard{Company: company, Score: score})
	}
	return cards
}

func companionMessage(name string, progress ProgressCard) string {
	if len(name) == 0 {
		name = "Your companion"
	}
	xpNeeded := progress.XPToNextLevel
	if xpNeeded <= 0 {
		return fmt.Sprintf("%s says: You're crushing it — keep the journey going!", name)
	}
	challenges := (xpNeeded + config.XPByDifficulty["medium"] - 1) / config.XPByDifficulty["medium"]
	return fmt.Sprintf("%s says: Only %d challenges until Level %d.", name, challenges, progress.CurrentLevel+1)
}

func defaultDailyQuests(streakActive bool) []DailyQuestCard {
	return []DailyQuestCard{
		{ID: "daily_question", Title: "Complete today's challenge", Progress: 0, Target: 1, Completed: false, RewardXP: 50, RewardGems: 10},
		{ID: "maintain_streak", Title: "Maintain your streak", Progress: boolToInt(streakActive), Target: 1, Completed: streakActive, RewardXP: 20, RewardGems: 5},
		{ID: "score_high", Title: "Score above 80% on a challenge", Progress: 0, Target: 1, Completed: false, RewardXP: 30, RewardGems: 5},
	}
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func extractBearer(header string) string {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
