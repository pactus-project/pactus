package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
)

func TestAccountChange(t *testing.T) {
	st, _ := mockState(t)
	cache := newCache(st.(*state).store)
	acc1 := account.GenerateTestAccount()

	acc := cache.Account(acc1.Address())
	assert.Nil(t, acc, nil)

	// update cache
	cache.UpdateAccount(acc1)
	acc11 := cache.Account(acc1.Address())
	assert.Equal(t, acc1, acc11)

	// update state
	acc2 := account.GenerateTestAccount()
	cache.UpdateAccount(acc2)
	acc22 := cache.Account(acc2.Address())
	assert.Equal(t, acc2, acc22)

	cache.reset()
	assert.Equal(t, cache.accChanges.Len(), 0)
	assert.Equal(t, cache.valChanges.Len(), 0)
}
