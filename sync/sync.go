package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/cache"
	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/logger"
)

// IMPORTANT NOTES
//
// Sync module is based on pulling, not pushing. It means,
// network doesn't update a node (push),
// a node itself should update itself (pull).
//
// Synchronizer should not have any locks to prevent dead lock situations.
// All submodules like state or consensus should be thread safe.

type synchronizer struct {
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

	if sync.config.HeartBeatTimer > 0 {
		sync.heartBeatTicker = time.NewTicker(sync.config.HeartBeatTimer)
		go sync.heartBeatTickerLoop()
	}

	sync.sayHello(false)
	sync.moveConsensusToNewHeight()

	return nil
}

func (sync *synchronizer) Stop() {
	sync.ctx.Done()
	if sync.heartBeatTicker != nil {
		sync.heartBeatTicker.Stop()
	}
}

func (sync *synchronizer) moveConsensusToNewHeight() {
	if sync.state.IsValidator(sync.signer.Address()) {
		if err := sync.network.JoinConsensusTopic(); err != nil {
			sync.logger.Error("error on joining consensus topic", "err", err)
		}
		sync.consensus.MoveToNewHeight()
	}
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
	if height > 0 {
		msg := message.NewHeartBeatMessage(height, round, sync.state.LastBlockHash())
		sync.broadcast(msg)
	}
}

func (sync *synchronizer) sayHello(helloAck bool) {
	flags := 0
	if sync.config.NodeNetwork {
		flags = util.SetFlag(flags, message.FlagNodeNetwork)
	}
	if helloAck {
		flags = util.SetFlag(flags, message.FlagHelloAck)
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
				if err := se.Reader.Close(); err != nil {
					sync.logger.Warn("error on closing stream", "err", err)
				}
			}

			err := sync.processIncomingBundle(bdl)
			if err != nil {
				sync.logger.Warn("error on parsing a bundle", "initiator", bdl.Initiator, "bundle", bdl, "err", err)
				sync.peerSet.IncreaseInvalidBundlesCounter(bdl.Initiator)
			}
		}
	}
}

func (sync *synchronizer) processIncomingBundle(bdl *bundle.Bundle) error {
	if bdl == nil {
		return nil
	}

	sync.logger.Info("received a bundle", "initiator", bdl.Initiator, "bundle", bdl)
	h := sync.handlers[bdl.Message.Type()]
	if h == nil {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid message type: %v", bdl.Message.Type())
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

// updateBlockchain checks if the node height is shorter than the network or not.
// If the node height is shorter than network more than two hours (720 blocks),
// it should start downloading the blocks from node networks,
// otherwise the node can request the latest blocks from the network.
func (sync *synchronizer) updateBlockchain() {
	// TODO: write test for me
	if sync.peerSet.HasAnyOpenSession() {
		sync.logger.Debug("we have open session")
		return
	}

	ourHeight := sync.state.LastBlockHeight()
	claimedHeight := sync.peerSet.MaxClaimedHeight()
	if claimedHeight > ourHeight {
		from := ourHeight
		// Make sure we have committed the latest blocks inside the cache
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
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagCarrierLibP2P)

		if sync.state.Params().IsMainnet() {
			bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
		}

		if sync.state.Params().IsTestnet() {
			bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
		}
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
			sync.logger.Error("error on sending bundle", "bundle", bdl, "err", err, "to", to)
		} else {
			sync.logger.Info("sending bundle to a peer", "bundle", bdl, "to", to)
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
			sync.logger.Error("error on broadcasting bundle", "bundle", bdl, "err", err)
		} else {
			sync.logger.Info("broadcasting new bundle", "bundle", bdl)
		}
	}
}

func (sync *synchronizer) SelfID() peer.ID {
	return sync.network.SelfID()
}

func (sync *synchronizer) Peers() []peerset.Peer {
	return sync.peerSet.GetPeerList()
}

