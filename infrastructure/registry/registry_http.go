package registry

import (
	"context"
	"net/http"

	"go-pipeline/infrastructure/httpserver"
	"go-pipeline/internal/ports"
)

type HTTPServerRegistry struct {
	server  ports.HTTPServer
	handler http.Handler
}

func NewHTTPServerRegistry(handler http.Handler) *HTTPServerRegistry {
	return &HTTPServerRegistry{
		server:  &httpserver.ServerHTTP{},
		handler: handler,
	}
}

func (s *HTTPServerRegistry) Start(ctx context.Context) error {
	if err := s.server.Start(ctx, s.handler); err != nil {
		return err
	}
	return nil
}

func (s *HTTPServerRegistry) Stop(ctx context.Context) error {
	if err := s.server.Stop(ctx); err != nil {
		return err
	}
	return nil
}
