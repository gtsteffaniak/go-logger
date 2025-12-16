package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	DISABLED LogLevel = 0
	ERROR    LogLevel = 1
	FATAL    LogLevel = 2
	WARNING  LogLevel = 3
	INFO     LogLevel = 4
	DEBUG    LogLevel = 5
	API      LogLevel = 10
	// COLORS
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	GRAY   = "\033[2;37m"
)

var (
	globalLogger Logger
)

// SetGlobalLogger sets the global logger instance for legacy compatibility
func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}

type levelConsts struct {
	INFO     string
	FATAL    string
	ERROR    string
	WARNING  string
	DEBUG    string
	API      string
	DISABLED string
}

var levels = levelConsts{
	INFO:     "INFO ", // with consistent space padding
	FATAL:    "FATAL",
	ERROR:    "ERROR",
	WARNING:  "WARN ", // with consistent space padding
	DEBUG:    "DEBUG",
	DISABLED: "DISABLED",
	API:      "API",
}

// stringToLevel maps string representation to LogLevel
var stringToLevel = map[string]LogLevel{
	"DEBUG":    DEBUG,
	"INFO ":    INFO, // with consistent space padding
	"ERROR":    ERROR,
	"DISABLED": DISABLED,
	"WARN ":    WARNING, // with consistent space padding
	"FATAL":    FATAL,
	"API":      API,
}

// logMessage is a helper function that uses the global logger if available
func logMessage(level string, message string, modernLog func(string)) {
	if globalLogger != nil {
		modernLog(message)
	} else {
		log.Printf("[%s] %s", level, message)
	}
}

// --- Sprintf-style logging functions ---

func Debugf(format string, a ...interface{}) {
	messageToSend := fmt.Sprintf(format, a...)
	logMessage(levels.DEBUG, messageToSend, func(msg string) { globalLogger.Debugf(format, a...) })
}

func Infof(format string, a ...interface{}) {
	messageToSend := fmt.Sprintf(format, a...)
	logMessage(levels.INFO, messageToSend, func(msg string) { globalLogger.Infof(format, a...) })
}

func Warningf(format string, a ...interface{}) {
	messageToSend := fmt.Sprintf(format, a...)
	logMessage(levels.WARNING, messageToSend, func(msg string) { globalLogger.Warnf(format, a...) })
}

func Errorf(format string, a ...interface{}) {
	messageToSend := fmt.Sprintf(format, a...)
	logMessage(levels.ERROR, messageToSend, func(msg string) { globalLogger.Errorf(format, a...) })
}

func Fatalf(format string, a ...interface{}) {
	messageToSend := fmt.Sprintf(format, a...)
	if globalLogger != nil {
		globalLogger.Fatalf(format, a...)
	} else {
		log.Println("[FATAL]", messageToSend)
		os.Exit(1)
	}
}

func Apif(statusCode int, format string, a ...interface{}) {
	if globalLogger != nil {
		globalLogger.APIf(statusCode, format, a...)
	} else {
		messageToSend := fmt.Sprintf(format, a...)
		log.Printf("[API] %s\n", messageToSend)
	}
}

// --- Sprint-style logging functions (space-separated arguments) ---

func sprintArgs(a ...interface{}) string {
	if len(a) == 0 {
		return ""
	}
	// fmt.Sprintln always adds a newline, so we trim it.
	return strings.TrimSuffix(fmt.Sprintln(a...), "\n")
}

func Debug(a ...interface{}) {
	logMessage(levels.DEBUG, sprintArgs(a...), func(msg string) { globalLogger.Debugf("%s", msg) })
}

func Info(a ...interface{}) {
	logMessage(levels.INFO, sprintArgs(a...), func(msg string) { globalLogger.Infof("%s", msg) })
}

func Warning(a ...interface{}) {
	logMessage(levels.WARNING, sprintArgs(a...), func(msg string) { globalLogger.Warnf("%s", msg) })
}

func Error(a ...interface{}) {
	logMessage(levels.ERROR, sprintArgs(a...), func(msg string) { globalLogger.Errorf("%s", msg) })
}

func Fatal(a ...interface{}) {
	messageToSend := sprintArgs(a...)
	if globalLogger != nil {
		globalLogger.Fatal(messageToSend)
	} else {
		log.Printf("[FATAL] %s\n", messageToSend)
		os.Exit(1)
	}
}

func Api(statusCode int, a ...interface{}) {
	messageToSend := sprintArgs(a...)
	if globalLogger != nil {
		globalLogger.API(statusCode, messageToSend)
	} else {
		log.Printf("[API] %s\n", messageToSend)
	}
}
