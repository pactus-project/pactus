package nanomsg

type Config struct {
	Listen string `toml:"listen"`
	Enable bool   `toml:"enable"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable: false,
		Listen: "",
	}
}

// SanityCheck performs basic checks on the configuration.
func (conf *Config) SanityCheck() error {
	return nil
}
