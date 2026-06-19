package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/config"
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

// CompanySkillWeight is a weighted skill importance for a company profile.
type CompanySkillWeight struct {
	Company   string
	SkillID   string
	SkillSlug string
	SkillName string
	Weight    int
}

// QuestionSkillContribution holds data needed to update mastery from an answer.
type QuestionSkillContribution struct {
	SkillID         string
	SkillWeight     float64
	ReadinessWeight float64
	Difficulty      string
}

// ReadinessStore handles skill mastery and company weight queries.
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

// GetUserSkillScore returns mastery for one user/skill pair.
func (s *ReadinessStore) GetUserSkillScore(ctx context.Context, userID, skillID string) (*UserSkillScore, error) {
	const q = `
		SELECT uss.user_id, uss.skill_id, sk.slug, sk.name,
		       uss.mastery, uss.attempts, uss.last_practiced_at, uss.source
		FROM user_skill_scores uss
		JOIN skills sk ON sk.id = uss.skill_id
		WHERE uss.user_id = $1 AND uss.skill_id = $2`

	var row UserSkillScore
	err := s.pool.QueryRow(ctx, q, userID, skillID).Scan(
		&row.UserID, &row.SkillID, &row.SkillSlug, &row.SkillName,
		&row.Mastery, &row.Attempts, &row.LastPracticedAt, &row.Source,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user skill score: %w", err)
	}
	return &row, nil
}

// UpsertUserSkillScore inserts or updates a user's skill mastery.
func (s *ReadinessStore) UpsertUserSkillScore(
	ctx context.Context,
	userID, skillID string,
	mastery, attempts int,
	practicedAt time.Time,
	source string,
) error {
	if len(userID) == 0 || len(skillID) == 0 {
		return fmt.Errorf("user id and skill id are required")
	}
	if len(source) == 0 {
		source = config.ReadinessSourceLive
	}

	const q = `
		INSERT INTO user_skill_scores (user_id, skill_id, mastery, attempts, last_practiced_at, source)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, skill_id) DO UPDATE SET
			mastery = EXCLUDED.mastery,
			attempts = EXCLUDED.attempts,
			last_practiced_at = EXCLUDED.last_practiced_at,
			source = EXCLUDED.source,
			updated_at = now()`

	_, err := s.pool.Exec(ctx, q, userID, skillID, mastery, attempts, practicedAt, source)
	if err != nil {
		return fmt.Errorf("upsert user skill score: %w", err)
	}
	return nil
}

// ListCompanySkillWeights returns weighted skills for a company profile.
func (s *ReadinessStore) ListCompanySkillWeights(ctx context.Context, company string) ([]CompanySkillWeight, error) {
	if len(company) == 0 {
		return nil, fmt.Errorf("company is required")
	}

	const q = `
		SELECT csw.company, csw.skill_id, sk.slug, sk.name, csw.weight
		FROM company_skill_weights csw
		JOIN skills sk ON sk.id = csw.skill_id
		WHERE csw.company = $1
		ORDER BY csw.weight DESC, sk.name`

	rows, err := s.pool.Query(ctx, q, company)
	if err != nil {
		return nil, fmt.Errorf("list company skill weights: %w", err)
	}
	defer rows.Close()

	weights := make([]CompanySkillWeight, 0)
	for rows.Next() {
		var row CompanySkillWeight
		if err := rows.Scan(&row.Company, &row.SkillID, &row.SkillSlug, &row.SkillName, &row.Weight); err != nil {
			return nil, fmt.Errorf("scan company skill weight: %w", err)
		}
		weights = append(weights, row)
	}
	return weights, rows.Err()
}

// ListUserTargetCompanies returns onboarding target companies for a user.
func (s *ReadinessStore) ListUserTargetCompanies(ctx context.Context, userID string) ([]string, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("user id is required")
	}

	const q = `SELECT company FROM user_targets WHERE user_id = $1 ORDER BY company`
	rows, err := s.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list user targets: %w", err)
	}
	defer rows.Close()

	companies := make([]string, 0)
	for rows.Next() {
		var company string
		if err := rows.Scan(&company); err != nil {
			return nil, fmt.Errorf("scan target company: %w", err)
		}
		companies = append(companies, company)
	}
	return companies, rows.Err()
}

// ListQuestionSkillContributions returns skill mappings and weights for a question.
func (s *ReadinessStore) ListQuestionSkillContributions(ctx context.Context, questionID string) ([]QuestionSkillContribution, error) {
	if len(questionID) == 0 {
		return nil, fmt.Errorf("question id is required")
	}

	const q = `
		SELECT qs.skill_id, qs.weight::float8, q.readiness_weight::float8, q.difficulty
		FROM question_skills qs
		JOIN questions q ON q.id = qs.question_id
		WHERE qs.question_id = $1`

	rows, err := s.pool.Query(ctx, q, questionID)
	if err != nil {
		return nil, fmt.Errorf("list question skill contributions: %w", err)
	}
	defer rows.Close()

	contributions := make([]QuestionSkillContribution, 0)
	for rows.Next() {
		var row QuestionSkillContribution
		if err := rows.Scan(&row.SkillID, &row.SkillWeight, &row.ReadinessWeight, &row.Difficulty); err != nil {
			return nil, fmt.Errorf("scan question skill contribution: %w", err)
		}
		contributions = append(contributions, row)
	}
	return contributions, rows.Err()
}
