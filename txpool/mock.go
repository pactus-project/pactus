package txpool

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
)

var _ TxPool = &MockTxPool{}

// MockTxPool is a testing mock
type MockTxPool struct {
	txs []*tx.Tx
}

func NewMockTxPool() *MockTxPool {
	return &MockTxPool{
		txs: make([]*tx.Tx, 0),
	}
}
func (m *MockTxPool) SetSandbox(sandbox sandbox.Sandbox) {

}
func (m *MockTxPool) PendingTx(id crypto.Hash) *tx.Tx {
	for _, t := range m.txs {
		if t.ID().EqualsTo(id) {
			return t
		}
	}
	return nil
}

func (m *MockTxPool) HasTx(id crypto.Hash) bool {
	for _, t := range m.txs {
		if t.ID().EqualsTo(id) {
			return true
		}
	}
	return false
}

func (m *MockTxPool) Size() int {
	return len(m.txs)
}

func (m *MockTxPool) Fingerprint() string {
	return ""
}

func (m *MockTxPool) AppendTxs(txs []*tx.Tx) {
	for _, t := range txs {
		m.txs = append(m.txs, t)
	}
}

func (m *MockTxPool) AppendTx(t *tx.Tx) error {
	m.txs = append(m.txs, t)
	return nil
}
func (m *MockTxPool) AppendTxAndBroadcast(t *tx.Tx) error {
	m.txs = append(m.txs, t)
	return nil
}

func (m *MockTxPool) RemoveTx(hash crypto.Hash) {
	// This pools is shared between different instances
	// Lets keep txs then
	//delete(m.txs, hash)
}

func (m *MockTxPool) AllTransactions() []*tx.Tx {
	return m.txs
}
