package service

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrSkillNotFound  = errors.New("skill not found")
)

// Lesson runtime errors.
var (
	ErrLessonNotFound       = errors.New("lesson not found")
	ErrLessonLocked         = errors.New("lesson is locked")
	ErrAttemptNotFound      = errors.New("attempt not found")
	ErrStepNotFound         = errors.New("step not found")
	ErrAttemptNotInProgress = errors.New("attempt is not in progress")
	ErrAttemptIncomplete    = errors.New("every graded step needs a correct answer before completing")
	ErrInvalidAnswer        = errors.New("invalid answer")
	ErrInvalidTry           = errors.New("invalid try")
)
