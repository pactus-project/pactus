package zmq

import (
	"fmt"
)

// Config defines parameters for the ZeroMQ publishers.
type Config struct {
	ZmqPubBlockInfo string `toml:"zmqpubblockinfo"`
	ZmqPubTxInfo    string `toml:"zmqpubtxinfo"`
	ZmqPubRawBlock  string `toml:"zmqpubrawblock"`
	ZmqPubRawTx     string `toml:"zmqpubrawtx"`
	ZmqPubHWM       int    `toml:"zmqpubhwm"`
}

func DefaultConfig() *Config {
	return &Config{
		ZmqPubBlockInfo: "",
		ZmqPubTxInfo:    "",
		ZmqPubRawBlock:  "",
		ZmqPubRawTx:     "",
		ZmqPubHWM:       1000,
	}
}

func (c *Config) BasicCheck() error {
	if c.ZmqPubHWM < 0 {
		return fmt.Errorf("invalid publisher hwm %d", c.ZmqPubHWM)
	}

	return nil
}
