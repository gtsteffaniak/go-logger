package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

// modernLogger implements the Logger interface with modern Go practices
type modernLogger struct {
	mu      sync.RWMutex
	configs []*LoggerConfig
	slog    *slog.Logger
}

// NewLogger creates a new Logger instance with modern features
func NewLogger(config JsonConfig) (Logger, error) {
	// Convert JsonConfig to LoggerConfig
	loggerConfig, err := convertJsonConfigToLoggerConfig(config)
	if err != nil {
		return nil, fmt.Errorf("convert config: %w", err)
	}

	// Create slog logger for structured logging
	var slogHandler slog.Handler
	var output io.Writer = os.Stdout

	if config.Output != "" && strings.ToUpper(config.Output) != "STDOUT" {
		file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		output = file
	}

	slogLevel := convertLogLevelsToSlogLevel(config.Levels)

	if config.Json {
		// Use JSON handler for JSON output
		slogHandler = slog.NewJSONHandler(output, &slog.HandlerOptions{
			AddSource: true,
			Level:     slogLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.SourceKey {
					source := a.Value.Any().(*slog.Source)
					source.File = stripProjectPath(source.File)
					source.Function = stripFunctionPath(source.Function)
				}
				return a
			},
		})
	} else {
		// Use custom handler for text output to maintain original format
		slogHandler = NewCustomHandler(output, slogLevel, loggerConfig)
	}

	// Create the logger instance with proper initialization
	loggerInstance, err := AddLogger(*loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("create logger instance: %w", err)
	}

	ml := &modernLogger{
		configs: []*LoggerConfig{loggerInstance},
		slog:    slog.New(slogHandler),
	}

	return ml, nil
}

// convertLogLevelsToSlogLevel converts the logger levels to slog level
func convertLogLevelsToSlogLevel(levels string) slog.Level {
	levelStrs := SplitByMultiple(levels)
	for _, levelStr := range levelStrs {
		upperLevel := strings.ToUpper(levelStr)
		switch upperLevel {
		case "DEBUG":
			return slog.LevelDebug
		case "INFO", "INFO ":
			return slog.LevelInfo
		case "WARNING", "WARN", "WARN ":
			return slog.LevelWarn
		case "ERROR":
			return slog.LevelError
		}
	}
	return slog.LevelInfo
}

// stripProjectPath removes the project path from file paths for cleaner output
func stripProjectPath(filePath string) string {
	// Simple implementation - just return the filename
	if idx := strings.LastIndex(filePath, "/"); idx != -1 {
		return filePath[idx+1:]
	}
	return filePath
}

// stripFunctionPath removes the package path from function names for cleaner output
func stripFunctionPath(functionPath string) string {
	// Remove package path and keep only the function name
	// Example: "github.com/gtsteffaniak/go-logger/logger.(*modernLogger).slogStructuredLog" -> "slogStructuredLog"
	if idx := strings.LastIndex(functionPath, "."); idx != -1 {
		return functionPath[idx+1:]
	}
	return functionPath
}

// convertJsonConfigToLoggerConfig converts JsonConfig to LoggerConfig
func convertJsonConfigToLoggerConfig(config JsonConfig) (*LoggerConfig, error) {
	upperLevels := []LogLevel{}
	for _, level := range SplitByMultiple(config.Levels) {
		if level == "" {
			break
		}
		upperLevel := strings.ToUpper(level)
		if upperLevel == "WARNING" || upperLevel == "WARN" {
			upperLevel = "WARN "
		}
		if upperLevel == "INFO" {
			upperLevel = "INFO "
		}
		level, ok := stringToLevel[upperLevel]
		if !ok {
			return nil, fmt.Errorf("invalid log level: %s", upperLevel)
		}
		upperLevels = append(upperLevels, level)
	}
	if len(upperLevels) == 0 {
		upperLevels = []LogLevel{INFO, ERROR, WARNING}
	}

	upperApiLevels := []LogLevel{}
	for _, level := range SplitByMultiple(config.ApiLevels) {
		if level == "" {
			break
		}
		upperLevel := strings.ToUpper(level)
		if upperLevel == "WARNING" || upperLevel == "WARN" {
			upperLevel = "WARN "
		}
		if upperLevel == "INFO" {
			upperLevel = "INFO "
		}
		level, ok := stringToLevel[strings.ToUpper(upperLevel)]
		if !ok {
			return nil, fmt.Errorf("invalid api log level: %s", upperLevel)
		}
		upperApiLevels = append(upperApiLevels, level)
	}
	if len(upperApiLevels) == 0 {
		upperApiLevels = []LogLevel{INFO, ERROR, WARNING}
	}

	outputStdout := strings.ToUpper(config.Output)
	if outputStdout == "STDOUT" {
		config.Output = ""
	}

	// JSON always enables structured logging, otherwise use the structured config (default: false)
	structuredOutput := config.Json || config.Structured

	return &LoggerConfig{
		Levels:       upperLevels,
		ApiLevels:    upperApiLevels,
		Stdout:       config.Output == "",
		Colors:       !config.NoColors,
		Disabled:     slices.Contains(upperLevels, DISABLED),
		DebugEnabled: slices.Contains(upperLevels, DEBUG),
		DisabledAPI:  slices.Contains(upperApiLevels, DISABLED),
		Utc:          config.Utc,
		FilePath:     config.Output,
		Structured:   structuredOutput,
		Json:         config.Json,
	}, nil
}

