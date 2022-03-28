package lastinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

// The best way to test this module, is writing a test code in `state.CommitBlock` function
// to restore state after each commit.
//
// Testing this part is not easy ;(

var tStore *store.MockStore
var tLastInfo *LastInfo

func setup(t *testing.T) {
	tStore = store.MockingStore()
	tLastInfo = NewLastInfo(tStore)

	pub0, _ := bls.GenerateTestKeyPair()
	pub1, _ := bls.GenerateTestKeyPair()
	pub2, _ := bls.GenerateTestKeyPair()
	pub3, _ := bls.GenerateTestKeyPair()
	pub4, prv4 := bls.GenerateTestKeyPair()
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

	val0.UpdateLastJoinedHeight(0)
	val1.UpdateLastJoinedHeight(0)
	val2.UpdateLastJoinedHeight(0)
	val3.UpdateLastJoinedHeight(0)
	val4.UpdateLastJoinedHeight(100)

	tStore.UpdateValidator(val0)
	tStore.UpdateValidator(val1)
	tStore.UpdateValidator(val2)
	tStore.UpdateValidator(val3)
	tStore.UpdateValidator(val4)

	// Last block
	committers := []int32{0, 1, 2, 3}
	trx := tx.NewSortitionTx(hash.GenerateTestStamp(), 1, pub4.Address(), sortition.GenerateRandomProof())
	signer.SignMsg(trx)
	prevHash := hash.GenerateTestHash()
	prevCert := block.GenerateTestCertificate(prevHash)
	lastHeight := util.RandInt32(100000)
	lastSeed := sortition.GenerateRandomSeed()
	lastBlock := block.MakeBlock(1, util.Now(), block.Txs{trx},
		prevHash,
		hash.GenerateTestHash(),
		prevCert, lastSeed, val2.Address())

	sig := signer.SignData([]byte("fatdog"))
	lastCert := block.NewCertificate(0, committers, []int32{}, sig.(*bls.Signature))
	tStore.SaveBlock(lastHeight, lastBlock, lastCert)
	assert.Equal(t, tStore.LastHeight, lastHeight)

	tLastInfo.SetSortitionSeed(lastSeed)
	tLastInfo.SetBlockHeight(lastHeight)
	tLastInfo.SetBlockHash(lastBlock.Hash())
	tLastInfo.SetCertificate(lastCert)
	tLastInfo.SetBlockTime(lastBlock.Header().Time())
}

func TestRestoreCommittee(t *testing.T) {
	setup(t)

	li := NewLastInfo(tStore)

	cmt, err := li.RestoreLastInfo(4)
	assert.NoError(t, err)

	assert.Equal(t, tLastInfo.SortitionSeed(), li.SortitionSeed())
	assert.Equal(t, tLastInfo.BlockHeight(), li.BlockHeight())
	assert.Equal(t, tLastInfo.BlockHash(), li.BlockHash())
	assert.Equal(t, tLastInfo.Certificate().Hash(), li.Certificate().Hash())
	assert.Equal(t, tLastInfo.BlockTime(), li.BlockTime())
	assert.Equal(t, cmt.Committers(), []int32{1, 4, 2, 3})
}

func TestRestoreFailed(t *testing.T) {
	t.Run("Unable to get validator from store", func(t *testing.T) {
		setup(t)

		li := NewLastInfo(tStore)

		tStore.Validators = make(map[crypto.Address]validator.Validator) // Reset Validators
		_, err := li.RestoreLastInfo(4)
		assert.Error(t, err)
	})

	t.Run("Unable to get block from store", func(t *testing.T) {
		setup(t)

		li := NewLastInfo(tStore)

		tStore.Blocks = make(map[int32]block.Block) // Reset Blocks
		_, err := li.RestoreLastInfo(4)
		assert.Error(t, err)
	})

}
