package grpc

type Config struct {
	Enable  bool
	Address string
	Gateway GatewayConfig
}

func DefaultConfig() *Config {
	return &Config{
		Enable:  true,
		Address: "[::]:9090",
		Gateway: GatewayConfig{
			Enable:  true,
			Address: "[::]:8080",
		},
	}
}

func TestConfig() *Config {
	return &Config{
		Enable:  true,
		Address: "[::]:0",
		Gateway: GatewayConfig{
			Enable:  true,
			Address: "[::]:0",
		},
	}
}
