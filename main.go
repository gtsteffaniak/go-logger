package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gtsteffaniak/go-logger/logger"
)

func main() {
	// Legacy usage (backward compatible)
	legacyExample()

	// Modern usage with dependency injection (recommended)
	modernExample()

	// Compatibility mode (discouraged but available)
	compatibilityModeExample()
}

func legacyExample() {
	println("=== Legacy Usage (Backward Compatible) ===")

	// example stdout logger
	config := logger.JsonConfig{
		Levels:    "INFO,DEBUG",
		ApiLevels: "INFO,ERROR",
		NoColors:  false,
	}
	err := logger.SetupLogger(config)
	if err != nil {
		logger.Errorf("failed to setup logger: %v", err)
	}
	config.Output = "./stdout.log"
	config.Utc = true
	config.NoColors = true
	err = logger.SetupLogger(config)
	if err != nil {
		logger.Errorf("failed to setup file logger: %v", err)
	}
	logger.Debugf("this is a debug format int value %d in message.", 400)
	logger.Info("this is a basic info message from the logger.")
	logger.Api(200, "api call successful")
	logger.Api(400, "api call warning")
	logger.Api(500, "api call error")
	// logger.Fatal("this is a fatal message, the program will exit 1") // Commented out to prevent exit
}

func modernExample() {
	println("\n=== Modern Usage (Dependency Injection - Recommended) ===")

	// Create logger instance (no global state)
	config := logger.JsonConfig{
		Levels:    "INFO,DEBUG,WARNING,ERROR",
		ApiLevels: "INFO,ERROR,WARNING",
		NoColors:  false,
		Json:      false, // Set to true for JSON output
	}

	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
		return
	}

	// Basic logging
	log.Info("Basic info message")
	log.Debugf("Debug message with value: %d", 42)

	// Structured logging with key-value pairs
	log.Info("User action",
		"user_id", 123,
		"action", "login",
		"ip", "192.168.1.1",
		"duration", time.Millisecond*150,
	)

	// Context-aware logging
	ctx := context.WithValue(context.Background(), "request_id", "req-123")
	log.InfoContext(ctx, "Processing request",
		"endpoint", "/api/users",
		"method", "GET",
	)

	// Logger with additional context
	userLogger := log.With("user_id", 456, "session", "sess-789")
	userLogger.Info("User performed action", "action", "view_profile")
	userLogger.Warn("Suspicious activity detected", "reason", "multiple_failed_logins")

	// Grouped logging
	apiLogger := log.WithGroup("api")
	apiLogger.Info("Request received", "path", "/api/data", "method", "POST")
	apiLogger.Error("Database error", "table", "users", "error", "connection_timeout")

	// API logging with context
	log.APIContext(ctx, 200, "Request completed successfully",
		"response_time", time.Millisecond*50,
		"bytes_sent", 1024,
	)

	// Dependency injection example
	dependencyInjectionExampleWithLogger(log)
}

func compatibilityModeExample() {
	println("\n=== Compatibility Mode (Global State - Discouraged) ===")

	// Enable compatibility mode (creates global state)
	config := logger.JsonConfig{
		Levels:    "INFO,DEBUG",
		ApiLevels: "INFO,ERROR",
		NoColors:  false,
	}

	err := logger.EnableCompatibilityMode(config)
	if err != nil {
		fmt.Printf("failed to enable compatibility mode: %v\n", err)
		return
	}

	// Now global functions work (but this is discouraged)
	logger.Info("This uses global state (not recommended)")
	logger.Debugf("Global debug: %d", 42)

	// Clean up
	logger.DisableCompatibilityMode()
}

func dependencyInjectionExampleWithLogger(log logger.Logger) {
	println("\n=== Dependency Injection Example ===")

	// Simulate a service that uses dependency injection
	service := NewService(log)
	service.ProcessData("example data")
}

// Example service that uses dependency injection
type Service struct {
	logger logger.Logger
}

func NewService(logger logger.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) ProcessData(data string) {
	s.logger.Info("Processing data", "data", data, "service", "data-processor")

	// Simulate some processing
	time.Sleep(time.Millisecond * 10)

	s.logger.Info("Data processed successfully",
		"data", data,
		"result", "success",
		"processing_time", time.Millisecond*10,
	)
}
