package logger

import (
	"context"
	"log/slog"
)

// contextHandler wraps an slog.Handler to automatically extract
// context values (request_id, is_debug) and inject them into log records.
type contextHandler struct {
	handler slog.Handler
}

func newContextHandler(inner slog.Handler) *contextHandler {
	return &contextHandler{handler: inner}
}

func (h *contextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Allow debug logs through if is_debug is set in context,
	// regardless of the global log level.
	if level == slog.LevelDebug && IsDebugFromContext(ctx) {
		return true
	}
	return h.handler.Enabled(ctx, level)
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := RequestIDFromContext(ctx); ok {
		r.AddAttrs(slog.String("request_id", id))
	}
	return h.handler.Handle(ctx, r)
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{handler: h.handler.WithGroup(name)}
}
