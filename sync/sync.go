package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/cache"
	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/logger"
)

// IMPORTANT NOTES:
//
// 1. The Sync module is based on pulling instead of pushing. This means that the network
// does not update a node (push); instead, a node should update itself (pull).
//
// 2. The Synchronizer should not have any locks to prevent deadlocks. All submodules,
// such as state or consensus, should be thread-safe.

type synchronizer struct {
	ctx             context.Context
	config          *Config
	signers         []crypto.Signer
	state           state.Facade
	consMgr         consensus.Manager
	peerSet         *peerset.PeerSet
	firewall        *firewall.Firewall
	cache           *cache.Cache
	handlers        map[message.Type]messageHandler
	broadcastCh     <-chan message.Message
	networkCh       <-chan network.Event
	network         network.Network
	heartBeatTicker *time.Ticker
	logger          *logger.SubLogger
}

func NewSynchronizer(
	conf *Config,
	signers []crypto.Signer,
	state state.Facade,
	consMgr consensus.Manager,
	net network.Network,
	broadcastCh <-chan message.Message) (Synchronizer, error) {
	sync := &synchronizer{
		ctx:         context.Background(), // TODO, set proper context
		config:      conf,
		signers:     signers,
		state:       state,
		consMgr:     consMgr,
		network:     net,
		broadcastCh: broadcastCh,
		networkCh:   net.EventChannel(),
	}

	peerSet := peerset.NewPeerSet(conf.SessionTimeout)
	logger := logger.NewSubLogger("_sync", sync)
	firewall := firewall.NewFirewall(conf.Firewall, net, peerSet, state, logger)
	cache, err := cache.NewCache(conf.CacheSize)
	if err != nil {
		return nil, err
	}

	sync.logger = logger
	sync.cache = cache
	sync.peerSet = peerSet
	sync.firewall = firewall

	handlers := make(map[message.Type]messageHandler)

	handlers[message.TypeHello] = newHelloHandler(sync)
	handlers[message.TypeHeartBeat] = newHeartBeatHandler(sync)
	handlers[message.TypeVote] = newVoteHandler(sync)
	handlers[message.TypeProposal] = newProposalHandler(sync)
	handlers[message.TypeTransactions] = newTransactionsHandler(sync)
	handlers[message.TypeHeartBeat] = newHeartBeatHandler(sync)
	handlers[message.TypeQueryVotes] = newQueryVotesHandler(sync)
	handlers[message.TypeQueryProposal] = newQueryProposalHandler(sync)
	handlers[message.TypeBlockAnnounce] = newBlockAnnounceHandler(sync)
	handlers[message.TypeBlocksRequest] = newBlocksRequestHandler(sync)
	handlers[message.TypeBlocksResponse] = newBlocksResponseHandler(sync)

	sync.handlers = handlers

	return sync, nil
}

func (sync *synchronizer) Start() error {
	if err := sync.network.JoinGeneralTopic(); err != nil {
		return err
	}
	// TODO: Not joining consensus topic when we are syncing
	if err := sync.network.JoinConsensusTopic(); err != nil {
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
	sync.consMgr.MoveToNewHeight()
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
	// Broadcast a random vote if we are inside the committee
	if sync.weAreInTheCommittee() {
		_, round := sync.consMgr.HeightRound()
		v := sync.consMgr.PickRandomVote(round)
		if v != nil {
			msg := message.NewVoteMessage(v)
			sync.broadcast(msg)
		}

		height, round := sync.consMgr.HeightRound()
		if height > 0 {
			msg := message.NewHeartBeatMessage(height, round, sync.state.LastBlockHash())
			sync.broadcast(msg)
		}
	} else {
		height := sync.state.LastBlockHeight()
		if height > 0 {
			msg := message.NewHeartBeatMessage(height, -1, sync.state.LastBlockHash())
			sync.broadcast(msg)
		}
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
		flags,
		sync.state.LastBlockHash(),
		sync.state.Genesis().Hash())

	for _, signer := range sync.signers {
		signer.SignMsg(msg)
		sync.broadcast(msg)
	}
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
					// TODO: write test for me
					sync.peerSet.IncreaseSendFailedCounter(se.Source)
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

	return h.ParseMessage(bdl.Message, bdl.Initiator)
}

func (sync *synchronizer) String() string {
	return fmt.Sprintf("{☍ %d ⛃ %d ⇈ %d ↑ %d}",
		sync.peerSet.Len(),
		sync.cache.Len(),
		sync.peerSet.MaxClaimedHeight(),
		sync.state.LastBlockHeight())
}

// updateBlockchain checks whether the node's height is shorter than the network's height or not.
// If the node's height is shorter than the network's height by more than two hours (720 blocks),
// it should start downloading blocks from the network's nodes.
// Otherwise, the node can request the latest blocks from the network.
func (sync *synchronizer) updateBlockchain() {
	// First, let's check if we have any open sessions.
	// If there are any open sessions, we should wait for them to be closed.
	// Otherwise, we can request the same blocks from different peers.
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
			sync.downloadBlocks(from, true)
		} else {
			sync.downloadBlocks(from, false)
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
		// Bundles will be carried through LibP2P.
		// In future we might support other libraries.
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagCarrierLibP2P)

		switch sync.state.Genesis().ChainType() {
		case genesis.Mainnet:
			bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
		case genesis.Testnet:
			bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
		default:
			// It's localnet and for testing purpose only
		}

		return bdl
	}
	return nil
}

