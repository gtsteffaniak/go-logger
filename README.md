# go-logger
Modern, feature-rich logger for Go applications with backward compatibility

## Features

- ✅ **Structured Logging** - Key-value pair logging with `slog` compatibility
- ✅ **Context Support** - Context-aware logging for request tracing
- ✅ **Dependency Injection** - Interface-based design for testability
- ✅ **Thread Safe** - Concurrent access protection
- ✅ **Multiple Outputs** - stdout, files, JSON, text formats
- ✅ **Color Coding** - Terminal color support
- ✅ **API Logging** - HTTP status code-based logging
- ✅ **Level Filtering** - Configurable log levels
- ✅ **Grouped Logging** - Hierarchical log organization

## Quick Start

### Recommended: Dependency Injection (New Code)

```go
package main

import (
    "context"
    "github.com/gtsteffaniak/go-logger/logger"
)

func main() {
    // Create logger instance (no global state)
    config := logger.JsonConfig{
        Levels:    "INFO,DEBUG",
        ApiLevels: "INFO,ERROR",
        NoColors:  false,
    }

    log, err := logger.NewLogger(config)
    if err != nil {
        panic(err)
    }

    // Basic logging
    log.Info("Hello, World!")
    log.Debugf("Debug value: %d", 42)

    // Structured logging
    log.Info("User action", "user_id", 123, "action", "login")

    // Context-aware logging
    ctx := context.WithValue(context.Background(), "request_id", "req-123")
    log.InfoContext(ctx, "Processing request", "endpoint", "/api/users")

    // Pass logger to services
    service := NewService(log)
    service.ProcessData("example")
}

type Service struct {
    logger logger.Logger
}

func NewService(logger logger.Logger) *Service {
    return &Service{logger: logger}
}

func (s *Service) ProcessData(data string) {
    s.logger.Info("Processing data", "data", data)
}
```

### 🔄 Compatibility Mode (Existing Code)

For existing code that uses global functions, enable compatibility mode:

```go
package main

import "github.com/gtsteffaniak/go-logger/logger"

func main() {
    // Enable compatibility mode (maintains global state for backward compatibility)
    config := logger.JsonConfig{
        Levels:    "INFO,DEBUG",
        ApiLevels: "INFO,ERROR",
        NoColors:  false,
    }

    err := logger.EnableCompatibilityMode(config)
    if err != nil {
        logger.Errorf("failed to setup logger: %v", err)
    }

    // All existing code works unchanged
    logger.Info("Hello, World!")
    logger.Debugf("Debug value: %d", 42)
    logger.Api(200, "API call successful")

    // Clean up when done (optional)
    logger.DisableCompatibilityMode()
}
```

### Legacy Usage (Deprecated)

```go
package main

import (
    "context"
    "github.com/gtsteffaniak/go-logger/logger"
)

func main() {
    // Setup modern logger
    config := logger.JsonConfig{
        Levels:    "INFO,DEBUG,WARNING,ERROR",
        ApiLevels: "INFO,ERROR,WARNING",
        Json:      false, // Set to true for JSON output
    }

    err := logger.SetupLoggerWithModern(config)
    if err != nil {
        logger.Errorf("failed to setup logger: %v", err)
    }

    // Structured logging
    logger.Info("User action",
        "user_id", 123,
        "action", "login",
        "ip", "192.168.1.1",
    )

    // Context-aware logging
    ctx := context.WithValue(context.Background(), "request_id", "req-123")
    logger.InfoContext(ctx, "Processing request",
        "endpoint", "/api/users",
        "method", "GET",
    )

    // Logger with additional context
    userLogger := logger.With("user_id", 456, "session", "sess-789")
    userLogger.Info("User performed action", "action", "view_profile")

    // Grouped logging
    apiLogger := logger.WithGroup("api")
    apiLogger.Info("Request received", "path", "/api/data", "method", "POST")
}
```

### Dependency Injection

