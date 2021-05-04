package sync

import (
	"time"

	"github.com/zarbchain/zarb-go/sync/firewall"
)

type Config struct {
	Moniker              string
	StartingTimeout      time.Duration
	HeartBeatTimeout     time.Duration
	SessionTimeout       time.Duration
	InitialBlockDownload bool
	BlockPerMessage      int
	RequestBlockInterval int
	MaximumOpenSessions  int
	CacheSize            int
	Firewall             *firewall.Config
}

func DefaultConfig() *Config {
	return &Config{
		StartingTimeout:      time.Second * 3,
		HeartBeatTimeout:     time.Second * 5,
		SessionTimeout:       time.Second * 30,
		InitialBlockDownload: true,
		BlockPerMessage:      120,
		RequestBlockInterval: 720,
		MaximumOpenSessions:  8,
		CacheSize:            500000,
		Firewall:             firewall.DefaultConfig(),
	}
}

func TestConfig() *Config {
	return &Config{
		Moniker:              "test",
		StartingTimeout:      0,
		HeartBeatTimeout:     time.Second * 1,
		SessionTimeout:       time.Second * 1,
		InitialBlockDownload: true,
		BlockPerMessage:      10,
		RequestBlockInterval: 20,
		MaximumOpenSessions:  4,
		CacheSize:            1000,
		Firewall:             firewall.TestConfig(),
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
