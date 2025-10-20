package jsonrpc

// Config defines parameters for the JSON-RPC server.
type Config struct {
	Enable  bool     `toml:"enable"`
	Listen  string   `toml:"listen"`
	Origins []string `toml:"origins"`
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
