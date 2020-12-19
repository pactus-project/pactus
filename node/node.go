package node

import (
	"time"

	"github.com/pkg/errors"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/www/capnp"
	"github.com/zarbchain/zarb-go/www/http"
)

type Node struct {
	genesisDoc *genesis.Genesis
	config     *config.Config
	state      state.State
	txPool     txpool.TxPool
	consensus  consensus.Consensus
	network    *network.Network
	sync       *sync.Synchronizer
	capnp      *capnp.Server
	http       *http.Server
}

func NewNode(genDoc *genesis.Genesis, conf *config.Config, signer crypto.Signer) (*Node, error) {

	// Init logger
	logger.InitLogger(conf.Logger)

	network, err := network.NewNetwork(conf.Network)
	if err != nil {
		return nil, err
	}
	broadcastCh := make(chan *message.Message, 100)

	txPool, err := txpool.NewTxPool(conf.TxPool, broadcastCh)
	if err != nil {
		return nil, err
	}

	state, err := state.LoadOrNewState(conf.State, genDoc, signer, txPool)
	if err != nil {
		return nil, err
	}

	consensus, err := consensus.NewConsensus(conf.Consensus, state, signer, broadcastCh)
	if err != nil {
		return nil, err
	}

	sync, err := sync.NewSynchronizer(conf.Sync, signer.Address(), state, consensus, txPool, network, broadcastCh)
	if err != nil {
		return nil, err
	}

	capnp, err := capnp.NewServer(conf.Capnp, state, txPool)
	if err != nil {
		return nil, errors.Wrap(err, "could not create Capnproto server")
	}

	http, err := http.NewServer(conf.Http)
	if err != nil {
		return nil, errors.Wrap(err, "could not create http server")
	}

	node := &Node{
		config:     conf,
		genesisDoc: genDoc,
		network:    network,
		state:      state,
		txPool:     txPool,
		consensus:  consensus,
		sync:       sync,
		capnp:      capnp,
		http:       http,
	}

	return node, nil
}

func (n *Node) Start() error {
	now := util.Now()
	genTime := n.genesisDoc.GenesisTime()
	if genTime.After(now) {
		logger.Info("💤 Genesis time is in the future. Sleeping until then...", "genTime", genTime)
		time.Sleep(genTime.Sub(now) - 1*time.Second)
	}

	n.network.Start()
	// Wait for network to started
	time.Sleep(1 * time.Second)

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

	return nil
}

func (n *Node) Stop() {
	logger.Info("Stopping Node")

	n.network.Stop()
	n.sync.Stop()
	n.state.Close()
	n.http.StopServer()
	n.capnp.StopServer()
}
