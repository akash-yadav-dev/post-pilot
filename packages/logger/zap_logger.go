package logger

import (
	"log/slog"
	"os"
	"strings"
)

type slogLogger struct {
	logger *slog.Logger
}

func newJSONLogger(cfg Config) (Logger, error) {
	level := parseLevel(cfg.Level)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     level,
	})

	base := slog.New(handler)
	base = base.With(
		slog.String("env", cfg.Environment),
		slog.String("service", cfg.ServiceName),
	)

	return &slogLogger{logger: base}, nil
}

func (l *slogLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fieldsToAttrs(fields)...)
}

func (l *slogLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fieldsToAttrs(fields)...)
}

func (l *slogLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fieldsToAttrs(fields)...)
}

func (l *slogLogger) Error(msg string, err error, fields ...Field) {
	attrs := fieldsToAttrs(fields)
	if err != nil {
		attrs = append(attrs, slog.Any("error", err))
	}
	l.logger.Error(msg, attrs...)
}

func (l *slogLogger) With(fields ...Field) Logger {
	return &slogLogger{logger: l.logger.With(fieldsToAttrs(fields)...)}
}

func (l *slogLogger) Sync() error {
	return nil
}

func parseLevel(level string) slog.Level {
	switch strings.TrimSpace(strings.ToLower(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func fieldsToAttrs(fields []Field) []any {
	if len(fields) == 0 {
		return nil
	}

	attrs := make([]any, 0, len(fields))
	for _, field := range fields {
		attrs = append(attrs, slog.Any(field.Key, field.Value))
	}

	return attrs
}
