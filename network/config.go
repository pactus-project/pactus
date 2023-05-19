package network

import (
	"fmt"
	"time"
)

type Config struct {
	Name        string           `toml:"name"`
	Listens     []string         `toml:"listens"`
	RelayAddrs  []string         `toml:"relays"`
	NetworkKey  string           `toml:"network_key"`
	EnableNAT   bool             `toml:"enable_nat"`
	EnableRelay bool             `toml:"enable_relay"`
	EnableMdns  bool             `toml:"enable_mdns"`
	Bootstrap   *BootstrapConfig `toml:"bootstrap"`
}

// BootstrapConfig holds all configuration options related to bootstrap nodes.
type BootstrapConfig struct {
	Addresses    []string      `toml:"addresses"`
	MinThreshold int           `toml:"min_threshold"`
	MaxThreshold int           `toml:"max_threshold"`
	Period       time.Duration `toml:"period"`
}

func DefaultConfig() *Config {
	nodes := []struct {
		ip   string
		port string
		id   string
	}{
		{
			ip:   "172.104.46.145",
			port: "21777",
			id:   "12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2epv75ATV3D8VD3Mq",
		},
	}

	addresses := []string{}
	for _, n := range nodes {
		addresses = append(addresses,
			fmt.Sprintf("/ip4/%s/tcp/%s/p2p/%s", n.ip, n.port, n.id))
	}

	return &Config{
		Name:        "pactus",
		Listens:     []string{"/ip4/0.0.0.0/tcp/21777", "/ip6/::/tcp/21777"},
		NetworkKey:  "network_key",
		EnableNAT:   true,
		EnableRelay: true,
		EnableMdns:  false,
		Bootstrap: &BootstrapConfig{
			Addresses:    addresses,
			MinThreshold: 8,
			MaxThreshold: 16,
			Period:       1 * time.Minute,
		},
	}
}

// SanityCheck performs basic checks on the configuration.
func (conf *Config) SanityCheck() error {
	return nil
}
