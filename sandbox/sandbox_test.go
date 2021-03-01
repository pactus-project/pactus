package sandbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/validator"
)

var tValSigners [5]crypto.Signer
var tSortitions [5]*sortition.Sortition
var tValset *validator.ValidatorSet
var tStore *store.MockStore
var tSandbox *SandboxConcrete

func init() {
	logger.InitLogger(logger.TestConfig())
}

func setup(t *testing.T) {
	var err error
	tStore = store.MockingStore()

	_, pub1, priv1 := crypto.GenerateTestKeyPair()
	_, pub2, priv2 := crypto.GenerateTestKeyPair()
	_, pub3, priv3 := crypto.GenerateTestKeyPair()
	_, pub4, priv4 := crypto.GenerateTestKeyPair()
	_, _, priv5 := crypto.GenerateTestKeyPair()
	tValSigners[0] = crypto.NewSigner(priv1)
	tValSigners[1] = crypto.NewSigner(priv2)
	tValSigners[2] = crypto.NewSigner(priv3)
	tValSigners[3] = crypto.NewSigner(priv4)
	tValSigners[4] = crypto.NewSigner(priv5)

	acc := account.NewAccount(crypto.TreasuryAddress, 0)
	acc.AddToBalance(21 * 1e14)
	val1 := validator.NewValidator(pub1, 0, 0)
	val2 := validator.NewValidator(pub2, 1, 0)
	val3 := validator.NewValidator(pub2, 2, 0)
	val4 := validator.NewValidator(pub3, 3, 0)
	val5 := validator.NewValidator(pub4, 4, 0)

	tStore.UpdateAccount(acc)
	tStore.UpdateValidator(val1)
	tStore.UpdateValidator(val2)
	tStore.UpdateValidator(val3)
	tStore.UpdateValidator(val4)
	tStore.UpdateValidator(val5)

	tSortitions[0] = sortition.NewSortition()
	tSortitions[1] = sortition.NewSortition()
	tSortitions[2] = sortition.NewSortition()
	tSortitions[3] = sortition.NewSortition()
	tSortitions[4] = sortition.NewSortition()
	tValset, err = validator.NewValidatorSet([]*validator.Validator{val1, val2, val3, val4}, 4, tValSigners[0].Address())
	assert.NoError(t, err)

	params := param.DefaultParams()
	tSandbox, err = NewSandbox(tStore, params, 0, tSortitions[0], tValset)
	assert.NoError(t, err)
	assert.Equal(t, tSandbox.MaxMemoLength(), params.MaximumMemoLength)
	assert.Equal(t, tSandbox.FeeFraction(), params.FeeFraction)
	assert.Equal(t, tSandbox.MinFee(), params.MinimumFee)
	assert.Equal(t, tSandbox.TransactionToLiveInterval(), params.TransactionToLiveInterval)
	assert.Equal(t, tSandbox.CommitteeSize(), params.CommitteeSize)
}

func TestLoadRecentBlocks(t *testing.T) {
	store := store.MockingStore()

	lastHeight := 21
	for i := 0; i <= lastHeight; i++ {
		b, _ := block.GenerateTestBlock(nil, nil)
		store.Blocks[i+1] = b
	}

	params := param.DefaultParams()
	params.TransactionToLiveInterval = 10

	sandbox, err := NewSandbox(store, params, lastHeight, nil, nil)
	assert.NoError(t, err)

	assert.Equal(t, sandbox.RecentBlockHeight(crypto.UndefHash), -1)
	assert.Equal(t, sandbox.RecentBlockHeight(store.Blocks[11].Hash()), -1)
	assert.Equal(t, sandbox.RecentBlockHeight(store.Blocks[12].Hash()), 12)
	assert.Equal(t, sandbox.RecentBlockHeight(store.Blocks[21].Hash()), 21)

	h := crypto.GenerateTestHash()
	sandbox.AppendNewBlock(h, 22)
	assert.Equal(t, sandbox.RecentBlockHeight(store.Blocks[12].Hash()), -1)
	assert.Equal(t, sandbox.RecentBlockHeight(store.Blocks[21].Hash()), 21)
}

