package logger

import (
	"context"
	"log/slog"
	"os"
)

var String = slog.String

// Logger wraps slog.Logger with context-aware convenience methods.
type Logger struct {
	*slog.Logger
}

// Option configures a Logger.
type Option func(*options)

type options struct {
	level      slog.Level
	addSource  bool
	writer     *os.File
	attrs      []slog.Attr
	replaceErr func([]string, slog.Attr) slog.Attr
}

// WithLevel sets the minimum log level.k
func WithLevel(level slog.Level) Option {
	return func(o *options) { o.level = level }
}

// WithSource enables source code location in log output (file:line).
func WithSource(enabled bool) Option {
	return func(o *options) { o.addSource = enabled }
}

// WithWriter sets the output writer. Defaults to os.Stdout.
func WithWriter(w *os.File) Option {
	return func(o *options) { o.writer = w }
}

// WithServiceName adds a "service" attribute to every log record.
func WithServiceName(name string) Option {
	return func(o *options) {
		o.attrs = append(o.attrs, slog.String("service", name))
	}
}

// WithAttr adds a default attribute to every log record.
func WithAttr(key, value string) Option {
	return func(o *options) {
		o.attrs = append(o.attrs, slog.String(key, value))
	}
}

// New creates a new Logger with the given options.
func New(opts ...Option) *Logger {
	cfg := &options{
		level:  slog.LevelInfo,
		writer: os.Stdout,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	jsonHandler := slog.NewJSONHandler(cfg.writer, &slog.HandlerOptions{
		Level:     cfg.level,
		AddSource: cfg.addSource,
	})

	var handler slog.Handler = jsonHandler
	if len(cfg.attrs) > 0 {
		handler = handler.WithAttrs(cfg.attrs)
	}

	return &Logger{
		Logger: slog.New(newContextHandler(handler)),
	}
}

// With returns a new Logger with the given attributes.
func (l *Logger) With(args ...any) *Logger {
	return &Logger{Logger: l.Logger.With(args...)}
}

// WithGroup returns a new Logger with the given group name.
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{Logger: l.Logger.WithGroup(name)}
}

// --- Context-aware methods ---

func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.Logger.DebugContext(ctx, msg, args...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.Logger.InfoContext(ctx, msg, args...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.Logger.WarnContext(ctx, msg, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.Logger.ErrorContext(ctx, msg, args...)
}

// --- Convenience helpers ---

// Err returns an slog.Attr for logging errors consistently.
func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.String("error", err.Error())
}

// --- Default logger ---

var defaultLogger = New()

// SetDefault sets the package-level default logger and also sets it
// as the default for the standard slog package.
func SetDefault(l *Logger) {
	defaultLogger = l
	slog.SetDefault(l.Logger)
}

// Default returns the current default logger.
func Default() *Logger {
	return defaultLogger
}

// --- Package-level functions ---

func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.DebugContext(ctx, msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.InfoContext(ctx, msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.WarnContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.ErrorContext(ctx, msg, args...)
}
