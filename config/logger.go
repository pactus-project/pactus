package config

type LoggerConfig struct {
	Levels map[string]string
}

func DefaultLoggerConfig() *LoggerConfig {
	def := &LoggerConfig{
		Levels: make(map[string]string),
	}

	def.Levels["default"] = "trace"
	def.Levels["network"] = "trace"
	def.Levels["consensus"] = "trace"
	def.Levels["state"] = "trace"
	return def
}