func TestAccountChange(t *testing.T) {
	setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr, _, _ := crypto.GenerateTestKeyPair()
		assert.Nil(t, tSandbox.Account(invAddr))
	})

	t.Run("Retrieve an account from store, modify it and commit it", func(t *testing.T) {
		acc1, _ := account.GenerateTestAccount(0)
		tStore.UpdateAccount(acc1)

		acc1a := tSandbox.Account(acc1.Address())
		assert.Equal(t, acc1, acc1a)

		acc1a.IncSequence()
		acc1a.AddToBalance(acc1a.Balance() + 1)

		assert.False(t, tSandbox.accounts[acc1a.Address()].Updated)
		tSandbox.UpdateAccount(acc1a)
		assert.True(t, tSandbox.accounts[acc1a.Address()].Updated)
	})

	t.Run("Make new account and reset the sandbox", func(t *testing.T) {
		addr, _, _ := crypto.GenerateTestKeyPair()
		acc2 := tSandbox.MakeNewAccount(addr)

		acc2.IncSequence()
		acc2.AddToBalance(acc2.Balance() + 1)

		tSandbox.UpdateAccount(acc2)
		acc22 := tSandbox.Account(acc2.Address())
		assert.Equal(t, acc2, acc22)

		tSandbox.Clear()
		assert.Equal(t, len(tSandbox.accounts), 0)
		assert.Nil(t, tSandbox.Account(addr))
	})
}

func TestValidatorChange(t *testing.T) {
	setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr, _, _ := crypto.GenerateTestKeyPair()
		assert.Nil(t, tSandbox.Validator(invAddr))
	})

	t.Run("Retrieve an validator from store, modify it and commit it", func(t *testing.T) {
		val1, _ := validator.GenerateTestValidator(0)
		tStore.UpdateValidator(val1)

		val1a := tSandbox.Validator(val1.Address())
		assert.Equal(t, val1.Hash(), val1a.Hash())

		val1a.IncSequence()
		val1a.AddToStake(val1a.Stake() + 1)

		assert.False(t, tSandbox.validators[val1a.Address()].Updated)
		assert.False(t, tSandbox.validators[val1a.Address()].AddToSet)
		tSandbox.UpdateValidator(val1a)
		assert.True(t, tSandbox.validators[val1a.Address()].Updated)
		assert.False(t, tSandbox.validators[val1a.Address()].AddToSet)
	})

	t.Run("Make new validator and reset the sandbox", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		val2 := tSandbox.MakeNewValidator(pub)

		val2.IncSequence()
		val2.AddToStake(val2.Stake() + 1)

		tSandbox.UpdateValidator(val2)
		val22 := tSandbox.Validator(val2.Address())
		assert.Equal(t, val2, val22)

		tSandbox.Clear()
		assert.Equal(t, len(tSandbox.validators), 0)
		assert.Nil(t, tSandbox.Validator(pub.Address()))
	})
}

func TestAddValidatorToSet(t *testing.T) {
	setup(t)

	val5, err := tStore.ValidatorByNumber(4)
	require.NoError(t, err)
	a1 := tValSigners[0].Address()
	a2 := tValSigners[1].Address()
	block11, _ := block.GenerateTestBlock(&a1, nil)
	assert.NoError(t, tValset.UpdateTheSet(0, nil))
	block12, _ := block.GenerateTestBlock(&a2, nil)
	assert.NoError(t, tValset.UpdateTheSet(0, []*validator.Validator{val5}))

	tStore.Blocks[11] = block11
	tStore.Blocks[12] = block12

	t.Run("Add unknown validator to the set, Should returns error", func(t *testing.T) {
		val, _ := validator.GenerateTestValidator(1)
		h := crypto.GenerateTestHash()
		assert.Error(t, tSandbox.AddToSet(h, val.Address()))
	})

	t.Run("Already in the set, Should returns error", func(t *testing.T) {
		h := crypto.GenerateTestHash()
		v := tSandbox.Validator(tValSigners[3].Address())
		assert.Error(t, tSandbox.AddToSet(h, v.Address()))
	})

	t.Run("In set at time of doing sortition, Should returns error", func(t *testing.T) {
		v := tSandbox.Validator(tValSigners[0].Address())
		assert.Error(t, tSandbox.AddToSet(block11.Hash(), v.Address()))
	})

	t.Run("Invalid block hash, Should returns error", func(t *testing.T) {
		v := tSandbox.Validator(tValSigners[0].Address())
		assert.Error(t, tSandbox.AddToSet(crypto.GenerateTestHash(), v.Address()))
	})

	t.Run("More than 1/3, Should returns error", func(t *testing.T) {
		tSandbox.params.CommitteeSize = 4

		_, pub1, _ := crypto.GenerateTestKeyPair()
		_, pub2, _ := crypto.GenerateTestKeyPair()
		val1 := tSandbox.MakeNewValidator(pub1)
		val2 := tSandbox.MakeNewValidator(pub2)
		assert.NoError(t, tSandbox.AddToSet(block11.Hash(), val1.Address()))
		assert.Error(t, tSandbox.AddToSet(block11.Hash(), val2.Address()))
	})

	t.Run("Update validator and add to set", func(t *testing.T) {
		tSandbox.Clear()
		addr1, pub1, _ := crypto.GenerateTestKeyPair()
		val1 := tSandbox.MakeNewValidator(pub1)
		assert.NoError(t, tSandbox.AddToSet(block11.Hash(), val1.Address()))
		seq := val1.Sequence()
		val1.IncSequence()
		tSandbox.UpdateValidator(val1)
		vs := tSandbox.validators[addr1]
		assert.True(t, vs.AddToSet)
		assert.True(t, vs.Updated)
		assert.Equal(t, vs.Validator.Sequence(), seq+1)
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

		tSandbox.Clear()
		assert.Equal(t, tSandbox.totalAccounts, 1)
		assert.Equal(t, tStore.TotalAccounts(), 1)

		acc = tSandbox.MakeNewAccount(addr)
		assert.Equal(t, tSandbox.totalAccounts, 2)

		tSandbox.UpdateAccount(acc)

		assert.Equal(t, tSandbox.totalAccounts, 2)
		assert.Equal(t, tStore.TotalAccounts(), 1)
	})
}

