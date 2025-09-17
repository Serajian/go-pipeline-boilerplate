package httpserver

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go-pipeline/config"
	"go-pipeline/internal/ports"
	"go-pipeline/pkg/logger"
)

type ServerHTTP struct {
	server *http.Server
}

// Start start httpserver server in blocking mode
func (s *ServerHTTP) Start(ctx context.Context, handler http.Handler) error {
	traceID := config.GetTraceID(ctx)

	port := config.Get().AppConfig.Port
	addr := ":" + strconv.Itoa(port)

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(config.Get().HTTPServer.Timeout.Read) * time.Second,
		WriteTimeout: time.Duration(config.Get().HTTPServer.Timeout.Write) * time.Second,
		IdleTimeout:  time.Duration(config.Get().HTTPServer.Timeout.Idle) * time.Second,
	}

	s.server = server

	logger.GetLogger().Info(&logger.Log{
		Event:      "start httpserver server",
		Error:      nil,
		TraceID:    traceID,
		Additional: map[string]interface{}{"port": config.Get().AppConfig.Port},
	})

	return server.ListenAndServe()
}

// Stop handles the httpserver server in graceful shutdown
func (s *ServerHTTP) Stop(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(
		ctx,
		time.Duration(config.Get().HTTPServer.Timeout.Shutdown)*time.Second,
	)
	defer cancel()

	s.server.SetKeepAlivesEnabled(false)

	if err := s.server.Shutdown(ctxTimeout); err != nil {
		return err
	}
	logger.GetLogger().Warn(&logger.Log{
		Event:   "stop httpserver server",
		Error:   nil,
		TraceID: config.GetTraceID(ctx),
	})
	return nil
}

var _ ports.HTTPServer = (*ServerHTTP)(nil)
