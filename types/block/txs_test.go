package block_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestTxsMerkle(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	txs := block.NewTxs()
	trx1, _ := ts.GenerateTestTransferTx()
	trx2, _ := ts.GenerateTestTransferTx()
	txs.Append(trx1)
	merkle := txs.Root()
	assert.Equal(t, merkle, trx1.ID())

	txs.Append(trx2)
	merkle = txs.Root()
	data := make([]byte, 64)
	copy(data[:32], trx1.ID().Bytes())
	copy(data[32:], trx2.ID().Bytes())
	assert.Equal(t, merkle, hash.CalcHash(data))
}

func TestAppendPrependRemove(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	txs := block.NewTxs()
	trx1, _ := ts.GenerateTestTransferTx()
	trx2, _ := ts.GenerateTestTransferTx()
	trx3, _ := ts.GenerateTestTransferTx()
	trx4, _ := ts.GenerateTestTransferTx()
	trx5, _ := ts.GenerateTestTransferTx()
	txs.Append(trx2)
	txs.Append(trx3)
	txs.Prepend(trx1)
	txs.Append(trx5)
	txs.Append(trx4)
	txs.Remove(3)

	assert.Equal(t, txs, block.Txs{trx1, trx2, trx3, trx4})
}

func TestIsEmpty(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	txs := block.NewTxs()
	assert.True(t, txs.IsEmpty())

	trx, _ := ts.GenerateTestTransferTx()
	txs.Append(trx)
	assert.False(t, txs.IsEmpty())
}

func TestGetTransaction(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	txs := block.NewTxs()
	trx1, _ := ts.GenerateTestTransferTx()
	trx2, _ := ts.GenerateTestTransferTx()
	txs.Append(trx1)
	txs.Append(trx2)
	assert.Equal(t, trx1, txs.Get(0))
}
