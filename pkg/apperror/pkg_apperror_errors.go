package apperror

import "errors"

// Core application/domain errors
var (
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrNotFound       = errors.New("not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrTooMany        = errors.New("too many requests")
	ErrUnavailable    = errors.New("service unavailable")
	ErrTimeout        = errors.New("timeout")
	ErrInternal       = errors.New("internal error")
)
