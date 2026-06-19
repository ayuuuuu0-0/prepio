package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// World is a journey world row.
type World struct {
	ID          string
	Slug        string
	Name        string
	Description string
	Theme       string
}

// JourneyNode is a node on a world path.
type JourneyNode struct {
	ID        string
	WorldID   string
	Slug      string
	Label     string
	NodeType  string
	SortOrder int
}

// JourneyStore handles worlds and journey nodes.
type JourneyStore struct {
	pool *pgxpool.Pool
}

// NewJourneyStore creates a JourneyStore.
func NewJourneyStore(pool *pgxpool.Pool) *JourneyStore {
	return &JourneyStore{pool: pool}
}

// GetWorldBySlug returns a world by slug.
func (s *JourneyStore) GetWorldBySlug(ctx context.Context, slug string) (*World, error) {
	const q = `SELECT id, slug, name, description, theme FROM worlds WHERE slug = $1`
	var w World
	err := s.pool.QueryRow(ctx, q, slug).Scan(&w.ID, &w.Slug, &w.Name, &w.Description, &w.Theme)
	if err != nil {
		return nil, fmt.Errorf("get world: %w", err)
	}
	return &w, nil
}

// ListNodesByWorld returns ordered nodes for a world.
func (s *JourneyStore) ListNodesByWorld(ctx context.Context, worldID string) ([]JourneyNode, error) {
	const q = `
		SELECT id, world_id, COALESCE(slug, ''), label, node_type, sort_order
		FROM journey_nodes
		WHERE world_id = $1
		ORDER BY sort_order ASC`

	rows, err := s.pool.Query(ctx, q, worldID)
	if err != nil {
		return nil, fmt.Errorf("list journey nodes: %w", err)
	}
	defer rows.Close()

	nodes := make([]JourneyNode, 0)
	for rows.Next() {
		var n JourneyNode
		if err := rows.Scan(&n.ID, &n.WorldID, &n.Slug, &n.Label, &n.NodeType, &n.SortOrder); err != nil {
			return nil, fmt.Errorf("scan journey node: %w", err)
		}
		nodes = append(nodes, n)
	}
	return nodes, rows.Err()
}

// GetNodeByID returns a journey node by ID.
func (s *JourneyStore) GetNodeByID(ctx context.Context, nodeID string) (*JourneyNode, error) {
	if len(nodeID) == 0 {
		return nil, fmt.Errorf("node id is required")
	}

	const q = `
		SELECT id, world_id, COALESCE(slug, ''), label, node_type, sort_order
		FROM journey_nodes
		WHERE id = $1`

	var n JourneyNode
	err := s.pool.QueryRow(ctx, q, nodeID).Scan(
		&n.ID, &n.WorldID, &n.Slug, &n.Label, &n.NodeType, &n.SortOrder,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get journey node: %w", err)
	}
	return &n, nil
}

// UpsertProgress marks a node complete for a user.
func (s *JourneyStore) UpsertProgress(ctx context.Context, userID, nodeID, status string) error {
	const q = `
		INSERT INTO user_journey_progress (user_id, node_id, status, completed_at)
		VALUES ($1, $2, $3, CASE WHEN $3 = 'done' THEN now() ELSE NULL END)
		ON CONFLICT (user_id, node_id) DO UPDATE
		SET status = EXCLUDED.status,
		    completed_at = CASE WHEN EXCLUDED.status = 'done' THEN now() ELSE user_journey_progress.completed_at END`

	_, err := s.pool.Exec(ctx, q, userID, nodeID, status)
	if err != nil {
		return fmt.Errorf("upsert journey progress: %w", err)
	}
	return nil
}
