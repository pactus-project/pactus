package sync

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/cache"
	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/peerset/session"
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
	ctx         context.Context
	config      *Config
	valKeys     []*bls.ValidatorKey
	state       state.Facade
	consMgr     consensus.Manager
	peerSet     *peerset.PeerSet
	firewall    *firewall.Firewall
	cache       *cache.Cache
	handlers    map[message.Type]messageHandler
	broadcastCh <-chan message.Message
	networkCh   <-chan network.Event
	network     network.Network
	logger      *logger.SubLogger
}

func NewSynchronizer(
	conf *Config,
	valKeys []*bls.ValidatorKey,
	st state.Facade,
	consMgr consensus.Manager,
	net network.Network,
	broadcastCh <-chan message.Message,
) (Synchronizer, error) {
	sync := &synchronizer{
		ctx:         context.Background(), // TODO, set proper context
		config:      conf,
		valKeys:     valKeys,
		state:       st,
		consMgr:     consMgr,
		network:     net,
		broadcastCh: broadcastCh,
		networkCh:   net.EventChannel(),
	}

	sync.peerSet = peerset.NewPeerSet(conf.SessionTimeout)
	sync.logger = logger.NewSubLogger("_sync", sync)
	sync.firewall = firewall.NewFirewall(conf.Firewall, net, sync.peerSet, st, sync.logger)

	cacheSize := conf.CacheSize()
	ca, err := cache.NewCache(conf.CacheSize())
	if err != nil {
		return nil, err
	}
	sync.cache = ca
	sync.logger.Info("cache setup", "size", cacheSize)

	handlers := make(map[message.Type]messageHandler)

	handlers[message.TypeHello] = newHelloHandler(sync)
	handlers[message.TypeHelloAck] = newHelloAckHandler(sync)
	handlers[message.TypeTransactions] = newTransactionsHandler(sync)
	handlers[message.TypeQueryProposal] = newQueryProposalHandler(sync)
	handlers[message.TypeProposal] = newProposalHandler(sync)
	handlers[message.TypeQueryVotes] = newQueryVotesHandler(sync)
	handlers[message.TypeVote] = newVoteHandler(sync)
	handlers[message.TypeBlockAnnounce] = newBlockAnnounceHandler(sync)
	handlers[message.TypeBlocksRequest] = newBlocksRequestHandler(sync)
	handlers[message.TypeBlocksResponse] = newBlocksResponseHandler(sync)

	sync.handlers = handlers

	return sync, nil
}

func (sync *synchronizer) Start() error {
	if err := sync.network.JoinGeneralTopic(sync.shouldPropagateGeneralMessage); err != nil {
		return err
	}
	// TODO: Not joining consensus topic when we are syncing
	if err := sync.network.JoinConsensusTopic(sync.shouldPropagateConsensusMessage); err != nil {
		return err
	}

	go sync.receiveLoop()
	go sync.broadcastLoop()

	return nil
}

func (sync *synchronizer) Stop() {
	sync.ctx.Done()
}

func (sync *synchronizer) stateHeight() uint32 {
	stateHeight := sync.state.LastBlockHeight()
	return stateHeight
}

