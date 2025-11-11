package node

import (
	"context"
	"time"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/consensus/manager"
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
	"github.com/pactus-project/pactus/util/pipeline"
	"github.com/pactus-project/pactus/version"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/pactus-project/pactus/www/html"
	"github.com/pactus-project/pactus/www/http"
	"github.com/pactus-project/pactus/www/jsonrpc"
	"github.com/pactus-project/pactus/www/zmq"
	"github.com/pkg/errors"
)

type Node struct {
	ctx           context.Context
	cancel        context.CancelFunc
	genesisDoc    *genesis.Genesis
	config        *config.Config
	state         state.Facade
	store         store.Store
	txPool        txpool.TxPool
	consV1Mgr     manager.Manager // Deprecated:: replaced by new consensus algorithm
	consV2Mgr     manager.Manager
	network       network.Network
	sync          sync.Synchronizer
	grpc          *grpc.Server
	html          *html.Server
	http          *http.Server
	jsonrpc       *jsonrpc.Server
	zeromq        *zmq.Server
	broadcastPipe pipeline.Pipeline[message.Message]
	networkPipe   pipeline.Pipeline[network.Event]
	eventPipe     pipeline.Pipeline[any]
}

func NewNode(genDoc *genesis.Genesis, conf *config.Config,
	valKeys []*bls.ValidatorKey, rewardAddrs []crypto.Address,
) (*Node, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize the logger
	logger.InitGlobalLogger(conf.Logger)

	chainType := genDoc.ChainType()

	broadcastPipe := pipeline.New[message.Message](ctx, "Broadcast Pipeline", 100)
	networkPipe := pipeline.New[network.Event](ctx, "Network Pipeline", 500)
	eventPipe := pipeline.New[any](ctx, "Event Pipeline", 100)

	store, err := store.NewStore(conf.Store)
	if err != nil {
		cancel()

		return nil, err
	}

	logger.Info("You are running a Pactus blockchain",
		"version", version.NodeVersion().StringWithAlias(),
		"network", chainType, "pruned", store.IsPruned())

	txPool := txpool.NewTxPool(conf.TxPool, store, broadcastPipe)

	state, err := state.LoadOrNewState(genDoc, valKeys, store, txPool, eventPipe)
	if err != nil {
		cancel()

		return nil, err
	}

	net, err := network.NewNetwork(ctx, conf.Network, networkPipe)
	if err != nil {
		cancel()

		return nil, err
	}

	consV1Mgr := manager.NewManagerV1(conf.Consensus, state, valKeys, rewardAddrs, broadcastPipe)
	consV2Mgr := manager.NewManagerV2(conf.ConsensusV2, state, valKeys, rewardAddrs, broadcastPipe)
	walletMgr := wallet.NewWalletManager(conf.WalletManager)

	if !store.IsPruned() {
		conf.Sync.Services.Append(service.FullNode)
	}
	sync, err := sync.NewSynchronizer(conf.Sync, valKeys, state, consV1Mgr, consV2Mgr, net, broadcastPipe, networkPipe)
	if err != nil {
		cancel()

		return nil, err
	}

	enableHTTPAuth := (conf.GRPC.BasicAuth != "")

	zeromqServer, err := zmq.New(ctx, conf.ZeroMq, eventPipe)
	if err != nil {
		cancel()

		return nil, err
	}
	curConsMgr := consV1Mgr
	if consV1Mgr.IsDeprecated() {
		curConsMgr = consV2Mgr
	}

	grpcServer := grpc.NewServer(ctx, conf.GRPC, state, sync, net, curConsMgr, walletMgr, zeromqServer.Publishers())
	htmlServer := html.NewServer(ctx, conf.HTML, enableHTTPAuth)
	httpServer := http.NewServer(ctx, conf.HTTP)
	jsonrpcServer := jsonrpc.NewServer(ctx, conf.JSONRPC)

	node := &Node{
		ctx:           ctx,
		cancel:        cancel,
		config:        conf,
		genesisDoc:    genDoc,
		network:       net,
		state:         state,
		txPool:        txPool,
		consV1Mgr:     consV1Mgr,
		consV2Mgr:     consV2Mgr,
		sync:          sync,
		store:         store,
		grpc:          grpcServer,
		html:          htmlServer,
		http:          httpServer,
		jsonrpc:       jsonrpcServer,
		zeromq:        zeromqServer,
		broadcastPipe: broadcastPipe,
		networkPipe:   networkPipe,
		eventPipe:     eventPipe,
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

	curConsMgr := n.consV1Mgr
	if n.consV1Mgr.IsDeprecated() {
		curConsMgr = n.consV2Mgr
	}
	curConsMgr.MoveToNewHeight()

	err := n.grpc.StartServer()
	if err != nil {
		return errors.Wrap(err, "could not start gRPC server")
	}

	err = n.html.StartServer(n.grpc.Address())
	if err != nil {
		return errors.Wrap(err, "could not start HTML server")
	}

	err = n.http.StartServer(n.grpc.Address())
	if err != nil {
		return errors.Wrap(err, "could not start HTTP-API server")
	}

	err = n.jsonrpc.StartServer(n.grpc.Address())
	if err != nil {
		return errors.Wrap(err, "could not start JSON-RPC server")
	}

	return nil
}

func (n *Node) Stop() {
	logger.Info("stopping Node")
	n.cancel()
	n.broadcastPipe.Close()
	n.networkPipe.Close()
	n.eventPipe.Close()

	// Wait for one second
	time.Sleep(1 * time.Second)

	n.network.Stop()

	// Wait for network to stop
	time.Sleep(1 * time.Second)

	n.sync.Stop()
	n.state.Close()
	n.store.Close()
	n.grpc.StopServer()
	n.html.StopServer()
	n.http.StopServer()
	n.jsonrpc.StopServer()
	n.zeromq.Close()
}

// these methods are using by GUI.

func (n *Node) ConsManager() manager.ManagerReader {
	return n.consV1Mgr
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
