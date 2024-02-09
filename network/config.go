package network

import (
	"fmt"

	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util"
)

type Config struct {
	NetworkKey           string   `toml:"network_key"`
	PublicAddrString     string   `toml:"public_addr"`
	ListenAddrStrings    []string `toml:"listen_addrs"`
	BootstrapAddrStrings []string `toml:"bootstrap_addrs"`
	MaxConns             int      `toml:"max_connections"`
	EnableUDP            bool     `toml:"enable_udp"`
	EnableNATService     bool     `toml:"enable_nat_service"`
	EnableUPnP           bool     `toml:"enable_upnp"`
	EnableRelay          bool     `toml:"enable_relay"`
	EnableRelayService   bool     `toml:"enable_relay_service"`
	EnableMdns           bool     `toml:"enable_mdns"`
	EnableMetrics        bool     `toml:"enable_metrics"`
	ForcePrivateNetwork  bool     `toml:"force_private_network"`

	// Private configs
	NetworkName                 string   `toml:"-"`
	DefaultPort                 int      `toml:"-"`
	DefaultBootstrapAddrStrings []string `toml:"-"`
	IsBootstrapper              bool     `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		NetworkKey:           "network_key",
		PublicAddrString:     "",
		ListenAddrStrings:    []string{},
		BootstrapAddrStrings: []string{},
		MaxConns:             64,
		EnableUDP:            false,
		EnableNATService:     false,
		EnableUPnP:           false,
		EnableRelay:          true,
		EnableRelayService:   false,
		EnableMdns:           false,
		EnableMetrics:        false,
		ForcePrivateNetwork:  false,
		DefaultPort:          0,
		IsBootstrapper:       false,
	}
}

func validateMultiAddr(addrs ...string) error {
	_, err := MakeMultiAddrs(addrs)
	if err != nil {
		return ConfigError{
			Reason: fmt.Sprintf("address is not valid: %s", err.Error()),
		}
	}

	return err
}

func validateAddrInfo(addrs ...string) error {
	_, err := MakeAddrInfos(addrs)
	if err != nil {
		return ConfigError{
			Reason: fmt.Sprintf("address is not valid: %s", err.Error()),
		}
	}

	return nil
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if conf.PublicAddrString != "" {
		if err := validateMultiAddr(conf.PublicAddrString); err != nil {
			return err
		}
	}
	if err := validateMultiAddr(conf.ListenAddrStrings...); err != nil {
		return err
	}
	if err := validateAddrInfo(conf.DefaultBootstrapAddrStrings...); err != nil {
		return err
	}
	if conf.EnableRelay && conf.EnableRelayService {
		return ConfigError{
			Reason: "both the relay and relay service cannot be active at the same time",
		}
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
	listenAddrs := conf.ListenAddrStrings
	if len(listenAddrs) == 0 {
		listenAddrs = []string{
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", conf.DefaultPort),
			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", conf.DefaultPort),
			fmt.Sprintf("/ip6/::/tcp/%d", conf.DefaultPort),
			fmt.Sprintf("/ip6/::/udp/%d/quic-v1", conf.DefaultPort),
		}
	}
	addrs, _ := MakeMultiAddrs(listenAddrs)

	return addrs
}

func (conf *Config) BootstrapAddrInfos() []lp2ppeer.AddrInfo {
	addrs := util.Merge(conf.DefaultBootstrapAddrStrings, conf.BootstrapAddrStrings)
	addrInfos, _ := MakeAddrInfos(addrs)

	return addrInfos
}

func (conf *Config) CheckIsBootstrapper(pid lp2pcore.PeerID) {
	addrInfos := conf.BootstrapAddrInfos()
	for _, ai := range addrInfos {
		if ai.ID == pid {
			conf.IsBootstrapper = true

			break
		}
	}
}

func (conf *Config) ScaledMaxConns() int {
	return util.LogScale(conf.MaxConns)
}

func (conf *Config) ScaledMinConns() int {
	return conf.ScaledMaxConns() / 4
}
