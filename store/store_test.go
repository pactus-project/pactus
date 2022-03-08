package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var tStore *store

func setup(t *testing.T) {
	conf := TestConfig()
	s, err := NewStore(conf, 21)
	require.NoError(t, err)

	tStore = s.(*store)
	SaveTestBlocks(t, 10)
}

func SaveTestBlocks(t *testing.T, num int) {
	lastHeight, _ := tStore.LastCertificate()
	for i := 0; i < num; i++ {
		b, txs := block.GenerateTestBlock(nil, nil)
		c := block.GenerateTestCertificate(b.Hash())

		tStore.SaveBlock(lastHeight+i+1, b, c)
		for _, tx := range txs {
			tStore.SaveTransaction(tx)
		}
		assert.NoError(t, tStore.WriteBatch())
	}
}

func TestReturnNilForNonExistingItems(t *testing.T) {
	setup(t)

	lastHeight, _ := tStore.LastCertificate()
	block, err := tStore.Block(lastHeight + 1)
	assert.Error(t, err)
	assert.Nil(t, block)

	block, err = tStore.Block(0)
	assert.Error(t, err)
	assert.Nil(t, block)

	height, err := tStore.BlockHeight(hash.GenerateTestHash())
	assert.Error(t, err)
	assert.Equal(t, height, -1)

	tx, err := tStore.Transaction(hash.GenerateTestHash())
	assert.Error(t, err)
	assert.Nil(t, tx)

	acc, err := tStore.Account(crypto.GenerateTestAddress())
	assert.Error(t, err)
	assert.Nil(t, acc)

	val, err := tStore.Validator(crypto.GenerateTestAddress())
	assert.Error(t, err)
	assert.Nil(t, val)

	assert.NoError(t, tStore.Close())
}

func TestRetrieveBlockAndTransactions(t *testing.T) {
	setup(t)

	height, cert := tStore.LastCertificate()

	h2, err := tStore.BlockHeight(cert.BlockHash())
	assert.NoError(t, err)
	b2, err := tStore.Block(h2)
	assert.NoError(t, err)
	assert.Equal(t, cert.BlockHash(), b2.Hash())
	assert.Equal(t, height, h2)

	for _, id := range b2.TxIDs().IDs() {
		trx2, err := tStore.Transaction(id)
		assert.NoError(t, err)

		assert.Equal(t, id, trx2.ID())
	}

	// After closing db, we should not crash
	assert.NoError(t, tStore.Close())
	assert.Error(t, tStore.WriteBatch())
	_, err = tStore.Block(h2 + 1)
	assert.Error(t, err)
}

func TestRetrieveAccount(t *testing.T) {
	setup(t)

	acc, _ := account.GenerateTestAccount(util.RandInt(10000))

	t.Run("Add account, should able to retrieve", func(t *testing.T) {
		assert.False(t, tStore.HasAccount(acc.Address()))
		tStore.UpdateAccount(acc)
		assert.NoError(t, tStore.WriteBatch())
		assert.True(t, tStore.HasAccount(acc.Address()))
		acc2, err := tStore.Account(acc.Address())
		assert.NoError(t, err)
		assert.Equal(t, acc, acc2)
	})

	t.Run("Update account, should update database", func(t *testing.T) {
		acc.AddToBalance(1)
		tStore.UpdateAccount(acc)
		assert.NoError(t, tStore.WriteBatch())
		acc2, err := tStore.Account(acc.Address())
		assert.NoError(t, err)
		assert.Equal(t, acc, acc2)
	})
	assert.Equal(t, tStore.TotalAccounts(), 1)

	// Should not crash
	assert.NoError(t, tStore.Close())
	_, err := tStore.Account(acc.Address())
	assert.Error(t, err)

}

func TestRetrieveValidator(t *testing.T) {
	setup(t)

	val, _ := validator.GenerateTestValidator(util.RandInt(1000))

	t.Run("Add validator, should able to retrieve", func(t *testing.T) {
		assert.False(t, tStore.HasValidator(val.Address()))
		tStore.UpdateValidator(val)
		assert.NoError(t, tStore.WriteBatch())
		assert.True(t, tStore.HasValidator(val.Address()))
		val2, err := tStore.Validator(val.Address())
		assert.NoError(t, err)
		assert.Equal(t, val.Hash(), val2.Hash())
	})

	t.Run("Update validator, should update database", func(t *testing.T) {
		val.AddToStake(1)
		tStore.UpdateValidator(val)
		assert.NoError(t, tStore.WriteBatch())
		val2, err := tStore.Validator(val.Address())
		assert.NoError(t, err)
		assert.Equal(t, val.Hash(), val2.Hash())
	})

	assert.Equal(t, tStore.TotalValidators(), 1)
	val2, _ := tStore.ValidatorByNumber(val.Number())
	assert.Equal(t, val.Hash(), val2.Hash())

	assert.NoError(t, tStore.Close())
	_, err := tStore.Validator(val.Address())
	assert.Error(t, err)
}

func TestIterateAccounts(t *testing.T) {
	setup(t)

	accs1 := []hash.Hash{}
	for i := 0; i < 10; i++ {
		acc, _ := account.GenerateTestAccount(i)
		tStore.UpdateAccount(acc)
		assert.NoError(t, tStore.WriteBatch())
		accs1 = append(accs1, acc.Hash())
	}

	stopped := false
	tStore.IterateAccounts(func(acc *account.Account) bool {
		if acc.Hash().EqualsTo(accs1[0]) {
			stopped = true
		}
		return stopped
	})
	assert.True(t, stopped)

	accs2 := []hash.Hash{}
	tStore.IterateAccounts(func(acc *account.Account) bool {
		accs2 = append(accs2, acc.Hash())
		return false
	})

	assert.ElementsMatch(t, accs1, accs2)
}

func TestIterateValidators(t *testing.T) {
	setup(t)

	vals1 := []hash.Hash{}
	for i := 0; i < 10; i++ {
		val, _ := validator.GenerateTestValidator(i)
		tStore.UpdateValidator(val)
		assert.NoError(t, tStore.WriteBatch())
		vals1 = append(vals1, val.Hash())
	}

	stopped := false
	tStore.IterateValidators(func(val *validator.Validator) bool {
		if val.Hash().EqualsTo(vals1[0]) {
			stopped = true
		}
		return stopped
	})
	assert.True(t, stopped)

	vals2 := []hash.Hash{}
	tStore.IterateValidators(func(val *validator.Validator) bool {
		vals2 = append(vals2, val.Hash())
		return false
	})

	assert.ElementsMatch(t, vals1, vals2)
}

// func TestReestoreLastInfo(t *testing.T) {
// 	setup(t)

// 	assert.Nil(t, tStore.RestoreLastInfo())
// 	tStore.SaveLastInfo([]byte{1})
// 	assert.Nil(t, tStore.RestoreLastInfo())
// 	assert.NoError(t, tStore.WriteBatch())
// 	assert.NotNil(t, tStore.RestoreLastInfo())
// }
