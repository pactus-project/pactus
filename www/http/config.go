package http

type Config struct {
	Enable  bool
	Address string
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
