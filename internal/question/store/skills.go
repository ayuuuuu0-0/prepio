package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SkillCategory is a row from skill_categories.
type SkillCategory struct {
	ID        string
	Slug      string
	Name      string
	SortOrder int
}

// Skill is a row from skills.
type Skill struct {
	ID          string
	CategoryID  string
	Slug        string
	Name        string
	Description string
	SortOrder   int
}

// Subskill is a row from subskills.
type Subskill struct {
	ID        string
	SkillID   string
	Slug      string
	Name      string
	SortOrder int
}

// SkillStore handles skill graph queries.
type SkillStore struct {
	pool *pgxpool.Pool
}

// NewSkillStore creates a SkillStore.
func NewSkillStore(pool *pgxpool.Pool) *SkillStore {
	return &SkillStore{pool: pool}
}

// ListCategories returns all skill categories ordered by sort_order.
func (s *SkillStore) ListCategories(ctx context.Context) ([]SkillCategory, error) {
	const q = `
		SELECT id, slug, name, sort_order
		FROM skill_categories
		ORDER BY sort_order, name`

	rows, err := s.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list skill categories: %w", err)
	}
	defer rows.Close()

	var categories []SkillCategory
	for rows.Next() {
		var category SkillCategory
		if err := rows.Scan(&category.ID, &category.Slug, &category.Name, &category.SortOrder); err != nil {
			return nil, fmt.Errorf("scan skill category: %w", err)
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

// ListSkillsByCategory returns skills for a category ordered by sort_order.
func (s *SkillStore) ListSkillsByCategory(ctx context.Context, categoryID string) ([]Skill, error) {
	const q = `
		SELECT id, category_id, slug, name, COALESCE(description, ''), sort_order
		FROM skills
		WHERE category_id = $1
		ORDER BY sort_order, name`

	rows, err := s.pool.Query(ctx, q, categoryID)
	if err != nil {
		return nil, fmt.Errorf("list skills by category: %w", err)
	}
	defer rows.Close()

	var skills []Skill
	for rows.Next() {
		var skill Skill
		if err := rows.Scan(
			&skill.ID, &skill.CategoryID, &skill.Slug, &skill.Name,
			&skill.Description, &skill.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("scan skill: %w", err)
		}
		skills = append(skills, skill)
	}
	return skills, rows.Err()
}

// GetSkillBySlug returns a skill by slug.
func (s *SkillStore) GetSkillBySlug(ctx context.Context, slug string) (*Skill, error) {
	if len(slug) == 0 {
		return nil, fmt.Errorf("skill slug is required")
	}

	const q = `
		SELECT id, category_id, slug, name, COALESCE(description, ''), sort_order
		FROM skills
		WHERE slug = $1`

	var skill Skill
	err := s.pool.QueryRow(ctx, q, slug).Scan(
		&skill.ID, &skill.CategoryID, &skill.Slug, &skill.Name,
		&skill.Description, &skill.SortOrder,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get skill by slug: %w", err)
	}
	return &skill, nil
}

// ListSubskillsBySkillID returns subskills for a skill ordered by sort_order.
func (s *SkillStore) ListSubskillsBySkillID(ctx context.Context, skillID string) ([]Subskill, error) {
	const q = `
		SELECT id, skill_id, slug, name, sort_order
		FROM subskills
		WHERE skill_id = $1
		ORDER BY sort_order, name`

	rows, err := s.pool.Query(ctx, q, skillID)
	if err != nil {
		return nil, fmt.Errorf("list subskills: %w", err)
	}
	defer rows.Close()

	var subskills []Subskill
	for rows.Next() {
		var subskill Subskill
		if err := rows.Scan(
			&subskill.ID, &subskill.SkillID, &subskill.Slug,
			&subskill.Name, &subskill.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("scan subskill: %w", err)
		}
		subskills = append(subskills, subskill)
	}
	return subskills, rows.Err()
}
