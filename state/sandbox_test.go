package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
)

func TestAccountChange(t *testing.T) {
	st, _ := mockState(t, nil)
	addr1, _, _ := crypto.GenerateTestKeyPair()
	addr2, _, _ := crypto.GenerateTestKeyPair()
	acc1 := account.NewAccount(addr1, 0)
	acc2 := account.NewAccount(addr2, 0)

	st.store.UpdateAccount(acc1)
	sb, err := newSandbox(st.store, st.params, 0)
	assert.NoError(t, err)

	acc1a := sb.Account(acc1.Address())
	assert.Equal(t, acc1, acc1a)

	acc1a.IncSequence()
	acc1a.SetBalance(acc1a.Balance() + 1)

	sb.UpdateAccount(acc1a)

	acc2a := sb.Account(acc2.Address())
	assert.Nil(t, acc2a)

	assert.NoError(t, sb.CommitAndClear(st.validatorSet))

	acc1b, err := sb.store.Account(acc1.Address())
	assert.NoError(t, err)
	assert.Equal(t, acc1a, acc1b)

	// update state
	sb.UpdateAccount(acc2)
	acc22 := sb.Account(acc2.Address())
	assert.Equal(t, acc2, acc22)

	sb.Clear()
	assert.Equal(t, len(sb.accounts), 0)
	assert.Equal(t, len(sb.validators), 0)

	acc1c := sb.Account(acc1.Address())
	assert.Equal(t, acc1b, acc1c)
}

func TestValidatorChange(t *testing.T) {
	st, _ := mockState(t, nil)
	_, pub1, _ := crypto.GenerateTestKeyPair()
	_, pub2, _ := crypto.GenerateTestKeyPair()
	val1 := validator.NewValidator(pub1, 0, 0)
	val2 := validator.NewValidator(pub2, 1, 0)

	st.store.UpdateValidator(val1)
	sb, _ := newSandbox(st.store, st.params, 0)

	val1a := sb.Validator(val1.Address())
	assert.Equal(t, val1.Hash(), val1a.Hash())

	val1a.IncSequence()
	val1a.AddToStake(+1)

	sb.UpdateValidator(val1a)

	val2a := sb.Validator(val2.Address())
	assert.Nil(t, val2a)

	assert.NoError(t, sb.CommitAndClear(st.validatorSet))

	val1b, err := sb.store.Validator(val1.Address())
	assert.NoError(t, err)
	assert.Equal(t, val1a, val1b)

	// update state
	sb.UpdateValidator(val2)
	val22 := sb.Validator(val2.Address())
	assert.Equal(t, val2, val22)

	sb.Clear()
	assert.Equal(t, len(sb.validators), 0)
	assert.Equal(t, len(sb.validators), 0)

	val1c := sb.Validator(val1.Address())
	assert.Equal(t, val1b, val1c)
}

func TestAddValidatorToSet(t *testing.T) {
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

		sb.Clear()
		assert.Equal(t, sb.totalAccounts, 1)
		assert.Equal(t, st.store.TotalAccounts(), 1)

		acc = sb.MakeNewAccount(addr)
		sb.UpdateAccount(acc)

		assert.NoError(t, sb.CommitAndClear(st.validatorSet))
		assert.Equal(t, sb.totalAccounts, 2)
		assert.Equal(t, st.store.TotalAccounts(), 2)
	})

	t.Run("Should return error for invalid account number", func(t *testing.T) {
		addr, _, _ := crypto.GenerateTestKeyPair()
		// account.NewAcoount is a ssolo account constructor and
		// doesn't update total account counter in sandbox.
		// We should always use sandbox.MakeNewAccount to update total account counter
		acc := account.NewAccount(addr, 2)
		sb.UpdateAccount(acc)

		assert.Error(t, sb.CommitAndClear(st.validatorSet))
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

		sb.Clear()
		assert.Equal(t, sb.totalValidators, 1)
		assert.Equal(t, st.store.TotalValidators(), 1)

		val = sb.MakeNewValidator(pub)
		sb.UpdateValidator(val)

		assert.NoError(t, sb.CommitAndClear(st.validatorSet))
		assert.Equal(t, sb.totalValidators, 2)
		assert.Equal(t, st.store.TotalValidators(), 2)
	})

	t.Run("Should return error for invalid account number", func(t *testing.T) {
		_, pub, _ := crypto.GenerateTestKeyPair()
		// validator.NewValidator is a ssolo validator constructor and
		// doesn't update total validator counter in sandbox.
		// We should always use sandbox.MakeNewValidator to update total validator counter
		val := validator.NewValidator(pub, 2, 0)
		sb.UpdateValidator(val)

		assert.Error(t, sb.CommitAndClear(st.validatorSet))
		assert.Equal(t, sb.totalValidators, 2)
		assert.Equal(t, st.store.TotalValidators(), 2)
	})
}
