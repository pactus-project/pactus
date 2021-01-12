package sandbox

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
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
var tSandbox1 *SandboxConcrete

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
	acc.AddToBalance(21000000000000)
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

	tSortitions[0] = sortition.NewSortition(tValSigners[0])
	tSortitions[1] = sortition.NewSortition(tValSigners[1])
	tSortitions[2] = sortition.NewSortition(tValSigners[2])
	tSortitions[3] = sortition.NewSortition(tValSigners[3])
	tSortitions[4] = sortition.NewSortition(tValSigners[4])

	tValset, err = validator.NewValidatorSet([]*validator.Validator{val1, val2, val3, val4}, 4, tValSigners[0].Address())
	assert.NoError(t, err)

	tParams := param.MainnetParams()
	tSandbox1, err = NewSandbox(tStore, tParams, 0, tSortitions[0], tValset)
	assert.NoError(t, err)
}

func TestLoadRecentBlocks(t *testing.T) {
	store := store.MockingStore()

	lastHeight := 21
	for i := 0; i <= lastHeight; i++ {
		b, _ := block.GenerateTestBlock(nil, nil)
		store.Blocks[i+1] = b
	}

	params := param.MainnetParams()
	params.TransactionToLiveInterval = 10

	sandbox, err := NewSandbox(store, params, lastHeight, nil, nil)
	assert.NoError(t, err)

	_, ok := sandbox.recentBlocks.Get(crypto.UndefHash)
	assert.False(t, ok)

	v, _ := sandbox.recentBlocks.Get(store.Blocks[21].Hash())
	assert.Equal(t, v, 21)

	_, ok = sandbox.recentBlocks.Get(store.Blocks[11].Hash())
	assert.False(t, ok)

	v, _ = sandbox.recentBlocks.Get(store.Blocks[12].Hash())
	assert.Equal(t, v, 12)
}

func TestAccountChange(t *testing.T) {
	setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr, _, _ := crypto.GenerateTestKeyPair()
		assert.Nil(t, tSandbox1.Account(invAddr))
	})

	t.Run("Retrieve an account from store, modify it and commit it", func(t *testing.T) {
		acc1, _ := account.GenerateTestAccount(0)
		tStore.UpdateAccount(acc1)

		acc1a := tSandbox1.Account(acc1.Address())
		assert.Equal(t, acc1, acc1a)

		acc1a.IncSequence()
		acc1a.AddToBalance(acc1a.Balance() + 1)

		assert.False(t, tSandbox1.accounts[acc1a.Address()].Updated)
		tSandbox1.UpdateAccount(acc1a)
		assert.True(t, tSandbox1.accounts[acc1a.Address()].Updated)
	})

	t.Run("Make new account and reset the sandbox", func(t *testing.T) {
		addr, _, _ := crypto.GenerateTestKeyPair()
		acc2 := tSandbox1.MakeNewAccount(addr)

		acc2.IncSequence()
		acc2.AddToBalance(acc2.Balance() + 1)

		tSandbox1.UpdateAccount(acc2)
		acc22 := tSandbox1.Account(acc2.Address())
		assert.Equal(t, acc2, acc22)

		tSandbox1.Clear()
		assert.Equal(t, len(tSandbox1.accounts), 0)
		assert.Nil(t, tSandbox1.Account(addr))
	})
}

