package network

import (
	"time"

	"github.com/zarbchain/zarb-go/util"
)

type Config struct {
	Name           string
	Address        string
	NodeKeyFile    string
	EnableMDNS     bool
	EnableKademlia bool
	Bootstrap      *BootstrapConfig
}

// BootstrapConfig holds all configuration options related to bootstrap nodes
type BootstrapConfig struct {
	// Peers to connect to if we fall below the threshold.
	Addresses []string
	// MinPeerThreshold is the number of connections it attempts to maintain.
	MinThreshold int
	// MaxThreshold is the threshold of maximum number of connections.
	MaxThreshold int
	// Period is the interval at which it periodically checks to see
	// if the threshold is maintained.
	Period time.Duration
	// ConnectionTimeout is how long to wait before timing out a connection attempt.
	Timeout time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		Name:           "zarb",
		Address:        "/ip4/0.0.0.0/tcp/0",
		NodeKeyFile:    "node_key",
		EnableMDNS:     true,
		EnableKademlia: true,
		Bootstrap: &BootstrapConfig{
			Addresses:    []string{},
			MinThreshold: 8,
			MaxThreshold: 16,
			Period:       1 * time.Minute,
			Timeout:      20 * time.Second,
		},
	}
}

func TestConfig() *Config {
	return &Config{
		Name:           "zarb-testnet",
		Address:        "/ip4/0.0.0.0/tcp/0",
		NodeKeyFile:    util.TempFilePath(),
		EnableMDNS:     true,
		EnableKademlia: true,
		Bootstrap: &BootstrapConfig{
			Addresses:    []string{},
			MinThreshold: 4,
			MaxThreshold: 8,
			Period:       1 * time.Minute,
			Timeout:      20 * time.Second,
		},
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
