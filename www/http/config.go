package http

type Config struct {
	Enable  bool
	Address string
}

func DefaultConfig() *Config {
	return &Config{
		Enable:  true,
		Address: "[::]:8080",
	}
}

func TestConfig() *Config {
	return &Config{
		Enable:  true,
		Address: "[::]:0",
	}
}
