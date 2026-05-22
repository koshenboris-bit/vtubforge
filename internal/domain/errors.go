package domain

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrValidation       = errors.New("validation error")
	ErrConflict         = errors.New("conflict")
	ErrInvalidToken     = errors.New("invalid token")
	ErrRefreshRevoked   = errors.New("refresh token revoked")
	ErrRefreshExpired   = errors.New("refresh token expired")
	ErrPasswordMismatch = errors.New("password mismatch")
)
