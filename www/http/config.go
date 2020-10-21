package http

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
