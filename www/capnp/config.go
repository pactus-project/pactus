package capnp

type Config struct {
	Enable  bool
	Address string
}

func DefaultConfig() *Config {
	return &Config{
		Enable:  true,
		Address: "[::]:37621",
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
