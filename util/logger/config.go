package logger

type Config struct {
	Colorful           bool              `toml:"colorful"`
	MaxBackups         int               `toml:"max_backups"`
	RotateLogAfterDays int               `toml:"rotate_log_after_days"`
	Compress           bool              `toml:"compress"`
	Levels             map[string]string `toml:"levels"`
}

func DefaultConfig() *Config {
	conf := &Config{
		Levels:             make(map[string]string),
		Colorful:           true,
		MaxBackups:         0,
		RotateLogAfterDays: 1,
		Compress:           true,
	}

	conf.Levels["default"] = "info"
	conf.Levels["_network"] = "error"
	conf.Levels["_consensus"] = "warn"
	conf.Levels["_state"] = "info"
	conf.Levels["_sync"] = "error"
	conf.Levels["_pool"] = "error"
	conf.Levels["_http"] = "error"
	conf.Levels["_grpc"] = "error"
	conf.Levels["_nonomsg"] = "error"

	return conf
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	return nil
}