func TestValidatorChange(t *testing.T) {
	setup(t)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr, _, _ := crypto.GenerateTestKeyPair()
		assert.Nil(t, tSandbox1.Validator(invAddr))
	})

	t.Run("Retrieve an validator from store, modify it and commit it", func(t *testing.T) {
		val1, _ := validator.GenerateTestValidator(0)
		tStore.UpdateValidator(val1)

		val1a := tSandbox1.Validator(val1.Address())
		assert.Equal(t, val1.Hash(), val1a.Hash())

		val1a.IncSequence()
		val1a.AddToStake(val1a.Stake() + 1)

		assert.False(t, tSandbox1.validators[val1a.Address()].Updated)
		assert.False(t, tSandbox1.validators[val1a.Address()].AddToSet)
		tSandbox1.UpdateValidator(val1a)
		assert.True(t, tSandbox1.validators[val1a.Address()].Updated)
		assert.False(t, tSandbox1.validators[val1a.Address()].AddToSet)
	})

	t.Run("Make new validator and reset the sandbox", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		val2 := tSandbox1.MakeNewValidator(pub)

		val2.IncSequence()
		val2.AddToStake(val2.Stake() + 1)

		tSandbox1.UpdateValidator(val2)
		val22 := tSandbox1.Validator(val2.Address())
		assert.Equal(t, val2, val22)

		tSandbox1.Clear()
		assert.Equal(t, len(tSandbox1.validators), 0)
		assert.Nil(t, tSandbox1.Validator(pub.Address()))
	})
}

func TestAddValidatorToSet(t *testing.T) {
	setup(t)

	val5, err := tStore.ValidatorByNumber(4)
	require.NoError(t, err)
	a1 := tValSigners[0].Address()
	a2 := tValSigners[1].Address()
	block1, _ := block.GenerateTestBlock(&a1, nil)
	assert.NoError(t, tValset.UpdateTheSet(0, nil))
	block2, _ := block.GenerateTestBlock(&a2, nil)
	assert.NoError(t, tValset.UpdateTheSet(0, []*validator.Validator{val5}))

	tStore.Blocks[1] = block1
	tStore.Blocks[2] = block2

	t.Run("Add unknown validator to the set, Should panic", func(t *testing.T) {
		val, _ := validator.GenerateTestValidator(1)
		h := crypto.GenerateTestHash()
		assert.Error(t, tSandbox1.AddToSet(h, val.Address()))
	})

	t.Run("Already in the set, Should returns error", func(t *testing.T) {
		h := crypto.GenerateTestHash()
		v := tSandbox1.Validator(tValSigners[3].Address())
		assert.Error(t, tSandbox1.AddToSet(h, v.Address()))
	})

	t.Run("In set at time of doing sortition, Should returns error", func(t *testing.T) {
		v := tSandbox1.Validator(tValSigners[0].Address())
		assert.Error(t, tSandbox1.AddToSet(block1.Hash(), v.Address()))
	})

	t.Run("More than 1/3, Should returns error", func(t *testing.T) {
		_, pub1, _ := crypto.GenerateTestKeyPair()
		_, pub2, _ := crypto.GenerateTestKeyPair()
		val1 := tSandbox1.MakeNewValidator(pub1)
		val2 := tSandbox1.MakeNewValidator(pub2)
		assert.NoError(t, tSandbox1.AddToSet(block1.Hash(), val1.Address()))
		assert.Error(t, tSandbox1.AddToSet(block1.Hash(), val2.Address()))
	})

}

func TestTotalAccountCounter(t *testing.T) {
	setup(t)

	t.Run("Should update total account counter", func(t *testing.T) {
		assert.Equal(t, tStore.TotalAccounts(), 1) // Sandbox has an account

		addr, _, _ := crypto.GenerateTestKeyPair()
		addr2, _, _ := crypto.GenerateTestKeyPair()
		acc := tSandbox1.MakeNewAccount(addr)
		assert.Equal(t, acc.Number(), 1)
		acc2 := tSandbox1.MakeNewAccount(addr2)
		assert.Equal(t, acc2.Number(), 2)
		assert.Equal(t, acc2.Balance(), int64(0))

		tSandbox1.Clear()
		assert.Equal(t, tSandbox1.totalAccounts, 1)
		assert.Equal(t, tStore.TotalAccounts(), 1)

		acc = tSandbox1.MakeNewAccount(addr)
		assert.Equal(t, tSandbox1.totalAccounts, 2)

		tSandbox1.UpdateAccount(acc)

		assert.Equal(t, tSandbox1.totalAccounts, 2)
		assert.Equal(t, tStore.TotalAccounts(), 1)
	})
}

