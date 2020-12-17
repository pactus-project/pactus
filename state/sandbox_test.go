package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

func TestAccountChange(t *testing.T) {
	setup(t)

	st, _ := mockState(t, nil)
	sb, err := newSandbox(st.store, st.params, 0, st.sortition, st.validatorSet)
	assert.NoError(t, err)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr, _, _ := crypto.GenerateTestKeyPair()
		assert.Nil(t, sb.Account(invAddr))
	})

	t.Run("Retreive an account from store, modify it and commit it", func(t *testing.T) {
		acc1, _ := account.GenerateTestAccount(0)
		st.store.UpdateAccount(acc1)

		acc1a := sb.Account(acc1.Address())
		assert.Equal(t, acc1, acc1a)

		acc1a.IncSequence()
		acc1a.AddToBalance(acc1a.Balance() + 1)

		sb.UpdateAccount(acc1a)

		sb.CommitChanges(0)

		acc1b, err := sb.store.Account(acc1.Address())
		assert.NoError(t, err)
		assert.Equal(t, acc1a, acc1b)
	})

	t.Run("Make new account and reset the sandbox", func(t *testing.T) {
		addr, _, _ := crypto.GenerateTestKeyPair()
		acc2 := sb.MakeNewAccount(addr)

		acc2.IncSequence()
		acc2.AddToBalance(acc2.Balance() + 1)

		sb.UpdateAccount(acc2)
		acc22 := sb.Account(acc2.Address())
		assert.Equal(t, acc2, acc22)

		sb.Clear()
		assert.Equal(t, len(sb.accounts), 0)
		assert.Nil(t, sb.Account(addr))
	})
}

func TestValidatorChange(t *testing.T) {
	setup(t)

	st, _ := mockState(t, nil)
	sb, err := newSandbox(st.store, st.params, 0, st.sortition, st.validatorSet)
	assert.NoError(t, err)

	t.Run("Should returns nil for invalid address", func(t *testing.T) {
		invAddr, _, _ := crypto.GenerateTestKeyPair()
		assert.Nil(t, sb.Validator(invAddr))
	})

	t.Run("Retreive an validator from store, modify it and commit it", func(t *testing.T) {
		val1, _ := validator.GenerateTestValidator(0)
		st.store.UpdateValidator(val1)

		val1a := sb.Validator(val1.Address())
		assert.Equal(t, val1.Hash(), val1a.Hash())

		val1a.IncSequence()
		val1a.AddToStake(val1a.Stake() + 1)

		sb.UpdateValidator(val1a)

		sb.CommitChanges(0)

		val1b, err := sb.store.Validator(val1.Address())
		assert.NoError(t, err)
		assert.Equal(t, val1a, val1b)
	})

	t.Run("Make new validator and reset the sandbox", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		val2 := sb.MakeNewValidator(pub)

		val2.IncSequence()
		val2.AddToStake(val2.Stake() + 1)

		sb.UpdateValidator(val2)
		val22 := sb.Validator(val2.Address())
		assert.Equal(t, val2, val22)

		sb.Clear()
		assert.Equal(t, len(sb.validators), 0)
		assert.Nil(t, sb.Validator(pub.Address()))
	})
}

func TestAddValidatorToSet(t *testing.T) {
	// setup(t)

	// st, _ := mockState(t, nil)
	// _, pub1, _ := crypto.GenerateTestKeyPair()

	// sb, _ := newSandbox(st.store, st.params, 0, st.sortition, st.validatorSet)

	// val1 := sb.MakeNewValidator(pub1)

	// sb.AddToSet(val1)
	// assert.Nil(t, st.validatorSet.Validator(val1.Address()), "Shouldn't be is not in set")
	// sb.CommitChanges(0)
	// assert.Nil(t, st.validatorSet.Validator(val1.Address()), "Should be is not in set")
}

func TestTotalAccountCounter(t *testing.T) {
	setup(t)

	st, _ := mockState(t, nil)
	sb, _ := newSandbox(st.store, st.params, 0, st.sortition, st.validatorSet)

	t.Run("Should update total account counter", func(t *testing.T) {
		assert.Equal(t, st.store.TotalAccounts(), 1) // Sandbox has an account

		addr, _, _ := crypto.GenerateTestKeyPair()
		acc := sb.MakeNewAccount(addr)
		assert.Equal(t, acc.Number(), 1)

		sb.Clear()
		assert.Equal(t, sb.totalAccounts, 1)
		assert.Equal(t, st.store.TotalAccounts(), 1)

		acc = sb.MakeNewAccount(addr)
		sb.UpdateAccount(acc)

		sb.CommitChanges(0)
		assert.Equal(t, sb.totalAccounts, 2)
		assert.Equal(t, st.store.TotalAccounts(), 2)
	})
}

