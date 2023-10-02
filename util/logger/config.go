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
	conf.Levels["_network"] = "info"
	conf.Levels["_consensus"] = "info"
	conf.Levels["_state"] = "info"
	conf.Levels["_sync"] = "warning"
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
