package network

import (
	"fmt"

	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/errors"
)

type Config struct {
	NetworkKey     string   `toml:"network_key"`
	Listens        []string `toml:"listens"`
	RelayAddrs     []string `toml:"relay_addresses"`
	BootstrapAddrs []string `toml:"bootstrap_addresses"`
	MinConns       int      `toml:"min_connections"`
	MaxConns       int      `toml:"max_connections"`
	EnableNAT      bool     `toml:"enable_nat"`
	EnableRelay    bool     `toml:"enable_relay"`
	EnableMdns     bool     `toml:"enable_mdns"`
	EnableMetrics  bool     `toml:"enable_metrics"`
	Bootstrapper   bool     `toml:"bootstrapper"`
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

	bootstrapAddrs := []string{}
	for _, n := range nodes {
		bootstrapAddrs = append(bootstrapAddrs,
			fmt.Sprintf("/ip4/%s/tcp/%s/p2p/%s", n.ip, n.port, n.id))
	}

	return &Config{
		NetworkKey: "network_key",
		Listens: []string{
			"/ip4/0.0.0.0/tcp/21888", "/ip6/::/tcp/21888",
			"/ip4/0.0.0.0/udp/21888/quic-v1", "/ip6/::/udp/21888/quic-v1",
		},
		RelayAddrs:     []string{},
		BootstrapAddrs: bootstrapAddrs,
		MinConns:       8,
		MaxConns:       16,
		EnableNAT:      true,
		EnableRelay:    false,
		EnableMdns:     false,
		EnableMetrics:  false,
		Bootstrapper:   false,
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
	if conf.EnableRelay {
		if len(conf.RelayAddrs) == 0 {
			return errors.Errorf(errors.ErrInvalidConfig, "at least one relay address should be defined")
		}
	}
	if err := validateAddresses(conf.RelayAddrs); err != nil {
		return err
	}
	if err := validateAddresses(conf.Listens); err != nil {
		return err
	}
	return validateAddresses(conf.BootstrapAddrs)
}
