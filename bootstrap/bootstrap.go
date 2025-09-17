package bootstrap

import (
	"context"
	"fmt"
	"go-pipeline/internal/di"
	"go-pipeline/internal/presentation/http"
	"go-pipeline/internal/presentation/mq"
	"os"
	"runtime/debug"
	"sync"

	"go-pipeline/config"
	"go-pipeline/infrastructure/registry"
	"go-pipeline/pkg/logger"
)

// App encapsulates the application's core services.
type App struct {
	sync.WaitGroup
	httpServer *registry.HTTPServerRegistry
	mq         *registry.MQRegistry
	pipelines  *di.Pipelines
	stages     *di.Stages
}

// Initialize sets up the application's core services.
func Initialize(ctx context.Context) (*App, error) {
	var err error
	traceID := config.GetTraceID(ctx)
	log := logger.New()

	app := &App{}

	// 1) initialize databases

	// 2) initialize message queue
	handler := mq.NewConsumerHandler()
	mqRegistry, err := registry.NewMQRegistry(ctx, handler)
	if err != nil {
		return nil, err
	}
	log.Info(&logger.Log{
		Event:   "initialize mq",
		TraceID: traceID,
	})
	app.mq = mqRegistry
	// 3) initialize stages
	app.stages = di.NewStagesContainer(app.mq.GetKafkaProducer())

	// 4) initialize pipelines
	app.pipelines = di.NewPipelines(app.stages)

	// 5) initialize httpserver server
	handlerHTTP := http.NewGinAdapter(app.pipelines.Registry)
	httpRegistry := registry.NewHTTPServerRegistry(handlerHTTP.Engin)
	app.httpServer = httpRegistry
	log.Info(&logger.Log{
		Event:   "initialize http server",
		TraceID: traceID,
	})

	// 6) initialize scheduler

	// 7) initialize worker pool

	log.Info(&logger.Log{
		Event:   "finish initializing app",
		TraceID: traceID,
	})

	return app, err
}

// Start begins the application's core services.
func (app *App) Start(ctx context.Context) {
	traceID := config.GetTraceID(ctx)
	logger.GetLogger().Info(&logger.Log{
		Event:   "start app",
		TraceID: traceID,
	})

	app.Add(2)
	go app.safeRun("httpserver server", traceID, func() {
		if err := app.httpServer.Start(ctx); err != nil {
			logger.GetLogger().Error(&logger.Log{
				Event:      "start app",
				Error:      err,
				TraceID:    traceID,
				Additional: map[string]interface{}{"msg": "failed to start httpserver server"},
			})
		}
	})
	go app.safeRun("kafka consumer", traceID, func() {
		if err := app.mq.GetKafkaConsumer().Consume(ctx); err != nil {
			logger.GetLogger().Error(&logger.Log{
				Event:      "start kafka consumer",
				Error:      err,
				TraceID:    traceID,
				Additional: map[string]interface{}{"msg": "failed to start kafka consumer"},
			})
		}
	})

}

// Stop gracefully shuts down the application's core services.
func (app *App) Stop(ctx context.Context) {
	traceID := config.GetTraceID(ctx)
	logger.GetLogger().Warn(&logger.Log{
		Event:      "stop app",
		Error:      nil,
		TraceID:    traceID,
		Additional: nil,
	})

	if app.httpServer != nil {
		if err := app.httpServer.Stop(ctx); err != nil {
			logger.GetLogger().Error(&logger.Log{
				Event:      "stop app",
				Error:      err,
				TraceID:    traceID,
				Additional: map[string]interface{}{"msg": "failed to stop httpserver server"},
			})
		}
	}
}

// GracefulShutdown handles the graceful shutdown of the application.
// It waits for an OS signal, stops the application services, and signals completion.
func (app *App) GracefulShutdown(
	ctx context.Context,
	quitSignal <-chan os.Signal,
	done chan<- bool,
) {
	traceID := config.GetTraceID(ctx)

	// wait for os signals
	<-quitSignal

	logger.GetLogger().Warn(&logger.Log{
		Event:      "start graceful shutdown",
		Error:      nil,
		TraceID:    traceID,
		Additional: nil,
	})

	app.Stop(ctx)

	logger.GetLogger().Info(&logger.Log{
		Event:      "end graceful shutdown",
		Error:      nil,
		TraceID:    traceID,
		Additional: map[string]interface{}{"msg": "gracefully shutdown complete"},
	})

	close(done)
}

// safeRun wraps a function to recover from panic and log it
func (app *App) safeRun(component, traceID string, fn func()) {
	msg := fmt.Sprintf("%s panic recovered", component)
	defer app.Done()
	defer func() {
		if r := recover(); r != nil {
			logger.GetLogger().Error(&logger.Log{
				Event:   "start app",
				Error:   fmt.Errorf("%v", r),
				TraceID: traceID,
				Additional: map[string]interface{}{
					"stack": string(debug.Stack()),
					"msg":   msg,
				},
			})
		}
	}()
	fn()
}