func TestTotalValidatorCounter(t *testing.T) {
	setup(t)

	t.Run("Should update total validator counter", func(t *testing.T) {
		assert.Equal(t, tStore.TotalValidators(), 4) // Sandbox has 5 validators

		_, pub, _ := crypto.GenerateTestKeyPair()
		_, pub2, _ := crypto.GenerateTestKeyPair()
		val := tSandbox.MakeNewValidator(pub)
		assert.Equal(t, val.Number(), 4)
		assert.Equal(t, val.BondingHeight(), tSandbox.CurrentHeight())
		val2 := tSandbox.MakeNewValidator(pub2)
		assert.Equal(t, val2.Number(), 5)
		assert.Equal(t, val2.BondingHeight(), tSandbox.CurrentHeight())
		assert.Equal(t, val2.Stake(), int64(0))

		tSandbox.Clear()
		assert.Equal(t, tSandbox.totalValidators, 4)
		assert.Equal(t, tStore.TotalValidators(), 4)

		val = tSandbox.MakeNewValidator(pub)
		tSandbox.UpdateValidator(val)
		assert.Equal(t, val.Number(), 4)
		assert.Equal(t, val.BondingHeight(), tSandbox.CurrentHeight())

		assert.Equal(t, tSandbox.totalValidators, 5)
		assert.Equal(t, tStore.TotalValidators(), 4)
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
		acc, _ := account.GenerateTestAccount(1)
		tSandbox.UpdateAccount(acc)
	})

	t.Run("Try update a validator from outside the sandbox, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		val, _ := validator.GenerateTestValidator(1)
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

func TestChangeToStake(t *testing.T) {
	setup(t)

	_, pub1, _ := crypto.GenerateTestKeyPair()
	_, pub2, _ := crypto.GenerateTestKeyPair()
	val1 := tSandbox.MakeNewValidator(pub1)
	val2 := tSandbox.MakeNewValidator(pub2)
	val3 := tSandbox.Validator(tValSigners[0].Address())

	val1.AddToStake(1000)
	val2.AddToStake(2000)
	val3.AddToStake(3000)

	tSandbox.UpdateValidator(val1)
	assert.Equal(t, tSandbox.TotalStakeChange(), int64(1000))

	val1.AddToStake(500)
	assert.Equal(t, tSandbox.TotalStakeChange(), int64(1000))

	tSandbox.UpdateValidator(val1)
	tSandbox.UpdateValidator(val2)
	tSandbox.UpdateValidator(val3)
	assert.Equal(t, tSandbox.TotalStakeChange(), int64(6500))

	tSandbox.Clear()
	assert.Equal(t, tSandbox.TotalStakeChange(), int64(0))
}
