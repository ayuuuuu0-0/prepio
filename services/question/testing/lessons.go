package testing

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/services/question/internal/dto"
	"github.com/prepio/prepio/services/question/internal/grading"
	"github.com/prepio/prepio/services/question/internal/lesson"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/prepio/prepio/services/question/internal/store"
)

// Re-exports so integration tests outside this module tree can build content and call the runtime.
type (
	Content      = lesson.Content
	World        = lesson.World
	Node         = lesson.Node
	Lesson       = lesson.Lesson
	SkillRef     = lesson.SkillRef
	Step         = lesson.Step
	Beat         = lesson.Beat
	Explanations = lesson.Explanations
	Rubric       = lesson.Rubric
	Concept      = lesson.RubricConcept

	Answer         = grading.Answer
	AnswerRequest  = dto.AnswerRequest
	AnswerResponse = dto.AnswerResponse
	AttemptResp    = dto.AttemptResponse
	CompleteResp   = dto.CompleteResponse
	PathResp       = dto.PathResponse
	SyncReport     = store.SyncReport
)

// Service errors tests assert on.
var (
	ErrLessonLocked      = service.ErrLessonLocked
	ErrLessonNotFound    = service.ErrLessonNotFound
	ErrAttemptNotFound   = service.ErrAttemptNotFound
	ErrStepNotFound      = service.ErrStepNotFound
	ErrAttemptIncomplete = service.ErrAttemptIncomplete
	ErrInvalidAnswer     = service.ErrInvalidAnswer
	ErrInvalidTry        = service.ErrInvalidTry
)

// NewLessonService wires the lesson runtime for integration tests (structural prose grading, no LLM).
func NewLessonService(pool *pgxpool.Pool, publisher service.EventPublisher) *service.LessonService {
	return service.NewLessonService(store.NewLessonStore(pool), service.NewPipelineEvaluator(nil), publisher)
}

// NewContentSync returns the content sync store.
func NewContentSync(pool *pgxpool.Pool) *store.ContentSyncStore {
	return store.NewContentSyncStore(pool)
}

// Validate validates content against the skill slugs in the database.
func Validate(c *Content, skills map[string]bool) []string {
	return lesson.Validate(c, lesson.Catalog{Skills: skills})
}

// LoadContent loads the authored content directory.
func LoadContent(root string) (*Content, error) {
	return lesson.Load(root)
}

// Step type and status constants for building fixtures.
const (
	StepIntro     = lesson.StepIntro
	StepMCQ       = lesson.StepMCQ
	StepTrueFalse = lesson.StepTrueFalse
	StepFillBlank = lesson.StepFillBlank
	StepArrange   = lesson.StepArrange
	StepProse     = lesson.StepProse

	StatusPublished = lesson.StatusPublished
)
