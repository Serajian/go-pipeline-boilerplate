package config

import (
	"context"
	"os"

	"go-pipeline/pkg/apperror"
	"go-pipeline/pkg/generate"
	"go-pipeline/pkg/logger"

	"github.com/Serajian/go-configmgr/configmgr"
)

// instance holds the singleton instance of the Config struct.
var instance *Config

// Config holds the overall configuration for the application.
type Config struct {
	AppConfig        AppConfig        `json:"app"         yaml:"app"`
	DBConfig         []DBConfig       `json:"databases"   yaml:"databases"`
	HTTPServer       HTTPServer       `json:"http_server" yaml:"http_server"`
	MQConfig         []MQConfig       `json:"mq"          yaml:"mq"`
	WorkerPoolConfig WorkerPoolConfig `json:"worker_pool" yaml:"worker_pool"`
}

// AppConfig holds configuration settings for the application.
type AppConfig struct {
	Name     string `json:"name"      validate:"required" yaml:"name"`
	Port     int    `json:"port"      validate:"required" yaml:"port"`
	Debug    bool   `json:"debug"                         yaml:"debug"     default:"true"`
	IsMaster bool   `json:"is_master" validate:"required" yaml:"is_master"`
}

// DBConfig holds configuration settings for a single database.
type DBConfig struct {
	Name     string `json:"name"     validate:"required" yaml:"name"`
	Type     string `json:"type"     validate:"required" yaml:"type"`
	User     string `json:"user"     validate:"required" yaml:"user"`
	Password string `json:"password" validate:"required" yaml:"password"`
	Host     string `json:"host"     validate:"required" yaml:"host"`
	Port     int    `json:"port"     validate:"required" yaml:"port"`
	SSL      bool   `json:"ssl"                          yaml:"ssl"`
}

// HTTPServer holds configuration settings for httpserver server.
type HTTPServer struct {
	Timeout Timeout `json:"timeout_second" yaml:"timeout_second"`
}

// Timeout holds configuration settings for a timeouts on httpserver server.
type Timeout struct {
	Write    int `json:"write"    yaml:"write"`
	Read     int `json:"read"     yaml:"read"`
	Idle     int `json:"idle"     yaml:"idle"`
	Shutdown int `json:"shutdown" yaml:"shutdown"`
}

// MQConfig holds configuration settings for the message queue.
type MQConfig struct {
	Name    string   `json:"name"     validate:"required" yaml:"name"`
	Type    string   `json:"type"     validate:"required" yaml:"type"`
	Port    int      `json:"port"     validate:"required" yaml:"port"`
	Address string   `json:"address"  validate:"required" yaml:"address"`
	GroupID string   `json:"group_id"                     yaml:"group_id"`
	Topics  []string `json:"topics"                       yaml:"topics"`
}

// WorkerPoolConfig holds configuration settings for the worker pool.
type WorkerPoolConfig struct {
	WorkerNum  int `json:"worker_num"  validate:"required" yaml:"worker_num"`
	QueueSize  int `json:"queue_size"  validate:"required" yaml:"queue_size"`
	RetryDelay int `json:"retry_delay" validate:"required" yaml:"retry_delay"`
	RetryMax   int `json:"retry_max"   validate:"required" yaml:"retry_max"`
}

// Get returns the singleton instance of the Config struct.
func Get() *Config {
	return instance
}

// init initializes the configuration by loading it from a file.
// It sets up the Config instance and logs the outcome.
func init() {
	log := logger.New()
	cm := configmgr.NewConfigManager()
	if os.Getenv("GO_ENV") == "test" {
		instance = &Config{}
		return
	}
	if err := cm.LoadFromFile(DefaultCFGPath); err != nil {
		log.Fatal(&logger.Log{
			Event:      "config",
			Error:      err,
			TraceID:    "config",
			Additional: map[string]interface{}{"msg": "failed to load config file"},
		})
	}

	var cfg Config
	if err := cm.Unmarshal(&cfg); err != nil {
		log.Fatal(&logger.Log{
			Event:      "config",
			Error:      err,
			TraceID:    "config",
			Additional: map[string]interface{}{"msg": "failed to unmarshal config"},
		})
	}

	instance = &cfg

	log.Info(&logger.Log{
		Event:      "config",
		Error:      nil,
		TraceID:    "config",
		Additional: map[string]interface{}{"msg": "successfully loaded config"},
	})
}

// GetTraceID retrieves the trace ID from the context.
// It returns the trace ID if found, otherwise logs an error and returns an empty string.
func GetTraceID(ctx context.Context) string {
	if val := ctx.Value(TraceIDKey); val != nil {
		if traceID, ok := val.(string); ok {
			return traceID
		}
	}
	logger.GetLogger().Error(&logger.Log{
		Event:      "get trace id",
		Error:      generate.Error("trace id not found in context", apperror.ErrInvalidInput),
		TraceID:    "config",
		Additional: map[string]interface{}{"key": TraceIDKey},
	})
	return ""
}
