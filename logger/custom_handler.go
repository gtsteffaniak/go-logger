package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
)

// customHandler implements slog.Handler with the original logger format
type customHandler struct {
	writer io.Writer
	level  slog.Level
	config *LoggerConfig
	colors bool
	utc    bool
}

// NewCustomHandler creates a custom slog handler that mimics the original logger format
func NewCustomHandler(writer io.Writer, level slog.Level, config *LoggerConfig) *customHandler {
	return &customHandler{
		writer: writer,
		level:  level,
		config: config,
		colors: config.Colors,
		utc:    config.Utc,
	}
}

// Enabled returns true if the level is enabled
func (h *customHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle processes a log record
func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	// Format timestamp
	var timestamp string
	if h.utc {
		timestamp = r.Time.UTC().Format("2006/01/02 15:04:05")
	} else {
		timestamp = r.Time.Local().Format("2006/01/02 15:04:05")
	}

	// Format level
	levelStr := h.formatLevel(r.Level)

	// Format source location
	source := h.formatSource(r.PC)

	// Format message
	msg := r.Message

	// Format attributes
	attrs := h.formatAttrs(r)

	// Build the final message
	var finalMsg string
	if attrs != "" {
		finalMsg = fmt.Sprintf("%s %s", msg, attrs)
	} else {
		finalMsg = msg
	}

	// Apply colors if enabled
	if h.colors {
		color := h.getLevelColor(r.Level)
		if color != "" {
			timestamp = timestamp + color
			finalMsg = finalMsg + "\033[0m"
		}
	}

	// Write the log entry
	prefix := fmt.Sprintf("%s [%s] %s: ", timestamp, levelStr, source)
	_, err := fmt.Fprintf(h.writer, "%s%s\n", prefix, finalMsg)

	return err
}

// WithAttrs returns a new handler with additional attributes
func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we'll just return the same handler
	// In a more complex implementation, we'd store these attributes
	return h
}

// WithGroup returns a new handler with a group name
func (h *customHandler) WithGroup(name string) slog.Handler {
	// For simplicity, we'll just return the same handler
	// In a more complex implementation, we'd store the group name
	return h
}

// formatLevel formats the log level to match the original format
func (h *customHandler) formatLevel(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "ERROR"
	case level >= slog.LevelWarn:
		return "WARN "
	case level >= slog.LevelInfo:
		return "INFO "
	case level >= slog.LevelDebug:
		return "DEBUG"
	default:
		return "INFO "
	}
}

// formatSource formats the source location to match the original format
func (h *customHandler) formatSource(pc uintptr) string {
	if pc == 0 {
		return ""
	}

	// Get the call stack and skip frames to reach the actual application code
	// The call stack is: application code -> logger wrapper -> slog -> customHandler
	callers := make([]uintptr, 10)
	n := runtime.Callers(0, callers)

	// Skip the first few frames to get to the application code
	// Frame 0: formatSource
	// Frame 1: Handle
	// Frame 2: slog internal
	// Frame 3: logger wrapper (Infof, Debugf, etc.)
	// Frame 4: application code <- This is what we want
	skipFrames := 5
	if n > skipFrames {
		frames := runtime.CallersFrames(callers[skipFrames:])
		frame, _ := frames.Next()

		// Extract just the filename and line number
		file := frame.File
		if idx := strings.LastIndex(file, "/"); idx != -1 {
			file = file[idx+1:]
		}

		return fmt.Sprintf("%s:%d", file, frame.Line)
	}

	// Fallback to the original PC if we can't get enough frames
	frames := runtime.CallersFrames([]uintptr{pc})
	frame, _ := frames.Next()

	// Extract just the filename and line number
	file := frame.File
	if idx := strings.LastIndex(file, "/"); idx != -1 {
		file = file[idx+1:]
	}

	return fmt.Sprintf("%s:%d", file, frame.Line)
}

// formatAttrs formats the attributes as key-value pairs
func (h *customHandler) formatAttrs(r slog.Record) string {
	if r.NumAttrs() == 0 {
		return ""
	}

	var parts []string
	r.Attrs(func(attr slog.Attr) bool {
		parts = append(parts, fmt.Sprintf("%s=%v", attr.Key, attr.Value.Any()))
		return true
	})

	return strings.Join(parts, " ")
}

// getLevelColor returns the color code for a log level
func (h *customHandler) getLevelColor(level slog.Level) string {
	if !h.colors {
		return ""
	}

	switch {
	case level >= slog.LevelError:
		return RED
	case level >= slog.LevelWarn:
		return YELLOW
	case level >= slog.LevelInfo:
		return ""
	case level >= slog.LevelDebug:
		return GRAY
	default:
		return ""
	}
}
