package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/shared/events"
)

// MasteryChange is one skill's mastery movement caused by a lesson attempt.
type MasteryChange struct {
	SkillID   string
	SkillSlug string
	SkillName string
	TopicSlug string
	TopicName string
	Before    int
	After     int
	Delta     int
	Accuracy  float64
}

// LessonOutcome is what applying one lesson.completed event changed.
type LessonOutcome struct {
	AttemptID       string
	LessonID        string
	FirstCompletion bool
	XPAwarded       int
	GemsAwarded     int
	LevelBefore     int
	LevelAfter      int
	TotalXP         int
	TotalGems       int
	Mastery         []MasteryChange
	CreatedAt       time.Time
}

// LessonStore applies lesson completions to mastery, XP, and gems, atomically.
type LessonStore struct {
	pool *pgxpool.Pool
}

// NewLessonStore creates a LessonStore.
func NewLessonStore(pool *pgxpool.Pool) *LessonStore {
	return &LessonStore{pool: pool}
}

// ApplyLessonCompleted records the rewards and mastery changes of a completed attempt.
// It is idempotent per attempt: a redelivered event changes nothing and returns
// applied=false. Mastery and rewards are granted only for a user's first completion
// of a lesson, so replaying a lesson cannot farm progress.
func (s *LessonStore) ApplyLessonCompleted(ctx context.Context, event events.LessonCompleted) (outcome *LessonOutcome, applied bool, err error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("begin lesson completion: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Serialise completions of the same lesson by the same user.
	if _, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1, 0))`, event.UserID+":"+event.LessonID); err != nil {
		return nil, false, fmt.Errorf("lock lesson completion: %w", err)
	}

	var already bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM lesson_rewards WHERE attempt_id = $1)`, event.AttemptID).Scan(&already); err != nil {
		return nil, false, fmt.Errorf("check processed attempt: %w", err)
	}
	if already {
		return nil, false, nil
	}

	var priorCompletions int
	if err := tx.QueryRow(ctx,
		`SELECT count(*) FROM lesson_rewards WHERE user_id = $1 AND lesson_id = $2`,
		event.UserID, event.LessonID).Scan(&priorCompletions); err != nil {
		return nil, false, fmt.Errorf("count prior completions: %w", err)
	}
	first := priorCompletions == 0

	accuracy := 0.0
	if event.GradedSteps > 0 {
		accuracy = float64(event.FirstTryCorrect) / float64(event.GradedSteps)
	}

	out := &LessonOutcome{AttemptID: event.AttemptID, LessonID: event.LessonID, FirstCompletion: first}

	if first {
		out.XPAwarded = config.LessonXP(event.Kind, event.Difficulty, accuracy)
		out.GemsAwarded = config.LessonGems(event.Difficulty, accuracy)
		for _, sk := range event.Skills {
			change, err := applyMastery(ctx, tx, event, sk, accuracy)
			if err != nil {
				return nil, false, err
			}
			out.Mastery = append(out.Mastery, *change)
		}
	}

	if err := applyRewards(ctx, tx, event, out); err != nil {
		return nil, false, err
	}

	if err := tx.QueryRow(ctx, `
		INSERT INTO lesson_rewards (attempt_id, user_id, lesson_id, first_completion, xp_awarded, gems_awarded)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at`,
		event.AttemptID, event.UserID, event.LessonID, first, out.XPAwarded, out.GemsAwarded).Scan(&out.CreatedAt); err != nil {
		return nil, false, fmt.Errorf("insert lesson reward: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, false, fmt.Errorf("commit lesson completion: %w", err)
	}
	return out, true, nil
}

