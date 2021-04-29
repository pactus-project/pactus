package network

import (
	"time"

	"github.com/zarbchain/zarb-go/util"
)

type Config struct {
	Name             string           `toml:"Name" comment:"Name which is use to dispay network name."`
	ListenAddress    []string         `toml:"ListenAddress" comment:"ListenAddress which support multiaddrs."`
	NodeKeyFile      string           `toml:"NodeKeyFile" comment:"NodeKeyFile contains the private key to use for node authentication in the p2p protocol."`
	EnableNATService bool             `toml:"EnableNATService" comment:"EnableNATService NAT allows many machines to share a single public address."`
	EnableRelay      bool             `toml:"EnableRelay" comment:"EnableRelay is a transport protocol that routes traffic between two peers over a third-party “relay” peer."`
	EnableMDNS       bool             `toml:"EnableMDNS" comment:"EnableMDNS is a protocol to discover local peers quickly and efficiently."`
	EnableKademlia   bool             `toml:"EnableKademlia" comment:"EnableKademlia Kademlia routing algorithm.which uses the dht routing table."`
	Bootstrap        *BootstrapConfig `toml:"Bootstrap" comment:"Bootstrap comma separated list of peers to be added to the peer store on startup bootstrap peers."`
}

// BootstrapConfig holds all configuration options related to bootstrap nodes
type BootstrapConfig struct {
	Addresses    []string      `toml:"Addresses" comment:"Addresses it is List of peers address needed for peer discovery."`
	MinThreshold int           `toml:"MinThreshold" comment:"MinPeerThreshold is the number of connections it attempts to maintain."`
	MaxThreshold int           `toml:"MaxThreshold" comment:"MaxThreshold is the threshold of maximum number of connections."`
	Period       time.Duration `toml:"Period" comment:"Period periodically checks to see if the threshold is maintained."`
}

func DefaultConfig() *Config {
	return &Config{
		Name:             "zarb",
		ListenAddress:    []string{"/ip4/0.0.0.0/tcp/0", "/ip6/::/tcp/0"},
		NodeKeyFile:      "node_key",
		EnableNATService: true,
		EnableRelay:      true,
		EnableMDNS:       true,
		EnableKademlia:   true,
		Bootstrap: &BootstrapConfig{
			Addresses:    []string{},
			MinThreshold: 8,
			MaxThreshold: 16,
			Period:       1 * time.Minute,
		},
	}
}

func TestConfig() *Config {
	return &Config{
		Name:             "zarb-testnet",
		ListenAddress:    []string{"/ip4/0.0.0.0/tcp/0", "/ip6/::/tcp/0"},
		NodeKeyFile:      util.TempFilePath(),
		EnableNATService: true,
		EnableRelay:      true,
		EnableMDNS:       true,
		EnableKademlia:   true,
		Bootstrap: &BootstrapConfig{
			Addresses:    []string{},
			MinThreshold: 4,
			MaxThreshold: 8,
			Period:       1 * time.Minute,
		},
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
