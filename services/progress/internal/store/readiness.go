package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ReadinessStore handles skill mastery queries.
type ReadinessStore struct {
	pool *pgxpool.Pool
}

// NewReadinessStore creates a ReadinessStore.
func NewReadinessStore(pool *pgxpool.Pool) *ReadinessStore {
	return &ReadinessStore{pool: pool}
}

// TopicSkillRow is one skill of a topic with the user's mastery, if any.
type TopicSkillRow struct {
	TopicSlug        string
	TopicName        string
	TopicDescription string
	SkillSlug        string
	SkillName        string
	Mastery          *int
	Attempts         int
}

// ListTopicSkills returns every skill that belongs to a topic, in topic and catalog order,
// with the user's mastery where they have practiced it.
func (s *ReadinessStore) ListTopicSkills(ctx context.Context, userID string) ([]TopicSkillRow, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("user id is required")
	}

	const q = `
		SELECT t.slug, t.name, t.description, sk.slug, sk.name, u.mastery, COALESCE(u.attempts, 0)
		FROM topics t
		JOIN skill_categories c ON c.topic_id = t.id
		JOIN skills sk ON sk.category_id = c.id
		LEFT JOIN user_skill_scores u ON u.skill_id = sk.id AND u.user_id = $1
		ORDER BY t.sort_order, t.name, c.sort_order, sk.sort_order, sk.name`

	rows, err := s.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list topic skills: %w", err)
	}
	defer rows.Close()

	var out []TopicSkillRow
	for rows.Next() {
		var r TopicSkillRow
		if err := rows.Scan(&r.TopicSlug, &r.TopicName, &r.TopicDescription, &r.SkillSlug, &r.SkillName, &r.Mastery, &r.Attempts); err != nil {
			return nil, fmt.Errorf("scan topic skill: %w", err)
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
