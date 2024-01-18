package grpc

type Config struct {
	Enable       bool          `toml:"enable"`
	EnableWallet bool          `toml:"enable_wallet"`
	Listen       string        `toml:"listen"`
	Gateway      GatewayConfig `toml:"gateway"`

	// Private config
	WalletsDir       string `toml:"-"`
	DefaluWalletName string `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable: false,
		Listen: "50051",
		Gateway: GatewayConfig{
			Enable:     false,
			Listen:     "8080",
			EnableCORS: false,
		},
	}
}
