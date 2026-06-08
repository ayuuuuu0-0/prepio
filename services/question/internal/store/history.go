package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AnswerRecord is a row from user_question_history.
type AnswerRecord struct {
	ID          string
	UserID      string
	QuestionID  string
	Correct     bool
	Score       int
	SubmittedAt time.Time
	SessionID   string
}

// CompanyPerformance aggregates answer stats per company tag.
type CompanyPerformance struct {
	Company  string
	Answered int
	Correct  int
	ScoreAvg int
}

// HistoryStore handles user_question_history queries.
type HistoryStore struct {
	pool *pgxpool.Pool
}

// NewHistoryStore creates a HistoryStore.
func NewHistoryStore(pool *pgxpool.Pool) *HistoryStore {
	return &HistoryStore{pool: pool}
}

// Insert records an answer submission with an evaluation score.
func (s *HistoryStore) Insert(ctx context.Context, userID, questionID, sessionID string, correct bool, score int, submittedAt time.Time) error {
	const q = `
		INSERT INTO user_question_history (user_id, question_id, correct, score, submitted_at, session_id)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.pool.Exec(ctx, q, userID, questionID, correct, score, submittedAt, sessionID)
	if err != nil {
		return fmt.Errorf("insert answer history: %w", err)
	}
	return nil
}

// ExistsForSession reports whether the user already submitted this question in the session.
func (s *HistoryStore) ExistsForSession(ctx context.Context, userID, questionID, sessionID string) (bool, error) {
	const q = `
		SELECT 1 FROM user_question_history
		WHERE user_id = $1 AND question_id = $2 AND session_id = $3`

	var one int
	err := s.pool.QueryRow(ctx, q, userID, questionID, sessionID).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check answer exists: %w", err)
	}
	return true, nil
}

// ListBySession returns answer records for a daily paper session.
func (s *HistoryStore) ListBySession(ctx context.Context, userID, sessionID string) ([]AnswerRecord, error) {
	const q = `
		SELECT id, user_id, question_id, correct, score, submitted_at, session_id
		FROM user_question_history
		WHERE user_id = $1 AND session_id = $2
		ORDER BY submitted_at ASC`

	rows, err := s.pool.Query(ctx, q, userID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("list session history: %w", err)
	}
	defer rows.Close()

	records := make([]AnswerRecord, 0)
	for rows.Next() {
		var rec AnswerRecord
		if err := rows.Scan(&rec.ID, &rec.UserID, &rec.QuestionID, &rec.Correct, &rec.Score, &rec.SubmittedAt, &rec.SessionID); err != nil {
			return nil, fmt.Errorf("scan session history: %w", err)
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

// HasAnswerToday reports whether the user submitted any answer today (UTC date).
func (s *HistoryStore) HasAnswerToday(ctx context.Context, userID string) (bool, error) {
	const q = `
		SELECT 1 FROM user_question_history
		WHERE user_id = $1 AND submitted_at::date = CURRENT_DATE
		LIMIT 1`

	var one int
	err := s.pool.QueryRow(ctx, q, userID).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check answer today: %w", err)
	}
	return true, nil
}

// AvgScoreByUser returns the average evaluation score across all answers.
func (s *HistoryStore) AvgScoreByUser(ctx context.Context, userID string) (int, error) {
	const q = `SELECT COALESCE(AVG(score), 0)::int FROM user_question_history WHERE user_id = $1`
	var avg int
	err := s.pool.QueryRow(ctx, q, userID).Scan(&avg)
	if err != nil {
		return 0, fmt.Errorf("avg score: %w", err)
	}
	return avg, nil
}

// CompanyPerformanceByUser aggregates correctness by company tag.
func (s *HistoryStore) CompanyPerformanceByUser(ctx context.Context, userID string) ([]CompanyPerformance, error) {
	const q = `
		SELECT qt.company,
		       COUNT(*)::int AS answered,
		       SUM(CASE WHEN h.correct THEN 1 ELSE 0 END)::int AS correct,
		       COALESCE(AVG(h.score), 0)::int AS score_avg
		FROM user_question_history h
		JOIN question_tags qt ON qt.question_id = h.question_id
		WHERE h.user_id = $1
		GROUP BY qt.company
		ORDER BY qt.company`

	rows, err := s.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("company performance: %w", err)
	}
	defer rows.Close()

	stats := make([]CompanyPerformance, 0)
	for rows.Next() {
		var row CompanyPerformance
		if err := rows.Scan(&row.Company, &row.Answered, &row.Correct, &row.ScoreAvg); err != nil {
			return nil, fmt.Errorf("scan company performance: %w", err)
		}
		stats = append(stats, row)
	}
	return stats, rows.Err()
}
