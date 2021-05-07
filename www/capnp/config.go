package capnp

type Config struct {
	Enable  bool   `toml:"" comment:"Enable Cap’n proto servers for client communication."`
	Address string `toml:"" comment:"Address to listen for incoming connections for Cap’n proto. Default port is 37621."`
}

func DefaultConfig() *Config {
	return &Config{
		Enable:  false,
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
