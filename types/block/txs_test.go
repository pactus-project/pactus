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
	trx1 := ts.GenerateTestTransferTx()
	trx2 := ts.GenerateTestTransferTx()
	txs.Append(trx1)
	merkle := txs.Root()
	assert.Equal(t, trx1.ID(), merkle)

	txs.Append(trx2)
	merkle = txs.Root()
	data := make([]byte, 64)
	copy(data[:32], trx1.ID().Bytes())
	copy(data[32:], trx2.ID().Bytes())
	assert.Equal(t, hash.CalcHash(data), merkle)
}

func TestAppendPrependRemove(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	txs := block.NewTxs()
	trx1 := ts.GenerateTestTransferTx()
	trx2 := ts.GenerateTestTransferTx()
	trx3 := ts.GenerateTestTransferTx()
	trx4 := ts.GenerateTestTransferTx()
	trx5 := ts.GenerateTestTransferTx()
	txs.Append(trx2)
	txs.Append(trx3)
	txs.Prepend(trx1)
	txs.Append(trx5)
	txs.Append(trx4)
	txs.Remove(3)

	assert.Equal(t, block.Txs{trx1, trx2, trx3, trx4}, txs)
}

func TestIsEmpty(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	txs := block.NewTxs()
	assert.True(t, txs.IsEmpty())

	trx := ts.GenerateTestTransferTx()
	txs.Append(trx)
	assert.False(t, txs.IsEmpty())
}

func TestGetTransaction(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	txs := block.NewTxs()
	trx1 := ts.GenerateTestTransferTx()
	trx2 := ts.GenerateTestTransferTx()
	txs.Append(trx1)
	txs.Append(trx2)
	assert.Equal(t, trx1, txs.Get(0))
}

func TestSubsidy(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blk, _ := ts.GenerateTestBlock(ts.RandHeight())
	assert.Equal(t, blk.Transactions()[0], blk.Transactions().Subsidy())
}
