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
	SubmittedAt time.Time
	SessionID   string
}

// HistoryStore handles user_question_history queries.
type HistoryStore struct {
	pool *pgxpool.Pool
}

// NewHistoryStore creates a HistoryStore.
func NewHistoryStore(pool *pgxpool.Pool) *HistoryStore {
	return &HistoryStore{pool: pool}
}

// Insert records an answer submission.
func (s *HistoryStore) Insert(ctx context.Context, userID, questionID, sessionID string, correct bool, submittedAt time.Time) error {
	const q = `
		INSERT INTO user_question_history (user_id, question_id, correct, submitted_at, session_id)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := s.pool.Exec(ctx, q, userID, questionID, correct, submittedAt, sessionID)
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
