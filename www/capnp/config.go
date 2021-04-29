package capnp

type Config struct {
	Enable  bool   `toml:"Enable" comment:"Enable CAPN servers for client communication."`
	Address string `toml:"Address" comment:"Address to listen for incoming connections for CAPNP. Default port is 37621."`
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
