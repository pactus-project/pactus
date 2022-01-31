package network

import (
	"time"

	"github.com/zarbchain/zarb-go/util"
)

type Config struct {
	Name             string           `toml:"" comment:"Name dispay network name ex. “zarb”."`
	ListenAddress    []string         `toml:"" comment:"ListenAddress which support multiaddrs and it is the listining address."`
	NodeKeyFile      string           `toml:"" comment:"NodeKeyFile contains the private key to use for node authentication in the p2p protocol."`
	EnableNATService bool             `toml:"" comment:"EnableNATService NAT allows many machines to share a single public address."`
	EnableRelay      bool             `toml:"" comment:"EnableRelay is a transport protocol that routes traffic between two peers over a third-party “relay” peer."`
	EnableMDNS       bool             `toml:"" comment:"EnableMDNS is a protocol to discover local peers quickly and efficiently."`
	EnableKademlia   bool             `toml:"" comment:"EnableKademlia which is used a routing algorithm and it uses the dht routing table."`
	EnablePing       bool             `toml:"" comment:"EnablePing which enables the ping service."`
	Bootstrap        *BootstrapConfig `toml:"" comment:"Bootstrap comma separated list of peers to be added to the peer store on startup bootstrap peers."`
}

// BootstrapConfig holds all configuration options related to bootstrap nodes
type BootstrapConfig struct {
	Addresses    []string      `toml:"" comment:"Addresses it is List of peers address needed for peer discovery."`
	MinThreshold int           `toml:"" comment:"MinPeerThreshold is the number of connections it attempts to maintain."`
	MaxThreshold int           `toml:"" comment:"MaxThreshold is the threshold of maximum number of connections."`
	Period       time.Duration `toml:"" comment:"Period periodically checks to see if the threshold is maintained."`
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
		EnablePing:       true,
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
		EnablePing:       true,
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
