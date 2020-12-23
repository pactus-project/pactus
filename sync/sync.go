package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/sync/stats"
	"github.com/zarbchain/zarb-go/txpool"
)

type Synchronizer struct {
	// Not: Synchronizer should not have any lock to prevent dead lock situation.
	// Other modules like state or consesnus are thread safe

	ctx             context.Context
	config          *Config
	signer          crypto.Signer
	state           state.State
	txPool          txpool.TxPool
	consensus       consensus.Consensus
	stats           *stats.Stats
	cache           *cache.Cache
	broadcastCh     <-chan *message.Message
	networkAPI      NetworkAPI
	heartBeatTicker *time.Ticker
	logger          *logger.Logger
}

func NewSynchronizer(
	conf *Config,
	signer crypto.Signer,
	state state.State,
	consensus consensus.Consensus,
	txPool txpool.TxPool,
	net *network.Network,
	broadcastCh <-chan *message.Message) (*Synchronizer, error) {
	syncer := &Synchronizer{
		ctx:         context.Background(),
		config:      conf,
		signer:      signer,
		state:       state,
		consensus:   consensus,
		txPool:      txPool,
		broadcastCh: broadcastCh,
	}

	logger := logger.NewLogger("_sync", syncer)

	api, err := newNetworkAPI(syncer.ctx, net, syncer.ParsMessage)
	if err != nil {
		return nil, err
	}

	cache, err := cache.NewCache(conf.CacheSize, state.StoreReader())
	if err != nil {
		return nil, err
	}

	syncer.logger = logger
	syncer.cache = cache
	syncer.stats = stats.NewStats(state.GenesisHash())
	syncer.networkAPI = api

	return syncer, nil
}

func (syncer *Synchronizer) Start() error {
	if err := syncer.networkAPI.Start(); err != nil {
		return err
	}

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

func (syncer *Synchronizer) Stop() {
	syncer.ctx.Done()
	syncer.networkAPI.Stop()
	syncer.heartBeatTicker.Stop()
}

func (syncer *Synchronizer) maybeSynced() {
	lastHeight := syncer.state.LastBlockHeight()
	networkHeight := syncer.stats.MaxClaimedHeight()

	if lastHeight >= networkHeight {
		syncer.logger.Info("We are synced", "height", lastHeight)
		syncer.informConsensusToMoveToNewHeight()
	}
}

func (syncer *Synchronizer) informConsensusToMoveToNewHeight() {
	syncer.consensus.MoveToNewHeight()
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

			switch msg.PayloadType() {
			case payload.PayloadTypeTxsReq:
				pld := msg.Payload.(*payload.TxsReqPayload)
				for i, id := range pld.IDs {
					trx := syncer.cache.GetTransaction(id)
					if trx != nil {
						if err := syncer.txPool.AppendTx(trx); err == nil {
							pld.IDs = append(pld.IDs[:i], pld.IDs[i+1:]...)
						}
					}

				}

				if len(pld.IDs) > 0 {
					syncer.publishMessage(msg)
				}

			default:
				syncer.publishMessage(msg)

			}
		}
	}
}
func (syncer *Synchronizer) Fingerprint() string {
	return fmt.Sprintf("{☍ %d ⛲ %d ↥ %d}",
		syncer.stats.PeersCount(),
		syncer.cache.Len(),
		syncer.stats.MaxClaimedHeight())
}
