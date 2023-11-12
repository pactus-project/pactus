package sync

import (
	"time"

	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset/service"
	"github.com/pactus-project/pactus/util"
)

type Config struct {
	Moniker        string           `toml:"moniker"`
	SessionTimeout time.Duration    `toml:"session_timeout"`
	NodeNetwork    bool             `toml:"node_network"`
	NodeGossip     bool             `toml:"node_gossip_experimental"`
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
		NodeGossip:          false,
		BlockPerMessage:     60,
		MaxOpenSessions:     8,
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

func (conf *Config) Services() service.Services {
	s := service.New()
	if conf.NodeNetwork {
		s.Append(service.Network)
	}
	if conf.NodeGossip {
		s.Append(service.Gossip)
	}
	return s
}
