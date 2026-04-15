package manager

import (
	"github.com/pactus-project/pactus/genesis"
)

// Config defines parameters for the wallet module.
type Config struct {
	LockMode bool `toml:"lock_mode"`

	// private config
	ChainType  genesis.ChainType `toml:"-"`
	WalletsDir string            `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		LockMode: true,
	}
}

func (*Config) BasicCheck() error {
	return nil
}