func (sync *synchronizer) sendTo(msg message.Message, to peer.ID, sessionID int) {
	bdl := sync.prepareBundle(msg)
	if bdl != nil {
		data, _ := bdl.Encode()
		sync.peerSet.UpdateLastSent(to)
		err := sync.network.SendTo(data, to)
		if err != nil {
			sync.logger.Warn("error on sending bundle", "bundle", bdl, "err", err, "to", to)
			sync.peerSet.IncreaseSendFailedCounter(to)

			// Let's close the session with this peer because we couldn't establish a connection.
			// This helps to free sessions and ask blocks from other peers.
			sync.peerSet.CloseSession(sessionID)
		} else {
			sync.logger.Info("sending bundle to a peer", "bundle", bdl, "to", to)
			sync.peerSet.IncreaseSendSuccessCounter(to)
		}

		sync.peerSet.IncreaseSentBytesCounter(msg.Type(), int64(len(data)), &to)
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
		sync.peerSet.IncreaseSentBytesCounter(msg.Type(), int64(len(data)), nil)
	}
}

func (sync *synchronizer) SelfID() peer.ID {
	return sync.network.SelfID()
}

func (sync *synchronizer) Moniker() string {
	return sync.config.Moniker
}

func (sync *synchronizer) PeerSet() *peerset.PeerSet {
	return sync.peerSet
}

// downloadBlocks starts downloading blocks from the network.
func (sync *synchronizer) downloadBlocks(from uint32, onlyNodeNetwork bool) {
	sync.logger.Debug("downloading blocks", "from", from)
	for i := sync.peerSet.NumberOfOpenSessions(); i < sync.config.MaxOpenSessions; i++ {
		p := sync.peerSet.GetRandomPeer()

		// Don't open a new session if we already have an open session with the same peer.
		// This helps us to get blocks from different peers.
		// TODO: write test for me
		if sync.peerSet.HasOpenSession(p.PeerID) {
			continue
		}

		if onlyNodeNetwork && !p.IsNodeNetwork() {
			continue
		}

		// TODO: write test for me
		if p.Height < from+1 {
			continue
		}

		count := LatestBlockInterval
		sync.logger.Debug("sending download request", "from", from+1, "count", count, "pid", p.PeerID)
		session := sync.peerSet.OpenSession(p.PeerID)
		msg := message.NewBlocksRequestMessage(session.SessionID(), from+1, count)
		sync.sendTo(msg, p.PeerID, session.SessionID())
		from += count
	}
}

// peerIsInTheCommittee checks if the peer is a member of the committee
// at the current height.
func (sync *synchronizer) peerIsInTheCommittee(pid peer.ID) bool {
	p := sync.peerSet.GetPeer(pid)
	if !p.IsKnownOrTrusty() {
		return false
	}

	for key := range p.ConsensusKeys {
		if sync.state.IsInCommittee(key.Address()) {
			return true
		}
	}

	return false
}

// weAreInTheCommittee checks if one of the validators is a member of the committee
// at the current height.
func (sync *synchronizer) weAreInTheCommittee() bool {
	return sync.consMgr.HasActiveInstance()
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

func (sync *synchronizer) prepareBlocks(from uint32, count uint32) [][]byte {
	ourHeight := sync.state.LastBlockHeight()

	if from > ourHeight {
		sync.logger.Debug("we don't have block at this height", "height", from)
		return nil
	}

	if from+count > ourHeight {
		count = ourHeight - from + 1
	}

	blocks := make([][]byte, 0, count)

	for h := from; h < from+count; h++ {
		b := sync.state.StoredBlock(h)
		if b == nil {
			sync.logger.Warn("unable to find a block", "height", h)
			return nil
		}

		blocks = append(blocks, b.Data)
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
		sync.peerSet.IncreaseSendFailedCounter(pid)
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
