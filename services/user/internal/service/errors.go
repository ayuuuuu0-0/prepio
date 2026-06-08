package service

import "errors"

var (
	ErrInvalidRequest      = errors.New("invalid request")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailTaken          = errors.New("email taken")
	ErrUsernameTaken       = errors.New("username taken")
	ErrInvalidToken        = errors.New("invalid token")
	ErrRefreshTokenInvalid = errors.New("refresh token invalid")
	ErrUserNotFound        = errors.New("user not found")
	ErrDeviceNotFound      = errors.New("device not found")
)
