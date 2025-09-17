// Package generate provides utility functions for generating.
package generate

import (
	"fmt"

	"github.com/google/uuid"
)

// TraceID generates a unique trace identifier.
// It returns a UUID string to be used for request tracking or logging purposes.
func TraceID() string {
	return uuid.New().String()
}

// Error wraps an existing error with additional context.
// It formats the provided text and error into a new error with the format "text: error".
func Error(txt string, err error) error {
	return fmt.Errorf("%s: %w", txt, err)
}
