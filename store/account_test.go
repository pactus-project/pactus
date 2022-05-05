package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/types/account"
	"github.com/zarbchain/zarb-go/util"
)

func TestAccountCounter(t *testing.T) {
	setup(t)

	acc, _ := account.GenerateTestAccount(util.RandInt32(1000))

	t.Run("Update count after adding new account", func(t *testing.T) {
		assert.Zero(t, tStore.TotalAccounts())

		tStore.UpdateAccount(acc)
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), int32(1))
	})

	t.Run("Update account, should not increase counter", func(t *testing.T) {
		acc.AddToBalance(1)
		tStore.UpdateAccount(acc)
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), int32(1))
	})

	t.Run("Get account", func(t *testing.T) {
		assert.True(t, tStore.HasAccount(acc.Address()))
		acc2, err := tStore.Account(acc.Address())
		assert.NoError(t, err)
		assert.Equal(t, acc2.Hash(), acc.Hash())
	})
}

func TestAccountBatchSaving(t *testing.T) {
	setup(t)

	t.Run("Add 100 accounts", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			acc, _ := account.GenerateTestAccount(int32(i))
			tStore.UpdateAccount(acc)
		}
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), int32(100))
	})
	t.Run("Close and load db", func(t *testing.T) {
		tStore.Close()
		store, _ := NewStore(tStore.config, 21)
		assert.Equal(t, store.TotalAccounts(), int32(100))
	})
}
