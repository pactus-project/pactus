package firewall

import (
	"net"
)

type Config struct {
	BannedNets []string `toml:"banned_nets"`
}

func DefaultConfig() *Config {
	return &Config{
		BannedNets: make([]string, 0),
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	for _, address := range conf.BannedNets {
		_, _, err := net.ParseCIDR(address)
		if err != nil {
			return err
		}
	}

	return nil
}
