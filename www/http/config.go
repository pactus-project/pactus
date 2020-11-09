package http

type Config struct {
	Enable      bool
	Address     string
	CapnpServer string
}

func DefaultConfig() *Config {
	return &Config{
		Enable:      true,
		Address:     "[::]:8081",
		CapnpServer: "[::]:37621",
	}
}

func TestConfig() *Config {
	return &Config{
		Enable: false,
	}
}
