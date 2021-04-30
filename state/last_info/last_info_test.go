package last_info

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

// The best way to test this module, is writting test code in `state.CommitBlock` function
// And try to sync with the main net and restore state after each commit.
// The restored info should be exactly same as currect info.
//
// Testing this module is not easy ;(

var tStore *store.MockStore
var tLastInfo *LastInfo
var tSortition *sortition.Sortition

func setup(t *testing.T) {
	tStore = store.MockingStore()
	tLastInfo = NewLastInfo(tStore)
	tSortition = sortition.NewSortition()

	setSortitionParams := func(hash crypto.Hash, seed sortition.Seed, committers []int) {
		totalStake := int64(0)
		tStore.IterateValidators(func(v *validator.Validator) (stop bool) {
			totalStake += v.Stake()
			return false
		})

		committeeStake := int64(0)
		for _, num := range committers {
			v, _ := tStore.ValidatorByNumber(num)
			committeeStake += v.Stake()
		}

		poolStake := totalStake - committeeStake
		fmt.Printf("hash: %v, total stake: %d, committee stake: %v, pool stake: %v\n",
			hash.Fingerprint(), totalStake, committeeStake, poolStake)

		tSortition.SetParams(hash, seed, poolStake)
	}

	signer0 := crypto.GenerateTestSigner()
	signer1 := crypto.GenerateTestSigner()
	signer2 := crypto.GenerateTestSigner()
	signer3 := crypto.GenerateTestSigner()
	signer4 := crypto.GenerateTestSigner()
	signer5 := crypto.GenerateTestSigner()
	signer6 := crypto.GenerateTestSigner()
	signer7 := crypto.GenerateTestSigner()

	val0 := validator.NewValidator(signer0.PublicKey(), 0, 0)
	val1 := validator.NewValidator(signer1.PublicKey(), 1, 0)
	val2 := validator.NewValidator(signer2.PublicKey(), 2, 0)
	val3 := validator.NewValidator(signer3.PublicKey(), 3, 0)

	val0.AddToStake(1000)
	val1.AddToStake(2000)
	val2.AddToStake(3000)
	val3.AddToStake(4000)

	tStore.UpdateValidator(val0)
	tStore.UpdateValidator(val1)
	tStore.UpdateValidator(val2)
	tStore.UpdateValidator(val3)

	sig := signer1.SignData([]byte("dummy"))

	committers1 := []int{0, 1, 2, 3}
	// Block 1
	trx1, _ := tx.GenerateTestSendTx()
	ids1 := block.NewTxIDs()
	ids1.Append(trx1.ID())
	seed1 := sortition.GenerateRandomSeed()
	block1 := block.MakeBlock(1, util.Now(), ids1,
		crypto.UndefHash,
		crypto.GenerateTestHash(),
		nil, seed1, val1.Address())

	cert1 := block.NewCertificate(block1.Hash(), 0, committers1, []int{}, sig)
	tStore.SaveBlock(1, block1)
	tStore.SaveTransaction(trx1)
	committers2 := []int{0, 1, 2, 3}
	setSortitionParams(block1.Hash(), seed1, committers2)

	// Block 2
	val4 := validator.NewValidator(signer4.PublicKey(), 4, 3)
	val4.AddToStake(4000)
	tStore.UpdateValidator(val4)
	trx2 := tx.NewBondTx(block1.Hash(), 1, signer1.Address(), val4.PublicKey(), 4000, 4000, "")
	ids2 := block.NewTxIDs()
	ids2.Append(trx2.ID())
	seed2 := sortition.GenerateRandomSeed()
	block2 := block.MakeBlock(1, util.Now(), ids2,
		block1.Hash(),
		crypto.GenerateTestHash(),
		cert1, seed2, val1.Address())

	cert2 := block.NewCertificate(block2.Hash(), 0, committers2, []int{}, sig)
	tStore.SaveBlock(2, block2)
	tStore.SaveTransaction(trx2)
	committers3 := []int{0, 1, 2, 3}
	setSortitionParams(block2.Hash(), seed2, committers3)

	// Block 3
	val5 := validator.NewValidator(signer5.PublicKey(), 5, 3)
	val5.AddToStake(5000)
	val6 := validator.NewValidator(signer6.PublicKey(), 6, 3)
	val5.AddToStake(5000)
	val4.UpdateLastJoinedHeight(3)
	tStore.UpdateValidator(val4)
	tStore.UpdateValidator(val5)
	tStore.UpdateValidator(val6)
	trx31 := tx.NewBondTx(block2.Hash(), 1, signer1.Address(), val5.PublicKey(), 5000, 5000, "")
	trx32 := tx.NewBondTx(block2.Hash(), 2, signer1.Address(), val6.PublicKey(), 5000, 5000, "")
	trx33 := tx.NewSortitionTx(block2.Hash(), 1, signer4.Address(), sortition.GenerateRandomProof())
	ids3 := block.NewTxIDs()
	ids3.Append(trx31.ID())
	ids3.Append(trx32.ID())
	ids3.Append(trx33.ID())
	seed3 := sortition.GenerateRandomSeed()
	block3 := block.MakeBlock(1, util.Now(), ids3,
		block2.Hash(),
		crypto.GenerateTestHash(),
		cert2, seed3, val1.Address())

	cert3 := block.NewCertificate(block3.Hash(), 0, committers2, []int{}, sig)
	tStore.SaveBlock(3, block3)
	tStore.SaveTransaction(trx31)
	tStore.SaveTransaction(trx32)
	tStore.SaveTransaction(trx33)
	committers4 := []int{4, 1, 2, 3}
	setSortitionParams(block3.Hash(), seed3, committers4)

	// Block 4
	val0.AddToStake(5000)
	val6.UpdateLastJoinedHeight(4)
	tStore.UpdateValidator(val0)
	tStore.UpdateValidator(val6)
	trx41 := tx.NewBondTx(block3.Hash(), 1, signer1.Address(), val0.PublicKey(), 5000, 5000, "")
	trx42 := tx.NewSortitionTx(block3.Hash(), 1, signer6.Address(), sortition.GenerateRandomProof())
	ids4 := block.NewTxIDs()
	ids4.Append(trx41.ID())
	ids4.Append(trx42.ID())
	seed4 := sortition.GenerateRandomSeed()
	block4 := block.MakeBlock(1, util.Now(), ids4,
		block3.Hash(),
		crypto.GenerateTestHash(),
		cert3, seed4, val1.Address())

	cert4 := block.NewCertificate(block4.Hash(), 0, committers4, []int{}, sig)
	tStore.SaveBlock(4, block4)
	tStore.SaveTransaction(trx41)
	tStore.SaveTransaction(trx42)
	committers5 := []int{4, 6, 2, 3}
	setSortitionParams(block4.Hash(), seed4, committers5)

	// Block 5
	val7 := validator.NewValidator(signer7.PublicKey(), 7, 5)
	val7.AddToStake(7000)
	val5.UpdateLastJoinedHeight(5)
	tStore.UpdateValidator(val5)
	tStore.UpdateValidator(val7)
	trx51 := tx.NewBondTx(block3.Hash(), 1, signer1.Address(), val7.PublicKey(), 7000, 7000, "")
	trx52 := tx.NewSortitionTx(block3.Hash(), 1, signer5.Address(), sortition.GenerateRandomProof())
	ids5 := block.NewTxIDs()
	ids5.Append(trx51.ID())
	ids5.Append(trx52.ID())
	seed5 := sortition.GenerateRandomSeed()
	block5 := block.MakeBlock(1, util.Now(), ids5,
		block4.Hash(),
		crypto.GenerateTestHash(),
		cert4, seed5, val1.Address())

	cert5 := block.NewCertificate(block5.Hash(), 0, committers5, []int{}, sig)
	tStore.SaveBlock(5, block5)
	tStore.SaveTransaction(trx51)
	tStore.SaveTransaction(trx52)
	committers6 := []int{5, 4, 6, 3}
	setSortitionParams(block5.Hash(), seed5, committers6)

	tLastInfo.SetSortitionSeed(seed5)
	tLastInfo.SetBlockHeight(5)
	tLastInfo.SetBlockHash(block5.Hash())
	tLastInfo.SetCertificate(cert5)
	tLastInfo.SetBlockTime(block5.Header().Time())
	tLastInfo.SaveLastInfo()
}

