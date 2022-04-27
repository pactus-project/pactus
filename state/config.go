package state

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

// Config holds the configuration of the state
type Config struct {
	RewardAddress string `toml:"reward_address"`
}

// DefaultConfig instantiates the default configuration for the state
func DefaultConfig() *Config {
	return &Config{}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	if conf.RewardAddress != "" {
		_, err := crypto.AddressFromString(conf.RewardAddress)
		if err != nil {
			return errors.Errorf(errors.ErrInvalidConfig, "invalid reward address: %v", err.Error())
		}
	}
	return nil
}
