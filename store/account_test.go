package store

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/account"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountCounter(t *testing.T) {
	td := setup(t)

	num := td.RandInt32(1000)
	acc, signer := td.GenerateTestAccount(num)

	t.Run("Add new account, should increase the total accounts number", func(t *testing.T) {
		assert.Zero(t, td.store.TotalAccounts())

		td.store.UpdateAccount(signer.Address(), acc)
		assert.NoError(t, td.store.WriteBatch())
		assert.Equal(t, td.store.TotalAccounts(), int32(1))
	})

	t.Run("Update account, should not increase the total accounts number", func(t *testing.T) {
		acc.AddToBalance(1)
		td.store.UpdateAccount(signer.Address(), acc)

		assert.NoError(t, td.store.WriteBatch())
		assert.Equal(t, td.store.TotalAccounts(), int32(1))
	})

	t.Run("Get account", func(t *testing.T) {
		acc1, err := td.store.Account(signer.Address())
		assert.NoError(t, err)

		acc2, err := td.store.AccountByNumber(num)
		assert.NoError(t, err)

		assert.Equal(t, acc1.Hash(), acc2.Hash())
		assert.Equal(t, td.store.TotalAccounts(), int32(1))
		assert.True(t, td.store.HasAccount(signer.Address()))
	})
}

func TestAccountBatchSaving(t *testing.T) {
	td := setup(t)

	total := td.RandInt32NonZero(100)
	t.Run("Add some accounts", func(t *testing.T) {
		for i := int32(0); i < total; i++ {
			acc, signer := td.GenerateTestAccount(i)
			td.store.UpdateAccount(signer.Address(), acc)
		}
		assert.NoError(t, td.store.WriteBatch())
		assert.Equal(t, td.store.TotalAccounts(), total)
	})

	t.Run("Close and load db", func(t *testing.T) {
		td.store.Close()
		store, _ := NewStore(td.store.config, 21)
		assert.Equal(t, store.TotalAccounts(), total)
	})
}

func TestAccountByNumber(t *testing.T) {
	td := setup(t)

	total := td.RandInt32NonZero(100)
	t.Run("Add some accounts", func(t *testing.T) {
		for i := int32(0); i < total; i++ {
			acc, signer := td.GenerateTestAccount(i)
			td.store.UpdateAccount(signer.Address(), acc)
		}
		assert.NoError(t, td.store.WriteBatch())
		assert.Equal(t, td.store.TotalAccounts(), total)
	})

	t.Run("Get a random account", func(t *testing.T) {
		num := td.RandInt32(total)
		acc, err := td.store.AccountByNumber(num)
		assert.NoError(t, err)
		require.NotNil(t, acc)
		assert.Equal(t, acc.Number(), num)
	})

	t.Run("negative number", func(t *testing.T) {
		acc, err := td.store.AccountByNumber(-1)
		assert.Error(t, err)
		assert.Nil(t, acc)
	})

	t.Run("Non existing account", func(t *testing.T) {
		acc, err := td.store.AccountByNumber(total + 1)
		assert.Error(t, err)
		assert.Nil(t, acc)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		td.store.Close()
		store, _ := NewStore(td.store.config, 21)

		num := td.RandInt32(total)
		acc, err := store.AccountByNumber(num)
		assert.NoError(t, err)
		require.NotNil(t, acc)
		assert.Equal(t, acc.Number(), num)

		acc, err = td.store.AccountByNumber(total + 1)
		assert.Error(t, err)
		assert.Nil(t, acc)
	})
}

func TestAccountByAddress(t *testing.T) {
	td := setup(t)

	total := td.RandInt32NonZero(100)
	var lastAddr crypto.Address
	t.Run("Add some accounts", func(t *testing.T) {
		for i := int32(0); i < total; i++ {
			acc, signer := td.GenerateTestAccount(i)
			td.store.UpdateAccount(signer.Address(), acc)

			lastAddr = signer.Address()
		}
		assert.NoError(t, td.store.WriteBatch())
		assert.Equal(t, td.store.TotalAccounts(), total)
	})

	t.Run("Get random account", func(t *testing.T) {
		acc, err := td.store.Account(lastAddr)
		assert.NoError(t, err)
		require.NotNil(t, acc)
		assert.Equal(t, acc.Number(), total-1)
	})

	t.Run("Unknown address", func(t *testing.T) {
		acc, err := td.store.Account(td.RandAddress())
		assert.Error(t, err)
		assert.Nil(t, acc)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		td.store.Close()
		store, _ := NewStore(td.store.config, 21)

		acc, err := store.Account(lastAddr)
		assert.NoError(t, err)
		require.NotNil(t, acc)
		assert.Equal(t, acc.Number(), total-1)
	})
}

func TestIterateAccounts(t *testing.T) {
	td := setup(t)

	total := td.RandInt32NonZero(100)
	accs1 := []hash.Hash{}
	for i := int32(0); i < total; i++ {
		acc, signer := td.GenerateTestAccount(i)
		td.store.UpdateAccount(signer.Address(), acc)
		accs1 = append(accs1, acc.Hash())
	}
	assert.NoError(t, td.store.WriteBatch())

	accs2 := []hash.Hash{}
	td.store.IterateAccounts(func(_ crypto.Address, acc *account.Account) bool {
		accs2 = append(accs2, acc.Hash())
		return false
	})
	assert.ElementsMatch(t, accs1, accs2)

	stopped := false
	td.store.IterateAccounts(func(addr crypto.Address, acc *account.Account) bool {
		if acc.Hash().EqualsTo(accs1[0]) {
			stopped = true
		}
		return stopped
	})
	assert.True(t, stopped)
}

func TestAccountDeepCopy(t *testing.T) {
	td := setup(t)

	num := td.RandInt32(1000)
	acc1, signer := td.GenerateTestAccount(num)
	td.store.UpdateAccount(signer.Address(), acc1)

	acc2, _ := td.store.AccountByNumber(num)
	acc2.AddToBalance(1)
	assert.NotEqual(t, td.store.accountStore.numberMap[num].Hash(), acc2.Hash())

	acc3, _ := td.store.Account(signer.Address())
	acc3.AddToBalance(1)
	assert.NotEqual(t, td.store.accountStore.numberMap[num].Hash(), acc3.Hash())
}
