package config

type HttpConfig struct {
	Enable  bool
	Address string
}

func DefaultHttpConfig() *HttpConfig {
	return &HttpConfig{
		Enable:  true,
		Address: "[::]:37621",
	}
}
