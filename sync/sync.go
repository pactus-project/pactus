package sync

import (
	"context"
	"fmt"
	"time"

	"github.com/ezex-io/gopkg/pipeline"
	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	consmgr "github.com/pactus-project/pactus/consensus/manager"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/cache"
	"github.com/pactus-project/pactus/sync/firewall"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
	"github.com/pactus-project/pactus/sync/peerset/session"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/ntp"
)

// IMPORTANT NOTES:
//
// 1. The Sync module is based on pulling instead of pushing. This means that the network
// does not update a node (push); instead, a node should update itself (pull).
//
// 2. The Synchronizer should not have any locks to prevent deadlocks. All submodules,
// such as state or consensus, should be thread-safe.

type synchronizer struct {
	config        *Config
	valKeys       []*bls.ValidatorKey
	state         state.Facade
	consV1Mgr     consmgr.Manager
	consV2Mgr     consmgr.Manager
	peerSet       *peerset.PeerSet
	firewall      *firewall.Firewall
	cache         *cache.Cache
	handlers      map[message.Type]messageHandler
	broadcastPipe pipeline.Pipeline[message.Message]
	networkPipe   pipeline.Pipeline[network.Event]
	network       network.Network
	logger        *logger.SubLogger
	ntp           *ntp.Checker
}

//nolint:revive // arguments can't be reduced.
func NewSynchronizer(
	ctx context.Context,
	conf *Config,
	valKeys []*bls.ValidatorKey,
	state state.Facade,
	consV1Mgr consmgr.Manager,
	consV2Mgr consmgr.Manager,
	network network.Network,
	broadcastPipe pipeline.Pipeline[message.Message],
	networkPipe pipeline.Pipeline[network.Event],
) (Synchronizer, error) {
	sync := &synchronizer{
		config:        conf,
		valKeys:       valKeys,
		state:         state,
		consV1Mgr:     consV1Mgr,
		consV2Mgr:     consV2Mgr,
		network:       network,
		broadcastPipe: broadcastPipe,
		networkPipe:   networkPipe,
		ntp:           ntp.NewNtpChecker(ctx),
	}

	sync.peerSet = peerset.NewPeerSet(conf.SessionTimeout())
	sync.logger = logger.NewSubLogger("_sync", sync)
	fw, err := firewall.NewFirewall(conf.Firewall, network, sync.peerSet, state)
	if err != nil {
		return nil, err
	}

	sync.firewall = fw

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
	handlers[message.TypeTransaction] = newTransactionsHandler(sync)
	handlers[message.TypeQueryProposal] = newQueryProposalHandler(sync)
	handlers[message.TypeProposal] = newProposalHandler(sync)
	handlers[message.TypeQueryVote] = newQueryVoteHandler(sync)
	handlers[message.TypeVote] = newVoteHandler(sync)
	handlers[message.TypeBlockAnnounce] = newBlockAnnounceHandler(sync)
	handlers[message.TypeBlocksRequest] = newBlocksRequestHandler(sync)
	handlers[message.TypeBlocksResponse] = newBlocksResponseHandler(sync)

	sync.handlers = handlers

	sync.networkPipe.RegisterReceiver(sync.processNetworkEvent)
	sync.broadcastPipe.RegisterReceiver(sync.broadcastMessage)

	return sync, nil
}

func (sync *synchronizer) Start() error {
	if err := sync.network.JoinTopic(network.TopicIDBlock, sync.blockTopicEvaluator); err != nil {
		return err
	}
	if err := sync.network.JoinTopic(network.TopicIDTransaction, sync.transactionTopicEvaluator); err != nil {
		return err
	}
	// TODO: Not joining consensus topic when we are syncing
	if err := sync.network.JoinTopic(network.TopicIDConsensus, sync.consensusTopicEvaluator); err != nil {
		return err
	}

	sync.ntp.Start()

	return nil
}

func (sync *synchronizer) Stop() {
	sync.ntp.Stop()
}

func (sync *synchronizer) ClockOffset() (time.Duration, error) {
	return sync.ntp.ClockOffset()
}

func (sync *synchronizer) IsClockOutOfSync() bool {
	return sync.ntp.IsOutOfSync()
}

