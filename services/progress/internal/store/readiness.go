package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserSkillScore is a row from user_skill_scores joined with skill metadata.
type UserSkillScore struct {
	UserID          string
	SkillID         string
	SkillSlug       string
	SkillName       string
	Mastery         int
	Attempts        int
	LastPracticedAt *time.Time
	Source          string
}

// ReadinessStore handles skill mastery queries.
type ReadinessStore struct {
	pool *pgxpool.Pool
}

// NewReadinessStore creates a ReadinessStore.
func NewReadinessStore(pool *pgxpool.Pool) *ReadinessStore {
	return &ReadinessStore{pool: pool}
}

// ListUserSkillScores returns all skill mastery rows for a user.
func (s *ReadinessStore) ListUserSkillScores(ctx context.Context, userID string) ([]UserSkillScore, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("user id is required")
	}

	const q = `
		SELECT uss.user_id, uss.skill_id, sk.slug, sk.name,
		       uss.mastery, uss.attempts, uss.last_practiced_at, uss.source
		FROM user_skill_scores uss
		JOIN skills sk ON sk.id = uss.skill_id
		WHERE uss.user_id = $1
		ORDER BY sk.name`

	rows, err := s.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list user skill scores: %w", err)
	}
	defer rows.Close()

	scores := make([]UserSkillScore, 0)
	for rows.Next() {
		var row UserSkillScore
		if err := rows.Scan(
			&row.UserID, &row.SkillID, &row.SkillSlug, &row.SkillName,
			&row.Mastery, &row.Attempts, &row.LastPracticedAt, &row.Source,
		); err != nil {
			return nil, fmt.Errorf("scan user skill score: %w", err)
		}
		scores = append(scores, row)
	}
	return scores, rows.Err()
}
