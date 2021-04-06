package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/sync/firewall"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/network_api"
	"github.com/zarbchain/zarb-go/sync/parser"
	"github.com/zarbchain/zarb-go/sync/parser/handler"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/vote"
)

// IMPORTANT NOTE
//
// Sync module is based on pulling, not pushing
// Means if a node is behind the network, we don't send him anything
// The node should request (pull) itself.

const FlagInitialBlockDownload = 0x1

type synchronizer struct {
	// Not: Synchronizer should not have any lock to prevent dead lock situation.
	// Other modules like state or consesnus are thread safe

	ctx             context.Context
	config          *Config
	signer          crypto.Signer
	state           state.StateFacade
	consensus       consensus.Consensus
	peerSet         *peerset.PeerSet
	cache           *cache.Cache
	parser          *parser.Parser
	broadcastCh     <-chan payload.Payload
	networkAPI      network_api.NetworkAPI
	heartBeatTicker *time.Ticker
	logger          *logger.Logger
}

func NewSynchronizer(
	conf *Config,
	signer crypto.Signer,
	state state.StateFacade,
	consensus consensus.Consensus,
	net *network.Network,
	broadcastCh <-chan payload.Payload) (Synchronizer, error) {
	syncer := &synchronizer{
		ctx: context.Background(),
	}

	peerSet := peerset.NewPeerSet(conf.SessionTimeout)
	firewall := firewall.NewFirewall(peerSet, state)
	api, err := network_api.NewNetworkAPI(syncer.ctx, net, firewall, syncer.ParsMessage)
	if err != nil {
		return nil, err
	}
	err = syncer.new(conf, signer, state, consensus, api, peerSet, broadcastCh)
	if err != nil {
		return nil, err
	}
	return syncer, nil
}

func (syncer *synchronizer) new(
	conf *Config,
	signer crypto.Signer,
	state state.StateFacade,
	consensus consensus.Consensus,
	api network_api.NetworkAPI,
	peerSet *peerset.PeerSet,
	broadcastCh <-chan payload.Payload) error {
	logger := logger.NewLogger("_sync", syncer)
	cache, err := cache.NewCache(conf.CacheSize, state)
	if err != nil {
		return err
	}

	syncer.config = conf
	syncer.signer = signer
	syncer.state = state
	syncer.consensus = consensus
	syncer.logger = logger
	syncer.cache = cache
	syncer.peerSet = peerSet
	syncer.networkAPI = api
	syncer.broadcastCh = broadcastCh

	ctx := handler.NewHandlerContext(
		state,
		consensus,
		cache,
		peerSet,
		conf.Moniker,
		signer.PublicKey(),
		api.SelfID(),
		syncer.broadcast,
		syncer.synced,
		conf.BlockPerMessage,
		conf.InitialBlockDownload,
		conf.RequestBlockInterval,
		logger,
	)
	syncer.parser = parser.NewParser(ctx)

	if conf.InitialBlockDownload {
		if err := syncer.joinDownloadTopic(); err != nil {
			return err
		}
	}
	return nil
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
	syncer.logger.Debug("We are synced", "height", syncer.state.LastBlockHeight())
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

		case pld := <-syncer.broadcastCh:
			syncer.broadcast(pld)
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
		if claimedHeight > ourHeight+syncer.config.RequestBlockInterval {
			syncer.logger.Info("We are far behind the network, Join download topic")
			// TODO:
			// If peer doesn't respond, we should leave the topic
			// A byzantine peer can send an invalid height, then all the nodes will join download topic.
			// We should find a way to avoid it.
			if err := syncer.joinDownloadTopic(); err != nil {
				syncer.logger.Info("We can't join download topic", "err", err)
			} else {
				syncer.RequestForMoreBlock()
			}
		} else {
			syncer.logger.Info("We are behind the network, Ask for more blocks")
			syncer.RequestForLatestBlock()
		}
	}
}

func (syncer *synchronizer) ParsMessage(msg *message.Message, from peer.ID) {
	syncer.logger.Debug("Received a message", "from", util.FingerprintPeerID(from), "message", msg)

	if err := syncer.parser.ParsMessage(msg); err != nil {
		syncer.logger.Warn("Error on parsing a message", "from", util.FingerprintPeerID(from), "message", msg, "err", err)
		return
	}

	syncer.sendBlocksRequestIfWeAreBehind()
}

// TODO: Refactor this:
func (syncer *synchronizer) broadcastHeartBeat() {
	// Broadcast a random vote if the validator is an active validator
	if syncer.weAreInTheCommittee() {
		v := syncer.consensus.PickRandomVote()
		if v != nil {
			pld := payload.NewVotePayload(*v)
			syncer.broadcast(pld)
		}
	}

	height, round := syncer.consensus.HeightRound()
	pld := payload.NewHeartBeatPayload(height, round, syncer.state.LastBlockHash())
	syncer.broadcast(pld)
}

func (syncer *synchronizer) broadcast(pld payload.Payload) {
	msg := message.NewMessage(syncer.networkAPI.SelfID(), pld)

	switch pld.Type() {
	case payload.PayloadTypeLatestBlocksResponse,
		payload.PayloadTypeDownloadResponse,
		payload.PayloadTypeTransactions:
		msg.CompressIt()

	case payload.PayloadTypeQueryVotes:
		p := pld.(*payload.QueryVotesPayload)
		syncer.queryVotes(p.Height, p.Round)
		return

	case payload.PayloadTypeQueryProposal:
		p := pld.(*payload.QueryProposalPayload)
		syncer.queryProposal(p.Height, p.Round)
		return

	case payload.PayloadTypeQueryTransactions:
		p := pld.(*payload.QueryTransactionsPayload)
		syncer.queryTransactions(p.IDs)
		return
	}

	err := syncer.networkAPI.PublishMessage(msg)
	if err != nil {
		syncer.logger.Error("Error on publishing message", "message", msg, "err", err)
	} else {
		syncer.logger.Debug("Publishing new message", "message", msg)
	}
}