// Basic logging methods
func (ml *modernLogger) Debug(msg string, args ...any) {
	ml.logWithLevel(DEBUG, msg, false, false, args...)
}

func (ml *modernLogger) Info(msg string, args ...any) {
	ml.logWithLevel(INFO, msg, false, false, args...)
}

func (ml *modernLogger) Warn(msg string, args ...any) {
	ml.logWithLevel(WARNING, msg, false, false, args...)
}

func (ml *modernLogger) Error(msg string, args ...any) {
	ml.logWithLevel(ERROR, msg, false, false, args...)
}

func (ml *modernLogger) Fatal(msg string, args ...any) {
	ml.logWithLevel(FATAL, msg, false, false, args...)
}

// Formatted logging methods
func (ml *modernLogger) Debugf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ml.logWithLevel(DEBUG, msg, true, false)
}

func (ml *modernLogger) Infof(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ml.logWithLevel(INFO, msg, true, false)
}

func (ml *modernLogger) Warnf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ml.logWithLevel(WARNING, msg, true, false)
}

func (ml *modernLogger) Errorf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ml.logWithLevel(ERROR, msg, true, false)
}

func (ml *modernLogger) Fatalf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ml.logWithLevel(FATAL, msg, true, false)
}

// Context-aware methods
func (ml *modernLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	ml.logWithLevelAndContext(DEBUG, msg, false, false, ctx, args...)
}

func (ml *modernLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	ml.logWithLevelAndContext(INFO, msg, false, false, ctx, args...)
}

func (ml *modernLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	ml.logWithLevelAndContext(WARNING, msg, false, false, ctx, args...)
}

func (ml *modernLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	ml.logWithLevelAndContext(ERROR, msg, false, false, ctx, args...)
}

func (ml *modernLogger) FatalContext(ctx context.Context, msg string, args ...any) {
	ml.logWithLevelAndContext(FATAL, msg, false, false, ctx, args...)
}

// Formatted context-aware methods
func (ml *modernLogger) DebugfContext(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ml.logWithLevelAndContext(DEBUG, msg, true, false, ctx)
}

func (ml *modernLogger) InfofContext(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ml.logWithLevelAndContext(INFO, msg, true, false, ctx)
}

func (ml *modernLogger) WarnfContext(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ml.logWithLevelAndContext(WARNING, msg, true, false, ctx)
}

func (ml *modernLogger) ErrorfContext(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ml.logWithLevelAndContext(ERROR, msg, true, false, ctx)
}

func (ml *modernLogger) FatalfContext(ctx context.Context, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ml.logWithLevelAndContext(FATAL, msg, true, false, ctx)
}

// Structured logging
func (ml *modernLogger) With(args ...any) Logger {
	ml.mu.RLock()
	defer ml.mu.RUnlock()

	newLogger := &modernLogger{
		configs: make([]*LoggerConfig, len(ml.configs)),
		slog:    ml.slog.With(args...),
	}
	copy(newLogger.configs, ml.configs)
	return newLogger
}

func (ml *modernLogger) WithGroup(name string) Logger {
	ml.mu.RLock()
	defer ml.mu.RUnlock()

	newLogger := &modernLogger{
		configs: make([]*LoggerConfig, len(ml.configs)),
		slog:    ml.slog.WithGroup(name),
	}
	copy(newLogger.configs, ml.configs)
	return newLogger
}

// API logging
func (ml *modernLogger) API(statusCode int, msg string, args ...any) {
	ml.logAPI(statusCode, msg, false, args...)
}

func (ml *modernLogger) APIf(statusCode int, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ml.logAPI(statusCode, msg, true)
}

func (ml *modernLogger) APIContext(ctx context.Context, statusCode int, msg string, args ...any) {
	ml.logAPIWithContext(statusCode, msg, false, ctx, args...)
}

func (ml *modernLogger) APIfContext(ctx context.Context, statusCode int, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	ml.logAPIWithContext(statusCode, msg, true, ctx)
}

