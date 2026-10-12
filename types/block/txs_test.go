package block_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Equal(t, trx2, txs.Get(1))
}

func TestSubsidy(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blk, _ := ts.GenerateTestBlock(ts.RandHeight())
	assert.Equal(t, blk.Transactions()[0], blk.Transactions().Subsidy())
}

// The merkle root of the transactions can be the same for two different lists of transactions
// (see `simplemerkle`): [a, b, c] and [a, b, c, c]. The hash of the block should not.
// Otherwise, a block with a duplicated transaction would share the hash (and the certificate)
// of a valid block.
func TestBlockHashWithDuplicatedLastTransaction(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx1 := ts.GenerateTestTransferTx()
	trx2 := ts.GenerateTestTransferTx()
	trx3 := ts.GenerateTestTransferTx()
	txs := block.Txs{trx1, trx2, trx3}
	dupTxs := block.Txs{trx1, trx2, trx3, trx3}

	require.Equal(t, txs.Root(), dupTxs.Root(), "the premise of this test is the same merkle root")

	blk, _ := ts.GenerateTestBlock(ts.RandHeight(), testsuite.BlockWithTransactions(txs))
	dupBlk := block.MakeBlock(blk.Header().Version(), blk.Header().Time(), dupTxs,
		blk.Header().PrevBlockHash(), blk.Header().StateRoot(), blk.PrevCertificate(),
		blk.Header().SortitionSeed(), blk.Header().ProposerAddress())

	assert.NotEqual(t, blk.Hash(), dupBlk.Hash())
}
