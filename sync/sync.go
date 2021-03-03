package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/sync/firewall"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/network_api"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
)

// IMPORTANT NOTE
//
// Sync module is based on pulling, not pushing
// Means if a node is behind the network, we don't send him anything
// The node should request (pull) itself.

const FlagInitialBlockDownload = 0x1

type publishMessageFn = func(msg *message.Message)
type syncedCallbackFn = func()

type synchronizer struct {
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
	consensusSync   *ConsensusSync
	stateSync       *StateSync
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
	broadcastCh <-chan *message.Message) (Synchronizer, error) {
	syncer := &synchronizer{
		ctx:         context.Background(),
		config:      conf,
		signer:      signer,
		state:       state,
		consensus:   consensus,
		txPool:      txPool,
		broadcastCh: broadcastCh,
	}

	logger := logger.NewLogger("_sync", syncer)
	peerSet := peerset.NewPeerSet(conf.SessionTimeout)
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

	syncer.consensusSync = NewConsensusSync(conf, net.ID(), consensus, logger, syncer.publishMessage)
	syncer.stateSync = NewStateSync(conf, net.ID(), cache, state, txPool, peerSet, logger, syncer.publishMessage, syncer.synced)

	if conf.InitialBlockDownload {
		if err := syncer.joinDownloadTopic(); err != nil {
			return nil, err
		}
	}

	return syncer, nil
}

func (syncer *synchronizer) Start() error {
	if err := syncer.networkAPI.Start(); err != nil {
		return err
	}

	go syncer.broadcastLoop()

	syncer.heartBeatTicker = time.NewTicker(syncer.config.HeartBeatTimeout)
	go syncer.heartBeatTickerLoop()

	syncer.BroadcastSalam()

	timer := time.NewTimer(syncer.config.StartingTimeout)
	go func() {
		<-timer.C
		syncer.maybeSynced()
	}()

	return nil
}

func (syncer *synchronizer) Stop() {
	syncer.ctx.Done()
	syncer.networkAPI.Stop()
	syncer.heartBeatTicker.Stop()
}

func (syncer *synchronizer) joinDownloadTopic() error {
	if err := syncer.networkAPI.JoinDownloadTopic(); err != nil {
		return err
	}

	return nil
}

func (syncer *synchronizer) maybeSynced() {
	ourHeight := syncer.state.LastBlockHeight()
	networkHeight := syncer.peerSet.MaxClaimedHeight()

	if ourHeight >= networkHeight-1 {
		syncer.synced()
	}
}

func (syncer *synchronizer) synced() {
	syncer.logger.Debug("We are synced", "hrs", syncer.consensus.HRS())
	syncer.consensus.MoveToNewHeight()
}

func (syncer *synchronizer) heartBeatTickerLoop() {
	for {
		select {
		case <-syncer.ctx.Done():
			return
		case <-syncer.heartBeatTicker.C:
			syncer.broadcastHeartBeat()
		}
	}
}

func (syncer *synchronizer) broadcastLoop() {
	for {
		select {
		case <-syncer.ctx.Done():
			return

		case msg := <-syncer.broadcastCh:

			switch msg.PayloadType() {
			case payload.PayloadTypeQueryTransactions:
				pld := msg.Payload.(*payload.QueryTransactionsPayload)
				syncer.queryTransactions(pld.IDs)

			case payload.PayloadTypeQueryProposal:
				pld := msg.Payload.(*payload.QueryProposalPayload)
				syncer.queryProposal(pld.Height, pld.Round)

			case payload.PayloadTypeBlockAnnounce:
				pld := msg.Payload.(*payload.BlockAnnouncePayload)
				syncer.announceBlock(pld.Height, pld.Block, pld.Commit)

			case payload.PayloadTypeVote,
				payload.PayloadTypeProposal,
				payload.PayloadTypeTransactions:
				syncer.publishMessage(msg)

			default:
				panic("Unexpected message to broadcast")

			}
		}
	}
}
func (syncer *synchronizer) Fingerprint() string {
	return fmt.Sprintf("{☍ %d ⛃ %d ⇈ %d ↑ %d}",
		syncer.peerSet.Len(),
		syncer.cache.Len(),
		syncer.peerSet.MaxClaimedHeight(),
		syncer.state.LastBlockHeight())
}

