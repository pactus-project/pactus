package grpc

type Config struct {
	Enable  bool          `toml:"enable"`
	Listen  string        `toml:"listen"`
	Gateway GatewayConfig `toml:"gateway"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable: false,
		Listen: "",
		Gateway: GatewayConfig{
			Enable:     false,
			Listen:    "",
			EnableCORS: false,
		},
	}
}
