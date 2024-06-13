package logging

import (
	"context"
	"log/slog"
)

type (
	// ctxLogAttrs is the context key for log attrs.
	ctxLogAttrs struct{}
	// ctxLogger is the context key for the logger.
	ctxLogger struct{}
)

// ContextHandler wraps slog json handler to propagate context values
// to the logger.
type ContextHandler struct {
	*slog.JSONHandler
}

// Handle adds slog fields from context to logger attributes.
func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(ctxLogAttrs{}).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}
	return h.JSONHandler.Handle(ctx, r)
}

// FromContext is a helper method to get logger from upstream context.
// If there is no existing logger in the context, it will create a new log
// instance and add it to the context.
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(ctxLogger{}).(*Logger); ok {
		return logger
	}
	return New()
}

// WithContext adds logger and (optional) log attributes to the context.
func WithContext(ctx context.Context, l *Logger, args ...any) context.Context {
	var v []slog.Attr
	if logAttrs, ok := ctx.Value(ctxLogAttrs{}).([]slog.Attr); ok {
		v = logAttrs
	}

	for _, attr := range args {
		v = append(v, attr.(slog.Attr))
	}
	c := context.WithValue(ctx, ctxLogAttrs{}, v)
	return context.WithValue(c, ctxLogger{}, l)
}
