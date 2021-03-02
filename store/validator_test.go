package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

func TestValidatorCounter(t *testing.T) {
	store, _ := newValidatorStore(util.TempDirPath())
	val, _ := validator.GenerateTestValidator(util.RandInt(1000))

	t.Run("Update count after adding new validator", func(t *testing.T) {
		assert.Equal(t, store.total, 0)

		assert.NoError(t, store.updateValidator(val))
		assert.Equal(t, store.total, 1)
	})

	t.Run("Update validator, should not increase counter", func(t *testing.T) {
		val.AddToStake(1)

		assert.NoError(t, store.updateValidator(val))
		assert.Equal(t, store.total, 1)
	})
}

func TestValidatorBatchSaving(t *testing.T) {
	path := util.TempDirPath()
	store, _ := newValidatorStore(path)

	t.Run("Add 100 validators", func(t *testing.T) {

		for i := 0; i < 100; i++ {
			val, _ := validator.GenerateTestValidator(util.RandInt(1000))
			assert.NoError(t, store.updateValidator(val))
		}

		assert.Equal(t, store.total, 100)
	})
	t.Run("Close and load db", func(t *testing.T) {
		store.close()
		store, _ = newValidatorStore(path)

		assert.Equal(t, store.total, 100)
	})
}

func TestValidatorByNumber(t *testing.T) {
	path := util.TempDirPath()
	store, _ := newValidatorStore(path)

	t.Run("Add some validators", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			val, _ := validator.GenerateTestValidator(i)
			assert.NoError(t, store.updateValidator(val))
		}

		v, err := store.validatorByNumber(5)
		assert.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, v.Number(), 5)

		v, err = store.validatorByNumber(11)
		assert.Error(t, err)
		assert.Nil(t, v)
	})

	t.Run("Reopen the store", func(t *testing.T) {
		store.close()
		store, _ := newValidatorStore(path)

		v, err := store.validatorByNumber(5)
		assert.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, v.Number(), 5)

		v, err = store.validatorByNumber(11)
		assert.Error(t, err)
		assert.Nil(t, v)
	})
}

func TestUpdateValidator(t *testing.T) {
	path := util.TempDirPath()
	store, _ := newValidatorStore(path)

	val1, _ := validator.GenerateTestValidator(0)
	assert.NoError(t, store.updateValidator(val1))

	val2, _ := store.validatorByNumber(val1.Number())
	assert.Equal(t, val1.Hash(), val2.Hash())

	val3, _ := store.validator(val1.Address())
	val3.AddToStake(10000)
	assert.NoError(t, store.updateValidator(val3))

	val4, _ := store.validatorByNumber(val1.Number())
	assert.Equal(t, val4.Hash(), val3.Hash())
}
