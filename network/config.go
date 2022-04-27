package network

import (
	"time"
)

type Config struct {
	Name        string           `toml:"name"`
	Listens     []string         `toml:"listens"`
	NodeKey     string           `toml:"node_key"`
	EnableDHT   bool             `toml:"enable_dht"`
	EnableNAT   bool             `toml:"enable_nat"`
	EnableRelay bool             `toml:"enable_relay"`
	EnablePing  bool             `toml:"enable_ping"`
	EnableMdns  bool             `toml:"enable_mdns"`
	Bootstrap   *BootstrapConfig `toml:"bootstrap"`
}

// BootstrapConfig holds all configuration options related to bootstrap nodes
type BootstrapConfig struct {
	Addresses    []string      `toml:"addresses"`
	MinThreshold int           `toml:"min_threshold"`
	MaxThreshold int           `toml:"max_threshold"`
	Period       time.Duration `toml:"period"`
}

func DefaultConfig() *Config {
	return &Config{
		Name:        "zarb",
		Listens:     []string{"/ip4/0.0.0.0/tcp/21888", "/ip6/::/tcp/21888"},
		NodeKey:     "node_key",
		EnableNAT:   true,
		EnableRelay: true,
		EnableMdns:  false,
		EnableDHT:   true,
		EnablePing:  true,
		Bootstrap: &BootstrapConfig{
			Addresses:    []string{"/ip4/172.104.169.94/tcp/21777/p2p/12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2epv75ATV3D8VD3Mq"},
			MinThreshold: 8,
			MaxThreshold: 16,
			Period:       1 * time.Minute,
		},
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
