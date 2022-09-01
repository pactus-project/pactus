package sync

import (
	"time"

	"github.com/zarbchain/zarb-go/sync/firewall"
)

var LatestBlockInterval = uint32(720) // 720 blocks is about two hours

type Config struct {
	Moniker          string           `toml:"moniker"`
	StartingTimeout  time.Duration    `toml:"starting_timeout"`
	HeartBeatTimeout time.Duration    `toml:"heartbeat_timeout"`
	SessionTimeout   time.Duration    `toml:"session_timeout"`
	MaxOpenSessions  int              `toml:"max_open_sessions"`
	BlockPerMessage  uint32           `toml:"block_per_message"`
	CacheSize        int              `toml:"cache_size"`
	NodeNetwork      bool             `toml:"node_network"`
	Firewall         *firewall.Config `toml:"firewall"`
}

func DefaultConfig() *Config {
	return &Config{
		StartingTimeout:  time.Second * 3,
		HeartBeatTimeout: time.Second * 5,
		SessionTimeout:   time.Second * 30,
		NodeNetwork:      true,
		BlockPerMessage:  60,
		MaxOpenSessions:  8,
		CacheSize:        50000,
		Firewall:         firewall.DefaultConfig(),
	}
}

// SanityCheck is a basic checks for config.
func (conf *Config) SanityCheck() error {
	return nil
}
