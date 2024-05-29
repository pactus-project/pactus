package txpool

import (
	"fmt"
	"github.com/pactus-project/pactus/types/amount"
	"sync"

	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/linkedlist"
	"github.com/pactus-project/pactus/util/linkedmap"
	"github.com/pactus-project/pactus/util/logger"
)

type txPool struct {
	lk sync.RWMutex

	config      *Config
	checker     *execution.Execution
	sandbox     sandbox.Sandbox
	pools       map[payload.Type]pool
	broadcastCh chan message.Message
	logger      *logger.SubLogger

	baseFee float64
}

func NewTxPool(conf *Config, baseFee float64, broadcastCh chan message.Message) TxPool {
	minValue := conf.minValue()

	pools := make(map[payload.Type]pool)
	pools[payload.TypeTransfer] = newPool(conf.transferPoolSize(), minValue)
	pools[payload.TypeBond] = newPool(conf.bondPoolSize(), 0)
	pools[payload.TypeUnbond] = newPool(conf.unbondPoolSize(), 0)
	pools[payload.TypeWithdraw] = newPool(conf.withdrawPoolSize(), minValue)
	pools[payload.TypeSortition] = newPool(conf.sortitionPoolSize(), 0)

	pool := &txPool{
		config:      conf,
		checker:     execution.NewChecker(),
		pools:       pools,
		broadcastCh: broadcastCh,
		baseFee:     baseFee,
	}

	pool.logger = logger.NewSubLogger("_pool", pool)

	return pool
}

func (p *txPool) SetNewSandboxAndRecheck(sb sandbox.Sandbox) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.sandbox = sb
	p.logger.Debug("set new sandbox")

	var next *linkedlist.Element[linkedmap.Pair[tx.ID, *tx.Tx]]
	for _, pool := range p.pools {
		for e := pool.list.HeadNode(); e != nil; e = next {
			next = e.Next
			trx := e.Data.Value

			if err := p.checkTx(trx); err != nil {
				p.logger.Debug("invalid transaction after rechecking", "id", trx.ID())
				pool.list.Remove(trx.ID())
			}
		}
	}
}

func (p *txPool) calculateDynamicFee() (amount.Amount, error) {
	totalSize := p.config.MaxSize
	currentSize := p.Size()
	usageRatio := float64(currentSize) / float64(totalSize)

	switch {
	case usageRatio >= 0.90:
		return amount.NewAmount(p.baseFee * 1000)
	case usageRatio >= 0.80:
		return amount.NewAmount(p.baseFee * 500)
	case usageRatio >= 0.75:
		return amount.NewAmount(p.baseFee * 200)
	case usageRatio >= 0.50:
		return amount.NewAmount(p.baseFee * 100)
	default:
		return amount.NewAmount(p.baseFee)
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
	payloadType := trx.Payload().Type()
	pool := p.pools[payloadType]
	if pool.list.Has(trx.ID()) {
		p.logger.Trace("transaction is already in pool", "id", trx.ID())

		return nil
	}

	if trx.Payload().Value() < pool.minValue {
		p.logger.Warn("low value transaction", "tx", trx, "minValue", pool.minValue)

		return AppendError{
			Err: fmt.Errorf("low value transaction, expected to be more than %s", pool.minValue),
		}
	}

	if err := p.checkTx(trx); err != nil {
		return AppendError{
			Err: err,
		}
	}

	pool.list.PushBack(trx.ID(), trx)
	p.logger.Debug("transaction appended into pool", "tx", trx)

	return nil
}

func (p *txPool) checkTx(trx *tx.Tx) error {
	if err := p.checker.Execute(trx, p.sandbox); err != nil {
		p.logger.Debug("invalid transaction", "tx", trx, "error", err)

		return err
	}

	return nil
}

func (p *txPool) RemoveTx(id tx.ID) {
	p.lk.Lock()
	defer p.lk.Unlock()

	for _, pool := range p.pools {
		if pool.list.Remove(id) {
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
		n := pool.list.GetNode(id)
		if n != nil {
			return n.Data.Value
		}
	}

	return nil
}

func (p *txPool) PrepareBlockTransactions() block.Txs {
	trxs := make([]*tx.Tx, 0, p.Size())

	p.lk.RLock()
	defer p.lk.RUnlock()

	// Appending one sortition transaction
	poolSortition := p.pools[payload.TypeSortition]
	for n := poolSortition.list.HeadNode(); n != nil; n = n.Next {
		trxs = append(trxs, n.Data.Value)
	}

	// Appending bond transactions
	poolBond := p.pools[payload.TypeBond]
	for n := poolBond.list.HeadNode(); n != nil; n = n.Next {
		trxs = append(trxs, n.Data.Value)
	}

	// Appending unbond transactions
	poolUnbond := p.pools[payload.TypeUnbond]
	for n := poolUnbond.list.HeadNode(); n != nil; n = n.Next {
		trxs = append(trxs, n.Data.Value)
	}

	// Appending withdraw transactions
	poolWithdraw := p.pools[payload.TypeWithdraw]
	for n := poolWithdraw.list.HeadNode(); n != nil; n = n.Next {
		trxs = append(trxs, n.Data.Value)
	}

	// Appending transfer transactions
	poolTransfer := p.pools[payload.TypeTransfer]
	for n := poolTransfer.list.HeadNode(); n != nil; n = n.Next {
		trxs = append(trxs, n.Data.Value)
	}

	return trxs
}

func (p *txPool) HasTx(id tx.ID) bool {
	p.lk.RLock()
	defer p.lk.RUnlock()

	for _, pool := range p.pools {
		if pool.list.Has(id) {
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
		size += pool.list.Size()
	}

	return size
}

func (p *txPool) String() string {
	return fmt.Sprintf("{💸 %v 🔐 %v 🔓 %v 🎯 %v 🧾 %v}",
		p.pools[payload.TypeTransfer].list.Size(),
		p.pools[payload.TypeBond].list.Size(),
		p.pools[payload.TypeUnbond].list.Size(),
		p.pools[payload.TypeSortition].list.Size(),
		p.pools[payload.TypeWithdraw].list.Size(),
	)
}
