package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

func TestRetreiveValidator(t *testing.T) {
	store, _ := newValidatorStore(util.TempDirPath())

	val, _ := validator.GenerateTestValidator(util.RandInt(1000))

	t.Run("Add validator, should able to retrieve", func(t *testing.T) {
		assert.False(t, store.hasValidator(val.Address()))
		store.updateValidator(val)
		assert.True(t, store.hasValidator(val.Address()))
		val2, err := store.validator(val.Address())
		assert.NoError(t, err)
		assert.Equal(t, val.Hash(), val2.Hash())
	})

	t.Run("Update validator, should update database", func(t *testing.T) {
		val.AddToStake(1)
		store.updateValidator(val)

		val2, err := store.validator(val.Address())
		assert.NoError(t, err)
		assert.Equal(t, val.Hash(), val2.Hash())
	})
}

func TestValidatorCounter(t *testing.T) {
	store, _ := newValidatorStore(util.TempDirPath())

	val, _ := validator.GenerateTestValidator(util.RandInt(1000))

	t.Run("Update count after adding new validator", func(t *testing.T) {
		assert.Equal(t, store.total, store.countValidators())
		assert.Equal(t, store.total, 0)

		store.updateValidator(val)
		assert.Equal(t, store.total, store.countValidators())
		assert.Equal(t, store.total, 1)
	})

	t.Run("Update validator, should not increatse counter", func(t *testing.T) {
		val.AddToStake(1)
		store.updateValidator(val)

		store.updateValidator(val)
		assert.Equal(t, store.total, store.countValidators())
		assert.Equal(t, store.total, 1)
	})
}

func TestValidatorBatchSaving(t *testing.T) {
	path := util.TempDirPath()
	store, _ := newValidatorStore(path)

	t.Run("Add 100 validators", func(t *testing.T) {

		for i := 0; i < 100; i++ {
			val, _ := validator.GenerateTestValidator(util.RandInt(1000))
			store.updateValidator(val)
		}

		assert.Equal(t, store.total, store.countValidators())
		assert.Equal(t, store.total, 100)
	})
	t.Run("Close and load db", func(t *testing.T) {
		store.close()
		store, _ = newValidatorStore(path)

		assert.Equal(t, store.total, store.countValidators())
		assert.Equal(t, store.total, 100)
	})
}
