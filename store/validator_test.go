package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

func TestValidatorCounter(t *testing.T) {
	store, _ := NewStore(TestConfig())
	val, _ := validator.GenerateTestValidator(util.RandInt(1000))

	t.Run("Update count after adding new validator", func(t *testing.T) {
		assert.Equal(t, store.TotalValidators(), 0)
		store.UpdateValidator(val)
		assert.NoError(t, store.WriteBatch())
		assert.Equal(t, store.TotalValidators(), 1)
	})

	t.Run("Update validator, should not increase counter", func(t *testing.T) {
		val.AddToStake(1)

		store.UpdateValidator(val)
		assert.NoError(t, store.WriteBatch())
		assert.Equal(t, store.TotalValidators(), 1)
	})
}

func TestValidatorBatchSaving(t *testing.T) {
	conf := TestConfig()
	store, _ := NewStore(conf)

	t.Run("Add 100 validators", func(t *testing.T) {

		for i := 0; i < 100; i++ {
			val, _ := validator.GenerateTestValidator(util.RandInt(1000))
			store.UpdateValidator(val)
			assert.NoError(t, store.WriteBatch())
		}

		assert.Equal(t, store.TotalValidators(), 100)
	})
	t.Run("Close and load db", func(t *testing.T) {
		store.Close()
		store, _ = NewStore(conf)

		assert.Equal(t, store.TotalValidators(), 100)
	})
}

func TestValidatorByNumber(t *testing.T) {
	conf := TestConfig()
	store, _ := NewStore(conf)

	t.Run("Add some validators", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			val, _ := validator.GenerateTestValidator(i)
			store.UpdateValidator(val)
			assert.NoError(t, store.WriteBatch())
		}

		v, err := store.ValidatorByNumber(5)
		assert.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, v.Number(), 5)

		v, err = store.ValidatorByNumber(11)
		assert.Error(t, err)
		assert.Nil(t, v)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		store.Close()
		store, _ := NewStore(conf)

		v, err := store.ValidatorByNumber(5)
		assert.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, v.Number(), 5)

		v, err = store.ValidatorByNumber(11)
		assert.Error(t, err)
		assert.Nil(t, v)
	})
}

func TestUpdateValidator(t *testing.T) {
	store, _ := NewStore(TestConfig())

	val1, _ := validator.GenerateTestValidator(0)
	store.UpdateValidator(val1)
	assert.NoError(t, store.WriteBatch())

	val2, _ := store.ValidatorByNumber(val1.Number())
	assert.Equal(t, val1.Hash(), val2.Hash())

	val3, _ := store.Validator(val1.Address())
	val3.AddToStake(10000)
	store.UpdateValidator(val3)
	assert.NoError(t, store.WriteBatch())
	val4, _ := store.ValidatorByNumber(val1.Number())
	assert.Equal(t, val4.Hash(), val3.Hash())
}
