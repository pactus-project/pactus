package http

type Config struct {
	Enable         bool   `toml:"enable"`
	Listen         string `toml:"listen"`
	EnableDebugger bool   `toml:"-"` // EnableDebugger is private configs
}

func DefaultConfig() *Config {
	return &Config{
		Enable:         false,
		Listen:         "",
		EnableDebugger: false,
	}
}

// BasicCheck performs basic checks on the configuration.
func (*Config) BasicCheck() error {
	return nil
}
