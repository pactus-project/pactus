package state

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

// Config holds the configuration of the node
type Config struct {
	MintbaseAddress string
}

// DefaultConfig instantiates the default configuration for the node
func DefaultConfig() *Config {
	return &Config{}
}

// TestConfig instantiates the test configuration
func TestConfig() *Config {
	return &Config{}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	if conf.MintbaseAddress != "" {
		_, err := crypto.AddressFromString(conf.MintbaseAddress)
		if err != nil {
			return errors.Errorf(errors.ErrInvalidConfig, "invalid mintbase address: %s", err.Error())
		}
	}
	return nil
}