```go
type Service struct {
    logger logger.Logger
}

func NewService(logger logger.Logger) *Service {
    return &Service{logger: logger}
}

func (s *Service) ProcessData(data string) {
    s.logger.Info("Processing data", "data", data, "service", "data-processor")
    // ... processing logic
    s.logger.Info("Data processed successfully",
        "data", data,
        "result", "success",
    )
}

// Usage
config := logger.JsonConfig{Levels: "INFO,ERROR"}
serviceLogger, err := logger.NewLogger(config)
if err != nil {
    log.Fatal(err)
}

service := NewService(serviceLogger)
service.ProcessData("example data")
```

## Configuration

```go
type JsonConfig struct {
    Levels     string `json:"levels"`     // "INFO|DEBUG|WARNING|ERROR"
    ApiLevels  string `json:"apiLevels"`  // "INFO|ERROR|WARNING"
    Output     string `json:"output"`     // "stdout" or "/path/to/file.log"
    NoColors   bool   `json:"noColors"`   // disable colors
    Json       bool   `json:"json"`       // JSON output format (enables structured logging)
    Structured bool   `json:"structured"` // enable structured logging (default: false)
    Utc        bool   `json:"utc"`        // UTC timestamps
}
```

### Structured Logging

The logger supports structured logging with key-value pairs. When `structured: true` or `json: true` is enabled, the logger uses Go's `slog` package for consistent structured output.

```go
// Enable structured logging
config := logger.JsonConfig{
    Levels:     "INFO,DEBUG",
    Structured: true,  // Enable structured logging
    NoColors:   true,
}

log, _ := logger.NewLogger(config)

// Structured logging with key-value pairs
log.Info("User action", "user_id", 123, "action", "login", "ip", "192.168.1.1")

// JSON output (automatically enables structured logging)
config := logger.JsonConfig{
    Levels: "INFO,DEBUG",
    Json:   true,  // JSON output + structured logging
}
```

**Note:** When `json: true` is set, structured logging is automatically enabled regardless of the `structured` setting.

## Migration Guide

### Why Migrate?

The global logger pattern has several drawbacks:
- **Testing difficulties**: Hard to mock or control logger behavior in tests
- **Concurrency issues**: Global state can cause race conditions
- **Hidden dependencies**: Code dependencies are not explicit
- **Configuration conflicts**: Multiple parts of code can't have different logger configs

### Migration Strategy

#### Phase 1: Enable Compatibility Mode (Immediate - No Breaking Changes)

Replace your existing `SetupLogger()` calls with `EnableCompatibilityMode()`:

```go
// Before
err := logger.SetupLogger(config)

// After
err := logger.EnableCompatibilityMode(config)
```

This gives you:
- ✅ All existing code works unchanged
- ✅ Access to modern features (structured logging, context support)
- ✅ Deprecation warnings to guide future migration

#### Phase 2: Gradual Migration (Recommended)

Start using dependency injection in new code while keeping existing code working:

```go
// New code - use dependency injection
func NewAPIHandler(logger logger.Logger) *APIHandler {
    return &APIHandler{logger: logger}
}

// Existing code - continues to work with compatibility mode
logger.Info("This still works")
```

#### Phase 3: Full Migration (Optional)

Replace global functions with explicit logger instances:

```go
// Before (global state)
func processData() {
    logger.Info("Processing data")
}

// After (dependency injection)
func processData(log logger.Logger) {
    log.Info("Processing data")
}
```

### Migration Examples

#### Simple Migration
```go
// Before
package main

import "github.com/gtsteffaniak/go-logger/logger"

func main() {
    config := logger.JsonConfig{Levels: "INFO,DEBUG"}
    logger.SetupLogger(config)

    logger.Info("Hello World")
    doSomething()
}

func doSomething() {
    logger.Info("Doing something")
}

// After
package main

import "github.com/gtsteffaniak/go-logger/logger"

func main() {
    config := logger.JsonConfig{Levels: "INFO,DEBUG"}
    err := logger.EnableCompatibilityMode(config)
    if err != nil {
        panic(err)
    }

    logger.Info("Hello World") // Still works!
    doSomething() // Still works!
}

func doSomething() {
    logger.Info("Doing something") // Still works!
}
```