func applyMastery(ctx context.Context, tx pgx.Tx, event events.LessonCompleted, sk events.LessonSkillWeight, accuracy float64) (*MasteryChange, error) {
	if _, err := tx.Exec(ctx,
		`INSERT INTO user_skill_scores (user_id, skill_id) VALUES ($1, $2) ON CONFLICT (user_id, skill_id) DO NOTHING`,
		event.UserID, sk.SkillID); err != nil {
		return nil, fmt.Errorf("ensure skill score: %w", err)
	}

	var before int
	if err := tx.QueryRow(ctx,
		`SELECT mastery FROM user_skill_scores WHERE user_id = $1 AND skill_id = $2 FOR UPDATE`,
		event.UserID, sk.SkillID).Scan(&before); err != nil {
		return nil, fmt.Errorf("lock skill score: %w", err)
	}

	delta := config.LessonMasteryDelta(before, event.Kind, event.Difficulty, accuracy, sk.Weight)
	after := before + delta

	if _, err := tx.Exec(ctx, `
		UPDATE user_skill_scores
		SET mastery = $3, attempts = attempts + 1, last_practiced_at = $4, source = 'live'
		WHERE user_id = $1 AND skill_id = $2`,
		event.UserID, sk.SkillID, after, event.CompletedAt); err != nil {
		return nil, fmt.Errorf("update skill score: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO mastery_ledger (user_id, skill_id, attempt_id, lesson_id, delta, mastery_before, mastery_after, accuracy)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		event.UserID, sk.SkillID, event.AttemptID, event.LessonID, delta, before, after, accuracy); err != nil {
		return nil, fmt.Errorf("insert mastery ledger: %w", err)
	}

	change := &MasteryChange{SkillID: sk.SkillID, SkillSlug: sk.SkillSlug, Before: before, After: after, Delta: delta, Accuracy: accuracy}
	if err := tx.QueryRow(ctx, `
		SELECT s.name, COALESCE(t.slug, ''), COALESCE(t.name, '')
		FROM skills s
		JOIN skill_categories c ON c.id = s.category_id
		LEFT JOIN topics t ON t.id = c.topic_id
		WHERE s.id = $1`, sk.SkillID).Scan(&change.SkillName, &change.TopicSlug, &change.TopicName); err != nil {
		return nil, fmt.Errorf("load skill names: %w", err)
	}
	return change, nil
}

func applyRewards(ctx context.Context, tx pgx.Tx, event events.LessonCompleted, out *LessonOutcome) error {
	var totalXP, gems int
	err := tx.QueryRow(ctx,
		`SELECT total_xp, gem_balance FROM user_progress WHERE user_id = $1 FOR UPDATE`, event.UserID).Scan(&totalXP, &gems)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("lock user progress: %w", err)
	}

	out.LevelBefore = config.CurrentLevel(totalXP)
	totalXP += out.XPAwarded
	gems += out.GemsAwarded
	out.LevelAfter = config.CurrentLevel(totalXP)
	out.TotalXP = totalXP
	out.TotalGems = gems

	if out.XPAwarded == 0 && out.GemsAwarded == 0 {
		return nil
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO user_progress (user_id, total_xp, current_level, gem_balance)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			total_xp = EXCLUDED.total_xp,
			current_level = EXCLUDED.current_level,
			gem_balance = EXCLUDED.gem_balance`,
		event.UserID, totalXP, out.LevelAfter, gems); err != nil {
		return fmt.Errorf("upsert user progress: %w", err)
	}
	if out.XPAwarded > 0 {
		if _, err := tx.Exec(ctx,
			`INSERT INTO xp_ledger (user_id, amount, reason, source_event_id) VALUES ($1, $2, 'lesson_completed', $3)`,
			event.UserID, out.XPAwarded, event.EventID); err != nil {
			return fmt.Errorf("insert xp ledger: %w", err)
		}
	}
	if out.GemsAwarded > 0 {
		if _, err := tx.Exec(ctx,
			`INSERT INTO gem_ledger (user_id, amount, reason, source_event_id) VALUES ($1, $2, 'lesson_completed', $3)`,
			event.UserID, out.GemsAwarded, event.EventID); err != nil {
			return fmt.Errorf("insert gem ledger: %w", err)
		}
	}
	return nil
}

// GetAttemptRewards returns what one of the user's attempts earned, or nil if
// Progress has not processed it (yet).
func (s *LessonStore) GetAttemptRewards(ctx context.Context, userID, attemptID string) (*LessonOutcome, error) {
	out := &LessonOutcome{AttemptID: attemptID}
	err := s.pool.QueryRow(ctx, `
		SELECT lesson_id, first_completion, xp_awarded, gems_awarded, created_at
		FROM lesson_rewards WHERE attempt_id = $1 AND user_id = $2`,
		attemptID, userID).Scan(&out.LessonID, &out.FirstCompletion, &out.XPAwarded, &out.GemsAwarded, &out.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get lesson reward: %w", err)
	}

	rows, err := s.pool.Query(ctx, `
		SELECT m.skill_id, s.slug, s.name, COALESCE(t.slug, ''), COALESCE(t.name, ''),
		       m.mastery_before, m.mastery_after, m.delta, m.accuracy::float8
		FROM mastery_ledger m
		JOIN skills s ON s.id = m.skill_id
		JOIN skill_categories c ON c.id = s.category_id
		LEFT JOIN topics t ON t.id = c.topic_id
		WHERE m.attempt_id = $1 AND m.user_id = $2
		ORDER BY s.slug`, attemptID, userID)
	if err != nil {
		return nil, fmt.Errorf("list mastery changes: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var c MasteryChange
		if err := rows.Scan(&c.SkillID, &c.SkillSlug, &c.SkillName, &c.TopicSlug, &c.TopicName, &c.Before, &c.After, &c.Delta, &c.Accuracy); err != nil {
			return nil, fmt.Errorf("scan mastery change: %w", err)
		}
		out.Mastery = append(out.Mastery, c)
	}
	return out, rows.Err()
}
