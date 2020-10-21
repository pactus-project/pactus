package logger

type Config struct {
	Levels map[string]string
}

func DefaultConfig() *Config {
	def := &Config{
		Levels: make(map[string]string),
	}

	def.Levels["default"] = "trace"
	def.Levels["network"] = "trace"
	def.Levels["consensus"] = "trace"
	def.Levels["state"] = "trace"
	def.Levels["sync"] = "trace"
	return def
}
