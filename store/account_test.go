package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/util"
)

func TestRetreiveAccount(t *testing.T) {
	store, _ := newAccountStore(util.TempDirPath())

	acc, _ := account.GenerateTestAccount(util.RandInt(10000))

	t.Run("Add account, should able to retrieve", func(t *testing.T) {
		assert.False(t, store.hasAccount(acc.Address()))
		store.updateAccount(acc)
		assert.True(t, store.hasAccount(acc.Address()))
		acc2, err := store.account(acc.Address())
		assert.NoError(t, err)
		assert.Equal(t, acc, acc2)
	})

	t.Run("Update account, should update database", func(t *testing.T) {
		acc.AddToBalance(1)
		store.updateAccount(acc)

		acc2, err := store.account(acc.Address())
		assert.NoError(t, err)
		assert.Equal(t, acc, acc2)
	})
}

func TestAccountCounter(t *testing.T) {
	store, _ := newAccountStore(util.TempDirPath())

	acc, _ := account.GenerateTestAccount(0)

	t.Run("Update count after adding new account", func(t *testing.T) {
		assert.Equal(t, store.total, store.countAccounts())
		assert.Equal(t, store.total, 0)

		store.updateAccount(acc)
		assert.Equal(t, store.total, store.countAccounts())
		assert.Equal(t, store.total, 1)
	})

	t.Run("Update account, should not increatse counter", func(t *testing.T) {
		acc.AddToBalance(1)
		store.updateAccount(acc)

		store.updateAccount(acc)
		assert.Equal(t, store.total, store.countAccounts())
		assert.Equal(t, store.total, 1)
	})
}

func TestAccountBatchSaving(t *testing.T) {
	path := util.TempDirPath()
	store, _ := newAccountStore(path)

	t.Run("Add 100 accounts", func(t *testing.T) {

		for i := 0; i < 100; i++ {
			acc, _ := account.GenerateTestAccount(i)
			store.updateAccount(acc)
		}

		assert.Equal(t, store.total, store.countAccounts())
		assert.Equal(t, store.total, 100)
	})
	t.Run("Close and load db", func(t *testing.T) {
		store.close()
		store, _ = newAccountStore(path)

		assert.Equal(t, store.total, store.countAccounts())
		assert.Equal(t, store.total, 100)
	})
}
