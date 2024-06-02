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

// BlackListAddresses returns the list of blacklisted addresses.
// It is a combination of user-defined addresses and pre-defined addresses in the `black_list.json` file.
func (conf *Config) GetBlackListAddresses() []string {
	var blacklisted []string

	err := json.Unmarshal(_defaultBlackListCidrs, &blacklisted)
	if err != nil {
		panic(err)
	}

	blacklisted = append(blacklisted, conf.BlackListAddresses...)

	return blacklisted
}
