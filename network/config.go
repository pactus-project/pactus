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
	Addresses        []string
	MinPeerThreshold int
	Period           time.Duration
}

func DefaultBootstrapConfig() *BootstrapConfig {
	return &BootstrapConfig{
		Addresses:        []string{},
		MinPeerThreshold: 0,
		Period:           1 * time.Minute,
	}
}

func TestBootstrapConfig() *BootstrapConfig {
	return &BootstrapConfig{
		Addresses:        []string{},
		MinPeerThreshold: 0,
		Period:           1 * time.Minute,
	}
}

func DefaultConfig() *Config {
	return &Config{
		Name:           "zarb-testnet",
		Address:        "/ip4/0.0.0.0/tcp/0",
		NodeKeyFile:    "node_key",
		EnableMDNS:     true,
		EnableKademlia: true,
		Bootstrap:      DefaultBootstrapConfig(),
	}
}

func TestConfig() *Config {
	return &Config{
		Name:           "zarb-testnet",
		Address:        "/ip4/0.0.0.0/tcp/0",
		NodeKeyFile:    util.TempFilePath(),
		EnableMDNS:     false,
		EnableKademlia: false,
		Bootstrap:      TestBootstrapConfig(),
	}
}
