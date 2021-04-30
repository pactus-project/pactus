package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var tValSigners [9]crypto.Signer
var tSortitions *sortition.Sortition
var tCommittee *committee.Committee
var tStore *store.MockStore
var tSandbox *Concrete

func init() {
	logger.InitLogger(logger.TestConfig())
}

func setup(t *testing.T) {
	var err error
	tStore = store.MockingStore()

	lastHeight := 21
	for i := 0; i <= lastHeight; i++ {
		b, _ := block.GenerateTestBlock(nil, nil)
		tStore.SaveBlock(i+1, b)
	}

	_, pub1, priv1 := crypto.GenerateTestKeyPair()
	_, pub2, priv2 := crypto.GenerateTestKeyPair()
	_, pub3, priv3 := crypto.GenerateTestKeyPair()
	_, pub4, priv4 := crypto.GenerateTestKeyPair()
	_, pub5, priv5 := crypto.GenerateTestKeyPair()
	_, pub6, priv6 := crypto.GenerateTestKeyPair()
	_, pub7, priv7 := crypto.GenerateTestKeyPair()
	_, pub8, priv8 := crypto.GenerateTestKeyPair()

	tValSigners[0] = crypto.NewSigner(priv1)
	tValSigners[1] = crypto.NewSigner(priv2)
	tValSigners[2] = crypto.NewSigner(priv3)
	tValSigners[3] = crypto.NewSigner(priv4)
	tValSigners[4] = crypto.NewSigner(priv5)
	tValSigners[5] = crypto.NewSigner(priv6)
	tValSigners[6] = crypto.NewSigner(priv7)
	tValSigners[7] = crypto.NewSigner(priv8)

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14)

	val1 := validator.NewValidator(pub1, 0, util.RandInt(1000))
	val2 := validator.NewValidator(pub2, 1, util.RandInt(1000))
	val3 := validator.NewValidator(pub3, 2, util.RandInt(1000))
	val4 := validator.NewValidator(pub4, 3, util.RandInt(1000))
	val5 := validator.NewValidator(pub5, 4, util.RandInt(1000))
	val6 := validator.NewValidator(pub6, 5, util.RandInt(1000))
	val7 := validator.NewValidator(pub7, 6, util.RandInt(1000))
	val8 := validator.NewValidator(pub8, 7, util.RandInt(1000))

	val1.AddToStake(1000)
	val2.AddToStake(2000)
	val3.AddToStake(3000)
	val4.AddToStake(4000)
	val5.AddToStake(1000)
	val6.AddToStake(3000)
	val7.AddToStake(1000)
	val8.AddToStake(100)

	tStore.UpdateAccount(acc)
	tStore.UpdateValidator(val1)
	tStore.UpdateValidator(val2)
	tStore.UpdateValidator(val3)
	tStore.UpdateValidator(val4)
	tStore.UpdateValidator(val5)
	tStore.UpdateValidator(val6)
	tStore.UpdateValidator(val7)
	tStore.UpdateValidator(val8)

	tSortitions = sortition.NewSortition()
	tCommittee, err = committee.NewCommittee([]*validator.Validator{val1, val2, val3, val4}, 4, tValSigners[0].Address())
	assert.NoError(t, err)

	params := param.DefaultParams()
	tSandbox = NewSandbox(tStore, params, tStore.LastBlockHeight(), tSortitions, tCommittee)
	assert.Equal(t, tSandbox.MaxMemoLength(), params.MaximumMemoLength)
	assert.Equal(t, tSandbox.FeeFraction(), params.FeeFraction)
	assert.Equal(t, tSandbox.MinFee(), params.MinimumFee)
	assert.Equal(t, tSandbox.TransactionToLiveInterval(), params.TransactionToLiveInterval)
	assert.Equal(t, tSandbox.CommitteeSize(), params.CommitteeSize)
}

func TestAccountChange(t *testing.T) {
	setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr, _, _ := crypto.GenerateTestKeyPair()
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
		addr, _, _ := crypto.GenerateTestKeyPair()
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
		invAddr, _, _ := crypto.GenerateTestKeyPair()
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
		assert.False(t, tSandbox.validators[val1a.Address()].JoinedCommittee)
		tSandbox.UpdateValidator(val1a)
		assert.True(t, tSandbox.validators[val1a.Address()].Updated)
		assert.False(t, tSandbox.validators[val1a.Address()].JoinedCommittee)
	})

	t.Run("Make new validator", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		val2 := tSandbox.MakeNewValidator(pub)

		val2.IncSequence()
		val2.AddToStake(val2.Stake() + 1)

		tSandbox.UpdateValidator(val2)
		val22 := tSandbox.Validator(val2.Address())
		assert.Equal(t, val2, val22)
	})
}