func (sync *synchronizer) moveConsensusToNewHeight() {
	stateHeight := sync.stateHeight()
	consHeight, _ := sync.consMgr.HeightRound()
	if stateHeight >= consHeight {
		sync.consMgr.MoveToNewHeight()
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

func (sync *synchronizer) sendTo(msg message.Message, to peer.ID) error {
	bdl := sync.prepareBundle(msg)
	if bdl != nil {
		data, _ := bdl.Encode()
		sync.peerSet.UpdateLastSent(to)
		sync.peerSet.IncreaseSentBytesCounter(msg.Type(), int64(len(data)), &to)

		err := sync.network.SendTo(data, to)
		if err != nil {
			return err
		}
		sync.logger.Info("sending bundle to a peer",
			"bundle", bdl, "to", to)
	}
	return nil
}

func (sync *synchronizer) shouldBroadcast(msg message.Message) bool {
	if msg.Type() == message.TypeBlockAnnounce {
		m := msg.(*message.BlockAnnounceMessage)
		if sync.cache.HasBlockInCache(m.Height()) {
			// We have received the block announcement from other peers before,
			// so we can simply ignore broadcasting it again.
			// This helps to reduce the network bandwidth.
			return false
		}
	}

	return true
}

func (sync *synchronizer) broadcast(msg message.Message) {
	if !sync.shouldBroadcast(msg) {
		return
	}
	bdl := sync.prepareBundle(msg)
	if bdl != nil {
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagBroadcasted)

		data, _ := bdl.Encode()
		err := sync.network.Broadcast(data, msg.Type().TopicID())
		if err != nil {
			sync.logger.Error("error on broadcasting bundle", "bundle", bdl, "error", err)
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

func (sync *synchronizer) sayHello(to peer.ID) error {
	msg := message.NewHelloMessage(
		sync.SelfID(),
		sync.config.Moniker,
		sync.stateHeight(),
		sync.config.Services(),
		sync.state.LastBlockHash(),
		sync.state.Genesis().Hash(),
	)
	msg.Sign(sync.valKeys)

	sync.logger.Info("sending Hello message", "to", to)
	return sync.sendTo(msg, to)
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
			switch e.Type() {
			case network.EventTypeGossip:
				ge := e.(*network.GossipMessage)
				sync.processGossipMessage(ge)

			case network.EventTypeStream:
				se := e.(*network.StreamMessage)
				sync.processStreamMessage(se)

			case network.EventTypeConnect:
				ce := e.(*network.ConnectEvent)
				sync.processConnectEvent(ce)

			case network.EventTypeDisconnect:
				de := e.(*network.DisconnectEvent)
				sync.processDisconnectEvent(de)
			}
		}
	}
}

func (sync *synchronizer) processGossipMessage(msg *network.GossipMessage) {
	bdl := sync.firewall.OpenGossipBundle(msg.Data, msg.From)
	err := sync.processIncomingBundle(bdl, msg.From)
	if err != nil {
		sync.logger.Warn("error on parsing a Gossip bundle",
			"from", msg.From, "bundle", bdl, "error", err)
		sync.peerSet.IncreaseInvalidBundlesCounter(msg.From)
	}
}

func (sync *synchronizer) processStreamMessage(msg *network.StreamMessage) {
	bdl := sync.firewall.OpenStreamBundle(msg.Reader, msg.From)
	if err := msg.Reader.Close(); err != nil {
		// TODO: write test for me
		sync.logger.Warn("error on closing stream", "error", err, "source", msg.From)
	}
	err := sync.processIncomingBundle(bdl, msg.From)
	if err != nil {
		sync.logger.Warn("error on parsing a Stream bundle",
			"source", msg.From, "bundle", bdl, "error", err)
		sync.peerSet.IncreaseInvalidBundlesCounter(msg.From)
	}
}

func (sync *synchronizer) processConnectEvent(ce *network.ConnectEvent) {
	sync.peerSet.UpdateStatus(ce.PeerID, peerset.StatusCodeConnected)
	sync.peerSet.UpdateAddress(ce.PeerID, ce.RemoteAddress)

	if err := sync.sayHello(ce.PeerID); err != nil {
		sync.logger.Warn("sending Hello message failed",
			"to", ce.PeerID, "error", err)
	}
}

func (sync *synchronizer) processDisconnectEvent(de *network.DisconnectEvent) {
	sync.peerSet.UpdateStatus(de.PeerID, peerset.StatusCodeDisconnected)
}

func (sync *synchronizer) processIncomingBundle(bdl *bundle.Bundle, from peer.ID) error {
	if bdl == nil {
		return nil
	}

	sync.logger.Info("received a bundle", "from", from, "bundle", bdl)
	h := sync.handlers[bdl.Message.Type()]
	if h == nil {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid message type: %v", bdl.Message.Type())
	}

	return h.ParseMessage(bdl.Message, from)
}

func (sync *synchronizer) String() string {
	return fmt.Sprintf("{☍ %d ⛃ %d}",
		sync.peerSet.Len(),
		sync.cache.Len())
}

// updateBlockchain checks whether the node's height is shorter than the network's height or not.
// If the node's height is shorter than the network's height by more than two hours (720 blocks),
// it should start downloading blocks from the network's nodes.
// Otherwise, the node can request the latest blocks from the network.
func (sync *synchronizer) updateBlockchain() {
	// Maybe we have some blocks inside the cache?
	_ = sync.tryCommitBlocks()

	// If we have the last block inside the cache but no certificate,
	// we need to download the next block.
	downloadHeight := sync.state.LastBlockHeight()
	downloadHeight++

	if sync.cache.HasBlockInCache(downloadHeight) {
		downloadHeight++
	}

	// Check if we have any expired sessions
	sync.peerSet.SetExpiredSessionsAsUncompleted()

	sync.peerSet.IterateSessions(func(ssn *session.Session) {
		if ssn.Status == session.Uncompleted {
			sync.logger.Debug("uncompleted block request, re-download",
				"sid", ssn.SessionID, "pid", ssn.PeerID,
				"stats", sync.peerSet.SessionStats())

			// Try to re-download the blocks from this closed session
			sync.sendBlockRequestToRandomPeer(ssn.From, ssn.Count, true)
		}
	})

	// First, let's check if we have any open sessions.
	// If there are any open sessions, we should wait for them to be closed.
	// Otherwise, we can request the same blocks from different peers.
	// TODO: write test for me
	if sync.peerSet.HasAnyOpenSession() {
		sync.logger.Debug("we have open session",
			"stats", sync.peerSet.SessionStats())

		return
	}

	sync.peerSet.RemoveAllSessions()

	blockInterval := sync.state.Params().BlockInterval()
	curTime := util.RoundNow(int(blockInterval.Seconds()))
	lastBlockTime := sync.state.LastBlockTime()
	diff := curTime.Sub(lastBlockTime)
	numOfBlocks := uint32(diff.Seconds() / blockInterval.Seconds())

	if numOfBlocks <= 1 {
		// We are sync
		return
	}

	sync.logger.Info("start syncing with the network",
		"numOfBlocks", numOfBlocks, "height", downloadHeight)
	if numOfBlocks > sync.config.LatestBlockInterval {
		sync.downloadBlocks(downloadHeight, true)
	} else {
		sync.downloadBlocks(downloadHeight, false)
	}
}

// downloadBlocks starts downloading blocks from the network.
func (sync *synchronizer) downloadBlocks(from uint32, onlyNodeNetwork bool) {
	sync.logger.Debug("downloading blocks", "from", from)

	for i := sync.peerSet.NumberOfSessions(); i < sync.config.MaxSessions; i++ {
		count := sync.config.LatestBlockInterval
		sent := sync.sendBlockRequestToRandomPeer(from, count, onlyNodeNetwork)
		if !sent {
			return
		}

		from += count
	}
}

func (sync *synchronizer) sendBlockRequestToRandomPeer(from, count uint32, onlyNodeNetwork bool) bool {
	for i := sync.peerSet.NumberOfSessions(); i < sync.config.MaxSessions; i++ {
		p := sync.peerSet.GetRandomPeer()
		if p == nil {
			break
		}

		// Don't open a new session if we already have an open session with the same peer.
		// This helps us to get blocks from different peers.
		// TODO: write test for me
		if sync.peerSet.HasOpenSession(p.PeerID) {
			continue
		}

		// We haven't completed the handshake with this peer.
		// Maybe it is a gossip peer.
		if !p.IsKnownOrTrusty() {
			continue
		}

		if onlyNodeNetwork && !p.HasNetworkService() {
			if onlyNodeNetwork {
				sync.network.CloseConnection(p.PeerID)
			}
			continue
		}

		for sync.cache.HasBlockInCache(from) {
			from++
			count--

			if count == 0 {
				// we have blocks inside the cache
				sync.logger.Debug("sending download request ignored", "from", from+1)
				return true
			}
		}

		sync.logger.Debug("sending download request", "from", from+1, "count", count, "pid", p.PeerID)
		ssn := sync.peerSet.OpenSession(p.PeerID, from, from+count-1)
		msg := message.NewBlocksRequestMessage(ssn.SessionID, from, count)
		err := sync.sendTo(msg, p.PeerID)
		if err != nil {
			sync.logger.Debug("sending download request failed",
				"from", from, "count", count, "pid", p.PeerID, "error", err)

			sync.peerSet.SetSessionUncompleted(ssn.SessionID)
		} else {
			return true
		}
	}

	sync.logger.Debug("unable to open a new session",
		"stats", sync.peerSet.SessionStats())
	return false
}

func (sync *synchronizer) tryCommitBlocks() error {
	height := sync.stateHeight() + 1
	for {
		blk := sync.cache.GetBlock(height)
		if blk == nil {
			break
		}
		cert := sync.cache.GetCertificate(height)
		if cert == nil {
			break
		}
		trxs := blk.Transactions()
		for i := 0; i < trxs.Len(); i++ {
			trx := trxs[i]
			if trx.IsPublicKeyStriped() {
				pub, err := sync.state.PublicKey(trx.Payload().Signer())
				if err != nil {
					return err
				}
				trx.SetPublicKey(pub)
			}
		}

		if err := blk.BasicCheck(); err != nil {
			return err
		}
		if err := cert.BasicCheck(); err != nil {
			return err
		}

		sync.logger.Trace("committing block", "height", height, "block", blk)
		if err := sync.state.CommitBlock(blk, cert); err != nil {
			sync.logger.Warn("committing block failed, removing block from the cache",
				"height", height, "block", blk, "error", err)

			sync.cache.RemoveBlock(height) // TODO: test me
			return err
		}
		height++
	}

	return nil
}

func (sync *synchronizer) prepareBlocks(from, count uint32) [][]byte {
	ourHeight := sync.stateHeight()

	if from > ourHeight {
		sync.logger.Debug("we don't have block at this height", "height", from)
		return nil
	}

	if from+count > ourHeight {
		count = ourHeight - from + 1
	}

	blocks := make([][]byte, 0, count)

	for h := from; h < from+count; h++ {
		b := sync.state.CommittedBlock(h)
		if b == nil {
			sync.logger.Warn("unable to find a block", "height", h)
			return nil
		}

		blocks = append(blocks, b.Data)
	}

	return blocks
}

// weAreInTheCommittee checks if one of the validators is a member of the committee
// at the current height.
func (sync *synchronizer) weAreInTheCommittee() bool {
	return sync.consMgr.HasActiveInstance()
}

func (sync *synchronizer) shouldPropagateMessage(msg *network.GossipMessage) bool {
	// Propagate our messages
	if msg.From == sync.SelfID() {
		return true
	}

	if sync.config.NodeGossip {
		return true
	}

	return sync.weAreInTheCommittee()
}

func (sync *synchronizer) shouldPropagateGeneralMessage(msg *network.GossipMessage) bool {
	return sync.shouldPropagateMessage(msg)
}

func (sync *synchronizer) shouldPropagateConsensusMessage(msg *network.GossipMessage) bool {
	return sync.shouldPropagateMessage(msg)
}
