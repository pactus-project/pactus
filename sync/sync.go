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
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
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
	firewall        *firewall.Firewall
	cache           *cache.Cache
	handlers        map[payload.PayloadType]payloadHandler
	broadcastCh     <-chan payload.Payload
	network         network.Network
	heartBeatTicker *time.Ticker
	logger          *logger.Logger
}

func NewSynchronizer(
	conf *Config,
	signer crypto.Signer,
	state state.StateFacade,
	consensus consensus.Consensus,
	net network.Network,
	broadcastCh <-chan payload.Payload) (Synchronizer, error) {
	sync := &synchronizer{
		ctx:         context.Background(),
		config:      conf,
		signer:      signer,
		state:       state,
		consensus:   consensus,
		network:     net,
		broadcastCh: broadcastCh,
	}

	peerSet := peerset.NewPeerSet(conf.SessionTimeout)
	firewall := firewall.NewFirewall(peerSet, state)
	logger := logger.NewLogger("_sync", sync)
	cache, err := cache.NewCache(conf.CacheSize, state)
	if err != nil {
		return nil, err
	}

	sync.logger = logger
	sync.cache = cache
	sync.peerSet = peerSet
	sync.firewall = firewall

	handlers := make(map[payload.PayloadType]payloadHandler)

	handlers[payload.PayloadTypeSalam] = newSalamHandler(sync)
	handlers[payload.PayloadTypeAleyk] = newAleykHandler(sync)
	handlers[payload.PayloadTypeHeartBeat] = newHeartBeatHandler(sync)
	handlers[payload.PayloadTypeVote] = newVoteHandler(sync)
	handlers[payload.PayloadTypeProposal] = newProposalHandler(sync)
	handlers[payload.PayloadTypeTransactions] = newTransactionsHandler(sync)
	handlers[payload.PayloadTypeHeartBeat] = newHeartBeatHandler(sync)
	handlers[payload.PayloadTypeQueryVotes] = newQueryVotesHandler(sync)
	handlers[payload.PayloadTypeQueryProposal] = newQueryProposalHandler(sync)
	handlers[payload.PayloadTypeQueryTransactions] = newQueryTransactionsHandler(sync)
	handlers[payload.PayloadTypeBlockAnnounce] = newBlockAnnounceHandler(sync)
	handlers[payload.PayloadTypeDownloadRequest] = newDownloadRequestHandler(sync)
	handlers[payload.PayloadTypeDownloadResponse] = newDownloadResponseHandler(sync)
	handlers[payload.PayloadTypeLatestBlocksRequest] = newLatestBlocksRequestHandler(sync)
	handlers[payload.PayloadTypeLatestBlocksResponse] = newLatestBlocksResponseHandler(sync)

	sync.handlers = handlers

	if conf.InitialBlockDownload {
		if err := sync.joinDownloadTopic(); err != nil {
			return nil, err
		}
	}
	return sync, nil
}

func (sync *synchronizer) Start() error {
	if err := sync.network.JoinTopics(sync.onReceiveData); err != nil {
		return err
	}

	go sync.broadcastLoop()

	sync.heartBeatTicker = time.NewTicker(sync.config.HeartBeatTimeout)
	go sync.heartBeatTickerLoop()

	timer := time.NewTimer(sync.config.StartingTimeout)
	go func() {
		<-timer.C
		sync.onStartingTimeout()
	}()

	sync.broadcastSalam()

	return nil
}

func (sync *synchronizer) Stop() {
	sync.ctx.Done()
	sync.heartBeatTicker.Stop()
}

func (sync *synchronizer) joinDownloadTopic() error {
	if err := sync.network.JoinDownloadTopic(); err != nil {
		return err
	}

	return nil
}

func (sync *synchronizer) onStartingTimeout() {
	ourHeight := sync.state.LastBlockHeight()
	networkHeight := sync.peerSet.MaxClaimedHeight()

	if ourHeight >= networkHeight-1 {
		sync.synced()
	}
}

func (sync *synchronizer) synced() {
	sync.logger.Debug("We are synced", "height", sync.state.LastBlockHeight())
	sync.consensus.MoveToNewHeight()
}

func (sync *synchronizer) heartBeatTickerLoop() {
	for {
		select {
		case <-sync.ctx.Done():
			return
		case <-sync.heartBeatTicker.C:
			sync.broadcastHeartBeat()
		}
	}
}

func (sync *synchronizer) broadcastHeartBeat() {
	// Broadcast a random vote if the validator is an active validator
	if sync.weAreInTheCommittee() {
		v := sync.consensus.PickRandomVote()
		if v != nil {
			pld := payload.NewVotePayload(v)
			sync.broadcast(pld)
		}
	}

	height, round := sync.consensus.HeightRound()
	pld := payload.NewHeartBeatPayload(height, round, sync.state.LastBlockHash())
	sync.broadcast(pld)
}

