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
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/sync/stats"
	"github.com/zarbchain/zarb-go/txpool"
)

type Synchronizer struct {
	// Not: Synchronizer should not have any lock to prevent dead lock situation.
	// Other modules like state or consesnus are thread safe

	ctx             context.Context
	config          *Config
	store           store.StoreReader
	state           state.State
	txPool          txpool.TxPool
	consensus       *consensus.Consensus
	stats           *stats.Stats
	blockPool       *BlockPool
	broadcastCh     <-chan *message.Message
	networkAPI      NetworkAPI
	heartBeatTicker *time.Ticker
	logger          *logger.Logger
}

func NewSynchronizer(
	conf *Config,
	addr crypto.Address,
	state state.State,
	consensus *consensus.Consensus,
	txPool txpool.TxPool,
	net *network.Network,
	broadcastCh <-chan *message.Message) (*Synchronizer, error) {
	syncer := &Synchronizer{
		ctx:         context.Background(),
		config:      conf,
		store:       state.StoreReader(),
		state:       state,
		consensus:   consensus,
		txPool:      txPool,
		broadcastCh: broadcastCh,
	}

	logger := logger.NewLogger("_sync", syncer)

	api, err := newNetworkAPI(syncer.ctx, addr, net, syncer.ParsMessage)
	if err != nil {
		return nil, err
	}

	syncer.logger = logger
	syncer.blockPool = NewBlockPool()
	syncer.stats = stats.NewStats(state.GenesisHash())
	syncer.networkAPI = api

	return syncer, nil
}

func (syncer *Synchronizer) Start() error {
	syncer.networkAPI.Start()

	go syncer.broadcastLoop()

	syncer.heartBeatTicker = time.NewTicker(syncer.config.HeartBeatTimeout)
	go syncer.heartBeatTickerLoop()

	syncer.broadcastSalam()

	timer := time.NewTimer(syncer.config.StartingTimeout)
	go func() {
		<-timer.C
		syncer.maybeSynced()
	}()

	return nil
}

func (syncer *Synchronizer) Stop() error {
	syncer.ctx.Done()
	syncer.networkAPI.Stop()
	syncer.heartBeatTicker.Stop()

	return nil
}

func (syncer *Synchronizer) maybeSynced() {
	lastHeight := syncer.state.LastBlockHeight()
	networkHeight := syncer.stats.MaxHeight()

	if lastHeight >= networkHeight {
		syncer.logger.Info("We are synced", "height", lastHeight)
		syncer.consensus.MoveToNewHeight()
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
			syncer.publishMessage(msg)
		}
	}
}
func (syncer *Synchronizer) Fingerprint() string {
	return fmt.Sprintf("{☍ %d ⛲ %d height %d}",
		syncer.stats.PeersCount(),
		syncer.blockPool.BlockLen(),
		syncer.stats.MaxHeight())
}