#### Advanced Migration with Dependency Injection
```go
// Before
type Service struct {
    // other fields
}

func (s *Service) ProcessData(data string) {
    logger.Info("Processing data", "data", data)
    // ... processing
    logger.Info("Data processed")
}

// After
type Service struct {
    logger logger.Logger
    // other fields
}

func NewService(logger logger.Logger) *Service {
    return &Service{logger: logger}
}

func (s *Service) ProcessData(data string) {
    s.logger.Info("Processing data", "data", data)
    // ... processing
    s.logger.Info("Data processed")
}

// Usage
func main() {
    config := logger.JsonConfig{Levels: "INFO,DEBUG"}
    log, err := logger.NewLogger(config)
    if err != nil {
        panic(err)
    }

    service := NewService(log)
    service.ProcessData("example")
}
```

## Best Practices

### ✅ Do
- Use dependency injection for new code
- Pass logger instances explicitly to functions
- Use structured logging with key-value pairs
- Leverage context for request tracing
- Use different logger configs for different components

### ❌ Don't
- Rely on global state in new code
- Mix global and instance-based logging in the same codebase
- Ignore deprecation warnings
- Use global functions in library code

## Troubleshooting

### Deprecation Warnings
If you see deprecation warnings, you're using global functions. Consider migrating to dependency injection:

```go
// This shows warnings
logger.Info("message")

// This doesn't
log := logger.NewLogger(config)
log.Info("message")
```

### Compatibility Mode Not Working
Make sure you're calling `EnableCompatibilityMode()` instead of `SetupLogger()`:

```go
// Wrong
err := logger.SetupLogger(config)

// Correct
err := logger.EnableCompatibilityMode(config)
```

## Examples

### JSON Logging

JSON logging provides structured output that's perfect for log aggregation systems like ELK stack, Fluentd, or cloud logging services.

```go
package main

import (
    "context"
    "github.com/gtsteffaniak/go-logger/logger"
)

func main() {
    // Configure for JSON output
    config := logger.JsonConfig{
        Levels:    "INFO,DEBUG,WARNING,ERROR",
        ApiLevels: "INFO,ERROR,WARNING",
        Json:      true,  // Enable JSON output
        Utc:       true,  // Use UTC timestamps
    }

    log, err := logger.NewLogger(config)
    if err != nil {
        panic(err)
    }

    // Basic JSON logging
    log.Info("Application started")
    // Output: {"time":"2025-09-18T19:14:56Z","level":"INFO","msg":"Application started"}

    // Structured JSON logging with key-value pairs
    log.Info("User login",
        "user_id", 12345,
        "email", "user@example.com",
        "ip_address", "192.168.1.100",
        "user_agent", "Mozilla/5.0...",
        "login_method", "oauth",
    )
    // Output: {"time":"2025-09-18T19:14:56Z","level":"INFO","msg":"User login","user_id":12345,"email":"user@example.com","ip_address":"192.168.1.100","user_agent":"Mozilla/5.0...","login_method":"oauth"}

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

    // Error logging with stack trace information
    log.Error("Database connection failed",
        "database", "users_db",
        "host", "db.example.com",
        "port", 5432,
        "error", "connection timeout",
        "retry_count", 3,
    )

    // Grouped logging for API endpoints
    apiLogger := log.WithGroup("api")
    apiLogger.Info("Request received",
        "method", "POST",
        "path", "/api/v1/users",
        "content_type", "application/json",
    )
}
```

### Structured Logging (Text Format)

Structured logging in text format provides human-readable output while maintaining key-value structure.

