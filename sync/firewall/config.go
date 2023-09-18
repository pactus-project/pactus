package firewall

import "github.com/libp2p/go-libp2p/core/peer"

type Config struct {
	Enabled                 bool      `toml:"enable"`
	OnlyAcceptFromValidator bool      `toml:"only_accept_from_validator"`
	TrustedPeers            []peer.ID `toml:"trusted_peers"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled:                 false,
		OnlyAcceptFromValidator: false,
		TrustedPeers:            []peer.ID{},
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	return nil
}
