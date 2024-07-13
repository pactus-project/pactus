package sync

import (
	"time"

	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
)

type Config struct {
	Moniker        string           `toml:"moniker"`
	SessionTimeout time.Duration    `toml:"session_timeout"`
	Firewall       *firewall.Config `toml:"firewall"`

	// Private configs
	MaxSessions         int              `toml:"-"`
	LatestBlockInterval uint32           `toml:"-"`
	BlockPerMessage     uint32           `toml:"-"`
	LatestSupportingVer version.Version  `toml:"-"`
	Services            service.Services `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		SessionTimeout:      time.Second * 10,
		Services:            service.New(service.PrunedNode),
		BlockPerMessage:     60,
		MaxSessions:         8,
		LatestBlockInterval: 10 * 8640, // 10 days, same as default retention blocks in prune node
		Firewall:            firewall.DefaultConfig(),
		LatestSupportingVer: version.Version{
			Major: 1,
			Minor: 1,
			Patch: 8,
		},
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	return conf.Firewall.BasicCheck()
}

func (conf *Config) CacheSize() int {
	return util.LogScale(
		int(conf.BlockPerMessage * conf.LatestBlockInterval))
}
