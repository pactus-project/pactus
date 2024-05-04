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
)

var (
	//go:embed example_config.toml
	exampleConfigBytes []byte

	//go:embed bootstrap.json
	bootstrapInfoBytes []byte

	//go:embed banned.json
	bannedBytes []byte
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
			return NodeConfigError{
				Reason: fmt.Sprintf("invalid reward address: %v", err.Error()),
			}
		}

		if !addr.IsAccountAddress() {
			return NodeConfigError{
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
	if err := json.Unmarshal(bootstrapInfoBytes, &bootstrapNodes); err != nil {
		panic(err)
	}

	bootstrapAddrs := []string{}
	for _, node := range bootstrapNodes {
		bootstrapAddrs = append(bootstrapAddrs, node.Address)
	}

	// The first item in the banned list is for testing.
	// The address is: "pc1p8slveave2zm9tgj7q260fgrdfu2ph8n7ezxhtt"
	// The associated private key: "SECRET1PP8SYQAH8JH8QLGEEX7L7T8WTU69K6T9AVSNMVCZ8DP6PPLWVYE3SVTHFR8"
	bannedList := make([]string, 0)
	if err := json.Unmarshal(bannedBytes, &bannedList); err != nil {
		panic(err)
	}

	bannedAddrs := make(map[crypto.Address]bool)
	for _, str := range bannedList {
		addr, err := crypto.AddressFromString(str)
		if err != nil {
			panic(err)
		}
		bannedAddrs[addr] = true
	}

	conf.Store.BannedAddrs = bannedAddrs
	conf.Network.MaxConns = 64
	conf.Network.EnableNATService = false
	conf.Network.EnableUPnP = false
	conf.Network.EnableRelay = true
	conf.Network.NetworkName = "pactus"
	conf.Network.DefaultPort = 21888
	conf.Network.DefaultBootstrapAddrStrings = bootstrapAddrs
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
		"/ip4/34.35.39.182/tcp/21777/p2p/12D3KooWJKYdHzWZGibnj74NSSgKRu4Ez6MijDWMfLfXxeL4un6v",
		"/ip4/128.140.41.234/tcp/21777/p2p/12D3KooWLV9Y2MzYMVMtqCuivFAzcgVM1H7US4Vud1NU6KFmJBkw",
		"/ip4/34.18.2.86/tcp/21777/p2p/12D3KooWLsAPSJ4xowd9thGbPmbweBT6sg3nEiPjDJccaWZacsUR",
		"/ip4/113.176.163.161/tcp/21777/p2p/12D3KooWCNE13y2qh9W4qZFd61Me8nDDANGx36j2zWTz3sscBJou",
		"/ip4/35.234.166.185/tcp/21777/p2p/12D3KooWQcDuFDMGsw6gG7oNFw7C4x7ozoMu69J7WEAojKCaNzji",
		"/ip4/49.13.205.202/tcp/21777/p2p/12D3KooWLQdsUExqTux9T6xKbXE3JbgsEFx3Gu9AhCSV3L1ZuMNn",
		"/ip4/104.234.1.82/tcp/21777/p2p/12D3KooWLnyM3SLvye7f2BW5AH1ZQQvkvZxKr9yaCNEPfFY4cbYx",
		"/ip4/207.180.250.57/tcp/21777/p2p/12D3KooWFpSTV3r6auPeWfKM1vb4FNug6FeAuJkgmzPTUJocxD39",
		"/ip4/202.182.126.146/tcp/21777/p2p/12D3KooWB4XaqGMDxAzaXypBBXBvwqZpgAAPcgkwuzy2LbKMRxpH",
		"/ip4/157.90.111.140/tcp/21777/p2p/12D3KooWKwbdkwfCFdrucHir69U7E9P2BbHctmsBj8Hcms5jP64h",
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
