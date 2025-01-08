package zmq

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	ZmqPubBlockInfo string `toml:"zmqpubblockinfo"`
	ZmqPubTxInfo    string `toml:"zmqpubtxinfo"`
	ZmqPubRawBlock  string `toml:"zmqpubrawblock"`
	ZmqPubRawTx     string `toml:"zmqpubrawtx"`
	ZmqPubHWM       int    `toml:"zmqpubhwm"`

	// Private config
	ZmqAutomaticReconnect bool          `toml:"-"`
	ZmqDialerRetryTime    time.Duration `toml:"-"`
	ZmqDialerMaxRetries   int           `toml:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		ZmqAutomaticReconnect: true,
		ZmqDialerMaxRetries:   10,
		ZmqDialerRetryTime:    250 * time.Millisecond,
		ZmqPubHWM:             1000,
	}
}

func (c *Config) BasicCheck() error {
	if c.ZmqPubBlockInfo != "" {
		if err := validateTopicSocket(c.ZmqPubBlockInfo); err != nil {
			return err
		}
	}

	if c.ZmqPubTxInfo != "" {
		if err := validateTopicSocket(c.ZmqPubTxInfo); err != nil {
			return err
		}
	}

	if c.ZmqPubRawBlock != "" {
		if err := validateTopicSocket(c.ZmqPubRawBlock); err != nil {
			return err
		}
	}

	if c.ZmqPubRawTx != "" {
		if err := validateTopicSocket(c.ZmqPubRawTx); err != nil {
			return err
		}
	}

	if c.ZmqPubHWM < 0 {
		return fmt.Errorf("invalid publisher hwm %d", c.ZmqPubHWM)
	}

	return nil
}

func validateTopicSocket(socket string) error {
	addr, err := url.Parse(socket)
	if err != nil {
		return errors.New("failed to parse ZmqPub value: " + err.Error())
	}

	if addr.Scheme != "tcp" {
		return errors.New("invalid scheme: zeromq socket schema")
	}

	if addr.Host == "" {
		return errors.New("invalid host: host is empty")
	}

	parts := strings.Split(addr.Host, ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return errors.New("invalid host: missing or malformed host/port")
	}

	port := parts[1]
	for _, r := range port {
		if r < '0' || r > '9' {
			return errors.New("invalid port: non-numeric characters detected")
		}
	}

	return nil
}
