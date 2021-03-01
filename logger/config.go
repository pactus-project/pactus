package logger

type Config struct {
	Levels    map[string]string
	Colorfull bool
}

func DefaultConfig() *Config {
	conf := &Config{
		Levels: make(map[string]string),
	}

	conf.Levels["default"] = "info"
	conf.Levels["_network"] = "error"
	conf.Levels["_consensus"] = "info"
	conf.Levels["_state"] = "info"
	conf.Levels["_sync"] = "warning"
	conf.Levels["_pool"] = "error"
	conf.Levels["_capnp"] = "error"
	conf.Levels["_http"] = "error"
	conf.Colorfull = true

	return conf
}

func TestConfig() *Config {
	conf := &Config{
		Levels: make(map[string]string),
	}

	conf.Levels["default"] = "debug"
	conf.Levels["_network"] = "debug"
	conf.Levels["_consensus"] = "debug"
	conf.Levels["_state"] = "debug"
	conf.Levels["_sync"] = "debug"
	conf.Levels["_pool"] = "debug"
	conf.Levels["_capnp"] = "debug"
	conf.Levels["_http"] = "debug"
	conf.Colorfull = true

	return conf
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
