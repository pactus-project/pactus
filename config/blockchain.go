package config

import "time"

type BlockchainConfig struct {
	BlockTime time.Duration
}

func DefaultBlockchainConfig() *BlockchainConfig {
	return &BlockchainConfig{
		BlockTime: 5 * time.Second,
	}
}
