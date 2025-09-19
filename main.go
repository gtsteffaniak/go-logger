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

	// JSON logging examples
	jsonLoggingExample()

	// Structured logging examples (text format)
	structuredLoggingExample()

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

func jsonLoggingExample() {
	println("\n=== JSON Logging Example ===")

	// Configure for JSON output
	config := logger.JsonConfig{
		Levels:    "INFO,DEBUG,WARNING,ERROR",
		ApiLevels: "INFO,ERROR,WARNING",
		Json:      true, // Enable JSON output
		Utc:       true, // Use UTC timestamps
	}

	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("failed to create JSON logger: %v\n", err)
		return
	}

	// Basic JSON logging
	log.Info("Application started", "version", "1.0.0", "environment", "production")

	// Structured JSON logging with key-value pairs
	log.Info("User login",
		"user_id", 12345,
		"email", "user@example.com",
		"ip_address", "192.168.1.100",
		"login_method", "oauth",
		"timestamp", time.Now().Unix(),
	)

	// Context-aware JSON logging
	ctx := context.WithValue(context.Background(), "request_id", "req-abc123")
	ctx = context.WithValue(ctx, "user_id", 12345)

	log.InfoContext(ctx, "Processing payment",
		"amount", 99.99,
		"currency", "USD",
		"payment_method", "credit_card",
		"transaction_id", "txn-xyz789",
	)

	// API logging with JSON
	log.APIContext(ctx, 200, "Payment processed successfully",
		"response_time_ms", 150,
		"bytes_sent", 1024,
	)

	// Error logging with detailed context
	log.Error("Database connection failed",
		"database", "users_db",
		"host", "db.example.com",
		"port", 5432,
		"error", "connection timeout",
		"retry_count", 3,
		"last_attempt", time.Now().Add(-time.Minute*5),
	)

	// Grouped logging for API endpoints
	apiLogger := log.WithGroup("api")
	apiLogger.Info("Request received",
		"method", "POST",
		"path", "/api/v1/users",
		"content_type", "application/json",
		"content_length", 1024,
	)
}

func structuredLoggingExample() {
	println("\n=== Structured Logging Example (Text Format) ===")

	// Configure for structured text output (not JSON)
	config := logger.JsonConfig{
		Levels:     "INFO,DEBUG,WARNING,ERROR",
		ApiLevels:  "INFO,ERROR,WARNING",
		Structured: true,  // Enable structured logging
		Json:       false, // Text format, not JSON
		NoColors:   false, // Enable colors for better readability
	}

	log, err := logger.NewLogger(config)
	if err != nil {
		fmt.Printf("failed to create structured logger: %v\n", err)
		return
	}

	// Basic structured logging
	log.Info("Server started",
		"port", 8080,
		"environment", "production",
		"version", "1.2.3",
		"pid", 12345,
	)

	// Business logic logging
	log.Info("Order created",
		"order_id", "ORD-12345",
		"customer_id", 67890,
		"total_amount", 149.99,
		"currency", "USD",
		"items_count", 3,
		"shipping_method", "express",
	)

	// Performance monitoring
	start := time.Now()
	// Simulate some operation
	time.Sleep(time.Millisecond * 50)
	duration := time.Since(start)

	log.Info("Database query completed",
		"query", "SELECT * FROM users WHERE active = true",
		"duration_ms", duration.Milliseconds(),
		"rows_returned", 1250,
		"cache_hit", false,
		"connection_pool_size", 10,
	)

	// Error logging with context
	log.Error("Failed to process payment",
		"payment_id", "pay-abc123",
		"error_code", "INSUFFICIENT_FUNDS",
		"error_message", "Account balance too low",
		"retry_after", "2025-09-19T10:00:00Z",
		"account_balance", 50.00,
		"required_amount", 99.99,
	)

	// Context-aware logging
	ctx := context.WithValue(context.Background(), "trace_id", "trace-xyz789")
	ctx = context.WithValue(ctx, "user_id", 12345)

	log.InfoContext(ctx, "User action performed",
		"action", "view_profile",
		"target_user_id", 67890,
		"timestamp", time.Now().Unix(),
		"session_duration", time.Minute*15,
	)

	// Service-specific logging
	userServiceLogger := log.With("service", "user-service", "version", "2.1.0")
	userServiceLogger.Info("User profile updated",
		"user_id", 12345,
		"fields_updated", []string{"email", "phone"},
		"updated_by", "self",
		"change_reason", "user_request",
	)

	// API endpoint logging
	apiLogger := log.WithGroup("api")
	apiLogger.Info("Request processed",
		"method", "GET",
		"path", "/api/v1/users/12345",
		"status_code", 200,
		"response_time_ms", 45,
		"user_agent", "MyApp/1.0",
		"request_size", 256,
		"response_size", 1024,
	)

	// Debug logging with detailed information
	log.Debug("Cache operation",
		"operation", "set",
		"key", "user:12345:profile",
		"ttl_seconds", 3600,
		"cache_size_mb", 128,
		"hit_rate", 0.85,
	)
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
