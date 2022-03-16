package block

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/tx"
)

func TestTxsMerkle(t *testing.T) {
	txs := NewTxs()
	trx1, _ := tx.GenerateTestSendTx()
	trx2, _ := tx.GenerateTestSendTx()
	txs.Append(trx1)
	merkle := txs.Root()
	assert.Equal(t, merkle, trx1.ID())

	txs.Append(trx2)
	merkle = txs.Root()
	data := make([]byte, 64)
	copy(data[:32], trx1.ID().RawBytes())
	copy(data[32:], trx2.ID().RawBytes())
	assert.Equal(t, merkle, hash.CalcHash(data))
}

func TestAppendPrependRemove(t *testing.T) {
	txs := NewTxs()
	trx1, _ := tx.GenerateTestSendTx()
	trx2, _ := tx.GenerateTestSendTx()
	trx3, _ := tx.GenerateTestSendTx()
	trx4, _ := tx.GenerateTestSendTx()
	trx5, _ := tx.GenerateTestSendTx()
	txs.Append(trx2)
	txs.Append(trx3)
	txs.Prepend(trx1)
	txs.Append(trx5)
	txs.Append(trx4)
	txs.Remove(3)

	assert.Equal(t, txs, Txs{trx1, trx2, trx3, trx4})
}

func TestIsEmpty(t *testing.T) {
	txs := NewTxs()
	assert.True(t, txs.IsEmpty())

	trx, _ := tx.GenerateTestSendTx()
	txs.Append(trx)
	assert.False(t, txs.IsEmpty())
}
