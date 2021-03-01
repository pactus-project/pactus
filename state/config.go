package state

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/store"
)

// Config holds the configuration of the node
type Config struct {
	MintbaseAddress string
	Store           *store.Config
}

// DefaultConfig instantiates the default configuration for the node
func DefaultConfig() *Config {
	return &Config{
		Store: store.DefaultConfig(),
	}
}

// TestConfig instantiates the test configuration
func TestConfig() *Config {
	return &Config{
		Store: store.TestConfig(),
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	if conf.MintbaseAddress != "" {
		_, err := crypto.AddressFromString(conf.MintbaseAddress)
		if err != nil {
			return errors.Errorf(errors.ErrInvalidConfig, "Invalid mintbase address: %s", err.Error())
		}
	}
	return nil
}
