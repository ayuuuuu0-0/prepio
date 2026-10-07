package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// uniqueViolation is the Postgres error code for unique_violation.
const uniqueViolation = "23505"

// Lesson is a row from lessons.
type Lesson struct {
	ID         string
	Slug       string
	NodeID     string
	Title      string
	Summary    []string
	Kind       string
	Difficulty string
	EstMinutes int
	Status     string
	Version    int
}

// LessonStep is a row from lesson_steps.
type LessonStep struct {
	ID       string
	Key      string
	Type     string
	Position int
	Payload  []byte
}

// LessonSkill is one skill weight of a lesson version.
type LessonSkill struct {
	SkillID   string
	SkillSlug string
	Weight    float64
}

// Attempt is a row from lesson_attempts.
type Attempt struct {
	ID                string
	UserID            string
	LessonID          string
	LessonVersion     int
	Status            string
	GradedSteps       int
	FirstTryCorrect   int
	TotalTries        int
	CompletionEventID string
	EventPublishedAt  *time.Time
	StartedAt         time.Time
	CompletedAt       *time.Time
}

// StepResult is a row from lesson_step_results joined to its step key.
type StepResult struct {
	StepID   string
	StepKey  string
	Try      int
	Answer   []byte
	Correct  bool
	Feedback []byte
}

// PathNode is one node of the learner's path with everything the journey needs.
type PathNode struct {
	WorldID          string
	WorldSlug        string
	WorldName        string
	WorldDescription string
	WorldTheme       string
	NodeID           string
	NodeSlug         string
	NodeLabel        string
	NodeType         string
	Requires         []string
	Lesson           Lesson
	Completed        bool
	InProgress       bool
}

// LessonStore handles lessons, attempts, and step results.
type LessonStore struct {
	pool *pgxpool.Pool
}

// NewLessonStore creates a LessonStore.
func NewLessonStore(pool *pgxpool.Pool) *LessonStore {
	return &LessonStore{pool: pool}
}

const lessonColumns = `l.id, l.slug, l.node_id, l.title, l.summary, l.kind, l.difficulty, l.est_minutes, l.status, l.version`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanLesson(row rowScanner) (*Lesson, error) {
	var l Lesson
	var summary []byte
	if err := row.Scan(&l.ID, &l.Slug, &l.NodeID, &l.Title, &summary, &l.Kind, &l.Difficulty, &l.EstMinutes, &l.Status, &l.Version); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(summary, &l.Summary); err != nil {
		return nil, fmt.Errorf("decode lesson summary: %w", err)
	}
	return &l, nil
}

// GetLesson returns a lesson by id regardless of status, or nil.
func (s *LessonStore) GetLesson(ctx context.Context, lessonID string) (*Lesson, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+lessonColumns+` FROM lessons l WHERE l.id = $1`, lessonID)
	l, err := scanLesson(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get lesson: %w", err)
	}
	return l, nil
}

// ListSteps returns the ordered steps of one lesson version.
func (s *LessonStore) ListSteps(ctx context.Context, lessonID string, version int) ([]LessonStep, error) {
	const q = `
		SELECT id, step_key, type, position, payload
		FROM lesson_steps
		WHERE lesson_id = $1 AND lesson_version = $2
		ORDER BY position`
	rows, err := s.pool.Query(ctx, q, lessonID, version)
	if err != nil {
		return nil, fmt.Errorf("list lesson steps: %w", err)
	}
	defer rows.Close()

	var steps []LessonStep
	for rows.Next() {
		var st LessonStep
		if err := rows.Scan(&st.ID, &st.Key, &st.Type, &st.Position, &st.Payload); err != nil {
			return nil, fmt.Errorf("scan lesson step: %w", err)
		}
		steps = append(steps, st)
	}
	return steps, rows.Err()
}

// ListSkills returns the skill weights of one lesson version.
func (s *LessonStore) ListSkills(ctx context.Context, lessonID string, version int) ([]LessonSkill, error) {
	const q = `
		SELECT ls.skill_id, sk.slug, ls.weight::float8
		FROM lesson_skills ls
		JOIN skills sk ON sk.id = ls.skill_id
		WHERE ls.lesson_id = $1 AND ls.lesson_version = $2
		ORDER BY sk.slug`
	rows, err := s.pool.Query(ctx, q, lessonID, version)
	if err != nil {
		return nil, fmt.Errorf("list lesson skills: %w", err)
	}
	defer rows.Close()

	var skills []LessonSkill
	for rows.Next() {
		var sk LessonSkill
		if err := rows.Scan(&sk.SkillID, &sk.SkillSlug, &sk.Weight); err != nil {
			return nil, fmt.Errorf("scan lesson skill: %w", err)
		}
		skills = append(skills, sk)
	}
	return skills, rows.Err()
}