func (sync *synchronizer) broadcastSalam() {
	flags := 0
	if sync.config.InitialBlockDownload {
		flags = util.SetFlag(flags, FlagInitialBlockDownload)
	}
	pld := payload.NewSalamPayload(
		sync.config.Moniker,
		sync.signer.PublicKey(),
		sync.state.GenesisHash(),
		sync.state.LastBlockHeight(),
		flags)

	sync.broadcast(pld)
}

func (sync *synchronizer) broadcastAleyk(code payload.ResponseCode, resMsg string) {
	flags := 0
	if sync.config.InitialBlockDownload {
		flags = util.SetFlag(flags, FlagInitialBlockDownload)
	}
	response := payload.NewAleykPayload(
		code,
		resMsg,
		sync.config.Moniker,
		sync.signer.PublicKey(),
		sync.state.LastBlockHeight(),
		flags)

	sync.broadcast(response)
}

func (sync *synchronizer) broadcastLoop() {
	for {
		select {
		case <-sync.ctx.Done():
			return

		case pld := <-sync.broadcastCh:
			sync.broadcast(pld)
		}
	}
}
func (sync *synchronizer) Fingerprint() string {
	return fmt.Sprintf("{☍ %d ⛃ %d ⇈ %d ↑ %d}",
		sync.peerSet.Len(),
		sync.cache.Len(),
		sync.peerSet.MaxClaimedHeight(),
		sync.state.LastBlockHeight())
}

func (sync *synchronizer) sendBlocksRequestIfWeAreBehind() {
	if sync.peerSet.HasAnyValidSession() {
		sync.logger.Debug("We have open seasson")
		return
	}

	ourHeight := sync.state.LastBlockHeight()
	claimedHeight := sync.peerSet.MaxClaimedHeight()
	if claimedHeight > ourHeight {
		if claimedHeight > ourHeight+sync.config.RequestBlockInterval {
			sync.logger.Info("We are far behind the network, Join download topic")
			// TODO:
			// If peer doesn't respond, we should leave the topic
			// A byzantine peer can send an invalid height, then all the nodes will join download topic.
			// We should find a way to avoid it.
			if err := sync.joinDownloadTopic(); err != nil {
				sync.logger.Info("We can't join download topic", "err", err)
			} else {
				sync.RequestForMoreBlock()
			}
		} else {
			sync.logger.Info("We are behind the network, Ask for more blocks")
			sync.RequestForLatestBlock()
		}
	}
}

func (sync *synchronizer) onReceiveData(data []byte, from peer.ID) {
	msg := sync.firewall.OpenMessage(data, from)
	if msg == nil {
		return
	}

	sync.logger.Debug("Received a message", "from", util.FingerprintPeerID(from), "message", msg)
	handler := sync.handlers[msg.Payload.Type()]
	if handler == nil {
		// TODO: mark as bad message
		sync.logger.Warn("Invalid payload type: %v", msg.Payload.Type())
		return
	}

	if err := handler.ParsPayload(msg.Payload, msg.Initiator); err != nil {
		// TODO: mark as bad message
		sync.logger.Warn("Error on parsing a message", "from", util.FingerprintPeerID(from), "message", msg, "err", err)
		return
	}

	sync.sendBlocksRequestIfWeAreBehind()
}

func (sync *synchronizer) broadcast(pld payload.Payload) {
	handler := sync.handlers[pld.Type()]
	if handler == nil {
		sync.logger.Warn("Invalid payload type: %v", pld.Type())
		return
	}
	msg := handler.PrepareMessage(pld)
	if msg != nil {
		err := sync.network.PublishMessage(msg)
		if err != nil {
			sync.logger.Error("Error on publishing message", "message", msg, "err", err)
		} else {
			sync.logger.Debug("Publishing new message", "message", msg)
		}
	}
}

func (sync *synchronizer) SelfID() peer.ID {
	return sync.network.SelfID()
}

func (sync *synchronizer) Peers() []*peerset.Peer {
	return sync.peerSet.GetPeerList()
}

// TODO:
// maximum nodes to query block should be 8
//
func (sync *synchronizer) RequestForMoreBlock() {
	from := sync.state.LastBlockHeight()
	l := sync.peerSet.GetPeerList()
	for _, p := range l {
		if p.InitialBlockDownload() {
			if p.Height() > from {
				to := from + sync.config.RequestBlockInterval
				if to > p.Height() {
					to = p.Height()
				}

				session := sync.peerSet.OpenSession(p.PeerID())
				pld := payload.NewDownloadRequestPayload(session.SessionID(), p.PeerID(), from+1, to)
				sync.broadcast(pld)
				from = to
			}
		}
	}
}

