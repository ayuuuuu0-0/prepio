package store

import (
	"context"
	"encoding/json"
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

// QuestionSkillMapping links a question to a skill and subskill with weight.
type QuestionSkillMapping struct {
	QuestionID   string
	SkillID      string
	SkillSlug    string
	SubskillID   string
	SubskillSlug string
	Weight       float64
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

// ListQuestionSkills returns skill mappings for a question.
func (s *SkillStore) ListQuestionSkills(ctx context.Context, questionID string) ([]QuestionSkillMapping, error) {
	if len(questionID) == 0 {
		return nil, fmt.Errorf("question id is required")
	}

	const q = `
		SELECT qs.question_id, qs.skill_id, sk.slug, qs.subskill_id, ss.slug, qs.weight
		FROM question_skills qs
		JOIN skills sk ON sk.id = qs.skill_id
		JOIN subskills ss ON ss.id = qs.subskill_id
		WHERE qs.question_id = $1
		ORDER BY qs.weight DESC`

	rows, err := s.pool.Query(ctx, q, questionID)
	if err != nil {
		return nil, fmt.Errorf("list question skills: %w", err)
	}
	defer rows.Close()

	var mappings []QuestionSkillMapping
	for rows.Next() {
		var mapping QuestionSkillMapping
		if err := rows.Scan(
			&mapping.QuestionID, &mapping.SkillID, &mapping.SkillSlug,
			&mapping.SubskillID, &mapping.SubskillSlug, &mapping.Weight,
		); err != nil {
			return nil, fmt.Errorf("scan question skill: %w", err)
		}
		mappings = append(mappings, mapping)
	}
	return mappings, rows.Err()
}

// QuestionHint is a structured hint on a question.
type QuestionHint struct {
	Order int    `json:"order"`
	Text  string `json:"text"`
}

// QuestionContentMetadata holds extended question fields from the A3 schema upgrade.
type QuestionContentMetadata struct {
	EvaluationType  string
	Explanation     string
	Hints           []QuestionHint
	Solution        string
	ReadinessWeight float64
	EstimatedTime   int
}

// GetQuestionContentMetadata returns extended content fields for a question.
func (s *SkillStore) GetQuestionContentMetadata(ctx context.Context, questionID string) (*QuestionContentMetadata, error) {
	if len(questionID) == 0 {
		return nil, fmt.Errorf("question id is required")
	}

	const q = `
		SELECT COALESCE(evaluation_type, ''), COALESCE(explanation, ''),
		       hints, COALESCE(solution, ''), readiness_weight, estimated_time
		FROM questions
		WHERE id = $1`

	var meta QuestionContentMetadata
	var hintsJSON []byte
	err := s.pool.QueryRow(ctx, q, questionID).Scan(
		&meta.EvaluationType, &meta.Explanation, &hintsJSON,
		&meta.Solution, &meta.ReadinessWeight, &meta.EstimatedTime,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get question content metadata: %w", err)
	}

	if len(hintsJSON) > 0 {
		if err := json.Unmarshal(hintsJSON, &meta.Hints); err != nil {
			return nil, fmt.Errorf("unmarshal hints: %w", err)
		}
	}
	return &meta, nil
}
