package manager

import (
	"github.com/pactus-project/pactus/genesis"
)

// Config defines parameters for the wallet module.
type Config struct {
	// private config
	ChainType         genesis.ChainType `toml:"-"`
	WalletsDir        string            `toml:"-"`
	DefaultWalletName string            `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{}
}
