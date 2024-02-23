package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	someDB, err := NewDB(":memory:")

	assert.Nil(t, err)
	assert.NotNil(t, someDB)
}

func TestInsert(t *testing.T) {
	t.Run("could not insert into address table", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")

		addr := &Address{
			Address:   "some-address",
			PublicKey: "some-public-key",
			Label:     "some-label",
			Path:      "some-path",
		}
		_, err := someDB.InsertIntoAddress(addr)
		assert.EqualError(t, ErrCouldNotInsertRecordIntoTable, err.Error())
	})

	t.Run("insert into address table", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
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
		someDB, _ := NewDB(":memory:")

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
		assert.EqualError(t, ErrCouldNotInsertRecordIntoTable, err.Error())
	})

	t.Run("insert into tranasction table", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
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
		someDB, _ := NewDB(":memory:")

		key, value := "key", "value"
		_, err := someDB.InsertIntoPair(key, value)

		assert.EqualError(t, ErrCouldNotInsertRecordIntoTable, err.Error())
	})

	t.Run("insert into pair table", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		key, value := "key", "value"
		actual, err := someDB.InsertIntoPair(key, value)

		assert.Nil(t, err)
		assert.Equal(t, key, actual.Key)
		assert.Equal(t, value, actual.Value)
	})
}

func TestGetById(t *testing.T) {
	t.Run("could not get address by id", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		addr := &Address{
			Address:   "some-address",
			PublicKey: "some-public-key",
			Label:     "some-label",
			Path:      "some-path",
		}
		_, _ = someDB.InsertIntoAddress(addr)

		actual, err := someDB.GetAddressByID(5)
		assert.Nil(t, actual)
		assert.EqualError(t, ErrCouldNotFindRecord, err.Error())
	})

	t.Run("get address by id", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		addr := &Address{
			Address:   "some-address",
			PublicKey: "some-public-key",
			Label:     "some-label",
			Path:      "some-path",
		}
		expected, _ := someDB.InsertIntoAddress(addr)

		actual, err := someDB.GetAddressByID(expected.ID)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("could not get transaction by id", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
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
		_, _ = someDB.InsertIntoTransaction(tr)

		actual, err := someDB.GetTransactionByID(10)

		assert.Nil(t, actual)
		assert.EqualError(t, ErrCouldNotFindRecord, err.Error())
	})

	t.Run("get transaction by id", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
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
		expected, _ := someDB.InsertIntoTransaction(tr)

		actual, err := someDB.GetTransactionByID(expected.ID)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("could not get pair by key", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		key, value := "key", "value"
		_, _ = someDB.InsertIntoPair(key, value)

		actual, err := someDB.GetPairByKey("some-thing-wrong")

		assert.Nil(t, actual)
		assert.EqualError(t, ErrCouldNotFindRecord, err.Error())
	})

	t.Run("get pair by key", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		key, value := "key", "value"
		expected, _ := someDB.InsertIntoPair(key, value)

		actual, err := someDB.GetPairByKey(expected.Key)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestAddress(t *testing.T) {
	t.Run("Could not get address by address", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		addr := &Address{
			Address:   "some-pactus-addr",
			PublicKey: "some-public-key",
			Label:     "some-label",
			Path:      "some-path",
		}
		expected, _ := someDB.InsertIntoAddress(addr)

		expected.Address = "some-other-pactus-addr"
		actual, err := someDB.GetAddressByAddress(expected.Address)

		assert.Nil(t, actual)
		assert.EqualError(t, ErrCouldNotFindRecord, err.Error())
	})

	t.Run("Get address by address", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		addr := &Address{
			Address:   "some-pactus-addr",
			PublicKey: "some-public-key",
			Label:     "some-label",
			Path:      "some-path",
		}
		expected, _ := someDB.InsertIntoAddress(addr)

		actual, err := someDB.GetAddressByAddress(expected.Address)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("could not get address by path", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		addr := &Address{
			Address:   "some-pactus-addr",
			PublicKey: "some-public-key",
			Label:     "some-label",
			Path:      "some-path",
		}
		expected, _ := someDB.InsertIntoAddress(addr)

		expected.Path = "some-other-path"
		actual, err := someDB.GetAddressByPath(expected.Path)

		assert.Nil(t, actual)
		assert.EqualError(t, ErrCouldNotFindRecord, err.Error())
	})

	t.Run("Get address by path", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		addr := &Address{
			Address:   "some-pactus-addr",
			PublicKey: "some-public-key",
			Label:     "some-label",
			Path:      "some-path",
		}
		expected, _ := someDB.InsertIntoAddress(addr)

		actual, err := someDB.GetAddressByPath(expected.Path)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("update label of address", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		addr := &Address{
			Address:    "some-pactus-addr",
			PublicKey:  "some-public-key",
			Label:      "some-label",
			Path:       "some-path",
			IsImported: true,
		}
		addr, _ = someDB.InsertIntoAddress(addr)

		addr.Label = "some-other-lable"
		_, _ = someDB.UpdateAddressLabel(addr)

		actual, err := someDB.GetAddressByAddress(addr.Address)

		assert.Nil(t, err)
		assert.Equal(t, addr.Label, actual.Label)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("get all addresses", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		addr := &Address{
			Address:    "some-address",
			PublicKey:  "some-public-key",
			Label:      "some-label",
			Path:       "some-path",
			IsImported: true,
		}
		someInsertOne, _ := someDB.InsertIntoAddress(addr)
		someInsertTwo, _ := someDB.InsertIntoAddress(addr)
		someInsertThree, _ := someDB.InsertIntoAddress(addr)

		expected := make([]Address, 0, 3)
		expected = append(expected, *someInsertThree, *someInsertTwo, *someInsertOne)

		acutal, err := someDB.GetAllAddresses()

		assert.Nil(t, err)
		assert.Equal(t, expected, acutal)
	})

	t.Run("get all addresses with total records", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		addr := &Address{
			Address:    "some-address",
			PublicKey:  "some-public-key",
			Label:      "some-label",
			Path:       "some-path",
			IsImported: true,
		}
		someInsertOne, _ := someDB.InsertIntoAddress(addr)
		someInsertTwo, _ := someDB.InsertIntoAddress(addr)
		someInsertThree, _ := someDB.InsertIntoAddress(addr)

		expected := make([]Address, 0, 3)
		expected = append(expected, *someInsertThree, *someInsertTwo, *someInsertOne)

		acutal, totalRecords, err := someDB.GetAllAddressesWithTotalRecords(1, 3)

		assert.Nil(t, err)
		assert.Equal(t, int64(3), totalRecords)
		assert.Equal(t, expected, acutal)
	})

	t.Run("get all transactions", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
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
		someInsertOne, _ := someDB.InsertIntoTransaction(tr)
		someInsertTwo, _ := someDB.InsertIntoTransaction(tr)
		someInsertThree, _ := someDB.InsertIntoTransaction(tr)

		expected := make([]Transaction, 0, 3)
		expected = append(expected, *someInsertThree, *someInsertTwo, *someInsertOne)

		acutal, err := someDB.GetAllTransactions()

		assert.Nil(t, err)
		assert.Equal(t, expected, acutal)
	})

	t.Run("get all transactions with total records", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
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
		someInsertOne, _ := someDB.InsertIntoTransaction(tr)
		someInsertTwo, _ := someDB.InsertIntoTransaction(tr)
		someInsertThree, _ := someDB.InsertIntoTransaction(tr)

		expected := make([]Transaction, 0, 3)
		expected = append(expected, *someInsertThree, *someInsertTwo, *someInsertOne)

		acutal, totalRecords, err := someDB.GetAllTransactionsWithTotalRecords(1, 3)

		assert.Nil(t, err)
		assert.Equal(t, int64(3), totalRecords)
		assert.Equal(t, expected, acutal)
	})
}

