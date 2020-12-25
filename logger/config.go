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
	conf.Levels["_consensus"] = "error"
	conf.Levels["_state"] = "info"
	conf.Levels["_sync"] = "error"
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
	conf.Levels["_network"] = "trace"
	conf.Levels["_consensus"] = "trace"
	conf.Levels["_state"] = "trace"
	conf.Levels["_sync"] = "trace"
	conf.Levels["_pool"] = "trace"
	conf.Levels["_capnp"] = "trace"
	conf.Levels["_http"] = "trace"
	conf.Colorfull = true

	return conf
}

// SanityCheck is a basic hecks for config
func (conf *Config) SanityCheck() error {
	return nil
}
