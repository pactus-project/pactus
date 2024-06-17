package firewall

import (
	"net"
)

type RateLimit struct {
	BlockTopic       int `toml:"block_topic"`
	TransactionTopic int `toml:"transaction_topic"`
	ConsensusTopic   int `toml:"consensus_topic"`
}

type Config struct {
	BannedNets               []string  `toml:"banned_nets"`
	RateLimit                RateLimit `toml:"rate_limit"`
	DisallowDuplicateAddress bool      `toml:"disallow_duplicate_address"`
}

func DefaultConfig() *Config {
	return &Config{
		BannedNets: make([]string, 0),
		RateLimit: RateLimit{
			BlockTopic:       0,
			TransactionTopic: 5,
			ConsensusTopic:   0,
		},
		DisallowDuplicateAddress: true,
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	for _, address := range conf.BannedNets {
		_, _, err := net.ParseCIDR(address)
		if err != nil {
			return err
		}
	}

	return nil
}
