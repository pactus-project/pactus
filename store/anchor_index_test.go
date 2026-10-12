package store

import (
	"bytes"
	"math"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func anchoredAccount(t *testing.T, td *testData, number int32) (crypto.Address, *account.Account) {
	t.Helper()

	addr, acc := td.GenerateTestAccount(testsuite.AccountWithNumber(number))
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytes.Repeat([]byte{byte(number)}, 32),
		LockedDeposit: 1,
	}))

	return addr, acc
}

func TestAnchorIndex(t *testing.T) {
	td := setup(t, nil)

	addr3, acc3 := anchoredAccount(t, td, 3)
	addr1, acc1 := anchoredAccount(t, td, 1)
	addr2, acc2 := anchoredAccount(t, td, 2)
	plainAddr, plainAcc := td.GenerateTestAccount(testsuite.AccountWithNumber(0))
	td.store.UpdateAccount(addr3, acc3)
	td.store.UpdateAccount(addr1, acc1)
	td.store.UpdateAccount(addr2, acc2)
	td.store.UpdateAccount(plainAddr, plainAcc)
	require.NoError(t, td.store.WriteBatch())

	t.Run("pages are ordered by account number", func(t *testing.T) {
		addrs, total := td.store.AnchorAddresses(0, 10)
		require.Equal(t, uint32(3), total)
		require.Equal(t, []crypto.Address{addr1, addr2, addr3}, addrs)

		addrs, total = td.store.AnchorAddresses(1, 1)
		require.Equal(t, uint32(3), total)
		require.Equal(t, []crypto.Address{addr2}, addrs)

		addrs, total = td.store.AnchorAddresses(1, math.MaxUint32)
		require.Equal(t, uint32(3), total)
		require.Equal(t, []crypto.Address{addr2, addr3}, addrs)
	})

	t.Run("empty pages", func(t *testing.T) {
		addrs, total := td.store.AnchorAddresses(3, 10)
		require.Equal(t, uint32(3), total)
		require.Empty(t, addrs)

		addrs, total = td.store.AnchorAddresses(0, 0)
		require.Equal(t, uint32(3), total)
		require.Empty(t, addrs)
	})

	t.Run("rewriting an anchored account keeps one entry", func(t *testing.T) {
		require.NoError(t, acc1.SetAnchor(account.AnchorData{
			RootHash:      bytes.Repeat([]byte{0x77}, 64),
			LockedDeposit: 2,
		}))
		td.store.UpdateAccount(addr1, acc1)
		td.store.UpdateAccount(addr1, acc1)
		require.NoError(t, td.store.WriteBatch())

		addrs, total := td.store.AnchorAddresses(0, 10)
		require.Equal(t, uint32(3), total)
		require.Equal(t, []crypto.Address{addr1, addr2, addr3}, addrs)
	})

	t.Run("updates keep the index in sync", func(t *testing.T) {
		acc2.ClearAnchor()
		td.store.UpdateAccount(addr2, acc2)
		td.store.UpdateAccount(addr2, acc2)
		require.NoError(t, td.store.WriteBatch())

		addrs, total := td.store.AnchorAddresses(0, 10)
		require.Equal(t, uint32(2), total)
		require.Equal(t, []crypto.Address{addr1, addr3}, addrs)

		require.NoError(t, plainAcc.SetAnchor(account.AnchorData{
			RootHash:      bytes.Repeat([]byte{0xAA}, 32),
			LockedDeposit: 1,
		}))
		td.store.UpdateAccount(plainAddr, plainAcc)
		require.NoError(t, td.store.WriteBatch())

		addrs, total = td.store.AnchorAddresses(0, 10)
		require.Equal(t, uint32(3), total)
		require.Equal(t, []crypto.Address{plainAddr, addr1, addr3}, addrs)
	})

	t.Run("the index is rebuilt when the store reopens", func(t *testing.T) {
		td.store.Close()
		reopened, err := NewStore(td.store.config)
		require.NoError(t, err)
		td.store = reopened.(*store)

		addrs, total := td.store.AnchorAddresses(0, 10)
		require.Equal(t, uint32(3), total)
		require.Equal(t, []crypto.Address{plainAddr, addr1, addr3}, addrs)
	})
}

// The anchor index is built from the stored accounts at startup. A record
// longer than a plain account that does not decode must stop the node.
func TestNewStorePanicsOnCorruptAnchorRecord(t *testing.T) {
	td := setup(t, nil)

	addr, acc := anchoredAccount(t, td, 5)
	raw, err := acc.Bytes()
	require.NoError(t, err)
	raw[12] = 0x02 // the anchor flag must be 1
	require.NoError(t, td.store.db.Put(accountKey(addr), raw, nil))
	td.store.Close()

	require.Panics(t, func() {
		_, _ = NewStore(td.store.config)
	})
}

// Only a missing key is ErrNotFound. A record that does not decode, or a
// closed database, is a real error.
func TestAccountReadErrors(t *testing.T) {
	td := setup(t, nil)

	short := td.RandAccAddress()
	require.NoError(t, td.store.db.Put(accountKey(short), []byte{0x01, 0x02, 0x03}, nil))
	_, err := td.store.Account(short)
	require.Error(t, err)
	require.NotErrorIs(t, err, ErrNotFound)

	td.store.Close()
	_, err = td.store.Account(td.RandAccAddress())
	require.Error(t, err)
	require.NotErrorIs(t, err, ErrNotFound)
}

func TestAccountNotFound(t *testing.T) {
	td := setup(t, nil)

	_, err := td.store.Account(td.RandAccAddress())
	require.ErrorIs(t, err, ErrNotFound)
}
