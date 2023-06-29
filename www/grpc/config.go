package grpc

type Config struct {
	Listen  string        `toml:"listen"`
	Gateway GatewayConfig `toml:"gateway"`
	Enable  bool          `toml:"enable"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable: false,
		Listen: "",
		Gateway: GatewayConfig{
			Enable:     false,
			Listen:     "",
			EnableCORS: false,
		},
	}
}
