package store

import (
	"context"
	"fmt"
)

// Topic is a user-facing grouping of skills; readiness is computed per topic.
type Topic struct {
	ID          string
	Slug        string
	Name        string
	Description string
	SortOrder   int
}

// ListTopics returns the topic catalog in display order.
func (s *SkillStore) ListTopics(ctx context.Context) ([]Topic, error) {
	const q = `SELECT id, slug, name, description, sort_order FROM topics ORDER BY sort_order, name`
	rows, err := s.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list topics: %w", err)
	}
	defer rows.Close()

	var topics []Topic
	for rows.Next() {
		var t Topic
		if err := rows.Scan(&t.ID, &t.Slug, &t.Name, &t.Description, &t.SortOrder); err != nil {
			return nil, fmt.Errorf("scan topic: %w", err)
		}
		topics = append(topics, t)
	}
	return topics, rows.Err()
}
