package logger

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

var stdOutLoggerExists bool

// NewLogger creates a new Logger instance with separate file and stdout loggers
func AddLogger(logger LoggerConfig) (*LoggerConfig, error) {
	var flags int
	if slices.Contains(logger.Levels, DEBUG) {
		flags |= log.Lshortfile
	}

	if slices.Contains(logger.Levels, DEBUG) {
		flags |= log.Lshortfile
	}

	if logger.Stdout {
		logger.logger = log.New(os.Stdout, "", flags)
		stdOutLoggerExists = true
	} else {
		file, err := os.OpenFile(logger.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
		logger.logger = log.New(file, "", flags)
	}
	return &logger, nil
}

// SetupLogger configures the logger with file and stdout options and their respective log levels
func SetupLogger(config JsonConfig) error {
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
		// Convert level strings to LogLevel
		level, ok := stringToLevel[upperLevel]
		if !ok {
			loggers = []*LoggerConfig{}
			return fmt.Errorf("invalid file log level: %s", upperLevel)
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
		// Convert level strings to LogLevel
		level, ok := stringToLevel[strings.ToUpper(upperLevel)]
		if !ok {
			return fmt.Errorf("invalid api log level: %s", upperLevel)
		}
		upperApiLevels = append(upperApiLevels, level)
	}
	if len(upperApiLevels) == 0 {
		upperApiLevels = []LogLevel{INFO, ERROR, WARNING}
	}
	if slices.Contains(upperLevels, DISABLED) && slices.Contains(upperApiLevels, DISABLED) {
		// both disabled, not creating a logger
		loggers = []*LoggerConfig{}
		return nil
	}
	outputStdout := strings.ToUpper(config.Output)
	if outputStdout == "STDOUT" {
		config.Output = ""
	}
	if config.Output == "" && stdOutLoggerExists {
		// stdout logger already exists... don't create another
		return fmt.Errorf("stdout logger already exists, could not set config levels=[%v] apiLevels=[%v] noColors=[%v]", levels, config.ApiLevels, config.NoColors)
	}
	// Create the logger
	logger, err := AddLogger(LoggerConfig{
		Levels:       upperLevels,
		ApiLevels:    upperApiLevels,
		Stdout:       config.Output == "",
		Colors:       !config.NoColors,
		Disabled:     slices.Contains(upperLevels, DISABLED),
		DebugEnabled: slices.Contains(upperLevels, DEBUG),
		DisabledAPI:  slices.Contains(upperApiLevels, DISABLED),
		Utc:          config.Utc,
		FilePath:     config.Output,
	})
	if err != nil {
		return err
	}
	loggers = append(loggers, logger)
	return nil
}

func SplitByMultiple(str string) []string {
	delimiters := []rune{'|', ',', ' '}
	return strings.FieldsFunc(str, func(r rune) bool {
		for _, d := range delimiters {
			if r == d {
				return true
			}
		}
		return false
	})
}
