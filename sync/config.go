package sync

import "time"

type Config struct {
	Moniker              string
	StartingTimeout      time.Duration
	HeartBeatTimeout     time.Duration
	SessionTimeout       time.Duration
	InitialBlockDownload bool
	BlockPerMessage      int
	RequestBlockInterval int
	CacheSize            int
}

func DefaultConfig() *Config {
	return &Config{
		StartingTimeout:      time.Second * 3,
		HeartBeatTimeout:     time.Second * 5,
		SessionTimeout:       time.Second * 3,
		InitialBlockDownload: true,
		BlockPerMessage:      10,
		RequestBlockInterval: 720,
		CacheSize:            100000,
	}
}

func TestConfig() *Config {
	return &Config{
		Moniker:              "test",
		StartingTimeout:      time.Millisecond * 100,
		HeartBeatTimeout:     time.Second * 1,
		SessionTimeout:       time.Second * 1,
		InitialBlockDownload: true,
		BlockPerMessage:      10,
		RequestBlockInterval: 20,
		CacheSize:            1000,
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
