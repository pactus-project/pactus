package node

import (
	"time"

	"gitlab.com/zarb-chain/zarb-go/network"

	"github.com/pkg/errors"
	"gitlab.com/zarb-chain/zarb-go/config"
	"gitlab.com/zarb-chain/zarb-go/consensus"
	"gitlab.com/zarb-chain/zarb-go/genesis"
	"gitlab.com/zarb-chain/zarb-go/logger"
	"gitlab.com/zarb-chain/zarb-go/state"
	"gitlab.com/zarb-chain/zarb-go/store"
	"gitlab.com/zarb-chain/zarb-go/txpool"
	"gitlab.com/zarb-chain/zarb-go/utils"
	"gitlab.com/zarb-chain/zarb-go/validator"
	"gitlab.com/zarb-chain/zarb-go/www/capnp"
	"gitlab.com/zarb-chain/zarb-go/www/http"
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
