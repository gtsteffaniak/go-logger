package logger

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"
)

type LogLevel int

const (
	DISABLED LogLevel = 0
	ERROR    LogLevel = 1
	FATAL    LogLevel = 1
	WARNING  LogLevel = 2
	INFO     LogLevel = 3
	DEBUG    LogLevel = 4
	API      LogLevel = 10
	// COLORS
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	GRAY   = "\033[2;37m"
)

var (
	loggers []*LoggerConfig
)

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

// Log prints a log message if its level is greater than or equal to the logger's levels
func Log(level string, msg string, prefix, api bool, color string) {
	LEVEL := stringToLevel[level]
	for _, logger := range loggers {
		if api {
			if logger.DisabledAPI || !slices.Contains(logger.ApiLevels, LEVEL) {
				continue
			}
		} else if level != levels.FATAL {
			if logger.Disabled || !slices.Contains(logger.Levels, LEVEL) {
				continue
			}
		}
		writeOut := msg
		var formattedTime string
		if logger.Utc {
			formattedTime = time.Now().UTC().Format("2006/01/02 15:04:05")
		} else {
			formattedTime = time.Now().Local().Format("2006/01/02 15:04:05")
		}
		if logger.Colors && color != "" {
			formattedTime = formattedTime + color
		}
		if prefix || logger.DebugEnabled {
			logger.logger.SetPrefix(fmt.Sprintf("%s [%s] ", formattedTime, level))
		} else {
			logger.logger.SetPrefix(formattedTime + " ")
		}
		if logger.Colors && color != "" {
			writeOut = writeOut + "\033[0m"
		}
		err := logger.logger.Output(3, writeOut) // 3 skips this function for correct file:line
		if err != nil {
			log.Printf("failed to log message '%v' with error `%v`", msg, err)
		}
	}
	if level == levels.FATAL {
		os.Exit(1)
	}
}

// --- Sprintf-style logging functions ---

func Debugf(format string, a ...interface{}) {
	messageToSend := fmt.Sprintf(format, a...)
	if len(loggers) > 0 {
		Log(levels.DEBUG, messageToSend, true, false, GRAY)
	} else {
		log.Println("[DEBUG]", messageToSend)
	}
}

func Infof(format string, a ...interface{}) {
	messageToSend := fmt.Sprintf(format, a...)
	if len(loggers) > 0 {
		Log(levels.INFO, messageToSend, true, false, "")
	} else {
		log.Println("[INFO]", messageToSend)
	}
}

func Warningf(format string, a ...interface{}) {
	messageToSend := fmt.Sprintf(format, a...)
	if len(loggers) > 0 {
		Log(levels.WARNING, messageToSend, true, false, YELLOW)
	} else {
		log.Println("[WARN ]", messageToSend)
	}
}

func Errorf(format string, a ...interface{}) {
	messageToSend := fmt.Sprintf(format, a...)
	if len(loggers) > 0 {
		Log(levels.ERROR, messageToSend, true, false, RED)
	} else {
		log.Println("[ERROR]", messageToSend)
	}
}

func Fatalf(format string, a ...interface{}) {
	messageToSend := fmt.Sprintf(format, a...)
	if len(loggers) > 0 {
		Log(levels.FATAL, messageToSend, true, false, RED)
	} else {
		log.Println("[FATAL]", messageToSend)
		os.Exit(1)
	}
}

func Apif(statusCode int, format string, a ...interface{}) {
	messageToSend := fmt.Sprintf(format, a...)
	var levelStr, colorStr string
	if statusCode > 304 && statusCode < 500 {
		levelStr, colorStr = levels.WARNING, YELLOW
	} else if statusCode >= 500 {
		levelStr, colorStr = levels.ERROR, RED
	} else {
		levelStr, colorStr = levels.INFO, GREEN
	}
	Log(levelStr, messageToSend, false, true, colorStr)
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
	messageToSend := sprintArgs(a...)
	if len(loggers) > 0 {
		Log(levels.DEBUG, messageToSend, true, false, GRAY)
	} else {
		log.Println("[DEBUG]", messageToSend)
	}
}

func Info(a ...interface{}) {
	messageToSend := sprintArgs(a...)
	if len(loggers) > 0 {
		Log(levels.INFO, messageToSend, true, false, "")
	} else {
		log.Println("[INFO]", messageToSend)
	}
}

func Warning(a ...interface{}) {
	messageToSend := sprintArgs(a...)
	if len(loggers) > 0 {
		Log(levels.WARNING, messageToSend, true, false, YELLOW)
	} else {
		log.Println("[WARN ]", messageToSend)
	}
}

func Error(a ...interface{}) {
	messageToSend := sprintArgs(a...)
	if len(loggers) > 0 {
		Log(levels.ERROR, messageToSend, true, false, RED)
	} else {
		log.Println("[ERROR]", messageToSend)
	}
}

func Fatal(a ...interface{}) {
	messageToSend := sprintArgs(a...)
	if len(loggers) > 0 {
		Log(levels.FATAL, messageToSend, true, false, RED)
	} else {
		log.Println("[FATAL]", messageToSend)
		os.Exit(1)
	}
}

func Api(statusCode int, a ...interface{}) {
	messageToSend := sprintArgs(a...)
	var levelStr, colorStr string
	if statusCode > 304 && statusCode < 500 {
		levelStr, colorStr = levels.WARNING, YELLOW
	} else if statusCode >= 500 {
		levelStr, colorStr = levels.ERROR, RED
	} else {
		levelStr, colorStr = levels.INFO, GREEN
	}
	Log(levelStr, messageToSend, false, true, colorStr)
}
