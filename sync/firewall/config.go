package firewall

import (
	_ "embed"
	"encoding/json"
	"net"
)

//go:embed black_list.json
var _defaultBlackListCidrs []byte

type RateLimit struct {
	BlockTopic       int `toml:"block_topic"`
	TransactionTopic int `toml:"transaction_topic"`
	ConsensusTopic   int `toml:"consensus_topic"`
}

type Config struct {
	BlackListAddresses []string  `toml:"blacklist_addresses"`
	RateLimit          RateLimit `toml:"rate_limit"`
}

type defaultBlackListCIDRs struct {
	Addresses []string `json:"addresses"`
}

func DefaultConfig() *Config {
	return &Config{
		BlackListAddresses: make([]string, 0),
		RateLimit: RateLimit{
			BlockTopic:       0,
			TransactionTopic: 3,
			ConsensusTopic:   0,
		},
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	for _, address := range conf.BlackListAddresses {
		_, _, err := net.ParseCIDR(address)
		if err != nil {
			return err
		}
	}

	return nil
}

// LoadDefaultBlackListAddresses loads default blacklist addresses from the `black_list.json` file.
func (conf *Config) LoadDefaultBlackListAddresses() error {
	var def defaultBlackListCIDRs

	err := json.Unmarshal(_defaultBlackListCidrs, &def)
	if err != nil {
		return err
	}

	for _, cidr := range def.Addresses {
		conf.BlackListAddresses = append(conf.BlackListAddresses, cidr)
	}

	for _, addr := range conf.BlackListAddresses {
		conf.BlackListAddresses = append(conf.BlackListAddresses, addr)
	}

	return nil
}
