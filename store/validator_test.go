package store

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatorCounter(t *testing.T) {
	td := setup(t, nil)

	num := td.RandInt32(1000)
	val, _ := td.GenerateTestValidator(num)

	t.Run("Add new validator, should increase the total validators number", func(t *testing.T) {
		assert.Zero(t, td.store.TotalValidators())

		td.store.UpdateValidator(val)
		assert.NoError(t, td.store.WriteBatch())
		assert.Equal(t, int32(1), td.store.TotalValidators())
	})

	t.Run("Update validator, should not increase the total validators number", func(t *testing.T) {
		val.AddToStake(1)
		td.store.UpdateValidator(val)

		assert.NoError(t, td.store.WriteBatch())
		assert.Equal(t, int32(1), td.store.TotalValidators())
	})

	t.Run("Get validator", func(t *testing.T) {
		val1, err := td.store.Validator(val.Address())
		assert.NoError(t, err)

		val2, err := td.store.ValidatorByNumber(num)
		assert.NoError(t, err)

		assert.Equal(t, val1.Hash(), val2.Hash())
		assert.Equal(t, int32(1), td.store.TotalValidators())
		assert.True(t, td.store.HasValidator(val.Address()))
	})
}

func TestValidatorBatchSaving(t *testing.T) {
	td := setup(t, nil)

	total := td.RandInt32NonZero(100)
	t.Run("Add some validators", func(t *testing.T) {
		for i := int32(0); i < total; i++ {
			val, _ := td.GenerateTestValidator(i)
			td.store.UpdateValidator(val)
		}
		assert.NoError(t, td.store.WriteBatch())
		assert.Equal(t, total, td.store.TotalValidators())
	})

	t.Run("Close and load db", func(t *testing.T) {
		td.store.Close()
		store, _ := NewStore(td.store.config)
		assert.Equal(t, total, store.TotalValidators())
	})
}

func TestValidatorAddresses(t *testing.T) {
	td := setup(t, nil)

	total := td.RandInt32NonZero(100)
	addrs1 := make([]crypto.Address, 0, total)

	for i := int32(0); i < total; i++ {
		val, _ := td.GenerateTestValidator(i)
		td.store.UpdateValidator(val)
		addrs1 = append(addrs1, val.Address())
	}

	addrs2 := td.store.ValidatorAddresses()
	assert.ElementsMatch(t, addrs1, addrs2)
}

func TestValidatorByNumber(t *testing.T) {
	td := setup(t, nil)

	total := td.RandInt32NonZero(100)
	t.Run("Add some validators", func(t *testing.T) {
		for i := int32(0); i < total; i++ {
			val, _ := td.GenerateTestValidator(i)
			td.store.UpdateValidator(val)
		}
		assert.NoError(t, td.store.WriteBatch())
		assert.Equal(t, total, td.store.TotalValidators())
	})

	t.Run("Get a random Validator", func(t *testing.T) {
		num := td.RandInt32(total)
		val, err := td.store.ValidatorByNumber(num)
		assert.NoError(t, err)
		require.NotNil(t, val)
		assert.Equal(t, num, val.Number())
	})

	t.Run("Negative number", func(t *testing.T) {
		val, err := td.store.ValidatorByNumber(-1)
		assert.Error(t, err)
		assert.Nil(t, val)
	})

	t.Run("Non existing validator", func(t *testing.T) {
		addr := td.RandValAddress()
		val, err := td.store.Validator(addr)
		has := td.store.HasValidator(addr)

		assert.False(t, has)
		assert.Error(t, err)
		assert.Nil(t, val)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		td.store.Close()
		store, _ := NewStore(td.store.config)

		num := td.RandInt32(total)
		val, err := store.ValidatorByNumber(num)
		assert.NoError(t, err)
		require.NotNil(t, val)
		assert.Equal(t, num, val.Number())

		val, err = td.store.ValidatorByNumber(total + 1)
		assert.Error(t, err)
		assert.Nil(t, val)
	})
}

func TestValidatorByAddress(t *testing.T) {
	td := setup(t, nil)

	total := td.RandInt32NonZero(100)
	t.Run("Add some validators", func(t *testing.T) {
		for i := int32(0); i < total; i++ {
			val, _ := td.GenerateTestValidator(i)
			td.store.UpdateValidator(val)
		}
		assert.NoError(t, td.store.WriteBatch())
		assert.Equal(t, total, td.store.TotalValidators())
	})

	t.Run("Get random validator", func(t *testing.T) {
		num := td.RandInt32(total)
		val0, _ := td.store.ValidatorByNumber(num)
		val, err := td.store.Validator(val0.Address())
		assert.NoError(t, err)
		require.NotNil(t, val)
		assert.Equal(t, num, val.Number())
	})

	t.Run("Unknown address", func(t *testing.T) {
		val, err := td.store.Validator(td.RandAccAddress())
		assert.Error(t, err)
		assert.Nil(t, val)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		td.store.Close()
		store, _ := NewStore(td.store.config)

		num := td.RandInt32(total)
		val0, _ := store.ValidatorByNumber(num)
		val, err := store.Validator(val0.Address())
		assert.NoError(t, err)
		require.NotNil(t, val)
		assert.Equal(t, num, val.Number())
	})
}

func TestIterateValidators(t *testing.T) {
	td := setup(t, nil)

	total := td.RandInt32NonZero(100)
	hashes1 := []hash.Hash{}
	for i := int32(0); i < total; i++ {
		val, _ := td.GenerateTestValidator(i)
		td.store.UpdateValidator(val)
		hashes1 = append(hashes1, val.Hash())
	}
	assert.NoError(t, td.store.WriteBatch())

	hashes2 := []hash.Hash{}
	td.store.IterateValidators(func(val *validator.Validator) bool {
		hashes2 = append(hashes2, val.Hash())

		return false
	})
	assert.ElementsMatch(t, hashes1, hashes2)

	stopped := false
	td.store.IterateValidators(func(val *validator.Validator) bool {
		if val.Hash() == hashes1[0] {
			stopped = true
		}

		return stopped
	})
	assert.True(t, stopped)
}

func TestValidatorDeepCopy(t *testing.T) {
	td := setup(t, nil)

	num := td.RandInt32NonZero(1000)
	val1, _ := td.GenerateTestValidator(num)
	td.store.UpdateValidator(val1)

	val2, _ := td.store.ValidatorByNumber(num)
	val2.AddToStake(1)
	assert.NotEqual(t, td.store.validatorStore.numberMap[num].Hash(), val2.Hash())

	val3, _ := td.store.Validator(val1.Address())
	val3.AddToStake(1)
	assert.NotEqual(t, td.store.validatorStore.numberMap[num].Hash(), val3.Hash())
}
