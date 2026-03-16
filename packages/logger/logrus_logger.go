package logger

import (
	"log/slog"
	"os"
)

func newTextLogger(cfg Config) (Logger, error) {
	level := parseLevel(cfg.Level)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
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
