package store

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestAccountCounter(t *testing.T) {
	setup(t)

	acc, signer := account.GenerateTestAccount(util.RandInt32(1000))

	t.Run("Update count after adding new account", func(t *testing.T) {
		assert.Zero(t, tStore.TotalAccounts())

		tStore.UpdateAccount(signer.Address(), acc)
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), int32(1))
	})

	t.Run("Update account, should not increase counter", func(t *testing.T) {
		acc.AddToBalance(1)
		tStore.UpdateAccount(signer.Address(), acc)
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), int32(1))
	})

	t.Run("Get account", func(t *testing.T) {
		assert.True(t, tStore.HasAccount(signer.Address()))
		acc2, err := tStore.Account(signer.Address())
		assert.NoError(t, err)
		assert.Equal(t, acc2.Hash(), acc.Hash())
	})
}

func TestAccountBatchSaving(t *testing.T) {
	setup(t)

	t.Run("Add 100 accounts", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			acc, signer := account.GenerateTestAccount(int32(i))
			tStore.UpdateAccount(signer.Address(), acc)
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

func TestAccountByNumber(t *testing.T) {
	setup(t)

	t.Run("Add some accounts", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			val, signer := account.GenerateTestAccount(int32(i))
			tStore.UpdateAccount(signer.Address(), val)
		}
		assert.NoError(t, tStore.WriteBatch())

		v, err := tStore.AccountByNumber(5)
		assert.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, v.Number(), int32(5))

		v, err = tStore.AccountByNumber(11)
		assert.Error(t, err)
		assert.Nil(t, v)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		tStore.Close()
		store, _ := NewStore(tStore.config, 21)

		v, err := store.AccountByNumber(5)
		assert.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, v.Number(), int32(5))

		v, err = tStore.AccountByNumber(11)
		assert.Error(t, err)
		assert.Nil(t, v)
	})
}
