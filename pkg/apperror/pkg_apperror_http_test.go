package apperror_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"go-pipeline/pkg/apperror"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestHTTPStatus(t *testing.T) {
	validationErr := validator.ValidationErrors{}
	unmarshalErr := &json.UnmarshalTypeError{}

	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{"NilError", nil, http.StatusOK},
		{"DuplicateEntry", apperror.ErrDuplicateEntry, http.StatusConflict},
		{"NotFound", apperror.ErrNotFound, http.StatusNotFound},
		{"InvalidInput", apperror.ErrInvalidInput, http.StatusBadRequest},
		{"Unauthorized", apperror.ErrUnauthorized, http.StatusUnauthorized},
		{"Forbidden", apperror.ErrForbidden, http.StatusForbidden},
		{"TooMany", apperror.ErrTooMany, http.StatusTooManyRequests},
		{"Unavailable", apperror.ErrUnavailable, http.StatusServiceUnavailable},
		{"TimeoutError", apperror.ErrTimeout, http.StatusGatewayTimeout},
		{"ContextCanceled", context.Canceled, http.StatusGatewayTimeout},
		{"DeadlineExceeded", context.DeadlineExceeded, http.StatusGatewayTimeout},
		{"Internal", apperror.ErrInternal, http.StatusInternalServerError},
		{
			"EOFError",
			io.EOF,
			http.StatusBadRequest,
		}, // invalid input
		{
			"ValidationError",
			validationErr,
			http.StatusBadRequest,
		}, // invalid input
		{
			"UnmarshalTypeError",
			unmarshalErr,
			http.StatusBadRequest,
		}, // invalid input
		{"UnknownError", errors.New("random error"), http.StatusInternalServerError}, // fallback
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := apperror.HTTPStatus(tt.err)
			assert.Equal(t, tt.expected, status)
		})
	}
}
