package sync

import "time"

type Config struct {
	Moniker              string        `toml:"Moniker" comment:"Moniker A custom human readable name for this node."`
	StartingTimeout      time.Duration `toml:"StartingTimeout" comment:"StartingTimeout is timeout to syncing the node."`
	HeartBeatTimeout     time.Duration `toml:"HeartBeatTimeout" comment:"HeartBeatTimeout timeout for detecting the tcp connections."`
	SessionTimeout       time.Duration `toml:"SessionTimeout" comment:"SessionTimeout session."`
	InitialBlockDownload bool          `toml:"InitialBlockDownload" comment:"InitialBlockDownload enable or disable for initial block downloading."`
	BlockPerMessage      int           `toml:"BlockPerMessage" comment:"BlockPerMessage receive the number of blocks per message. Default is 10."`
	RequestBlockInterval int           `toml:"RequestBlockInterval" comment:"RequestBlockInterval max duration for a request, block interval."`
	CacheSize            int           `toml:"CacheSize" comment:"CacheSize Size of the cache in transactions."`
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
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
