package txpool

import (
	"slices"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

var _ TxPool = &MockTxPool{}

// MockTxPool is a testing mock.
type MockTxPool struct {
	Trxs []*tx.Tx
}

func MockingTxPool() *MockTxPool {
	return &MockTxPool{
		Trxs: make([]*tx.Tx, 0),
	}
}
func (*MockTxPool) SetNewSandboxAndRecheck(_ sandbox.Sandbox) {}
func (m *MockTxPool) PendingTx(id tx.ID) *tx.Tx {
	for _, t := range m.Trxs {
		if t.ID() == id {
			return t
		}
	}

	return nil
}

func (m *MockTxPool) QueryTx(id tx.ID) *tx.Tx {
	return m.PendingTx(id)
}

func (m *MockTxPool) HasTx(id tx.ID) bool {
	for _, t := range m.Trxs {
		if t.ID() == id {
			return true
		}
	}

	return false
}

func (m *MockTxPool) Size() int {
	return len(m.Trxs)
}

func (*MockTxPool) String() string {
	return ""
}

func (m *MockTxPool) AppendTx(trx *tx.Tx) error {
	m.Trxs = append(m.Trxs, trx)

	return nil
}

func (m *MockTxPool) AppendTxAndBroadcast(trx *tx.Tx) error {
	m.Trxs = append(m.Trxs, trx)

	return nil
}

func (m *MockTxPool) RemoveTx(id hash.Hash) {
	for i, trx := range m.Trxs {
		if trx.ID() == id {
			m.Trxs = slices.Delete(m.Trxs, i, i+1)

			return
		}
	}
}

func (m *MockTxPool) PrepareBlockTransactions() block.Txs {
	txs := make([]*tx.Tx, m.Size())
	copy(txs, m.Trxs)

	return txs
}

func (*MockTxPool) EstimatedFee(_ amount.Amount, _ payload.Type) amount.Amount {
	return amount.Amount(0.1e9)
}

func (m *MockTxPool) AllPendingTxs() []*tx.Tx {
	return make([]*tx.Tx, m.Size())
}
