package network

import (
	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pactus-project/pactus/util"
)

type Config struct {
	NetworkKey           string   `toml:"network_key"`
	PublicAddrString     string   `toml:"public_addr"`
	ListenAddrStrings    []string `toml:"listen_addrs"`
	RelayAddrStrings     []string `toml:"relay_addrs"`
	BootstrapAddrStrings []string `toml:"bootstrap_addrs"`
	MaxConns             int      `toml:"max_connections"`
	EnableNATService     bool     `toml:"enable_nat_service"`
	EnableUPnP           bool     `toml:"enable_upnp"`
	EnableRelay          bool     `toml:"enable_relay"`
	EnableMdns           bool     `toml:"enable_mdns"`
	EnableMetrics        bool     `toml:"enable_metrics"`
	ForcePrivateNetwork  bool     `toml:"force_private_network"`

	// Private configs
	NetworkName                 string   `toml:"-"`
	DefaultPort                 int      `toml:"-"`
	DefaultRelayAddrStrings     []string `toml:"-"`
	DefaultBootstrapAddrStrings []string `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		NetworkKey:           "network_key",
		PublicAddrString:     "",
		ListenAddrStrings:    []string{},
		RelayAddrStrings:     []string{},
		BootstrapAddrStrings: []string{},
		MaxConns:             64,
		EnableNATService:     false,
		EnableUPnP:           false,
		EnableRelay:          false,
		EnableMdns:           false,
		EnableMetrics:        false,
		ForcePrivateNetwork:  false,
		DefaultPort:          21888,
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
	if err := validateAddrInfo(conf.DefaultRelayAddrStrings...); err != nil {
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
	addrs := util.Merge(conf.DefaultRelayAddrStrings, conf.RelayAddrStrings)
	addrInfos, _ := MakeAddrInfos(addrs)
	return addrInfos
}

func (conf *Config) BootstrapAddrInfos() []lp2ppeer.AddrInfo {
	addrs := util.Merge(conf.DefaultBootstrapAddrStrings, conf.BootstrapAddrStrings)
	addrInfos, _ := MakeAddrInfos(addrs)
	return addrInfos
}

func (conf *Config) IsBootstrapper(pid lp2pcore.PeerID) bool {
	addrInfos := conf.BootstrapAddrInfos()
	for _, ai := range addrInfos {
		if ai.ID == pid {
			return true
		}
	}

	return false
}

func (conf *Config) ScaledMaxConns() int {
	return util.LogScale(conf.MaxConns)
}

func (conf *Config) ScaledMinConns() int {
	return conf.ScaledMaxConns() / 4
}

func (conf *Config) ConnsThreshold() int {
	return conf.ScaledMaxConns() / 8
}
