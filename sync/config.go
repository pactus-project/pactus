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
	EnableFirewall       bool
}

func DefaultConfig() *Config {
	return &Config{
		StartingTimeout:      time.Second * 3,
		HeartBeatTimeout:     time.Second * 5,
		SessionTimeout:       time.Second * 3,
		InitialBlockDownload: true,
		BlockPerMessage:      10,
		RequestBlockInterval: 720,
		CacheSize:            500000,
		EnableFirewall:       true,
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
		CacheSize:            1000,
		EnableFirewall:       false,
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
