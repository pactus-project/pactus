package sync

import (
	"context"
	"fmt"
	"io"
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
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

// IMPORTANT NOTE
//
// Sync module is based on pulling, not pushing
// Means if a node is shorter than the network, we don't send him anything
// The node should request (pull) itself.

const FlagInitialBlockDownload = 0x1

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
	handlers        map[payload.Type]payloadHandler
	broadcastCh     <-chan payload.Payload
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

	handlers := make(map[payload.Type]payloadHandler)

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
	handlers[payload.PayloadTypeBlocksRequest] = newBlocksRequestHandler(sync)
	handlers[payload.PayloadTypeBlocksResponse] = newBlocksResponseHandler(sync)

	sync.handlers = handlers

	return sync, nil
}

func (sync *synchronizer) Start() error {
	sync.network.SetCallback(sync.onReceiveData)
	if err := sync.network.JoinGeneralTopic(); err != nil {
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

func (sync *synchronizer) onStartingTimeout() {
	ourHeight := sync.state.LastBlockHeight()
	networkHeight := sync.peerSet.MaxClaimedHeight()

	if ourHeight >= networkHeight-1 {
		sync.synced()
	}
}

func (sync *synchronizer) synced() {
	sync.logger.Debug("We are synced", "height", sync.state.LastBlockHeight())
	if err := sync.network.JoinConsensusTopic(); err != nil {
		sync.logger.Error("Error on joining consensus topic", "err", err)
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
		sync.signer.SignData(sync.signer.PublicKey().RawBytes()),
		sync.state.LastBlockHeight(),
		flags, sync.state.GenesisHash())

	sync.broadcast(pld)
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

// updateBlokchain checks if the node height is shorter than the network or not.
// If the node height is shorter than network more than two hours (720 blocks),
// it should join the download topic and start downloading the blocks,
// otherwise the node can request the latest blocks from the network.
func (sync *synchronizer) updateBlokchain() {
	// TODO: write test for me
	if sync.peerSet.HasAnyOpenSession() {
		sync.logger.Debug("We have open seasson")
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

		sync.logger.Info("We are not synced with the network.")
		if claimedHeight > ourHeight+LatestBlockInterval {
			sync.downloadBlocks(from)
		} else {
			sync.queryLatestBlocks(from)
		}
	}
}

func (sync *synchronizer) onReceiveData(reader io.Reader, source peer.ID, from peer.ID) {
	msg := sync.firewall.OpenMessage(reader, source, from)
	if msg == nil {
		return
	}

	sync.logger.Debug("Received a message", "from", util.FingerprintPeerID(from), "message", msg)
	handler := sync.handlers[msg.Payload.Type()]
	if handler == nil {
		sync.logger.Error("Invalid payload type: %v", msg.Payload.Type())
		return
	}

	if err := handler.ParsPayload(msg.Payload, msg.Initiator); err != nil {
		peer := sync.peerSet.MustGetPeer(from)
		peer.IncreaseInvalidMessage()
		sync.logger.Warn("Error on parsing a message", "from", util.FingerprintPeerID(from), "message", msg, "err", err)
		return
	}
}

func (sync *synchronizer) prepareMessage(pld payload.Payload) *message.Message {
	handler := sync.handlers[pld.Type()]
	if handler == nil {
		sync.logger.Warn("Invalid payload type: %v", pld.Type())
		return nil
	}
	msg := handler.PrepareMessage(pld)
	if msg != nil {
		if err := msg.SanityCheck(); err != nil {
			sync.logger.Error("Broadcasting an invalid message", "message", msg, "err", err)
			return nil
		}
		return msg
	}
	return nil
}

func (sync *synchronizer) sendTo(pld payload.Payload, to peer.ID) {
	msg := sync.prepareMessage(pld)
	if msg != nil {
		data, _ := msg.Encode()
		err := sync.network.SendTo(data, to)
		if err != nil {
			sync.logger.Error("Error on sending message", "message", msg, "err", err)
		} else {
			sync.logger.Debug("Sending message to a peer", "message", msg, "to", to)
		}
	}
}

func (sync *synchronizer) broadcast(pld payload.Payload) {
	msg := sync.prepareMessage(pld)
	if msg != nil {
		msg.Flags = util.SetFlag(msg.Flags, message.FlagBroadcasted)
		data, _ := msg.Encode()
		err := sync.network.Broadcast(data, pld.Type().TopicID())
		if err != nil {
			sync.logger.Error("Error on broadcasting message", "message", msg, "err", err)
		} else {
			sync.logger.Debug("Broadcasting new message", "message", msg)
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
			sync.logger.Debug("We reached maximum number of open sessions")
			break
		}

		// TODO: write test for me
		if !peer.InitialBlockDownload() {
			continue
		}

		// TODO: write test for me
		if peer.Status() != peerset.StatusCodeOK {
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

		sync.logger.Debug("Sending download request", "from", from+1, "to", to, "pid", util.FingerprintPeerID(peer.PeerID()))
		session := sync.peerSet.OpenSession(peer.PeerID())
		pld := payload.NewBlocksRequestPayload(session.SessionID(), from+1, to)
		sync.sendTo(pld, peer.PeerID())
		from = to
	}
}

func (sync *synchronizer) queryLatestBlocks(from int) {
	randPeer := sync.peerSet.GetRandomPeer()

	// TODO: write test for me
	if randPeer == nil || randPeer.Status() != peerset.StatusCodeOK {
		return
	}

	// TODO: write test for me
	to := randPeer.Height()
	if randPeer.Height() < from+1 {
		return
	}

	sync.logger.Debug("Querying the latest blocks", "from", from+1, "to", to, "pid", util.FingerprintPeerID(randPeer.PeerID()))
	session := sync.peerSet.OpenSession(randPeer.PeerID())
	pld := payload.NewBlocksRequestPayload(session.SessionID(), from+1, to)
	sync.sendTo(pld, randPeer.PeerID())
}

// peerIsInTheCommittee checks if the peer is a member of committee
func (sync *synchronizer) peerIsInTheCommittee(id peer.ID) bool {
	p := sync.peerSet.GetPeer(id)
	if p == nil {
		return false
	}

	if !p.HasPublicKey() {
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

func (sync *synchronizer) updateSession(sessionID int, pid peer.ID, code payload.ResponseCode) {
	s := sync.peerSet.FindSession(sessionID)
	if s == nil {
		sync.logger.Debug("Session not found or closed", "session-id", sessionID)
		return
	}

	if s.PeerID() != pid {
		sync.logger.Debug("Peer ID is not known", "session-id", sessionID, "pid", pid)
		return
	}

	s.SetLastResponseCode(code)

	switch code {
	case payload.ResponseCodeRejected:
		sync.logger.Debug("session rejected, close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.updateBlokchain()

	case payload.ResponseCodeBusy:
		sync.logger.Debug("Peer is busy. close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.updateBlokchain()

	case payload.ResponseCodeMoreBlocks:
		sync.logger.Debug("Peer responding us. keep session open", "session-id", sessionID)

	case payload.ResponseCodeNoMoreBlocks:
		sync.logger.Debug("Peer has no more block. close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.updateBlokchain()

	case payload.ResponseCodeSynced:
		sync.logger.Debug("Peer infomed us we are synced. close session", "session-id", sessionID)
		sync.peerSet.CloseSession(sessionID)
		sync.synced()
	}

}
