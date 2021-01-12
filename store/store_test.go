package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var tStore *Store

func setup(t *testing.T) {
	conf := TestConfig()
	store, err := NewStore(conf)
	assert.NoError(t, err)

	tStore = store
}

func TestReturnNilForNonExistingItems(t *testing.T) {
	setup(t)

	b, txs := block.GenerateTestBlock(nil, nil)
	h := util.RandInt(10000)
	block, err := tStore.Block(h)
	assert.Error(t, err)
	assert.Nil(t, block)

	height, err := tStore.BlockHeight(b.Hash())
	assert.Error(t, err)
	assert.Equal(t, height, -1)

	tx, err := tStore.Transaction(txs[0].ID())
	assert.Error(t, err)
	assert.Nil(t, tx)

	acc, err := tStore.Account(b.Header().ProposerAddress())
	assert.Error(t, err)
	assert.Nil(t, acc)

	val, err := tStore.Validator(b.Header().ProposerAddress())
	assert.Error(t, err)
	assert.Nil(t, val)

	assert.NoError(t, tStore.Close())
}
func TestRetrieveBlockAndTransactions(t *testing.T) {
	setup(t)

	b, txs := block.GenerateTestBlock(nil, nil)
	h := util.RandInt(10000)
	err := tStore.SaveBlock(*b, h)
	assert.NoError(t, err)

	for _, trx := range txs {
		r := trx.GenerateReceipt(tx.Ok, b.Hash())
		ctrx := tx.CommittedTx{Tx: trx, Receipt: r}
		tStore.SaveTransaction(ctrx)
	}

	h2, err := tStore.BlockHeight(b.Hash())
	assert.NoError(t, err)
	b2, err := tStore.Block(h2)
	assert.NoError(t, err)
	assert.Equal(t, b.Hash(), b2.Hash())
	bz1, _ := b.Encode()
	bz2, _ := b2.Encode()
	assert.Equal(t, bz1, bz2)
	assert.Equal(t, h, h2)

	for _, trx := range txs {
		r := trx.GenerateReceipt(tx.Ok, b.Hash())
		ctrx2, err := tStore.Transaction(trx.ID())
		assert.NoError(t, err)

		assert.Equal(t, trx.ID(), ctrx2.Tx.ID())
		assert.Equal(t, r, ctrx2.Receipt)
	}

	assert.NoError(t, tStore.Close())
}

func TestRetrieveAccount(t *testing.T) {
	setup(t)

	acc, _ := account.GenerateTestAccount(util.RandInt(10000))

	t.Run("Add account, should able to retrieve", func(t *testing.T) {
		assert.False(t, tStore.HasAccount(acc.Address()))
		tStore.UpdateAccount(acc)
		assert.True(t, tStore.HasAccount(acc.Address()))
		acc2, err := tStore.Account(acc.Address())
		assert.NoError(t, err)
		assert.Equal(t, acc, acc2)
	})

	t.Run("Update account, should update database", func(t *testing.T) {
		acc.AddToBalance(1)
		tStore.UpdateAccount(acc)

		acc2, err := tStore.Account(acc.Address())
		assert.NoError(t, err)
		assert.Equal(t, acc, acc2)
	})
}

func TestRetrieveValidator(t *testing.T) {
	setup(t)

	val, _ := validator.GenerateTestValidator(util.RandInt(1000))

	t.Run("Add validator, should able to retrieve", func(t *testing.T) {
		assert.False(t, tStore.HasValidator(val.Address()))
		tStore.UpdateValidator(val)
		assert.True(t, tStore.HasValidator(val.Address()))
		val2, err := tStore.Validator(val.Address())
		assert.NoError(t, err)
		assert.Equal(t, val.Hash(), val2.Hash())
	})

	t.Run("Update validator, should update database", func(t *testing.T) {
		val.AddToStake(1)
		tStore.UpdateValidator(val)

		val2, err := tStore.Validator(val.Address())
		assert.NoError(t, err)
		assert.Equal(t, val.Hash(), val2.Hash())
	})
}

func TestIterateAccounts(t *testing.T) {
	setup(t)

	accs1 := []crypto.Hash{}
	for i := 0; i < 10; i++ {
		acc, _ := account.GenerateTestAccount(i)
		tStore.UpdateAccount(acc)

		accs1 = append(accs1, acc.Hash())
	}

	accs2 := []crypto.Hash{}
	tStore.IterateAccounts(func(acc *account.Account) bool {
		accs2 = append(accs2, acc.Hash())
		return false
	})

	assert.ElementsMatch(t, accs1, accs2)
}

func TestIterateValidators(t *testing.T) {
	setup(t)

	vals1 := []crypto.Hash{}
	for i := 0; i < 10; i++ {
		val, _ := validator.GenerateTestValidator(i)
		tStore.UpdateValidator(val)

		vals1 = append(vals1, val.Hash())
	}

	vals2 := []crypto.Hash{}
	tStore.IterateValidators(func(val *validator.Validator) bool {
		vals2 = append(vals2, val.Hash())
		return false
	})

	assert.ElementsMatch(t, vals1, vals2)
}
