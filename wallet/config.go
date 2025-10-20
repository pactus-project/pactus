package wallet

import (
	"github.com/pactus-project/pactus/genesis"
)

// Config defines parameters for the wallet module.
type Config struct {
	// private config
	WalletsDir string            `toml:"-"`
	ChainType  genesis.ChainType `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{}
}
