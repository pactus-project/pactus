package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/libs/linkedmap"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/validator"
)

var tValAddress [6]crypto.Address
var tSortitions *sortition.Sortition
var tCommittee *committee.Committee
var tStore *store.MockStore
var tSandbox *sandbox

func init() {
	logger.InitLogger(logger.TestConfig())
}

func setup(t *testing.T) {
	var err error
	tStore = store.MockingStore()
	params := param.DefaultParams()
	params.TransactionToLiveInterval = 64
	latestBlocks := linkedmap.NewLinkedMap(params.TransactionToLiveInterval)

	lastHeight := 124
	for i := 0; i <= lastHeight; i++ {
		b, _ := block.GenerateTestBlock(nil, nil)
		tStore.SaveBlock(i+1, b)
		latestBlocks.PushBack(b.Stamp(), &BlockInfo{height: i + 1, hash: b.Hash()})
	}

	pub1, _ := bls.GenerateTestKeyPair()
	pub2, _ := bls.GenerateTestKeyPair()
	pub3, _ := bls.GenerateTestKeyPair()
	pub4, _ := bls.GenerateTestKeyPair()
	pub5, _ := bls.GenerateTestKeyPair()
	pub6, _ := bls.GenerateTestKeyPair()

	tValAddress[0] = pub1.Address()
	tValAddress[1] = pub2.Address()
	tValAddress[2] = pub3.Address()
	tValAddress[3] = pub4.Address()
	tValAddress[4] = pub5.Address()
	tValAddress[5] = pub6.Address()

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14)

	val1 := validator.NewValidator(pub1, 0)
	val2 := validator.NewValidator(pub2, 1)
	val3 := validator.NewValidator(pub3, 2)
	val4 := validator.NewValidator(pub4, 3)
	val5 := validator.NewValidator(pub5, 4)
	val6 := validator.NewValidator(pub6, 5)

	val1.AddToStake(2000)
	val2.AddToStake(2500)
	val3.AddToStake(2000)
	val4.AddToStake(2500)
	val5.AddToStake(1000)
	val6.AddToStake(3000)

	tStore.UpdateAccount(acc)
	tStore.UpdateValidator(val1)
	tStore.UpdateValidator(val2)
	tStore.UpdateValidator(val3)
	tStore.UpdateValidator(val4)
	tStore.UpdateValidator(val5)
	tStore.UpdateValidator(val6)

	tSortitions = sortition.NewSortition()
	tCommittee, err = committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, tValAddress[0])
	assert.NoError(t, err)

	tSandbox = NewSandbox(tStore, params, latestBlocks, tSortitions, tCommittee)
	assert.Equal(t, tSandbox.MaxMemoLength(), params.MaximumMemoLength)
	assert.Equal(t, tSandbox.FeeFraction(), params.FeeFraction)
	assert.Equal(t, tSandbox.MinFee(), params.MinimumFee)
	assert.Equal(t, tSandbox.TransactionToLiveInterval(), params.TransactionToLiveInterval)
	assert.Equal(t, tSandbox.CommitteeSize(), params.CommitteeSize)
}

func TestAccountChange(t *testing.T) {
	setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr := crypto.GenerateTestAddress()
		assert.Nil(t, tSandbox.Account(invAddr))
	})

	t.Run("Retrieve an account from store, modify it and commit it", func(t *testing.T) {
		acc1, _ := account.GenerateTestAccount(888)
		tStore.UpdateAccount(acc1)

		acc1a := tSandbox.Account(acc1.Address())
		assert.Equal(t, acc1, acc1a)

		acc1a.IncSequence()
		acc1a.AddToBalance(acc1a.Balance() + 1)

		assert.False(t, tSandbox.accounts[acc1a.Address()].Updated)
		tSandbox.UpdateAccount(acc1a)
		assert.True(t, tSandbox.accounts[acc1a.Address()].Updated)
	})

	t.Run("Make new account", func(t *testing.T) {
		addr := crypto.GenerateTestAddress()
		acc2 := tSandbox.MakeNewAccount(addr)

		acc2.IncSequence()
		acc2.AddToBalance(acc2.Balance() + 1)

		tSandbox.UpdateAccount(acc2)
		acc22 := tSandbox.Account(acc2.Address())
		assert.Equal(t, acc2, acc22)
	})
}

