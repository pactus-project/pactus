package txpool

import (
	"fmt"
	"time"

	"github.com/zarbchain/zarb-go/execution"

	"github.com/zarbchain/zarb-go/sandbox"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/libs/linkedmap"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/tx"
)

type txPool struct {
	lk deadlock.RWMutex

	config      *Config
	checker     *execution.Execution
	pendings    *linkedmap.LinkedMap
	appendTxCh  chan *tx.Tx
	broadcastCh chan *message.Message
	logger      *logger.Logger
}

func NewTxPool(
	conf *Config,
	broadcastCh chan *message.Message) (TxPool, error) {
	pool := &txPool{
		config:      conf,
		pendings:    linkedmap.NewLinkedMap(conf.MaxSize),
		appendTxCh:  make(chan *tx.Tx, 5),
		broadcastCh: broadcastCh,
	}

	pool.logger = logger.NewLogger("_pool", pool)
	return pool, nil
}

func (pool *txPool) SetSandbox(sb sandbox.Sandbox) {
	pool.checker = execution.NewExecution(sb)
}

func (pool *txPool) AppendTxs(trxs []tx.Tx) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	for _, trx := range trxs {
		if err := pool.appendTx(trx); err != nil {
			pool.logger.Info("Error on appending a transaction", "err", err)
		}
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
		return errors.Errorf(errors.ErrInvalidTx, "Transaction is already in pool. hash: %v", trx.Hash())
	}

	if err := pool.checker.Execute(&trx); err != nil {
		pool.logger.Error("Invalid transaction", "tx", trx, "err", err)
		return err
	}

	pool.pendings.PushBack(trx.Hash(), &trx)

	return nil
}

func (pool *txPool) RemoveTx(hash crypto.Hash) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	pool.pendings.Remove(hash)
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
