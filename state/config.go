package state

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/store"
)

// Config holds the configuration of the node
type Config struct {
	MintbaseAddress *crypto.Address
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
