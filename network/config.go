package network

import (
	"fmt"
	"time"

	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/errors"
)

type Config struct {
	Listens            []string         `toml:"listens"`
	NetworkKey         string           `toml:"network_key"`
	EnableNAT          bool             `toml:"enable_nat"`
	EnableRelay        bool             `toml:"enable_relay"`
	EnableHolePunching bool             `toml:"enable_hole_punching"`
	RelayAddrs         []string         `toml:"relay_addresses"`
	EnableMdns         bool             `toml:"enable_mdns"`
	EnableMetrics      bool             `toml:"enable_metrics"`
	Bootstrap          *BootstrapConfig `toml:"bootstrap"`
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
		Listens:            []string{"/ip4/0.0.0.0/tcp/21777", "/ip6/::/tcp/21777"},
		NetworkKey:         "network_key",
		EnableNAT:          true,
		EnableRelay:        false,
		EnableHolePunching: true,
		EnableMdns:         false,
		EnableMetrics:      false,
		Bootstrap: &BootstrapConfig{
			Addresses:    addresses,
			MinThreshold: 8,
			MaxThreshold: 16,
			Period:       1 * time.Minute,
		},
	}
}

func validateAddresses(address []string) error {
	for _, addr := range address {
		_, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return err
		}
	}
	return nil
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if conf.EnableRelay || conf.EnableHolePunching {
		if len(conf.RelayAddrs) == 0 {
			return errors.Errorf(errors.ErrInvalidConfig, "at least one relay address should be defined")
		}
	}
	if err := validateAddresses(conf.RelayAddrs); err != nil {
		return err
	}
	return validateAddresses(conf.Listens)
}
