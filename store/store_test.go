package store

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tStore *store

func setup(t *testing.T) {
	conf := &Config{
		Path: util.TempDirPath(),
	}
	s, err := NewStore(conf, 21)
	require.NoError(t, err)

	tStore = s.(*store)
	SaveTestBlocks(t, 10)
}

func SaveTestBlocks(t *testing.T, num int) {
	lastHeight, _ := tStore.LastCertificate()
	for i := 0; i < num; i++ {
		b := block.GenerateTestBlock(nil, nil)
		c := block.GenerateTestCertificate(b.Hash())

		tStore.SaveBlock(lastHeight+uint32(i+1), b, c)
		assert.NoError(t, tStore.WriteBatch())
	}
}

func TestReturnNilForNonExistingItems(t *testing.T) {
	setup(t)

	lastHeight, _ := tStore.LastCertificate()

	assert.Equal(t, tStore.BlockHash(lastHeight+1), hash.UndefHash)
	assert.Equal(t, tStore.BlockHash(0), hash.UndefHash)

	block, err := tStore.Block(hash.GenerateTestHash())
	assert.Error(t, err)
	assert.Nil(t, block)

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
func TestWriteAndClosePeacefully(t *testing.T) {
	setup(t)

	// After closing db, we should not crash
	assert.NoError(t, tStore.Close())
	assert.Error(t, tStore.WriteBatch())
}
func TestRetrieveBlockAndTransactions(t *testing.T) {
	setup(t)

	height, _ := tStore.LastCertificate()
	storedBlock, err := tStore.Block(tStore.BlockHash(height))
	assert.NoError(t, err)
	assert.Equal(t, height, storedBlock.Height)
	block := storedBlock.ToBlock()
	for _, trx := range block.Transactions() {
		storedTx, err := tStore.Transaction(trx.ID())
		assert.NoError(t, err)
		assert.Equal(t, storedTx.TxID, trx.ID())
		assert.Equal(t, storedTx.BlockTime, uint32(block.Header().Time().Unix()))
		assert.Equal(t, storedTx.ToTx().ID(), trx.ID())
	}
}

func TestRetrieveAccount(t *testing.T) {
	setup(t)

	acc, _ := account.GenerateTestAccount(util.RandInt32(10000))

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
	assert.Equal(t, tStore.TotalAccounts(), int32(1))

	// Should not crash
	assert.NoError(t, tStore.Close())
	_, err := tStore.Account(acc.Address())
	assert.Error(t, err)
}

func TestRetrieveValidator(t *testing.T) {
	setup(t)

	val, _ := validator.GenerateTestValidator(util.RandInt32(1000))

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

	assert.Equal(t, tStore.TotalValidators(), int32(1))
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
		acc, _ := account.GenerateTestAccount(int32(i))
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
		val, _ := validator.GenerateTestValidator(int32(i))
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

func TestFindBlockHashByStamp(t *testing.T) {
	setup(t)

	hash1 := tStore.BlockHash(1)

	h, ok := tStore.FindBlockHashByStamp(hash1.Stamp())
	assert.Equal(t, h, hash1)
	assert.True(t, ok)

	height, ok := tStore.FindBlockHeightByStamp(hash1.Stamp())
	assert.Equal(t, height, uint32(1))
	assert.True(t, ok)

	// Reopen the store
	tStore.Close()
	s, _ := NewStore(tStore.config, 21)
	tStore = s.(*store)

	h, ok = tStore.FindBlockHashByStamp(hash1.Stamp())
	assert.Equal(t, h, hash1)
	assert.True(t, ok)

	height, ok = tStore.FindBlockHeightByStamp(hash1.Stamp())
	assert.Equal(t, height, uint32(1))
	assert.True(t, ok)

	// Saving more blocks
	SaveTestBlocks(t, 12)
	hash2 := tStore.BlockHash(2)
	hash14 := tStore.BlockHash(14)
	hash22 := tStore.BlockHash(22)

	// First block should be removed from the list
	h, ok = tStore.FindBlockHashByStamp(hash1.Stamp())
	assert.Equal(t, h, hash.UndefHash)
	assert.False(t, ok)

	height, ok = tStore.FindBlockHeightByStamp(hash1.Stamp())
	assert.Zero(t, height)
	assert.False(t, ok)

	h, ok = tStore.FindBlockHashByStamp(hash2.Stamp())
	assert.Equal(t, h, hash2)
	assert.True(t, ok)

	h, ok = tStore.FindBlockHashByStamp(hash14.Stamp())
	assert.Equal(t, h, hash14)
	assert.True(t, ok)

	h, ok = tStore.FindBlockHashByStamp(hash22.Stamp())
	assert.Equal(t, h, hash22)
	assert.True(t, ok)

	height, ok = tStore.FindBlockHeightByStamp(hash2.Stamp())
	assert.Equal(t, height, uint32(2))
	assert.True(t, ok)

	height, ok = tStore.FindBlockHeightByStamp(hash14.Stamp())
	assert.Equal(t, height, uint32(14))
	assert.True(t, ok)

	height, ok = tStore.FindBlockHeightByStamp(hash22.Stamp())
	assert.Equal(t, height, uint32(22))
	assert.True(t, ok)

	// Reopen the store
	tStore.Close()
	s, _ = NewStore(tStore.config, 21)
	tStore = s.(*store)

	height, ok = tStore.FindBlockHeightByStamp(hash2.Stamp())
	assert.Equal(t, height, uint32(2))
	assert.True(t, ok)

	SaveTestBlocks(t, 1)

	// Second block should bre removed from th list
	height, ok = tStore.FindBlockHeightByStamp(hash2.Stamp())
	assert.Zero(t, height)
	assert.False(t, ok)

	// Genesis block
	h, ok = tStore.FindBlockHashByStamp(hash.UndefHash.Stamp())
	assert.Equal(t, h, hash.UndefHash)
	assert.True(t, ok)

	height, ok = tStore.FindBlockHeightByStamp(hash.UndefHash.Stamp())
	assert.Zero(t, h, height)
	assert.True(t, ok)
}