func TestTotalValidatorCounter(t *testing.T) {
	setup(t)

	st, _ := mockState(t, nil)
	sb, _ := newSandbox(st.store, st.params, 0, st.sortition, st.validatorSet)

	t.Run("Should update total validator counter", func(t *testing.T) {
		assert.Equal(t, st.store.TotalValidators(), 1) // Sandbox has a validator

		_, pub, _ := crypto.GenerateTestKeyPair()
		val := sb.MakeNewValidator(pub)
		assert.Equal(t, val.Number(), 1)

		sb.Clear()
		assert.Equal(t, sb.totalValidators, 1)
		assert.Equal(t, st.store.TotalValidators(), 1)

		val = sb.MakeNewValidator(pub)
		sb.UpdateValidator(val)

		sb.CommitChanges(0)
		assert.Equal(t, sb.totalValidators, 2)
		assert.Equal(t, st.store.TotalValidators(), 2)
	})
}

func TestCreateDuplicated(t *testing.T) {
	setup(t)

	st, _ := mockState(t, nil)
	sb, _ := newSandbox(st.store, st.params, 0, st.sortition, st.validatorSet)

	t.Run("Try creating duplicated account, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		addr := gen.Accounts()[0].Address()
		sb.MakeNewAccount(addr)
	})

	t.Run("Try creating duplicated validator, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		pub := gen.Validators()[0].PublicKey()
		sb.MakeNewValidator(pub)
	})
}

func TestUpdateFromOutsideTheSandbox(t *testing.T) {
	setup(t)

	st, _ := mockState(t, nil)
	sb, _ := newSandbox(st.store, st.params, 0, st.sortition, st.validatorSet)

	t.Run("Try update an account from outside the sandbox, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		acc, _ := account.GenerateTestAccount(1)
		sb.UpdateAccount(acc)
	})

	t.Run("Try update a validator from outside the sandbox, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()
		val, _ := validator.GenerateTestValidator(1)
		sb.UpdateValidator(val)
	})
}

func TestDeepCopy(t *testing.T) {
	setup(t)

	st, _ := mockState(t, nil)
	sb, err := newSandbox(st.store, st.params, 0, st.sortition, st.validatorSet)
	assert.NoError(t, err)

	addr, pub, _ := crypto.GenerateTestKeyPair()
	acc1 := sb.MakeNewAccount(addr)
	val1 := sb.MakeNewValidator(pub)

	acc2 := sb.Account(addr)
	val2 := sb.Validator(pub.Address())

	assert.Equal(t, acc1, acc2)
	assert.Equal(t, val1.Hash(), val2.Hash())

	acc1.IncSequence()
	val1.IncSequence()

	acc2.AddToBalance(1)
	val2.AddToStake(1)

	assert.NotEqual(t, acc1.Hash(), acc2.Hash())
	assert.NotEqual(t, val1.Hash(), val2.Hash())

	acc3 := sb.accounts[addr]
	val3 := sb.validators[pub.Address()]

	assert.NotEqual(t, acc1.Hash(), acc3.account.Hash())
	assert.NotEqual(t, val1.Hash(), val3.validator.Hash())

	assert.NotEqual(t, acc2.Hash(), acc3.account.Hash())
	assert.NotEqual(t, val2.Hash(), val3.validator.Hash())
}

func TestChangeToStake(t *testing.T) {
	setup(t)

	st, _ := mockState(t, nil)
	sb, err := newSandbox(st.store, st.params, 0, st.sortition, st.validatorSet)
	assert.NoError(t, err)

	_, pub1, _ := crypto.GenerateTestKeyPair()
	_, pub2, _ := crypto.GenerateTestKeyPair()
	val1 := sb.MakeNewValidator(pub1)
	val2 := sb.MakeNewValidator(pub2)

	val1.AddToStake(1000)
	val2.AddToStake(2000)
	sb.UpdateValidator(val1)

	assert.Equal(t, sb.changeToStake, int64(1000))
	val1.AddToStake(500)
	assert.Equal(t, sb.changeToStake, int64(1000))

	sb.UpdateValidator(val1)
	sb.UpdateValidator(val2)
	assert.Equal(t, sb.changeToStake, int64(3500))
	val2.WithdrawStake()
	assert.Equal(t, sb.changeToStake, int64(3500))

	sb.UpdateValidator(val2)
	assert.Equal(t, sb.changeToStake, int64(1500))
}
