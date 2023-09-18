package nanomsg

type Config struct {
	Enable bool   `toml:"enable"`
	Listen string `toml:"listen"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable: false,
		Listen: "40899",
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	return nil
}
