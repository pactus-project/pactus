package state

import (
	"github.com/pactus-project/pactus/crypto"
)

type Config struct {
	FoundationAddress []crypto.Address `toml:"-"`
	RewardForkHeight  uint32           `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		FoundationAddress: []crypto.Address{},
		RewardForkHeight:  0, // 4_888_000,
	}
}
