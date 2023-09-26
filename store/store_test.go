package store

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	store *store
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	conf := &Config{
		Path: util.TempDirPath(),
	}
	s, err := NewStore(conf)
	require.NoError(t, err)

	td := &testData{
		TestSuite: ts,
		store:     s.(*store),
	}

	td.saveTestBlocks(t, 10)

	return td
}

func (td *testData) saveTestBlocks(t *testing.T, num int) {
	t.Helper()

	lastHeight, _ := td.store.LastCertificate()
	for i := 0; i < num; i++ {
		b := td.GenerateTestBlock(nil)
		c := td.GenerateTestCertificate()

		td.store.SaveBlock(lastHeight+uint32(i+1), b, c)
		assert.NoError(t, td.store.WriteBatch())
	}
}

func TestBlockHash(t *testing.T) {
	td := setup(t)

	sb, _ := td.store.Block(1)

	assert.Equal(t, td.store.BlockHash(0), hash.UndefHash)
	assert.Equal(t, td.store.BlockHash(util.MaxUint32), hash.UndefHash)
	assert.Equal(t, td.store.BlockHash(1), sb.BlockHash)
}

func TestBlockHeight(t *testing.T) {
	td := setup(t)

	sb, _ := td.store.Block(1)

	assert.Equal(t, td.store.BlockHeight(hash.UndefHash), uint32(0))
	assert.Equal(t, td.store.BlockHeight(td.RandHash()), uint32(0))
	assert.Equal(t, td.store.BlockHeight(sb.BlockHash), uint32(1))
}

func TestUnknownTransactionID(t *testing.T) {
	td := setup(t)

	tx, err := td.store.Transaction(td.RandHash())
	assert.Error(t, err)
	assert.Nil(t, tx)
}

func TestWriteAndClosePeacefully(t *testing.T) {
	td := setup(t)

	// After closing db, we should not crash
	assert.NoError(t, td.store.Close())
	assert.Error(t, td.store.WriteBatch())
}

func TestRetrieveBlockAndTransactions(t *testing.T) {
	td := setup(t)

	height, _ := td.store.LastCertificate()
	committedBlock, err := td.store.Block(height)
	assert.NoError(t, err)
	assert.Equal(t, height, committedBlock.Height)
	block, _ := committedBlock.ToBlock()
	for _, trx := range block.Transactions() {
		committedTx, err := td.store.Transaction(trx.ID())
		assert.NoError(t, err)
		assert.Equal(t, committedTx.TxID, trx.ID())
		assert.Equal(t, committedTx.BlockTime, block.Header().UnixTime())
		trx2, _ := committedTx.ToTx()
		assert.Equal(t, trx2.ID(), trx.ID())
	}
}

func TestIndexingPublicKeys(t *testing.T) {
	td := setup(t)

	committedBlock, _ := td.store.Block(1)
	blk, _ := committedBlock.ToBlock()
	for _, trx := range blk.Transactions() {
		addr := trx.Payload().Signer()
		pub, found := td.store.PublicKey(addr)

		assert.NoError(t, found)
		assert.Equal(t, pub.Address(), addr)
	}

	pub, found := td.store.PublicKey(td.RandAddress())
	assert.Error(t, found)
	assert.Nil(t, pub)
}

func TestCommittedBlockToBlock(t *testing.T) {
	td := setup(t)

	// Use a tricky way to save transactions from the first block again.
	committedBlock1, _ := td.store.Block(1)
	committedBlock2, _ := td.store.Block(2)
	blk1, _ := committedBlock1.ToBlock()
	blk2, _ := committedBlock2.ToBlock()
	td.store.SaveBlock(11, blk1, blk2.PrevCertificate())
	err := td.store.WriteBatch()
	assert.NoError(t, err)

	// Ensure that the committed block can obtain the public key.
	committedBlock11, err := td.store.Block(11)
	assert.NoError(t, err)

	blk11, err := committedBlock11.ToBlock()
	assert.NoError(t, err)

	err = blk11.BasicCheck()
	assert.NoError(t, err)

	// Ensure that the committed transactions can obtain the public key.
	committedTrx, err := td.store.Transaction(blk11.Transactions()[0].ID())
	assert.NoError(t, err)

	trx, err := committedTrx.ToTx()
	assert.NoError(t, err)

	err = trx.BasicCheck()
	assert.NoError(t, err)
}
