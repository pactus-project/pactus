package txpool

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
)

var _ TxPool = &MockTxPool{}

// MockTxPool is a testing mock.
type MockTxPool struct {
	Txs []*tx.Tx
}

func MockingTxPool() *MockTxPool {
	return &MockTxPool{
		Txs: make([]*tx.Tx, 0),
	}
}
func (m *MockTxPool) SetNewSandboxAndRecheck(_ sandbox.Sandbox) {}
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

func (m *MockTxPool) Size() int {
	return len(m.Txs)
}

func (m *MockTxPool) String() string {
	return ""
}

func (m *MockTxPool) AppendTx(trx *tx.Tx) error {
	m.Txs = append(m.Txs, trx)
	return nil
}
func (m *MockTxPool) AppendTxAndBroadcast(trx *tx.Tx) error {
	m.Txs = append(m.Txs, trx)
	return nil
}

func (m *MockTxPool) RemoveTx(_ hash.Hash) {
	// This test pools is shared between different test objects
	//delete(m.Txs, id)
}

func (m *MockTxPool) PrepareBlockTransactions() block.Txs {
	txs := make([]*tx.Tx, m.Size())
	copy(txs, m.Txs)
	return txs
}
