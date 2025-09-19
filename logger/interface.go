package logger

import (
	"context"
)

// Logger interface for dependency injection and modern Go practices
type Logger interface {
	// Basic logging methods
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)

	// Formatted logging methods
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)

	// Context-aware methods
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	FatalContext(ctx context.Context, msg string, args ...any)

	// Formatted context-aware methods
	DebugfContext(ctx context.Context, format string, args ...any)
	InfofContext(ctx context.Context, format string, args ...any)
	WarnfContext(ctx context.Context, format string, args ...any)
	ErrorfContext(ctx context.Context, format string, args ...any)
	FatalfContext(ctx context.Context, format string, args ...any)

	// Structured logging
	With(args ...any) Logger
	WithGroup(name string) Logger

	// API logging
	API(statusCode int, msg string, args ...any)
	APIf(statusCode int, format string, args ...any)
	APIContext(ctx context.Context, statusCode int, msg string, args ...any)
	APIfContext(ctx context.Context, statusCode int, format string, args ...any)
}

// slog.Logger interface compatibility
type SlogCompatible interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) SlogCompatible
	WithGroup(name string) SlogCompatible
}

// Constructor interface
type LoggerConstructor interface {
	NewLogger(config JsonConfig) (Logger, error)
	SetupLogger(config JsonConfig) error
}
