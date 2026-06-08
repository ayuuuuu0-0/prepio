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
	questionURL string
	client      *http.Client
}

// NewService creates a dashboard aggregation service.
func NewService(userURL, progressURL, streakURL, questionURL string) *Service {
	return &Service{
		userURL:     userURL,
		progressURL: progressURL,
		streakURL:   streakURL,
		questionURL: questionURL,
		client:      &http.Client{},
	}
}

// HomeResponse is returned by GET /api/v1/dashboard/home.
type HomeResponse struct {
	Streak           StreakCard       `json:"streak"`
	Progress         ProgressCard     `json:"progress"`
	Companion        CompanionCard    `json:"companion"`
	Readiness        []ReadinessCard  `json:"readiness"`
	League           LeagueCard       `json:"league"`
	DailyQuests      []DailyQuestCard `json:"daily_quests"`
	CompanionMessage string           `json:"companion_message"`
	OnboardingNeeded bool             `json:"onboarding_needed"`
}

// StreakCard summarizes streak state.
type StreakCard struct {
	CurrentStreak     int  `json:"current_streak"`
	LongestStreak     int  `json:"longest_streak"`
	FreezeCount       int  `json:"freeze_count"`
	StreakActiveToday bool `json:"streak_active_today"`
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

// LeagueCard summarizes league placement (placeholder until Phase 9).
type LeagueCard struct {
	Tier      string `json:"tier"`
	Rank      int    `json:"rank"`
	Label     string `json:"label"`
	Available bool   `json:"available"`
}

// DailyQuestCard is a daily quest entry.
type DailyQuestCard struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Progress    int    `json:"progress"`
	Target      int    `json:"target"`
	Completed   bool   `json:"completed"`
	RewardXP    int    `json:"reward_xp"`
	RewardGems  int    `json:"reward_gems"`
	ComingSoon  bool   `json:"coming_soon"`
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

	readinessStats, _ := s.fetchReadinessStats(ctx, token)

	resp := &HomeResponse{
		Progress:         progress,
		Streak:           streak,
		OnboardingNeeded: !profile.OnboardingCompleted,
		League: LeagueCard{
			Tier:      "",
			Rank:      0,
			Label:     "Leagues launching soon",
			Available: false,
		},
		DailyQuests: comingSoonQuests(),
	}

	if profile.Companion != nil {
		resp.Companion = CompanionCard{
			ID:      profile.Companion.ID,
			Name:    profile.Companion.Name,
			Species: profile.Companion.Species,
		}
	}

	resp.Readiness = computeReadiness(profile.TargetCompanies, readinessStats)
	resp.CompanionMessage = companionMessage(resp.Companion.Name, progress)

	return resp, nil
}

type profilePayload struct {
	OnboardingCompleted bool           `json:"onboarding_completed"`
	TargetCompanies     []string       `json:"target_companies"`
	Companion           *CompanionCard `json:"companion"`
}

type companyStatsPayload struct {
	Company  string `json:"company"`
	Answered int    `json:"answered"`
	Correct  int    `json:"correct"`
	ScoreAvg int    `json:"score_avg"`
}

type readinessStatsPayload struct {
	ByCompany []companyStatsPayload `json:"by_company"`
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

func (s *Service) fetchReadinessStats(ctx context.Context, token string) (*readinessStatsPayload, error) {
	if len(s.questionURL) == 0 {
		return &readinessStatsPayload{}, nil
	}
	body, err := s.get(ctx, s.questionURL+"/api/v1/questions/stats/readiness", token)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data readinessStatsPayload `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("decode readiness stats: %w", err)
	}
	return &envelope.Data, nil
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

// computeReadiness derives company readiness from actual answer history.
func computeReadiness(targets []string, stats *readinessStatsPayload) []ReadinessCard {
	if len(targets) == 0 {
		return []ReadinessCard{}
	}

	byCompany := map[string]companyStatsPayload{}
	if stats != nil {
		for _, row := range stats.ByCompany {
			byCompany[row.Company] = row
		}
	}

	cards := make([]ReadinessCard, 0, len(targets))
	for _, company := range targets {
		row, ok := byCompany[company]
		score := 0
		if ok && row.Answered > 0 {
			score = (row.Correct * 100) / row.Answered
			if row.ScoreAvg > 0 {
				score = (score + row.ScoreAvg) / 2
			}
			if score > 95 {
				score = 95
			}
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
	level := progress.CurrentLevel
	challenges := 0
	if xpNeeded > 0 {
		challenges = (xpNeeded + config.XPByDifficulty["medium"] - 1) / config.XPByDifficulty["medium"]
	}

	switch {
	case progress.TotalXP == 0:
		return fmt.Sprintf("%s is ready. Let's see what you can do.", name)
	case challenges <= 1:
		return fmt.Sprintf("One challenge from Level %d. Don't stop now.", level+1)
	case challenges <= 3:
		return fmt.Sprintf("%d challenges from Level %d. Google Readiness is watching.", challenges, level+1)
	case progress.CurrentLevel < 5:
		return fmt.Sprintf("Level %d. The real prep starts around Level 10 — keep going.", level)
	default:
		return fmt.Sprintf("Level %d. %d challenges from Level %d. Consistency compounds.", level, challenges, level+1)
	}
}

func comingSoonQuests() []DailyQuestCard {
	return []DailyQuestCard{
		{ID: "daily_question", Title: "Complete today's challenge", Progress: 0, Target: 1, Completed: false, RewardXP: 50, RewardGems: 10, ComingSoon: true},
		{ID: "maintain_streak", Title: "Keep the streak alive", Progress: 0, Target: 1, Completed: false, RewardXP: 20, RewardGems: 5, ComingSoon: true},
		{ID: "score_high", Title: "Score above 80% on a challenge", Progress: 0, Target: 1, Completed: false, RewardXP: 30, RewardGems: 5, ComingSoon: true},
	}
}

func extractBearer(header string) string {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
