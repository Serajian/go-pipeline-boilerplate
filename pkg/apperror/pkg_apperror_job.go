package apperror

import "errors"

// JobCode maps app errors to Job worker codes (string)
func JobCode(err error) string {
	if err == nil {
		return "ok"
	}

	switch {
	case errors.Is(err, ErrDuplicateEntry):
		return "conflict"
	case errors.Is(err, ErrNotFound):
		return "not_found"
	case errors.Is(err, ErrInvalidInput):
		return "invalid_input"
	case errors.Is(err, ErrUnauthorized):
		return "unauthorized"
	case errors.Is(err, ErrForbidden):
		return "forbidden"
	case errors.Is(err, ErrTooMany):
		return "rate_limited"
	case errors.Is(err, ErrUnavailable):
		return "unavailable"
	case errors.Is(err, ErrTimeout):
		return "timeout"
	default:
		return "internal"
	}
}
