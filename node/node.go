package node

import (
	"time"

	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/txpool"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/version"
	"github.com/pactus-project/pactus/www/grpc"
	"github.com/pactus-project/pactus/www/http"
	"github.com/pactus-project/pactus/www/nanomsg"
	"github.com/pactus-project/pactus/www/nanomsg/event"
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
	nanomsg    *nanomsg.Server
}

func NewNode(genDoc *genesis.Genesis, conf *config.Config,
	signers []crypto.Signer, rewardAddrs []crypto.Address,
) (*Node, error) {
	// Initialize the logger
	logger.InitGlobalLogger(conf.Logger)

	logger.Info("You are running a Pactus blockchain",
		"version", version.Version(),
		"network", genDoc.ChainType())

	networkName := genDoc.NetworkName()
	network, err := network.NewNetwork(networkName, conf.Network)
	if err != nil {
		return nil, err
	}
	messageCh := make(chan message.Message, 500)
	eventCh := make(chan event.Event, 500)
	if !conf.Nanomsg.Enable {
		eventCh = nil
	}

	txPool := txpool.NewTxPool(conf.TxPool, messageCh)

	store, err := store.NewStore(conf.Store, int(genDoc.Params().TransactionToLiveInterval))
	if err != nil {
		return nil, err
	}

	state, err := state.LoadOrNewState(genDoc, signers, store, txPool, eventCh)
	if err != nil {
		return nil, err
	}

	consMgr := consensus.NewManager(conf.Consensus, state, signers, rewardAddrs, messageCh)

	sync, err := sync.NewSynchronizer(conf.Sync, signers, state, consMgr, network, messageCh)
	if err != nil {
		return nil, err
	}

	http := http.NewServer(conf.HTTP)
	grpc := grpc.NewServer(conf.GRPC, state, sync, consMgr)
	nanomsg := nanomsg.NewServer(conf.Nanomsg, eventCh)

	node := &Node{
		config:     conf,
		genesisDoc: genDoc,
		network:    network,
		state:      state,
		txPool:     txPool,
		consMgr:    consMgr,
		sync:       sync,
		store:      store,
		http:       http,
		grpc:       grpc,
		nanomsg:    nanomsg,
	}

	return node, nil
}

func (n *Node) Start() error {
	now := util.Now()
	genTime := n.genesisDoc.GenesisTime()
	if genTime.After(now) {
		logger.Info("💤 Genesis time is in the future. Sleeping until then...",
			"duration", genTime.Sub(now), "genesis_hash", n.genesisDoc.Hash())
		time.Sleep(genTime.Sub(now) - 1*time.Second)
	}

	if err := n.network.Start(); err != nil {
		return err
	}
	// Wait for network to started
	time.Sleep(1 * time.Second)

	if err := n.consMgr.Start(); err != nil {
		return err
	}

	if err := n.sync.Start(); err != nil {
		return err
	}

	err := n.grpc.StartServer()
	if err != nil {
		return errors.Wrap(err, "could not start grpc server")
	}

	err = n.http.StartServer(n.grpc.Address())
	if err != nil {
		return errors.Wrap(err, "could not start http server")
	}

	err = n.nanomsg.StartServer()
	if err != nil {
		return errors.Wrap(err, "could not start nanomsg server")
	}

	return nil
}

func (n *Node) Stop() {
	logger.Info("stopping Node")

	n.consMgr.Stop()
	n.network.Stop()
	n.sync.Stop()
	n.state.Close()
	n.store.Close()
	n.http.StopServer()
	n.grpc.StopServer()
	n.nanomsg.StopServer()
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
