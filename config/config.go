package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"

	toml "github.com/pelletier/go-toml"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/www/capnp"
	"github.com/zarbchain/zarb-go/www/grpc"
	"github.com/zarbchain/zarb-go/www/http"
)

type Config struct {
	State     *state.Config     `toml:"" comment:"State contains the state of the blockchain."`
	Store     *store.Config     `toml:"" comment:"Store which write and store the blockchin data using golevel db. "`
	TxPool    *txpool.Config    `toml:"" comment:"TxPool is pool of unconfirmed transaction."`
	Consensus *consensus.Config `toml:"" comment:"Consensus configuration."`
	Network   *network.Config   `toml:"" comment:"Network contains all details of network configuration. Zarb uses lip2p protocol."`
	Logger    *logger.Config    `toml:"" comment:"Logger contains Output level for logging."`
	Sync      *sync.Config      `toml:"" comment:"Sync is used for peer to peer connection and synchronizing blockchain and it also contains monkier and its details."`
	Capnp     *capnp.Config     `toml:"" comment:"Cap’n Proto is an insanely fast data interchange format and capability-based RPC system."`
	HTTP      *http.Config      `toml:"" comment:"Http configuration."`
	GRPC      *grpc.Config      `toml:"" comment:"GRPC configuration."`
}

func DefaultConfig() *Config {
	conf := &Config{
		State:     state.DefaultConfig(),
		Store:     store.DefaultConfig(),
		TxPool:    txpool.DefaultConfig(),
		Consensus: consensus.DefaultConfig(),
		Network:   network.DefaultConfig(),
		Sync:      sync.DefaultConfig(),
		Logger:    logger.DefaultConfig(),
		Capnp:     capnp.DefaultConfig(),
		HTTP:      http.DefaultConfig(),
		GRPC:      grpc.DefaultConfig(),
	}

	return conf
}

func TestConfig() *Config {
	conf := &Config{
		State:     state.TestConfig(),
		Store:     store.TestConfig(),
		TxPool:    txpool.TestConfig(),
		Consensus: consensus.TestConfig(),
		Network:   network.TestConfig(),
		Sync:      sync.TestConfig(),
		Logger:    logger.TestConfig(),
		Capnp:     capnp.TestConfig(),
		HTTP:      http.TestConfig(),
		GRPC:      grpc.TestConfig(),
	}

	return conf
}

func FromTOML(t string) (*Config, error) {
	conf := DefaultConfig()

	if err := toml.Unmarshal([]byte(t), conf); err != nil {
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

	return nil, errors.Errorf(errors.ErrInvalidConfig, "invalid suffix for the config file")
}

func (conf *Config) SaveToFile(file string) error {
	var dat []byte
	if strings.HasSuffix(file, "toml") {
		dat, _ = conf.ToTOML()
	} else if strings.HasSuffix(file, "json") {
		dat, _ = conf.ToJSON()
	} else {
		return errors.Errorf(errors.ErrInvalidConfig, "invalid suffix for the config file")
	}
	return util.WriteFile(file, dat)
}

func (conf *Config) SanityCheck() error {
	if err := conf.State.SanityCheck(); err != nil {
		return err
	}
	if err := conf.Store.SanityCheck(); err != nil {
		return err
	}
	if err := conf.TxPool.SanityCheck(); err != nil {
		return err
	}
	if err := conf.Consensus.SanityCheck(); err != nil {
		return err
	}
	if err := conf.Network.SanityCheck(); err != nil {
		return err
	}
	if err := conf.Logger.SanityCheck(); err != nil {
		return err
	}
	if err := conf.Sync.SanityCheck(); err != nil {
		return err
	}
	if err := conf.Capnp.SanityCheck(); err != nil {
		return err
	}
	return conf.HTTP.SanityCheck()
}
