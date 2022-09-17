package sync

import (
	"time"

	"github.com/pactus-project/pactus/sync/firewall"
)

var LatestBlockInterval = uint32(720) // 720 blocks is about two hours
var StartingTimeout = 3 * time.Second // TODO: We can remove this timer if we know we have connected to the network (min_threshould?)

type Config struct {
	Moniker         string           `toml:"moniker"`
	HeartBeatTimer  time.Duration    `toml:"heartbeat_timer"`
	SessionTimeout  time.Duration    `toml:"session_timeout"`
	MaxOpenSessions int              `toml:"max_open_sessions"`
	BlockPerMessage uint32           `toml:"block_per_message"`
	CacheSize       int              `toml:"cache_size"`
	NodeNetwork     bool             `toml:"node_network"`
	Firewall        *firewall.Config `toml:"firewall"`
}

func DefaultConfig() *Config {
	return &Config{
		HeartBeatTimer:  time.Second * 5,
		SessionTimeout:  time.Second * 30,
		NodeNetwork:     true,
		BlockPerMessage: 60,
		MaxOpenSessions: 8,
		CacheSize:       50000,
		Firewall:        firewall.DefaultConfig(),
	}
}

// SanityCheck is a basic checks for config.
func (conf *Config) SanityCheck() error {
	return nil
}
