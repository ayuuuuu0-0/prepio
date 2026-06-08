package service

import "errors"

var (
	ErrInvalidRequest         = errors.New("invalid request")
	ErrQuestionNotFound       = errors.New("question not found")
	ErrSessionNotFound        = errors.New("session not found")
	ErrQuestionNotInSession   = errors.New("question not in session")
	ErrAnswerAlreadySubmitted = errors.New("answer already submitted")
	ErrUserNotFound           = errors.New("user not found")
)
