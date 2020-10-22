package http

type Config struct {
	Enable      bool
	Address     string
	CapnpServer string
}

func DefaultConfig() *Config {
	return &Config{
		Enable:      true,
		Address:     "[::]:8080",
		CapnpServer: "[::]:37621",
	}
}
