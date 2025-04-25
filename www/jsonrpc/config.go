package jsonrpc

type Config struct {
	Enable  bool     `toml:"enable"`
	Listen  string   `toml:"listen"`
	Origins []string `toml:"origins"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable: false,
		Listen: "localhost:8545",
	}
}

// BasicCheck performs basic checks on the configuration.
func (*Config) BasicCheck() error {
	return nil
}
