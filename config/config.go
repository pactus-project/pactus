package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/pactus-project/pactus/www/http"
	"github.com/pactus-project/pactus/www/nanomsg"
	toml "github.com/pelletier/go-toml"
)

//go:embed example_config.toml
var exampleConfigBytes []byte

type Config struct {
	Node      *NodeConfig       `toml:"node"`
	Store     *store.Config     `toml:"store"`
	Network   *network.Config   `toml:"network"`
	Sync      *sync.Config      `toml:"sync"`
	TxPool    *txpool.Config    `toml:"tx_pool"`
	Consensus *consensus.Config `toml:"consensus"`
	Logger    *logger.Config    `toml:"logger"`
	GRPC      *grpc.Config      `toml:"grpc"`
	HTTP      *http.Config      `toml:"http"`
	Nanomsg   *nanomsg.Config   `toml:"nanomsg"`
}

type NodeConfig struct {
	NumValidators   int      `toml:"num_validators"` // TODO: we can remove this now
	RewardAddresses []string `toml:"reward_addresses"`
}

func DefaultNodeConfig() *NodeConfig {
	return &NodeConfig{
		NumValidators: 7,
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *NodeConfig) BasicCheck() error {
	if conf.NumValidators < 1 || conf.NumValidators > 32 {
		return errors.Errorf(errors.ErrInvalidConfig, "number of validators must be between 1 and 32")
	}

	if len(conf.RewardAddresses) > 0 &&
		len(conf.RewardAddresses) != conf.NumValidators {
		return errors.Errorf(errors.ErrInvalidConfig, "reward addresses should be %v", conf.NumValidators)
	}

	for _, addrStr := range conf.RewardAddresses {
		addr, err := crypto.AddressFromString(addrStr)
		if err != nil {
			return errors.Errorf(errors.ErrInvalidConfig, "invalid reward address: %v", err.Error())
		}

		if !addr.IsAccountAddress() {
			return errors.Errorf(errors.ErrInvalidConfig, "reward address is not an account address: %s", addrStr)
		}
	}
	return nil
}

func DefaultConfig() *Config {
	conf := &Config{
		Node:      DefaultNodeConfig(),
		Store:     store.DefaultConfig(),
		Network:   network.DefaultConfig(),
		Sync:      sync.DefaultConfig(),
		TxPool:    txpool.DefaultConfig(),
		Consensus: consensus.DefaultConfig(),
		Logger:    logger.DefaultConfig(),
		GRPC:      grpc.DefaultConfig(),
		HTTP:      http.DefaultConfig(),
		Nanomsg:   nanomsg.DefaultConfig(),
	}

	return conf
}

func SaveMainnetConfig(path string, numValidators int) error {
	conf := string(exampleConfigBytes)
	conf = strings.Replace(conf, "%num_validators%",
		fmt.Sprintf("%v", numValidators), 1)

	return util.WriteFile(path, []byte(conf))
}

func SaveTestnetConfig(path string, numValidators int) error {
	conf := DefaultConfig()
	conf.Node.NumValidators = numValidators
	conf.Network.Listens = []string{
		"/ip4/0.0.0.0/tcp/21777", "/ip4/0.0.0.0/udp/21777/quic-v1",
		"/ip6/::/tcp/21777", "/ip6/::/udp/21777/quic-v1",
	}
	conf.Network.BootstrapAddrs = []string{
		"/ip4/94.101.184.118/tcp/21777/p2p/12D3KooWCwQZt8UriVXobQHPXPR8m83eceXVoeT6brPNiBHomebc",
		"/ip4/94.101.184.118/udp/21777/quic-v1/p2p/12D3KooWCwQZt8UriVXobQHPXPR8m83eceXVoeT6brPNiBHomebc",
		"/ip4/172.104.46.145/tcp/21777/p2p/12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2epv75ATV3D8VD3Mq",
		"/ip4/172.104.46.145/udp/21777/quic-v1/p2p/12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2epv75ATV3D8VD3Mq",
		"/ip6/2400:8901::f03c:93ff:fe1c:c3ec/tcp/21777/p2p/12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2epv75ATV3D8VD3Mq",
		"/ip6/2400:8901::f03c:93ff:fe1c:c3ec/udp/21777/quic-v1/p2p/12D3KooWNYD4bB82YZRXv6oNyYPwc5ozabx2epv75ATV3D8VD3Mq",
		"/ip4/13.115.190.71/tcp/21777/p2p/12D3KooWBGNEH8NqdK1UddSnPV1yRHGLYpaQUcnujC24s7YNWPiq",
		"/ip4/13.115.190.71/udp/21777/quic-v1/p2p/12D3KooWBGNEH8NqdK1UddSnPV1yRHGLYpaQUcnujC24s7YNWPiq",
		"/ip4/51.158.118.181/tcp/21777/p2p/12D3KooWDF8a4goNCHriP1y922y4jagaPwHdX4eSrG5WtQpjzS6k",
		"/ip4/51.158.118.181/udp/21777/quic-v1/p2p/12D3KooWDF8a4goNCHriP1y922y4jagaPwHdX4eSrG5WtQpjzS6k",
		"/ip6/2001:bc8:700:8017::1/tcp/21777/p2p/12D3KooWDF8a4goNCHriP1y922y4jagaPwHdX4eSrG5WtQpjzS6k",
		"/ip6/2001:bc8:700:8017::1/udp/21777/quic-v1/p2p/12D3KooWDF8a4goNCHriP1y922y4jagaPwHdX4eSrG5WtQpjzS6k",
	}
	conf.Network.MinConns = 8
	conf.Network.MaxConns = 16
	conf.Network.EnableRelay = true
	conf.Network.RelayAddrs = []string{
		"/ip4/139.162.153.10/tcp/4002/p2p/12D3KooWNR79jqHVVNhNVrqnDbxbJJze4VjbEsBjZhz6mkvinHAN",
		"/ip4/139.162.153.10/udp/4002/quic/p2p/12D3KooWNR79jqHVVNhNVrqnDbxbJJze4VjbEsBjZhz6mkvinHAN",
		"/ip6/2a01:7e01::f03c:93ff:fed2:84c5/tcp/4002/p2p/12D3KooWNR79jqHVVNhNVrqnDbxbJJze4VjbEsBjZhz6mkvinHAN",
		"/ip6/2a01:7e01::f03c:93ff:fed2:84c5/udp/4002/quic/p2p/12D3KooWNR79jqHVVNhNVrqnDbxbJJze4VjbEsBjZhz6mkvinHAN",
		"/ip4/188.121.102.178/tcp/4002/p2p/12D3KooWCRHn8vjrKNBEQcut8uVCYX5q77RKidPaE6iMK31qEVHb",
		"/ip4/188.121.102.178/udp/4002/quic/p2p/12D3KooWCRHn8vjrKNBEQcut8uVCYX5q77RKidPaE6iMK31qEVHb",
		"/ip6/2a07:3900:1:1::113/tcp/4002/p2p/12D3KooWCRHn8vjrKNBEQcut8uVCYX5q77RKidPaE6iMK31qEVHb",
		"/ip6/2a07:3900:1:1::113/udp/4002/quic/p2p/12D3KooWCRHn8vjrKNBEQcut8uVCYX5q77RKidPaE6iMK31qEVHb",
	}
	conf.GRPC.Enable = true
	conf.GRPC.Listen = "[::]:50052"
	conf.GRPC.Gateway.Enable = true
	conf.GRPC.Gateway.Listen = "[::]:8080"
	conf.HTTP.Enable = false
	conf.HTTP.Listen = "[::]:80"
	conf.Nanomsg.Enable = false
	conf.Nanomsg.Listen = "tcp://127.0.0.1:40799"

	return util.WriteFile(path, conf.toTOML())
}

func SaveLocalnetConfig(path string, numValidators int) error {
	conf := DefaultConfig()
	conf.Node.NumValidators = numValidators
	conf.Network.Listens = []string{}
	conf.Network.EnableNAT = false
	conf.Network.BootstrapAddrs = []string{}
	conf.Network.MinConns = 0
	conf.Network.MaxConns = 0
	conf.GRPC.Enable = true
	conf.GRPC.Listen = "[::]:0"
	conf.GRPC.Gateway.Enable = true
	conf.GRPC.Gateway.Listen = "[::]:0"
	conf.HTTP.Enable = false
	conf.HTTP.Listen = "[::]:0"
	conf.Nanomsg.Enable = true
	conf.Nanomsg.Listen = "tcp://127.0.0.1:0"

	return util.WriteFile(path, conf.toTOML())
}

func (conf *Config) toTOML() []byte {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	encoder.Order(toml.OrderPreserve)
	err := encoder.Encode(conf)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func LoadFromFile(file string, strict bool) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	conf := DefaultConfig()
	buf := bytes.NewBuffer(data)
	decoder := toml.NewDecoder(buf)
	decoder.Strict(strict)
	if err := decoder.Decode(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// BasicCheck performs basic checks on the configuration.
func (conf *Config) BasicCheck() error {
	if err := conf.Node.BasicCheck(); err != nil {
		return err
	}
	if err := conf.Store.BasicCheck(); err != nil {
		return err
	}
	if err := conf.TxPool.BasicCheck(); err != nil {
		return err
	}
	if err := conf.Consensus.BasicCheck(); err != nil {
		return err
	}
	if err := conf.Network.BasicCheck(); err != nil {
		return err
	}
	if err := conf.Logger.BasicCheck(); err != nil {
		return err
	}
	if err := conf.Sync.BasicCheck(); err != nil {
		return err
	}
	if err := conf.Nanomsg.BasicCheck(); err != nil {
		return err
	}
	return conf.HTTP.BasicCheck()
}