func TestAddValidatorToCommittee(t *testing.T) {
	setup(t)

	stamp := tSandbox.LastBlockHash()
	height := tSandbox.CurrentHeight()

	t.Run("Add unknown validator to the committee, Should returns error", func(t *testing.T) {
		val, _ := validator.GenerateTestValidator(777)
		assert.Error(t, tSandbox.EnterCommittee(stamp, val.Address()))
	})

	t.Run("Already in the committee, Should returns error", func(t *testing.T) {
		v := tSandbox.Validator(tValSigners[3].Address())
		assert.Error(t, tSandbox.EnterCommittee(stamp, v.Address()))
	})

	t.Run("Invalid block hash, Should returns error", func(t *testing.T) {
		_, pub1, _ := crypto.GenerateTestKeyPair()
		val := tSandbox.MakeNewValidator(pub1)
		assert.Error(t, tSandbox.EnterCommittee(crypto.GenerateTestHash(), val.Address()))
	})

	t.Run("More than 1/3 stake, Should returns error", func(t *testing.T) {
		tSandbox.params.CommitteeSize = 4

		val5 := tSandbox.Validator(tValSigners[4].Address())
		val6 := tSandbox.Validator(tValSigners[5].Address())
		val7 := tSandbox.Validator(tValSigners[6].Address())

		assert.Equal(t, tSandbox.committee.TotalStake(), int64(10000), "Total stake should be 10000")
		assert.NoError(t, tSandbox.EnterCommittee(stamp, val5.Address()))
		assert.Error(t, tSandbox.EnterCommittee(stamp, val5.Address()), "Duplicated entry")
		assert.Error(t, tSandbox.EnterCommittee(stamp, val6.Address()), "More than 1/3 of stake is going to change")
		assert.NoError(t, tSandbox.EnterCommittee(stamp, val7.Address()))
	})

	t.Run("In committee at time of sending sortition, Should returns error", func(t *testing.T) {
		tSandbox.params.CommitteeSize = 4

		b, _ := tStore.Block(height - 3)
		num := b.LastCertificate().Committers()[2]
		addr, pub, _ := crypto.GenerateTestKeyPair()
		val := validator.NewValidator(pub, num, 0)
		tStore.UpdateValidator(val)
		assert.Equal(t, tSandbox.Validator(addr), val)

		assert.Error(t, tSandbox.EnterCommittee(stamp, val.Address()))
	})

	t.Run("Not leaving committee before proposing a block", func(t *testing.T) {

		vals := tSandbox.committee.Validators()
		vals[0].UpdateLastJoinedHeight(height - 4)
		vals[1].UpdateLastJoinedHeight(height - 3)
		vals[2].UpdateLastJoinedHeight(height - 2)
		vals[3].UpdateLastJoinedHeight(height - 1)

		val8 := tSandbox.Validator(tValSigners[7].Address())

		assert.Error(t, tSandbox.EnterCommittee(stamp, val8.Address()))
	})

	t.Run("Update validator and add to committee", func(t *testing.T) {
		tSandbox.params.CommitteeSize = 8

		addr1, pub1, _ := crypto.GenerateTestKeyPair()
		val1 := tSandbox.MakeNewValidator(pub1)
		assert.NoError(t, tSandbox.EnterCommittee(stamp, val1.Address()))
		seq := val1.Sequence()
		val1.IncSequence()
		tSandbox.UpdateValidator(val1)
		val := tSandbox.validators[addr1]
		assert.True(t, val.JoinedCommittee)
		assert.True(t, val.Updated)
		assert.Equal(t, val.Validator.Sequence(), seq+1)
	})
}

func TestTotalAccountCounter(t *testing.T) {
	setup(t)

	t.Run("Should update total account counter", func(t *testing.T) {
		assert.Equal(t, tStore.TotalAccounts(), 1) // Sandbox has an account

		addr, _, _ := crypto.GenerateTestKeyPair()
		addr2, _, _ := crypto.GenerateTestKeyPair()
		acc := tSandbox.MakeNewAccount(addr)
		assert.Equal(t, acc.Number(), 1)
		acc2 := tSandbox.MakeNewAccount(addr2)
		assert.Equal(t, acc2.Number(), 2)
		assert.Equal(t, acc2.Balance(), int64(0))
	})
}

func TestTotalValidatorCounter(t *testing.T) {
	setup(t)

	t.Run("Should update total validator counter", func(t *testing.T) {
		assert.Equal(t, tStore.TotalValidators(), 8)

		_, pub, _ := crypto.GenerateTestKeyPair()
		_, pub2, _ := crypto.GenerateTestKeyPair()
		val := tSandbox.MakeNewValidator(pub)
		assert.Equal(t, val.Number(), 8)
		assert.Equal(t, val.BondingHeight(), tSandbox.CurrentHeight())
		val2 := tSandbox.MakeNewValidator(pub2)
		assert.Equal(t, val2.Number(), 9)
		assert.Equal(t, val2.BondingHeight(), tSandbox.CurrentHeight())
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
		pub := tValSigners[3].PublicKey()
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

	addr, pub, _ := crypto.GenerateTestKeyPair()
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
