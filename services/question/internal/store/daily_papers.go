package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DailyPaper is a generated daily paper for a user.
type DailyPaper struct {
	ID        string
	UserID    string
	SessionID string
	PaperDate time.Time
}

// DailyPaperStore handles daily_papers table queries.
type DailyPaperStore struct {
	pool *pgxpool.Pool
}

// NewDailyPaperStore creates a DailyPaperStore.
func NewDailyPaperStore(pool *pgxpool.Pool) *DailyPaperStore {
	return &DailyPaperStore{pool: pool}
}

// GetByUserAndDate returns an existing paper for the user and date.
func (s *DailyPaperStore) GetByUserAndDate(ctx context.Context, userID string, paperDate time.Time) (*DailyPaper, []Question, error) {
	const paperQ = `
		SELECT id, user_id, session_id, paper_date
		FROM daily_papers
		WHERE user_id = $1 AND paper_date = $2`

	var paper DailyPaper
	err := s.pool.QueryRow(ctx, paperQ, userID, paperDate.Format("2006-01-02")).Scan(
		&paper.ID, &paper.UserID, &paper.SessionID, &paper.PaperDate,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("get daily paper: %w", err)
	}

	const questionsQ = `
		SELECT q.id, q.body, q.round_type, q.difficulty, q.answer_guide, q.status, q.is_weekend
		FROM daily_paper_questions dpq
		JOIN questions q ON q.id = dpq.question_id
		WHERE dpq.daily_paper_id = $1
		ORDER BY dpq.position`

	rows, err := s.pool.Query(ctx, questionsQ, paper.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("get daily paper questions: %w", err)
	}
	defer rows.Close()

	questionStore := NewQuestionStore(s.pool)
	var questions []Question
	for rows.Next() {
		var question Question
		if err := rows.Scan(
			&question.ID, &question.Body, &question.RoundType, &question.Difficulty,
			&question.AnswerGuide, &question.Status, &question.IsWeekend,
		); err != nil {
			return nil, nil, fmt.Errorf("scan daily question: %w", err)
		}
		tags, err := questionStore.loadTags(ctx, question.ID)
		if err != nil {
			return nil, nil, err
		}
		question.CompanyTags = tags
		questions = append(questions, question)
	}
	return &paper, questions, rows.Err()
}

// Create inserts a daily paper and its questions.
func (s *DailyPaperStore) Create(ctx context.Context, userID, sessionID string, paperDate time.Time, questions []Question) (*DailyPaper, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	const paperQ = `
		INSERT INTO daily_papers (user_id, session_id, paper_date)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, session_id, paper_date`

	var paper DailyPaper
	err = tx.QueryRow(ctx, paperQ, userID, sessionID, paperDate.Format("2006-01-02")).Scan(
		&paper.ID, &paper.UserID, &paper.SessionID, &paper.PaperDate,
	)
	if err != nil {
		return nil, fmt.Errorf("insert daily paper: %w", err)
	}

	const linkQ = `INSERT INTO daily_paper_questions (daily_paper_id, question_id, position) VALUES ($1, $2, $3)`
	for i, question := range questions {
		if _, err := tx.Exec(ctx, linkQ, paper.ID, question.ID, i+1); err != nil {
			return nil, fmt.Errorf("link question: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit daily paper: %w", err)
	}
	return &paper, nil
}

// QuestionInSession reports whether a question belongs to the session's daily paper.
func (s *DailyPaperStore) QuestionInSession(ctx context.Context, sessionID, questionID string) (bool, error) {
	const q = `
		SELECT 1
		FROM daily_papers dp
		JOIN daily_paper_questions dpq ON dpq.daily_paper_id = dp.id
		WHERE dp.session_id = $1 AND dpq.question_id = $2`

	var one int
	err := s.pool.QueryRow(ctx, q, sessionID, questionID).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check question in session: %w", err)
	}
	return true, nil
}
