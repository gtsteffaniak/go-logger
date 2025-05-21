package logger

import (
	"bytes"
	"fmt"
	"log"
	"regexp" // Used by logger package
	"strings"
	"testing"
	// Logger uses time
)

// setupForPrimaryLoggerTest (no changes from your previous version)
func setupForPrimaryLoggerTest(t *testing.T, buf *bytes.Buffer, activeLevels []LogLevel, noColors bool) {
	t.Helper()
	savedLoggers := append([]*Logger(nil), loggers...)
	savedStdOutLoggerExists := stdOutLoggerExists
	t.Cleanup(func() {
		loggers = savedLoggers
		stdOutLoggerExists = savedStdOutLoggerExists
	})
	loggers = nil
	stdOutLoggerExists = false
	l, err := NewLogger("", activeLevels, []LogLevel{}, noColors)
	if err != nil {
		t.Fatalf("Failed to create test logger: %v", err)
	}
	l.logger.SetOutput(buf)
	loggers = []*Logger{l}
}

// setupForFallbackTest (no changes from your previous version)
func setupForFallbackTest(t *testing.T, buf *bytes.Buffer) {
	t.Helper()
	savedLoggers := append([]*Logger(nil), loggers...)
	savedStdOutLoggerExists := stdOutLoggerExists
	// Assuming direct control via log.SetOutput for fallback testing.
	// For a more robust save/restore of global log output, os.Stderr could be an assumption or a more complex capture.
	currentGlobalLogOutput := log.Writer()
	savedGlobalLogFlags := log.Flags()
	t.Cleanup(func() {
		loggers = savedLoggers
		stdOutLoggerExists = savedStdOutLoggerExists
		log.SetOutput(currentGlobalLogOutput)
		log.SetFlags(savedGlobalLogFlags)
	})
	loggers = nil
	stdOutLoggerExists = false
	log.SetOutput(buf)
	log.SetFlags(0)
}

// --- New/Updated Regex Patterns ---
// These patterns now aim to capture the message part in group 1 (or the last group).

// For INFO (not in debug mode): "YYYY/MM/DD HH:MM:SS MESSAGE\n"
var infoPrimaryLogPattern_NonDebug = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s(.*\n)$`)

// For INFO (when logger is in debug mode): "YYYY/MM/DD HH:MM:SS [INFO ] file:line: MESSAGE\n"
var infoPrimaryLogPattern_DebugMode = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s\[INFO\s+\]\s[a-zA-Z0-9._/-]+:\d+:\s(.*\n)$`)

// For DEBUG: "YYYY/MM/DD HH:MM:SS [DEBUG] file:line: MESSAGE\n"
var debugPrimaryLogPattern = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s\[DEBUG\]\s[a-zA-Z0-9._/-]+:\d+:\s(.*\n)$`)

// For WARNING (not in debug mode): "YYYY/MM/DD HH:MM:SS [WARN ] MESSAGE\n"
var warnPrimaryLogPattern_NonDebug = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s\[WARN\s+\]\s(.*\n)$`)

