package logger

import (
	"errors"
	"os"
	"strings"
)

type Field struct {
	Key   string
	Value any
}

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, err error, fields ...Field)
	With(fields ...Field) Logger
	Sync() error
}

type Config struct {
	Level       string
	Format      string
	Environment string
	ServiceName string
	AddSource   bool
}

func DefaultConfig() Config {
	return Config{
		Level:       "info",
		Format:      "json",
		Environment: "production",
		ServiceName: "api",
		AddSource:   true,
	}
}

func New(config Config) (Logger, error) {
	format := strings.TrimSpace(strings.ToLower(config.Format))
	switch format {
	case "", "json":
		return newJSONLogger(config)
	case "text":
		return newTextLogger(config)
	default:
		return nil, errors.New("unsupported logger format")
	}
}

func NewFromEnv() (Logger, error) {
	cfg := DefaultConfig()
	cfg.Level = getEnv("LOG_LEVEL", cfg.Level)
	cfg.Format = getEnv("LOG_FORMAT", cfg.Format)
	cfg.Environment = getEnv("APP_ENV", cfg.Environment)
	cfg.ServiceName = getEnv("SERVICE_NAME", cfg.ServiceName)

	if getEnv("LOG_SOURCE", "true") == "false" {
		cfg.AddSource = false
	}

	return New(cfg)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
