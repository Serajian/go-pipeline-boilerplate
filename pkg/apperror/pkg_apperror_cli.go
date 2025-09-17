package apperror

import (
	"errors"
)

// ExitCode maps app errors to CLI exit codes
func ExitCode(err error) int {
	if err == nil {
		return 0 // success
	}

	switch {
	case errors.Is(err, ErrInvalidInput):
		return 2 // invalid arguments
	case errors.Is(err, ErrUnauthorized):
		return 3 // permission denied
	case errors.Is(err, ErrTimeout):
		return 4 // timeout
	case errors.Is(err, ErrUnavailable):
		return 5 // unavailable
	default:
		return 1 // generic error
	}
}
