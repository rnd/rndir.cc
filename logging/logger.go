package logging

import (
	"context"
	"log/slog"
	"os"
)

// Logger is a wrapper of slog logger package.
type Logger struct {
	*slog.Logger
}

const (
	LevelTrace = slog.Level(slog.LevelDebug - 4)
	LevelFatal = slog.Level(slog.LevelError + 4)
)

func fmtLevel(attr slog.Attr) slog.Attr {
	if attr.Key != slog.LevelKey {
		return attr // nothing to do
	}

	level := attr.Value.Any().(slog.Level)
	switch level {
	case LevelTrace:
		attr.Value = slog.StringValue("TRACE")
	case LevelFatal:
		attr.Value = slog.StringValue("FATAL")
	}
	return attr
}

func New() *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			switch attr.Key {
			case slog.LevelKey:
				return fmtLevel(attr)
			}
			return attr
		},
	}
	handler := &ContextHandler{
		slog.NewJSONHandler(
			os.Stdout,
			opts,
		),
	}
	return &Logger{
		slog.New(handler),
	}
}

// Fatal logs at custom fatal error and exit with error.
func (l *Logger) Fatal(msg string, args ...any) {
	l.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

// FatalContext logs at custom fatal error and exit with error with context.
func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelFatal, msg, args...)
	os.Exit(1)
}