func TestTotalValidatorCounter(t *testing.T) {
	setup(t)

	t.Run("Should update total validator counter", func(t *testing.T) {
		assert.Equal(t, tStore.TotalValidators(), 4) // Sandbox has 5 validators

		_, pub, _ := crypto.GenerateTestKeyPair()
		_, pub2, _ := crypto.GenerateTestKeyPair()
		val := tSandbox1.MakeNewValidator(pub)
		assert.Equal(t, val.Number(), 4)
		assert.Equal(t, val.BondingHeight(), tSandbox1.CurrentHeight())
		val2 := tSandbox1.MakeNewValidator(pub2)
		assert.Equal(t, val2.Number(), 5)
		assert.Equal(t, val2.BondingHeight(), tSandbox1.CurrentHeight())
		assert.Equal(t, val2.Stake(), int64(0))

		tSandbox1.Clear()
		assert.Equal(t, tSandbox1.totalValidators, 4)
		assert.Equal(t, tStore.TotalValidators(), 4)

		val = tSandbox1.MakeNewValidator(pub)
		tSandbox1.UpdateValidator(val)
		assert.Equal(t, val.Number(), 4)
		assert.Equal(t, val.BondingHeight(), tSandbox1.CurrentHeight())

		assert.Equal(t, tSandbox1.totalValidators, 5)
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
		tSandbox1.MakeNewAccount(addr)
	})

	t.Run("Try creating duplicated validator, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		pub := tValSigners[3].PublicKey()
		tSandbox1.MakeNewValidator(pub)
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
		tSandbox1.UpdateAccount(acc)
	})

	t.Run("Try update a validator from outside the sandbox, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		val, _ := validator.GenerateTestValidator(1)
		tSandbox1.UpdateValidator(val)
	})
}

func TestDeepCopy(t *testing.T) {
	setup(t)

	addr, pub, _ := crypto.GenerateTestKeyPair()
	acc1 := tSandbox1.MakeNewAccount(addr)
	val1 := tSandbox1.MakeNewValidator(pub)

	acc2 := tSandbox1.Account(addr)
	val2 := tSandbox1.Validator(pub.Address())

	assert.Equal(t, acc1, acc2)
	assert.Equal(t, val1.Hash(), val2.Hash())

	acc1.IncSequence()
	val1.IncSequence()

	acc2.AddToBalance(1)
	val2.AddToStake(1)

	assert.NotEqual(t, acc1.Hash(), acc2.Hash())
	assert.NotEqual(t, val1.Hash(), val2.Hash())

	acc3 := tSandbox1.accounts[addr]
	val3 := tSandbox1.validators[pub.Address()]

	assert.NotEqual(t, acc1.Hash(), acc3.Account.Hash())
	assert.NotEqual(t, val1.Hash(), val3.Validator.Hash())

	assert.NotEqual(t, acc2.Hash(), acc3.Account.Hash())
	assert.NotEqual(t, val2.Hash(), val3.Validator.Hash())
}

func TestChangeToStake(t *testing.T) {
	setup(t)

	_, pub1, _ := crypto.GenerateTestKeyPair()
	_, pub2, _ := crypto.GenerateTestKeyPair()
	val1 := tSandbox1.MakeNewValidator(pub1)
	val2 := tSandbox1.MakeNewValidator(pub2)

	val1.AddToStake(1000)
	val2.AddToStake(2000)
	tSandbox1.UpdateValidator(val1)

	assert.Equal(t, tSandbox1.changeToStake, int64(1000))
	val1.AddToStake(500)
	assert.Equal(t, tSandbox1.changeToStake, int64(1000))

	tSandbox1.UpdateValidator(val1)
	tSandbox1.UpdateValidator(val2)
	assert.Equal(t, tSandbox1.changeToStake, int64(3500))
}
