package txpool

import (
	"container/list"
	"fmt"
	"time"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/tx"
)

// TODO: We need to have LRU cache for mempool.
// We need to prune stale transactions
// A transaction might valid at heigh M, but invalid at height N (N > M)
type TxPoolReader interface {
	PendingTx(hash crypto.Hash) (*tx.Tx, bool)
}

type TxPool struct {
	lk deadlock.RWMutex

	config       *Config
	pendingsList *list.List
	pendingsMap  map[crypto.Hash]*list.Element
	appendTxCh   chan *tx.Tx
	broadcastCh  chan *message.Message
	logger       *logger.Logger
}

func NewTxPool(
	conf *Config,
	broadcastCh chan *message.Message) (*TxPool, error) {
	pool := &TxPool{
		config:       conf,
		pendingsList: list.New(),
		pendingsMap:  make(map[crypto.Hash]*list.Element),
		appendTxCh:   make(chan *tx.Tx, 5),
		broadcastCh:  broadcastCh,
	}

	pool.logger = logger.NewLogger("_pool", pool)
	return pool, nil
}

func (pool *TxPool) AppendTxs(txs []tx.Tx) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	for _, tx := range txs {
		pool.appendTx(tx)
	}
}

func (pool *TxPool) AppendTx(tx tx.Tx) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	pool.appendTx(tx)

	pool.appendTxCh <- &tx
}

func (pool *TxPool) AppendTxAndBroadcast(trx tx.Tx) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	pool.appendTx(trx)

	msg := message.NewTxsMessage([]tx.Tx{trx})
	pool.broadcastCh <- msg
}

func (pool *TxPool) appendTx(tx tx.Tx) {
	if pool.pendingsList.Len() >= pool.config.MaxSize {
		pool.logger.Warn("Tx pool is full")
	}

	_, found := pool.pendingsMap[tx.Hash()]
	if found {
		pool.logger.Debug("We already have this transaction", "hash", tx.Hash())
		return
	}
	// TODO:
	// validate transaction

	el := pool.pendingsList.PushFront(&tx)
	pool.pendingsMap[tx.Hash()] = el
}

func (pool *TxPool) RemoveTx(hash crypto.Hash) *tx.Tx {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	el, found := pool.pendingsMap[hash]
	if !found {
		return nil
	}

	pool.pendingsList.Remove(el)
	delete(pool.pendingsMap, hash)

	return el.Value.(*tx.Tx)

}

func (pool *TxPool) PendingTx(hash crypto.Hash) *tx.Tx {
	pool.lk.RLock()

	el, found := pool.pendingsMap[hash]
	if found {
		tx := el.Value.(*tx.Tx)
		pool.lk.RUnlock()
		return tx
	}

	pool.logger.Debug("Request transaction from peers", "hash", hash)
	pool.lk.RUnlock()

	msg := message.NewTxsReqMessage([]crypto.Hash{hash})
	pool.broadcastCh <- msg

	timeout := time.NewTimer(pool.config.WaitingTimeout)

	for {
		select {
		case <-timeout.C:
			pool.logger.Warn("Transaction not received", "hash", hash, "timeout", pool.config.WaitingTimeout)
			return nil
		case tx := <-pool.appendTxCh:
			pool.logger.Debug("Transaction found", "hash", hash)
			if tx.Hash().EqualsTo(hash) {
				return tx
			}
		}
	}
}

func (pool *TxPool) HasTx(hash crypto.Hash) bool {
	pool.lk.RLock()
	defer pool.lk.RUnlock()

	_, found := pool.pendingsMap[hash]
	return found
}

func (pool *TxPool) Size(hash crypto.Hash) int {
	pool.lk.RLock()
	defer pool.lk.RUnlock()

	return pool.pendingsList.Len()
}

func (pool *TxPool) Fingerprint() string {
	return fmt.Sprintf("{%v}", pool.pendingsList.Len())
}
