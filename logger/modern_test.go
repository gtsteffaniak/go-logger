package logger

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestModernLogger_StructuredLogging(t *testing.T) {
	var buf bytes.Buffer

	config := JsonConfig{
		Levels:    "INFO,DEBUG",
		ApiLevels: "INFO,ERROR",
		NoColors:  true,
		Json:      false,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Override the slog output for testing
	logger.(*modernLogger).slog = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Test structured logging
	logger.Info("User action",
		"user_id", 123,
		"action", "login",
		"ip", "192.168.1.1",
	)

	output := buf.String()
	if !strings.Contains(output, "User action") {
		t.Errorf("Expected 'User action' in output, got: %s", output)
	}
	if !strings.Contains(output, "user_id=123") {
		t.Errorf("Expected 'user_id=123' in output, got: %s", output)
	}
	if !strings.Contains(output, "action=login") {
		t.Errorf("Expected 'action=login' in output, got: %s", output)
	}
}

func TestModernLogger_ContextAwareLogging(t *testing.T) {
	var buf bytes.Buffer

	config := JsonConfig{
		Levels:    "INFO",
		ApiLevels: "INFO",
		NoColors:  true,
		Json:      false,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Override the slog output for testing
	logger.(*modernLogger).slog = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Test context-aware logging
	ctx := context.WithValue(context.Background(), "request_id", "req-123")
	logger.InfoContext(ctx, "Processing request",
		"endpoint", "/api/users",
		"method", "GET",
	)

	output := buf.String()
	if !strings.Contains(output, "Processing request") {
		t.Errorf("Expected 'Processing request' in output, got: %s", output)
	}
}

func TestModernLogger_WithMethod(t *testing.T) {
	var buf bytes.Buffer

	config := JsonConfig{
		Levels:    "INFO",
		ApiLevels: "INFO",
		NoColors:  true,
		Json:      false,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Override the slog output for testing
	logger.(*modernLogger).slog = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Test With method
	userLogger := logger.With("user_id", 456, "session", "sess-789")
	userLogger.Info("User performed action", "action", "view_profile")

	output := buf.String()
	if !strings.Contains(output, "User performed action") {
		t.Errorf("Expected 'User performed action' in output, got: %s", output)
	}
}

func TestModernLogger_WithGroup(t *testing.T) {
	var buf bytes.Buffer

	config := JsonConfig{
		Levels:    "INFO",
		ApiLevels: "INFO",
		NoColors:  true,
		Json:      false,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Override the slog output for testing
	logger.(*modernLogger).slog = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Test WithGroup method
	apiLogger := logger.WithGroup("api")
	apiLogger.Info("Request received", "path", "/api/data", "method", "POST")

	output := buf.String()
	if !strings.Contains(output, "Request received") {
		t.Errorf("Expected 'Request received' in output, got: %s", output)
	}
}

func TestModernLogger_API_Logging(t *testing.T) {
	var buf bytes.Buffer

	config := JsonConfig{
		Levels:    "INFO,WARNING,ERROR",
		ApiLevels: "INFO,WARNING,ERROR",
		NoColors:  true,
		Json:      false,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Override the slog output for testing
	logger.(*modernLogger).slog = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Test API logging
	logger.API(200, "Request successful", "response_time", time.Millisecond*50)

	output := buf.String()
	if !strings.Contains(output, "Request successful") {
		t.Errorf("Expected 'Request successful' in output, got: %s", output)
	}
}

func TestModernLogger_APIContext(t *testing.T) {
	var buf bytes.Buffer

	config := JsonConfig{
		Levels:    "INFO,WARNING,ERROR",
		ApiLevels: "INFO,WARNING,ERROR",
		NoColors:  true,
		Json:      false,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Override the slog output for testing
	logger.(*modernLogger).slog = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Test API context logging
	ctx := context.WithValue(context.Background(), "request_id", "req-456")
	logger.APIContext(ctx, 404, "Resource not found", "path", "/api/users/999")

	output := buf.String()
	if !strings.Contains(output, "Resource not found") {
		t.Errorf("Expected 'Resource not found' in output, got: %s", output)
	}
}

func TestModernLogger_BackwardCompatibility(t *testing.T) {
	config := JsonConfig{
		Levels:    "INFO,DEBUG",
		ApiLevels: "INFO,ERROR",
		NoColors:  true,
		Json:      false,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test that basic methods still work
	logger.Info("Basic info message")
	logger.Debugf("Debug message with value: %d", 42)
	logger.Warn("Warning message")
	logger.Error("Error message")

	// Test API methods
	logger.API(200, "API call successful")
	logger.APIf(404, "API call failed: %s", "not found")

	// Test that Fatal methods exist (but don't call them as they would exit)
	// logger.Fatalf("This would exit, but we're testing the method exists")
}

func TestCompatibilityLayer_GlobalFunctions(t *testing.T) {
	// Test that global functions work without setup (fallback behavior)
	// These should not panic even without a global logger set

	// Test global context functions (should fall back to legacy)
	ctx := context.WithValue(context.Background(), "test", "value")
	InfoContext(ctx, "Global context info", "key", "value")
	DebugfContext(ctx, "Global context debug: %s", "test")

	// Test global structured logging (should fall back to legacy)
	Info("Global structured", "key1", "value1", "key2", "value2")

	// Test global With functions (should return no-op logger)
	userLogger := With("user_id", 123)
	userLogger.Info("User action", "action", "test")

	// Test global WithGroup (should return no-op logger)
	apiLogger := WithGroup("api")
	apiLogger.Info("API action", "endpoint", "/test")
}

func TestNoOpLogger(t *testing.T) {
	// Test no-op logger when no global logger is set
	globalLogger = nil

	noOp := &noOpLogger{}

	// These should not panic
	noOp.Info("test")
	noOp.Debug("test")
	noOp.Warn("test")
	noOp.Error("test")
	noOp.Fatal("test")

	noOp.Infof("test %s", "value")
	noOp.Debugf("test %s", "value")
	noOp.Warnf("test %s", "value")
	noOp.Errorf("test %s", "value")
	noOp.Fatalf("test %s", "value")

	ctx := context.Background()
	noOp.InfoContext(ctx, "test")
	noOp.DebugContext(ctx, "test")
	noOp.WarnContext(ctx, "test")
	noOp.ErrorContext(ctx, "test")
	noOp.FatalContext(ctx, "test")

	noOp.InfofContext(ctx, "test %s", "value")
	noOp.DebugfContext(ctx, "test %s", "value")
	noOp.WarnfContext(ctx, "test %s", "value")
	noOp.ErrorfContext(ctx, "test %s", "value")
	noOp.FatalfContext(ctx, "test %s", "value")

	noOp.API(200, "test")
	noOp.APIf(200, "test %s", "value")
	noOp.APIContext(ctx, 200, "test")
	noOp.APIfContext(ctx, 200, "test %s", "value")

	// With methods should return the same no-op logger
	result := noOp.With("key", "value")
	if result != noOp {
		t.Error("With should return the same no-op logger")
	}

	result = noOp.WithGroup("group")
	if result != noOp {
		t.Error("WithGroup should return the same no-op logger")
	}
}

func TestDependencyInjection(t *testing.T) {
	// Test dependency injection pattern
	config := JsonConfig{
		Levels:    "INFO,ERROR",
		ApiLevels: "ERROR",
		NoColors:  true,
		Json:      false,
	}

	serviceLogger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create service logger: %v", err)
	}

	// Test service with dependency injection
	service := &TestService{logger: serviceLogger}
	service.ProcessData("test data")

	// This should not panic
	service.ProcessDataWithContext(context.Background(), "test data with context")
}

type TestService struct {
	logger Logger
}

func (s *TestService) ProcessData(data string) {
	s.logger.Info("Processing data", "data", data, "service", "test-service")
	s.logger.Info("Data processed", "data", data, "result", "success")
}

func (s *TestService) ProcessDataWithContext(ctx context.Context, data string) {
	s.logger.InfoContext(ctx, "Processing data with context", "data", data, "service", "test-service")
	s.logger.InfoContext(ctx, "Data processed with context", "data", data, "result", "success")
}

func TestErrorHandling(t *testing.T) {
	// Test error handling in configuration
	invalidConfig := JsonConfig{
		Levels:    "INVALID_LEVEL",
		ApiLevels: "INFO",
		NoColors:  true,
	}

	_, err := NewLogger(invalidConfig)
	if err == nil {
		t.Error("Expected error for invalid log level, got nil")
	}

	// Test error handling in API levels
	invalidApiConfig := JsonConfig{
		Levels:    "INFO",
		ApiLevels: "INVALID_API_LEVEL",
		NoColors:  true,
	}

	_, err = NewLogger(invalidApiConfig)
	if err == nil {
		t.Error("Expected error for invalid API log level, got nil")
	}
}

func TestThreadSafety(t *testing.T) {
	// Test that the logger is thread-safe
	config := JsonConfig{
		Levels:    "INFO",
		ApiLevels: "INFO",
		NoColors:  true,
		Json:      false,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Run concurrent logging
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				logger.Info("Concurrent log", "goroutine", id, "iteration", j)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// If we get here without panicking, thread safety is working
}
