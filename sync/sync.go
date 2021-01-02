package sync

import (
	"context"
	"fmt"
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/sync/firewall"
	"github.com/zarbchain/zarb-go/sync/network_api"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
)

type PublishMessageFn = func(msg *message.Message)

const FlagInitialBlockDownload = 0x1

type Synchronizer struct {
	// Not: Synchronizer should not have any lock to prevent dead lock situation.
	// Other modules like state or consesnus are thread safe

	ctx             context.Context
	config          *Config
	signer          crypto.Signer
	state           state.State
	txPool          txpool.TxPool
	consensus       consensus.Consensus
	peerSet         *peerset.PeerSet
	firewall        *firewall.Firewall
	cache           *cache.Cache
	consensusTopic  *ConsensusTopic
	generalTopic    *GeneralTopic
	dataTopic       *DataTopic
	broadcastCh     <-chan *message.Message
	networkAPI      network_api.NetworkAPI
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
	peerSet := peerset.NewPeerSet()
	firewall := firewall.NewFirewall(peerSet, state)

	api, err := network_api.NewNetworkAPI(syncer.ctx, net, firewall, syncer.ParsMessage)
	if err != nil {
		return nil, err
	}

	cache, err := cache.NewCache(conf.CacheSize, state.StoreReader())
	if err != nil {
		return nil, err
	}

	syncer.logger = logger
	syncer.cache = cache
	syncer.peerSet = peerSet
	syncer.firewall = firewall
	syncer.networkAPI = api

	syncer.generalTopic = NewGeneralTopic(conf, api.SelfID(), signer.PublicKey(), peerSet, state, logger, syncer.PublishMessage)
	syncer.consensusTopic = NewConsensusTopic(conf, consensus, logger, syncer.PublishMessage)
	syncer.dataTopic = NewDataTopic(conf, cache, state, logger, syncer.PublishMessage)

	return syncer, nil
}

func (syncer *Synchronizer) Start() error {
	if err := syncer.networkAPI.Start(); err != nil {
		return err
	}

	go syncer.broadcastLoop()

	syncer.heartBeatTicker = time.NewTicker(syncer.config.HeartBeatTimeout)
	go syncer.heartBeatTickerLoop()

	syncer.generalTopic.BroadcastSalam()

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
	networkHeight := syncer.peerSet.MaxClaimedHeight()

	if lastHeight >= networkHeight-1 {
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
			// Check if we have transaction in the cache
			case payload.PayloadTypeTransactionsRequest:
				pld := msg.Payload.(*payload.TransactionsRequestPayload)
				for i, id := range pld.IDs {
					trx := syncer.cache.GetTransaction(id)
					if trx != nil {
						if err := syncer.txPool.AppendTx(trx); err == nil {
							pld.IDs = append(pld.IDs[:i], pld.IDs[i+1:]...)
						}
					}

				}

				if len(pld.IDs) > 0 {
					syncer.PublishMessage(msg)
				}

			default:
				syncer.PublishMessage(msg)

			}
		}
	}
}
func (syncer *Synchronizer) Fingerprint() string {
	return fmt.Sprintf("{☍ %d ⛃ %d ↥ %d}",
		syncer.peerSet.Len(),
		syncer.cache.Len(),
		syncer.peerSet.MaxClaimedHeight())
}

func (syncer *Synchronizer) sendBlocksRequestIfWeAreBehind() {
	ourHeight := syncer.state.LastBlockHeight()
	claimedHeight := syncer.peerSet.MaxClaimedHeight()
	if claimedHeight > ourHeight {
		syncer.logger.Debug("Ask for more blocks", "our_height", ourHeight, "claimed_height", claimedHeight)
		hash := syncer.state.LastBlockHash()
		syncer.dataTopic.BroadcastLatestBlocksRequest(ourHeight+1, hash)
	}
}

func (syncer *Synchronizer) ParsMessage(msg *message.Message, from peer.ID) {
	syncer.logger.Debug("Received a message", "from", util.FingerprintPeerID(from), "message", msg)

	switch msg.PayloadType() {
	case payload.PayloadTypeSalam:
		pld := msg.Payload.(*payload.SalamPayload)
		syncer.generalTopic.ProcessSalamPayload(pld)

	case payload.PayloadTypeAleyk:
		pld := msg.Payload.(*payload.AleykPayload)
		syncer.generalTopic.ProcessAleykPayload(pld)

	case payload.PayloadTypeHeartBeat:
		pld := msg.Payload.(*payload.HeartBeatPayload)
		syncer.processHeartBeatPayload(pld)

	case payload.PayloadTypeLatestBlocksRequest:
		pld := msg.Payload.(*payload.LatestBlocksRequestPayload)
		syncer.dataTopic.ProcessLatestBlocksRequestPayload(pld)

	case payload.PayloadTypeLatestBlocks:
		pld := msg.Payload.(*payload.LatestBlocksPayload)
		syncer.dataTopic.ProcessLatestBlocksPayload(pld)
		syncer.informConsensusToMoveToNewHeight()

	case payload.PayloadTypeTransactionsRequest:
		pld := msg.Payload.(*payload.TransactionsRequestPayload)
		syncer.dataTopic.ProcessTransactionsRequestPayload(pld)

	case payload.PayloadTypeTransactions:
		pld := msg.Payload.(*payload.TransactionsPayload)
		syncer.dataTopic.ProcessTransactionsPayload(pld)

	case payload.PayloadTypeProposalRequest:
		pld := msg.Payload.(*payload.ProposalRequestPayload)
		syncer.consensusTopic.ProcessProposalRequestPayload(pld)

	case payload.PayloadTypeProposal:
		pld := msg.Payload.(*payload.ProposalPayload)
		syncer.consensusTopic.ProcessProposalPayload(pld)

	case payload.PayloadTypeVote:
		pld := msg.Payload.(*payload.VotePayload)
		syncer.consensusTopic.ProcessVotePayload(pld)

	case payload.PayloadTypeVoteSet:
		pld := msg.Payload.(*payload.VoteSetPayload)
		syncer.consensusTopic.ProcessVoteSetPayload(pld)

	default:
		syncer.logger.Error("Unknown message type", "type", msg.PayloadType())
	}

	syncer.sendBlocksRequestIfWeAreBehind()
}

func (syncer *Synchronizer) broadcastHeartBeat() {
	hrs := syncer.consensus.HRS()

	// Probable we are syncing
	if !hrs.IsValid() {
		return
	}

	msg := message.NewHeartBeatMessage(syncer.state.LastBlockHash(), hrs)
	syncer.PublishMessage(msg)
}

func (syncer *Synchronizer) PublishMessage(msg *message.Message) {

	err := syncer.networkAPI.PublishMessage(msg)

	if err != nil {
		syncer.logger.Error("Error on publishing message", "message", msg, "err", err)
	} else {
		syncer.logger.Debug("Publishing new message", "message", msg)
	}
}

func (syncer *Synchronizer) processHeartBeatPayload(pld *payload.HeartBeatPayload) {
	syncer.logger.Trace("Process heartbeat payload", "pld", pld)

	hrs := syncer.consensus.HRS()
	if pld.Pulse.Height() == hrs.Height() {
		if pld.Pulse.GreaterThan(hrs) {
			syncer.logger.Trace("Our consensus is behind of this peer.")
			// Let's ask for more votes
			syncer.consensusTopic.BroadcastVoteSet()
		} else if pld.Pulse.LessThan(hrs) {
			syncer.logger.Trace("Our consensus is ahead of this peer.")
		} else {
			syncer.logger.Trace("Our consensus is at the same step with this peer.")
		}
	}
}
