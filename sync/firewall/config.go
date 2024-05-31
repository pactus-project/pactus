package firewall

import (
	_ "embed"
	"encoding/json"
)

//go:embed black_list.json
var _defaultBlackListIPs []byte

type Config struct {
	Enabled      bool     `toml:"enable"`
	BlackListIPs []string `toml:"blacklist_ips"`
}

type defaultBlackListIPs struct {
	IPs []string `json:"ips"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled:      false,
		BlackListIPs: make([]string, 0),
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	return nil
}

// LoadDefaultBlackListIPs load default blacklist ips from black_list.json
func (conf *Config) LoadDefaultBlackListIPs() error {
	var def defaultBlackListIPs

	if err := json.Unmarshal(_defaultBlackListIPs, &def); err != nil {
		return err
	}

	conf.BlackListIPs = append(conf.BlackListIPs, def.IPs...)
	return nil
}
