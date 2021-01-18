package sync

import (
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/txpool"
	"github.com/zarbchain/zarb-go/util"
)

var LatestBlockInterval = 720
var DownloadBlockInterval = 1000

type StateSync struct {
	config    *Config
	selfID    peer.ID
	cache     *cache.Cache
	state     state.State
	txPool    txpool.TxPool
	peerSet   *peerset.PeerSet
	logger    *logger.Logger
	publishFn PublishMessageFn
	syncedFN  SyncedCallbackFn
}

func NewStateSync(
	conf *Config,
	selfID peer.ID,
	cache *cache.Cache,
	state state.State,
	txPool txpool.TxPool,
	peerSet *peerset.PeerSet,
	logger *logger.Logger,
	publishFn PublishMessageFn,
	syncedFN SyncedCallbackFn) *StateSync {
	return &StateSync{
		config:    conf,
		selfID:    selfID,
		cache:     cache,
		state:     state,
		txPool:    txPool,
		peerSet:   peerSet,
		logger:    logger,
		publishFn: publishFn,
		syncedFN:  syncedFN,
	}
}

func (ss *StateSync) BroadcastLatestBlocksRequest(target peer.ID, from int) {
	s := ss.openNewSession(target)
	msg := message.NewLatestBlocksRequestMessage(ss.selfID, target, s.SessionID, from)
	ss.publishFn(msg)
}

func (ss *StateSync) BroadcastLatestBlocksResponse(code payload.ResponseCode, target peer.ID, sessionID, from int, blocks []*block.Block, trxs []*tx.Tx, commit *block.Commit) {
	msg := message.NewLatestBlocksResponseMessage(code, ss.selfID, target, sessionID, from, blocks, trxs, commit)
	msg.CompressIt()
	ss.publishFn(msg)
}

func (ss *StateSync) BroadcastTransactions(txs []*tx.Tx) {
	msg := message.NewTransactionsMessage(txs)
	ss.publishFn(msg)
}

func (ss *StateSync) BroadcastDownloadRequest(target peer.ID, from, to int) {
	s := ss.openNewSession(target)
	msg := message.NewDownloadRequestMessage(ss.selfID, target, s.SessionID, from, to)
	msg.CompressIt()
	ss.publishFn(msg)
}

func (ss *StateSync) BroadcastDownloadResponse(code payload.ResponseCode, target peer.ID, sessionID, from int, blocks []*block.Block, trxs []*tx.Tx) {
	msg := message.NewDownloadResponseMessage(code, ss.selfID, target, sessionID, from, blocks, trxs)
	msg.CompressIt()
	ss.publishFn(msg)
}

func (ss *StateSync) openNewSession(target peer.ID) *peerset.Session {
	s := ss.peerSet.OpenSession(target)
	ss.logger.Debug("Open new session", "session", s)
	return s
}
func (ss *StateSync) updateSession(code payload.ResponseCode, sessionID int, initiator peer.ID, target peer.ID) {
	if target != ss.selfID {
		return
	}

	s := ss.peerSet.FindSession(sessionID)
	if s == nil {
		ss.logger.Debug("Session not found or closed before")
	} else {
		s.LastResponseCode = code
		s.LastActivityAt = util.Now()
	}

	switch code {
	case payload.ResponseCodeRejected:
		ss.logger.Debug("session rejected, close session")
		ss.peerSet.CloseSession(sessionID)

	case payload.ResponseCodeBusy:
		// TODO: handle this situation
		ss.logger.Debug("Peer is busy. close session")
		ss.peerSet.CloseSession(sessionID)

	case payload.ResponseCodeMoreBlocks:
		ss.logger.Debug("Peer responding us. keep session open")

	case payload.ResponseCodeNoMoreBlocks:
		ss.logger.Debug("Peer has no more block. close session")
		ss.peerSet.CloseSession(sessionID)

	case payload.ResponseCodeSynced:
		ss.logger.Debug("Peer infomed us we are synced. close session")
		ss.peerSet.CloseSession(sessionID)
		ss.syncedFN()
	}
}

