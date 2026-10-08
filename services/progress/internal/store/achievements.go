package store

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/config"
)

// EarnedAchievement is an achievement a learner holds.
type EarnedAchievement struct {
	Slug       string
	UnlockedAt time.Time
}

// AchievementStore reads earned achievements. Awarding happens inside the transactions that
// change the facts (lesson rewards, streak bonuses), through unlockAchievements.
type AchievementStore struct {
	pool *pgxpool.Pool
}

// NewAchievementStore creates an AchievementStore.
func NewAchievementStore(pool *pgxpool.Pool) *AchievementStore {
	return &AchievementStore{pool: pool}
}

// List returns every achievement the learner has earned.
func (s *AchievementStore) List(ctx context.Context, userID string) ([]EarnedAchievement, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT achievement_slug, unlocked_at FROM user_achievements WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("list achievements: %w", err)
	}
	defer rows.Close()
	var out []EarnedAchievement
	for rows.Next() {
		var a EarnedAchievement
		if err := rows.Scan(&a.Slug, &a.UnlockedAt); err != nil {
			return nil, fmt.Errorf("scan achievement: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// AwardStreak stores the streak achievements a streak of days earns. Idempotent.
func (s *AchievementStore) AwardStreak(ctx context.Context, userID string, days int) ([]string, error) {
	return unlockAchievements(ctx, s.pool, userID, config.AchievementFacts{Streak: days}, nil)
}

// unlockAchievements stores every achievement the facts earn and returns the ones that are
// new. Already-earned achievements are left untouched, so calling it again is safe.
func unlockAchievements(ctx context.Context, q querier, userID string, facts config.AchievementFacts, attemptID *string) ([]string, error) {
	var unlocked []string
	for _, slug := range config.EarnedAchievements(facts) {
		tag, err := q.Exec(ctx, `
			INSERT INTO user_achievements (user_id, achievement_slug, attempt_id) VALUES ($1, $2, $3)
			ON CONFLICT (user_id, achievement_slug) DO NOTHING`, userID, slug, attemptID)
		if err != nil {
			return nil, fmt.Errorf("unlock achievement %s: %w", slug, err)
		}
		if tag.RowsAffected() == 1 {
			unlocked = append(unlocked, slug)
		}
	}
	return unlocked, nil
}

// lessonFacts gathers what a lesson completion can earn: lessons completed, level, best
// topic readiness, and whether the learner has ever been promoted in a league.
func lessonFacts(ctx context.Context, q querier, userID string, level int) (config.AchievementFacts, error) {
	facts := config.AchievementFacts{Level: level}
	if err := q.QueryRow(ctx,
		`SELECT count(*) FROM lesson_rewards WHERE user_id = $1 AND first_completion`, userID).Scan(&facts.Lessons); err != nil {
		return facts, fmt.Errorf("count lessons: %w", err)
	}
	if err := q.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM league_memberships WHERE user_id = $1 AND tier > 0)`, userID).Scan(&facts.Promoted); err != nil {
		return facts, fmt.Errorf("check promotion: %w", err)
	}

	// Topic readiness exactly as GET /progress/topics computes it: the rounded mean mastery
	// of the started skills (attempts > 0) in each topic.
	rows, err := q.Query(ctx, `
		SELECT sum(u.mastery), count(*)
		FROM user_skill_scores u
		JOIN skills sk ON sk.id = u.skill_id
		JOIN skill_categories c ON c.id = sk.category_id
		WHERE u.user_id = $1 AND u.attempts > 0 AND c.topic_id IS NOT NULL
		GROUP BY c.topic_id`, userID)
	if err != nil {
		return facts, fmt.Errorf("topic readiness: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var sum, started int
		if err := rows.Scan(&sum, &started); err != nil {
			return facts, fmt.Errorf("scan topic readiness: %w", err)
		}
		if mean := int(math.Round(float64(sum) / float64(started))); mean > facts.BestTopicMastery {
			facts.BestTopicMastery = mean
		}
	}
	return facts, rows.Err()
}

// AttemptAchievements returns the achievements a lesson attempt earned.
func (s *AchievementStore) AttemptAchievements(ctx context.Context, userID, attemptID string) ([]string, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT achievement_slug FROM user_achievements
		WHERE user_id = $1 AND attempt_id = $2`, userID, attemptID)
	if err != nil {
		return nil, fmt.Errorf("attempt achievements: %w", err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		out = append(out, slug)
	}
	return out, rows.Err()
}
