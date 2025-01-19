package node

import (
	"context"
	"time"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/version"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/pactus-project/pactus/www/http"
	"github.com/pactus-project/pactus/www/jsonrpc"
	"github.com/pactus-project/pactus/www/zmq"
	"github.com/pkg/errors"
)

type Node struct {
	genesisDoc *genesis.Genesis
	config     *config.Config
	state      state.Facade
	store      store.Store
	txPool     txpool.TxPool
	consMgr    consensus.Manager
	network    network.Network
	sync       sync.Synchronizer
	http       *http.Server
	grpc       *grpc.Server
	jsonrpc    *jsonrpc.Server
	zeromq     *zmq.Server
	eventCh    chan any
}

func NewNode(genDoc *genesis.Genesis, conf *config.Config,
	valKeys []*bls.ValidatorKey, rewardAddrs []crypto.Address,
) (*Node, error) {
	// Initialize the logger
	logger.InitGlobalLogger(conf.Logger)

	chainType := genDoc.ChainType()

	logger.Info("You are running a Pactus blockchain",
		"version", version.NodeVersion.StringWithAlias(),
		"network", chainType)

	messageCh := make(chan message.Message, 500)
	eventCh := make(chan any)

	store, err := store.NewStore(conf.Store)
	if err != nil {
		return nil, err
	}

	txPool := txpool.NewTxPool(conf.TxPool, store, messageCh)

	state, err := state.LoadOrNewState(genDoc, valKeys, store, txPool, eventCh)
	if err != nil {
		return nil, err
	}

	net, err := network.NewNetwork(conf.Network)
	if err != nil {
		return nil, err
	}

	consMgr := consensus.NewManager(conf.Consensus, state, valKeys, rewardAddrs, messageCh)
	walletMgr := wallet.NewWalletManager(conf.WalletManager)

	if !store.IsPruned() {
		conf.Sync.Services.Append(service.FullNode)
	}
	syn, err := sync.NewSynchronizer(conf.Sync, valKeys, state, consMgr, net, messageCh)
	if err != nil {
		return nil, err
	}

	enableHTTPAuth := false
	if conf.GRPC.BasicAuth != "" {
		enableHTTPAuth = true
	}

	zeromqServer, err := zmq.New(context.TODO(), conf.ZeroMq, eventCh)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer(conf.GRPC, state, syn, net, consMgr, walletMgr, zeromqServer)
	httpServer := http.NewServer(conf.HTTP, enableHTTPAuth)
	jsonrpcServer := jsonrpc.NewServer(conf.JSONRPC)

	node := &Node{
		config:     conf,
		genesisDoc: genDoc,
		network:    net,
		state:      state,
		txPool:     txPool,
		consMgr:    consMgr,
		sync:       syn,
		store:      store,
		http:       httpServer,
		grpc:       grpcServer,
		jsonrpc:    jsonrpcServer,
		zeromq:     zeromqServer,
		eventCh:    eventCh,
	}

	return node, nil
}

func (n *Node) Start() error {
	now := time.Now()
	genTime := n.genesisDoc.GenesisTime()
	if genTime.After(now) {
		logger.Info("💤 Genesis time is in the future. Sleeping until then...",
			"duration", genTime.Sub(now), "genesis_hash", n.genesisDoc.Hash())
		time.Sleep(genTime.Sub(now) - 1*time.Second)
	}

	if err := n.network.Start(); err != nil {
		return errors.Wrap(err, "could not start Network")
	}
	// Wait for network to start
	time.Sleep(1 * time.Second)

	if err := n.sync.Start(); err != nil {
		return errors.Wrap(err, "could not start Sync")
	}

	if err := n.consMgr.Start(); err != nil {
		return errors.Wrap(err, "could not start Consensus manager")
	}

	err := n.grpc.StartServer()
	if err != nil {
		return errors.Wrap(err, "could not start gRPC server")
	}

	err = n.http.StartServer(n.grpc.Address())
	if err != nil {
		return errors.Wrap(err, "could not start HTTP server")
	}

	err = n.jsonrpc.StartServer(n.grpc.Address())
	if err != nil {
		return errors.Wrap(err, "could not start JSON-RPC server")
	}

	return nil
}

func (n *Node) Stop() {
	logger.Info("stopping Node")

	n.network.Stop()

	// Wait for network to stop
	time.Sleep(1 * time.Second)

	n.consMgr.Stop()
	n.sync.Stop()
	n.state.Close()
	n.store.Close()
	n.grpc.StopServer()
	n.http.StopServer()
	n.jsonrpc.StopServer()
	n.zeromq.Close()

	close(n.eventCh)
}

// these methods are using by GUI.

func (n *Node) ConsManager() consensus.ManagerReader {
	return n.consMgr
}

func (n *Node) Sync() sync.Synchronizer {
	return n.sync
}

func (n *Node) State() state.Facade {
	return n.state
}

func (n *Node) GRPC() *grpc.Server {
	return n.grpc
}

func (n *Node) Network() network.Network {
	return n.network
}
