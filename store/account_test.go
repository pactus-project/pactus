package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
)

func TestAccountCounter(t *testing.T) {
	setup(t)

	acc, _ := account.GenerateTestAccount(0)

	t.Run("Update count after adding new account", func(t *testing.T) {
		assert.Equal(t, tStore.TotalAccounts(), 0)

		tStore.UpdateAccount(acc)
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), 1)
	})

	t.Run("Update account, should not increase counter", func(t *testing.T) {
		acc.AddToBalance(1)
		tStore.UpdateAccount(acc)
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), 1)
	})
}

func TestAccountBatchSaving(t *testing.T) {
	setup(t)

	t.Run("Add 100 accounts", func(t *testing.T) {

		for i := 0; i < 100; i++ {
			acc, _ := account.GenerateTestAccount(i)
			tStore.UpdateAccount(acc)
		}
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), 100)
	})
	t.Run("Close and load db", func(t *testing.T) {
		tStore.Close()
		store, _ := NewStore(tStore.config, 21)
		assert.Equal(t, store.TotalAccounts(), 100)
	})
}
