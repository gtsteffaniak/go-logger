package logger

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	// Logger uses time, slices, os
)

// setupForPrimaryLoggerTest (no changes needed)
func setupForPrimaryLoggerTest(t *testing.T, buf *bytes.Buffer, activeLevels []LogLevel, noColors bool) {
	t.Helper()
	savedLoggers := append([]*LoggerConfig(nil), loggers...)
	savedStdOutLoggerExists := stdOutLoggerExists
	t.Cleanup(func() {
		loggers = savedLoggers
		stdOutLoggerExists = savedStdOutLoggerExists
	})
	loggers = nil
	stdOutLoggerExists = false
	config := LoggerConfig{
		Stdout:    true,
		Levels:    activeLevels,
		ApiLevels: []LogLevel{},
		Colors:    !noColors,
	}
	// Assuming NewLogger is defined in your logger package (e.g. setup.go)
	// and Logger struct has fields like 'logger', 'apiLevels', 'disabledAPI', etc.
	l, err := AddLogger(config)
	if err != nil {
		t.Fatalf("Failed to create test logger: %v", err)
	}
	l.logger.SetOutput(buf)
	loggers = []*LoggerConfig{l}
}

// setupForFallbackTest (no changes needed)
func setupForFallbackTest(t *testing.T, buf *bytes.Buffer) {
	t.Helper()
	savedLoggers := append([]*LoggerConfig(nil), loggers...)
	savedStdOutLoggerExists := stdOutLoggerExists
	currentGlobalLogOutput := log.Writer() // May need adjustment if global log output isn't os.Stderr
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

// Regex patterns for prefix stripping (no changes needed)
var infoPrimaryLogPattern_NonDebug = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s(.*\n)$`)
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

func TestInfo_PrimaryLogger(t *testing.T) {
	var buf bytes.Buffer
	noColors := true

	t.Run("InfoNonDebugMode", func(t *testing.T) {
		setupForPrimaryLoggerTest(t, &buf, []LogLevel{INFO}, noColors)
		buf.Reset()

		// USE Infof for formatting
		Infof("Hello %s, number %d", "Alice", 100)
		expectedMessage := "Hello Alice, number 100\n"
		actualOutput := buf.String()
		// The Log function for Infof (and Info) now has prefix=true.
		// If logger.debugEnabled is false, prefix for INFO is "timestamp [INFO ] "
		// This regex needs to match that if infoPrimaryLogPattern_NonDebug is "timestamp "
		// Let's assume infoPrimaryLogPattern_NonDebug is for when prefix=false in Log OR (prefix=true AND level is not one that gets [LEVEL] like INFO)
		// Your Log function: if prefix || logger.debugEnabled { fmt.Sprintf("%s [%s] ", formattedTime, level) } else { formattedTime + " " }
		// For Infof, prefix is true. So it will always be "YYYY/MM/DD HH:MM:SS [INFO ] "
		// Let's define a generic pattern for this or adjust.
		infoPrefixPattern := regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\s\[INFO\s+\]\s(.*\n)$`)
		actualMessage := extractMessage(t, actualOutput, infoPrefixPattern)

		if actualMessage != expectedMessage {
			t.Errorf("Expected log message '%s', got '%s'. Full output: '%s'", expectedMessage, actualMessage, actualOutput)
		}
	})

	t.Run("InfoWithLoggerDebugMode", func(t *testing.T) {
		setupForPrimaryLoggerTest(t, &buf, []LogLevel{INFO, DEBUG}, noColors)
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
		setupForPrimaryLoggerTest(t, &buf, []LogLevel{INFO}, noColors)
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

func TestInfo_FallbackLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForFallbackTest(t, &buf)

	// USE Infof for formatting
	Infof("Message for %s", "fallback_user")
	// Fallback for Infof is log.Println("[INFO]", messageToSend)
	expectedOutput := "[INFO] Message for fallback_user\n"
	actualOutput := buf.String()

	if actualOutput != expectedOutput {
		t.Errorf("Expected fallback log output '%s', got '%s'", expectedOutput, actualOutput)
	}

	buf.Reset()
	// USE Infof for single string formatted
	Infof("Simple info fallback")
	expectedSimple := "[INFO] Simple info fallback\n"
	actualSimple := buf.String()
	if actualSimple != expectedSimple {
		t.Errorf("Expected simple fallback log output '%s', got '%s'", expectedSimple, actualSimple)
	}

	buf.Reset()
	// Test Info (Sprint) fallback
	Info("Simple", "sprint", "fallback")
	expectedSprintFallback := "[INFO] Simple sprint fallback\n"
	actualSprintFallback := buf.String()
	if actualSprintFallback != expectedSprintFallback {
		t.Errorf("Expected Sprint fallback log output '%s', got '%s'", expectedSprintFallback, actualSprintFallback)
	}
}

func TestDebug_PrimaryLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForPrimaryLoggerTest(t, &buf, []LogLevel{DEBUG}, true)

	// USE Debugf for formatting
	Debugf("Processing %s, item %d", "data_set", 77)
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

	// USE Debugf for formatting
	Debugf("Fallback debug %s", "message")
	expectedOutput := "[DEBUG] Fallback debug message\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected fallback log output '%s', got '%s'", expectedOutput, actualOutput)
	}
}

