package logger

import "path/filepath"

const (
	LogDirectory = "logs"
	LogFilename  = "pactus.log"
)

type Config struct {
	// ConsoleLoggingEnabled Enable console logging
	ConsoleLoggingEnabled bool `toml:"console_logging_enabled"`
	// EncodeLogsAsJSON makes the log framework log JSON
	EncodeLogsAsJSON bool `toml:"encode_logs_as_json"`
	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool `toml:"file_logging_enabled"`
	// Directory to log to when file logging is enabled
	Directory string `toml:"directory"`
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string `toml:"filename"`
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `toml:"max_size"`
	// MaxBackups the max number of rolled files to keep
	MaxBackups int `toml:"max_backups"`
	// MaxAge the max age in days to keep a logfile
	MaxAge int `toml:"max_age"`

	Colorful bool              `toml:"colorful"`
	Levels   map[string]string `toml:"levels"`
}

func DefaultConfig(workingDir ...string) *Config {
	var logDir string
	var logFilename string

	if len(workingDir) > 0 {
		logDir = filepath.Join(workingDir[0], LogDirectory)
		logFilename = filepath.Join(logDir, LogFilename)
	}

	conf := &Config{
		ConsoleLoggingEnabled: true,
		EncodeLogsAsJSON:      true,
		Directory:             logDir,
		MaxSize:               500, // megabytes
		MaxBackups:            3,   // files
		MaxAge:                30,  // days
		Filename:              logFilename,
		Levels:                make(map[string]string),
		Colorful:              true,
	}

	conf.Levels["default"] = "info"
	conf.Levels["_network"] = "info"
	conf.Levels["_consensus"] = "info"
	conf.Levels["_state"] = "info"
	conf.Levels["_sync"] = "warning"
	conf.Levels["_pool"] = "error"
	conf.Levels["_http"] = "error"
	conf.Levels["_grpc"] = "error"

	return conf
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	return nil
}
