package firewall

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/pactus-project/pactus/util/addr"
)

//go:embed black_list.json
var _defaultBlackListAddresses []byte

type RateLimit struct {
	BlockTopic       int `toml:"block_topic"`
	TransactionTopic int `toml:"transaction_topic"`
	ConsensusTopic   int `toml:"consensus_topic"`
}

type Config struct {
	BlackListAddresses []string  `toml:"blacklist_addresses"`
	RateLimit          RateLimit `toml:"rate_limit"`

	blackListAddrSet map[string]any
}

type defaultBlackListIPs struct {
	Addresses []string `json:"addresses"`
}

func DefaultConfig() *Config {
	return &Config{
		BlackListAddresses: make([]string, 0),
		blackListAddrSet:   make(map[string]any),
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
		// TODO: use libp2p library (multi-address)
		// TODO: address should only contain protocol + address like: "/ip4/1.1.1.1"
		_, err := addr.Parse(address)
		if err != nil {
			return fmt.Errorf("invalid blacklist address format: %s", address)
		}
	}

	return nil
}

// LoadDefaultBlackListAddresses loads default blacklist addresses from the `black_list.json` file.
func (conf *Config) LoadDefaultBlackListAddresses() {
	var def defaultBlackListIPs

	_ = json.Unmarshal(_defaultBlackListAddresses, &def)

	for _, a := range def.Addresses {
		ma, _ := addr.Parse(a)
		conf.blackListAddrSet[ma.Address()] = true
	}

	for _, a := range conf.BlackListAddresses {
		ma, _ := addr.Parse(a)
		conf.blackListAddrSet[ma.Address()] = true
	}
}
