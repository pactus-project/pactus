package lastinfo

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

// The best way to test this module is by writing test code in the `state.CommitBlock` function
// to restore the state after each commit.
//
// Testing this part is challenging. :(

type testData struct {
	*testsuite.TestSuite

	store    *store.MockStore
	lastInfo *LastInfo
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	store := store.MockingStore(ts)
	lastInfo := NewLastInfo()

	pub0, _ := ts.RandBLSKeyPair()
	pub1, _ := ts.RandBLSKeyPair()
	pub2, _ := ts.RandBLSKeyPair()
	pub3, _ := ts.RandBLSKeyPair()
	pub4, prv4 := ts.RandBLSKeyPair()
	signer := crypto.NewSigner(prv4)

	val0 := validator.NewValidator(pub0, 0)
	val1 := validator.NewValidator(pub1, 1)
	val2 := validator.NewValidator(pub2, 2)
	val3 := validator.NewValidator(pub3, 3)
	val4 := validator.NewValidator(pub4, 4)

	val0.AddToStake(100)
	val1.AddToStake(100)
	val2.AddToStake(100)
	val3.AddToStake(100)
	val4.AddToStake(100)

	val0.UpdateLastSortitionHeight(0)
	val1.UpdateLastSortitionHeight(0)
	val2.UpdateLastSortitionHeight(0)
	val3.UpdateLastSortitionHeight(0)
	val4.UpdateLastSortitionHeight(100)

	store.UpdateValidator(val0)
	store.UpdateValidator(val1)
	store.UpdateValidator(val2)
	store.UpdateValidator(val3)
	store.UpdateValidator(val4)

	// Last block
	committers := []int32{0, 1, 2, 3}
	trx := tx.NewSortitionTx(ts.RandStamp(), 1, pub4.Address(), ts.RandProof())
	signer.SignMsg(trx)
	prevHash := ts.RandHash()
	prevCert := ts.GenerateTestCertificate()
	lastHeight := ts.RandHeight()
	lastSeed := ts.RandSeed()
	lastBlock := block.MakeBlock(1, util.Now(), block.Txs{trx},
		prevHash,
		ts.RandHash(),
		prevCert, lastSeed, val2.Address())

	sig := signer.SignData([]byte("fatdog"))
	lastCert := certificate.NewCertificate(lastHeight, 0, committers, []int32{}, sig.(*bls.Signature))
	store.SaveBlock(lastHeight, lastBlock, lastCert)
	assert.Equal(t, store.LastHeight, lastHeight)

	lastInfo.UpdateSortitionSeed(lastSeed)
	lastInfo.UpdateBlockHeight(lastHeight)
	lastInfo.UpdateBlockHash(lastBlock.Hash())
	lastInfo.UpdateCertificate(lastCert)
	lastInfo.UpdateBlockTime(lastBlock.Header().Time())
	lastInfo.UpdateValidators([]*validator.Validator{val0, val1, val2, val3})

	return &testData{
		TestSuite: ts,
		store:     store,
		lastInfo:  lastInfo,
	}
}

func TestRestoreCommittee(t *testing.T) {
	td := setup(t)

	li := NewLastInfo()

	cmt, err := li.RestoreLastInfo(td.store, 4)
	assert.NoError(t, err)

	val0, _ := td.store.ValidatorByNumber(0)
	val1, _ := td.store.ValidatorByNumber(1)
	val2, _ := td.store.ValidatorByNumber(2)
	val3, _ := td.store.ValidatorByNumber(3)

	assert.Equal(t, td.lastInfo.SortitionSeed(), li.SortitionSeed())
	assert.Equal(t, td.lastInfo.BlockHeight(), li.BlockHeight())
	assert.Equal(t, td.lastInfo.BlockHash(), li.BlockHash())
	assert.Equal(t, td.lastInfo.Certificate().Hash(), li.Certificate().Hash())
	assert.Equal(t, td.lastInfo.BlockTime(), li.BlockTime())
	assert.Equal(t, td.lastInfo.Validators(), []*validator.Validator{val0, val1, val2, val3})
	assert.Equal(t, cmt.Committers(), []int32{1, 4, 2, 3})
}

func TestRestoreFailed(t *testing.T) {
	td := setup(t)

	t.Run("Unable to get validator from store", func(t *testing.T) {
		setup(t)

		li := NewLastInfo()

		td.store.Validators = make(map[crypto.Address]validator.Validator) // Reset Validators
		_, err := li.RestoreLastInfo(td.store, 4)
		assert.Error(t, err)
	})

	t.Run("Unable to get block from store", func(t *testing.T) {
		setup(t)

		li := NewLastInfo()

		td.store.Blocks = make(map[uint32]block.Block) // Reset Blocks
		_, err := li.RestoreLastInfo(td.store, 4)
		assert.Error(t, err)
	})
}
