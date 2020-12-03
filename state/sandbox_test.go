package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/validator"
)

func TestAccountChange(t *testing.T) {
	st, _ := mockState(t, nil)
	acc1, _ := account.GenerateTestAccount()
	acc2, _ := account.GenerateTestAccount()

	st.store.UpdateAccount(acc1)
	sb := newSandbox(st.store, st)

	acc1a := sb.Account(acc1.Address())
	assert.Equal(t, acc1, acc1a)

	acc1a.IncSequence()
	acc1a.SetBalance(acc1a.Balance() + 1)

	sb.UpdateAccount(acc1a)

	acc2a := sb.Account(acc2.Address())
	assert.Nil(t, acc2a)

	sb.commit(st.validatorSet)

	acc1b, err := sb.store.Account(acc1.Address())
	assert.NoError(t, err)
	assert.Equal(t, acc1a, acc1b)

	// update state
	sb.UpdateAccount(acc2)
	acc22 := sb.Account(acc2.Address())
	assert.Equal(t, acc2, acc22)

	sb.reset()
	assert.Equal(t, len(sb.accounts), 0)
	assert.Equal(t, len(sb.validators), 0)

	acc1c := sb.Account(acc1.Address())
	assert.Equal(t, acc1b, acc1c)
}

func TestValidatorChange(t *testing.T) {
	st, _ := mockState(t, nil)
	val1, _ := validator.GenerateTestValidator()
	val2, _ := validator.GenerateTestValidator()

	st.store.UpdateValidator(val1)
	sb := newSandbox(st.store, st)

	val1a := sb.Validator(val1.Address())
	assert.Equal(t, val1.Hash(), val1a.Hash())

	val1a.IncSequence()
	val1a.AddToStake(+1)

	sb.UpdateValidator(val1a)

	val2a := sb.Validator(val2.Address())
	assert.Nil(t, val2a)

	sb.commit(st.validatorSet)

	val1b, err := sb.store.Validator(val1.Address())
	assert.NoError(t, err)
	assert.Equal(t, val1a, val1b)

	// update state
	sb.UpdateValidator(val2)
	val22 := sb.Validator(val2.Address())
	assert.Equal(t, val2, val22)

	sb.reset()
	assert.Equal(t, len(sb.validators), 0)
	assert.Equal(t, len(sb.validators), 0)

	val1c := sb.Validator(val1.Address())
	assert.Equal(t, val1b, val1c)
}

func TestAddValidatorToSet(t *testing.T) {
	st, _ := mockState(t, nil)
	val1, _ := validator.GenerateTestValidator()
	val2, _ := validator.GenerateTestValidator()

	st.validatorSet.Join(val1)
	sb := newSandbox(st.store, st)

	sb.AddToSet(val2)
	// Still is not in set
	assert.Nil(t, st.validatorSet.Validator(val2.Address()))

	sb.commit(st.validatorSet)
	// Still is not in set
	assert.Nil(t, st.validatorSet.Validator(val2.Address()))

}
