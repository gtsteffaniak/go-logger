package main

import (
	"github.com/gtsteffaniak/go-logger/logger"
)

func main() {
	// example stdout logger
	config := logger.JsonConfig{
		Levels:    "INFO,DEBUG",
		ApiLevels: "INFO,ERROR",
		NoColors:  false,
	}
	err := logger.SetupLogger(config)
	if err != nil {
		logger.Errorf("failed to setup logger: %v", err)
	}
	config.Output = "./stdout.log"
	config.Utc = true
	config.NoColors = true
	err = logger.SetupLogger(config)
	if err != nil {
		logger.Errorf("failed to setup file logger: %v", err)
	}
	logger.Debugf("this is a debug format int value %d in message.", 400)
	logger.Info("this is a basic info message from the logger.")
	logger.Api(200, "api call successful")
	logger.Api(400, "api call warning")
	logger.Api(500, "api call error")
	logger.Fatal("this is a fatal message, the program will exit 1")
}
