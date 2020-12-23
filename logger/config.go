package logger

type Config struct {
	Levels map[string]string
}

func DefaultConfig() *Config {
	def := &Config{
		Levels: make(map[string]string),
	}

	def.Levels["default"] = "info"
	def.Levels["_network"] = "error"
	def.Levels["_consensus"] = "error"
	def.Levels["_state"] = "info"
	def.Levels["_sync"] = "error"
	def.Levels["_pool"] = "error"
	def.Levels["_capnp"] = "error"
	def.Levels["_http"] = "error"
	return def
}

func TestConfig() *Config {
	def := &Config{
		Levels: make(map[string]string),
	}

	def.Levels["default"] = "trace"
	def.Levels["_network"] = "trace"
	def.Levels["_consensus"] = "trace"
	def.Levels["_state"] = "trace"
	def.Levels["_sync"] = "trace"
	def.Levels["_pool"] = "trace"
	def.Levels["_capnp"] = "trace"
	def.Levels["_http"] = "trace"
	return def
}

// SanityCheck is a basic hecks for config
func (conf *Config) SanityCheck() error {
	return nil
}
