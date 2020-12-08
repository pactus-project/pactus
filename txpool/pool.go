package txpool

import (
	"fmt"
	"time"

	"github.com/zarbchain/zarb-go/sandbox"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/libs/linkedmap"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/tx"
)

// TODO: We need to have LRU cache for mempool.
// We need to prune stale transactions
// A transaction might valid at heigh M, but invalid at height N (N > M)

type txPool struct {
	lk deadlock.RWMutex

	config      *Config
	sandbox     sandbox.Sandbox
	pendings    *linkedmap.LinkedMap
	appendTxCh  chan *tx.Tx
	broadcastCh chan *message.Message
	isSyncing   bool
	logger      *logger.Logger
}

func NewTxPool(
	conf *Config,
	broadcastCh chan *message.Message) (TxPool, error) {
	pool := &txPool{
		config:      conf,
		pendings:    linkedmap.NewLinkedMap(conf.MaxSize),
		appendTxCh:  make(chan *tx.Tx, 5),
		isSyncing:   true,
		broadcastCh: broadcastCh,
	}

	pool.logger = logger.NewLogger("_pool", pool)
	return pool, nil
}

func (pool *txPool) SetSandbox(sandbox sandbox.Sandbox) {
	pool.sandbox = sandbox
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
	if pool.pendings.Has(trx.Hash()) {
		return errors.Errorf(errors.ErrInvalidTx, "Transaction is alreasy in pool. hash: %v", trx.Hash())
	}

	// When we are syncing we should ignore validating transaction
	// Some transactions might belong to blocks which we don't have them yet
	if !pool.isSyncing {
		if err := pool.validateTx(&trx); err != nil {
			pool.logger.Error("Invalid transaction", "tx", trx, "err", err)
			return err
		}
	}

	pool.pendings.PushBack(trx.Hash(), &trx)

	return nil
}

func (pool *txPool) SetIsSyncing(syncing bool) {
	pool.isSyncing = syncing
}

func (pool *txPool) RemoveTx(hash crypto.Hash) *tx.Tx {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	val := pool.pendings.Remove(hash)
	if val != nil {
		return val.(*tx.Tx)
	}

	return nil
}

func (pool *txPool) PendingTx(hash crypto.Hash) *tx.Tx {
	pool.lk.RLock()

	val, found := pool.pendings.Get(hash)
	if found {
		trx := val.(*tx.Tx)
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

	return pool.pendings.Has(hash)
}

func (pool *txPool) Size() int {
	pool.lk.RLock()
	defer pool.lk.RUnlock()

	return pool.pendings.Size()
}

func (pool *txPool) Fingerprint() string {
	return fmt.Sprintf("{%v}", pool.pendings.Size())
}