func (syncer *synchronizer) sendBlocksRequestIfWeAreBehind() {
	if syncer.peerSet.HasAnyValidSession() {
		syncer.logger.Debug("We have open seasson")
		return
	}

	ourHeight := syncer.state.LastBlockHeight()
	claimedHeight := syncer.peerSet.MaxClaimedHeight()
	if claimedHeight > ourHeight {
		if claimedHeight > ourHeight+LatestBlockInterval {
			syncer.logger.Info("We are far behind the network, Join download topic")
			// TODO:
			// If peer doesn't respond, we should leave the topic
			// A byzantine peer can send an invalid height, then all the nodes will join download topic.
			// We should find a way to avoid it.
			if err := syncer.joinDownloadTopic(); err != nil {
				syncer.logger.Info("We can't join download topic", "err", err)
			} else {
				syncer.stateSync.RequestForMoreBlock()
			}
		} else {
			syncer.logger.Info("We are behind the network, Ask for more blocks")
			syncer.stateSync.RequestForLatestBlock()
		}
	}
}

func (syncer *synchronizer) ParsMessage(msg *message.Message, from peer.ID) {
	syncer.logger.Debug("Received a message", "from", util.FingerprintPeerID(from), "message", msg)

	switch msg.PayloadType() {
	case payload.PayloadTypeSalam:
		pld := msg.Payload.(*payload.SalamPayload)
		syncer.ProcessSalamPayload(pld)

	case payload.PayloadTypeAleyk:
		pld := msg.Payload.(*payload.AleykPayload)
		syncer.ProcessAleykPayload(pld)

	case payload.PayloadTypeHeartBeat:
		pld := msg.Payload.(*payload.HeartBeatPayload)
		syncer.processHeartBeatPayload(pld)

	case payload.PayloadTypeLatestBlocksRequest:
		pld := msg.Payload.(*payload.LatestBlocksRequestPayload)
		syncer.stateSync.ProcessLatestBlocksRequestPayload(pld)

	case payload.PayloadTypeLatestBlocksResponse:
		pld := msg.Payload.(*payload.LatestBlocksResponsePayload)
		syncer.stateSync.ProcessLatestBlocksResponsePayload(pld)

	case payload.PayloadTypeQueryTransactions:
		pld := msg.Payload.(*payload.QueryTransactionsPayload)
		if syncer.isPeerActiveValidator(pld.Querier) {
			syncer.stateSync.ProcessQueryTransactionsPayload(pld)
		}

	case payload.PayloadTypeTransactions:
		pld := msg.Payload.(*payload.TransactionsPayload)
		syncer.stateSync.ProcessTransactionsPayload(pld)

	case payload.PayloadTypeBlockAnnounce:
		pld := msg.Payload.(*payload.BlockAnnouncePayload)
		syncer.stateSync.ProcessBlockAnnouncePayload(pld)

	case payload.PayloadTypeQueryProposal:
		pld := msg.Payload.(*payload.QueryProposalPayload)
		if syncer.isPeerActiveValidator(pld.Querier) {
			syncer.consensusSync.ProcessQueryProposalPayload(pld)
		}

	case payload.PayloadTypeProposal:
		pld := msg.Payload.(*payload.ProposalPayload)
		syncer.consensusSync.ProcessProposalPayload(pld)
		syncer.cache.AddProposal(pld.Proposal)

	case payload.PayloadTypeVote:
		pld := msg.Payload.(*payload.VotePayload)
		syncer.consensusSync.ProcessVotePayload(pld)

	case payload.PayloadTypeQueryVotes:
		pld := msg.Payload.(*payload.QueryVotesPayload)
		if syncer.isPeerActiveValidator(pld.Querier) {
			syncer.consensusSync.ProcessQueryVotesPayload(pld)
		}

	case payload.PayloadTypeDownloadRequest:
		pld := msg.Payload.(*payload.DownloadRequestPayload)
		syncer.stateSync.ProcessDownloadRequestPayload(pld)

	case payload.PayloadTypeDownloadResponse:
		pld := msg.Payload.(*payload.DownloadResponsePayload)
		syncer.stateSync.ProcessDownloadResponsePayload(pld)

	default:
		syncer.logger.Error("Unknown message type", "type", msg.PayloadType())
	}

	syncer.sendBlocksRequestIfWeAreBehind()
}

