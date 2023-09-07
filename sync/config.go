package sync

import (
	"time"

	"github.com/pactus-project/pactus/sync/firewall"
)

var LatestBlockInterval = uint32(720) // 720 blocks is about two hours

type Config struct {
	Moniker         string           `toml:"moniker"`
	SessionTimeout  time.Duration    `toml:"session_timeout"`
	BlockPerMessage uint32           `toml:"block_per_message"` // TODO: Does the user need to change it?
	CacheSize       int              `toml:"cache_size"`        // TODO: Does the user need to change it?
	NodeNetwork     bool             `toml:"node_network"`
	Firewall        *firewall.Config `toml:"firewall"`
}

func DefaultConfig() *Config {
	return &Config{
		SessionTimeout:  time.Second * 10,
		NodeNetwork:     true,
		BlockPerMessage: 60,
		CacheSize:       50000,
		Firewall:        firewall.DefaultConfig(),
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	return nil
}
