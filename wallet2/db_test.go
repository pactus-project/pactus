package wallet2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	someDB, err := newDB(":memory:")

	assert.Nil(t, err)
	assert.NotNil(t, someDB)
}

func TestInsert(t *testing.T) {
	t.Run("could not insert into address table", func(t *testing.T) {
		someDB, _ := newDB(":memory:")

		addr := &Address{
			Address:   "some-address",
			PublicKey: "some-public-key",
			Label:     "some-label",
			Path:      "some-path",
		}
		_, err := someDB.InsertIntoAddress(addr)
		assert.EqualError(t, ErrCouldNotInsertIntoTable, err.Error())
	})

	t.Run("insert into address table", func(t *testing.T) {
		someDB, _ := newDB(":memory:")
		_ = someDB.CreateTables()

		addr := &Address{
			Address:   "some-address",
			PublicKey: "some-public-key",
			Label:     "some-label",
			Path:      "some-path",
		}
		actual, err := someDB.InsertIntoAddress(addr)

		assert.Nil(t, err)
		assert.Equal(t, 1, actual.ID)
		assert.Equal(t, addr.Address, actual.Address)
	})

	t.Run("could not insert into tranasction table", func(t *testing.T) {
		someDB, _ := newDB(":memory:")

		tr := &Transaction{
			TxID:        "some-txid",
			BlockHeight: 4,
			BlockTime:   5,
			PayloadType: "something",
			Data:        "some-data",
			Description: "some-description",
			Amount:      50,
			Status:      1,
		}
		_, err := someDB.InsertIntoTransaction(tr)
		assert.EqualError(t, ErrCouldNotInsertIntoTable, err.Error())
	})

	t.Run("insert into tranasction table", func(t *testing.T) {
		someDB, _ := newDB(":memory:")
		_ = someDB.CreateTables()

		tr := &Transaction{
			TxID:        "some-txid",
			BlockHeight: 4,
			BlockTime:   5,
			PayloadType: "something",
			Data:        "some-data",
			Description: "some-description",
			Amount:      50,
			Status:      1,
		}
		actual, err := someDB.InsertIntoTransaction(tr)

		assert.Nil(t, err)
		assert.Equal(t, 1, actual.ID)
		assert.Equal(t, tr.BlockHeight, actual.BlockHeight)
	})

	t.Run("could not insert into pair table", func(t *testing.T) {
		someDB, _ := newDB(":memory:")

		key, value := "key", "value"
		_, err := someDB.InsertIntoPair(key, value)

		assert.EqualError(t, ErrCouldNotInsertIntoTable, err.Error())
	})

	t.Run("insert into pair table", func(t *testing.T) {
		someDB, _ := newDB(":memory:")
		_ = someDB.CreateTables()

		key, value := "key", "value"
		actual, err := someDB.InsertIntoPair(key, value)

		assert.Nil(t, err)
		assert.Equal(t, key, actual.Key)
		assert.Equal(t, value, actual.Value)
	})
}
