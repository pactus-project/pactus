package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
)

func TestAccountChange(t *testing.T) {
	st, _ := mockState(t, nil)
	sb := newSandbox(st.store, st)
	acc1 := account.GenerateTestAccount()

	acc := sb.Account(acc1.Address())
	assert.Nil(t, acc, nil)

	// update sb
	sb.UpdateAccount(acc1)
	acc11 := sb.Account(acc1.Address())
	assert.Equal(t, acc1, acc11)

	// update state
	acc2 := account.GenerateTestAccount()
	sb.UpdateAccount(acc2)
	acc22 := sb.Account(acc2.Address())
	assert.Equal(t, acc2, acc22)

	sb.reset()
	assert.Equal(t, len(sb.accounts), 0)
	assert.Equal(t, len(sb.validators), 0)
}
