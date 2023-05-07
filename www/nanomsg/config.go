package nanomsg

type Config struct {
	Enable bool   `toml:"enable"`
	Listen string `toml:"listen"`
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
