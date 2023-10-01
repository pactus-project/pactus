package store

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
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
		b := td.GenerateTestBlock()
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

	lastHeight, _ := td.store.LastCertificate()
	committedBlock, err := td.store.Block(lastHeight)
	assert.NoError(t, err)
	assert.Equal(t, lastHeight, committedBlock.Height)
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

		if addr.IsAccountAddress() {
			assert.Equal(t, pub.AccountAddress(), addr)
		} else if addr.IsValidatorAddress() {
			assert.Equal(t, pub.ValidatorAddress(), addr)
		}
	}

	pub, found := td.store.PublicKey(td.RandValAddress())
	assert.Error(t, found)
	assert.Nil(t, pub)
}

func TestStrippedPublicKey(t *testing.T) {
	td := setup(t)

	// Find a public key that we have already indexed in the database.
	committedBlock1, _ := td.store.Block(1)
	blk1, _ := committedBlock1.ToBlock()
	trx0PubKey := blk1.Transactions()[0].PublicKey()
	assert.NotNil(t, trx0PubKey)
	knownPubKey := trx0PubKey.(*bls.PublicKey)

	lastHeight, _ := td.store.LastCertificate()
	lockTime := lastHeight
	randPubkey, _ := td.RandBLSKeyPair()

	trx0 := tx.NewTransferTx(lockTime, knownPubKey.AccountAddress(), td.RandAccAddress(), 1, 1, "")
	trx1 := tx.NewTransferTx(lockTime, randPubkey.AccountAddress(), td.RandAccAddress(), 1, 1, "")
	trx2 := tx.NewTransferTx(lockTime, randPubkey.AccountAddress(), td.RandAccAddress(), 1, 1, "")

	trx0.SetSignature(td.RandBLSSignature())
	trx1.SetSignature(td.RandBLSSignature())
	trx2.SetSignature(td.RandBLSSignature())

	trx0.StripPublicKey()
	trx1.SetPublicKey(randPubkey)
	trx2.StripPublicKey()

	tests := []struct {
		trx    *tx.Tx
		failed bool
	}{
		{trx0, false}, // indexed public key and stripped
		{trx1, false}, // not stripped
		{trx2, true},  // unknown public key and stripped
	}

	for _, test := range tests {
		trxs := block.Txs{test.trx}

		// Make a block
		prevCert := td.GenerateTestCertificate()
		blk := block.MakeBlock(1, util.Now(), trxs, td.RandHash(), td.RandHash(),
			prevCert, td.RandSeed(), td.RandValAddress())

		trxData, _ := test.trx.Bytes()
		blkData, _ := blk.Bytes()

		committedTrx := CommittedTx{
			store:  td.store,
			TxID:   test.trx.ID(),
			Height: lastHeight + 1,
			Data:   trxData,
		}
		committedBlock := CommittedBlock{
			store:     td.store,
			BlockHash: blk.Hash(),
			Height:    lastHeight + 1,
			Data:      blkData,
		}

		//
		if test.failed {
			_, err := committedBlock.ToBlock()
			assert.ErrorIs(t, err, PublicKeyNotFoundError{
				Address: test.trx.Payload().Signer(),
			})

			_, err = committedTrx.ToTx()
			assert.ErrorIs(t, err, PublicKeyNotFoundError{
				Address: test.trx.Payload().Signer(),
			})
		} else {
			_, err := committedBlock.ToBlock()
			assert.NoError(t, err)

			_, err = committedTrx.ToTx()
			assert.NoError(t, err)
		}
	}
}
