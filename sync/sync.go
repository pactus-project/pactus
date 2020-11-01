package sync

import (
	"context"
	"fmt"
	"time"

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
	state           state.State
	consensus       *consensus.Consensus
	txPool          *txpool.TxPool
	stats           *stats.Stats
	blockPool       *BlockPool
	txkPool         map[crypto.Hash]*tx.Tx
	broadcastCh     <-chan message.Message
	networkApi      NetworkApi
	heartBeatTicker *time.Ticker
	logger          *logger.Logger
}

func NewSynchronizer(
	conf *Config,
	addr crypto.Address,
	state state.State,
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
		txkPool:     make(map[crypto.Hash]*tx.Tx),
		broadcastCh: broadcastCh,
	}

	logger := logger.NewLogger("_sync", syncer)

	api, err := newNetworkApi(syncer.ctx, addr, net, syncer.ParsMessage, logger)
	if err != nil {
		return nil, err
	}

	syncer.logger = logger
	syncer.blockPool = NewBlockPool(logger)
	syncer.stats = stats.NewStats(logger)
	syncer.networkApi = api

	return syncer, nil
}

func (syncer *Synchronizer) Start() error {
	syncer.networkApi.Start()

	go syncer.broadcastLoop()

	syncer.heartBeatTicker = time.NewTicker(syncer.config.HeartBeatTimeout)
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
	syncer.ctx.Done()
	syncer.networkApi.Stop()
	syncer.heartBeatTicker.Stop()

	return nil
}

func (syncer *Synchronizer) maybeSyncing() {
	lastHeight := syncer.state.LastBlockHeight()
	networkHeight := syncer.stats.MaxHeight()

	if lastHeight >= networkHeight-1 {
		syncer.consensus.ScheduleNewHeight()
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

func (syncer *Synchronizer) broadcastLoop() {
	for {
		select {
		case <-syncer.ctx.Done():
			return

		case msg := <-syncer.broadcastCh:
			syncer.networkApi.PublishMessage(msg)
		}
	}
}
func (syncer *Synchronizer) Fingerprint() string {
	return fmt.Sprintf("{☍ %d ⛲ %d }", syncer.stats.PeersCount(), syncer.blockPool.BlockLen())
}