func TestRestore(t *testing.T) {
	setup(t)

	li := NewLastInfo(tStore)
	srt := sortition.NewSortition()

	cmt, err := li.RestoreLastInfo(4, srt)
	assert.NoError(t, err)

	assert.Equal(t, tLastInfo.SortitionSeed(), li.SortitionSeed())
	assert.Equal(t, tLastInfo.BlockHeight(), li.BlockHeight())
	assert.Equal(t, tLastInfo.BlockHash(), li.BlockHash())
	assert.Equal(t, tLastInfo.Certificate().Hash(), li.Certificate().Hash())
	assert.Equal(t, tLastInfo.BlockTime(), li.BlockTime())
	assert.Equal(t, cmt.Committers(), []int{5, 4, 6, 3})

	for i := 1; i < 6; i++ {
		b, _ := tStore.Block(i)
		seed1, stake1 := srt.GetParams(b.Hash())
		seed2, stake2 := tSortition.GetParams(b.Hash())

		assert.Equal(t, seed1, seed2, "Invalid seed for block %v", i)
		assert.Equal(t, stake1, stake2, "Invalid stake for block %v", i)
	}
}

func TestRestoreFailed(t *testing.T) {
	t.Run("Unable to get validator from store", func(t *testing.T) {
		setup(t)

		li := NewLastInfo(tStore)
		srt := sortition.NewSortition()

		tStore.Validators = make(map[crypto.Address]validator.Validator) // Reset Validators
		_, err := li.RestoreLastInfo(4, srt)
		assert.Error(t, err)
	})

	t.Run("Unable to get transaction from store", func(t *testing.T) {
		setup(t)

		li := NewLastInfo(tStore)
		srt := sortition.NewSortition()

		tStore.Transactions = make(map[crypto.Hash]tx.Tx) // Reset transactions
		_, err := li.RestoreLastInfo(4, srt)
		assert.Error(t, err)
	})

	t.Run("Unable to get block from store", func(t *testing.T) {
		setup(t)

		li := NewLastInfo(tStore)
		srt := sortition.NewSortition()

		tStore.Blocks = make(map[int]*block.Block) // Reset Blocks
		_, err := li.RestoreLastInfo(4, srt)
		assert.Error(t, err)
	})

}
