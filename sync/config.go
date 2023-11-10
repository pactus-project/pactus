package sync

import (
	"time"

	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/util"
)

type Config struct {
	Moniker        string           `toml:"moniker"`
	SessionTimeout time.Duration    `toml:"session_timeout"`
	NodeNetwork    bool             `toml:"node_network"`
	Firewall       *firewall.Config `toml:"firewall"`

	// Private configs
	MaxOpenSessions     int    `toml:"-"`
	LatestBlockInterval uint32 `toml:"-"`
	BlockPerMessage     uint32 `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		SessionTimeout:      time.Second * 10,
		NodeNetwork:         true,
		BlockPerMessage:     60,
		MaxOpenSessions:     16,
		LatestBlockInterval: 720,
		Firewall:            firewall.DefaultConfig(),
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	return nil
}

func (conf *Config) CacheSize() int {
	return util.LogScale(
		int(conf.BlockPerMessage * conf.LatestBlockInterval))
}
