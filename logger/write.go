package logger

import (
	"fmt"
	"log"
	"os"
	"slices"
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
	loggers []*Logger
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
// It now accepts args ...interface{} for Sprintf-like behavior
func Log(level string, format string, prefix, api bool, color string, args ...interface{}) {
	LEVEL := stringToLevel[level]
	for _, logger := range loggers {
		if api {
			if logger.disabledAPI || !slices.Contains(logger.apiLevels, LEVEL) {
				continue
			}
		} else {
			if logger.disabled || !slices.Contains(logger.levels, LEVEL) {
				continue
			}
		}

		var writeOut string
		if len(args) > 0 {
			writeOut = fmt.Sprintf(format, args...)
		} else {
			writeOut = format
		}

		formattedTime := time.Now().Format("2006/01/02 15:04:05")
		if logger.colors && color != "" {
			formattedTime = formattedTime + color
		}
		if prefix || logger.debugEnabled {
			logger.logger.SetPrefix(fmt.Sprintf("%s [%s] ", formattedTime, level))
		} else {
			logger.logger.SetPrefix(formattedTime + " ")
		}
		if logger.colors && color != "" {
			writeOut = writeOut + "\033[0m"
		}
		err := logger.logger.Output(3, writeOut) // 3 skips this function for correct file:line
		if err != nil {
			log.Printf("failed to log message '%v' with error `%v`", format, err)
		}
		if level == levels.FATAL {
			os.Exit(1)
		}
	}
}

func Api(msg string, statusCode int) {
	// redirects are not warnings anymore
	// content not modified is not a warning anymore
	if statusCode > 304 && statusCode < 500 {
		Log(levels.WARNING, msg, false, true, YELLOW)
	} else if statusCode >= 500 {
		Log(levels.ERROR, msg, false, true, RED)
	} else {
		Log(levels.INFO, msg, false, true, GREEN)
	}
}

// Helper methods for specific log levels
func Debug(format string, args ...interface{}) {
	if len(loggers) > 0 {
		Log(levels.DEBUG, format, true, false, GRAY, args...)
	} else {
		if len(args) > 0 {
			log.Println("[DEBUG]", fmt.Sprintf(format, args...))
		} else {
			log.Println("[DEBUG]", format)
		}
	}
}

func Info(format string, args ...interface{}) {
	if len(loggers) > 0 {
		Log(levels.INFO, format, false, false, "", args...)
	} else {
		if len(args) > 0 {
			log.Println("[DEBUG]", fmt.Sprintf(format, args...))
		} else {
			log.Println("[DEBUG]", format)
		}
	}
}

func Warning(format string, args ...interface{}) {
	if len(loggers) > 0 {
		Log(levels.WARNING, format, true, false, YELLOW, args...)
	} else {
		if len(args) > 0 {
			log.Println("[WARN ]", fmt.Sprintf(format, args...))
		} else {
			log.Println("[WARN ]", format)
		}
	}
}

func Error(format string, args ...interface{}) {
	if len(loggers) > 0 {
		Log(levels.ERROR, format, true, false, RED, args...)
	} else {
		if len(args) > 0 {
			log.Println("[ERROR]", fmt.Sprintf(format, args...))
		} else {
			log.Println("[ERROR]", format)
		}
	}
}

func Fatal(format string, args ...interface{}) {
	if len(loggers) > 0 {
		Log(levels.FATAL, format, true, false, RED, args...)
	} else {
		if len(args) > 0 {
			log.Fatal("[FATAL]", fmt.Sprintf(format, args...))
		} else {
			log.Fatal("[FATAL]", format)
		}
	}
}