// ListPath returns every published node (with its published lesson) in path order,
// annotated with the user's attempt state.
func (s *LessonStore) ListPath(ctx context.Context, userID string) ([]PathNode, error) {
	const q = `
		SELECT w.id, w.slug, w.name, w.description, w.theme,
		       n.id, n.slug, n.label, n.node_type,
		       COALESCE((SELECT array_agg(p.requires_node_id::text) FROM node_prerequisites p WHERE p.node_id = n.id), '{}'),
		       ` + lessonColumns + `,
		       EXISTS (SELECT 1 FROM lesson_attempts a WHERE a.user_id = $1 AND a.lesson_id = l.id AND a.status = 'completed'),
		       EXISTS (SELECT 1 FROM lesson_attempts a WHERE a.user_id = $1 AND a.lesson_id = l.id AND a.status = 'in_progress')
		FROM worlds w
		JOIN journey_nodes n ON n.world_id = w.id
		JOIN lessons l ON l.node_id = n.id
		WHERE w.status = 'published' AND n.status = 'published' AND l.status = 'published'
		ORDER BY w.sort_order, w.slug, n.sort_order`
	rows, err := s.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list path: %w", err)
	}
	defer rows.Close()

	var nodes []PathNode
	for rows.Next() {
		var pn PathNode
		var summary []byte
		l := &pn.Lesson
		if err := rows.Scan(
			&pn.WorldID, &pn.WorldSlug, &pn.WorldName, &pn.WorldDescription, &pn.WorldTheme,
			&pn.NodeID, &pn.NodeSlug, &pn.NodeLabel, &pn.NodeType, &pn.Requires,
			&l.ID, &l.Slug, &l.NodeID, &l.Title, &summary, &l.Kind, &l.Difficulty, &l.EstMinutes, &l.Status, &l.Version,
			&pn.Completed, &pn.InProgress,
		); err != nil {
			return nil, fmt.Errorf("scan path node: %w", err)
		}
		if err := json.Unmarshal(summary, &l.Summary); err != nil {
			return nil, fmt.Errorf("decode lesson summary: %w", err)
		}
		nodes = append(nodes, pn)
	}
	return nodes, rows.Err()
}

// PrerequisitesDone reports whether every published prerequisite of the node has
// a completed attempt by the user. Prerequisites that are not published are ignored.
func (s *LessonStore) PrerequisitesDone(ctx context.Context, userID, nodeID string) (bool, error) {
	const q = `
		SELECT NOT EXISTS (
			SELECT 1
			FROM node_prerequisites p
			JOIN journey_nodes pn ON pn.id = p.requires_node_id AND pn.status = 'published'
			JOIN lessons pl ON pl.node_id = pn.id AND pl.status = 'published'
			WHERE p.node_id = $1
			  AND NOT EXISTS (
				SELECT 1 FROM lesson_attempts a
				WHERE a.user_id = $2 AND a.lesson_id = pl.id AND a.status = 'completed'
			  )
		)`
	var done bool
	if err := s.pool.QueryRow(ctx, q, nodeID, userID).Scan(&done); err != nil {
		return false, fmt.Errorf("check prerequisites: %w", err)
	}
	return done, nil
}

const attemptColumns = `id, user_id, lesson_id, lesson_version, status,
	COALESCE(graded_steps, 0), COALESCE(first_try_correct, 0), COALESCE(total_tries, 0),
	COALESCE(completion_event_id::text, ''), event_published_at, started_at, completed_at`

func scanAttempt(row rowScanner) (*Attempt, error) {
	var a Attempt
	if err := row.Scan(&a.ID, &a.UserID, &a.LessonID, &a.LessonVersion, &a.Status,
		&a.GradedSteps, &a.FirstTryCorrect, &a.TotalTries,
		&a.CompletionEventID, &a.EventPublishedAt, &a.StartedAt, &a.CompletedAt); err != nil {
		return nil, err
	}
	return &a, nil
}

// GetAttempt returns an attempt by id, or nil.
func (s *LessonStore) GetAttempt(ctx context.Context, attemptID string) (*Attempt, error) {
	row := s.pool.QueryRow(ctx, `SELECT `+attemptColumns+` FROM lesson_attempts WHERE id = $1`, attemptID)
	a, err := scanAttempt(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get attempt: %w", err)
	}
	return a, nil
}

