package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

// Logger wraps a Zerolog logger with custom configuration
type Logger struct {
	zLog zerolog.Logger
}

// Log defines the fixed structure for log entries
type Log struct {
	Event      Event
	Error      error
	TraceID    string
	Additional map[string]interface{}
}

// singleton instance
var (
	loggerInstance *Logger
	once           sync.Once
)

// New creates and configures a new Logger instance based on the environment.
//
// The logging mode (console vs. JSON) is automatically chosen by inspecting
// common environment variables in this order:
//  1. BART_ENV
//  2. APP_ENV
//  3. GO_ENV
//
// If the value is "local" or "dev", logs will be printed in a human-friendly
// console format with colors and emojis, making it easier to read during
// development.
//
// For any other value (including when none of the above environment variables
// are set), logs will be emitted in structured JSON format, which is suitable
// for production environments and log aggregators like ELK, Loki, or Datadog.
//
// The log level is also configurable via the LOG_LEVEL environment variable.
// Supported levels: "debug", "info", "warn", "error", "fatal", "panic".
// If LOG_LEVEL is not set, the default level is:
//   - DebugLevel for "local"/"dev"
//   - InfoLevel for everything else
//
// Example usage:
//
//	// Initialize a logger (singleton is preferred via GetLogger())
//	log := logger.New()
//
//	// Log an event
//	log.Info(&logger.Log{
//	    Event:   "server started",
//	    TraceID: "12345",
//	})
//
// In most cases, applications should use GetLogger() instead of calling New()
// directly, to ensure a singleton logger instance is reused throughout the app.
func New() *Logger {
	env := detectEnv() // detect dev/prod
	var zLogger zerolog.Logger

	if env == "local" || env == "dev" {
		zLogger = zerolog.New(zerolog.ConsoleWriter{
			Out:           os.Stdout,
			PartsOrder:    []string{"level", "message", "caller"},
			PartsExclude:  []string{"time"},
			FieldsOrder:   []string{"additional"},
			FieldsExclude: []string{"trace_id"},
			FormatCaller: func(i interface{}) string {
				if caller, ok := i.(string); ok {
					_, file := filepath.Split(caller)
					return fmt.Sprintf("caller=%s", file)
				}
				return fmt.Sprintf("caller=%v", i)
			},
		}).With().CallerWithSkipFrameCount(3).Timestamp().Logger()
	} else {
		// JSON output for non-local environments
		zLogger = zerolog.New(os.Stdout).With().CallerWithSkipFrameCount(3).Timestamp().Logger()
	}
	level := detectLevel(env)
	zerolog.SetGlobalLevel(level)

	return &Logger{zLog: zLogger}
}

// GetLogger returns the singleton Logger instance
func GetLogger() *Logger {
	once.Do(func() {
		loggerInstance = New()
	})
	return loggerInstance
}

// Fatal (LEVEL:0) logs a message at Fatal level with the fixed Log structure
func (l *Logger) Fatal(log *Log) {
	event := l.zLog.Fatal()
	l.addFixedFields(event, log)
	msg := log.Event.upper().addIcon(zerolog.FatalLevel)
	event.Msg(msg)
}

// Panic (LEVEL:1) logs a message at Panic level with the fixed Log structure
func (l *Logger) Panic(log *Log) {
	event := l.zLog.Panic()
	l.addFixedFields(event, log)
	msg := log.Event.upper().addIcon(zerolog.PanicLevel)
	event.Msg(msg)
}

// Error (LEVEL:2) logs a message at Error level with the fixed Log structure
func (l *Logger) Error(log *Log) {
	event := l.zLog.Error()
	l.addFixedFields(event, log)
	msg := log.Event.upper().addIcon(zerolog.ErrorLevel)
	event.Msg(msg)
}

// Warn (LEVEL:3) logs a message at Warn level with the fixed Log structure
func (l *Logger) Warn(log *Log) {
	event := l.zLog.Warn()
	l.addFixedFields(event, log)
	msg := log.Event.upper().addIcon(zerolog.WarnLevel)
	event.Msg(msg)
}

// Info (LEVEL:4) logs a message at Info level with the fixed Log structure
func (l *Logger) Info(log *Log) {
	event := l.zLog.Info()
	l.addFixedFields(event, log)
	msg := log.Event.upper().addIcon(zerolog.InfoLevel)
	event.Msg(msg)
}

// Debug (LEVEL:5) logs a message at Debug level with the fixed Log structure
func (l *Logger) Debug(log *Log) {
	event := l.zLog.Debug()
	l.addFixedFields(event, log)
	msg := log.Event.upper().addIcon(zerolog.DebugLevel)
	event.Msg(msg)
}

// addFixedFields adds the fixed fields from the Log struct to the event
func (l *Logger) addFixedFields(event *zerolog.Event, log *Log) {
	event.Str("trace_id", log.TraceID)
	if log.Error != nil {
		event.Err(log.Error)
	}
	// Add additional fields from the Log struct
	for key, value := range log.Additional {
		event.Interface(key, value)
	}
}

// Event represents a logical log event name.
// It can be transformed into uppercase and decorated with icons.
type Event string

func (m *Event) upper() *Event {
	words := strings.Split(string(*m), " ")
	for i, s := range words {
		words[i] = strings.ToUpper(s)
	}
	msg := strings.Join(words, "_")
	*m = Event(msg)
	return m
}

func (m *Event) addIcon(level zerolog.Level) string {
	switch level {
	case zerolog.FatalLevel:
		return "💀" + string(*m)
	case zerolog.PanicLevel:
		return "😱" + string(*m)
	case zerolog.ErrorLevel:
		return "🚨" + string(*m)
	case zerolog.WarnLevel:
		return "🚧" + string(*m)
	case zerolog.InfoLevel:
		return "✅" + string(*m)
	case zerolog.DebugLevel:
		return "👨‍💻" + string(*m)
	case zerolog.NoLevel:
		return "🛎" + string(*m)
	default:
		return ""
	}
}

// detectEnv checks multiple env vars
func detectEnv() string {
	keys := []string{"BART_ENV", "APP_ENV", "GO_ENV"}
	for _, key := range keys {
		if val := strings.ToLower(os.Getenv(key)); val != "" {
			return val
		}
	}
	return "prod" // fallback
}

// detectLevel picks log level (override with LOG_LEVEL)
func detectLevel(env string) zerolog.Level {
	// check LOG_LEVEL env
	if lvl := os.Getenv("LOG_LEVEL"); lvl != "" {
		if parsed, err := zerolog.ParseLevel(strings.ToLower(lvl)); err == nil {
			return parsed
		}
	}

	// fallback by env
	if env == "local" || env == "dev" {
		return zerolog.DebugLevel
	}
	return zerolog.InfoLevel
}