func (sync *synchronizer) stateHeight() uint32 {
	stateHeight := sync.state.LastBlockHeight()

	return stateHeight
}

func (sync *synchronizer) moveConsensusToNewHeight() {
	stateHeight := sync.stateHeight()
	consHeight, _ := sync.getConsMgr().HeightRound()
	if stateHeight >= consHeight {
		sync.getConsMgr().MoveToNewHeight()
	}
}

func (sync *synchronizer) prepareBundle(msg message.Message) *bundle.Bundle {
	h := sync.handlers[msg.Type()]
	bdl := h.PrepareBundle(msg)

	// Bundles will be carried through LibP2P.
	// In future we might support other libraries.
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagCarrierLibP2P)

	switch sync.state.Genesis().ChainType() {
	case genesis.Mainnet:
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkMainnet)
	case genesis.Testnet:
		bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagNetworkTestnet)
	case genesis.Localnet:
		// It's localnet and for testing purpose only
	}

	return bdl
}

func (sync *synchronizer) sendTo(msg message.Message, pid peer.ID) {
	bdl := sync.prepareBundle(msg)
	data, _ := bdl.Encode()

	sync.network.SendTo(data, pid)
	sync.peerSet.UpdateLastSent(pid)
	sync.peerSet.UpdateSentMetric(&pid, msg.Type(), int64(len(data)))

	sync.logger.Debug("bundle sent", "bundle", bdl, "pid", pid)
}