func TestWarning_PrimaryLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForPrimaryLoggerTest(t, &buf, []LogLevel{WARNING, INFO}, true)

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

func TestWarning_FallbackLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForFallbackTest(t, &buf)

	// USE Warningf for formatting
	Warningf("Fallback warning: %s", "check this")
	expectedOutput := "[WARN ] Fallback warning: check this\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected fallback log output '%s', got '%s'", expectedOutput, actualOutput)
	}
}

func TestError_PrimaryLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForPrimaryLoggerTest(t, &buf, []LogLevel{ERROR, INFO}, true)

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

func TestError_FallbackLogger(t *testing.T) {
	var buf bytes.Buffer
	setupForFallbackTest(t, &buf)

	// USE Errorf for formatting
	Errorf("Fallback error: %s", "system down")
	expectedOutput := "[ERROR] Fallback error: system down\n"
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected fallback log output '%s', got '%s'", expectedOutput, actualOutput)
	}
}

func TestFormatting_VariousArgTypes_Primary(t *testing.T) {
	var buf bytes.Buffer
	setupForPrimaryLoggerTest(t, &buf, []LogLevel{INFO}, true)

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
	setupForPrimaryLoggerTest(t, &buf, []LogLevel{WARNING}, true) // Only WARNING is active

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

// You would also add tests for Api and Apif functions similarly.
// For example:
func TestApi_PrimaryLogger(t *testing.T) {
	var buf bytes.Buffer
	// Assuming API level is distinct and needs to be activated.
	// And that API logs use their own levels (INFO, WARNING, ERROR based on status code)
	// The setup should activate levels that Apif will use, e.g., INFO, WARNING, ERROR for the 'apiLevels'
	// For simplicity, using general levels here. Adjust `setupForPrimaryLoggerTest` or provide
	// a specific setup if `apiLevels` needs separate handling.
	// The current setupForPrimaryLoggerTest sets 'levels', not 'apiLevels'.
	// You might need to adjust NewLogger or setup to set 'apiLevels' appropriately.
	// For now, let's assume apiLevels in NewLogger gets populated with general levels for testing.

	// This test needs logger.apiLevels to be set appropriately by setupForPrimaryLoggerTest/NewLogger
	// For example, NewLogger("", []LogLevel{}, []LogLevel{INFO, WARNING, ERROR}, noColors)
	// The current setupForPrimaryLoggerTest passes empty []LogLevel{} for apiLevels.
	// This means Apif calls might not log anything unless NewLogger has a default for apiLevels.
	// Let's assume NewLogger makes apiLevels default to something reasonable if empty, or test will fail.
	// For now, I'll write the test assuming apiLevels will allow these messages.

	noColors := true
	// Create a logger instance specifically for API tests if apiLevels are distinct
	loggers = nil // Clear global loggers
	stdOutLoggerExists = false

	// Ensure logger/setup.go's NewLogger can handle distinct apiLevels
	l, err := AddLogger(LoggerConfig{
		Stdout:    true,
		Colors:    !noColors,
		ApiLevels: []LogLevel{INFO, WARNING, ERROR},
		Levels:    []LogLevel{INFO, WARNING, ERROR},
		Disabled:  false,
	}) // Activate relevant levels for API
	if err != nil {
		t.Fatalf("Failed to create test logger for API: %v", err)
	}
	l.logger.SetOutput(&buf)
	loggers = []*LoggerConfig{l}
	t.Cleanup(func() { loggers = nil })

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
