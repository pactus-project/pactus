package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/bundle"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/sync/firewall"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

// IMPORTANT NOTE
//
// Sync module is based on pulling, not pushing.
// Means: Network doesn't update the node (push),
// The node itself should send request (pull).

type synchronizer struct {
	// Not: Synchronizer should not have any lock to prevent dead lock situation.
	// Other modules like state or consesnus are thread safe

	ctx             context.Context
	config          *Config
	signer          crypto.Signer
	state           state.Facade
	consensus       consensus.Consensus
	peerSet         *peerset.PeerSet
	firewall        *firewall.Firewall
	cache           *cache.Cache
	handlers        map[message.Type]messageHandler
	broadcastCh     <-chan message.Message
	networkCh       <-chan network.Event
	network         network.Network
	heartBeatTicker *time.Ticker
	logger          *logger.Logger
}

func NewSynchronizer(
	conf *Config,
	signer crypto.Signer,
	state state.Facade,
	consensus consensus.Consensus,
	net network.Network,
	broadcastCh <-chan message.Message) (Synchronizer, error) {
	sync := &synchronizer{
		ctx:         context.Background(), // TODO, set proper context
		config:      conf,
		signer:      signer,
		state:       state,
		consensus:   consensus,
		network:     net,
		broadcastCh: broadcastCh,
		networkCh:   net.EventChannel(),
	}

	peerSet := peerset.NewPeerSet(conf.SessionTimeout)
	logger := logger.NewLogger("_sync", sync)
	firewall := firewall.NewFirewall(conf.Firewall, net, peerSet, state, logger)
	cache, err := cache.NewCache(conf.CacheSize, state)
	if err != nil {
		return nil, err
	}

	sync.logger = logger
	sync.cache = cache
	sync.peerSet = peerSet
	sync.firewall = firewall

	handlers := make(map[message.Type]messageHandler)

	handlers[message.MessageTypeHello] = newHelloHandler(sync)
	handlers[message.MessageTypeHeartBeat] = newHeartBeatHandler(sync)
	handlers[message.MessageTypeVote] = newVoteHandler(sync)
	handlers[message.MessageTypeProposal] = newProposalHandler(sync)
	handlers[message.MessageTypeTransactions] = newTransactionsHandler(sync)
	handlers[message.MessageTypeHeartBeat] = newHeartBeatHandler(sync)
	handlers[message.MessageTypeQueryVotes] = newQueryVotesHandler(sync)
	handlers[message.MessageTypeQueryProposal] = newQueryProposalHandler(sync)
	handlers[message.MessageTypeQueryTransactions] = newQueryTransactionsHandler(sync)
	handlers[message.MessageTypeBlockAnnounce] = newBlockAnnounceHandler(sync)
	handlers[message.MessageTypeBlocksRequest] = newBlocksRequestHandler(sync)
	handlers[message.MessageTypeBlocksResponse] = newBlocksResponseHandler(sync)

	sync.handlers = handlers

	return sync, nil
}

func (sync *synchronizer) Start() error {
	if err := sync.network.JoinGeneralTopic(); err != nil {
		return err
	}

	go sync.receiveLoop()
	go sync.broadcastLoop()

	sync.heartBeatTicker = time.NewTicker(sync.config.HeartBeatTimeout)
	go sync.heartBeatTickerLoop()

	timer := time.NewTimer(sync.config.StartingTimeout)
	go func() {
		<-timer.C
		sync.onStartingTimeout()
	}()

	sync.sayHello(true)

	return nil
}

func (sync *synchronizer) Stop() {
	sync.ctx.Done()
	sync.heartBeatTicker.Stop()
}

func (sync *synchronizer) onStartingTimeout() {
	ourHeight := sync.state.LastBlockHeight()
	networkHeight := sync.peerSet.MaxClaimedHeight()

	if ourHeight >= networkHeight-1 {
		sync.synced()
	}
}

func (sync *synchronizer) synced() {
	sync.logger.Debug("we are synced", "height", sync.state.LastBlockHeight())
	if err := sync.network.JoinConsensusTopic(); err != nil {
		sync.logger.Error("error on joining consensus topic", "err", err)
	}
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
			msg := message.NewVoteMessage(v)
			sync.broadcast(msg)
		}
	}

	height, round := sync.consensus.HeightRound()
	msg := message.NewHeartBeatMessage(height, round, sync.state.LastBlockHash())
	sync.broadcast(msg)
}

