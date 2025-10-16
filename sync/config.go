package sync

import (
	"time"

	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
)

type Config struct {
	Moniker           string           `toml:"moniker"`
	SessionTimeoutStr string           `toml:"session_timeout"`
	Firewall          *firewall.Config `toml:"firewall"`

	// Private configs
	MaxSessions         int              `toml:"-"`
	BlockPerSession     uint32           `toml:"-"`
	BlockPerMessage     uint32           `toml:"-"`
	PruneWindow         uint32           `toml:"-"`
	LatestSupportingVer version.Version  `toml:"-"`
	Services            service.Services `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		SessionTimeoutStr: "10s",
		Services:          service.New(service.PrunedNode),
		MaxSessions:       8,
		BlockPerSession:   720,
		BlockPerMessage:   60,
		PruneWindow:       86_400, // Default retention blocks in prune mode
		Firewall:          firewall.DefaultConfig(),
		// v1.9.0 is the hard-fork for Split-Reward support.
		LatestSupportingVer: version.Version{
			Major: 1,
			Minor: 9,
			Patch: 0,
		},
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	_, err := time.ParseDuration(conf.SessionTimeoutStr)
	if err != nil {
		return err
	}

	return conf.Firewall.BasicCheck()
}

func (conf *Config) CacheSize() int {
	return util.LogScale(
		int(conf.BlockPerMessage * conf.BlockPerSession))
}

func (conf *Config) SessionTimeout() time.Duration {
	timeout, _ := time.ParseDuration(conf.SessionTimeoutStr)

	return timeout
}