func TestValidatorChange(t *testing.T) {
	setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr := crypto.GenerateTestAddress()
		assert.Nil(t, tSandbox.Validator(invAddr))
	})

	t.Run("Retrieve an validator from store, modify it and commit it", func(t *testing.T) {
		val1, _ := validator.GenerateTestValidator(888)
		tStore.UpdateValidator(val1)

		val1a := tSandbox.Validator(val1.Address())
		assert.Equal(t, val1.Hash(), val1a.Hash())

		val1a.IncSequence()
		val1a.AddToStake(val1a.Stake() + 1)

		assert.False(t, tSandbox.validators[val1a.Address()].Updated)
		tSandbox.UpdateValidator(val1a)
		assert.True(t, tSandbox.validators[val1a.Address()].Updated)
	})

	t.Run("Make new validator", func(t *testing.T) {
		pub, _ := bls.GenerateTestKeyPair()
		val2 := tSandbox.MakeNewValidator(pub)

		val2.IncSequence()
		val2.AddToStake(val2.Stake() + 1)

		tSandbox.UpdateValidator(val2)
		val22 := tSandbox.Validator(val2.Address())
		assert.Equal(t, val2, val22)
	})
}

// func TestAddValidatorToCommittee(t *testing.T) {
// 	setup(t)

// 	tSandbox.params.CommitteeSize = 4

// 	t.Run("Add unknown validator to the committee, Should returns error", func(t *testing.T) {
// 		assert.Error(t, tSandbox.JoinCommittee(crypto.GenerateTestAddress()))
// 	})

// 	t.Run("More than 1/3 stake, Should returns error", func(t *testing.T) {
// 		val6 := tSandbox.Validator(tValAddress[5]) // val_stake: 3000, total_stake: 9000

// 		assert.Equal(t, tSandbox.committee.TotalStake(), int64(9000), "Total stake should be 9000")
// 		assert.Equal(t, val6.Stake(), int64(3000), "Validator stake should be 3000")
// 		assert.Error(t, tSandbox.JoinCommittee(val6.Address()), "More than 1/3 of stake is going to change")
// 	})

// 	t.Run("Update validator and add to committee", func(t *testing.T) {
// 		val6 := tSandbox.Validator(tValAddress[4]) // val_stake: 1000, total_stake: 9000

// 		assert.Equal(t, val6.Stake(), int64(1000), "Validator stake should be 1000")
// 		assert.NoError(t, tSandbox.JoinCommittee(val6.Address()))
// 	})

// 	t.Run("Only one validator can enter committee per height", func(t *testing.T) {
// 		val6 := tSandbox.Validator(tValAddress[4]) // val_stake: 1000, total_stake: 9000

// 		assert.Error(t, tSandbox.JoinCommittee(val6.Address()))
// 	})
// }

func TestTotalAccountCounter(t *testing.T) {
	setup(t)

	t.Run("Should update total account counter", func(t *testing.T) {
		assert.Equal(t, tStore.TotalAccounts(), 1) // Sandbox has an account

		addr1 := crypto.GenerateTestAddress()
		addr2 := crypto.GenerateTestAddress()
		acc := tSandbox.MakeNewAccount(addr1)
		assert.Equal(t, acc.Number(), 1)
		acc2 := tSandbox.MakeNewAccount(addr2)
		assert.Equal(t, acc2.Number(), 2)
		assert.Equal(t, acc2.Balance(), int64(0))
	})
}

