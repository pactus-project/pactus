package store

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
)

func TestAccountCounter(t *testing.T) {
	store, _ := NewStore(TestConfig())

	acc, _ := account.GenerateTestAccount(0)

	t.Run("Update count after adding new account", func(t *testing.T) {
		assert.Equal(t, store.TotalAccounts(), 0)

		store.UpdateAccount(acc)
		assert.Equal(t, store.TotalAccounts(), 1)
	})

	t.Run("Update account, should not increase counter", func(t *testing.T) {
		acc.AddToBalance(1)
		store.UpdateAccount(acc)
		assert.Equal(t, store.TotalAccounts(), 1)
	})
}

func TestAccountBatchSaving(t *testing.T) {
	store, _ := NewStore(TestConfig())

	t.Run("Add 100 accounts", func(t *testing.T) {

		for i := 0; i < 100; i++ {
			acc, _ := account.GenerateTestAccount(i)
			store.UpdateAccount(acc)
		}
		assert.Equal(t, store.TotalAccounts(), 100)
	})
	t.Run("Close and load db", func(t *testing.T) {
		store.Close()
		fmt.Println(TestConfig().Path)
		store, _ := NewStore(TestConfig())
		fmt.Println(TestConfig().Path)
		assert.Equal(t, store.TotalAccounts(), 100)
	})
}
