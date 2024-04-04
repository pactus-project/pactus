package grpc

type Config struct {
	Enable       bool          `toml:"enable"`
	EnableWallet bool          `toml:"enable_wallet"`
	Listen       string        `toml:"listen"`
	BasicAuth    string        `toml:"basic_auth"`
	Gateway      GatewayConfig `toml:"gateway"`

	// Private config
	WalletsDir        string `toml:"-"`
	DefaultWalletName string `toml:"-"`
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
