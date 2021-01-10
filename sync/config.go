package sync

import "time"

type Config struct {
	Moniker              string
	StartingTimeout      time.Duration
	HeartBeatTimeout     time.Duration
	InitialBlockDownload bool
	BlockPerMessage      int
	CacheSize            int
}

func DefaultConfig() *Config {
	return &Config{
		StartingTimeout:      time.Second * 10,
		HeartBeatTimeout:     time.Second * 5,
		InitialBlockDownload: true,
		BlockPerMessage:      10,
		CacheSize:            10000,
	}
}

func TestConfig() *Config {
	return &Config{
		Moniker:              "kitty",
		StartingTimeout:      time.Second * 1,
		HeartBeatTimeout:     time.Second * 1,
		InitialBlockDownload: true,
		BlockPerMessage:      10,
		CacheSize:            100,
	}
}

// SanityCheck is a basic hecks for config
func (conf *Config) SanityCheck() error {
	return nil
}
