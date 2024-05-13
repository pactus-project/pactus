package http

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

// BasicCheck performs basic checks on the configuration.
func (*Config) BasicCheck() error {
	return nil
}
