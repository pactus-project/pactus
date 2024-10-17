package txpool

import (
	"fmt"
	"sync"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/store"

	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/linkedlist"
	"github.com/pactus-project/pactus/util/linkedmap"
	"github.com/pactus-project/pactus/util/logger"
)

type txPool struct {
	lk sync.RWMutex

	config         *Config
	sandbox        sandbox.Sandbox
	pools          map[payload.Type]pool
	consumptionMap map[crypto.Address]uint32
	broadcastCh    chan message.Message
	strReader      store.Reader
	logger         *logger.SubLogger
}

func NewTxPool(conf *Config, broadcastCh chan message.Message, storeReader store.Reader) TxPool {
	pools := make(map[payload.Type]pool)
	pools[payload.TypeTransfer] = newPool(conf.transferPoolSize(), conf.minFee())
	pools[payload.TypeBond] = newPool(conf.bondPoolSize(), conf.minFee())
	pools[payload.TypeUnbond] = newPool(conf.unbondPoolSize(), 0)
	pools[payload.TypeWithdraw] = newPool(conf.withdrawPoolSize(), conf.minFee())
	pools[payload.TypeSortition] = newPool(conf.sortitionPoolSize(), 0)

	pool := &txPool{
		config:         conf,
		pools:          pools,
		consumptionMap: make(map[crypto.Address]uint32),
		strReader:      storeReader,
		broadcastCh:    broadcastCh,
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
	payloadPool := p.pools[payloadType]
	if payloadPool.list.Has(trx.ID()) {
		p.logger.Trace("transaction is already in pool", "id", trx.ID())

		return nil
	}

	if !trx.IsFreeTx() {
		if trx.Fee() < payloadPool.estimatedFee() {
			p.logger.Warn("low fee transaction", "tx", trx, "minFee", payloadPool.estimatedFee())

			return AppendError{
				Err: fmt.Errorf("low fee transaction, expected to be more than %s", payloadPool.estimatedFee()),
			}
		}
	}

	if err := p.checkTx(trx); err != nil {
		return AppendError{
			Err: err,
		}
	}

	payloadPool.list.PushBack(trx.ID(), trx)
	p.logger.Debug("transaction appended into pool", "tx", trx)

	return nil
}

func (p *txPool) checkTx(trx *tx.Tx) error {
	if err := execution.CheckAndExecute(trx, p.sandbox, false); err != nil {
		p.logger.Debug("invalid transaction", "tx", trx, "error", err)

		return err
	}

	return nil
}

func (p *txPool) HandleCommittedBlock(blk *block.Block) error {
	p.lk.Lock()
	defer p.lk.Unlock()

	for _, trx := range blk.Transactions() {
		p.removeTx(trx.ID())

		p.handleIncreaseConsumption(trx)
	}

	return p.handleDecreaseConsumption(blk.Height())
}

func (p *txPool) handleIncreaseConsumption(trx *tx.Tx) {
	if trx.IsTransferTx() || trx.IsBondTx() || trx.IsWithdrawTx() {
		signer := trx.Payload().Signer()

		// retrieve existing consumption or start with 0
		p.consumptionMap[signer] = p.consumptionMap[signer] + uint32(trx.SerializeSize())
	}
}

func (p *txPool) handleDecreaseConsumption(height uint32) error {
	if height <= p.config.ConsumptionBlocks {
		return nil
	}

	oldConsumptionHeight := height - p.config.ConsumptionBlocks
	committedBlock, err := p.strReader.Block(oldConsumptionHeight)
	if err != nil {
		return err
	}

	blk, err := committedBlock.ToBlock()
	if err != nil {
		return err
	}

	for _, trx := range blk.Transactions() {
		if trx.IsTransferTx() || trx.IsBondTx() || trx.IsWithdrawTx() {
			signer := trx.Payload().Signer()
			if v, ok := p.consumptionMap[signer]; ok {
				p.consumptionMap[signer] = v - uint32(trx.SerializeSize())
			}
		}
	}

	return nil
}

func (p *txPool) removeTx(id tx.ID) {
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

func (p *txPool) EstimatedFee(_ amount.Amount, payloadType payload.Type) amount.Amount {
	p.lk.RLock()
	defer p.lk.RUnlock()

	payloadPool, ok := p.pools[payloadType]
	if !ok {
		return 0
	}

	return payloadPool.estimatedFee()
}

func (p *txPool) AllPendingTxs() []*tx.Tx {
	p.lk.RLock()
	defer p.lk.RUnlock()

	txs := make([]*tx.Tx, 0, p.Size())

	var next *linkedlist.Element[linkedmap.Pair[tx.ID, *tx.Tx]]
	for _, pool := range p.pools {
		for e := pool.list.HeadNode(); e != nil; e = next {
			next = e.Next
			trx := e.Data.Value

			txs = append(txs, trx)
		}
	}

	return txs
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
