package txpool

type Config struct {
	MaxSize int `toml:"max_size"`
}

func DefaultConfig() *Config {
	return &Config{
		MaxSize: 2000,
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if conf.MaxSize == 0 {
		return ConfigError{
			Reason: "maxSize can't be negative or zero",
		}
	}

	return nil
}

func (conf *Config) sortitionPoolSize() int {
	return int(float32(conf.MaxSize) * 0.05)
}

func (conf *Config) bondPoolSize() int {
	return int(float32(conf.MaxSize) * 0.05)
}

func (conf *Config) unbondPoolSize() int {
	return int(float32(conf.MaxSize) * 0.05)
}

func (conf *Config) withdrawPoolSize() int {
	return int(float32(conf.MaxSize) * 0.05)
}

func (conf *Config) sendPoolSize() int {
	return int(float32(conf.MaxSize) * 0.8)
}
