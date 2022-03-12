package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

func TestValidatorCounter(t *testing.T) {
	setup(t)
	val, _ := validator.GenerateTestValidator(util.RandInt(1000))

	t.Run("Update count after adding new validator", func(t *testing.T) {
		assert.Equal(t, tStore.TotalValidators(), 0)
		tStore.UpdateValidator(val)
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalValidators(), 1)
	})

	t.Run("Update validator, should not increase counter", func(t *testing.T) {
		val.AddToStake(1)

		tStore.UpdateValidator(val)
		assert.NoError(t, tStore.WriteBatch())
		assert.Equal(t, tStore.TotalValidators(), 1)
	})

	t.Run("Get validator", func(t *testing.T) {
		assert.True(t, tStore.HasValidator(val.Address()))
		val2, err := tStore.Validator(val.Address())
		assert.NoError(t, err)
		assert.Equal(t, val2.Hash(), val.Hash())
	})
}

func TestValidatorBatchSaving(t *testing.T) {
	setup(t)

	t.Run("Add 100 validators", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			val, _ := validator.GenerateTestValidator(util.RandInt(1000))
			tStore.UpdateValidator(val)
			assert.NoError(t, tStore.WriteBatch())
		}

		assert.Equal(t, tStore.TotalValidators(), 100)
	})
	t.Run("Close and load db", func(t *testing.T) {
		tStore.Close()
		store, _ := NewStore(tStore.config, 21)
		assert.Equal(t, store.TotalValidators(), 100)
	})
}

func TestValidatorByNumber(t *testing.T) {
	setup(t)

	t.Run("Add some validators", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			val, _ := validator.GenerateTestValidator(i)
			tStore.UpdateValidator(val)
			assert.NoError(t, tStore.WriteBatch())
		}

		v, err := tStore.ValidatorByNumber(5)
		assert.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, v.Number(), 5)

		v, err = tStore.ValidatorByNumber(11)
		assert.Error(t, err)
		assert.Nil(t, v)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		tStore.Close()
		store, _ := NewStore(tStore.config, 21)

		v, err := store.ValidatorByNumber(5)
		assert.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, v.Number(), 5)

		v, err = tStore.ValidatorByNumber(11)
		assert.Error(t, err)
		assert.Nil(t, v)
	})
}

func TestUpdateValidator(t *testing.T) {
	setup(t)

	val1, _ := validator.GenerateTestValidator(0)
	tStore.UpdateValidator(val1)
	assert.NoError(t, tStore.WriteBatch())

	val2, _ := tStore.ValidatorByNumber(val1.Number())
	assert.Equal(t, val1.Hash(), val2.Hash())

	val3, _ := tStore.Validator(val1.Address())
	val3.AddToStake(10000)
	tStore.UpdateValidator(val3)
	assert.NoError(t, tStore.WriteBatch())
	val4, _ := tStore.ValidatorByNumber(val1.Number())
	assert.Equal(t, val4.Hash(), val3.Hash())
}
