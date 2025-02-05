package kqueuey

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
)

const (
	timeFormat = "[15:05:05.000]"
)

type Logging struct {
	Level  string
	Format string
}

type LogHandler struct {
	slog.Handler
	l *log.Logger
	j jsoniter.API
}


func (l *Logging) NewLogger() *slog.Logger {
	opts := &slog.HandlerOptions{}
	opts.Level = l.getLevel()
	logHandler := l.getFormatHandler(opts)
	color.NoColor = false

	return slog.New(
		&LogHandler{
			logHandler,
			log.New(os.Stderr, "", 0),
			jsoniter.ConfigFastest,
		},
	)
}

func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	var b string
	var err error
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	// Later do benchmarking to see if using pool for map to reduce memory pressure from logging always making a map.
	fields := make(map[string]any, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	if len(fields) > 0 {
		b, err = h.j.MarshalToString(fields)
		if err != nil {
			return err
		}
	}

	h.l.Println(r.Time.Format(timeFormat), level, r.Message, b)

	return nil
}

func (l *Logging) getLevel() slog.Level {
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

func (l *Logging) getFormatHandler(opts *slog.HandlerOptions) slog.Handler {
	switch format := l.Format; format {
	case "text":
		return slog.NewTextHandler(os.Stderr, opts)
	case "json":
		return slog.NewJSONHandler(os.Stderr, opts)
	default:
		return slog.NewJSONHandler(os.Stderr, opts)
	}
}
