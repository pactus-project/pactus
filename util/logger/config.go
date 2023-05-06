package logger

type Config struct {
	Colorful bool              `toml:"colorful"`
	Levels   map[string]string `toml:"levels"`
}

func DefaultConfig() *Config {
	conf := &Config{
		Levels:   make(map[string]string),
		Colorful: true,
	}

	conf.Levels["default"] = "info"
	conf.Levels["_network"] = "error"
	conf.Levels["_consensus"] = "info"
	conf.Levels["_state"] = "info"
	conf.Levels["_sync"] = "warning"
	conf.Levels["_pool"] = "error"
	conf.Levels["_http"] = "error"
	conf.Levels["_grpc"] = "error"

	return conf
}

// SanityCheck performs basic checks on the configuration.
func (conf *Config) SanityCheck() error {
	return nil
}
