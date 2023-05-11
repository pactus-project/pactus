package txpool

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/linkedmap"
	"github.com/pactus-project/pactus/util/logger"
)

type txPool struct {
	lk sync.RWMutex

	config      *Config
	checker     *execution.Execution
	sandbox     sandbox.Sandbox
	pools       map[payload.Type]*linkedmap.LinkedMap
	broadcastCh chan message.Message
	logger      *logger.Logger
}

func NewTxPool(conf *Config, broadcastCh chan message.Message) TxPool {
	pendings := make(map[payload.Type]*linkedmap.LinkedMap)

	pendings[payload.PayloadTypeTransfer] = linkedmap.NewLinkedMap(conf.sendPoolSize())
	pendings[payload.PayloadTypeBond] = linkedmap.NewLinkedMap(conf.bondPoolSize())
	pendings[payload.PayloadTypeUnbond] = linkedmap.NewLinkedMap(conf.unbondPoolSize())
	pendings[payload.PayloadTypeWithdraw] = linkedmap.NewLinkedMap(conf.withdrawPoolSize())
	pendings[payload.PayloadTypeSortition] = linkedmap.NewLinkedMap(conf.sortitionPoolSize())

	pool := &txPool{
		config:      conf,
		checker:     execution.NewChecker(),
		pools:       pendings,
		broadcastCh: broadcastCh,
	}

	pool.logger = logger.NewLogger("_pool", pool)
	return pool
}

func (p *txPool) SetNewSandboxAndRecheck(sb sandbox.Sandbox) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.sandbox = sb
	p.logger.Debug("set new sandbox")

	var next *list.Element
	for _, pool := range p.pools {
		for e := pool.FirstElement(); e != nil; e = next {
			next = e.Next()
			trx := e.Value.(*linkedmap.Pair).Second.(*tx.Tx)

			if err := p.checkTx(trx); err != nil {
				p.logger.Debug("invalid transaction after rechecking", "id", trx.ID())
				pool.Remove(trx.ID())
			}
		}
	}
}

// AppendTx validates the transaction and add it into the transaction pool
// without broadcast it.
func (p *txPool) AppendTx(trx *tx.Tx) error {
	p.lk.Lock()
	defer p.lk.Unlock()

	return p.appendTx(trx)
}

// AppendTxAndBroadcast validates the transaction, add it into the transaction pool
// and broadcast it.
func (p *txPool) AppendTxAndBroadcast(trx *tx.Tx) error {
	p.lk.Lock()
	defer p.lk.Unlock()

	if err := p.appendTx(trx); err != nil {
		return err
	}

	go func(t *tx.Tx) {
		p.broadcastCh <- message.NewTransactionsMessage([]*tx.Tx{t})
	}(trx)

	return nil
}

func (p *txPool) appendTx(trx *tx.Tx) error {
	pool := p.pools[trx.Payload().Type()]
	if pool.Has(trx.ID()) {
		p.logger.Trace("transaction is already in pool", "id", trx.ID())
		return nil
	}

	if err := p.checkTx(trx); err != nil {
		return err
	}

	pool.PushBack(trx.ID(), trx)
	p.logger.Debug("transaction appended into pool", "tx", trx)

	return nil
}

func (p *txPool) checkTx(trx *tx.Tx) error {
	if err := p.checker.Execute(trx, p.sandbox); err != nil {
		p.logger.Debug("invalid transaction", "tx", trx, "err", err)
		return err
	}
	return nil
}

func (p *txPool) RemoveTx(id tx.ID) {
	p.lk.Lock()
	defer p.lk.Unlock()

	for _, pool := range p.pools {
		if pool.Remove(id) {
			break
		}
	}
}

// PendingTx searches inside the transaction pool and returns the associated transaction.
// If transaction doesn't exist inside the pool, it returns nil.
func (p *txPool) PendingTx(id tx.ID) *tx.Tx {
	p.lk.Lock()
	defer p.lk.Unlock()

	for _, pool := range p.pools {
		val, found := pool.Get(id)
		if found {
			trx := val.(*tx.Tx)
			return trx
		}
	}

	return nil
}

func (p *txPool) PrepareBlockTransactions() block.Txs {
	trxs := make([]*tx.Tx, 0, p.Size())

	p.lk.RLock()
	defer p.lk.RUnlock()

	// Appending one sortition transaction
	poolSortition := p.pools[payload.PayloadTypeSortition]
	for e := poolSortition.FirstElement(); e != nil; e = e.Next() {
		trx := e.Value.(*linkedmap.Pair).Second.(*tx.Tx)
		trxs = append(trxs, trx)
	}

	// Appending bond transactions
	poolBond := p.pools[payload.PayloadTypeBond]
	for e := poolBond.FirstElement(); e != nil; e = e.Next() {
		trx := e.Value.(*linkedmap.Pair).Second.(*tx.Tx)
		trxs = append(trxs, trx)
	}

	// Appending unbond transactions
	poolUnbond := p.pools[payload.PayloadTypeUnbond]
	for e := poolUnbond.FirstElement(); e != nil; e = e.Next() {
		trx := e.Value.(*linkedmap.Pair).Second.(*tx.Tx)
		trxs = append(trxs, trx)
	}

	// Appending withdraw transactions
	poolWithdraw := p.pools[payload.PayloadTypeWithdraw]
	for e := poolWithdraw.FirstElement(); e != nil; e = e.Next() {
		trx := e.Value.(*linkedmap.Pair).Second.(*tx.Tx)
		trxs = append(trxs, trx)
	}

	// Appending send transactions
	poolSend := p.pools[payload.PayloadTypeTransfer]
	for e := poolSend.FirstElement(); e != nil; e = e.Next() {
		trx := e.Value.(*linkedmap.Pair).Second.(*tx.Tx)
		trxs = append(trxs, trx)
	}

	return trxs
}

func (p *txPool) HasTx(id tx.ID) bool {
	p.lk.RLock()
	defer p.lk.RUnlock()

	for _, pool := range p.pools {
		if pool.Has(id) {
			return true
		}
	}

	return false
}

func (p *txPool) Size() int {
	p.lk.RLock()
	defer p.lk.RUnlock()

	size := 0
	for _, pool := range p.pools {
		size += pool.Size()
	}
	return size
}

func (p *txPool) Fingerprint() string {
	return fmt.Sprintf("{💸 %v 🔐 %v 🔓 %v 🎯 %v 🧾 %v}",
		p.pools[payload.PayloadTypeTransfer].Size(),
		p.pools[payload.PayloadTypeBond].Size(),
		p.pools[payload.PayloadTypeUnbond].Size(),
		p.pools[payload.PayloadTypeSortition].Size(),
		p.pools[payload.PayloadTypeWithdraw].Size(),
	)
}
