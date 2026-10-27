package dashboard

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/prepio/prepio/config"
)

// Service aggregates dashboard data from upstream services. Aggregation does not change
// ownership: every number here is computed and owned by the service it came from.
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
	Streak           StreakCard    `json:"streak"`
	Progress         ProgressCard  `json:"progress"`
	Companion        CompanionCard `json:"companion"`
	Topics           []TopicCard   `json:"topics"`
	FocusTopics      []string      `json:"focus_topics"`
	NextLesson       *NextLesson   `json:"next_lesson"`
	League           LeagueCard    `json:"league"`
	CompanionMessage string        `json:"companion_message"`
	OnboardingNeeded bool          `json:"onboarding_needed"`
}

// LeagueCard is the learner's weekly league at a glance. Until they earn XP this week
// Joined is false and Rank, CohortSize, WeeklyXP, and Zone are zero values.
type LeagueCard struct {
	TierIndex  int    `json:"tier_index"`
	TierSlug   string `json:"tier_slug"`
	TierName   string `json:"tier_name"`
	Joined     bool   `json:"joined"`
	Rank       int    `json:"rank"`
	CohortSize int    `json:"cohort_size"`
	WeeklyXP   int    `json:"weekly_xp"`
	Zone       string `json:"zone"`
	EndsAt     string `json:"ends_at"`
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

// TopicCard is the learner's readiness in one topic. Mastery is nil until a skill in the
// topic has been practiced.
type TopicCard struct {
	Slug          string `json:"slug"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Mastery       *int   `json:"mastery"`
	SkillsStarted int    `json:"skills_started"`
	SkillsTotal   int    `json:"skills_total"`
	Focused       bool   `json:"focused"`
}

// NextLesson is the lesson the Continue button opens.
type NextLesson struct {
	LessonID   string `json:"lesson_id"`
	Title      string `json:"title"`
	NodeLabel  string `json:"node_label"`
	WorldName  string `json:"world_name"`
	EstMinutes int    `json:"est_minutes"`
	XPPreview  int    `json:"xp_preview"`
	InProgress bool   `json:"in_progress"`
	Kind       string `json:"kind"`
}

// GetHome aggregates dashboard data for the authenticated user. Upstream calls run in
// parallel; if any fails the dashboard fails, so it never shows partial or invented data.
func (s *Service) GetHome(ctx context.Context, token string) (*HomeResponse, error) {
	var (
		wg      sync.WaitGroup
		profile *profilePayload
		prog    ProgressCard
		streak  StreakCard
		topics  []TopicCard
		path    *pathPayload
		league  LeagueCard
		errs    = make([]error, 5)
	)

	run := func(i int, fn func() error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[i] = fn()
		}()
	}
	// The path is ordered by the learner's focus topics, so it follows the profile; the
	// other upstreams run alongside.
	run(0, func() (err error) {
		if profile, err = s.fetchProfile(ctx, token); err != nil {
			return err
		}
		path, err = s.fetchPath(ctx, token, profile.FocusTopics)
		return err
	})
	run(1, func() (err error) { prog, err = s.fetchProgress(ctx, token); return })
	run(2, func() (err error) { streak, err = s.fetchStreak(ctx, token); return })
	run(3, func() (err error) { topics, err = s.fetchTopics(ctx, token); return })
	run(4, func() (err error) { league, err = s.fetchLeague(ctx, token); return })
	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	resp := &HomeResponse{
		Progress:         prog,
		Streak:           streak,
		Topics:           orderTopics(topics, profile.FocusTopics),
		FocusTopics:      profile.FocusTopics,
		NextLesson:       nextLesson(path),
		League:           league,
		OnboardingNeeded: !profile.OnboardingCompleted,
	}
	if resp.FocusTopics == nil {
		resp.FocusTopics = []string{}
	}
	if profile.Companion != nil {
		resp.Companion = CompanionCard{
			ID:      profile.Companion.ID,
			Name:    profile.Companion.Name,
			Species: profile.Companion.Species,
		}
	}
	resp.CompanionMessage = companionMessage(resp.Companion.Name, prog)
	return resp, nil
}

// orderTopics puts the learner's focus topics first, in their chosen priority, then the rest
// in catalog order. Focus never hides or locks a topic.
func orderTopics(topics []TopicCard, focus []string) []TopicCard {
	focused := make(map[string]bool, len(focus))
	for _, slug := range focus {
		focused[slug] = true
	}
	ordered := make([]TopicCard, 0, len(topics))
	for _, slug := range focus {
		for _, t := range topics {
			if t.Slug == slug {
				t.Focused = true
				ordered = append(ordered, t)
			}
		}
	}
	for _, t := range topics {
		if !focused[t.Slug] {
			ordered = append(ordered, t)
		}
	}
	return ordered
}

// nextLesson picks the node the path marks "current". Nil means there is nothing to continue
// (no worlds yet, or everything is finished).
func nextLesson(path *pathPayload) *NextLesson {
	if path == nil {
		return nil
	}
	for _, w := range path.Worlds {
		for _, n := range w.Nodes {
			if n.Status == "current" {
				return &NextLesson{
					LessonID:   n.LessonID,
					Title:      n.Title,
					NodeLabel:  n.Label,
					WorldName:  w.Name,
					EstMinutes: n.EstMinutes,
					XPPreview:  n.XPPreview,
					InProgress: n.InProgress,
					Kind:       n.Kind,
				}
			}
		}
	}
	return nil
}

type profilePayload struct {
	OnboardingCompleted bool           `json:"onboarding_completed"`
	Companion           *CompanionCard `json:"companion"`
	FocusTopics         []string       `json:"focus_topics"`
}

type pathPayload struct {
	Worlds []struct {
		Name  string `json:"name"`
		Nodes []struct {
			Label      string `json:"label"`
			Status     string `json:"status"`
			LessonID   string `json:"lesson_id"`
			Title      string `json:"title"`
			Kind       string `json:"kind"`
			EstMinutes int    `json:"est_minutes"`
			XPPreview  int    `json:"xp_preview"`
			InProgress bool   `json:"in_progress"`
		} `json:"nodes"`
	} `json:"worlds"`
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

func (s *Service) fetchTopics(ctx context.Context, token string) ([]TopicCard, error) {
	body, err := s.get(ctx, s.progressURL+"/api/v1/progress/topics", token)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data []TopicCard `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("decode topics: %w", err)
	}
	return envelope.Data, nil
}

