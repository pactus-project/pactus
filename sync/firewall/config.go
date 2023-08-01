package firewall

import "github.com/libp2p/go-libp2p/core/peer"

type Config struct {
	Enabled      bool      `toml:"enable"`
	TrustedPeers []peer.ID `toml:"trusted_peers"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled: false,
		TrustedPeers: []peer.ID{},
	}
}

// SanityCheck performs basic checks on the configuration.
func (conf *Config) SanityCheck() error {
	return nil
}
