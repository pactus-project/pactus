package wallet2

import (
	"context"
	"testing"

	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet2/db"
	"github.com/stretchr/testify/assert"
)

func TestNewHistory(t *testing.T) {
	ctx := context.Background()
	someDB, _ := db.NewDB(ctx, ":memory:")
	h := newHistory(someDB)

	assert.NotNil(t, h)
}

func TestAddTransaction(t *testing.T) {
	t.Run("could not add transaction", func(t *testing.T) {
		ts := testsuite.NewTestSuite(t)
		ctx := context.Background()
		someDB, _ := db.NewDB(ctx, ":memory:")

		h := newHistory(someDB)
		someHash := ts.RandHash()

		err := h.addTransaction("some-addr", 3, someHash, 1, 1, payload.TypeTransfer.String(), []byte{1}, db.Confirmed)

		assert.Error(t, db.ErrCouldNotInsertRecordIntoTable, err.Error())
	})

	t.Run("add transaction ok", func(t *testing.T) {
		ts := testsuite.NewTestSuite(t)
		ctx := context.Background()
		someDB, _ := db.NewDB(ctx, ":memory:")
		_ = someDB.CreateTables()

		h := newHistory(someDB)
		someHash := ts.RandHash()

		err := h.addTransaction("some-addr", 3, someHash, 1, 1, payload.TypeTransfer.String(), []byte{1}, db.Confirmed)

		assert.NoError(t, err)
	})
}

func TestHasTransaction(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	ctx := context.Background()
	someDB, _ := db.NewDB(ctx, ":memory:")
	_ = someDB.CreateTables()

	h := newHistory(someDB)

	someHash := ts.RandHash()

	_ = h.addTransaction("some-addr", 3, someHash, 1, 1, payload.TypeTransfer.String(), []byte{1}, db.Confirmed)

	ok := h.hasTransaction(someHash.String())
	assert.True(t, ok)
}

func TestGetAddrHistory(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	ctx := context.Background()
	someDB, _ := db.NewDB(ctx, ":memory:")
	_ = someDB.CreateTables()

	h := newHistory(someDB)

	rndHash1 := ts.RandHash()
	rndHash2 := ts.RandHash()
	rndHash3 := ts.RandHash()

	_ = h.addTransaction("addr-1", 3, rndHash1, 1, 1, payload.TypeTransfer.String(), []byte{1}, db.Confirmed)
	_ = h.addTransaction("addr-1", 3, rndHash2, 1, 1, payload.TypeTransfer.String(), []byte{1}, db.Confirmed)
	_ = h.addTransaction("addr-2", 3, rndHash3, 1, 1, payload.TypeTransfer.String(), []byte{1}, db.Confirmed)

	history, err := h.getAddrHistory("addr-1")
	assert.NoError(t, err)
	assert.Equal(t, rndHash2.String(), history[0].TxID)
	assert.Equal(t, rndHash1.String(), history[1].TxID)
}
