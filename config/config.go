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
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/pactus-project/pactus/www/http"
	"github.com/pactus-project/pactus/www/jsonrpc"
	"github.com/pactus-project/pactus/www/zmq"
	"github.com/pelletier/go-toml/v2"
)

var (
	//go:embed example_config.toml
	exampleConfigBytes []byte

	//go:embed bootstrap.json
	bootstrapInfoBytes []byte

	//go:embed banned_addrs.json
	bannedAddrBytes []byte
)

type Config struct {
	Node          *NodeConfig       `toml:"node"`
	Store         *store.Config     `toml:"store"`
	Network       *network.Config   `toml:"network"`
	Sync          *sync.Config      `toml:"sync"`
	TxPool        *txpool.Config    `toml:"tx_pool"`
	Consensus     *consensus.Config `toml:"-"`
	Logger        *logger.Config    `toml:"logger"`
	GRPC          *grpc.Config      `toml:"grpc"`
	JSONRPC       *jsonrpc.Config   `toml:"jsonrpc"`
	HTTP          *http.Config      `toml:"http"`
	WalletManager *wallet.Config    `toml:"-"`
	ZeroMq        *zmq.Config       `toml:"zeromq"`
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
		Node:          DefaultNodeConfig(),
		Store:         store.DefaultConfig(),
		Network:       network.DefaultConfig(),
		Sync:          sync.DefaultConfig(),
		TxPool:        txpool.DefaultConfig(),
		Consensus:     consensus.DefaultConfig(),
		Logger:        logger.DefaultConfig(),
		GRPC:          grpc.DefaultConfig(),
		JSONRPC:       jsonrpc.DefaultConfig(),
		HTTP:          http.DefaultConfig(),
		ZeroMq:        zmq.DefaultConfig(),
		WalletManager: wallet.DefaultConfig(),
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

	bannedList := make([]string, 0)
	if err := json.Unmarshal(bannedAddrBytes, &bannedList); err != nil {
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
	conf.HTTP.EnablePprof = false

	return conf
}

func DefaultConfigTestnet() *Config {
	conf := defaultConfig()
	conf.Network.DefaultBootstrapAddrStrings = []string{
		"/dns/testnet1.pactus.org/tcp/21777/p2p/12D3KooWR7ZB3nGih1Fz7Yg83Zap8Cpxr73T6PPihBsEpTG5BZyk",
		"/dns/testnet2.pactus.org/tcp/21777/p2p/12D3KooWQcDuFDMGsw6gG7oNFw7C4x7ozoMu69J7WEAojKCaNzji",
		"/dns/testnet3.pactus.org/tcp/21777/p2p/12D3KooWLsAPSJ4xowd9thGbPmbweBT6sg3nEiPjDJccaWZacsUR",
		"/dns/testnet4.pactus.org/tcp/21777/p2p/12D3KooWJKYdHzWZGibnj74NSSgKRu4Ez6MijDWMfLfXxeL4un6v",
		"/ip4/65.108.211.187/tcp/21777/p2p/12D3KooWB42BLfzxSF5SMhSTSEyfJ6yhSM8togLfExrRWFMJeb5u",
		"/ip4/103.27.206.208/tcp/21777/p2p/12D3KooWMTDwDTBMaf2Sem5tWRe1dB6PFY8LeqkZ2e5drrbbPTDn", // andrut.pactus.testnet
		"/ip4/65.108.142.81/tcp/21777/p2p/12D3KooWAdRga2NCbaPfVgSEzAAZW2psfJmPi3PFJzF81qbccJsR",  // CherryValidator
		"/ip4/95.217.89.202/tcp/21777/p2p/12D3KooWH3S9gMYybr1pd4K5o3CBLbZLQ1REKBsPWt6NWPi4bgPn",  // CodeBlockLab
		"/ip4/65.109.234.125/tcp/21777/p2p/12D3KooWPWE8QwZxd32ui9DL115vmdcFo3cudGfVgHTX9Bo4HFEB", // Sensifai
		"/ip4/77.37.122.54/tcp/22194/p2p/12D3KooWKwbdkwfCFdrucHir69U7E9P2BbHctmsBj8Hcms5jP64h",   // Javad
		"/ip4/188.213.198.83/tcp/21777/p2p/12D3KooWJ5kSyD3VQb1ewhgRcmAPPg2zus1rYPbgnBMMGBtC9pr5",
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
	conf.HTTP.EnablePprof = false

	return conf
}

func DefaultConfigLocalnet() *Config {
	conf := defaultConfig()
	conf.Network.EnableRelay = false
	conf.Network.EnableNATService = false
	conf.Network.EnableUPnP = false
	conf.Network.BootstrapAddrStrings = []string{}
	conf.Network.MaxConns = 16
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
	conf.HTTP.EnablePprof = true
	conf.ZeroMq.ZmqPubBlockInfo = "tcp://127.0.0.1:28332"
	conf.ZeroMq.ZmqPubTxInfo = "tcp://127.0.0.1:28333"
	conf.ZeroMq.ZmqPubRawBlock = "tcp://127.0.0.1:28334"
	conf.ZeroMq.ZmqPubRawTx = "tcp://127.0.0.1:28335"
	conf.ZeroMq.ZmqPubHWM = 1000

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
	encoder.SetIndentTables(true)
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
	if strict {
		decoder.DisallowUnknownFields()
	}
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
	if err := conf.JSONRPC.BasicCheck(); err != nil {
		return err
	}
	if err := conf.GRPC.BasicCheck(); err != nil {
		return err
	}
	if err := conf.ZeroMq.BasicCheck(); err != nil {
		return err
	}

	return conf.HTTP.BasicCheck()
}
