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
	ErrInvalidOTP          = errors.New("invalid or expired verification code")
	ErrOTPCooldown         = errors.New("please wait before requesting another code")
	ErrTooManyAttempts     = errors.New("too many failed attempts, please request a new code")
)
