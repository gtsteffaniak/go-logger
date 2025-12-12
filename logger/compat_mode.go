package logger

import (
	"fmt"
)

// CompatibilityMode enables global logger functions for backward compatibility.
//
// This creates global state which is considered bad practice in Go.
// Consider using dependency injection instead:
//
//	logger := logger.NewLogger(config)
//	logger.Info("message")
//
// This function is provided for backward compatibility only.
// New code should prefer explicit logger instances.
func EnableCompatibilityMode(config JsonConfig) error {
	existingLogger := GetCompatibilityLogger()
	if existingLogger != nil {
		// Try to add this config to the existing logger
		if ml, ok := existingLogger.(*modernLogger); ok {
			err := ml.addConfig(config)
			if err != nil {
				return fmt.Errorf("add config to existing logger: %w", err)
			}

			// Don't call SetupLogger when using modern logger - it populates the legacy loggers array
			// which causes the global logger functions to use the legacy path instead of the modern logger
			return nil
		}
	}

	logger, err := NewLogger(config)
	if err != nil {
		return fmt.Errorf("enable compatibility mode: %w", err)
	}
	SetGlobalLogger(logger)

	// Don't call SetupLogger when using modern logger - it populates the legacy loggers array
	// which causes the global logger functions to use the legacy path instead of the modern logger
	// The modern logger handles all logging through the globalLogger interface

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