// Internal logging methods
func (ml *modernLogger) logWithLevel(level LogLevel, msg string, formatted bool, api bool, args ...any) {
	ml.mu.RLock()
	defer ml.mu.RUnlock()

	// Use structured logging if enabled and args are provided
	if len(args) > 0 && ml.configs[0].Structured {
		ml.slogStructuredLog(level, msg, args...)
		return
	}

	// Use slog for all logging if structured output is enabled (even without args)
	if ml.configs[0].Structured {
		ml.slogStructuredLog(level, msg, args...)
		return
	}

	// Fall back to original logging for backward compatibility
	levelStr := levelToString(level)
	color := getColorForLevel(level)

	for _, config := range ml.configs {
		if api {
			if config.DisabledAPI || !slices.Contains(config.ApiLevels, level) {
				continue
			}
		} else if level != FATAL {
			if config.Disabled || !slices.Contains(config.Levels, level) {
				continue
			}
		}

		ml.writeToConfig(config, levelStr, msg, formatted, api, color)
	}

	if level == FATAL {
		os.Exit(1)
	}
}

func (ml *modernLogger) logWithLevelAndContext(level LogLevel, msg string, formatted bool, api bool, ctx context.Context, args ...any) {
	ml.mu.RLock()
	defer ml.mu.RUnlock()

	// Use structured logging with context if enabled
	if ml.configs[0].Structured {
		ml.slogStructuredLogWithContext(ctx, level, msg, args...)
		return
	}

	// Fall back to context-less logging for backward compatibility
	ml.logWithLevel(level, msg, formatted, api)
}

func (ml *modernLogger) logAPI(statusCode int, msg string, formatted bool, args ...any) {
	level, _ := getAPILevelAndColor(statusCode)
	ml.logWithLevel(level, msg, formatted, true, args...)
}

func (ml *modernLogger) logAPIWithContext(statusCode int, msg string, formatted bool, ctx context.Context, args ...any) {
	level, _ := getAPILevelAndColor(statusCode)
	ml.logWithLevelAndContext(level, msg, formatted, true, ctx, args...)
}

func (ml *modernLogger) slogStructuredLog(level LogLevel, msg string, args ...any) {
	// Convert args to key-value pairs
	attrs := make([]any, 0, len(args))
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			attrs = append(attrs, args[i], args[i+1])
		}
	}

	// Use slog for structured logging
	switch level {
	case DEBUG:
		ml.slog.Debug(msg, attrs...)
	case INFO:
		ml.slog.Info(msg, attrs...)
	case WARNING:
		ml.slog.Warn(msg, attrs...)
	case ERROR:
		ml.slog.Error(msg, attrs...)
	case FATAL:
		ml.slog.Error(msg, attrs...)
		os.Exit(1)
	}
}

func (ml *modernLogger) slogStructuredLogWithContext(ctx context.Context, level LogLevel, msg string, args ...any) {
	// Convert args to key-value pairs
	attrs := make([]any, 0, len(args))
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			attrs = append(attrs, args[i], args[i+1])
		}
	}

	// Use slog for structured logging with context
	switch level {
	case DEBUG:
		ml.slog.DebugContext(ctx, msg, attrs...)
	case INFO:
		ml.slog.InfoContext(ctx, msg, attrs...)
	case WARNING:
		ml.slog.WarnContext(ctx, msg, attrs...)
	case ERROR:
		ml.slog.ErrorContext(ctx, msg, attrs...)
	case FATAL:
		ml.slog.ErrorContext(ctx, msg, attrs...)
		os.Exit(1)
	}
}

func (ml *modernLogger) writeToConfig(config *LoggerConfig, levelStr, msg string, formatted, api bool, color string) {
	writeOut := msg
	var formattedTime string
	if config.Utc {
		formattedTime = time.Now().UTC().Format("2006/01/02 15:04:05")
	} else {
		formattedTime = time.Now().Local().Format("2006/01/02 15:04:05")
	}

	if config.Colors && color != "" {
		formattedTime = formattedTime + color
	}

	if formatted || config.DebugEnabled {
		config.logger.SetPrefix(fmt.Sprintf("%s [%s] ", formattedTime, levelStr))
	} else {
		config.logger.SetPrefix(formattedTime + " ")
	}

	if config.Colors && color != "" {
		writeOut = writeOut + "\033[0m"
	}

	err := config.logger.Output(4, writeOut)
	if err != nil {
		log.Printf("failed to log message '%v' with error `%v`", msg, err)
	}
}

// Helper functions
func levelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO "
	case WARNING:
		return "WARN "
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	case API:
		return "API"
	default:
		return "UNKNOWN"
	}
}

func getColorForLevel(level LogLevel) string {
	switch level {
	case DEBUG:
		return GRAY
	case WARNING:
		return YELLOW
	case ERROR, FATAL:
		return RED
	case INFO:
		return ""
	default:
		return ""
	}
}

func getAPILevelAndColor(statusCode int) (LogLevel, string) {
	if statusCode > 304 && statusCode < 500 {
		return WARNING, YELLOW
	} else if statusCode >= 500 {
		return ERROR, RED
	} else {
		return INFO, GREEN
	}
}