func (sync *synchronizer) sayHello(needResponse bool) {
	flags := 0
	if sync.config.InitialBlockDownload {
		flags = util.SetFlag(flags, message.FlagInitialBlockDownload)
	}
	if needResponse {
		flags = util.SetFlag(flags, message.FlagNeedResponse)
	}
	msg := message.NewHelloMessage(
		sync.SelfID(),
		sync.config.Moniker,
		sync.state.LastBlockHeight(),
		flags, sync.state.GenesisHash())
	sync.signer.SignMsg(msg)

	sync.broadcast(msg)
}

func (sync *synchronizer) broadcastLoop() {
	for {
		select {
		case <-sync.ctx.Done():
			return

		case msg := <-sync.broadcastCh:
			sync.broadcast(msg)
		}
	}
}
func (sync *synchronizer) receiveLoop() {
	for {
		select {
		case <-sync.ctx.Done():
			return

		case e := <-sync.networkCh:
			var bdl *bundle.Bundle
			switch e.Type() {
			case network.EventTypeGossip:
				ge := e.(*network.GossipMessage)
				bdl = sync.firewall.OpenGossipBundle(ge.Data, ge.Source, ge.From)

			case network.EventTypeStream:
				se := e.(*network.StreamMessage)
				bdl = sync.firewall.OpenStreamBundle(se.Reader, se.Source)
			}

			err := sync.processIncomingBundle(bdl)
			if err != nil {
				sync.logger.Warn("error on parsing a message", "initiator", util.FingerprintPeerID(bdl.Initiator), "message", bdl, "err", err)
				peer := sync.peerSet.MustGetPeer(bdl.Initiator)
				peer.IncreaseInvalidBundlesCounter()
			}
		}
	}
}

func (sync *synchronizer) processIncomingBundle(bdl *bundle.Bundle) error {
	if bdl == nil {
		return nil
	}

	sync.logger.Debug("received a message", "initiator", util.FingerprintPeerID(bdl.Initiator), "bundle", bdl)
	h := sync.handlers[bdl.Message.Type()]
	if h == nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid message type: %v", bdl.Message.Type())
	}

	return h.ParsMessage(bdl.Message, bdl.Initiator)
}

func (sync *synchronizer) Fingerprint() string {
	return fmt.Sprintf("{☍ %d ⛃ %d ⇈ %d ↑ %d}",
		sync.peerSet.Len(),
		sync.cache.Len(),
		sync.peerSet.MaxClaimedHeight(),
		sync.state.LastBlockHeight())
}

// updateBlokchain checks if the node height is shorter than the network or not.
// If the node height is shorter than network more than two hours (720 blocks),
// it should start downloading the blocks,
// otherwise the node can request the latest blocks from the network.
func (sync *synchronizer) updateBlokchain() {
	// TODO: write test for me
	if sync.peerSet.HasAnyOpenSession() {
		sync.logger.Debug("we have open session")
		return
	}

	ourHeight := sync.state.LastBlockHeight()
	claimedHeight := sync.peerSet.MaxClaimedHeight()
	if claimedHeight > ourHeight {
		from := ourHeight
		// Make sure we done have the start block inside cache
		for sync.cache.HasBlockInCache(from + 1) {
			from++
		}

		sync.logger.Info("start syncing with the network")
		if claimedHeight > ourHeight+LatestBlockInterval {
			sync.downloadBlocks(from)
		} else {
			sync.queryLatestBlocks(from)
		}
	}
}

func (sync *synchronizer) prepareBundle(msg message.Message) *bundle.Bundle {
	h := sync.handlers[msg.Type()]
	if h == nil {
		sync.logger.Warn("invalid message type: %v", msg.Type())
		return nil
	}
	bdl := h.PrepareBundle(msg)
	if bdl != nil {
		if err := bdl.SanityCheck(); err != nil {
			sync.logger.Error("broadcasting an invalid bundle", "bundle", bdl, "err", err)
			return nil
		}

		// Bundles will be carried through LibP2P.
		// In future we might support other libraries.
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkLibP2P)
		return bdl
	}
	return nil
}

func (sync *synchronizer) sendTo(msg message.Message, to peer.ID) {
	bdl := sync.prepareBundle(msg)
	if bdl != nil {
		data, _ := bdl.Encode()
		err := sync.network.SendTo(data, to)
		if err != nil {
			sync.logger.Error("error on sending message", "message", bdl, "err", err)
		} else {
			sync.logger.Debug("sending message to a peer", "message", bdl, "to", to)
		}
	}
}

