package zmq

import (
	"fmt"
)

// Config defines parameters for the ZeroMQ publishers.
type Config struct {
	ZmqPubBlockInfo  string `toml:"zmqpubblockinfo"`
	ZmqPubTxInfo     string `toml:"zmqpubtxinfo"`
	ZmqPubRawBlock   string `toml:"zmqpubrawblock"`
	ZmqPubRawTx      string `toml:"zmqpubrawtx"`
	ZmqPubAnchorInfo string `toml:"zmqpubanchorinfo"`
	ZmqPubHWM        int    `toml:"zmqpubhwm"`
}

func DefaultConfig() *Config {
	return &Config{
		ZmqPubBlockInfo:  "",
		ZmqPubTxInfo:     "",
		ZmqPubRawBlock:   "",
		ZmqPubRawTx:      "",
		ZmqPubAnchorInfo: "",
		ZmqPubHWM:        1000,
	}
}

func (c *Config) BasicCheck() error {
	// The high water mark bounds the messages queued for a subscriber that doesn't read them.
	// Zero means no bound: such a subscriber would make the queue grow without limit.
	if c.ZmqPubHWM <= 0 {
		return fmt.Errorf("invalid publisher hwm %d: it must be positive", c.ZmqPubHWM)
	}

	return nil
}