func (sync *synchronizer) broadcast(msg message.Message) {
	if msg.Type() == message.TypeBlockAnnounce {
		m := msg.(*message.BlockAnnounceMessage)
		if sync.cache.HasBlockInCache(m.Height()) {
			// We have received the block announcement from other peers before,
			// so we can simply ignore broadcasting it again.
			// This helps to reduce the network bandwidth.
			return
		}
	}

	bdl := sync.prepareBundle(msg)
	bdl.Flags = util.SetFlag(bdl.Flags, bundle.BundleFlagBroadcasted)

	data, _ := bdl.Encode()
	sync.network.Broadcast(data, msg.TopicID())
	sync.peerSet.UpdateSentMetric(nil, msg.Type(), int64(len(data)))

	sync.logger.Debug("bundle broadcasted", "bundle", bdl)
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

func (sync *synchronizer) Services() service.Services {
	return sync.config.Services
}

func (sync *synchronizer) sayHello(pid peer.ID) {
	peer := sync.peerSet.GetPeer(pid)
	if peer.Status.IsKnown() {
		return
	}

	msg := message.NewHelloMessage(
		sync.SelfID(),
		sync.config.Moniker,
		sync.config.Services,
		sync.stateHeight(),
		sync.state.LastBlockHash(),
		sync.state.Genesis().Hash(),
	)
	msg.Sign(sync.valKeys)

	sync.sendTo(msg, pid)
}

func (sync *synchronizer) broadcastMessage(msg message.Message) {
	sync.broadcast(msg)
}

func (sync *synchronizer) processNetworkEvent(evt network.Event) {
	switch evt.Type() {
	case network.EventTypeGossip:
		ge := evt.(*network.GossipMessage)
		sync.processGossipMessage(ge)

	case network.EventTypeStream:
		se := evt.(*network.StreamMessage)
		sync.processStreamMessage(se)

	case network.EventTypeConnect:
		ce := evt.(*network.ConnectEvent)
		sync.processConnectEvent(ce)

	case network.EventTypeDisconnect:
		de := evt.(*network.DisconnectEvent)
		sync.processDisconnectEvent(de)

	case network.EventTypeProtocols:
		pe := evt.(*network.ProtocolsEvents)
		sync.processProtocolsEvent(pe)
	}
}

func (sync *synchronizer) processGossipMessage(msg *network.GossipMessage) {
	sync.logger.Debug("processing gossip message", "pid", msg.From)

	bdl, err := sync.firewall.OpenGossipBundle(msg.Data, msg.From)
	if err != nil {
		sync.logger.Debug("error on parsing a Gossip bundle",
			"from", msg.From, "bundle", bdl, "error", err)

		return
	}
	sync.processIncomingBundle(bdl, msg.From)
}

func (sync *synchronizer) processStreamMessage(msg *network.StreamMessage) {
	sync.logger.Debug("processing stream message", "pid", msg.From)

	defer func() {
		_ = msg.Reader.Close()
	}()

	bdl, err := sync.firewall.OpenStreamBundle(msg.Reader, msg.From)
	if err != nil {
		sync.logger.Debug("error on parsing a Stream bundle",
			"from", msg.From, "bundle", bdl, "error", err)

		return
	}

	sync.processIncomingBundle(bdl, msg.From)
}

func (sync *synchronizer) processConnectEvent(eve *network.ConnectEvent) {
	sync.logger.Debug("processing connect event", "pid", eve.PeerID)

	sync.peerSet.UpdateAddress(eve.PeerID, eve.RemoteAddress, eve.Direction)
	sync.peerSet.UpdateStatus(eve.PeerID, status.StatusConnected)
}

func (sync *synchronizer) processProtocolsEvent(eve *network.ProtocolsEvents) {
	sync.logger.Debug("processing protocols event", "pid", eve.PeerID, "protocols", eve.Protocols)

	sync.peerSet.UpdateProtocols(eve.PeerID, eve.Protocols)

	peer := sync.peerSet.GetPeer(eve.PeerID)
	if peer.Direction == lp2pnetwork.DirOutbound {
		sync.logger.Info("sending Hello message (outbound)", "to", eve.PeerID)
		sync.sayHello(eve.PeerID)

		// Mark that we've sent the hello message to the inbound peer
		sync.peerSet.UpdateOutboundHelloSent(eve.PeerID, true)
	}
}

func (sync *synchronizer) processDisconnectEvent(de *network.DisconnectEvent) {
	sync.logger.Debug("processing disconnect event", "pid", de.PeerID)

	sync.peerSet.UpdateStatus(de.PeerID, status.StatusDisconnected)
}

func (sync *synchronizer) processIncomingBundle(bdl *bundle.Bundle, from peer.ID) {
	sync.logger.Debug("received a bundle", "from", from, "bundle", bdl)
	handler := sync.handlers[bdl.Message.Type()]
	if handler == nil {
		sync.logger.Error("invalid message type", "type", bdl.Message.Type())

		return
	}

	handler.ParseMessage(bdl.Message, from)
}

// LogString returns a concise string representation intended for use in logs.
func (sync *synchronizer) LogString() string {
	return fmt.Sprintf("{☍ %d ⛃ %d}",
		sync.peerSet.Len(),
		sync.cache.Len())
}

// updateBlockchain checks whether the node's height is shorter than the network's height or not.
// If the node's height is shorter than the network's height by more than two hours (720 blocks),
// it should start downloading blocks from the network's nodes.
// Otherwise, the node can request the latest blocks from any nodes.
func (sync *synchronizer) updateBlockchain() {
	// Maybe we have some blocks inside the cache?
	sync.tryCommitBlocks()

	// Check if we have any expired sessions
	sync.peerSet.SetExpiredSessionsAsUncompleted()

	// Try to re-download the blocks for uncompleted sessions
	sessions := sync.peerSet.Sessions()
	for _, ssn := range sessions {
		if ssn.Status == session.Uncompleted {
			sync.logger.Info("uncompleted block request, re-download",
				"sid", ssn.SessionID, "pid", ssn.PeerID,
				"stats", sync.peerSet.SessionStats())

			sent := sync.sendBlockRequestToRandomPeer(ssn.From, ssn.Count, false)
			if !sent {
				break
			}
		}
	}

	// Check if there are any open sessions.
	// If open sessions exist, we should wait for them to close.
	// Otherwise, we might request to download the same blocks from different peers.
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

	downloadHeight := sync.state.LastBlockHeight()
	downloadHeight++

	if sync.cache.HasBlockInCache(downloadHeight) {
		// The last block exists inside the cache, without the certificate.
		// Ignore downloading this block again.
		downloadHeight++
	}

	sync.logger.Info("start syncing with the network",
		"numOfBlocks", numOfBlocks, "height", downloadHeight)

	if numOfBlocks > sync.config.PruneWindow {
		// Don't have blocks for more than 10 days
		sync.downloadBlocks(downloadHeight, true)
	} else {
		sync.downloadBlocks(downloadHeight, false)
	}
}

// downloadBlocks starts downloading blocks from the network.
func (sync *synchronizer) downloadBlocks(from uint32, onlyFullNodes bool) {
	sync.logger.Debug("downloading blocks", "from", from)

	for i := sync.peerSet.NumberOfSessions(); i < sync.config.MaxSessions; i++ {
		count := sync.config.BlockPerSession
		sent := sync.sendBlockRequestToRandomPeer(from, count, onlyFullNodes)
		if !sent {
			return
		}

		from += count
	}
}

func (sync *synchronizer) sendBlockRequestToRandomPeer(from, count uint32, onlyFullNodes bool) bool {
	// Prevent downloading blocks that might be cached before
	for sync.cache.HasBlockInCache(from) {
		from++
		count--

		if count == 0 {
			// we have blocks inside the cache
			sync.logger.Debug("sending download request ignored", "from", from+1)

			return true
		}
	}

	for i := sync.peerSet.NumberOfSessions(); i < sync.config.MaxSessions; i++ {
		peer := sync.peerSet.GetRandomPeer()
		if peer == nil {
			break
		}

		// Don't open a new session if we already have an open session with the same peer.
		// This helps us to get blocks from different peers.
		if sync.peerSet.HasOpenSession(peer.PeerID) {
			continue
		}

		// We haven't completed the handshake with this peer.
		if !peer.Status.IsKnown() {
			if onlyFullNodes {
				sync.network.CloseConnection(peer.PeerID)
			}

			continue
		}

		if onlyFullNodes && !peer.IsFullNode() {
			if onlyFullNodes {
				sync.network.CloseConnection(peer.PeerID)
			}

			continue
		}

		sid := sync.peerSet.OpenSession(peer.PeerID, from, count)
		msg := message.NewBlocksRequestMessage(sid, from, count)
		sync.sendTo(msg, peer.PeerID)

		sync.logger.Info("blocks request sent",
			"from", from+1, "count", count, "pid", peer.PeerID, "sid", sid)

		return true
	}

	sync.logger.Warn("unable to open a new session",
		"stats", sync.peerSet.SessionStats())

	return false
}

func (sync *synchronizer) tryCommitBlocks() {
	onError := func(height uint32, err error) {
		sync.logger.Warn("committing block failed, removing block from the cache",
			"height", height, "error", err)

		sync.cache.RemoveBlock(height)
	}

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
					onError(height, err)

					return
				}
				trx.SetPublicKey(pub)
			}
		}

		if err := blk.BasicCheck(); err != nil {
			onError(height, err)

			return
		}
		if err := cert.BasicCheck(); err != nil {
			onError(height, err)

			return
		}

		sync.logger.Trace("committing block", "height", height, "block", blk)
		if err := sync.state.CommitBlock(blk, cert); err != nil {
			onError(height, err)

			return
		}
		height++
	}
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

	for height := from; height < from+count; height++ {
		cBlk, err := sync.state.CommittedBlock(height)
		if err != nil {
			sync.logger.Warn("unable to find a block", "height", height)

			return nil
		}

		blocks = append(blocks, cBlk.Data)
	}

	return blocks
}

func (sync *synchronizer) blockTopicEvaluator(msg *network.GossipMessage) network.PropagationPolicy {
	return sync.firewall.AllowBlockRequest(msg)
}

func (sync *synchronizer) transactionTopicEvaluator(msg *network.GossipMessage) network.PropagationPolicy {
	return sync.firewall.AllowTransactionRequest(msg)
}

func (sync *synchronizer) consensusTopicEvaluator(msg *network.GossipMessage) network.PropagationPolicy {
	return sync.firewall.AllowConsensusRequest(msg)
}

// getConsMgr returns consensus manager based on the upgrade condition.
// After the chain is fully upgraded, we can remove this function.
func (sync *synchronizer) getConsMgr() consmgr.Manager {
	if sync.consV1Mgr.IsDeprecated() {
		return sync.consV2Mgr
	}

	return sync.consV1Mgr
}
