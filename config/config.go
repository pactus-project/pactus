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
	State     *state.Config     `toml:"state" comment:"State contians the state of the accounts."`
	Store     *store.Config     `toml:"store" comment:"Store db which write and store the blockchin data.Default golevel db. "`
	TxPool    *txpool.Config    `toml:"txPool" comment:"TxPool is blockchain mempool.Limit the total size of all txs in the txPool."`
	Consensus *consensus.Config `toml:"Consensus" comment:"Consensus contains proposer(block creator) and validator(block validator) configuration."`
	Network   *network.Config   `toml:"Network" comment:"Network contains all details of network confgiuration. Zarb uses lip2p protocal configuration."`
	Logger    *logger.Config    `toml:"Logger" comment:"Logger contains Output level for logging."`
	Sync      *sync.Config      `toml:"Sync" comment:"Sync is used for peer connection and synchronising blockchain and it also contains monkier and its details."`
	Capnp     *capnp.Config     `toml:"Capnp" comment:"Cap’n Proto is an insanely fast data interchange format and capability-based RPC system."`
	Http      *http.Config      `toml:"Http" comment:"TCP or UNIX socket address for the tcp server to listen on. Default port is 8081."`
	GRPC      *grpc.Config      `toml:"GRPC" comment:"TCP or UNIX socket address for zarb to listen on for connections from an external process. Default port is 9090."`
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
		Http:      http.DefaultConfig(),
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
		Http:      http.TestConfig(),
		GRPC:      grpc.TestConfig(),
	}

	return conf
}

func FromTOML(t string) (*Config, error) {
	conf := DefaultConfig()

	if _, err := toml.Load(t); err != nil {
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
	if err := util.WriteFile(file, dat); err != nil {
		return err
	}

	return nil
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
	if err := conf.Http.SanityCheck(); err != nil {
		return err
	}
	return nil
}