func (syncer *synchronizer) broadcastHeartBeat() {
	hrs := syncer.consensus.HRS()

	// Probable we are syncing
	if !hrs.IsValid() {
		return
	}

	// Broadcast a random vote if the validator is an active validator
	if syncer.isThisActiveValidator() {
		if syncer.consensus.RoundProposal(hrs.Round()) == nil {
			// We don't have proposal yet
			syncer.queryProposal(hrs.Height(), hrs.Round())
		}
		v := syncer.consensus.PickRandomVote(hrs.Round())
		if v != nil {
			syncer.consensusSync.BroadcastVote(v)
		}
	}

	msg := message.NewHeartBeatMessage(syncer.networkAPI.SelfID(), syncer.state.LastBlockHash(), hrs)
	syncer.publishMessage(msg)
}

func (syncer *synchronizer) publishMessage(msg *message.Message) {
	err := syncer.networkAPI.PublishMessage(msg)

	if err != nil {
		syncer.logger.Error("Error on publishing message", "message", msg, "err", err)
	} else {
		syncer.logger.Debug("Publishing new message", "message", msg)
	}
}

func (syncer *synchronizer) processHeartBeatPayload(pld *payload.HeartBeatPayload) {
	syncer.logger.Trace("Process heartbeat payload", "pld", pld)

	hrs := syncer.consensus.HRS()

	if pld.Pulse.Height() == hrs.Height() {
		if pld.Pulse.GreaterThan(hrs) {
			if syncer.isThisActiveValidator() {
				syncer.logger.Info("Our consensus is behind of this peer.", "ours", hrs, "peer", pld.Pulse, "address", syncer.signer.Address().Fingerprint())

				syncer.queryVotes(hrs.Height(), hrs.Round())
			}
		} else if pld.Pulse.LessThan(hrs) {
			syncer.logger.Trace("Our consensus is ahead of this peer.")
		} else {
			syncer.logger.Trace("Our consensus is at the same step with this peer.")
		}
	}

	p := syncer.peerSet.MustGetPeer(pld.PeerID)
	p.UpdateHeight(pld.Pulse.Height() - 1)
	syncer.peerSet.UpdateMaxClaimedHeight(pld.Pulse.Height() - 1)

	syncer.sendBlocksRequestIfWeAreBehind()
}

func (syncer *synchronizer) BroadcastSalam() {
	flags := 0
	if syncer.config.InitialBlockDownload {
		flags = util.SetFlag(flags, FlagInitialBlockDownload)
	}
	msg := message.NewSalamMessage(
		syncer.config.Moniker,
		syncer.signer.PublicKey(),
		syncer.networkAPI.SelfID(),
		syncer.state.GenesisHash(),
		syncer.state.LastBlockHeight(),
		flags)

	syncer.publishMessage(msg)
}

func (syncer *synchronizer) BroadcastAleyk(code payload.ResponseCode, resMsg string) {
	flags := 0
	if syncer.config.InitialBlockDownload {
		flags = util.SetFlag(flags, FlagInitialBlockDownload)
	}
	msg := message.NewAleykMessage(
		code,
		resMsg,
		syncer.config.Moniker,
		syncer.signer.PublicKey(),
		syncer.networkAPI.SelfID(),
		syncer.state.LastBlockHeight(),
		flags)

	syncer.publishMessage(msg)
}

func (syncer *synchronizer) ProcessSalamPayload(pld *payload.SalamPayload) {
	syncer.logger.Trace("Process salam payload", "pld", pld)

	if !pld.GenesisHash.EqualsTo(syncer.state.GenesisHash()) {
		syncer.logger.Info("Received a message from different chain", "genesis_hash", pld.GenesisHash)
		// Reply salam
		syncer.BroadcastAleyk(payload.ResponseCodeRejected, "Invalid genesis hash")
		return
	}

	p := syncer.peerSet.MustGetPeer(pld.PeerID)
	p.UpdateMoniker(pld.Moniker)
	p.UpdateHeight(pld.Height)
	p.UpdateNodeVersion(pld.NodeVersion)
	p.UpdatePublicKey(pld.PublicKey)
	p.UpdateInitialBlockDownload(util.IsFlagSet(pld.Flags, FlagInitialBlockDownload))

	syncer.peerSet.UpdateMaxClaimedHeight(pld.Height)

	// Reply salam
	syncer.BroadcastAleyk(payload.ResponseCodeOK, "Welcome!")
}