```go
package main

import (
    "context"
    "time"
    "github.com/gtsteffaniak/go-logger/logger"
)

func main() {
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
        panic(err)
    }

    // Basic structured logging
    log.Info("Server started",
        "port", 8080,
        "environment", "production",
        "version", "1.2.3",
    )
    // Output: 2025/09/18 19:14:56 [INFO ] Server started port=8080 environment=production version=1.2.3

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
    // ... some operation ...
    duration := time.Since(start)
    
    log.Info("Database query completed",
        "query", "SELECT * FROM users WHERE active = true",
        "duration_ms", duration.Milliseconds(),
        "rows_returned", 1250,
        "cache_hit", false,
    )

    // Error logging with context
    log.Error("Failed to process payment",
        "payment_id", "pay-abc123",
        "error_code", "INSUFFICIENT_FUNDS",
        "error_message", "Account balance too low",
        "retry_after", "2025-09-19T10:00:00Z",
    )

    // Context-aware logging
    ctx := context.WithValue(context.Background(), "trace_id", "trace-xyz789")
    ctx = context.WithValue(ctx, "user_id", 12345)
    
    log.InfoContext(ctx, "User action performed",
        "action", "view_profile",
        "target_user_id", 67890,
        "timestamp", time.Now().Unix(),
    )

    // Service-specific logging
    userServiceLogger := log.With("service", "user-service", "version", "2.1.0")
    userServiceLogger.Info("User profile updated",
        "user_id", 12345,
        "fields_updated", []string{"email", "phone"},
        "updated_by", "self",
    )

    // API endpoint logging
    apiLogger := log.WithGroup("api")
    apiLogger.Info("Request processed",
        "method", "GET",
        "path", "/api/v1/users/12345",
        "status_code", 200,
        "response_time_ms", 45,
        "user_agent", "MyApp/1.0",
    )
}
```

### Comparison: JSON vs Structured Text

```go
package main

import (
    "github.com/gtsteffaniak/go-logger/logger"
)

func main() {
    // JSON Configuration
    jsonConfig := logger.JsonConfig{
        Levels: "INFO,DEBUG",
        Json:   true,
        Utc:    true,
    }
    
    jsonLog, _ := logger.NewLogger(jsonConfig)
    jsonLog.Info("User action", "user_id", 123, "action", "login")
    // Output: {"time":"2025-09-18T19:14:56Z","level":"INFO","msg":"User action","user_id":123,"action":"login"}

    // Structured Text Configuration
    textConfig := logger.JsonConfig{
        Levels:     "INFO,DEBUG",
        Structured: true,
        Json:       false,
        NoColors:   false,
    }
    
    textLog, _ := logger.NewLogger(textConfig)
    textLog.Info("User action", "user_id", 123, "action", "login")
    // Output: 2025/09/18 19:14:56 [INFO ] User action user_id=123 action=login
}
```

### File Output Examples

```go
package main

import (
    "github.com/gtsteffaniak/go-logger/logger"
)

func main() {
    // JSON logging to file
    jsonConfig := logger.JsonConfig{
        Levels: "INFO,DEBUG,WARNING,ERROR",
        Json:   true,
        Output: "/var/log/app.json",  // Log to file
        Utc:    true,
    }
    
    jsonLog, _ := logger.NewLogger(jsonConfig)
    jsonLog.Info("Application started", "version", "1.0.0")

    // Structured text logging to file
    textConfig := logger.JsonConfig{
        Levels:     "INFO,DEBUG,WARNING,ERROR",
        Structured: true,
        Json:       false,
        Output:     "/var/log/app.log",  // Log to file
        NoColors:   true,  // No colors in file output
    }
    
    textLog, _ := logger.NewLogger(textConfig)
    textLog.Info("Application started", "version", "1.0.0")
}
```

See [main.go](./main.go) for comprehensive examples of both legacy and modern usage patterns.

## Linting

```bash
go mod tidy
go tool golangci-lint run
```

## Testing

With dependency injection, testing becomes much easier:

```go
func TestMyService(t *testing.T) {
    // Create a test logger with specific configuration
    config := logger.JsonConfig{
        Levels: "DEBUG,INFO,WARNING,ERROR",
        Json: true,
    }

    testLogger, err := logger.NewLogger(config)
    if err != nil {
        t.Fatal(err)
    }

    service := NewService(testLogger)
    service.ProcessData("test data")

    // You can now easily:
    // - Assert on logger output
    // - Use different logger configs for different tests
    // - Mock the logger interface if needed
}
```

Run tests:
```bash
go test ./...
```