// For ERROR (not in debug mode): "YYYY/MM/DD HH:MM:SS [ERROR] MESSAGE\n"
var errorPrimaryLogPattern_NonDebug = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s\[ERROR\]\s(.*\n)$`)

// Helper function to extract message using regex
func extractMessage(t *testing.T, output string, pattern *regexp.Regexp) string {
	t.Helper()
	matches := pattern.FindStringSubmatch(output)
	if len(matches) < 2 { // Ensure pattern matched and captured the message group
		t.Fatalf("Log output did not match expected pattern '%s'. Full output: '%s'", pattern.String(), output)
		return "" // Should not reach here due to t.Fatalf
	}
	return matches[1] // Captured message is the last group
}

func TestInfo_PrimaryLogger(t *testing.T) {
	var buf bytes.Buffer
	noColors := true

	t.Run("InfoNonDebugMode", func(t *testing.T) {
		setupForPrimaryLoggerTest(t, &buf, []LogLevel{INFO}, noColors)
		buf.Reset()

		Info("Hello %s, number %d", "Alice", 100)
		expectedMessage := "Hello Alice, number 100\n"
		actualOutput := buf.String()
		actualMessage := extractMessage(t, actualOutput, infoPrimaryLogPattern_NonDebug)

		if actualMessage != expectedMessage {
			t.Errorf("Expected log message '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
		}
	})

	t.Run("InfoWithLoggerDebugMode", func(t *testing.T) {
		setupForPrimaryLoggerTest(t, &buf, []LogLevel{INFO, DEBUG}, noColors)
		buf.Reset()

		Info("Hello %s, debug number %d", "Bob", 200)
		expectedMessage := "Hello Bob, debug number 200\n"
		actualOutput := buf.String()
		actualMessage := extractMessage(t, actualOutput, infoPrimaryLogPattern_DebugMode)

		if actualMessage != expectedMessage {
			t.Errorf("Expected log message (debug mode) '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
		}
	})
}

func TestInfo_FallbackLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForFallbackTest(t, &buf)

	Info("Message for %s", "fallback_user")
	expectedOutput := "[DEBUG] Message for fallback_user\n" // Note: current code has Info fallback logging as [DEBUG]
	actualOutput := buf.String()

	if actualOutput != expectedOutput {
		t.Errorf("Expected fallback log output '%s', got '%s'", expectedOutput, actualOutput)
	}

	buf.Reset()
	Info("Simple info fallback")
	expectedSimple := "[DEBUG] Simple info fallback\n"
	actualSimple := buf.String()
	if actualSimple != expectedSimple {
		t.Errorf("Expected simple fallback log output '%s', got '%s'", expectedSimple, actualSimple)
	}
}

func TestDebug_PrimaryLogger(t *testing.T) {
	var buf bytes.Buffer
	// DEBUG level needs to be active for Debug() to log.
	// Setting DEBUG level also activates Lshortfile and logger.debugEnabled.
	setupForPrimaryLoggerTest(t, &buf, []LogLevel{DEBUG}, true) // noColors = true

	Debug("Processing %s, item %d", "data_set", 77)
	expectedMessage := "Processing data_set, item 77\n"
	actualOutput := buf.String()
	actualMessage := extractMessage(t, actualOutput, debugPrimaryLogPattern)

	if actualMessage != expectedMessage {
		t.Errorf("Expected log message '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
	}
}

func TestDebug_FallbackLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForFallbackTest(t, &buf)

	Debug("Fallback debug %s", "message")
	expectedOutput := "[DEBUG] Fallback debug message\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected fallback log output '%s', got '%s'", expectedOutput, actualOutput)
	}
}

func TestWarning_PrimaryLogger(t *testing.T) {
	var buf bytes.Buffer
	// Test Warning when logger is NOT in debug mode (DEBUG level not active)
	setupForPrimaryLoggerTest(t, &buf, []LogLevel{WARNING, INFO}, true) // noColors = true

	Warning("Potential issue with %s", "config_value")
	expectedMessage := "Potential issue with config_value\n"
	actualOutput := buf.String()
	actualMessage := extractMessage(t, actualOutput, warnPrimaryLogPattern_NonDebug)

	if actualMessage != expectedMessage {
		t.Errorf("Expected log message '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
	}
}

func TestWarning_FallbackLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForFallbackTest(t, &buf)

	Warning("Fallback warning: %s", "check this")
	expectedOutput := "[WARN ] Fallback warning: check this\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected fallback log output '%s', got '%s'", expectedOutput, actualOutput)
	}
}

func TestError_PrimaryLogger(t *testing.T) {
	var buf bytes.Buffer
	// Test Error when logger is NOT in debug mode
	setupForPrimaryLoggerTest(t, &buf, []LogLevel{ERROR, INFO}, true) // noColors = true

	errVal := fmt.Errorf("critical failure")
	Error("System error: %v", errVal)
	expectedMessage := fmt.Sprintf("System error: %v\n", errVal)
	actualOutput := buf.String()
	actualMessage := extractMessage(t, actualOutput, errorPrimaryLogPattern_NonDebug)

	if actualMessage != expectedMessage {
		t.Errorf("Expected log message '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
	}
}

func TestError_FallbackLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForFallbackTest(t, &buf)

	Error("Fallback error: %s", "system down")
	expectedOutput := "[ERROR] Fallback error: system down\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected fallback log output '%s', got '%s'", expectedOutput, actualOutput)
	}
}

func TestFormatting_VariousArgTypes_Primary(t *testing.T) {
	var buf bytes.Buffer
	setupForPrimaryLoggerTest(t, &buf, []LogLevel{INFO}, true) // noColors = true

	type MyStruct struct{ Name string }
	s := MyStruct{Name: "DataObject"}
	ptr := &s

	Info("String: '%s', Int: %d, Bool: %t, Float: %.2f, Struct: %v, StructPtr: %v, PtrAddr: %p",
		"test string", 987, false, 123.4567, s, ptr, ptr)

	output := buf.String()
	actualMessage := extractMessage(t, output, infoPrimaryLogPattern_NonDebug)

	// %p output for PtrAddr is environment-dependent.
	// We verify the known parts of the string.
	// Create a regex that matches the message, allowing for the variable pointer address.
	expectedMessagePatternStr := fmt.Sprintf("^String: 'test string', Int: 987, Bool: false, Float: 123.46, Struct: %v, StructPtr: %v, PtrAddr: 0x[0-9a-f]+\n$", s, ptr)
	expectedMessagePattern := regexp.MustCompile(expectedMessagePatternStr)

	if !expectedMessagePattern.MatchString(actualMessage) {
		t.Errorf("Formatted message for various types did not match expected pattern.\nExpected pattern: %s\nActual message: %sFull output: %s",
			expectedMessagePattern.String(), actualMessage, output)
	}
}

func TestLogging_LevelNotActive(t *testing.T) {
	var buf bytes.Buffer
	setupForPrimaryLoggerTest(t, &buf, []LogLevel{WARNING}, true) // Only WARNING is active

	Info("This INFO message should not appear")
	Debug("This DEBUG message should not appear")
	Error("This ERROR message should not appear")
	Warning("This WARNING message SHOULD appear") // This one should log

	output := buf.String()

	if strings.Contains(output, "INFO message should not appear") {
		t.Errorf("INFO message logged when INFO level was not active. Output: %s", output)
	}
	if strings.Contains(output, "DEBUG message should not appear") {
		t.Errorf("DEBUG message logged when DEBUG level was not active. Output: %s", output)
	}
	if strings.Contains(output, "ERROR message should not appear") {
		t.Errorf("ERROR message logged when ERROR level was not active. Output: %s", output)
	}

	// Check that only the warning message is present
	// Extract the message part of the warning log. If no warning, this will fail.
	if strings.Contains(output, "WARNING message SHOULD appear") {
		actualWarningMessage := extractMessage(t, output, warnPrimaryLogPattern_NonDebug)
		expectedWarningMessage := "This WARNING message SHOULD appear\n"
		if actualWarningMessage != expectedWarningMessage {
			t.Errorf("Warning message content mismatch. Expected '%s', got '%s'. Full output: '%s'",
				expectedWarningMessage, actualWarningMessage, output)
		}
		// Count non-empty lines. There should be only one.
		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) != 1 {
			t.Errorf("Expected only one log line (Warning), but got %d lines. Output: %s", len(lines), output)
		}

	} else {
		t.Errorf("WARNING message NOT logged when WARNING level WAS active. Output: %s", output)
	}
}
