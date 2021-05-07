package sync

import (
	"time"

	"github.com/zarbchain/zarb-go/sync/firewall"
)

type Config struct {
	Moniker              string        `toml:"" comment:"Moniker A custom human readable name for this node."`
	StartingTimeout      time.Duration `toml:"" comment:"StartingTimeout is time taken for syncing the node."`
	HeartBeatTimeout     time.Duration `toml:"" comment:"HeartBeatTimeout timeout for broadcasting heartbeat message to network."`
	SessionTimeout       time.Duration `toml:"" comment:"SessionTimeout timeout for session of node."`
	InitialBlockDownload bool          `toml:"" comment:"InitialBlockDownload enable or disable for initial block downloading."`
	BlockPerMessage      int           `toml:"" comment:"BlockPerMessage the number of blocks per message. Default is 120."`
	RequestBlockInterval int
	CacheSize            int              `toml:"" comment:"CacheSize is the total capacity of the cache"`
	Firewall             *firewall.Config `toml:"" comment:"Setting for firewall"`
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
		CacheSize:            1000,
		Firewall:             firewall.TestConfig(),
	}
}

// SanityCheck is a basic checks for config
func (conf *Config) SanityCheck() error {
	return nil
}
