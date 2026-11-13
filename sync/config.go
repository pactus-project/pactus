package sync

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/moniker"
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

	// The peers drop a Hello message with such a moniker without any answer,
	// so the node would never complete a handshake.
	if len(conf.Moniker) > moniker.MaxLength {
		return ConfigError{
			Reason: fmt.Sprintf("moniker is too long: %d bytes, the maximum is %d bytes",
				len(conf.Moniker), moniker.MaxLength),
		}
	}

	if !utf8.ValidString(conf.Moniker) {
		return ConfigError{
			Reason: "moniker is not a valid UTF-8 string",
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
