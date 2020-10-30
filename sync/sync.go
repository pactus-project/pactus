package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/stats"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
)

type Synchronizer struct {
	ctx             context.Context
	config          *Config
	store           state.StoreReader
	state           *state.State
	consensus       *consensus.Consensus
	txPool          *txpool.TxPool
	stats           *stats.Stats
	selfID          peer.ID
	selfAddress     crypto.Address
	blockPool       *BlockPool
	txkPool         map[crypto.Hash]*tx.Tx
	broadcastCh     <-chan message.Message
	generalTopic    *pubsub.Topic
	txTopic         *pubsub.Topic
	blockTopic      *pubsub.Topic
	consensusTopic  *pubsub.Topic
	generalSub      *pubsub.Subscription
	txSub           *pubsub.Subscription
	blockSub        *pubsub.Subscription
	consensusSub    *pubsub.Subscription
	heartBeatTicker *time.Ticker
	logger          *logger.Logger
}

func NewSynchronizer(
	conf *Config,
	addr crypto.Address,
	state *state.State,
	consensus *consensus.Consensus,
	txpool *txpool.TxPool,
	net *network.Network,
	broadcastCh <-chan message.Message) (*Synchronizer, error) {
	syncer := &Synchronizer{
		ctx:         context.Background(),
		config:      conf,
		store:       state.StoreReader(),
		state:       state,
		consensus:   consensus,
		txPool:      txpool,
		selfAddress: addr,
		txkPool:     make(map[crypto.Hash]*tx.Tx),
		broadcastCh: broadcastCh,
	}
	generalTopic, err := net.JoinTopic("general")
	if err != nil {
		return nil, err
	}
	generalSub, err := generalTopic.Subscribe()
	if err != nil {
		return nil, err
	}
	txTopic, err := net.JoinTopic("tx")
	if err != nil {
		return nil, err
	}
	txSub, err := txTopic.Subscribe()
	if err != nil {
		return nil, err
	}
	blockTopic, err := net.JoinTopic("block")
	if err != nil {
		return nil, err
	}
	blockSub, err := blockTopic.Subscribe()
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

	logger := logger.NewLogger("_sync", syncer)
	syncer.selfID = net.ID()
	syncer.txTopic = txTopic
	syncer.txSub = txSub
	syncer.blockTopic = blockTopic
	syncer.blockSub = blockSub
	syncer.blockTopic = blockTopic
	syncer.generalTopic = generalTopic
	syncer.generalSub = generalSub
	syncer.consensusTopic = consensusTopic
	syncer.consensusSub = consensusSub
	syncer.logger = logger
	syncer.blockPool = NewBlockPool(logger)
	syncer.stats = stats.NewStats(logger)

	return syncer, nil
}

func (syncer *Synchronizer) Start() error {
	syncer.heartBeatTicker = time.NewTicker(syncer.config.HeartBeatTimeout)

	go syncer.txLoop()
	go syncer.blockLoop()
	go syncer.generalLoop()
	go syncer.consensusLoop()
	go syncer.broadcastLoop()
	go syncer.heartBeatTickerLoop()

	syncer.broadcastSalam()

	timer := time.NewTimer(syncer.config.StartingTimeout)
	go func() {
		<-timer.C
		syncer.maybeSyncing()
	}()

	return nil
}

func (syncer *Synchronizer) Stop() error {
	syncer.heartBeatTicker.Stop()

	syncer.ctx.Done()
	syncer.txTopic.Close()
	syncer.txSub.Cancel()
	syncer.blockTopic.Close()
	syncer.blockSub.Cancel()
	syncer.generalTopic.Close()
	syncer.generalSub.Cancel()
	syncer.consensusTopic.Close()
	syncer.consensusSub.Cancel()
	return nil
}

func (syncer *Synchronizer) maybeSyncing() {
	lastHeight := syncer.state.LastBlockHeight()
	networkHeight := syncer.stats.MaxHeight()

	if lastHeight >= networkHeight-1 {
		syncer.consensus.ScheduleNewHeight()
	}
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

func (syncer *Synchronizer) blockLoop() {
	for {
		m, err := syncer.blockSub.Next(syncer.ctx)
		if err != nil {
			syncer.logger.Error("readLoop error", "err", err)
			return
		}

		syncer.parsMessage(m)
	}
}

func (syncer *Synchronizer) generalLoop() {
	for {
		m, err := syncer.generalSub.Next(syncer.ctx)
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

func (syncer *Synchronizer) heartBeatTickerLoop() {
	for {
		select {
		case <-syncer.ctx.Done():
			return
		case <-syncer.heartBeatTicker.C:
			syncer.broadcastHeartBeat()
		}
	}
}

func (syncer *Synchronizer) Fingerprint() string {
	return fmt.Sprintf("{☍ %d ⛲ %d}", syncer.stats.PeersCount(), syncer.blockPool.Size())
}