func (sync *synchronizer) RequestForLatestBlock() {
	p := sync.peerSet.FindHighestPeer()
	if p != nil {
		session := sync.peerSet.OpenSession(p.PeerID())
		ourHeight := sync.state.LastBlockHeight()
		pld := payload.NewLatestBlocksRequestPayload(session.SessionID(), p.PeerID(), ourHeight+1, p.Height())
		sync.broadcast(pld)
	}
}

// peerIsInTheCommittee checks if the peer is a member of committee
func (sync *synchronizer) peerIsInTheCommittee(id peer.ID) bool {
	p := sync.peerSet.GetPeer(id)
	if p == nil {
		return false
	}

	addr := p.PublicKey().Address()
	return sync.state.IsInCommittee(addr)
}

// weAreInTheCommittee checks if we are a member of committee
func (sync *synchronizer) weAreInTheCommittee() bool {
	return sync.state.IsInCommittee(sync.signer.PublicKey().Address())
}

func (sync *synchronizer) tryCommitBlocks() {
	for {
		ourHeight := sync.state.LastBlockHeight()
		b := sync.cache.GetBlock(ourHeight + 1)
		if b == nil {
			break
		}
		c := sync.cache.GetCertificate(b.Hash())
		if c == nil {
			break
		}
		for _, id := range b.TxIDs().IDs() {
			if tx := sync.cache.GetTransaction(id); tx != nil {
				if err := sync.state.AddPendingTx(tx); err != nil {
					sync.logger.Trace("Error on appending pending transaction", "err", err)
				}
			}
		}
		sync.logger.Trace("Committing block", "height", ourHeight+1, "block", b)
		if err := sync.state.CommitBlock(ourHeight+1, b, c); err != nil {
			sync.logger.Warn("Committing block failed", "block", b, "err", err, "height", ourHeight+1)
			// We will ask peers to send this block later ...
			break
		}
	}
}

func (sync *synchronizer) prepareBlocksAndTransactions(from, count int) ([]*block.Block, []*tx.Tx) {
	ourHeight := sync.state.LastBlockHeight()

	if from > ourHeight {
		sync.logger.Debug("We don't have block at this height", "height", from)
		return nil, nil
	}

	if from+count > ourHeight {
		count = ourHeight - from + 1
	}

	blocks := make([]*block.Block, 0, count)
	trxs := make([]*tx.Tx, 0)

	for h := from; h < from+count; h++ {
		b := sync.cache.GetBlock(h)
		if b == nil {
			sync.logger.Warn("Unable to find a block", "height", h)
			return nil, nil
		}
		for _, id := range b.TxIDs().IDs() {
			trx := sync.cache.GetTransaction(id)
			if trx != nil {
				trxs = append(trxs, trx)
			} else {
				sync.logger.Debug("Unable to find a transaction", "id", id.Fingerprint())
				return nil, nil
			}
		}

		blocks = append(blocks, b)
	}

	return blocks, trxs
}

func (sync *synchronizer) prepareTransactions(ids []tx.ID) []*tx.Tx {
	trxs := make([]*tx.Tx, 0, len(ids))

	for _, id := range ids {
		trx := sync.cache.GetTransaction(id)
		if trx == nil {
			sync.logger.Debug("Unable to find a transaction", "id", id.Fingerprint())
			continue
		}
		trxs = append(trxs, trx)
	}
	return trxs
}

func (sync *synchronizer) updateSession(code payload.ResponseCode, sessionID int, initiator peer.ID, target peer.ID) {
	if target != sync.network.SelfID() {
		return
	}

	switch code {
	case payload.ResponseCodeRejected:
		sync.logger.Debug("session rejected, close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)

	case payload.ResponseCodeBusy:
		// TODO: handle this situation
		sync.logger.Debug("Peer is busy. close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)

	case payload.ResponseCodeMoreBlocks:
		sync.logger.Debug("Peer responding us. keep session open", "session-id", sessionID)

	case payload.ResponseCodeNoMoreBlocks:
		sync.logger.Debug("Peer has no more block. close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)

	case payload.ResponseCodeSynced:
		sync.logger.Debug("Peer infomed us we are synced. close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.synced()
	}

	s := sync.peerSet.FindSession(sessionID)
	if s == nil {
		sync.logger.Debug("Session not found or closed", "session-id", sessionID)
	} else {
		s.SetLastResponseCode(code)
	}
}