func (ss *StateSync) ProcessLatestBlocksRequestPayload(pld *payload.LatestBlocksRequestPayload) {
	ss.logger.Trace("Process latest blocks request payload", "pld", pld)

	peer := ss.peerSet.MustGetPeer(pld.Initiator)
	peer.UpdateHeight(pld.From)

	if pld.Target != ss.selfID {
		return
	}

	ourHeight := ss.state.LastBlockHeight()
	if pld.From < ourHeight-LatestBlockInterval {
		// TODO: Mark peer as bad peer
		ss.logger.Warn("The request height is not acceptable", "from", pld.From)
		return
	}

	from := pld.From
	count := ss.config.BlockPerMessage

	// Help peer to catch up
	for {
		blocks, trxs := ss.prepareBlocksAndTransactions(from, count)
		if len(blocks) == 0 {
			break
		}

		ss.BroadcastLatestBlocksResponse(payload.ResponseCodeMoreBlocks, pld.Initiator, pld.SessionID, from, blocks, trxs, nil)
		from += len(blocks)
	}

	if from < ourHeight {
		// TODO: write a test
		// Should not happen...
		ss.BroadcastLatestBlocksResponse(payload.ResponseCodeSynced, pld.Initiator, pld.SessionID, 0, nil, nil, nil)
	} else {
		lastCommit := ss.state.LastCommit()
		ss.BroadcastLatestBlocksResponse(payload.ResponseCodeSynced, pld.Initiator, pld.SessionID, 0, nil, nil, lastCommit)
	}
}

func (ss *StateSync) ProcessDownloadRequestPayload(pld *payload.DownloadRequestPayload) {
	ss.logger.Trace("Process download request payload", "pld", pld)

	peer := ss.peerSet.MustGetPeer(pld.Initiator)
	peer.UpdateHeight(pld.From)

	if pld.Target != ss.selfID {
		return
	}

	if pld.To-pld.From > DownloadBlockInterval {
		// TODO: Mark peer as bad peer
		return
	}

	ourHeight := ss.state.LastBlockHeight()
	from := pld.From
	count := ss.config.BlockPerMessage

	for {
		if from+count >= pld.To {
			// Last packet has one extra block, for confirming last block
			count++
		}
		blocks, trxs := ss.prepareBlocksAndTransactions(from, count)
		if len(blocks) == 0 {
			break
		}

		ss.BroadcastDownloadResponse(payload.ResponseCodeMoreBlocks, pld.Initiator, pld.SessionID, from, blocks, trxs)
		from += len(blocks)
		if from >= pld.To {
			break
		}
	}

	if from < ourHeight {
		ss.BroadcastDownloadResponse(payload.ResponseCodeNoMoreBlocks, pld.Initiator, pld.SessionID, 0, nil, nil)
	} else {
		ss.BroadcastDownloadResponse(payload.ResponseCodeSynced, pld.Initiator, pld.SessionID, 0, nil, nil)
	}
}

func (ss *StateSync) ProcessBlockAnnouncePayload(pld *payload.BlockAnnouncePayload) {
	ss.logger.Trace("Process block announce payload", "pld", pld)

	ss.cache.AddCommit(pld.Block.Hash(), pld.Commit)
	ss.cache.AddBlock(pld.Height, pld.Block)
	ss.tryCommitBlocks()
}

func (ss *StateSync) ProcessLatestBlocksResponsePayload(pld *payload.LatestBlocksResponsePayload) {
	ss.logger.Trace("Process latest blocks payload", "pld", pld)

	ourHeight := ss.state.LastBlockHeight()
	if pld.To() == 0 || ourHeight < pld.To() {
		if pld.LastCommit != nil {
			ss.cache.AddCommit(
				pld.LastCommit.BlockHash(),
				pld.LastCommit)
		}

		ss.addBlocksToCache(pld.From, pld.Blocks)
		ss.addTransactionsToCache(pld.Transactions)
		ss.tryCommitBlocks()
	}
	ss.updateSession(pld.ResponseCode, pld.SessionID, pld.Initiator, pld.Target)

}

func (ss *StateSync) ProcessDownloadResponsePayload(pld *payload.DownloadResponsePayload) {
	ss.logger.Trace("Process download response payload", "pld", pld)

	ourHeight := ss.state.LastBlockHeight()
	if pld.To() == 0 || ourHeight < pld.To() {
		ss.addBlocksToCache(pld.From, pld.Blocks)
		ss.addTransactionsToCache(pld.Transactions)
		ss.tryCommitBlocks()
	}
	ss.updateSession(pld.ResponseCode, pld.SessionID, pld.Initiator, pld.Target)
}

