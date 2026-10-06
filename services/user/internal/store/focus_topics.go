package store

import (
	"context"
	"errors"
	"fmt"
)

// ErrUnknownTopic is returned when a focus topic slug is not in the topic catalog.
var ErrUnknownTopic = errors.New("unknown topic")

// SetFocusTopics replaces the user's focus topics. slugs are ordered by priority and must
// all exist in the topic catalog; the whole replacement is atomic.
func (s *UserStore) SetFocusTopics(ctx context.Context, userID string, slugs []string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin focus topics: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, `DELETE FROM user_focus_topics WHERE user_id = $1`, userID); err != nil {
		return fmt.Errorf("clear focus topics: %w", err)
	}
	for i, slug := range slugs {
		tag, err := tx.Exec(ctx, `
			INSERT INTO user_focus_topics (user_id, topic_id, position)
			SELECT $1, id, $3 FROM topics WHERE slug = $2`, userID, slug, i+1)
		if err != nil {
			return fmt.Errorf("insert focus topic %q: %w", slug, err)
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("%w: %s", ErrUnknownTopic, slug)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit focus topics: %w", err)
	}
	return nil
}

// ListFocusTopics returns the user's focus topic slugs in priority order.
func (s *UserStore) ListFocusTopics(ctx context.Context, userID string) ([]string, error) {
	const q = `
		SELECT t.slug
		FROM user_focus_topics f
		JOIN topics t ON t.id = f.topic_id
		WHERE f.user_id = $1
		ORDER BY f.position`
	rows, err := s.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list focus topics: %w", err)
	}
	defer rows.Close()

	slugs := []string{}
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, fmt.Errorf("scan focus topic: %w", err)
		}
		slugs = append(slugs, slug)
	}
	return slugs, rows.Err()
}
