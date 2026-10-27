package store

import "errors"

// ErrInsufficientGems is returned when a gem deduction would overdraw balance.
var ErrInsufficientGems = errors.New("insufficient gems")
