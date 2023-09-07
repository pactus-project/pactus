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
	"github.com/pactus-project/pactus/sync/services"
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
	broadcastCh <-chan message.Message,
) (Synchronizer, error) {
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
	if err := sync.network.JoinGeneralTopic(); err != nil {
		return err
	}
	// TODO: Not joining consensus topic when we are syncing
	if err := sync.network.JoinConsensusTopic(); err != nil {
		return err
	}

	go sync.receiveLoop()
	go sync.broadcastLoop()

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

func (sync *synchronizer) sayHello(to peer.ID) error {
	nodeServices := []int{}
	if sync.config.NodeNetwork {
		nodeServices = append(nodeServices, services.Network)
	}

	msg := message.NewHelloMessage(
		sync.SelfID(),
		sync.config.Moniker,
		sync.state.LastBlockHeight(),
		services.New(nodeServices...),
		sync.state.LastBlockHash(),
		sync.state.Genesis().Hash(),
	)
	msg.Sign(sync.signers...)

	sync.peerSet.UpdateStatus(to, peerset.StatusCodeConnected)

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
				bdl := sync.firewall.OpenGossipBundle(ge.Data, ge.Source, ge.From)
				err := sync.processIncomingBundle(bdl)
				if err != nil {
					sync.logger.Warn("error on parsing a Gossip bundle", "initiator", bdl.Initiator, "bundle", bdl, "err", err)
					sync.peerSet.IncreaseInvalidBundlesCounter(bdl.Initiator)
				}

			case network.EventTypeStream:
				se := e.(*network.StreamMessage)
				bdl := sync.firewall.OpenStreamBundle(se.Reader, se.Source)
				if err := se.Reader.Close(); err != nil {
					// TODO: write test for me
					sync.logger.Warn("error on closing stream", "err", err)
				}
				err := sync.processIncomingBundle(bdl)
				if err != nil {
					sync.logger.Warn("error on parsing a Stream bundle", "initiator", bdl.Initiator, "bundle", bdl, "err", err)
					sync.peerSet.IncreaseInvalidBundlesCounter(bdl.Initiator)
				}
			case network.EventTypeConnect:
				ce := e.(*network.ConnectEvent)
				if err := sync.sayHello(ce.PeerID); err != nil {
					sync.logger.Warn("sending Hello message failed", "to", ce.PeerID, "error", err)
				}
			case network.EventTypeDisconnect:
				de := e.(*network.DisconnectEvent)
				sync.peerSet.UpdateStatus(de.PeerID, peerset.StatusCodeDisconnected)
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
	return fmt.Sprintf("{☍ %d ⛃ %d}",
		sync.peerSet.Len(),
		sync.cache.Len())
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

	blockInterval := sync.state.Params().BlockInterval()
	curTime := util.RoundNow(int(blockInterval.Seconds()))
	lastBlockTime := sync.state.LastBlockTime()
	diff := curTime.Sub(lastBlockTime)
	numOfBlocks := uint32(diff.Seconds() / blockInterval.Seconds())

	if numOfBlocks <= 1 {
		// We are sync
		return
	}

	// Make sure we have committed the latest blocks inside the cache
	LastBlockHeight := sync.state.LastBlockHeight()
	for sync.cache.HasBlockInCache(LastBlockHeight + 1) {
		LastBlockHeight++
	}

	sync.logger.Info("start syncing with the network", "numOfBlocks", numOfBlocks)
	if numOfBlocks > LatestBlockInterval {
		sync.downloadBlocks(LastBlockHeight, true)
	} else {
		sync.downloadBlocks(LastBlockHeight, false)
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
			sync.logger.Warn("error on sending bundle", "bundle", bdl, "err", err, "to", to)

			return err
		}
		sync.logger.Info("sending bundle to a peer", "bundle", bdl, "to", to)
	}
	return nil
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

	sync.peerSet.IteratePeers(func(p *peerset.Peer) {
		// Don't open a new session if we already have an open session with the same peer.
		// This helps us to get blocks from different peers.
		// TODO: write test for me
		if sync.peerSet.HasOpenSession(p.PeerID) {
			return
		}

		if !p.IsKnownOrTrusty() {
			return
		}

		if onlyNodeNetwork && !p.HasNetworkService() {
			if onlyNodeNetwork {
				sync.network.CloseConnection(p.PeerID)
			}
			return
		}

		count := LatestBlockInterval
		sync.logger.Debug("sending download request", "from", from+1, "count", count, "pid", p.PeerID)
		session := sync.peerSet.OpenSession(p.PeerID)
		msg := message.NewBlocksRequestMessage(session.SessionID(), from+1, count)
		err := sync.sendTo(msg, p.PeerID)
		if err != nil {
			sync.peerSet.CloseSession(session.SessionID())
		} else {
			from += count
		}
	})
}

// peerIsInTheCommittee checks if the peer is a member of the committee
// at the current height.
func (sync *synchronizer) peerIsInTheCommittee(pid peer.ID) bool {
	p := sync.peerSet.GetPeer(pid)
	if !p.IsKnownOrTrusty() {
		return false
	}

	for _, key := range p.ConsensusKeys {
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
		b := sync.state.CommittedBlock(h)
		if b == nil {
			sync.logger.Warn("unable to find a block", "height", h)
			return nil
		}

		blocks = append(blocks, b.Data)
	}

	return blocks
}
