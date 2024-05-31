package firewall

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/pactus-project/pactus/util/addr"
)

//go:embed black_list.json
var _defaultBlackListAddresses []byte

type Config struct {
	Enabled            bool     `toml:"enable"`
	BlackListAddresses []string `toml:"blacklist_ips"`
	blackListAddrSet   map[string]struct{}
}

type defaultBlackListIPs struct {
	Addresses []string `json:"addresses"`
}

func DefaultConfig() *Config {
	return &Config{
		Enabled:            false,
		BlackListAddresses: make([]string, 0),
		blackListAddrSet:   make(map[string]struct{}),
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	for _, address := range conf.BlackListAddresses {
		parsed, err := addr.Parse(address)
		if err != nil {
			return fmt.Errorf("invalid blacklist address format: %s", address)
		}
		conf.blackListAddrSet[parsed.Address()] = struct{}{}
	}

	return nil
}

// LoadDefaultBlackListAddresses load default blacklist addresses from black_list.json
func (conf *Config) LoadDefaultBlackListAddresses() error {
	var def defaultBlackListIPs

	if err := json.Unmarshal(_defaultBlackListAddresses, &def); err != nil {
		return err
	}

	conf.BlackListAddresses = append(conf.BlackListAddresses, def.Addresses...)
	return nil
}
