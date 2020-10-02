package node

import (
	"time"

	"github.com/zarbchain/zarb-go/network"

	"github.com/pkg/errors"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/utils"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/www/capnp"
	"github.com/zarbchain/zarb-go/www/http"
)

type Node struct {
	genesisDoc *genesis.Genesis
	config     *config.Config
	state      *state.State
	store      *store.Store
	txPool     *txpool.TxPool
	consensus  *consensus.Consensus
	network    *network.Network

	capnp *capnp.Server
	http  *http.Server
}

func NewNode(genDoc *genesis.Genesis, conf *config.Config, privValidator *validator.PrivValidator) (*Node, error) {

	// Init logger
	logger.InitLogger(conf)

	network, err := network.NewNetwork(conf)
	if err != nil {
		return nil, err
	}

	store, err := store.NewStore(conf)
	if err != nil {
		return nil, err
	}

	txPool, err := txpool.NewTxPool(conf, network)
	if err != nil {
		return nil, err
	}

	state, err := state.LoadStateOrNewState(conf, genDoc, network, store, txPool)
	if err != nil {
		return nil, err
	}

	consensus, err := consensus.NewConsensus(conf, state, network, store, privValidator)
	if err != nil {
		return nil, err
	}

	capnp, err := capnp.NewServer(store, conf)
	if err != nil {
		return nil, errors.Wrap(err, "could not create Capnproto server")
	}

	http, err := http.NewServer(conf)
	if err != nil {
		return nil, errors.Wrap(err, "could not create http server")
	}

	node := &Node{
		config:     conf,
		genesisDoc: genDoc,
		network:    network,
		state:      state,
		store:      store,
		txPool:     txPool,
		consensus:  consensus,
		capnp:      capnp,
		http:       http,
	}

	return node, nil
}

func (n *Node) Start() error {
	now := utils.Now()
	genTime := n.genesisDoc.GenesisTime()
	if genTime.After(now) {
		logger.Info("Genesis time is in the future. Sleeping until then...", "genTime", genTime)
		time.Sleep(genTime.Sub(now))
	}

	n.consensus.Start()
	n.txPool.Start()
	n.network.Start()

	// Wait for network to started
	time.Sleep(1 * time.Second)

	err := n.capnp.StartServer()
	if err != nil {
		return errors.Wrap(err, "could not start Capnproto server")
	}

	err = n.http.StartServer()
	if err != nil {
		return errors.Wrap(err, "could not start http server")
	}

	return nil
}

func (n *Node) Stop() {
	logger.Info("Stopping Node")

	n.consensus.Stop()
	n.network.Stop()
}
