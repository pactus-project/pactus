package grpc

import "github.com/pactus-project/pactus/util/htpasswd"

type Config struct {
	Enable       bool   `toml:"enable"`
	EnableWallet bool   `toml:"enable_wallet"`
	Listen       string `toml:"listen"`
	BasicAuth    string `toml:"basic_auth"`

	// Private config
	WalletsDir        string `toml:"-"`
	DefaultWalletName string `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		Enable: false,
		Listen: "",
	}
}

func (c *Config) BasicCheck() error {
	if c.BasicAuth != "" {
		if _, _, err := htpasswd.ExtractBasicAuth(c.BasicAuth); err != nil {
			return err
		}
	}

	return nil
}
