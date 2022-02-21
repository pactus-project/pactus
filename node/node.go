package node

import (
	"time"

	"github.com/pkg/errors"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/www/capnp"
	"github.com/zarbchain/zarb-go/www/grpc"
	"github.com/zarbchain/zarb-go/www/http"
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
}

func NewNode(genDoc *genesis.Genesis, conf *config.Config, signer crypto.Signer) (*Node, error) {
	// Init logger
	logger.InitLogger(conf.Logger)

	network, err := network.NewNetwork(conf.Network)
	if err != nil {
		return nil, err
	}
	broadcastCh := make(chan payload.Payload, 100)

	txPool, err := txpool.NewTxPool(conf.TxPool, broadcastCh)
	if err != nil {
		return nil, err
	}

	store, err := store.NewStore(conf.Store)
	if err != nil {
		return nil, err
	}

	state, err := state.LoadOrNewState(conf.State, genDoc, signer, store, txPool)
	if err != nil {
		return nil, err
	}

	consensus, err := consensus.NewConsensus(conf.Consensus, state, signer, broadcastCh)
	if err != nil {
		return nil, err
	}

	sync, err := sync.NewSynchronizer(conf.Sync, signer, state, consensus, network, broadcastCh)
	if err != nil {
		return nil, err
	}

	capnp, err := capnp.NewServer(conf.Capnp, state, sync)
	if err != nil {
		return nil, errors.Wrap(err, "could not create Capnproto server")
	}

	http, err := http.NewServer(conf.HTTP)
	if err != nil {
		return nil, errors.Wrap(err, "could not create http server")
	}

	grpc, err := grpc.NewServer(conf.GRPC, state, sync)
	if err != nil {
		return nil, errors.Wrap(err, "could not create grpc server")
	}

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
	}

	return node, nil
}

func (n *Node) Start() error {
	now := util.Now()
	genTime := n.genesisDoc.GenesisTime()
	if genTime.After(now) {
		logger.Info("💤 Genesis time is in the future. Sleeping until then...", "duration", genTime.Sub(now), "genesis_hash", n.genesisDoc.Hash())
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