func (syncer *synchronizer) ProcessAleykPayload(pld *payload.AleykPayload) {
	syncer.logger.Trace("Process Aleyk payload", "pld", pld)

	if pld.ResponseCode != payload.ResponseCodeOK {
		syncer.logger.Warn("Our Salam is not welcomed!", "message", pld.ResponseMessage)
	} else {
		p := syncer.peerSet.MustGetPeer(pld.PeerID)
		p.UpdateMoniker(pld.Moniker)
		p.UpdateHeight(pld.Height)
		p.UpdateNodeVersion(pld.NodeVersion)
		p.UpdatePublicKey(pld.PublicKey)
		p.UpdateInitialBlockDownload(util.IsFlagSet(pld.Flags, FlagInitialBlockDownload))

		syncer.peerSet.UpdateMaxClaimedHeight(pld.Height)
	}
}

// isPeerActiveValidator checks if the peer is an active validator
func (syncer *synchronizer) isPeerActiveValidator(id peer.ID) bool {
	p := syncer.peerSet.GetPeer(id)
	if p == nil {
		return false
	}

	addr := p.PublicKey().Address()
	valSet := syncer.state.ValidatorSet()

	return valSet.Contains(addr)
}

// isThisActiveValidator checks if we are an active validator
func (syncer *synchronizer) isThisActiveValidator() bool {
	valSet := syncer.state.ValidatorSet()
	return valSet.Contains(syncer.signer.Address())
}

// queryTransactions queries for a missed transactions if we don't have it in the cache
// Only active validators can send this messsage
func (syncer *synchronizer) queryTransactions(ids []tx.ID) {
	missed := []crypto.Hash{}
	for i, id := range ids {
		trx := syncer.cache.GetTransaction(id)
		if trx != nil {
			if err := syncer.txPool.AppendTx(trx); err != nil {
				syncer.logger.Warn("Query for an invalid transaction", "tx", trx)
			}
		} else {
			missed = append(missed, ids[i])
		}
	}

	if len(missed) > 0 {
		if syncer.isThisActiveValidator() {
			syncer.consensusSync.BroadcastQueryTransaction(missed)
		}
	}
}

// queryProposal queries for proposal if we don't have it in the cache
// Only active validators can send this messsage
func (syncer *synchronizer) queryProposal(height, round int) {
	p := syncer.consensus.RoundProposal(round)
	if p == nil {
		p = syncer.cache.GetProposal(height, round)
		if p != nil {
			// We have the proposal inside the cache
			syncer.consensus.SetProposal(p)
		} else {
			if syncer.isThisActiveValidator() {
				syncer.consensusSync.BroadcastQueryProposal(height, round)
			} else {
				syncer.logger.Debug("queryProposal ignored. Not an active validator")
			}
		}
	}
}

// queryVotes asks other peers to send us some votes randomly
// Only active validators can send this messsage
func (syncer *synchronizer) queryVotes(height, round int) {
	if syncer.isThisActiveValidator() {
		syncer.consensusSync.BroadcastQueryVotes(height, round)
	} else {
		syncer.logger.Debug("queryVotes ignored. Not an active validator")
	}
}

func (syncer *synchronizer) announceBlock(height int, block *block.Block, commit *block.Commit) {
	if !syncer.isThisActiveValidator() {
		return
	}

	syncer.stateSync.BroadcastBlockAnnounce(height, block, commit)
}

func (syncer *synchronizer) PeerID() peer.ID {
	return syncer.networkAPI.SelfID()
}

func (syncer *synchronizer) Peers() []*peerset.Peer {
	return syncer.peerSet.GetPeerList()
}