// GetInProgressAttempt returns the user's active attempt for a lesson, or nil.
func (s *LessonStore) GetInProgressAttempt(ctx context.Context, userID, lessonID string) (*Attempt, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT `+attemptColumns+` FROM lesson_attempts WHERE user_id = $1 AND lesson_id = $2 AND status = 'in_progress'`,
		userID, lessonID)
	a, err := scanAttempt(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get in-progress attempt: %w", err)
	}
	return a, nil
}

// StartAttempt creates an attempt, or returns the one already in progress when a
// concurrent request created it first.
func (s *LessonStore) StartAttempt(ctx context.Context, userID, lessonID string, version int) (*Attempt, error) {
	row := s.pool.QueryRow(ctx,
		`INSERT INTO lesson_attempts (user_id, lesson_id, lesson_version) VALUES ($1, $2, $3) RETURNING `+attemptColumns,
		userID, lessonID, version)
	a, err := scanAttempt(row)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return s.GetInProgressAttempt(ctx, userID, lessonID)
		}
		return nil, fmt.Errorf("start attempt: %w", err)
	}
	return a, nil
}

// ListResults returns every recorded try of an attempt, oldest first.
func (s *LessonStore) ListResults(ctx context.Context, attemptID string) ([]StepResult, error) {
	const q = `
		SELECT r.step_id, st.step_key, r.try_no, r.answer, r.correct, r.feedback
		FROM lesson_step_results r
		JOIN lesson_steps st ON st.id = r.step_id
		WHERE r.attempt_id = $1
		ORDER BY st.position, r.try_no`
	rows, err := s.pool.Query(ctx, q, attemptID)
	if err != nil {
		return nil, fmt.Errorf("list step results: %w", err)
	}
	defer rows.Close()

	var results []StepResult
	for rows.Next() {
		var r StepResult
		if err := rows.Scan(&r.StepID, &r.StepKey, &r.Try, &r.Answer, &r.Correct, &r.Feedback); err != nil {
			return nil, fmt.Errorf("scan step result: %w", err)
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

// InsertResult records one try. It returns false when that (attempt, step, try)
// was already recorded, so replays and concurrent duplicates never double-count.
func (s *LessonStore) InsertResult(ctx context.Context, attemptID, stepID string, try int, answer []byte, correct bool, feedback []byte) (bool, error) {
	const q = `
		INSERT INTO lesson_step_results (attempt_id, step_id, try_no, answer, correct, feedback)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (attempt_id, step_id, try_no) DO NOTHING`
	tag, err := s.pool.Exec(ctx, q, attemptID, stepID, try, answer, correct, feedback)
	if err != nil {
		return false, fmt.Errorf("insert step result: %w", err)
	}
	return tag.RowsAffected() == 1, nil
}

// GetResult returns one recorded try, or nil.
func (s *LessonStore) GetResult(ctx context.Context, attemptID, stepID string, try int) (*StepResult, error) {
	const q = `
		SELECT r.step_id, st.step_key, r.try_no, r.answer, r.correct, r.feedback
		FROM lesson_step_results r
		JOIN lesson_steps st ON st.id = r.step_id
		WHERE r.attempt_id = $1 AND r.step_id = $2 AND r.try_no = $3`
	var r StepResult
	err := s.pool.QueryRow(ctx, q, attemptID, stepID, try).Scan(&r.StepID, &r.StepKey, &r.Try, &r.Answer, &r.Correct, &r.Feedback)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get step result: %w", err)
	}
	return &r, nil
}

// CompleteAttempt marks an in-progress attempt completed. It returns the updated
// attempt, or nil when the attempt was not in progress.
func (s *LessonStore) CompleteAttempt(ctx context.Context, attemptID, eventID string, gradedSteps, firstTryCorrect, totalTries int) (*Attempt, error) {
	const q = `
		UPDATE lesson_attempts
		SET status = 'completed', completed_at = now(),
		    graded_steps = $2, first_try_correct = $3, total_tries = $4, completion_event_id = $5
		WHERE id = $1 AND status = 'in_progress'
		RETURNING ` + attemptColumns
	a, err := scanAttempt(s.pool.QueryRow(ctx, q, attemptID, gradedSteps, firstTryCorrect, totalTries, eventID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("complete attempt: %w", err)
	}
	return a, nil
}

// MarkEventPublished records that lesson.completed was published for an attempt.
func (s *LessonStore) MarkEventPublished(ctx context.Context, attemptID string) error {
	_, err := s.pool.Exec(ctx, `UPDATE lesson_attempts SET event_published_at = now() WHERE id = $1`, attemptID)
	if err != nil {
		return fmt.Errorf("mark event published: %w", err)
	}
	return nil
}
