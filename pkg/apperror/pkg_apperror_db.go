package apperror

import (
	"database/sql"
	"errors"
	"strings"
)

// DBError maps low-level DB errors to domain errors
func DBError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound
	case strings.Contains(err.Error(), "duplicate key"):
		return ErrDuplicateEntry
	case strings.Contains(err.Error(), "timeout"):
		return ErrTimeout
	default:
		return ErrInternal
	}
}
