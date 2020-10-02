package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"gitlab.com/zarb-chain/zarb-go/errors"
	"gitlab.com/zarb-chain/zarb-go/utils"
)

type Config struct {
	Store     *StoreConfig
	TxPool    *TxPoolConfig
	Consensus *ConsensusConfig
	Network   *NetworkConfig
	Logger    *LoggerConfig
	Capnp     *CapnpConfig
	Http      *HttpConfig
	// This are private and set by mint-account
	blockchain *BlockchainConfig
}

func DefaultConfig() *Config {

	conf := &Config{
		Store:      DefaultStoreConfig(),
		TxPool:     DefaultTxPoolConfig(),
		Consensus:  DefaultConsensusConfig(),
		Network:    DefaultNetworkConfig(),
		Logger:     DefaultLoggerConfig(),
		Capnp:      DefaultCapnpConfig(),
		Http:       DefaultHttpConfig(),
		blockchain: DefaultBlockchainConfig(),
	}

	return conf
}

func FromTOML(t string) (*Config, error) {
	conf := DefaultConfig()

	if _, err := toml.Decode(t, conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func (conf *Config) ToTOML() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	err := encoder.Encode(conf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func FromJSON(t string) (*Config, error) {
	conf := DefaultConfig()
	if err := json.Unmarshal([]byte(t), conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func (conf *Config) ToJSON() ([]byte, error) {
	return json.MarshalIndent(conf, "", "  ")
}

func LoadFromFile(file string) (*Config, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(file, "toml") {
		return FromTOML(string(dat))
	} else if strings.HasSuffix(file, "json") {
		return FromJSON(string(dat))
	}

	return nil, errors.Errorf(errors.ErrInvalidConfig, "Invalid suffix for the config file")
}

func (conf *Config) SaveToFile(file string) error {
	var dat []byte
	if strings.HasSuffix(file, "toml") {
		dat, _ = conf.ToTOML()
	} else if strings.HasSuffix(file, "json") {
		dat, _ = conf.ToJSON()
	} else {
		return errors.Errorf(errors.ErrInvalidConfig, "Invalid suffix for the config file")
	}
	if err := utils.WriteFile(file, dat); err != nil {
		return err
	}

	return nil
}

func (conf *Config) Check() error {
	return nil
}

func (conf *Config) BlockTime() time.Duration {
	return conf.blockchain.BlockTime
}
