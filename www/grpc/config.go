package grpc

type Config struct {
	Enable  bool          `toml:""   comment:"Enable gRPC servers for client communication."`
	Address string        `toml:""  comment:"Address to listen for incoming connections for gRPC.Default port is 9090."`
	Gateway GatewayConfig `toml:""  comment:"Gateway  server which translates a RESTful HTTP API into gRPC."`
}

func DefaultConfig() *Config {
	return &Config{
		Enable:  true,
		Address: "[::]:9090",
		Gateway: GatewayConfig{
			Enable:     true,
			Address:    "[::]:8080",
			EnableCORS: false,
		},
	}
}

func TestConfig() *Config {
	return &Config{
		Enable:  true,
		Address: "[::]:0",
		Gateway: GatewayConfig{
			Enable:     true,
			Address:    "[::]:0",
			EnableCORS: false,
		},
	}
}
