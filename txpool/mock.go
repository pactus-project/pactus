package txpool

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
)

var _ TxPool = &MockTxPool{}

// MockTxPool is a testing mock
type MockTxPool struct {
	Txs []*tx.Tx
}

func MockingTxPool() *MockTxPool {
	return &MockTxPool{
		Txs: make([]*tx.Tx, 0),
	}
}
func (m *MockTxPool) SetSandbox(sandbox sandbox.Sandbox) {

}
func (m *MockTxPool) PendingTx(id tx.ID) *tx.Tx {
	for _, t := range m.Txs {
		if t.ID().EqualsTo(id) {
			return t
		}
	}
	return nil
}

func (m *MockTxPool) QueryTx(id tx.ID) *tx.Tx {
	return m.PendingTx(id)
}

func (m *MockTxPool) HasTx(id tx.ID) bool {
	for _, t := range m.Txs {
		if t.ID().EqualsTo(id) {
			return true
		}
	}
	return false
}

func (m *MockTxPool) Recheck() {}
func (m *MockTxPool) Size() int {
	return len(m.Txs)
}

func (m *MockTxPool) Fingerprint() string {
	return ""
}

func (m *MockTxPool) AppendTx(t *tx.Tx) error {
	m.Txs = append(m.Txs, t)
	return nil
}
func (m *MockTxPool) AppendTxAndBroadcast(t *tx.Tx) error {
	m.Txs = append(m.Txs, t)
	return nil
}

func (m *MockTxPool) RemoveTx(hash crypto.Hash) {
	// This pools is shared between different instances
	// Lets keep txs then
	//delete(m.txs, hash)
}

func (m *MockTxPool) AllTransactions() []*tx.Tx {
	return m.Txs
}
