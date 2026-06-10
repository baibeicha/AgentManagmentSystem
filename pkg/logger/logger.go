package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

var LogLevel = new(slog.LevelVar)

type ContextHandler struct {
	slog.Handler
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		r.AddAttrs(slog.String("request_id", reqID))
	}
	return h.Handler.Handle(ctx, r)
}

func SetLevel(levelStr string) {
	switch strings.ToLower(levelStr) {
	case "debug":
		LogLevel.Set(slog.LevelDebug)
	case "info":
		LogLevel.Set(slog.LevelInfo)
	case "warn", "warning":
		LogLevel.Set(slog.LevelWarn)
	case "error", "err":
		LogLevel.Set(slog.LevelError)
	default:
		LogLevel.Set(slog.LevelInfo)
	}
}

func SetupLogger(env string, levelStr string) *slog.Logger {
	SetLevel(levelStr)

	var handler slog.Handler

	if env == "production" || env == "prod" {
		opts := &slog.HandlerOptions{
			Level:     LogLevel,
			AddSource: true,
		}
		handler = slog.NewJSONHandler(os.Stdout, opts)

	} else {
		opts := &slog.HandlerOptions{
			Level: LogLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.String(a.Key, a.Value.Time().Format(time.DateTime))
				}
				if a.Key == slog.SourceKey {
					source := a.Value.Any().(*slog.Source)
					source.File = filepath.Base(source.File)
				}
				if a.Key == slog.LevelKey {
					level := a.Value.Any().(slog.Level)
					var color string
					switch level {
					case slog.LevelDebug:
						color = "\033[36m"
					case slog.LevelInfo:
						color = "\033[32m"
					case slog.LevelWarn:
						color = "\033[33m"
					case slog.LevelError:
						color = "\033[31m"
					default:
						color = "\033[0m"
					}
					return slog.String(a.Key, color+level.String()+"\033[0m")
				}
				return a
			},
		}
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(ContextHandler{Handler: handler})
	slog.SetDefault(logger)

	return logger
}
