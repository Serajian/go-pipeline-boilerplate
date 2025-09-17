package ports

import (
	"context"
	"net/http"
)

type HTTPServer interface {
	Start(ctx context.Context, handler http.Handler) error
	Stop(ctx context.Context) error
}
