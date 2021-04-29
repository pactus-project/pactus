package http

type Config struct {
	Enable  bool   `toml:"Enable" comment:"Enable Http server for client communication."`
	Address string `toml:"Address" comment:"Address to listen for incoming connections for Http.Default port is 8081."`
}

func DefaultConfig() *Config {
	return &Config{
		Enable:  true,
		Address: "[::]:8081",
	}
}

func TestConfig() *Config {
	return &Config{
		Enable:  true,
		Address: "[::]:0",
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
