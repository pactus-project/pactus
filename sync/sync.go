package sync

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/sync/stats"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
)

type Synchronizer struct {
	ctx            context.Context
	config         *Config
	store          *store.Store
	state          *state.State
	consensus      *consensus.Consensus
	txPool         *txpool.TxPool
	stats          *stats.Stats
	selfID         peer.ID
	selfAddress    crypto.Address
	blockPool      map[int]*block.Block
	txkPool        map[crypto.Hash]*tx.Tx
	broadcastCh    <-chan message.Message
	txTopic        *pubsub.Topic
	stateTopic     *pubsub.Topic
	consensusTopic *pubsub.Topic
	txSub          *pubsub.Subscription
	stateSub       *pubsub.Subscription
	consensusSub   *pubsub.Subscription
	logger         *logger.Logger
}

func NewSynchronizer(
	conf *Config,
	addr crypto.Address,
	state *state.State,
	store *store.Store,
	consensus *consensus.Consensus,
	txpool *txpool.TxPool,
	net *network.Network,
	broadcastCh <-chan message.Message) (*Synchronizer, error) {
	syncer := &Synchronizer{
		ctx:         context.Background(),
		config:      conf,
		state:       state,
		store:       store,
		consensus:   consensus,
		txPool:      txpool,
		selfAddress: addr,
		blockPool:   make(map[int]*block.Block),
		txkPool:     make(map[crypto.Hash]*tx.Tx),
		broadcastCh: broadcastCh,
	}
	txTopic, err := net.JoinTopic("tx")
	if err != nil {
		return nil, err
	}
	txSub, err := txTopic.Subscribe()
	if err != nil {
		return nil, err
	}
	stateTopic, err := net.JoinTopic("state")
	if err != nil {
		return nil, err
	}
	stateSub, err := stateTopic.Subscribe()
	if err != nil {
		return nil, err
	}
	consensusTopic, err := net.JoinTopic("consensus")
	if err != nil {
		return nil, err
	}
	consensusSub, err := consensusTopic.Subscribe()
	if err != nil {
		return nil, err
	}

	logger := logger.NewLogger("_syncer", syncer)
	syncer.selfID = net.ID()
	syncer.txTopic = txTopic
	syncer.txSub = txSub
	syncer.stateTopic = stateTopic
	syncer.stateSub = stateSub
	syncer.consensusTopic = consensusTopic
	syncer.consensusSub = consensusSub
	syncer.logger = logger
	syncer.stats = stats.NewStats(logger)

	return syncer, nil
}

func (syncer *Synchronizer) Start() error {
	go syncer.txLoop()
	go syncer.stateLoop()
	go syncer.consensusLoop()
	go syncer.broadcastLoop()

	return nil
}

func (syncer *Synchronizer) Stop() error {
	syncer.ctx.Done()
	syncer.txTopic.Close()
	syncer.txSub.Cancel()
	syncer.stateTopic.Close()
	syncer.stateSub.Cancel()
	syncer.consensusTopic.Close()
	syncer.consensusSub.Cancel()
	return nil
}

func (syncer *Synchronizer) txLoop() {
	for {
		m, err := syncer.txSub.Next(syncer.ctx)
		if err != nil {
			syncer.logger.Error("readLoop error", "err", err)
			return
		}

		syncer.parsMessage(m)
	}
}

func (syncer *Synchronizer) stateLoop() {
	for {
		m, err := syncer.stateSub.Next(syncer.ctx)
		if err != nil {
			syncer.logger.Error("readLoop error", "err", err)
			return
		}

		syncer.parsMessage(m)
	}
}

func (syncer *Synchronizer) consensusLoop() {
	for {
		m, err := syncer.consensusSub.Next(syncer.ctx)
		if err != nil {
			syncer.logger.Error("readLoop error", "err", err)
			return
		}

		syncer.parsMessage(m)
	}
}

func (syncer *Synchronizer) Fingerprint() string {
	return fmt.Sprintf("{peers: %d}", syncer.stats.PeersCount())
}
