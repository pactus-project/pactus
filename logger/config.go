package logger

type Config struct {
	Levels map[string]string
}

func DefaultConfig() *Config {
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
