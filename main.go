package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go-pipeline/bootstrap"
	"go-pipeline/config"
	"go-pipeline/pkg/generate"
	"go-pipeline/pkg/logger"
	_ "net/http/pprof"
)

func main() {
	// Generate Main TraceID for APP
	traceID := generate.TraceID()

	// pprof
	go func() {
		logger.GetLogger().Info(&logger.Log{
			Event:   "pprof",
			TraceID: traceID,
			Additional: map[string]interface{}{
				"msg": "pprof listening on :6061",
				"cmd": "go tool pprof 'http://localhost:6061/debug/pprof/profile?seconds=10'",
			},
		})
		log.Fatal(http.ListenAndServe("localhost:6061", nil))
	}()

	// Bootstrap the application context
	ctxRAW := context.Background()
	ctxWithValue := context.WithValue(ctxRAW, config.TraceIDKey, traceID)
	ctx, cancel := context.WithCancel(ctxWithValue)
	defer cancel()

	// Initialize the bootstrap struct
	appInstance, err := bootstrap.Initialize(ctx)
	if err != nil {
		logger.GetLogger().Panic(&logger.Log{
			Event:      "main",
			Error:      err,
			TraceID:    traceID,
			Additional: nil,
		})
		os.Exit(1)
	}

	// Setup signal handling
	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Handle graceful shutdown
	go appInstance.GracefulShutdown(ctx, signals, done)

	// Start application
	appInstance.Start(ctx)

	// Wait until GracefulShutdown tells us we're done
	<-done

	// Cleanup context
	cancel()

	// Wait for background routines to finish
	appInstance.Wait()
}
