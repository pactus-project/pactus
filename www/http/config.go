package http

type Config struct {
	Enable  bool   `toml:"" comment:"Enable HTTP server for client communication."`
	Address string `toml:"" comment:"Address of HTTP server."`
}

func DefaultConfig() *Config {
	return &Config{
		Enable:  false,
		Address: "[::]:8081",
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
