package sync

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/cache"
	"github.com/zarbchain/zarb-go/tx"
)

var LatestBlockInterval = 720

type DataTopic struct {
	config    *Config
	publishFn PublishMessageFn
	cache     *cache.Cache
	state     state.State
	logger    *logger.Logger
}

func NewDataTopic(
	conf *Config,
	cache *cache.Cache,
	state state.State,
	logger *logger.Logger,
	publishFn PublishMessageFn) *DataTopic {
	return &DataTopic{
		config:    conf,
		cache:     cache,
		state:     state,
		logger:    logger,
		publishFn: publishFn,
	}
}

func (dt *DataTopic) BroadcastLatestBlocksRequest(from int, hash crypto.Hash) {
	msg := message.NewLatestBlocksRequestMessage(from, hash)
	dt.publishFn(msg)
}

func (dt *DataTopic) BroadcastLatestBlocks(from int, blocks []*block.Block, txs []*tx.Tx, commit *block.Commit) {
	msg := message.NewLatestBlocksMessage(from, blocks, txs, commit)
	dt.publishFn(msg)
}

func (dt *DataTopic) BroadcastTransactions(txs []*tx.Tx) {
	msg := message.NewTransactionsMessage(txs)
	dt.publishFn(msg)
}

func (dt *DataTopic) ProcessLatestBlocksRequestPayload(pld *payload.LatestBlocksRequestPayload) {
	dt.logger.Trace("Process blocks request payload", "pld", pld)

	ourHeight := dt.state.LastBlockHeight()
	if ourHeight < pld.From {
		dt.logger.Warn("We don't have block at this height", "from", pld.From)
		return
	}

	if pld.From < ourHeight-LatestBlockInterval {
		dt.logger.Warn("We can can only send blocks form last two hours", "from", pld.From)
		return
	}

	b := dt.cache.GetBlock(pld.From)
	if b != nil {
		if !b.Header().LastBlockHash().EqualsTo(pld.LastBlockHash) {
			dt.logger.Info("We don't have previous block hash",
				"height", pld.From-1,
				"ourHash", b.Header().LastBlockHash(),
				"peerHash", pld.LastBlockHash)

			return
		}
	}

	dt.sendBlocks(pld.From)
}

func (dt *DataTopic) ProcessLatestBlocksPayload(pld *payload.LatestBlocksPayload) {
	dt.logger.Trace("Process blocks payload", "pld", pld)

	ourHeight := dt.state.LastBlockHeight()
	if ourHeight >= pld.To() {
		return
	}

	if pld.Commit != nil {
		dt.cache.AddCommit(
			pld.Blocks[len(pld.Blocks)-1].Header().Hash(),
			pld.Commit)
	}

	height := pld.From
	for _, block := range pld.Blocks {
		dt.cache.AddCommit(
			block.Header().LastBlockHash(),
			block.LastCommit())

		dt.cache.AddBlock(height, block)

		height = height + 1
	}

	for _, trx := range pld.Transactions {
		dt.cache.AddTransaction(trx)
	}

	dt.tryCommitBlocks()

}

func (dt *DataTopic) ProcessTransactionsRequestPayload(pld *payload.TransactionsRequestPayload) {
	dt.logger.Trace("Process txs request Payload", "pld", pld)

	dt.sendTransactions(pld.IDs)
}

func (dt *DataTopic) ProcessTransactionsPayload(pld *payload.TransactionsPayload) {
	dt.logger.Trace("Process txs payload", "pld", pld)

	for _, trx := range pld.Transactions {
		dt.cache.AddTransaction(trx)
	}
}

func (dt *DataTopic) sendBlocks(from int) {
	ourHeight := dt.state.LastBlockHeight()

	if from > ourHeight {
		dt.logger.Warn("Invalid start point", "from", from)
		return
	}

	// Help peer to catch up
	trxs := make([]*tx.Tx, 0)
	blocks := make([]*block.Block, 0, ourHeight-from+1)
	for h := from; h <= ourHeight; h++ {
		b := dt.cache.GetBlock(h)
		if b == nil {
			dt.logger.Warn("Block can't find", "height", h)
			break
		}
		for _, id := range b.TxIDs().IDs() {
			trx := dt.cache.GetTransaction(id)
			if trx != nil {
				trxs = append(trxs, trx)
			} else {
				dt.logger.Debug("Transaction for block can't find", "id", id.Fingerprint())
			}
		}

		blocks = append(blocks, b)
	}

	commit := dt.state.LastCommit()

	dt.BroadcastLatestBlocks(from, blocks, trxs, commit)
}

func (dt *DataTopic) sendTransactions(ids []crypto.Hash) {
	trxs := make([]*tx.Tx, 0, len(ids))
	for _, id := range ids {
		trx := dt.cache.GetTransaction(id)
		if trx != nil {
			trxs = append(trxs, trx)
		} else {
			dt.logger.Debug("Transaction can't find", "id", id.Fingerprint())
		}
	}

	if len(trxs) > 0 {
		dt.BroadcastTransactions(trxs)
	}
}

func (dt *DataTopic) tryCommitBlocks() {
	for {
		ourHeight := dt.state.LastBlockHeight()
		b := dt.cache.GetBlock(ourHeight + 1)
		if b == nil {
			break
		}
		c := dt.cache.GetCommit(b.Hash())
		if c == nil {
			break
		}
		dt.logger.Trace("Committing block", "height", ourHeight+1, "block", b)
		if err := dt.state.ApplyBlock(ourHeight+1, *b, *c); err != nil {
			dt.logger.Error("Committing block failed", "block", b, "err", err, "height", ourHeight+1)
			// We will ask peers to send this block later ...
			break
		}
	}
}
