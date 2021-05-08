package sync

import (
	"time"

	"github.com/zarbchain/zarb-go/sync/firewall"
)

var LatestBlockInterval = 720 // 720 blocks is about two hours

type Config struct {
	Moniker              string
	StartingTimeout      time.Duration
	HeartBeatTimeout     time.Duration
	SessionTimeout       time.Duration
	InitialBlockDownload bool
	BlockPerMessage      int
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
		MaximumOpenSessions:  4,
		CacheSize:            1000,
		Firewall:             firewall.TestConfig(),
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
