package logger

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

// AddLogger creates a new logger configuration for internal use by modernLogger
func AddLogger(logger LoggerConfig) (*LoggerConfig, error) {
	var flags int
	if slices.Contains(logger.Levels, DEBUG) {
		flags |= log.Lshortfile
	}

	if logger.Stdout {
		logger.logger = log.New(os.Stdout, "", flags)
	} else {
		file, err := os.OpenFile(logger.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
		logger.logger = log.New(file, "", flags)
	}
	return &logger, nil
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
