package apperror

import "errors"

// MQError maps message queue errors to domain errors
func MQError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrTimeout):
		return ErrTimeout // requeue
	case errors.Is(err, ErrUnavailable):
		return ErrUnavailable // broker down
	default:
		return ErrInternal
	}
}
