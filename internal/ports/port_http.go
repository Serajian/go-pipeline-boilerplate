package ports

import (
	"context"
	"net/http"
)

// HTTPServer defines an abstraction for an HTTP server.
// It allows starting and gracefully stopping the server,
// decoupled from the actual implementation (e.g., net/http).
//
// Start runs the server with the given handler and blocks
// until the context is canceled or the server is stopped.
// Stop gracefully shuts down the server, releasing all resources.
type HTTPServer interface {
	Start(ctx context.Context, handler http.Handler) error
	Stop(ctx context.Context) error
}
