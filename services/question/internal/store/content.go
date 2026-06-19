package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// QuestionPool is a curated question set scoped to a skill.
type QuestionPool struct {
	ID          string
	SkillID     string
	SkillSlug   string
	SkillName   string
	Slug        string
	Name        string
	Description string
	SortOrder   int
}

// NodeSkillBinding links a journey node to a skill.
type NodeSkillBinding struct {
	NodeID    string
	SkillID   string
	SkillSlug string
	SkillName string
	IsPrimary bool
}

// NodePoolBinding links a journey node to a question pool.
type NodePoolBinding struct {
	NodeID             string
	PoolID             string
	PoolSlug           string
	PoolName           string
	SkillSlug          string
	SelectionStrategy  string
	QuestionsRequired  int
}

// NodeContent summarizes skills and pools attached to a journey node.
type NodeContent struct {
	NodeID string
	Skills []NodeSkillBinding
	Pools  []NodePoolBinding
}

// ContentStore handles question pools and node content bindings.
type ContentStore struct {
	pool *pgxpool.Pool
}

// NewContentStore creates a ContentStore.
func NewContentStore(pool *pgxpool.Pool) *ContentStore {
	return &ContentStore{pool: pool}
}

// GetPoolBySlug returns a question pool by slug.
func (s *ContentStore) GetPoolBySlug(ctx context.Context, slug string) (*QuestionPool, error) {
	if len(slug) == 0 {
		return nil, fmt.Errorf("pool slug is required")
	}

	const q = `
		SELECT p.id, p.skill_id, sk.slug, sk.name, p.slug, p.name, p.description, p.sort_order
		FROM question_pools p
		JOIN skills sk ON sk.id = p.skill_id
		WHERE p.slug = $1`

	var pool QuestionPool
	err := s.pool.QueryRow(ctx, q, slug).Scan(
		&pool.ID, &pool.SkillID, &pool.SkillSlug, &pool.SkillName,
		&pool.Slug, &pool.Name, &pool.Description, &pool.SortOrder,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get pool by slug: %w", err)
	}
	return &pool, nil
}

// ListPoolQuestionIDs returns approved question IDs in a pool ordered by sort_order.
func (s *ContentStore) ListPoolQuestionIDs(ctx context.Context, poolID string) ([]string, error) {
	if len(poolID) == 0 {
		return nil, fmt.Errorf("pool id is required")
	}

	const q = `
		SELECT pq.question_id
		FROM pool_questions pq
		JOIN questions q ON q.id = pq.question_id
		WHERE pq.pool_id = $1 AND q.status = 'approved'
		ORDER BY pq.sort_order, pq.question_id`

	rows, err := s.pool.Query(ctx, q, poolID)
	if err != nil {
		return nil, fmt.Errorf("list pool questions: %w", err)
	}
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan pool question id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// ListNodeSkills returns skills bound to a journey node.
func (s *ContentStore) ListNodeSkills(ctx context.Context, nodeID string) ([]NodeSkillBinding, error) {
	if len(nodeID) == 0 {
		return nil, fmt.Errorf("node id is required")
	}

	const q = `
		SELECT ns.node_id, ns.skill_id, sk.slug, sk.name, ns.is_primary
		FROM node_skills ns
		JOIN skills sk ON sk.id = ns.skill_id
		WHERE ns.node_id = $1
		ORDER BY ns.is_primary DESC, sk.name`

	rows, err := s.pool.Query(ctx, q, nodeID)
	if err != nil {
		return nil, fmt.Errorf("list node skills: %w", err)
	}
	defer rows.Close()

	bindings := make([]NodeSkillBinding, 0)
	for rows.Next() {
		var binding NodeSkillBinding
		if err := rows.Scan(
			&binding.NodeID, &binding.SkillID, &binding.SkillSlug,
			&binding.SkillName, &binding.IsPrimary,
		); err != nil {
			return nil, fmt.Errorf("scan node skill: %w", err)
		}
		bindings = append(bindings, binding)
	}
	return bindings, rows.Err()
}

// ListNodePools returns pools bound to a journey node.
func (s *ContentStore) ListNodePools(ctx context.Context, nodeID string) ([]NodePoolBinding, error) {
	if len(nodeID) == 0 {
		return nil, fmt.Errorf("node id is required")
	}

	const q = `
		SELECT np.node_id, np.pool_id, p.slug, p.name, sk.slug,
		       np.selection_strategy, np.questions_required
		FROM node_pools np
		JOIN question_pools p ON p.id = np.pool_id
		JOIN skills sk ON sk.id = p.skill_id
		WHERE np.node_id = $1
		ORDER BY p.sort_order, p.name`

	rows, err := s.pool.Query(ctx, q, nodeID)
	if err != nil {
		return nil, fmt.Errorf("list node pools: %w", err)
	}
	defer rows.Close()

	bindings := make([]NodePoolBinding, 0)
	for rows.Next() {
		var binding NodePoolBinding
		if err := rows.Scan(
			&binding.NodeID, &binding.PoolID, &binding.PoolSlug, &binding.PoolName,
			&binding.SkillSlug, &binding.SelectionStrategy, &binding.QuestionsRequired,
		); err != nil {
			return nil, fmt.Errorf("scan node pool: %w", err)
		}
		bindings = append(bindings, binding)
	}
	return bindings, rows.Err()
}

// GetNodeContent returns skill and pool bindings for a journey node.
func (s *ContentStore) GetNodeContent(ctx context.Context, nodeID string) (*NodeContent, error) {
	skills, err := s.ListNodeSkills(ctx, nodeID)
	if err != nil {
		return nil, err
	}
	pools, err := s.ListNodePools(ctx, nodeID)
	if err != nil {
		return nil, err
	}
	return &NodeContent{
		NodeID: nodeID,
		Skills: skills,
		Pools:  pools,
	}, nil
}

// ListPoolsBySkillID returns pools for a skill.
func (s *ContentStore) ListPoolsBySkillID(ctx context.Context, skillID string) ([]QuestionPool, error) {
	if len(skillID) == 0 {
		return nil, fmt.Errorf("skill id is required")
	}

	const q = `
		SELECT p.id, p.skill_id, sk.slug, sk.name, p.slug, p.name, p.description, p.sort_order
		FROM question_pools p
		JOIN skills sk ON sk.id = p.skill_id
		WHERE p.skill_id = $1
		ORDER BY p.sort_order, p.name`

	rows, err := s.pool.Query(ctx, q, skillID)
	if err != nil {
		return nil, fmt.Errorf("list pools by skill: %w", err)
	}
	defer rows.Close()

	pools := make([]QuestionPool, 0)
	for rows.Next() {
		var pool QuestionPool
		if err := rows.Scan(
			&pool.ID, &pool.SkillID, &pool.SkillSlug, &pool.SkillName,
			&pool.Slug, &pool.Name, &pool.Description, &pool.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("scan pool: %w", err)
		}
		pools = append(pools, pool)
	}
	return pools, rows.Err()
}
