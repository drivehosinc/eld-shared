package logger

import "context"

type requestIDKey struct{}
type isDebugKey struct{}

// WithRequestID returns a new context with the given request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// WithDebug returns a new context with debug logging enabled.
// This allows debug-level logs to be emitted for this specific
// request even when the global log level is higher.
func WithDebug(ctx context.Context) context.Context {
	return context.WithValue(ctx, isDebugKey{}, true)
}

// RequestIDFromContext extracts the request ID from context, if present.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey{}).(string)
	return id, ok && id != ""
}

// IsDebugFromContext checks whether debug logging is enabled in the context.
func IsDebugFromContext(ctx context.Context) bool {
	v, _ := ctx.Value(isDebugKey{}).(bool)
	return v
}
