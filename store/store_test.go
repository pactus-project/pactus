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

func TestBlockHash(t *testing.T) {
	setup(t)

	sb, _ := tStore.Block(1)

	assert.Equal(t, tStore.BlockHash(0), hash.UndefHash)
	assert.Equal(t, tStore.BlockHash(util.MaxUint32), hash.UndefHash)
	assert.Equal(t, tStore.BlockHash(1), sb.BlockHash)
}

func TestBlockHeight(t *testing.T) {
	setup(t)

	sb, _ := tStore.Block(1)

	assert.Equal(t, tStore.BlockHeight(hash.UndefHash), uint32(0))
	assert.Equal(t, tStore.BlockHeight(hash.GenerateTestHash()), uint32(0))
	assert.Equal(t, tStore.BlockHeight(sb.BlockHash), uint32(1))
}

func TestReturnNilForNonExistingItems(t *testing.T) {
	setup(t)

	lastHeight, _ := tStore.LastCertificate()

	assert.Equal(t, tStore.BlockHash(lastHeight+1), hash.UndefHash)
	assert.Equal(t, tStore.BlockHash(0), hash.UndefHash)

	block, err := tStore.Block(lastHeight + 1)
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
	storedBlock, err := tStore.Block(height)
	assert.NoError(t, err)
	assert.Equal(t, height, storedBlock.Height)
	block := storedBlock.ToBlock()
	for _, trx := range block.Transactions() {
		storedTx, err := tStore.Transaction(trx.ID())
		assert.NoError(t, err)
		assert.Equal(t, storedTx.TxID, trx.ID())
		assert.Equal(t, storedTx.BlockTime, block.Header().UnixTime())
		assert.Equal(t, storedTx.ToTx().ID(), trx.ID())
	}
}

func TestRetrieveAccount(t *testing.T) {
	setup(t)

	acc, signer := account.GenerateTestAccount(util.RandInt32(10000))

	t.Run("Add account, should able to retrieve", func(t *testing.T) {
		assert.False(t, tStore.HasAccount(signer.Address()))
		tStore.UpdateAccount(signer.Address(), acc)
		assert.NoError(t, tStore.WriteBatch())
		assert.True(t, tStore.HasAccount(signer.Address()))
		acc2, err := tStore.Account(signer.Address())
		assert.NoError(t, err)
		assert.Equal(t, acc, acc2)
	})

	t.Run("Update account, should update database", func(t *testing.T) {
		acc.AddToBalance(1)
		tStore.UpdateAccount(signer.Address(), acc)
		assert.NoError(t, tStore.WriteBatch())
		acc2, err := tStore.Account(signer.Address())
		assert.NoError(t, err)
		assert.Equal(t, acc, acc2)
	})
	assert.Equal(t, tStore.TotalAccounts(), int32(1))

	// Should not crash
	assert.NoError(t, tStore.Close())
	_, err := tStore.Account(signer.Address())
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
		acc, signer := account.GenerateTestAccount(int32(i))
		tStore.UpdateAccount(signer.Address(), acc)
		assert.NoError(t, tStore.WriteBatch())
		accs1 = append(accs1, acc.Hash())
	}

	stopped := false
	tStore.IterateAccounts(func(addr crypto.Address, acc *account.Account) bool {
		if acc.Hash().EqualsTo(accs1[0]) {
			stopped = true
		}
		return stopped
	})
	assert.True(t, stopped)

	accs2 := []hash.Hash{}
	tStore.IterateAccounts(func(addr crypto.Address, acc *account.Account) bool {
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

func TestRecentBlockByStamp(t *testing.T) {
	setup(t)

	hash1 := tStore.BlockHash(1)

	h, b := tStore.RecentBlockByStamp(hash.UndefHash.Stamp())
	assert.Zero(t, h)
	assert.Nil(t, b)

	h, b = tStore.RecentBlockByStamp(hash1.Stamp())
	assert.Equal(t, h, uint32(1))
	assert.Equal(t, b.Hash(), hash1)

	// Saving more blocks, blocks 11 to 22
	SaveTestBlocks(t, 12)
	hash2 := tStore.BlockHash(2)
	hash14 := tStore.BlockHash(14)
	hash22 := tStore.BlockHash(22)

	// First block should remove from the list
	h, b = tStore.RecentBlockByStamp(hash1.Stamp())
	assert.Zero(t, h)
	assert.Nil(t, b)

	h, b = tStore.RecentBlockByStamp(hash2.Stamp())
	assert.Equal(t, h, uint32(2))
	assert.Equal(t, b.Hash(), hash2)

	h, b = tStore.RecentBlockByStamp(hash14.Stamp())
	assert.Equal(t, h, uint32(14))
	assert.Equal(t, b.Hash(), hash14)

	h, b = tStore.RecentBlockByStamp(hash22.Stamp())
	assert.Equal(t, h, uint32(22))
	assert.Equal(t, b.Hash(), hash22)

	// Reopen the store
	tStore.Close()
	s, _ := NewStore(tStore.config, 21)
	tStore = s.(*store)

	h, b = tStore.RecentBlockByStamp(hash2.Stamp())
	assert.Equal(t, h, uint32(2))
	assert.Equal(t, b.Hash(), hash2)

	// Saving one more blocks, block 23
	SaveTestBlocks(t, 1)

	// Second block should remove from the list
	h, b = tStore.RecentBlockByStamp(hash2.Stamp())
	assert.Zero(t, h)
	assert.Nil(t, b)

	// Genesis block
	h, b = tStore.RecentBlockByStamp(hash.UndefHash.Stamp())
	assert.Zero(t, h)
	assert.Nil(t, b)
}
