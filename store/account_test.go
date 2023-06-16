package store

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountCounter(t *testing.T) {
	setup(t)

	num := util.RandInt32(1000)
	acc, signer := account.GenerateTestAccount(num)

	t.Run("Add new account, should increase the total accounts number", func(t *testing.T) {
		assert.Zero(t, tStore.TotalAccounts())

		tStore.UpdateAccount(signer.Address(), acc)
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), int32(1))
	})

	t.Run("Update account, should not increase the total accounts number", func(t *testing.T) {
		acc.AddToBalance(1)
		tStore.UpdateAccount(signer.Address(), acc)

		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), int32(1))
	})

	t.Run("Get account", func(t *testing.T) {
		acc1, err := tStore.Account(signer.Address())
		assert.NoError(t, err)

		acc2, err := tStore.AccountByNumber(num)
		assert.NoError(t, err)

		assert.Equal(t, acc1.Hash(), acc2.Hash())
		assert.Equal(t, tStore.TotalAccounts(), int32(1))
		assert.True(t, tStore.HasAccount(signer.Address()))
	})
}

func TestAccountBatchSaving(t *testing.T) {
	setup(t)

	total := util.RandInt32(100) + 1
	t.Run("Add some accounts", func(t *testing.T) {
		for i := int32(0); i < total; i++ {
			acc, signer := account.GenerateTestAccount(i)
			tStore.UpdateAccount(signer.Address(), acc)
		}
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), total)
	})

	t.Run("Close and load db", func(t *testing.T) {
		tStore.Close()
		store, _ := NewStore(tStore.config, 21)
		assert.Equal(t, store.TotalAccounts(), total)
	})
}

func TestAccountByNumber(t *testing.T) {
	setup(t)

	total := util.RandInt32(100) + 1
	t.Run("Add some accounts", func(t *testing.T) {
		for i := int32(0); i < total; i++ {
			acc, signer := account.GenerateTestAccount(i)
			tStore.UpdateAccount(signer.Address(), acc)
		}
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), total)
	})

	t.Run("Get a random account", func(t *testing.T) {
		num := util.RandInt32(total)
		acc, err := tStore.AccountByNumber(num)
		assert.NoError(t, err)
		require.NotNil(t, acc)
		assert.Equal(t, acc.Number(), num)
	})

	t.Run("negative number", func(t *testing.T) {
		acc, err := tStore.AccountByNumber(-1)
		assert.Error(t, err)
		assert.Nil(t, acc)
	})

	t.Run("Non existing account", func(t *testing.T) {
		acc, err := tStore.AccountByNumber(total + 1)
		assert.Error(t, err)
		assert.Nil(t, acc)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		tStore.Close()
		store, _ := NewStore(tStore.config, 21)

		num := util.RandInt32(total)
		acc, err := store.AccountByNumber(num)
		assert.NoError(t, err)
		require.NotNil(t, acc)
		assert.Equal(t, acc.Number(), num)

		acc, err = tStore.AccountByNumber(total + 1)
		assert.Error(t, err)
		assert.Nil(t, acc)
	})
}

func TestAccountByAddress(t *testing.T) {
	setup(t)

	total := util.RandInt32(100) + 1
	t.Run("Add some accounts", func(t *testing.T) {
		for i := int32(0); i < total; i++ {
			acc, signer := account.GenerateTestAccount(i)
			tStore.UpdateAccount(signer.Address(), acc)
		}
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalAccounts(), total)
	})

	t.Run("Get random account", func(t *testing.T) {
		// num := util.RandInt32(total)
		// acc0, _ := tStore.AccountByNumber(num)
		// acc, err := tStore.Account(acc0.Address())
		// assert.NoError(t, err)
		// require.NotNil(t, acc)
		// assert.Equal(t, acc.Number(), num)
	})

	t.Run("Unknown address", func(t *testing.T) {
		acc, err := tStore.Account(crypto.GenerateTestAddress())
		assert.Error(t, err)
		assert.Nil(t, acc)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		// tStore.Close()
		// store, _ := NewStore(tStore.config, 21)

		// num := util.RandInt32(total)
		// acc0, _ := store.AccountByNumber(num)
		// acc, err := store.Account(acc0.Address())
		// assert.NoError(t, err)
		// require.NotNil(t, acc)
		// assert.Equal(t, acc.Number(), num)
	})
}
