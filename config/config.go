package config

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/pactus-project/pactus/www/http"
	"github.com/pactus-project/pactus/www/jsonrpc"
	"github.com/pactus-project/pactus/www/nanomsg"
	"github.com/pelletier/go-toml"
)

var (
	//go:embed example_config.toml
	exampleConfigBytes []byte

	//go:embed bootstrap.json
	bootstrapInfosBytes []byte
)

type Config struct {
	Node      *NodeConfig       `toml:"node"`
	Store     *store.Config     `toml:"store"`
	Network   *network.Config   `toml:"network"`
	Sync      *sync.Config      `toml:"sync"`
	TxPool    *txpool.Config    `toml:"tx_pool"`
	Consensus *consensus.Config `toml:"-"`
	Logger    *logger.Config    `toml:"logger"`
	GRPC      *grpc.Config      `toml:"grpc"`
	JSONRPC   *jsonrpc.Config   `toml:"jsonrpc"`
	HTTP      *http.Config      `toml:"http"`
	Nanomsg   *nanomsg.Config   `toml:"nanomsg"`
}

type BootstrapInfo struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Website string `json:"website"`
	Address string `json:"address"`
}

type NodeConfig struct {
	RewardAddresses []string `toml:"reward_addresses"`
}

func DefaultNodeConfig() *NodeConfig {
	return &NodeConfig{
		RewardAddresses: []string{},
	}
}

// BasicCheck performs basic checks on the configuration.
func (conf *NodeConfig) BasicCheck() error {
	for _, addrStr := range conf.RewardAddresses {
		addr, err := crypto.AddressFromString(addrStr)
		if err != nil {
			return Error{
				Reason: fmt.Sprintf("invalid reward address: %v", err.Error()),
			}
		}

		if !addr.IsAccountAddress() {
			return Error{
				Reason: fmt.Sprintf("reward address is not an account address: %s", addrStr),
			}
		}
	}

	return nil
}

func defaultConfig() *Config {
	conf := &Config{
		Node:      DefaultNodeConfig(),
		Store:     store.DefaultConfig(),
		Network:   network.DefaultConfig(),
		Sync:      sync.DefaultConfig(),
		TxPool:    txpool.DefaultConfig(),
		Consensus: consensus.DefaultConfig(),
		Logger:    logger.DefaultConfig(),
		GRPC:      grpc.DefaultConfig(),
		JSONRPC:   jsonrpc.DefaultConfig(),
		HTTP:      http.DefaultConfig(),
		Nanomsg:   nanomsg.DefaultConfig(),
	}

	return conf
}

func DefaultConfigMainnet() *Config {
	conf := defaultConfig()

	bootstrapNodes := make([]BootstrapInfo, 0)
	if err := json.Unmarshal(bootstrapInfosBytes, &bootstrapNodes); err != nil {
		panic(err)
	}

	for _, node := range bootstrapNodes {
		conf.Network.DefaultBootstrapAddrStrings = append(conf.Network.DefaultBootstrapAddrStrings, node.Address)
	}

	conf.Network.MaxConns = 64
	conf.Network.EnableNATService = false
	conf.Network.EnableUPnP = false
	conf.Network.EnableRelay = true
	conf.Network.NetworkName = "pactus"
	conf.Network.DefaultPort = 21888
	conf.GRPC.Enable = true
	conf.GRPC.Listen = "127.0.0.1:50051"
	conf.GRPC.BasicAuth = ""
	conf.GRPC.Gateway.Enable = false
	conf.GRPC.Gateway.Listen = "127.0.0.1:8080"
	conf.JSONRPC.Enable = false
	conf.JSONRPC.Listen = "127.0.0.1:8545"
	conf.HTTP.Enable = false
	conf.HTTP.Listen = "127.0.0.1:80"
	conf.Nanomsg.Enable = false
	conf.Nanomsg.Listen = "tcp://127.0.0.1:40899"

	return conf
}

func DefaultConfigTestnet() *Config {
	conf := defaultConfig()
	conf.Network.DefaultBootstrapAddrStrings = []string{
		"/dns/testnet1.pactus.org/tcp/21777/p2p/12D3KooWR7ZB3nGih1Fz7Yg83Zap8Cpxr73T6PPihBsEpTG5BZyk",
		"/dns/testnet2.pactus.org/tcp/21777/p2p/12D3KooWQcDuFDMGsw6gG7oNFw7C4x7ozoMu69J7WEAojKCaNzji",
		"/dns/testnet3.pactus.org/tcp/21777/p2p/12D3KooWLsAPSJ4xowd9thGbPmbweBT6sg3nEiPjDJccaWZacsUR",
		"/dns/testnet4.pactus.org/tcp/21777/p2p/12D3KooWJKYdHzWZGibnj74NSSgKRu4Ez6MijDWMfLfXxeL4un6v",
	}
	conf.Network.MaxConns = 64
	conf.Network.EnableNATService = false
	conf.Network.EnableUPnP = false
	conf.Network.EnableRelay = true
	conf.Network.NetworkName = "pactus-testnet"
	conf.Network.DefaultPort = 21777
	conf.GRPC.Enable = true
	conf.GRPC.Listen = "[::]:50052"
	conf.GRPC.Gateway.Enable = true
	conf.GRPC.Gateway.Listen = "[::]:8080"
	conf.JSONRPC.Enable = false
	conf.JSONRPC.Listen = "127.0.0.1:8545"
	conf.HTTP.Enable = false
	conf.HTTP.Listen = "[::]:80"
	conf.Nanomsg.Enable = false
	conf.Nanomsg.Listen = "tcp://[::]:40799"

	return conf
}

func DefaultConfigLocalnet() *Config {
	conf := defaultConfig()
	conf.Network.EnableRelay = false
	conf.Network.EnableNATService = false
	conf.Network.EnableUPnP = false
	conf.Network.BootstrapAddrStrings = []string{}
	conf.Network.MaxConns = 0
	conf.Network.NetworkName = "pactus-localnet"
	conf.Network.DefaultPort = 0
	conf.Network.ForcePrivateNetwork = true
	conf.Network.EnableMdns = true
	conf.Sync.Moniker = "localnet-1"
	conf.GRPC.Enable = true
	conf.GRPC.EnableWallet = true
	conf.GRPC.Listen = "[::]:50052"
	conf.GRPC.Gateway.Enable = true
	conf.GRPC.Gateway.Listen = "[::]:8080"
	conf.JSONRPC.Enable = true
	conf.JSONRPC.Listen = "127.0.0.1:8545"
	conf.HTTP.Enable = true
	conf.HTTP.Listen = "[::]:0"
	conf.Nanomsg.Enable = true
	conf.Nanomsg.Listen = "tcp://[::]:40799"

	return conf
}

func SaveMainnetConfig(path string) error {
	conf := string(exampleConfigBytes)

	return util.WriteFile(path, []byte(conf))
}

func (conf *Config) Save(path string) error {
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

func LoadFromFile(file string, strict bool, defaultConfig *Config) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	conf := defaultConfig
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
	if err := conf.JSONRPC.BasicCheck(); err != nil {
		return err
	}

	return conf.HTTP.BasicCheck()
}