// downloadBlocks starts downloading blocks from the network.
func (sync *synchronizer) downloadBlocks(from uint32) {
	sync.logger.Debug("downloading blocks", "from", from)
	l := sync.peerSet.GetPeerList()
	for _, peer := range l {
		// TODO: write test for me
		if sync.peerSet.NumberOfOpenSessions() > sync.config.MaxOpenSessions {
			sync.logger.Debug("we reached maximum number of open sessions")
			break
		}

		// TODO: write test for me
		if !peer.IsNodeNetwork() {
			continue
		}

		// TODO: write test for me
		if peer.Status != peerset.StatusCodeKnown {
			continue
		}

		// TODO: write test for me
		if peer.Height < from+1 {
			continue
		}
		to := from + LatestBlockInterval
		if to > peer.Height {
			to = peer.Height
		}

		sync.logger.Debug("sending download request", "from", from+1, "to", to, "pid", peer.PeerID)
		session := sync.peerSet.OpenSession(peer.PeerID)
		msg := message.NewBlocksRequestMessage(session.SessionID(), from+1, to)
		sync.sendTo(msg, peer.PeerID)
		from = to
	}
}

func (sync *synchronizer) queryLatestBlocks(from uint32) {
	randPeer := sync.peerSet.GetRandomPeer()

	// TODO: write test for me
	if !randPeer.IsKnownOrTrusty() {
		return
	}

	// TODO: write test for me
	to := randPeer.Height
	if randPeer.Height < from+1 {
		return
	}

	sync.logger.Debug("querying the latest blocks", "from", from+1, "to", to, "pid", randPeer.PeerID)
	session := sync.peerSet.OpenSession(randPeer.PeerID)
	msg := message.NewBlocksRequestMessage(session.SessionID(), from+1, to)
	sync.sendTo(msg, randPeer.PeerID)
}

// peerIsInTheCommittee checks if the peer is a member of the committee
// at the current height.
func (sync *synchronizer) peerIsInTheCommittee(id peer.ID) bool {
	p := sync.peerSet.GetPeer(id)
	if !p.IsKnownOrTrusty() {
		return false
	}

	return sync.state.IsInCommittee(p.Address())
}

// weAreInTheCommittee checks if we are a member of the committee
// at the current height.
func (sync *synchronizer) weAreInTheCommittee() bool {
	return sync.state.IsInCommittee(sync.signer.PublicKey().Address())
}

func (sync *synchronizer) tryCommitBlocks() {
	height := sync.state.LastBlockHeight() + 1
	for {
		b := sync.cache.GetBlock(height)
		if b == nil {
			break
		}
		c := sync.cache.GetCertificate(height)
		if c == nil {
			break
		}
		sync.logger.Trace("committing block", "height", height, "block", b)
		if err := sync.state.CommitBlock(height, b, c); err != nil {
			sync.logger.Warn("committing block failed", "block", b, "err", err, "height", height)
			// We will ask network to re-send this block again ...
			break
		}
		height = height + 1
	}
}

func (sync *synchronizer) prepareBlocks(from uint32, count uint32) []*block.Block {
	ourHeight := sync.state.LastBlockHeight()

	if from > ourHeight {
		sync.logger.Debug("we don't have block at this height", "height", from)
		return nil
	}

	if from+count > ourHeight {
		count = ourHeight - from + 1
	}

	blocks := make([]*block.Block, 0, count)

	for h := from; h < from+count; h++ {
		b := sync.cache.GetBlock(h)
		if b == nil {
			sync.logger.Warn("unable to find a block", "height", h)
			return nil
		}

		blocks = append(blocks, b)
	}

	return blocks
}

func (sync *synchronizer) updateSession(sessionID int, pid peer.ID, code message.ResponseCode) {
	s := sync.peerSet.FindSession(sessionID)
	if s == nil {
		sync.logger.Debug("session not found or closed", "session-id", sessionID)
		return
	}

	if s.PeerID() != pid {
		sync.logger.Warn("peer ID is not known", "session-id", sessionID, "pid", pid)
		return
	}

	s.SetLastResponseCode(code)

	switch code {
	case message.ResponseCodeRejected:
		sync.logger.Debug("session rejected, close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.updateBlockchain()

	case message.ResponseCodeBusy:
		sync.logger.Debug("peer is busy. close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.updateBlockchain()

	case message.ResponseCodeMoreBlocks:
		sync.logger.Debug("peer responding us. keep session open", "session-id", sessionID)

	case message.ResponseCodeNoMoreBlocks:
		sync.logger.Debug("peer has no more block. close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.updateBlockchain()

	case message.ResponseCodeSynced:
		sync.logger.Debug("peer informed us we are synced. close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.moveConsensusToNewHeight()
	}
}
