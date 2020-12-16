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
	sb, err := newSandbox(st.store, st.params, 0)
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

		assert.NoError(t, sb.CommitAndClear(st.validatorSet))

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
	sb, err := newSandbox(st.store, st.params, 0)
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

		assert.NoError(t, sb.CommitAndClear(st.validatorSet))

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
	setup(t)

	st, _ := mockState(t, nil)
	_, pub1, _ := crypto.GenerateTestKeyPair()

	sb, _ := newSandbox(st.store, st.params, 0)

	val1 := sb.MakeNewValidator(pub1)

	sb.AddToSet(val1)
	assert.Nil(t, st.validatorSet.Validator(val1.Address()), "Shouldn't be is not in set")
	assert.NoError(t, sb.CommitAndClear(st.validatorSet))
	assert.Nil(t, st.validatorSet.Validator(val1.Address()), "Should be is not in set")
}

func TestTotalAccountCounter(t *testing.T) {
	st, _ := mockState(t, nil)
	sb, _ := newSandbox(st.store, st.params, 0)

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

		assert.NoError(t, sb.CommitAndClear(st.validatorSet))
		assert.Equal(t, sb.totalAccounts, 2)
		assert.Equal(t, st.store.TotalAccounts(), 2)
	})
}

func TestTotalValidatorCounter(t *testing.T) {
	st, _ := mockState(t, nil)
	sb, _ := newSandbox(st.store, st.params, 0)

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

		assert.NoError(t, sb.CommitAndClear(st.validatorSet))
		assert.Equal(t, sb.totalValidators, 2)
		assert.Equal(t, st.store.TotalValidators(), 2)
	})
}

func TestCreateduplicated(t *testing.T) {
	st, _ := mockState(t, nil)
	sb, _ := newSandbox(st.store, st.params, 0)

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

func TestUpdateFromOutsideThesandbox(t *testing.T) {
	st, _ := mockState(t, nil)
	sb, _ := newSandbox(st.store, st.params, 0)

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
