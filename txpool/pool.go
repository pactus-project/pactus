package txpool

import (
	"container/list"
	"fmt"
	"time"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/tx"
)

// TODO: We need to have LRU cache for mempool.
// We need to prune stale transactions
// A transaction might valid at heigh M, but invalid at height N (N > M)

type txPool struct {
	lk deadlock.RWMutex

	config        *Config
	pendingsList  *list.List
	pendingsMap   map[crypto.Hash]*list.Element
	appendTxCh    chan *tx.Tx
	broadcastCh   chan *message.Message
	maxMemoLenght int
	feeFraction   float64
	minFee        int64
	logger        *logger.Logger
}

func NewTxPool(
	conf *Config,
	broadcastCh chan *message.Message) (TxPool, error) {
	pool := &txPool{
		config:       conf,
		pendingsList: list.New(),
		pendingsMap:  make(map[crypto.Hash]*list.Element),
		appendTxCh:   make(chan *tx.Tx, 5),
		broadcastCh:  broadcastCh,
	}

	pool.logger = logger.NewLogger("_pool", pool)
	return pool, nil
}

func (pool *txPool) UpdateMaxMemoLenght(maxMemoLenght int) {
	pool.maxMemoLenght = maxMemoLenght
}

func (pool *txPool) UpdateFeeFraction(feeFraction float64) {
	pool.feeFraction = feeFraction
}

func (pool *txPool) UpdateMinFee(minFee int64) {
	pool.minFee = minFee
}

func (pool *txPool) AppendTxs(trxs []tx.Tx) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	for _, trx := range trxs {
		pool.appendTx(trx)
	}
}

func (pool *txPool) AppendTx(trx tx.Tx) error {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	if err := pool.appendTx(trx); err != nil {
		return err
	}

	pool.appendTxCh <- &trx

	return nil
}

func (pool *txPool) AppendTxAndBroadcast(trx tx.Tx) error {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	if err := pool.appendTx(trx); err != nil {
		return err
	}

	msg := message.NewTxsMessage([]tx.Tx{trx})
	pool.broadcastCh <- msg

	return nil
}

func (pool *txPool) appendTx(trx tx.Tx) error {
	if pool.pendingsList.Len() >= pool.config.MaxSize {
		return errors.Errorf(errors.ErrInvalidTx, "Transaction pool is full. Size %v", pool.pendingsList.Len())
	}

	_, found := pool.pendingsMap[trx.Hash()]
	if found {
		return errors.Errorf(errors.ErrInvalidTx, "Transaction is alreasy in pool. hash: %v", trx.Hash())
	}

	if err := pool.validateTx(&trx); err != nil {
		return err
	}

	el := pool.pendingsList.PushFront(&trx)
	pool.pendingsMap[trx.Hash()] = el

	return nil
}

func (pool *txPool) RemoveTx(hash crypto.Hash) *tx.Tx {
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

func (pool *txPool) PendingTx(hash crypto.Hash) *tx.Tx {
	pool.lk.RLock()

	el, found := pool.pendingsMap[hash]
	if found {
		trx := el.Value.(*tx.Tx)
		pool.lk.RUnlock()
		return trx
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
		case trx := <-pool.appendTxCh:
			pool.logger.Debug("Transaction found", "hash", hash)
			if trx.Hash().EqualsTo(hash) {
				return trx
			}
		}
	}
}

func (pool *txPool) HasTx(hash crypto.Hash) bool {
	pool.lk.RLock()
	defer pool.lk.RUnlock()

	_, found := pool.pendingsMap[hash]
	return found
}

func (pool *txPool) Size() int {
	pool.lk.RLock()
	defer pool.lk.RUnlock()

	return pool.pendingsList.Len()
}

func (pool *txPool) Fingerprint() string {
	return fmt.Sprintf("{%v}", pool.pendingsList.Len())
}
