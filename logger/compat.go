package logger

import (
	"context"
	"fmt"
)

// SetupLoggerWithModern creates a modern logger instance and sets it as global
func SetupLoggerWithModern(config JsonConfig) error {
	logger, err := NewLogger(config)
	if err != nil {
		return fmt.Errorf("setup modern logger: %w", err)
	}
	SetGlobalLogger(logger) // Set it in write.go for legacy compatibility

	// Don't setup legacy logger if it already exists
	// The modern logger will handle both modern and legacy calls
	return nil
}

// GetGlobalLogger returns the current global logger instance
func GetGlobalLogger() Logger {
	return globalLogger
}

// Modern logging functions that use the global logger if available
func DebugContext(ctx context.Context, msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.DebugContext(ctx, msg, args...)
	} else {
		// Fall back to legacy logging
		Debug(msg)
	}
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.InfoContext(ctx, msg, args...)
	} else {
		// Fall back to legacy logging
		Info(msg)
	}
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.WarnContext(ctx, msg, args...)
	} else {
		// Fall back to legacy logging
		Warning(msg)
	}
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.ErrorContext(ctx, msg, args...)
	} else {
		// Fall back to legacy logging
		Error(msg)
	}
}

func FatalContext(ctx context.Context, msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.FatalContext(ctx, msg, args...)
	} else {
		// Fall back to legacy logging
		Fatal(msg)
	}
}

// Formatted context-aware functions
func DebugfContext(ctx context.Context, format string, args ...any) {
	if globalLogger != nil {
		globalLogger.DebugfContext(ctx, format, args...)
	} else {
		// Fall back to legacy logging
		Debugf(format, args...)
	}
}

func InfofContext(ctx context.Context, format string, args ...any) {
	if globalLogger != nil {
		globalLogger.InfofContext(ctx, format, args...)
	} else {
		// Fall back to legacy logging
		Infof(format, args...)
	}
}

func WarnfContext(ctx context.Context, format string, args ...any) {
	if globalLogger != nil {
		globalLogger.WarnfContext(ctx, format, args...)
	} else {
		// Fall back to legacy logging
		Warningf(format, args...)
	}
}

func ErrorfContext(ctx context.Context, format string, args ...any) {
	if globalLogger != nil {
		globalLogger.ErrorfContext(ctx, format, args...)
	} else {
		// Fall back to legacy logging
		Errorf(format, args...)
	}
}

func FatalfContext(ctx context.Context, format string, args ...any) {
	if globalLogger != nil {
		globalLogger.FatalfContext(ctx, format, args...)
	} else {
		// Fall back to legacy logging
		Fatalf(format, args...)
	}
}

// Structured logging functions
func With(args ...any) Logger {
	if globalLogger != nil {
		return globalLogger.With(args...)
	}
	// Return a no-op logger if no global logger is set
	return &noOpLogger{}
}

func WithGroup(name string) Logger {
	if globalLogger != nil {
		return globalLogger.WithGroup(name)
	}
	// Return a no-op logger if no global logger is set
	return &noOpLogger{}
}

// API context-aware functions
func APIContext(ctx context.Context, statusCode int, msg string, args ...any) {
	if globalLogger != nil {
		globalLogger.APIContext(ctx, statusCode, msg, args...)
	} else {
		// Fall back to legacy logging
		Api(statusCode, msg)
	}
}

func APIfContext(ctx context.Context, statusCode int, format string, args ...any) {
	if globalLogger != nil {
		globalLogger.APIfContext(ctx, statusCode, format, args...)
	} else {
		// Fall back to legacy logging
		Apif(statusCode, format, args...)
	}
}

// noOpLogger is a no-op implementation for when no global logger is set
type noOpLogger struct{}

func (n *noOpLogger) Debug(msg string, args ...any)                                               {}
func (n *noOpLogger) Info(msg string, args ...any)                                                {}
func (n *noOpLogger) Warn(msg string, args ...any)                                                {}
func (n *noOpLogger) Error(msg string, args ...any)                                               {}
func (n *noOpLogger) Fatal(msg string, args ...any)                                               {}
func (n *noOpLogger) Debugf(format string, args ...any)                                           {}
func (n *noOpLogger) Infof(format string, args ...any)                                            {}
func (n *noOpLogger) Warnf(format string, args ...any)                                            {}
func (n *noOpLogger) Errorf(format string, args ...any)                                           {}
func (n *noOpLogger) Fatalf(format string, args ...any)                                           {}
func (n *noOpLogger) DebugContext(ctx context.Context, msg string, args ...any)                   {}
func (n *noOpLogger) InfoContext(ctx context.Context, msg string, args ...any)                    {}
func (n *noOpLogger) WarnContext(ctx context.Context, msg string, args ...any)                    {}
func (n *noOpLogger) ErrorContext(ctx context.Context, msg string, args ...any)                   {}
func (n *noOpLogger) FatalContext(ctx context.Context, msg string, args ...any)                   {}
func (n *noOpLogger) DebugfContext(ctx context.Context, format string, args ...any)               {}
func (n *noOpLogger) InfofContext(ctx context.Context, format string, args ...any)                {}
func (n *noOpLogger) WarnfContext(ctx context.Context, format string, args ...any)                {}
func (n *noOpLogger) ErrorfContext(ctx context.Context, format string, args ...any)               {}
func (n *noOpLogger) FatalfContext(ctx context.Context, format string, args ...any)               {}
func (n *noOpLogger) With(args ...any) Logger                                                     { return n }
func (n *noOpLogger) WithGroup(name string) Logger                                                { return n }
func (n *noOpLogger) API(statusCode int, msg string, args ...any)                                 {}
func (n *noOpLogger) APIf(statusCode int, format string, args ...any)                             {}
func (n *noOpLogger) APIContext(ctx context.Context, statusCode int, msg string, args ...any)     {}
func (n *noOpLogger) APIfContext(ctx context.Context, statusCode int, format string, args ...any) {}
