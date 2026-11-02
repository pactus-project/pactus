package firewall

import (
	"net"
)

type RateLimit struct {
	BlockTopic       int `toml:"block_topic"`
	TransactionTopic int `toml:"transaction_topic"`
	ConsensusTopic   int `toml:"consensus_topic"`
}

// Config defines parameters for the firewall module.
type Config struct {
	BannedNets []string  `toml:"banned_nets"`
	RateLimit  RateLimit `toml:"rate_limit"`
}

func DefaultConfig() *Config {
	return &Config{
		BannedNets: make([]string, 0),
		RateLimit: RateLimit{
			BlockTopic:       1,
			TransactionTopic: 5,
			ConsensusTopic:   0,
		},
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if conf.RateLimit.BlockTopic < 0 ||
		conf.RateLimit.TransactionTopic < 0 ||
		conf.RateLimit.ConsensusTopic < 0 {
		return ConfigError{
			Reason: "rate limits can't be negative",
		}
	}

	for _, address := range conf.BannedNets {
		_, _, err := net.ParseCIDR(address)
		if err != nil {
			return err
		}
	}

	return nil
}
