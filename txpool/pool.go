package txpool

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/zarbchain/zarb-go/execution"
	"github.com/zarbchain/zarb-go/libs/linkedmap"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/tx"
)

type txPool struct {
	lk sync.RWMutex

	config      *Config
	checker     *execution.Execution
	sandbox     sandbox.Sandbox
	pendings    *linkedmap.LinkedMap
	broadcastCh chan message.Message
	logger      *logger.Logger
}

func NewTxPool(
	conf *Config,
	broadcastCh chan message.Message) (TxPool, error) {
	pool := &txPool{
		config:      conf,
		checker:     execution.NewChecker(),
		pendings:    linkedmap.NewLinkedMap(conf.MaxSize),
		broadcastCh: broadcastCh,
	}

	pool.logger = logger.NewLogger("_pool", pool)
	return pool, nil
}

func (pool *txPool) SetNewSandboxAndRecheck(sb sandbox.Sandbox) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	pool.sandbox = sb

	pool.logger.Debug("set new sandbox")

	var next *list.Element
	for e := pool.pendings.FirstElement(); e != nil; e = next {
		next = e.Next()
		trx := e.Value.(*linkedmap.Pair).Second.(*tx.Tx)

		if err := pool.checkTx(trx); err != nil {
			pool.logger.Debug("invalid transaction after rechecking", "id", trx.ID())
			pool.pendings.Remove(trx.ID())
		}
	}
}

/// AppendTx validates the transaction and add it into the transaction pool
/// without broadcast it.
func (pool *txPool) AppendTx(trx *tx.Tx) error {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	return pool.appendTx(trx)
}

/// AppendTxAndBroadcast validates the transaction, add it into the transaction pool
/// and broadcast it.
func (pool *txPool) AppendTxAndBroadcast(trx *tx.Tx) error {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	if err := pool.appendTx(trx); err != nil {
		return err
	}

	go func(t *tx.Tx) {
		pool.broadcastCh <- message.NewTransactionsMessage([]*tx.Tx{t})
	}(trx)

	return nil
}

func (pool *txPool) appendTx(trx *tx.Tx) error {
	if pool.pendings.Has(trx.ID()) {
		pool.logger.Trace("transaction is already in pool", "id", trx.ID())
		return nil
	}

	if err := pool.checkTx(trx); err != nil {
		return err
	}

	pool.pendings.PushBack(trx.ID(), trx)
	pool.logger.Debug("transaction appended into pool", "tx", trx)

	return nil
}

func (pool *txPool) checkTx(trx *tx.Tx) error {
	if err := pool.checker.Execute(trx, pool.sandbox); err != nil {
		pool.logger.Debug("invalid transaction", "tx", trx, "err", err)
		return err
	}
	return nil
}

func (pool *txPool) RemoveTx(id tx.ID) {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	pool.pendings.Remove(id)
}

/// PendingTx searches inside the transaction pool and returns the associated transaction.
/// If transaction doesn't exist inside the pool, it returns nil.
func (pool *txPool) PendingTx(id tx.ID) *tx.Tx {
	pool.lk.Lock()
	defer pool.lk.Unlock()

	val, found := pool.pendings.Get(id)
	if found {
		trx := val.(*tx.Tx)
		return trx
	}

	return nil
}

/// QueryTx searches inside the transaction pool and returns the associated transaction.
/// If transaction doesn't exist inside the pool, it queries from other nodes.
func (pool *txPool) QueryTx(id tx.ID) *tx.Tx {
	trx := pool.PendingTx(id)
	if trx != nil {
		return trx
	}

	pool.logger.Debug("querying transaction from the network", "id", id)
	pool.broadcastCh <- message.NewQueryTransactionsMessage([]tx.ID{id})

	duration := time.Millisecond * 500
	timeout := time.NewTicker(duration)
	counter := 0

	for i := 0; i < 4; i++ {
		<-timeout.C
		trx := pool.PendingTx(id)
		if trx != nil {
			return trx
		}
	}

	pool.logger.Warn("querying transaction failed", "id", id, "duration", duration*time.Duration(counter))
	return nil
}

func (pool *txPool) AllTransactions() []*tx.Tx {
	pool.lk.RLock()
	defer pool.lk.RUnlock()

	trxs := make([]*tx.Tx, 0, pool.pendings.Size())
	for e := pool.pendings.FirstElement(); e != nil; e = e.Next() {
		trx := e.Value.(*linkedmap.Pair).Second.(*tx.Tx)
		trxs = append(trxs, trx)
	}

	return trxs
}

func (pool *txPool) HasTx(id tx.ID) bool {
	pool.lk.RLock()
	defer pool.lk.RUnlock()

	return pool.pendings.Has(id)
}

func (pool *txPool) Size() int {
	pool.lk.RLock()
	defer pool.lk.RUnlock()

	return pool.pendings.Size()
}

func (pool *txPool) Fingerprint() string {
	return fmt.Sprintf("{%v}", pool.pendings.Size())
}
