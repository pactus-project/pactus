package sync

import (
	"time"

	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
)

// ConfigError is returned when the sync configuration is invalid.
type ConfigError struct {
	Reason string
}

func (e ConfigError) Error() string {
	return e.Reason
}

// Config defines parameters for the sync module.
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
		LatestSupportingVer: version.Version{
			Major: 1,
			Minor: 15,
			Patch: 1,
		},
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	timeout, err := time.ParseDuration(conf.SessionTimeoutStr)
	if err != nil {
		return err
	}

	if timeout <= 0 {
		return ConfigError{
			Reason: "session timeout should be positive",
		}
	}

	return conf.Firewall.BasicCheck()
}

func (conf *Config) CacheSize() int {
	return util.LogScale(
		int(conf.BlockPerMessage * conf.BlockPerSession),
	)
}

func (conf *Config) SessionTimeout() time.Duration {
	timeout, _ := time.ParseDuration(conf.SessionTimeoutStr)

	return timeout
}