func TestTotalValidatorCounter(t *testing.T) {
	setup(t)

	t.Run("Should update total validator counter", func(t *testing.T) {
		assert.Equal(t, tStore.TotalValidators(), 6)

		pub, _ := bls.GenerateTestKeyPair()
		pub2, _ := bls.GenerateTestKeyPair()
		val1 := tSandbox.MakeNewValidator(pub)
		val1.UpdateLastBondingHeight(tSandbox.CurrentHeight())
		assert.Equal(t, val1.Number(), 6)
		assert.Equal(t, val1.LastBondingHeight(), tSandbox.CurrentHeight())

		val2 := tSandbox.MakeNewValidator(pub2)
		val2.UpdateLastBondingHeight(tSandbox.CurrentHeight() + 1)
		assert.Equal(t, val2.Number(), 7)
		assert.Equal(t, val2.LastBondingHeight(), tSandbox.CurrentHeight()+1)
		assert.Equal(t, val2.Stake(), int64(0))
	})
}

func TestCreateDuplicated(t *testing.T) {
	setup(t)

	t.Run("Try creating duplicated account, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		addr := crypto.TreasuryAddress
		tSandbox.MakeNewAccount(addr)
	})

	t.Run("Try creating duplicated validator, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		pub := tSandbox.Validator(tValAddress[3]).PublicKey()
		tSandbox.MakeNewValidator(pub)
	})
}

func TestUpdateFromOutsideTheSandbox(t *testing.T) {
	setup(t)

	t.Run("Try update an account from outside the sandbox, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		acc, _ := account.GenerateTestAccount(999)
		tSandbox.UpdateAccount(acc)
	})

	t.Run("Try update a validator from outside the sandbox, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		val, _ := validator.GenerateTestValidator(999)
		tSandbox.UpdateValidator(val)
	})
}

func TestDeepCopy(t *testing.T) {
	setup(t)

	addr := crypto.GenerateTestAddress()
	pub, _ := bls.GenerateTestKeyPair()
	acc1 := tSandbox.MakeNewAccount(addr)
	val1 := tSandbox.MakeNewValidator(pub)

	acc2 := tSandbox.Account(addr)
	val2 := tSandbox.Validator(pub.Address())

	assert.Equal(t, acc1, acc2)
	assert.Equal(t, val1.Hash(), val2.Hash())

	acc1.IncSequence()
	val1.IncSequence()

	acc2.AddToBalance(1)
	val2.AddToStake(1)

	assert.NotEqual(t, acc1.Hash(), acc2.Hash())
	assert.NotEqual(t, val1.Hash(), val2.Hash())

	acc3 := tSandbox.accounts[addr]
	val3 := tSandbox.validators[pub.Address()]

	assert.NotEqual(t, acc1.Hash(), acc3.Account.Hash())
	assert.NotEqual(t, val1.Hash(), val3.Validator.Hash())

	assert.NotEqual(t, acc2.Hash(), acc3.Account.Hash())
	assert.NotEqual(t, val2.Hash(), val3.Validator.Hash())
}

func TestFindBlockInfoByStamp(t *testing.T) {
	setup(t)

	height, _ := tSandbox.FindBlockInfoByStamp(hash.GenerateTestStamp())
	assert.Equal(t, height, -1)

	latestBlockHeight := tStore.LastBlockHeight()
	latestBlock := tStore.Blocks[latestBlockHeight]
	height, hash := tSandbox.FindBlockInfoByStamp(latestBlock.Stamp())
	assert.Equal(t, height, latestBlockHeight)
	assert.Equal(t, hash, latestBlock.Hash())

	anotherBlockHeight := tStore.LastBlockHeight() - 14
	anotherBlock := tStore.Blocks[anotherBlockHeight]
	height, hash = tSandbox.FindBlockInfoByStamp(anotherBlock.Stamp())
	assert.Equal(t, height, anotherBlockHeight)
	assert.Equal(t, hash, anotherBlock.Hash())
}
