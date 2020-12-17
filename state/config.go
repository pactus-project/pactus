package state

import (
	"github.com/zarbchain/zarb-go/store"
)

// Config holds the configuration of the node
type Config struct {
	Store *store.Config
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
