package service

import "errors"

var (
	ErrInvalidRequest   = errors.New("invalid request")
	ErrInsufficientGems = errors.New("insufficient gems")
)
