package node

import (
	"time"

	"github.com/pkg/errors"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/node/config"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/genesis"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/util/logger"
	"github.com/zarbchain/zarb-go/www/capnp"
	"github.com/zarbchain/zarb-go/www/grpc"
	"github.com/zarbchain/zarb-go/www/http"
	"github.com/zarbchain/zarb-go/www/zmq"
)

type Node struct {
	genesisDoc *genesis.Genesis
	config     *config.Config
	state      state.Facade
	store      store.Store
	txPool     txpool.TxPool
	consensus  consensus.Consensus
	network    network.Network
	sync       sync.Synchronizer
	capnp      *capnp.Server
	http       *http.Server
	grpc       *grpc.Server
	zmq        *zmq.Server
}

func NewNode(genDoc *genesis.Genesis, conf *config.Config, signer crypto.Signer) (*Node, error) {
	// Initialize the logger
	logger.InitLogger(conf.Logger)

	network, err := network.NewNetwork(conf.Network)
	if err != nil {
		return nil, err
	}
	broadcastCh := make(chan message.Message, 100)

	txPool := txpool.NewTxPool(conf.TxPool, broadcastCh)

	store, err := store.NewStore(conf.Store, int(genDoc.Params().TransactionToLiveInterval))
	if err != nil {
		return nil, err
	}

	state, err := state.LoadOrNewState(conf.State, genDoc, signer, store, txPool)
	if err != nil {
		return nil, err
	}

	consensus := consensus.NewConsensus(conf.Consensus, state, signer, broadcastCh)

	sync, err := sync.NewSynchronizer(conf.Sync, signer, state, consensus, network, broadcastCh)
	if err != nil {
		return nil, err
	}

	capnp := capnp.NewServer(conf.Capnp, state, sync, consensus)
	http := http.NewServer(conf.HTTP)
	grpc := grpc.NewServer(conf.GRPC, state, sync)
	zmq := zmq.NewServer(conf.Zmq, state, sync)

	node := &Node{
		config:     conf,
		genesisDoc: genDoc,
		network:    network,
		state:      state,
		txPool:     txPool,
		consensus:  consensus,
		sync:       sync,
		store:      store,
		capnp:      capnp,
		http:       http,
		grpc:       grpc,
		zmq:        zmq,
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

	if err := n.consensus.Start(); err != nil {
		return err
	}

	if err := n.sync.Start(); err != nil {
		return err
	}

	err := n.capnp.StartServer()
	if err != nil {
		return errors.Wrap(err, "could not start Capnproto server")
	}

	err = n.http.StartServer(n.capnp.Address())
	if err != nil {
		return errors.Wrap(err, "could not start http server")
	}

	err = n.grpc.StartServer()
	if err != nil {
		return errors.Wrap(err, "could not start grpc server")
	}
	err = n.zmq.StartServer()
	if err != nil {
		return errors.Wrap(err, "could not start grpc server")
	}

	return nil
}

func (n *Node) Stop() {
	logger.Info("stopping Node")

	n.consensus.Stop()
	n.network.Stop()
	n.sync.Stop()
	n.state.Close()
	n.store.Close()
	n.http.StopServer()
	n.capnp.StopServer()
	n.grpc.StopServer()
}

func (n *Node) Consensus() consensus.Reader {
	return n.consensus
}
func (n *Node) Sync() sync.Synchronizer {
	return n.sync
}
func (n *Node) State() state.Facade {
	return n.state
}
