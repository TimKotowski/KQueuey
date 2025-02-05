package kqueuey

import (
	"log/slog"
	"os"

	"go.uber.org/zap"
)

type Logging struct {
	Level  string
	Format string
}

func (l *Logging) NewLogger() *zap.Logger {
}

func (l *Logging) getLevel() zap.AtomicLevel {
	switch level := l.Level; level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}

	return slog.LevelInfo
}

func (l *Logging) getFormatLogFormat() slog.Handler {
	switch format := l.Format; format {
	case "text":
		return slog.NewTextHandler(os.Stderr, opts)
	case "json":
		return slog.NewJSONHandler(os.Stderr, opts)
	default:
		return slog.NewJSONHandler(os.Stderr, opts)
	}
}
