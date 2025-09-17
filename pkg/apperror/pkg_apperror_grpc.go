package apperror

import (
	"errors"

	"google.golang.org/grpc/codes"
)

// GRPCCode maps app errors to gRPC codes
func GRPCCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	switch {
	case errors.Is(err, ErrDuplicateEntry):
		return codes.AlreadyExists
	case errors.Is(err, ErrNotFound):
		return codes.NotFound
	case errors.Is(err, ErrInvalidInput):
		return codes.InvalidArgument
	case errors.Is(err, ErrUnauthorized):
		return codes.Unauthenticated
	case errors.Is(err, ErrForbidden):
		return codes.PermissionDenied
	case errors.Is(err, ErrTooMany):
		return codes.ResourceExhausted
	case errors.Is(err, ErrUnavailable):
		return codes.Unavailable
	case errors.Is(err, ErrTimeout):
		return codes.DeadlineExceeded
	default:
		return codes.Internal
	}
}