func TestTotalRecords(t *testing.T) {
	t.Run("could not find total records", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		addr := &Address{
			Address:    "some-address",
			PublicKey:  "some-public-key",
			Label:      "some-label",
			Path:       "some-path",
			IsImported: true,
		}
		_, _ = someDB.InsertIntoAddress(addr)
		_, _ = someDB.InsertIntoAddress(addr)
		_, _ = someDB.InsertIntoAddress(addr)

		totalRecords, err := someDB.GetTotalRecords("some-table")

		assert.Equal(t, int64(0), totalRecords)
		assert.EqualError(t, ErrCouldNotFindTotalRecords, err.Error())
	})

	t.Run("ok", func(t *testing.T) {
		someDB, _ := NewDB(":memory:")
		_ = someDB.CreateTables()

		addr := &Address{
			Address:    "some-address",
			PublicKey:  "some-public-key",
			Label:      "some-label",
			Path:       "some-path",
			IsImported: true,
		}
		_, _ = someDB.InsertIntoAddress(addr)
		_, _ = someDB.InsertIntoAddress(addr)
		_, _ = someDB.InsertIntoAddress(addr)

		totalRecords, err := someDB.GetTotalRecords("addresses")

		assert.Nil(t, err)
		assert.Equal(t, int64(3), totalRecords)
	})
}
