package store

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatorCounter(t *testing.T) {
	setup(t)

	num := util.RandInt32(1000)
	val, _ := validator.GenerateTestValidator(num)

	t.Run("Add new validator, should increase the total validators number", func(t *testing.T) {
		assert.Zero(t, tStore.TotalValidators())

		tStore.UpdateValidator(val)
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalValidators(), int32(1))
	})

	t.Run("Update validator, should not increase the total validators number", func(t *testing.T) {
		val.AddToStake(1)
		tStore.UpdateValidator(val)

		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalValidators(), int32(1))
	})

	t.Run("Get validator", func(t *testing.T) {
		val1, err := tStore.Validator(val.Address())
		assert.NoError(t, err)

		val2, err := tStore.ValidatorByNumber(num)
		assert.NoError(t, err)

		assert.Equal(t, val1.Hash(), val2.Hash())
		assert.Equal(t, tStore.TotalValidators(), int32(1))
		assert.True(t, tStore.HasValidator(val.Address()))
	})
}

func TestValidatorBatchSaving(t *testing.T) {
	setup(t)

	total := util.RandInt32(100) + 1
	t.Run("Add some validators", func(t *testing.T) {
		for i := int32(0); i < total; i++ {
			val, _ := validator.GenerateTestValidator(i)
			tStore.UpdateValidator(val)
		}
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalValidators(), total)
	})

	t.Run("Close and load db", func(t *testing.T) {
		tStore.Close()
		store, _ := NewStore(tStore.config, 21)
		assert.Equal(t, store.TotalValidators(), total)
	})
}

func TestValidatorByNumber(t *testing.T) {
	setup(t)

	total := util.RandInt32(100) + 1
	t.Run("Add some validators", func(t *testing.T) {
		for i := int32(0); i < total; i++ {
			val, _ := validator.GenerateTestValidator(i)
			tStore.UpdateValidator(val)
		}
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalValidators(), total)
	})

	t.Run("Get a random Validator", func(t *testing.T) {
		num := util.RandInt32(total)
		val, err := tStore.ValidatorByNumber(num)
		assert.NoError(t, err)
		require.NotNil(t, val)
		assert.Equal(t, val.Number(), num)
	})

	t.Run("negative number", func(t *testing.T) {
		val, err := tStore.ValidatorByNumber(-1)
		assert.Error(t, err)
		assert.Nil(t, val)
	})

	t.Run("Non existing validator", func(t *testing.T) {
		val, err := tStore.ValidatorByNumber(total + 1)
		assert.Error(t, err)
		assert.Nil(t, val)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		tStore.Close()
		store, _ := NewStore(tStore.config, 21)

		num := util.RandInt32(total)
		val, err := store.ValidatorByNumber(num)
		assert.NoError(t, err)
		require.NotNil(t, val)
		assert.Equal(t, val.Number(), num)

		val, err = tStore.ValidatorByNumber(total + 1)
		assert.Error(t, err)
		assert.Nil(t, val)
	})
}

func TestValidatorByAddress(t *testing.T) {
	setup(t)

	total := util.RandInt32(100) + 1
	t.Run("Add some validators", func(t *testing.T) {
		for i := int32(0); i < total; i++ {
			val, _ := validator.GenerateTestValidator(i)
			tStore.UpdateValidator(val)
		}
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalValidators(), total)
	})

	t.Run("Get random validator", func(t *testing.T) {
		num := util.RandInt32(total)
		val0, _ := tStore.ValidatorByNumber(num)
		val, err := tStore.Validator(val0.Address())
		assert.NoError(t, err)
		require.NotNil(t, val)
		assert.Equal(t, val.Number(), num)
	})

	t.Run("Unknown address", func(t *testing.T) {
		val, err := tStore.Validator(crypto.GenerateTestAddress())
		assert.Error(t, err)
		assert.Nil(t, val)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		tStore.Close()
		store, _ := NewStore(tStore.config, 21)

		num := util.RandInt32(total)
		val0, _ := store.ValidatorByNumber(num)
		val, err := store.Validator(val0.Address())
		assert.NoError(t, err)
		require.NotNil(t, val)
		assert.Equal(t, val.Number(), num)
	})
}

func TestUpdateValidator(t *testing.T) {
	setup(t)

	num := util.RandInt32(1000)
	val1, _ := validator.GenerateTestValidator(num)
	tStore.UpdateValidator(val1)

	val2, _ := tStore.ValidatorByNumber(num)
	assert.Equal(t, val1.Hash(), val2.Hash())

	assert.NoError(t, tStore.WriteBatch())
	val3, _ := tStore.ValidatorByNumber(num)
	assert.Equal(t, val1.Hash(), val3.Hash())

	val4, _ := tStore.Validator(val1.Address())
	val4.AddToStake(1)
	tStore.UpdateValidator(val4)

	val5, _ := tStore.ValidatorByNumber(num)
	assert.Equal(t, val4.Hash(), val5.Hash())
}
