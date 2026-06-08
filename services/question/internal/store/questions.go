package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Question is a row from the questions table.
type Question struct {
	ID          string
	Body        string
	RoundType   string
	Difficulty  string
	AnswerGuide string
	Status      string
	IsWeekend   bool
	CompanyTags []string
}

// QuestionStore handles question bank queries.
type QuestionStore struct {
	pool *pgxpool.Pool
}

// NewQuestionStore creates a QuestionStore.
func NewQuestionStore(pool *pgxpool.Pool) *QuestionStore {
	return &QuestionStore{pool: pool}
}

// SelectUnseenByDifficulty returns approved questions the user has not answered.
func (s *QuestionStore) SelectUnseenByDifficulty(ctx context.Context, userID, difficulty string, limit int, weekendOnly bool) ([]Question, error) {
	const q = `
		SELECT q.id, q.body, q.round_type, q.difficulty, q.answer_guide, q.status, q.is_weekend
		FROM questions q
		WHERE q.status = 'approved'
		  AND q.difficulty = $2
		  AND q.is_weekend = $4
		  AND NOT EXISTS (
		      SELECT 1 FROM user_question_history h
		      WHERE h.user_id = $1 AND h.question_id = q.id
		  )
		ORDER BY random()
		LIMIT $3`

	rows, err := s.pool.Query(ctx, q, userID, difficulty, limit, weekendOnly)
	if err != nil {
		return nil, fmt.Errorf("select unseen questions: %w", err)
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var question Question
		if err := rows.Scan(
			&question.ID, &question.Body, &question.RoundType, &question.Difficulty,
			&question.AnswerGuide, &question.Status, &question.IsWeekend,
		); err != nil {
			return nil, fmt.Errorf("scan question: %w", err)
		}
		tags, err := s.loadTags(ctx, question.ID)
		if err != nil {
			return nil, err
		}
		question.CompanyTags = tags
		questions = append(questions, question)
	}
	return questions, rows.Err()
}

// SelectRandomApproved returns random approved questions as a fallback.
func (s *QuestionStore) SelectRandomApproved(ctx context.Context, userID string, limit int, weekendOnly bool) ([]Question, error) {
	const q = `
		SELECT q.id, q.body, q.round_type, q.difficulty, q.answer_guide, q.status, q.is_weekend
		FROM questions q
		WHERE q.status = 'approved'
		  AND q.is_weekend = $3
		  AND NOT EXISTS (
		      SELECT 1 FROM user_question_history h
		      WHERE h.user_id = $1 AND h.question_id = q.id
		  )
		ORDER BY random()
		LIMIT $2`

	rows, err := s.pool.Query(ctx, q, userID, limit, weekendOnly)
	if err != nil {
		return nil, fmt.Errorf("select random questions: %w", err)
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var question Question
		if err := rows.Scan(
			&question.ID, &question.Body, &question.RoundType, &question.Difficulty,
			&question.AnswerGuide, &question.Status, &question.IsWeekend,
		); err != nil {
			return nil, fmt.Errorf("scan question: %w", err)
		}
		tags, err := s.loadTags(ctx, question.ID)
		if err != nil {
			return nil, err
		}
		question.CompanyTags = tags
		questions = append(questions, question)
	}
	return questions, rows.Err()
}

// GetByID returns a question by ID.
func (s *QuestionStore) GetByID(ctx context.Context, id string) (*Question, error) {
	const q = `
		SELECT id, body, round_type, difficulty, answer_guide, status, is_weekend
		FROM questions WHERE id = $1`

	var question Question
	err := s.pool.QueryRow(ctx, q, id).Scan(
		&question.ID, &question.Body, &question.RoundType, &question.Difficulty,
		&question.AnswerGuide, &question.Status, &question.IsWeekend,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get question: %w", err)
	}
	tags, err := s.loadTags(ctx, question.ID)
	if err != nil {
		return nil, err
	}
	question.CompanyTags = tags
	return &question, nil
}

// ListCompanies returns distinct company tags.
func (s *QuestionStore) ListCompanies(ctx context.Context) ([]string, error) {
	const q = `SELECT DISTINCT company FROM question_tags ORDER BY company`
	rows, err := s.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list companies: %w", err)
	}
	defer rows.Close()

	var companies []string
	for rows.Next() {
		var company string
		if err := rows.Scan(&company); err != nil {
			return nil, fmt.Errorf("scan company: %w", err)
		}
		companies = append(companies, company)
	}
	return companies, rows.Err()
}

func (s *QuestionStore) loadTags(ctx context.Context, questionID string) ([]string, error) {
	const q = `SELECT company FROM question_tags WHERE question_id = $1 ORDER BY company`
	rows, err := s.pool.Query(ctx, q, questionID)
	if err != nil {
		return nil, fmt.Errorf("load tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}
