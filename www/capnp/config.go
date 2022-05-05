package capnp

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

// SanityCheck is a basic checks for config.
func (conf *Config) SanityCheck() error {
	return nil
}
