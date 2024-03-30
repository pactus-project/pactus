package jsonrpc

type Config struct {
	Enable bool   `toml:"enable"`
	Listen string `toml:"listen"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable: false,
		Listen: "localhost:8545",
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	return nil
}
