package logger

type Config struct {
	Levels    map[string]string `toml:"" comment:"Levels contains trace,debug,info,warning,error type."`
	Colorfull bool              `toml:"" comment:"Colorfull Output format can be enable or disable. Default is true."`
}

func DefaultConfig() *Config {
	conf := &Config{
		Levels:    make(map[string]string),
		Colorfull: true,
	}

	conf.Levels["default"] = "info"
	conf.Levels["_network"] = "error"
	conf.Levels["_consensus"] = "info"
	conf.Levels["_state"] = "info"
	conf.Levels["_sync"] = "warning"
	conf.Levels["_pool"] = "error"
	conf.Levels["_capnp"] = "error"
	conf.Levels["_http"] = "error"
	conf.Levels["_grpc"] = "error"

	return conf
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