// queryTransactions queries for a missed transactions if we don't have it in the cache
// Only active validators can send this messsage
func (syncer *synchronizer) queryTransactions(ids []tx.ID) {
	missed := []crypto.Hash{}
	for i, id := range ids {
		trx := syncer.cache.GetTransaction(id)
		if trx != nil {
			if err := syncer.state.AddPendingTx(trx); err != nil {
				syncer.logger.Warn("Query for an invalid transaction", "tx", trx, "err", err)
			}
		} else {
			missed = append(missed, ids[i])
		}
	}

	if len(missed) > 0 {
		if syncer.weAreInTheCommittee() {
			syncer.BroadcastQueryTransaction(missed)
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
			if syncer.weAreInTheCommittee() {
				syncer.BroadcastQueryProposal(height, round)
			} else {
				syncer.logger.Debug("queryProposal ignored. Not an active validator")
			}
		}
	}
}

// queryVotes asks other peers to send us some votes randomly
// Only active validators can send this messsage
func (syncer *synchronizer) queryVotes(height, round int) {
	if syncer.weAreInTheCommittee() {
		syncer.BroadcastQueryVotes(height, round)
	} else {
		syncer.logger.Debug("queryVotes ignored. Not an active validator")
	}
}
func (syncer *synchronizer) BroadcastSalam() {
	flags := 0
	if syncer.config.InitialBlockDownload {
		flags = util.SetFlag(flags, FlagInitialBlockDownload)
	}
	pld := payload.NewSalamPayload(
		syncer.config.Moniker,
		syncer.signer.PublicKey(),
		syncer.state.GenesisHash(),
		syncer.state.LastBlockHeight(),
		flags)

	syncer.broadcast(pld)
}

func (syncer *synchronizer) BroadcastQueryProposal(height, round int) {
	msg := message.NewMessage(syncer.networkAPI.SelfID(), payload.NewQueryProposalPayload(height, round))
	err := syncer.networkAPI.PublishMessage(msg)
	if err != nil {
		syncer.logger.Error("Error on publishing message", "message", msg, "err", err)
	} else {
		syncer.logger.Debug("Publishing new message", "message", msg)
	}
}

func (syncer *synchronizer) BroadcastQueryVotes(height, round int) {
	msg := message.NewMessage(syncer.networkAPI.SelfID(), payload.NewQueryVotesPayload(height, round))
	err := syncer.networkAPI.PublishMessage(msg)
	if err != nil {
		syncer.logger.Error("Error on publishing message", "message", msg, "err", err)
	} else {
		syncer.logger.Debug("Publishing new message", "message", msg)
	}
}

func (syncer *synchronizer) BroadcastQueryTransaction(ids []tx.ID) {
	msg := message.NewMessage(syncer.networkAPI.SelfID(), payload.NewQueryTransactionsPayload(ids))
	err := syncer.networkAPI.PublishMessage(msg)
	if err != nil {
		syncer.logger.Error("Error on publishing message", "message", msg, "err", err)
	} else {
		syncer.logger.Debug("Publishing new message", "message", msg)
	}
}

// TODO: Move to test!?!??!
func (syncer *synchronizer) BroadcastProposal(p *proposal.Proposal) {
	pld := payload.NewProposalPayload(*p)
	syncer.broadcast(pld)
}

func (syncer *synchronizer) BroadcastVote(v *vote.Vote) {
	pld := payload.NewVotePayload(*v)
	syncer.broadcast(pld)
}

// weAreInTheCommittee checks if we are an active validator
func (syncer *synchronizer) weAreInTheCommittee() bool {
	return syncer.state.IsInCommittee(syncer.signer.Address())
}

func (syncer *synchronizer) PeerID() peer.ID {
	return syncer.networkAPI.SelfID()
}

func (syncer *synchronizer) Peers() []*peerset.Peer {
	return syncer.peerSet.GetPeerList()
}

// TODO:
// maximum nodes to query block should be 8
//
func (syncer *synchronizer) RequestForMoreBlock() {
	from := syncer.state.LastBlockHeight()
	l := syncer.peerSet.GetPeerList()
	for _, p := range l {
		if p.InitialBlockDownload() {
			if p.Height() > from {
				to := from + syncer.config.RequestBlockInterval
				if to > p.Height() {
					to = p.Height()
				}

				session := syncer.peerSet.OpenSession(p.PeerID())
				pld := payload.NewDownloadRequestPayload(session.SessionID, p.PeerID(), from+1, to)
				syncer.broadcast(pld)
				from = to
			}
		}
	}
}

func (syncer *synchronizer) RequestForLatestBlock() {
	p := syncer.peerSet.FindHighestPeer()
	if p != nil {
		session := syncer.peerSet.OpenSession(p.PeerID())
		ourHeight := syncer.state.LastBlockHeight()
		pld := payload.NewLatestBlocksRequestPayload(session.SessionID, p.PeerID(), ourHeight+1, p.Height())
		syncer.broadcast(pld)
	}
}
