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

	t.Run("Update count after adding new validator", func(t *testing.T) {
		assert.Zero(t, tStore.TotalValidators())

		tStore.UpdateValidator(val)
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalValidators(), int32(1))
	})

	t.Run("Update validator, should not increase counter", func(t *testing.T) {
		val.AddToStake(1)
		tStore.UpdateValidator(val)

		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalValidators(), int32(1))
	})

	t.Run("Get validator", func(t *testing.T) {
		assert.True(t, tStore.HasValidator(val.Address()))
		val1, err := tStore.Validator(val.Address())
		assert.NoError(t, err)

		val2, err := tStore.ValidatorByNumber(val.Number())
		assert.NoError(t, err)

		assert.Equal(t, val1.Hash(), val2.Hash())
	})
}

func TestValidatorBatchSaving(t *testing.T) {
	setup(t)

	t.Run("Add 100 validators", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			val, _ := validator.GenerateTestValidator(int32(i))
			tStore.UpdateValidator(val)
		}
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalValidators(), int32(100))
	})
	t.Run("Close and load db", func(t *testing.T) {
		tStore.Close()
		store, _ := NewStore(tStore.config, 21)
		assert.Equal(t, store.TotalValidators(), int32(100))
	})
}

func TestValidatorMap(t *testing.T) {
	setup(t)

	t.Run("Add 100 validators", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			val, _ := validator.GenerateTestValidator(int32(i))
			tStore.UpdateValidator(val)
		}
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalValidators(), int32(100))
	})

	t.Run("Get random validator", func(t *testing.T) {
		valNum := util.RandInt32(100)
		val, err := tStore.ValidatorByNumber(valNum)
		assert.NoError(t, err)
		require.NotNil(t, val)
		assert.Equal(t, val.Number(), valNum)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		tStore.Close()
		store, _ := NewStore(tStore.config, 21)

		valNum := util.RandInt32(100)
		val1, err := store.ValidatorByNumber(valNum)
		assert.NoError(t, err)
		assert.Equal(t, val1.Number(), valNum)

		val2, err := store.Validator(val1.Address())
		assert.NoError(t, err)
		assert.Equal(t, val1.Hash(), val2.Hash())
	})

	t.Run("Unknown validators", func(t *testing.T) {
		val, err := tStore.ValidatorByNumber(123456)
		assert.Error(t, err)
		assert.Nil(t, val)

		val, err = tStore.Validator(crypto.GenerateTestAddress())
		assert.Error(t, err)
		assert.Nil(t, val)
	})
}

func TestUpdateValidator(t *testing.T) {
	setup(t)

	num := util.RandInt32(1000)
	val1, _ := validator.GenerateTestValidator(num)
	tStore.UpdateValidator(val1)
	assert.NoError(t, tStore.WriteBatch())

	val2, _ := tStore.ValidatorByNumber(num)
	assert.Equal(t, val1.Hash(), val2.Hash())

	val3, _ := tStore.Validator(val1.Address())
	val3.AddToStake(10000)
	tStore.UpdateValidator(val3)
	assert.NoError(t, tStore.WriteBatch())

	val4, _ := tStore.ValidatorByNumber(num)
	assert.Equal(t, val4.Hash(), val3.Hash())
}
