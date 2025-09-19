# Migration Guide: From Global State to Dependency Injection

## Overview

This guide helps you migrate from the global logger pattern to the recommended dependency injection pattern.

## Why Migrate?

The global logger pattern has several drawbacks:
- **Testing difficulties**: Hard to mock or control logger behavior in tests
- **Concurrency issues**: Global state can cause race conditions
- **Hidden dependencies**: Code dependencies are not explicit
- **Configuration conflicts**: Multiple parts of code can't have different logger configs

## Migration Steps

### Step 1: Current Code (Global Pattern)
```go
package main

import "github.com/gtsteffaniak/go-logger/logger"

func main() {
    // Setup global logger
    config := logger.JsonConfig{
        Levels: "INFO,DEBUG",
        NoColors: false,
    }
    logger.SetupLoggerWithModern(config)
    
    // Use global functions
    logger.Info("This is a message")
    logger.Debugf("Debug value: %d", 42)
    
    // In other functions
    doSomething()
}

func doSomething() {
    logger.Info("Doing something") // Hidden dependency!
}
```

### Step 2: Migrated Code (Dependency Injection)
```go
package main

import "github.com/gtsteffaniak/go-logger/logger"

func main() {
    // Create logger instance
    config := logger.JsonConfig{
        Levels: "INFO,DEBUG",
        NoColors: false,
    }
    log := logger.NewLogger(config)
    
    // Pass logger explicitly
    log.Info("This is a message")
    log.Debugf("Debug value: %d", 42)
    
    // Pass logger to functions
    doSomething(log)
}

func doSomething(log logger.Logger) {
    log.Info("Doing something") // Explicit dependency!
}
```

### Step 3: Service Pattern
```go
type MyService struct {
    logger logger.Logger
}

func NewMyService(logger logger.Logger) *MyService {
    return &MyService{logger: logger}
}

func (s *MyService) ProcessData(data string) {
    s.logger.Info("Processing data", "data", data)
    // ... processing logic
    s.logger.Info("Data processed successfully")
}
```

## Backward Compatibility

The global functions will continue to work, but consider them deprecated:

```go
// ❌ Deprecated (but still works)
logger.Info("message")

// ✅ Recommended
log := logger.NewLogger(config)
log.Info("message")
```

## Testing Benefits

With dependency injection, testing becomes much easier:

```go
func TestMyService(t *testing.T) {
    // Create a test logger that captures output
    testLogger := logger.NewTestLogger()
    
    service := NewMyService(testLogger)
    service.ProcessData("test")
    
    // Assert on logger output
    assert.Contains(t, testLogger.Output(), "Processing data")
}
```

## Gradual Migration Strategy

1. **Phase 1**: Keep using global functions, but start passing logger instances to new code
2. **Phase 2**: Migrate high-level functions to accept logger parameters
3. **Phase 3**: Remove global logger usage entirely

## Configuration Per Component

With dependency injection, different components can have different logger configurations:

```go
// API logger with JSON output
apiLogger := logger.NewLogger(logger.JsonConfig{
    Levels: "INFO,ERROR",
    Json: true,
})

// Debug logger with verbose output
debugLogger := logger.NewLogger(logger.JsonConfig{
    Levels: "DEBUG,INFO,WARNING,ERROR",
    Json: false,
})

// Use appropriate logger for each component
apiHandler := NewAPIHandler(apiLogger)
debugService := NewDebugService(debugLogger)
```
