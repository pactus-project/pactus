package wallet

import (
	"github.com/pactus-project/pactus/genesis"
)

type Config struct {
	// private config
	WalletsDir string            `toml:"-"`
	ChainType  genesis.ChainType `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{}
}
