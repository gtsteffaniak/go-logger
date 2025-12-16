package logger

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"
	// Logger uses time, slices, os
)

// setupForModernLoggerTest sets up a modern logger for testing global functions
func setupForModernLoggerTest(t *testing.T, buf *bytes.Buffer, levels string, noColors bool) Logger {
	t.Helper()
	savedGlobalLogger := globalLogger
	t.Cleanup(func() {
		globalLogger = savedGlobalLogger
	})

	config := JsonConfig{
		Output:   "STDOUT",
		Levels:   levels,
		NoColors: noColors,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create test logger: %v", err)
	}

	// Set the logger as global for testing package-level functions
	SetGlobalLogger(logger)

	// Redirect output to buffer (this is a bit hacky but works for testing)
	// We need to access the underlying logger's output - this will work for non-JSON loggers
	ml := logger.(*modernLogger)
	if len(ml.configs) > 0 {
		ml.configs[0].logger.SetOutput(buf)
	}

	return logger
}

// Regex patterns for prefix stripping (no changes needed)
var infoPrimaryLogPattern_DebugMode = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s\[INFO\s+\]\s[a-zA-Z0-9._/-]+:\d+:\s(.*\n)$`)
var debugPrimaryLogPattern = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s\[DEBUG\]\s[a-zA-Z0-9._/-]+:\d+:\s(.*\n)$`)
var warnPrimaryLogPattern_NonDebug = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s\[WARN\s+\]\s(.*\n)$`)
var errorPrimaryLogPattern_NonDebug = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s\[ERROR\]\s(.*\n)$`)

// extractMessage helper (no changes needed)
func extractMessage(t *testing.T, output string, pattern *regexp.Regexp) string {
	t.Helper()
	matches := pattern.FindStringSubmatch(output)
	if len(matches) < 2 {
		t.Fatalf("Log output did not match expected pattern '%s'. Full output: '%s'", pattern.String(), output)
		return ""
	}
	return matches[1]
}

func TestInfo_ModernLogger(t *testing.T) {
	var buf bytes.Buffer
	noColors := true

	t.Run("InfoNonDebugMode", func(t *testing.T) {
		setupForModernLoggerTest(t, &buf, "INFO", noColors)
		buf.Reset()

		// USE Infof for formatting
		Infof("Hello %s, number %d", "Alice", 100)
		expectedMessage := "Hello Alice, number 100\n"
		actualOutput := buf.String()
		infoPrefixPattern := regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s\[INFO\s+\]\s(.*\n)$`)
		actualMessage := extractMessage(t, actualOutput, infoPrefixPattern)

		if actualMessage != expectedMessage {
			t.Errorf("Expected log message '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
		}
	})

	t.Run("InfoWithLoggerDebugMode", func(t *testing.T) {
		setupForModernLoggerTest(t, &buf, "INFO,DEBUG", noColors)
		buf.Reset()

		// USE Infof for formatting
		Infof("Hello %s, debug number %d", "Bob", 200)
		expectedMessage := "Hello Bob, debug number 200\n"
		actualOutput := buf.String()
		actualMessage := extractMessage(t, actualOutput, infoPrimaryLogPattern_DebugMode) // This pattern expects file:line

		if actualMessage != expectedMessage {
			t.Errorf("Expected log message (debug mode) '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
		}
	})

	t.Run("InfoSprint", func(t *testing.T) {
		setupForModernLoggerTest(t, &buf, "INFO", noColors)
		buf.Reset()
		Info("Hello", "Alice", 100)            // Uses Info (Sprint)
		expectedMessage := "Hello Alice 100\n" // fmt.Sprint behavior
		actualOutput := buf.String()
		infoPrefixPattern := regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s\[INFO\s+\]\s(.*\n)$`)
		actualMessage := extractMessage(t, actualOutput, infoPrefixPattern)
		if actualMessage != expectedMessage {
			t.Errorf("Expected Sprint log message '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
		}
	})
}

func TestDebug_ModernLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForModernLoggerTest(t, &buf, "DEBUG", true)

	// USE Debugf for formatting
	Debugf("Processing %s, item %d", "data_set", 77)
	expectedMessage := "Processing data_set, item 77\n"
	actualOutput := buf.String()
	actualMessage := extractMessage(t, actualOutput, debugPrimaryLogPattern)

	if actualMessage != expectedMessage {
		t.Errorf("Expected log message '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
	}
}

func TestWarning_ModernLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForModernLoggerTest(t, &buf, "WARNING,INFO", true)

	// USE Warningf for formatting
	Warningf("Potential issue with %s", "config_value")
	expectedMessage := "Potential issue with config_value\n"
	actualOutput := buf.String()
	// Warningf uses prefix=true. If not in debug mode, it's "timestamp [WARN ] "
	actualMessage := extractMessage(t, actualOutput, warnPrimaryLogPattern_NonDebug)

	if actualMessage != expectedMessage {
		t.Errorf("Expected log message '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
	}
}

