package sync

import (
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

var LatestBlockInterval = 720
var DownloadBlockInterval = 1000

type StateSync struct {
	config    *Config
	selfID    peer.ID
	cache     *cache.Cache
	state     state.State
	peerSet   *peerset.PeerSet
	logger    *logger.Logger
	publishFn PublishMessageFn
}

func NewStateSync(
	conf *Config,
	selfID peer.ID,
	cache *cache.Cache,
	state state.State,
	peerSet *peerset.PeerSet,
	logger *logger.Logger,
	publishFn PublishMessageFn) *StateSync {
	return &StateSync{
		config:    conf,
		selfID:    selfID,
		cache:     cache,
		state:     state,
		peerSet:   peerSet,
		logger:    logger,
		publishFn: publishFn,
	}
}

func (ss *StateSync) BroadcastLatestBlocksRequest(from int, hash crypto.Hash) {
	msg := message.NewLatestBlocksRequestMessage(from, hash)
	ss.publishFn(msg)
}

func (ss *StateSync) BroadcastLatestBlocks(from int, blocks []*block.Block, trxs []*tx.Tx, commit *block.Commit) {
	msg := message.NewLatestBlocksMessage(from, blocks, trxs, commit)
	msg.CompressIt()
	ss.publishFn(msg)
}

func (ss *StateSync) BroadcastTransactions(txs []*tx.Tx) {
	msg := message.NewTransactionsMessage(txs)
	ss.publishFn(msg)
}

func (ss *StateSync) BroadcastDownloadResponse(status, from int, blocks []*block.Block, trxs []*tx.Tx) {
	msg := message.NewDownloadResponseMessage(status, from, blocks, trxs)
	msg.CompressIt()
	ss.publishFn(msg)
}

func (ss *StateSync) ProcessLatestBlocksRequestPayload(pld *payload.LatestBlocksRequestPayload) {
	ss.logger.Trace("Process blocks request payload", "pld", pld)

	ourHeight := ss.state.LastBlockHeight()
	if pld.From < ourHeight-LatestBlockInterval {
		ss.logger.Warn("We can can only send blocks form last two hours", "from", pld.From)
		return
	}

	b := ss.cache.GetBlock(pld.From)
	if b != nil {
		if !b.Header().LastBlockHash().EqualsTo(pld.LastBlockHash) {
			ss.logger.Info("Previous block hash is invalid",
				"height", pld.From-1,
				"ourHash", b.Header().LastBlockHash(),
				"peerHash", pld.LastBlockHash)

			return
		}
	}

	blocks, trxs := ss.prepareBlocksAndTransactions(pld.From, ourHeight)
	if len(blocks) > 0 {
		commit := ss.state.LastCommit()

		// Help peer to catch up
		ss.BroadcastLatestBlocks(pld.From, blocks, trxs, commit)
	}
}

func (ss *StateSync) ProcessLatestBlocksPayload(pld *payload.LatestBlocksPayload) {
	ss.logger.Trace("Process blocks payload", "pld", pld)

	ourHeight := ss.state.LastBlockHeight()
	if ourHeight >= pld.To() {
		return
	}

	if pld.LastCommit != nil {
		ss.cache.AddCommit(
			pld.Blocks[len(pld.Blocks)-1].Header().Hash(),
			pld.LastCommit)
	}

	ss.appendBlocksToCache(pld.From, pld.Blocks)
	ss.appendTransactionsToCache(pld.Transactions)
	ss.tryCommitBlocks()
}

func (ss *StateSync) ProcessDownloadRequestMessage(pld *payload.DownloadRequestPayload) {
	if pld.PeerID != ss.selfID {
		return
	}

	if pld.To-pld.From > DownloadBlockInterval {
		return
	}

	blocks, trxs := ss.prepareBlocksAndTransactions(pld.From, pld.To)

	if len(blocks) > 0 {
		ourHeight := ss.state.LastBlockHeight()
		status := payload.DownloadResponseCodeOK
		if ourHeight <= pld.To {
			status = payload.DownloadResponseCodeNoMoreBlock
		}
		ss.BroadcastDownloadResponse(status, pld.From, blocks, trxs)
	}
}

func (ss *StateSync) ProcessDownloadResponseMessage(pld *payload.DownloadResponsePayload) {
	ss.appendBlocksToCache(pld.From, pld.Blocks)
	ss.appendTransactionsToCache(pld.Transactions)
	ss.tryCommitBlocks()

	if pld.Status == payload.DownloadResponseCodeOK {
		ss.RequestForMoreBlock()
	}
}

func (ss *StateSync) ProcessTransactionsRequestPayload(pld *payload.TransactionsRequestPayload) {
	ss.logger.Trace("Process txs request Payload", "pld", pld)

	trxs := ss.prepareTransactions(pld.IDs)

	if len(trxs) > 0 {
		ss.BroadcastTransactions(trxs)
	}
}

func (ss *StateSync) ProcessTransactionsPayload(pld *payload.TransactionsPayload) {
	ss.logger.Trace("Process txs payload", "pld", pld)

	ss.appendTransactionsToCache(pld.Transactions)
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

func (ss *StateSync) prepareBlocksAndTransactions(from, to int) ([]*block.Block, []*tx.Tx) {
	ourHeight := ss.state.LastBlockHeight()

	// `from` is smaller than `to`
	// `from` is smaller than `ourheight`
	if from > ourHeight {
		ss.logger.Warn("We don't have block at this height", "height", from)
		return nil, nil
	}
	to = util.Min(to, ourHeight)

	blocks := make([]*block.Block, 0, to-from+1)
	trxs := make([]*tx.Tx, 0)

	for h := from; h <= ourHeight; h++ {
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

func (ss *StateSync) appendBlocksToCache(from int, blocks []*block.Block) {
	for _, block := range blocks {
		ss.cache.AddCommit(
			block.Header().LastBlockHash(),
			block.LastCommit())

		ss.cache.AddBlock(from, block)

		from++
	}
}

func (ss *StateSync) appendTransactionsToCache(trxs []*tx.Tx) {
	for _, trx := range trxs {
		ss.cache.AddTransaction(trx)
	}
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
	start := ss.state.LastBlockHeight()

	l := ss.peerSet.GetPeerList()
	for _, p := range l {
		if p.InitialBlockDownload() {
			msg := message.NewDownloadRequestMessage(p.PeerID(), start, start+1000)
			ss.publishFn(msg)
			start += DownloadBlockInterval
		}
	}
}