func (sync *synchronizer) broadcast(msg message.Message) {
	bdl := sync.prepareBundle(msg)
	if bdl != nil {
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagBroadcasted)
		data, _ := bdl.Encode()
		err := sync.network.Broadcast(data, msg.Type().TopicID())
		if err != nil {
			sync.logger.Error("error on broadcasting message", "message", msg, "err", err)
		} else {
			sync.logger.Debug("broadcasting new message", "message", msg)
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
func (sync *synchronizer) downloadBlocks(from int) {
	l := sync.peerSet.GetPeerList()
	for _, peer := range l {
		// TODO: write test for me
		if sync.peerSet.NumberOfOpenSessions() > sync.config.MaximumOpenSessions {
			sync.logger.Debug("we reached maximum number of open sessions")
			break
		}

		// TODO: write test for me
		if !peer.InitialBlockDownload() {
			continue
		}

		// TODO: write test for me
		if peer.Status() != peerset.StatusCodeKnown {
			continue
		}

		// TODO: write test for me
		if peer.Height() < from+1 {
			continue
		}
		to := from + LatestBlockInterval
		if to > peer.Height() {
			to = peer.Height()
		}

		sync.logger.Debug("sending download request", "from", from+1, "to", to, "pid", util.FingerprintPeerID(peer.PeerID()))
		session := sync.peerSet.OpenSession(peer.PeerID())
		msg := message.NewBlocksRequestMessage(session.SessionID(), from+1, to)
		sync.sendTo(msg, peer.PeerID())
		from = to
	}
}

func (sync *synchronizer) queryLatestBlocks(from int) {
	randPeer := sync.peerSet.GetRandomPeer()

	// TODO: write test for me
	if randPeer == nil || !randPeer.IsKnownOrTrusted() {
		return
	}

	// TODO: write test for me
	to := randPeer.Height()
	if randPeer.Height() < from+1 {
		return
	}

	sync.logger.Debug("querying the latest blocks", "from", from+1, "to", to, "pid", util.FingerprintPeerID(randPeer.PeerID()))
	session := sync.peerSet.OpenSession(randPeer.PeerID())
	msg := message.NewBlocksRequestMessage(session.SessionID(), from+1, to)
	sync.sendTo(msg, randPeer.PeerID())
}

/// peerIsInTheCommittee checks if the peer is a member of committee
func (sync *synchronizer) peerIsInTheCommittee(id peer.ID) bool {
	p := sync.peerSet.GetPeer(id)
	if p == nil {
		return false
	}

	if !p.IsKnownOrTrusted() {
		return false
	}

	return sync.state.IsInCommittee(p.Address())
}

/// weAreInTheCommittee checks if we are a member of committee
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
					sync.logger.Trace("error on appending a transaction", "err", err)
				}
			}
		}
		sync.logger.Trace("committing block", "height", ourHeight+1, "block", b)
		if err := sync.state.CommitBlock(ourHeight+1, b, c); err != nil {
			sync.logger.Warn("committing block failed", "block", b, "err", err, "height", ourHeight+1)
			// We will ask network to re-send this block again ...
			break
		}
	}
}

func (sync *synchronizer) prepareBlocksAndTransactions(from, count int) ([]*block.Block, []*tx.Tx) {
	ourHeight := sync.state.LastBlockHeight()

	if from > ourHeight {
		sync.logger.Debug("we don't have block at this height", "height", from)
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
			sync.logger.Warn("unable to find a block", "height", h)
			return nil, nil
		}
		for _, id := range b.TxIDs().IDs() {
			trx := sync.cache.GetTransaction(id)
			if trx != nil {
				trxs = append(trxs, trx)
			} else {
				sync.logger.Debug("unable to find a transaction", "id", id.Fingerprint())
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
			sync.logger.Debug("unable to find a transaction", "id", id.Fingerprint())
			continue
		}
		trxs = append(trxs, trx)
	}
	return trxs
}

func (sync *synchronizer) updateSession(sessionID int, pid peer.ID, code message.ResponseCode) {
	s := sync.peerSet.FindSession(sessionID)
	if s == nil {
		sync.logger.Debug("session not found or closed", "session-id", sessionID)
		return
	}

	if s.PeerID() != pid {
		sync.logger.Debug("peer ID is not known", "session-id", sessionID, "pid", pid)
		return
	}

	s.SetLastResponseCode(code)

	switch code {
	case message.ResponseCodeRejected:
		sync.logger.Debug("session rejected, close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.updateBlokchain()

	case message.ResponseCodeBusy:
		sync.logger.Debug("peer is busy. close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.updateBlokchain()

	case message.ResponseCodeMoreBlocks:
		sync.logger.Debug("peer responding us. keep session open", "session-id", sessionID)

	case message.ResponseCodeNoMoreBlocks:
		sync.logger.Debug("peer has no more block. close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.updateBlokchain()

	case message.ResponseCodeSynced:
		sync.logger.Debug("peer infomed us we are synced. close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.synced()
	}
}