func TestError_ModernLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForModernLoggerTest(t, &buf, "ERROR,INFO", true)

	errVal := fmt.Errorf("critical failure")
	// USE Errorf for formatting
	Errorf("System error: %v", errVal)
	expectedMessage := fmt.Sprintf("System error: %v\n", errVal)
	actualOutput := buf.String()
	// Errorf uses prefix=true. If not in debug mode, it's "timestamp [ERROR] "
	actualMessage := extractMessage(t, actualOutput, errorPrimaryLogPattern_NonDebug)

	if actualMessage != expectedMessage {
		t.Errorf("Expected log message '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
	}
}

func TestFormatting_VariousArgTypes_Modern(t *testing.T) {
	var buf bytes.Buffer
	setupForModernLoggerTest(t, &buf, "INFO", true)

	type MyStruct struct{ Name string }
	s := MyStruct{Name: "DataObject"}
	ptr := &s

	// USE Infof for formatting
	Infof("String: '%s', Int: %d, Bool: %t, Float: %.2f, Struct: %v, StructPtr: %v, PtrAddr: %p",
		"test string", 987, false, 123.4567, s, ptr, ptr)

	output := buf.String()
	// Infof uses prefix=true. If not in debug mode, "timestamp [INFO ]"
	infoPrefixPattern := regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s\[INFO\s+\]\s(.*\n)$`)
	actualMessage := extractMessage(t, output, infoPrefixPattern)

	expectedMessagePatternStr := fmt.Sprintf("^String: 'test string', Int: 987, Bool: false, Float: 123.46, Struct: %v, StructPtr: %v, PtrAddr: 0x[0-9a-f]+\n$", s, ptr)
	expectedMessagePattern := regexp.MustCompile(expectedMessagePatternStr)

	if !expectedMessagePattern.MatchString(actualMessage) {
		t.Errorf("Formatted message for various types did not match expected pattern.\nExpected pattern: %s\nActual message: %sFull output: %s",
			expectedMessagePattern.String(), actualMessage, output)
	}
}

func TestLogging_LevelNotActive(t *testing.T) {
	var buf bytes.Buffer
	setupForModernLoggerTest(t, &buf, "WARNING", true) // Only WARNING is active

	// Use respective non-f or f functions
	Infof("This INFO message should not appear")   // or Info() if that's the test
	Debugf("This DEBUG message should not appear") // or Debug()
	Errorf("This ERROR message should not appear") // or Error()
	Warningf("This WARNING message SHOULD appear") // This one should log

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

	if strings.Contains(output, "WARNING message SHOULD appear") {
		// Warningf uses prefix=true. If not in debug mode, "timestamp [WARN ] "
		actualWarningMessage := extractMessage(t, output, warnPrimaryLogPattern_NonDebug)
		expectedWarningMessage := "This WARNING message SHOULD appear\n"
		if actualWarningMessage != expectedWarningMessage {
			t.Errorf("Warning message content mismatch. Expected '%s', got '%s'. Full output: '%s'",
				expectedWarningMessage, actualWarningMessage, output)
		}
		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) != 1 {
			t.Errorf("Expected only one log line (Warning), but got %d lines. Output: %s", len(lines), output)
		}
	} else {
		t.Errorf("WARNING message NOT logged when WARNING level WAS active. Output: %s", output)
	}
}

func TestApi_ModernLogger(t *testing.T) {
	var buf bytes.Buffer
	noColors := true

	// Set up modern logger with API levels
	config := JsonConfig{
		Output:    "STDOUT",
		Levels:    "INFO,WARNING,ERROR",
		ApiLevels: "INFO,WARNING,ERROR",
		NoColors:  noColors,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create test logger for API: %v", err)
	}
	SetGlobalLogger(logger)
	t.Cleanup(func() { SetGlobalLogger(nil) })

	// Redirect output to buffer
	ml := logger.(*modernLogger)
	if len(ml.configs) > 0 {
		ml.configs[0].logger.SetOutput(&buf)
	}

	t.Run("ApifInfo", func(t *testing.T) {
		buf.Reset()
		Apif(200, "API call successful: %s", "GET /health")
		expectedMessage := "API call successful: GET /health\n"
		actualOutput := buf.String()
		// Apif logs have prefix=false, so just "timestamp message"
		// and color is GREEN for info (200 status)
		// The regex should be: timestamp (message)
		apiInfoPattern := regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s(.*\n)$`)
		actualMessage := extractMessage(t, actualOutput, apiInfoPattern)
		if actualMessage != expectedMessage {
			t.Errorf("Expected API INFO message '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
		}
	})

	t.Run("ApiSprint", func(t *testing.T) {
		buf.Reset()
		Api(404, "Resource not found:", "/users/123")
		expectedMessage := "Resource not found: /users/123\n" // fmt.Sprint behavior
		actualOutput := buf.String()
		apiWarningPattern := regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s(.*\n)$`)
		actualMessage := extractMessage(t, actualOutput, apiWarningPattern)
		if actualMessage != expectedMessage {
			t.Errorf("Expected API Sprint (Warning) message '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
		}
	})
}
