package apperror

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// HTTPStatus maps app errors to HTTP status codes
func HTTPStatus(err error) int {
	if err == nil {
		return http.StatusOK // 200
	}

	switch {
	case errors.Is(err, ErrDuplicateEntry):
		return http.StatusConflict // 409
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound // 404
	case errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest // 400
	case isInvalidInput(err):
		return http.StatusBadRequest // 400
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized // 401
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden // 403
	case errors.Is(err, ErrTooMany):
		return http.StatusTooManyRequests // 429
	case errors.Is(err, ErrInternal):
		return http.StatusInternalServerError // 500
	case errors.Is(err, ErrUnavailable):
		return http.StatusServiceUnavailable // 503
	case errors.Is(err, ErrTimeout):
		return http.StatusGatewayTimeout // 504
	case errors.Is(err, context.Canceled):
		return http.StatusGatewayTimeout // 504
	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusGatewayTimeout // 504
	default:
		e := DBError(err)
		return HTTPStatus(e)
	}
}

// isInvalidInput checks if an error is a validation or binding error
func isInvalidInput(err error) bool {
	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeErr) {
		return true
	}
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		return true
	}
	if errors.Is(err, io.EOF) {
		return true
	}
	return false
}
