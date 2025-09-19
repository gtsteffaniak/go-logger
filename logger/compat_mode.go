package logger

import (
	"fmt"
	"strings"
)

// CompatibilityMode enables global logger functions for backward compatibility.
//
// ⚠️  WARNING: This creates global state which is considered bad practice in Go.
// Consider using dependency injection instead:
//
//	logger := logger.NewLogger(config)
//	logger.Info("message")
//
// This function is provided for backward compatibility only.
// New code should prefer explicit logger instances.
func EnableCompatibilityMode(config JsonConfig) error {
	logger, err := NewLogger(config)
	if err != nil {
		return fmt.Errorf("enable compatibility mode: %w", err)
	}
	SetGlobalLogger(logger)

	// Try to setup legacy logger, but don't fail if it already exists
	err = SetupLogger(config)
	if err != nil && !strings.Contains(err.Error(), "stdout logger already exists") {
		return fmt.Errorf("setup legacy logger: %w", err)
	}

	return nil
}

// DisableCompatibilityMode removes the global logger and falls back to basic logging.
func DisableCompatibilityMode() {
	SetGlobalLogger(nil)
	// Note: We don't clear the legacy loggers array to maintain existing behavior
}

// IsCompatibilityModeEnabled returns true if global logger functions are available.
func IsCompatibilityModeEnabled() bool {
	return globalLogger != nil
}

// GetCompatibilityLogger returns the current global logger instance.
// Returns nil if compatibility mode is not enabled.
func GetCompatibilityLogger() Logger {
	return globalLogger
}
