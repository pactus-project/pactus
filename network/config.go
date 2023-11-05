package network

import (
	"fmt"

	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util/errors"
)

type Config struct {
	NetworkKey           string   `toml:"network_key"`
	PublicAddrString     string   `toml:"public_address"`
	ListenAddrStrings    []string `toml:"listen_addresses"`
	RelayAddrStrings     []string `toml:"relay_addresses"`
	BootstrapAddrStrings []string `toml:"bootstrap_addresses"`
	MinConns             int      `toml:"min_connections"`
	MaxConns             int      `toml:"max_connections"`
	EnableNAT            bool     `toml:"enable_nat"`
	EnableRelay          bool     `toml:"enable_relay"`
	EnableMdns           bool     `toml:"enable_mdns"`
	EnableMetrics        bool     `toml:"enable_metrics"`
	ForcePrivateNetwork  bool     `toml:"force_private_network"`
	Bootstrapper         bool     `toml:"bootstrapper"` // TODO: detect it automatically
	DefaultPort          int      `toml:"-"`
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
		NetworkKey:       "network_key",
		PublicAddrString: "",
		ListenAddrStrings: []string{
			"/ip4/0.0.0.0/tcp/21888", "/ip6/::/tcp/21888",
			"/ip4/0.0.0.0/udp/21888/quic-v1", "/ip6/::/udp/21888/quic-v1",
		},
		RelayAddrStrings:     []string{},
		BootstrapAddrStrings: bootstrapAddrs,
		MinConns:             8,
		MaxConns:             16,
		EnableNAT:            true,
		EnableRelay:          false,
		EnableMdns:           false,
		EnableMetrics:        false,
		ForcePrivateNetwork:  false,
		Bootstrapper:         false,
		DefaultPort:          21777,
	}
}

func validateMultiAddr(addrs ...string) error {
	_, err := MakeMultiAddrs(addrs)
	return err
}

func validateAddrInfo(addrs ...string) error {
	_, err := MakeAddrInfos(addrs)
	return err
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if conf.EnableRelay {
		if len(conf.RelayAddrStrings) == 0 {
			return errors.Errorf(errors.ErrInvalidConfig, "at least one relay address should be defined")
		}
	}
	if conf.PublicAddrString != "" {
		if err := validateMultiAddr(conf.PublicAddrString); err != nil {
			return err
		}
	}
	if err := validateMultiAddr(conf.ListenAddrStrings...); err != nil {
		return err
	}
	if err := validateAddrInfo(conf.RelayAddrStrings...); err != nil {
		return err
	}
	return validateAddrInfo(conf.BootstrapAddrStrings...)
}

func (conf *Config) PublicAddr() multiaddr.Multiaddr {
	if conf.PublicAddrString != "" {
		addr, _ := multiaddr.NewMultiaddr(conf.PublicAddrString)
		return addr
	}
	return nil
}

func (conf *Config) ListenAddrs() []multiaddr.Multiaddr {
	addrs, _ := MakeMultiAddrs(conf.ListenAddrStrings)
	return addrs
}

func (conf *Config) RelayAddrInfos() []lp2ppeer.AddrInfo {
	addrInfos, _ := MakeAddrInfos(conf.RelayAddrStrings)
	return addrInfos
}

func (conf *Config) BootstrapAddrInfos() []lp2ppeer.AddrInfo {
	addrInfos, _ := MakeAddrInfos(conf.BootstrapAddrStrings)
	return addrInfos
}
