package service

import "errors"

var (
	ErrFreezeMaxHeld         = errors.New("max streak freezes held")
	ErrInsufficientGems      = errors.New("insufficient gems")
	ErrInvalidRequest        = errors.New("invalid request")
)
