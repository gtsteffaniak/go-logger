package logger

import "log"

// friendly config for yaml/json interfaces
type JsonConfig struct {
	Levels    string `json:"levels"`    // separated list of log levels to enable. (eg. "info|warning|error|debug")
	ApiLevels string `json:"apiLevels"` // separated list of log levels to enable for the API. (eg. "info|warning|error")
	Output    string `json:"output"`    // output location. (eg. "stdout" or "path/to/file.log")
	NoColors  bool   `json:"noColors"`  // disable colors in the output
	Json      bool   `json:"json"`      // output in json format, currently not supported
	Utc       bool   `json:"utc"`       // use UTC time in the output instead of local time
}

// go logger log config
type LoggerConfig struct {
	Levels       []LogLevel
	ApiLevels    []LogLevel
	Stdout       bool
	Disabled     bool
	DebugEnabled bool
	DisabledAPI  bool
	Colors       bool
	Utc          bool
	FilePath     string

	// not exposed
	logger *log.Logger
}