func (s *Service) fetchPath(ctx context.Context, token string, focus []string) (*pathPayload, error) {
	u := s.questionURL + "/api/v1/path"
	if len(focus) > 0 {
		u += "?focus=" + url.QueryEscape(strings.Join(focus, ","))
	}
	body, err := s.get(ctx, u, token)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data pathPayload `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("decode path: %w", err)
	}
	return &envelope.Data, nil
}

// fetchLeague summarizes the learner's weekly league from Progress. Names and companions
// are not needed here, so the User service is not involved.
func (s *Service) fetchLeague(ctx context.Context, token string) (LeagueCard, error) {
	body, err := s.get(ctx, s.progressURL+"/api/v1/progress/league", token)
	if err != nil {
		return LeagueCard{}, err
	}
	var envelope struct {
		Data struct {
			EndsAt string `json:"ends_at"`
			Joined bool   `json:"joined"`
			Tier   struct {
				Index int    `json:"index"`
				Slug  string `json:"slug"`
				Name  string `json:"name"`
			} `json:"tier"`
			MyRank    int `json:"my_rank"`
			Standings []struct {
				WeeklyXP int    `json:"weekly_xp"`
				Zone     string `json:"zone"`
				IsMe     bool   `json:"is_me"`
			} `json:"standings"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return LeagueCard{}, fmt.Errorf("decode league: %w", err)
	}
	d := envelope.Data
	card := LeagueCard{TierIndex: d.Tier.Index, TierSlug: d.Tier.Slug, TierName: d.Tier.Name, Joined: d.Joined, EndsAt: d.EndsAt}
	if d.Joined {
		card.Rank, card.CohortSize = d.MyRank, len(d.Standings)
		for _, st := range d.Standings {
			if st.IsMe {
				card.WeeklyXP, card.Zone = st.WeeklyXP, st.Zone
			}
		}
	}
	return card, nil
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

func companionMessage(name string, progress ProgressCard) string {
	if len(name) == 0 {
		name = "Your companion"
	}
	xpNeeded := progress.XPToNextLevel
	level := progress.CurrentLevel
	lessons := 0
	if xpNeeded > 0 {
		perLesson := config.LessonXPByDifficulty["medium"]
		lessons = (xpNeeded + perLesson - 1) / perLesson
	}

	switch {
	case progress.TotalXP == 0:
		return fmt.Sprintf("%s is ready. Let's see what you can do.", name)
	case lessons <= 1:
		return fmt.Sprintf("One lesson from Level %d. Don't stop now.", level+1)
	case lessons <= 3:
		return fmt.Sprintf("%d lessons from Level %d. Your companion is watching.", lessons, level+1)
	case progress.CurrentLevel < 5:
		return fmt.Sprintf("Level %d. The real prep starts around Level 10 — keep going.", level)
	default:
		return fmt.Sprintf("Level %d. %d lessons from Level %d. Consistency compounds.", level, lessons, level+1)
	}
}