func (ss *StateSync) ProcessQueryTransactionsPayload(pld *payload.QueryTransactionsPayload) {
	ss.logger.Trace("Process transactions request payload", "pld", pld)

	trxs := ss.prepareTransactions(pld.IDs)

	if len(trxs) > 0 {
		ss.BroadcastTransactions(trxs)
	}
}

func (ss *StateSync) ProcessTransactionsPayload(pld *payload.TransactionsPayload) {
	ss.logger.Trace("Process transactions payload", "pld", pld)

	ss.addTransactionsToCache(pld.Transactions)

	for _, trx := range pld.Transactions {
		if err := ss.txPool.AppendTx(trx); err != nil {
			// TODO: set peer as bad peer?
		}
	}
}

func (ss *StateSync) prepareTransactions(ids []crypto.Hash) []*tx.Tx {
	trxs := make([]*tx.Tx, 0, len(ids))

	for _, id := range ids {
		trx := ss.cache.GetTransaction(id)
		if trx == nil {
			ss.logger.Debug("Unable to find a transaction", "id", id.Fingerprint())
			continue
		}
		trxs = append(trxs, trx)
	}
	return trxs
}

func (ss *StateSync) prepareBlocksAndTransactions(from, count int) ([]*block.Block, []*tx.Tx) {
	ourHeight := ss.state.LastBlockHeight()

	if from > ourHeight {
		ss.logger.Warn("We don't have block at this height", "height", from)
		return nil, nil
	}

	if from+count > ourHeight {
		count = ourHeight - from + 1
	}

	blocks := make([]*block.Block, 0, count)
	trxs := make([]*tx.Tx, 0)

	for h := from; h < from+count; h++ {
		b := ss.cache.GetBlock(h)
		if b == nil {
			ss.logger.Warn("Unable to find a block", "height", h)
			return nil, nil
		}
		for _, id := range b.TxIDs().IDs() {
			trx := ss.cache.GetTransaction(id)
			if trx != nil {
				trxs = append(trxs, trx)
			} else {
				ss.logger.Debug("Unable to find a transaction", "id", id.Fingerprint())
				return nil, nil
			}
		}

		blocks = append(blocks, b)
	}

	return blocks, trxs
}

func (ss *StateSync) addBlocksToCache(from int, blocks []*block.Block) {
	for _, block := range blocks {
		ss.cache.AddCommit(
			block.Header().LastBlockHash(),
			block.LastCommit())

		ss.cache.AddBlock(from, block)

		from++
	}
}

func (ss *StateSync) addTransactionsToCache(trxs []*tx.Tx) {
	ss.cache.AddTransactions(trxs)
}

func (ss *StateSync) tryCommitBlocks() {
	for {
		ourHeight := ss.state.LastBlockHeight()
		b := ss.cache.GetBlock(ourHeight + 1)
		if b == nil {
			break
		}
		c := ss.cache.GetCommit(b.Hash())
		if c == nil {
			break
		}
		ss.logger.Trace("Committing block", "height", ourHeight+1, "block", b)
		if err := ss.state.ApplyBlock(ourHeight+1, *b, *c); err != nil {
			ss.logger.Error("Committing block failed", "block", b, "err", err, "height", ourHeight+1)
			// We will ask peers to send this block later ...
			break
		}
	}
}

func (ss *StateSync) RequestForMoreBlock() {
	from := ss.state.LastBlockHeight()
	l := ss.peerSet.GetPeerList()
	for _, p := range l {
		if p.InitialBlockDownload() {
			if p.Height() > from {
				to := from + DownloadBlockInterval
				if to > p.Height() {
					to = p.Height()
				}

				ss.BroadcastDownloadRequest(p.PeerID(), from+1, to)
				from = to
			}
		}
	}
}

func (ss *StateSync) RequestForLatestBlock() {
	ourHeight := ss.state.LastBlockHeight()
	p := ss.peerSet.FindHighestPeer()
	ss.BroadcastLatestBlocksRequest(p.PeerID(), ourHeight+1)
}
