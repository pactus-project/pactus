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
		Enable: false,
	}
}
