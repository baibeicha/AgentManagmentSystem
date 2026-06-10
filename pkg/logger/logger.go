package logger

import (
	"context"
	"fmt"
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

type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: handlers}
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

func SetupLogger(env, levelStr, filename string) (*slog.Logger, *os.File, error) {
	SetLevel(levelStr)

	var handlers []slog.Handler

	var consoleHandler slog.Handler
	if env == "production" {
		consoleHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: LogLevel, AddSource: true})
	} else {
		consoleHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
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
						color = "\033[36m" // Cyan
					case slog.LevelInfo:
						color = "\033[32m" // Green
					case slog.LevelWarn:
						color = "\033[33m" // Yellow
					case slog.LevelError:
						color = "\033[31m" // Red
					default:
						color = "\033[0m"
					}
					return slog.String(a.Key, color+level.String()+"\033[0m")
				}
				return a
			},
		})
	}
	handlers = append(handlers, consoleHandler)

	var logFile *os.File
	var err error
	if filename != "" {
		logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			err = fmt.Errorf("can not open logs file: %w", err)
		}

		fileHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level:     LogLevel,
			AddSource: true,
		})

		handlers = append(handlers, fileHandler)
	}

	logger := slog.New(ContextHandler{Handler: NewMultiHandler(handlers...)})
	slog.SetDefault(logger)

	return logger, logFile, err
}
