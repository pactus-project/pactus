package txpool

import (
	"fmt"
	"sync"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/store"
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
	sbx            sandbox.Sandbox
	pools          map[payload.Type]pool
	consumptionMap map[crypto.Address]uint32
	broadcastCh    chan message.Message
	store          store.Reader
	logger         *logger.SubLogger
}

// NewTxPool constructs a new transaction pool with various sub-pools for different transaction types.
// The transaction pool also maintains a consumption map for tracking byte usage per address.
func NewTxPool(conf *Config, storeReader store.Reader, broadcastCh chan message.Message) TxPool {
	pools := make(map[payload.Type]pool)
	pools[payload.TypeTransfer] = newPool(conf.transferPoolSize(), conf.fixedFee())
	pools[payload.TypeBond] = newPool(conf.bondPoolSize(), conf.fixedFee())
	pools[payload.TypeUnbond] = newPool(conf.unbondPoolSize(), 0)
	pools[payload.TypeWithdraw] = newPool(conf.withdrawPoolSize(), conf.fixedFee())
	pools[payload.TypeSortition] = newPool(conf.sortitionPoolSize(), 0)

	pool := &txPool{
		config:         conf,
		pools:          pools,
		consumptionMap: make(map[crypto.Address]uint32),
		store:          storeReader,
		broadcastCh:    broadcastCh,
	}

	pool.logger = logger.NewSubLogger("_pool", pool)

	return pool
}

// SetNewSandboxAndRecheck updates the sandbox and rechecks all transactions,
// removing expired or invalid ones.
func (p *txPool) SetNewSandboxAndRecheck(sbx sandbox.Sandbox) {
	p.lk.Lock()
	defer p.lk.Unlock()

	p.sbx = sbx
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
		minFee := p.estimatedMinimumFee(trx)

		if trx.Fee() < minFee {
			p.logger.Warn("low fee transaction", "txs", trx, "minFee", minFee)

			return AppendError{
				Err: fmt.Errorf("low fee transaction, expected to be more than %s", minFee),
			}
		}
	}

	if err := p.checkTx(trx); err != nil {
		return AppendError{
			Err: err,
		}
	}

	payloadPool.list.PushBack(trx.ID(), trx)
	p.logger.Debug("transaction appended into pool", "trx", trx)

	return nil
}

func (p *txPool) checkTx(trx *tx.Tx) error {
	if err := execution.CheckAndExecute(trx, p.sbx, false); err != nil {
		p.logger.Debug("invalid transaction", "trx", trx, "error", err)

		return err
	}

	return nil
}

func (p *txPool) EstimatedFee(_ amount.Amount, payloadType payload.Type) amount.Amount {
	selectedPool, ok := p.pools[payloadType]
	if !ok {
		return 0
	}

	return selectedPool.estimatedFee()
}

func (p *txPool) HandleCommittedBlock(blk *block.Block) {
	p.lk.Lock()
	defer p.lk.Unlock()

	for _, trx := range blk.Transactions() {
		p.removeTx(trx.ID())
	}

	if p.config.CalculateConsumption() {
		for _, trx := range blk.Transactions() {
			p.handleIncreaseConsumption(trx)
		}
		p.handleDecreaseConsumption(blk.Height())
	}
}

func (p *txPool) handleIncreaseConsumption(trx *tx.Tx) {
	if trx.IsTransferTx() || trx.IsBondTx() || trx.IsWithdrawTx() {
		signer := trx.Payload().Signer()

		p.consumptionMap[signer] += uint32(trx.SerializeSize())
	}
}

func (p *txPool) handleDecreaseConsumption(height uint32) {
	// If height is less than or equal to ConsumptionWindow, nothing to do.
	if height <= p.config.ConsumptionWindow {
		return
	}

	// Calculate the block height that has passed out of the consumption window.
	windowedBlockHeight := height - p.config.ConsumptionWindow
	committedBlock, err := p.store.Block(windowedBlockHeight)
	if err != nil {
		p.logger.Error("failed to read block", "height", windowedBlockHeight, "err", err)

		return
	}

	blk, err := committedBlock.ToBlock()
	if err != nil {
		p.logger.Error("failed to parse block", "height", windowedBlockHeight, "err", err)

		return
	}

	for _, trx := range blk.Transactions() {
		if !trx.IsFreeTx() {
			signer := trx.Payload().Signer()
			if consumption, ok := p.consumptionMap[signer]; ok {
				// Decrease the consumption by the size of the transaction
				consumption -= uint32(trx.SerializeSize())

				if consumption == 0 {
					// If the new value is zero, remove the signer from the consumptionMap
					delete(p.consumptionMap, signer)
				} else {
					// Otherwise, update the map with the new value
					p.consumptionMap[signer] = consumption
				}
			}
		}
	}
}

func (p *txPool) removeTx(txID tx.ID) {
	for _, pool := range p.pools {
		if pool.list.Remove(txID) {
			break
		}
	}
}

// PendingTx searches inside the transaction pool and returns the associated transaction.
// If transaction doesn't exist inside the pool, it returns nil.
func (p *txPool) PendingTx(txID tx.ID) *tx.Tx {
	p.lk.Lock()
	defer p.lk.Unlock()

	for _, pool := range p.pools {
		n := pool.list.GetNode(txID)
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

func (p *txPool) HasTx(txID tx.ID) bool {
	p.lk.RLock()
	defer p.lk.RUnlock()

	for _, pool := range p.pools {
		if pool.list.Has(txID) {
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

func (p *txPool) estimatedMinimumFee(trx *tx.Tx) amount.Amount {
	return p.fixedFee() + p.consumptionalFee(trx)
}

func (p *txPool) fixedFee() amount.Amount {
	return p.config.fixedFee()
}

// consumptionalFee calculates based on the amount of data each address consumes daily.
func (p *txPool) consumptionalFee(trx *tx.Tx) amount.Amount {
	var consumption uint32
	signer := trx.Payload().Signer()
	txSize := uint32(trx.SerializeSize())

	if !p.store.HasPublicKey(signer) {
		consumption = p.config.Fee.DailyLimit
	} else {
		consumption = p.consumptionMap[signer] + txSize + p.getPendingConsumption(signer)
	}

	coefficient := consumption / p.config.Fee.DailyLimit

	consumptionalFee, _ := amount.NewAmount(float64(coefficient) * float64(consumption) * p.config.Fee.UnitPrice)

	return consumptionalFee
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

func (p *txPool) getPendingConsumption(signer crypto.Address) uint32 {
	totalSize := uint32(0)

	// TODO: big o is "o(n * m)"
	var next *linkedlist.Element[linkedmap.Pair[tx.ID, *tx.Tx]]
	for ptype, pool := range p.pools {
		if ptype == payload.TypeTransfer || ptype == payload.TypeBond || ptype == payload.TypeWithdraw {
			for e := pool.list.HeadNode(); e != nil; e = next {
				next = e.Next
				if e.Data.Value.Payload().Signer() == signer {
					totalSize += uint32(e.Data.Value.SerializeSize())
				}
			}
		}
	}

	return totalSize
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
