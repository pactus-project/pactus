package executor

import (
	"bytes"
	"encoding/hex"
	"math"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestAnchorRejectedByExecutor(t *testing.T) {
	td := setup(t)
	td.sbx.FakeBlockVersion = protocol.ProtocolVersion4

	acc, addr := td.addTestAccount(t)
	before, err := acc.Bytes()
	require.NoError(t, err)
	require.Len(t, before, 12)
	require.False(t, acc.HasAnchor())

	pub, prv := td.RandBLSKeyPair()
	from := pub.AccountAddress()
	fee := td.RandFee()
	setTx := tx.NewAnchorTx(td.sbx.CurrentHeight(), from, payload.AnchorActionSet,
		bytes.Repeat([]byte{0x11}, 32), "abc", 0x04, 1, fee)
	td.HelperSignTransaction(prv, setTx)
	_, err = MakeExecutor(setTx, td.sbx)
	require.ErrorIs(t, err, InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})

	delTx := tx.NewAnchorTx(td.sbx.CurrentHeight(), from, payload.AnchorActionDelete, nil, "", 0, 0, fee)
	td.HelperSignTransaction(prv, delTx)
	_, err = MakeExecutor(delTx, td.sbx)
	require.ErrorIs(t, err, InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})

	after, err := acc.Bytes()
	require.NoError(t, err)
	require.Equal(t, before, after)
	require.False(t, acc.HasAnchor())

	receiver := td.RandAccAddress()
	amt, transferFee := td.randAmountFee(acc.Balance())
	transfer := tx.NewTransferTx(td.sbx.CurrentHeight(), addr, receiver, amt, transferFee)
	td.check(t, transfer, true, nil)
	td.execute(t, transfer)
	require.Equal(t, amt, td.sbx.Account(receiver).Balance())

	valPub, _ := td.RandBLSKeyPair()
	bondTx := tx.NewBondTx(td.sbx.CurrentHeight(), addr, valPub.ValidatorAddress(), valPub, 1, 1)
	_, err = MakeExecutor(bondTx, td.sbx)
	require.NoError(t, err)
}

func (td *testData) useAnchor(unix uint32) {
	td.sbx.FakeBlockVersion = protocol.Version(5)
	td.sbx.FakeUnixTime = unix
}

func (td *testData) setAnchorTx(addr crypto.Address, root []byte, uri string, kind uint8,
	deposit, fee amount.Amount,
) *tx.Tx {
	return tx.NewAnchorTx(td.sbx.CurrentHeight(), addr, payload.AnchorActionSet, root, uri, kind, deposit, fee)
}

func (td *testData) deleteAnchorTx(addr crypto.Address, fee amount.Amount) *tx.Tx {
	return tx.NewAnchorTx(td.sbx.CurrentHeight(), addr, payload.AnchorActionDelete, nil, "", 0, 0, fee)
}

func accountBytes(t *testing.T, acc *account.Account) []byte {
	t.Helper()

	raw, err := acc.Bytes()
	require.NoError(t, err)

	return raw
}

func requireUnchanged(t *testing.T, acc *account.Account, before []byte) {
	t.Helper()

	require.Equal(t, before, accountBytes(t, acc))
}

func TestAnchorVersionGate(t *testing.T) {
	for _, version := range []protocol.Version{0, 1, 4} {
		td := setup(t)
		td.sbx.FakeBlockVersion = version
		pub, prv := td.RandBLSKeyPair()
		_, addr := td.addTestAccount(t, testsuite.AccountWithAddress(pub.AccountAddress()),
			testsuite.AccountWithBalance(10*MinAnchorDeposit))
		acc := td.sbx.Account(addr)
		before := accountBytes(t, acc)

		setTx := td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 32), "abc", 0x04, MinAnchorDeposit, 1)
		td.HelperSignTransaction(prv, setTx)
		_, err := MakeExecutor(setTx, td.sbx)
		require.ErrorIs(t, err, InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})

		delTx := td.deleteAnchorTx(addr, 1)
		td.HelperSignTransaction(prv, delTx)
		_, err = MakeExecutor(delTx, td.sbx)
		require.ErrorIs(t, err, InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})
		requireUnchanged(t, acc, before)
		require.Len(t, before, 12)
	}

	td := setup(t)
	td.useAnchor(50)
	pub, prv := td.RandBLSKeyPair()
	_, addr := td.addTestAccount(t, testsuite.AccountWithAddress(pub.AccountAddress()),
		testsuite.AccountWithBalance(10*MinAnchorDeposit))
	setTx := td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 32), "abc", 0x04, MinAnchorDeposit, 1)
	td.HelperSignTransaction(prv, setTx)
	td.check(t, setTx, true, nil)
	td.check(t, setTx, false, nil)
	td.execute(t, setTx)
	require.True(t, td.sbx.Account(addr).HasAnchor())
}

func TestAnchorCreate(t *testing.T) {
	root := bytes.Repeat([]byte{0x11}, 32)

	t.Run("below minimum", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(10)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(10*MinAnchorDeposit))
		before := accountBytes(t, acc)
		trx := td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit-1, 1)
		td.check(t, trx, true, ErrAnchorDepositTooSmall)
		td.check(t, trx, false, ErrAnchorDepositTooSmall)
		requireUnchanged(t, acc, before)
	})

	t.Run("zero deposit", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(10)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(10*MinAnchorDeposit))
		before := accountBytes(t, acc)
		trx := td.setAnchorTx(addr, root, "", 0, 0, 1)
		td.check(t, trx, true, ErrAnchorDepositTooSmall)
		requireUnchanged(t, acc, before)
	})

	t.Run("minimum and one above", func(t *testing.T) {
		for _, deposit := range []amount.Amount{MinAnchorDeposit, MinAnchorDeposit + 1} {
			td := setup(t)
			td.useAnchor(10)
			td.sbx.FakeHeight = 7
			_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(deposit+1))
			other, _ := td.addTestAccount(t, testsuite.AccountWithBalance(3))
			otherHash := other.Hash()
			fee := amount.Amount(1)
			trx := td.setAnchorTx(addr, root, "", 0, deposit, fee)
			td.check(t, trx, true, nil)
			td.execute(t, trx)

			got := td.sbx.Account(addr)
			require.True(t, got.HasAnchor())
			require.Equal(t, deposit, got.LockedDeposit())
			require.Equal(t, amount.Amount(0), got.Balance())
			require.Equal(t, root, got.RootHash())
			require.Equal(t, 32, cap(got.RootHash()))
			require.Equal(t, types.Height(7), got.CreatedAtHeight())
			require.Equal(t, types.Height(7), got.UpdatedAtHeight())
			require.Equal(t, uint32(10), got.CreatedAtTime())
			require.Equal(t, uint32(10), got.UpdatedAtTime())
			require.Equal(t, otherHash, other.Hash())
			td.checkTotalCoin(t, fee)
		}
	})

	t.Run("unknown account", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(10)
		from := td.RandAccAddress()
		before := len(td.sbx.FakeAccounts)
		trx := td.setAnchorTx(from, root, "", 0, MinAnchorDeposit, 1)
		_, err := MakeExecutor(trx, td.sbx)
		require.ErrorIs(t, err, AccountNotFoundError{Address: from})
		require.Len(t, td.sbx.FakeAccounts, before)
		require.Nil(t, td.sbx.Account(from))
	})

	t.Run("exact balance and fee zero", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit))
		trx := td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0)
		td.check(t, trx, false, nil)
		td.execute(t, trx)
		require.Equal(t, amount.Amount(0), acc.Balance())
		require.Equal(t, MinAnchorDeposit, acc.LockedDeposit())
		td.checkTotalCoin(t, 0)
	})

	t.Run("one nano short", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit))
		before := accountBytes(t, acc)
		trx := td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 1)
		td.check(t, trx, true, ErrInsufficientFunds)
		requireUnchanged(t, acc, before)
	})

	t.Run("timestamps at the edges", func(t *testing.T) {
		for _, unixTime := range []uint32{0, math.MaxUint32} {
			td := setup(t)
			td.useAnchor(unixTime)
			td.sbx.FakeHeight = 1
			_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit))
			td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
			got := td.sbx.Account(addr)
			require.Equal(t, types.Height(1), got.CreatedAtHeight())
			require.Equal(t, unixTime, got.CreatedAtTime())
		}
	})

	t.Run("height zero is stored", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(9)
		td.sbx.FakeHeight = 0
		_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		require.Equal(t, types.Height(0), td.sbx.Account(addr).CreatedAtHeight())
	})

	t.Run("accepted shapes", func(t *testing.T) {
		shapes := []struct {
			root []byte
			uri  string
			kind uint8
		}{
			{bytes.Repeat([]byte{0x11}, 32), "", 0x00},
			{bytes.Repeat([]byte{0x22}, 64), strings.Repeat("a", 128), 0x04},
			{make([]byte, 32), "javascript:alert(1)", 0xFF},
		}
		for _, shape := range shapes {
			td := setup(t)
			td.useAnchor(3)
			_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit))
			td.execute(t, td.setAnchorTx(addr, shape.root, shape.uri, shape.kind, MinAnchorDeposit, 0))
			got := td.sbx.Account(addr)
			require.Equal(t, shape.root, got.RootHash())
			require.Equal(t, len(shape.root), cap(got.RootHash()))
			require.Equal(t, shape.uri, got.ManifestURI())
			require.Equal(t, shape.kind, got.AnchorType())
		}
	})

	t.Run("payload hash is copied", func(t *testing.T) {
		for _, n := range []int{32, 64} {
			td := setup(t)
			td.useAnchor(3)
			_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit))
			trx := td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, n), "", 0, MinAnchorDeposit, 0)
			td.execute(t, trx)
			trx.Payload().(*payload.AnchorPayload).RootHash[0] = 0xFF
			require.Equal(t, bytes.Repeat([]byte{0x11}, n), td.sbx.Account(addr).RootHash())
		}
	})

	t.Run("ceiling", func(t *testing.T) {
		maxAmt := amount.Amount(amount.MaxNanoPAC)
		td := setup(t)
		td.useAnchor(1)
		_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(maxAmt))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, maxAmt, 0))
		require.Equal(t, maxAmt, td.sbx.Account(addr).LockedDeposit())
		require.Equal(t, amount.Amount(0), td.sbx.Account(addr).Balance())
		td.checkTotalCoin(t, 0)

		td = setup(t)
		td.useAnchor(1)
		_, addr = td.addTestAccount(t, testsuite.AccountWithBalance(maxAmt))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, maxAmt-1, 1))
		require.Equal(t, maxAmt-1, td.sbx.Account(addr).LockedDeposit())
		require.Equal(t, amount.Amount(0), td.sbx.Account(addr).Balance())

		td = setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(maxAmt))
		before := accountBytes(t, acc)
		trx := td.setAnchorTx(addr, root, "", 0, maxAmt, 1)
		td.check(t, trx, true, ErrAmountOverflow)
		requireUnchanged(t, acc, before)

		td = setup(t)
		td.useAnchor(1)
		acc, addr = td.addTestAccount(t, testsuite.AccountWithBalance(maxAmt))
		before = accountBytes(t, acc)
		trx = td.setAnchorTx(addr, root, "", 0, maxAmt, maxAmt)
		td.check(t, trx, true, ErrAmountOverflow)
		requireUnchanged(t, acc, before)
	})
}

func TestAnchorUpdate(t *testing.T) {
	t.Run("zero deposit keeps the lock", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(10)
		td.sbx.FakeHeight = 4
		_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(10*MinAnchorDeposit))
		first := bytes.Repeat([]byte{0x11}, 32)
		td.execute(t, td.setAnchorTx(addr, first, "old", 0x04, 2*MinAnchorDeposit, 1))
		createdH := td.sbx.Account(addr).CreatedAtHeight()
		createdT := td.sbx.Account(addr).CreatedAtTime()

		td.sbx.FakeHeight = 9
		td.sbx.FakeUnixTime = 80
		fee := amount.Amount(1)
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x22}, 32), "new", 0x00, 0, fee))
		got := td.sbx.Account(addr)
		require.Equal(t, 2*MinAnchorDeposit, got.LockedDeposit())
		require.Equal(t, 10*MinAnchorDeposit-2*MinAnchorDeposit-2, got.Balance())
		require.Equal(t, bytes.Repeat([]byte{0x22}, 32), got.RootHash())
		require.Equal(t, "new", got.ManifestURI())
		require.Equal(t, uint8(0x00), got.AnchorType())
		require.Equal(t, createdH, got.CreatedAtHeight())
		require.Equal(t, createdT, got.CreatedAtTime())
		require.Equal(t, types.Height(9), got.UpdatedAtHeight())
		require.Equal(t, uint32(80), got.UpdatedAtTime())
		td.checkTotalCoin(t, 2)
	})

	t.Run("delta is added", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(12*MinAnchorDeposit))
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 32), "", 0, 10*MinAnchorDeposit, 0))
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x22}, 32), "", 0, MinAnchorDeposit, 0))
		require.Equal(t, 11*MinAnchorDeposit, td.sbx.Account(addr).LockedDeposit())
		td.checkTotalCoin(t, 0)
	})

	t.Run("one nano on an existing anchor is not a create", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit+1))
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 32), "", 0, MinAnchorDeposit, 0))
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x22}, 32), "", 0, 1, 0))
		require.Equal(t, MinAnchorDeposit+1, acc.LockedDeposit())
		require.Equal(t, amount.Amount(0), acc.Balance())
	})

	t.Run("lock plus deposit at the ceiling", func(t *testing.T) {
		maxAmt := amount.Amount(amount.MaxNanoPAC)
		extra := MinAnchorDeposit
		td := setup(t)
		td.useAnchor(1)
		_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(maxAmt))
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 32), "", 0, maxAmt-extra, 0))
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x22}, 32), "", 0, extra, 0))
		require.Equal(t, maxAmt, td.sbx.Account(addr).LockedDeposit())

		td = setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(maxAmt+extra))
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 32), "", 0, maxAmt, 0))
		before := accountBytes(t, acc)
		trx := td.setAnchorTx(addr, bytes.Repeat([]byte{0x22}, 32), "", 0, 1, 0)
		td.check(t, trx, true, ErrAmountOverflow)
		td.check(t, trx, false, ErrAmountOverflow)
		requireUnchanged(t, acc, before)
	})

	t.Run("second set sees the first", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(3*MinAnchorDeposit))
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 32), "one", 1, MinAnchorDeposit, 0))
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x22}, 32), "two", 2, MinAnchorDeposit, 0))
		got := td.sbx.Account(addr)
		require.Equal(t, bytes.Repeat([]byte{0x22}, 32), got.RootHash())
		require.Equal(t, "two", got.ManifestURI())
		require.Equal(t, 2*MinAnchorDeposit, got.LockedDeposit())
	})

	t.Run("second set fails and the first stays", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit+1))
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 32), "one", 1, MinAnchorDeposit, 0))
		kept := acc.Hash()
		trx := td.setAnchorTx(addr, bytes.Repeat([]byte{0x22}, 32), "two", 2, MinAnchorDeposit, 0)
		td.check(t, trx, true, ErrInsufficientFunds)
		require.Equal(t, kept, acc.Hash())
		require.Equal(t, "one", acc.ManifestURI())
	})

	t.Run("same hash still costs a fee and refreshes the time", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(10)
		td.sbx.FakeHeight = 1
		root := bytes.Repeat([]byte{0x11}, 32)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit+2))
		td.execute(t, td.setAnchorTx(addr, root, "abc", 1, MinAnchorDeposit, 1))
		td.sbx.FakeHeight = 2
		td.sbx.FakeUnixTime = 11
		before := acc.Hash()
		td.execute(t, td.setAnchorTx(addr, root, "abc", 1, 0, 1))
		require.NotEqual(t, before, acc.Hash())
		require.Equal(t, types.Height(1), acc.CreatedAtHeight())
		require.Equal(t, types.Height(2), acc.UpdatedAtHeight())
		require.Equal(t, uint32(11), acc.UpdatedAtTime())
		require.Equal(t, MinAnchorDeposit, acc.LockedDeposit())
	})

	t.Run("changing only one clock changes the hash", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(10)
		td.sbx.FakeHeight = 1
		root := bytes.Repeat([]byte{0x11}, 32)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit+3))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		first := acc.Hash()

		td.sbx.FakeHeight = 2
		td.execute(t, td.setAnchorTx(addr, root, "", 0, 0, 1))
		require.NotEqual(t, first, acc.Hash())
		heightHash := acc.Hash()

		td.sbx.FakeUnixTime = 11
		td.execute(t, td.setAnchorTx(addr, root, "", 0, 0, 1))
		require.NotEqual(t, heightHash, acc.Hash())
		require.Equal(t, types.Height(1), acc.CreatedAtHeight())
	})

	t.Run("shorter replacement drops the old tail", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(3*MinAnchorDeposit))
		longURI := strings.Repeat("x", 128)
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0xAB}, 64), longURI, 0xFF, MinAnchorDeposit, 0))
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 32), "ab", 0x00, 0, 1))
		raw := accountBytes(t, acc)
		require.False(t, bytes.Contains(raw, bytes.Repeat([]byte{0xAB}, 32)))
		require.False(t, bytes.Contains(raw, []byte(longURI)))
		require.Len(t, acc.RootHash(), 32)
		require.Equal(t, 32, cap(acc.RootHash()))
		require.Equal(t, "ab", acc.ManifestURI())
		require.Equal(t, uint8(0), acc.AnchorType())
	})

	t.Run("fee drains the balance", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit+1))
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 32), "", 0, MinAnchorDeposit, 0))
		require.Equal(t, amount.Amount(1), acc.Balance())
		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 32), "", 0, 0, 1))
		require.Equal(t, amount.Amount(0), acc.Balance())
		require.Equal(t, MinAnchorDeposit, acc.LockedDeposit())
		kept := acc.Hash()
		trx := td.setAnchorTx(addr, bytes.Repeat([]byte{0x22}, 32), "no", 1, 0, 1)
		td.check(t, trx, true, ErrInsufficientFunds)
		require.Equal(t, kept, acc.Hash())
		require.Equal(t, MinAnchorDeposit, acc.LockedDeposit())
	})
}

func TestAnchorDelete(t *testing.T) {
	root := bytes.Repeat([]byte{0x11}, 32)

	t.Run("missing anchor", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit))
		before := accountBytes(t, acc)
		hashBefore := acc.Hash()
		trx := td.deleteAnchorTx(addr, 1)
		td.check(t, trx, true, ErrAnchorNotFound)
		td.check(t, trx, false, ErrAnchorNotFound)
		requireUnchanged(t, acc, before)
		require.Equal(t, hashBefore, acc.Hash())
		require.Len(t, before, 12)
	})

	t.Run("unknown account", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		from := td.RandAccAddress()
		_, err := MakeExecutor(td.deleteAnchorTx(from, 1), td.sbx)
		require.ErrorIs(t, err, AccountNotFoundError{Address: from})
	})

	t.Run("fee comes from the refund", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithNumber(4),
			testsuite.AccountWithBalance(MinAnchorDeposit))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		require.Equal(t, amount.Amount(0), acc.Balance())
		fee := amount.Amount(1)
		td.check(t, td.deleteAnchorTx(addr, fee), true, nil)
		td.execute(t, td.deleteAnchorTx(addr, fee))
		require.False(t, acc.HasAnchor())
		require.Equal(t, MinAnchorDeposit-fee, acc.Balance())
		require.Len(t, accountBytes(t, acc), 12)

		plain := account.NewAccount(acc.Number())
		plain.AddToBalance(acc.Balance())
		require.Equal(t, accountBytes(t, plain), accountBytes(t, acc))
		td.checkTotalCoin(t, fee)
	})

	t.Run("lock equals the fee", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		td.execute(t, td.deleteAnchorTx(addr, MinAnchorDeposit))
		require.Equal(t, amount.Amount(0), acc.Balance())
		require.False(t, acc.HasAnchor())
		require.Len(t, accountBytes(t, acc), 12)
	})

	t.Run("lock just below the fee", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		before := accountBytes(t, acc)
		td.check(t, td.deleteAnchorTx(addr, MinAnchorDeposit+1), true, ErrInsufficientFunds)
		requireUnchanged(t, acc, before)
	})

	t.Run("zero fee returns the whole lock", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit+5))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		td.execute(t, td.deleteAnchorTx(addr, 0))
		require.Equal(t, MinAnchorDeposit+5, acc.Balance())
		require.False(t, acc.HasAnchor())
		td.checkTotalCoin(t, 0)
	})

	t.Run("refund would pass the ceiling", func(t *testing.T) {
		maxAmt := amount.Amount(amount.MaxNanoPAC)
		td := setup(t)
		td.useAnchor(1)
		_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit))
		funder, funderAddr := td.addTestAccount(t, testsuite.AccountWithBalance(maxAmt))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		td.check(t, tx.NewTransferTx(td.sbx.CurrentHeight(), funderAddr, addr, maxAmt, 0), true, nil)
		td.execute(t, tx.NewTransferTx(td.sbx.CurrentHeight(), funderAddr, addr, maxAmt, 0))
		require.Equal(t, maxAmt, td.sbx.Account(addr).Balance())
		require.Equal(t, amount.Amount(0), funder.Balance())
		before := accountBytes(t, td.sbx.Account(addr))
		trx := td.deleteAnchorTx(addr, 0)
		td.check(t, trx, true, ErrAmountOverflow)
		exe, err := MakeExecutor(trx, td.sbx)
		require.NoError(t, err)
		exe.Execute(td.sbx)
		requireUnchanged(t, td.sbx.Account(addr), before)
		require.True(t, td.sbx.Account(addr).HasAnchor())
	})

	t.Run("refund that lands on the ceiling", func(t *testing.T) {
		maxAmt := amount.Amount(amount.MaxNanoPAC)
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(maxAmt))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		require.Equal(t, maxAmt-MinAnchorDeposit, acc.Balance())
		td.execute(t, td.deleteAnchorTx(addr, 0))
		require.False(t, acc.HasAnchor())
		require.Equal(t, maxAmt, acc.Balance())
		require.Len(t, accountBytes(t, acc), 12)
	})

	t.Run("payload deposit is ignored", func(t *testing.T) {
		fee := amount.Amount(1)
		for _, deposit := range []amount.Amount{1, amount.Amount(amount.MaxNanoPAC), -1} {
			td := setup(t)
			td.useAnchor(1)
			acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(2*MinAnchorDeposit))
			td.execute(t, td.setAnchorTx(addr, root, "keep", 4, MinAnchorDeposit, 0))
			del := tx.NewAnchorTx(td.sbx.CurrentHeight(), addr, payload.AnchorActionDelete,
				bytes.Repeat([]byte{0xAB}, 32), "junk", 9, deposit, fee)
			td.execute(t, del)
			require.False(t, acc.HasAnchor())
			require.Equal(t, 2*MinAnchorDeposit-fee, acc.Balance())
			require.Len(t, accountBytes(t, acc), 12)
			require.False(t, bytes.Contains(accountBytes(t, acc), []byte{0xAB}))
		}
	})

	t.Run("second delete does not pay twice", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		td.execute(t, td.deleteAnchorTx(addr, 1))
		balance := acc.Balance()
		before := accountBytes(t, acc)
		td.check(t, td.deleteAnchorTx(addr, 1), true, ErrAnchorNotFound)
		require.Equal(t, balance, acc.Balance())
		requireUnchanged(t, acc, before)
	})
}

func TestAnchorSequences(t *testing.T) {
	root := bytes.Repeat([]byte{0x11}, 32)

	t.Run("delete then create", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(10)
		td.sbx.FakeHeight = 3
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(5*MinAnchorDeposit))
		td.execute(t, td.setAnchorTx(addr, root, "first", 1, MinAnchorDeposit, 1))
		td.execute(t, td.deleteAnchorTx(addr, 1))
		require.False(t, acc.HasAnchor())

		td.sbx.FakeHeight = 8
		td.sbx.FakeUnixTime = 90
		before := accountBytes(t, acc)
		td.check(t, td.setAnchorTx(addr, root, "second", 1, MinAnchorDeposit-1, 1), true, ErrAnchorDepositTooSmall)
		requireUnchanged(t, acc, before)

		td.execute(t, td.setAnchorTx(addr, bytes.Repeat([]byte{0x22}, 32), "second", 2, MinAnchorDeposit, 1))
		require.Equal(t, types.Height(8), acc.CreatedAtHeight())
		require.Equal(t, uint32(90), acc.CreatedAtTime())
		require.Equal(t, types.Height(8), acc.UpdatedAtHeight())
		require.Equal(t, MinAnchorDeposit, acc.LockedDeposit())
		td.checkTotalCoin(t, 3)
	})

	t.Run("spend the balance then delete", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit+4))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 1))
		left := acc.Balance()
		receiver := td.RandAccAddress()
		td.check(t, tx.NewTransferTx(td.sbx.CurrentHeight(), addr, receiver, left, 0), true, nil)
		td.execute(t, tx.NewTransferTx(td.sbx.CurrentHeight(), addr, receiver, left, 0))
		require.Equal(t, amount.Amount(0), acc.Balance())
		require.True(t, acc.HasAnchor())
		td.execute(t, td.deleteAnchorTx(addr, 1))
		require.False(t, acc.HasAnchor())
		require.Equal(t, MinAnchorDeposit-1, acc.Balance())
		require.Len(t, accountBytes(t, acc), 12)
	})

	t.Run("transfer cannot spend the lock", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(MinAnchorDeposit+5))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		locked := acc.LockedDeposit()
		left := acc.Balance()
		receiver := td.RandAccAddress()
		over := tx.NewTransferTx(td.sbx.CurrentHeight(), addr, receiver, left+1, 0)
		td.check(t, over, true, ErrInsufficientFunds)
		td.check(t, over, false, ErrInsufficientFunds)
		require.Equal(t, locked, acc.LockedDeposit())

		td.execute(t, tx.NewTransferTx(td.sbx.CurrentHeight(), addr, receiver, left, 0))
		require.Equal(t, amount.Amount(0), acc.Balance())
		require.Equal(t, locked, acc.LockedDeposit())
		require.True(t, acc.HasAnchor())
	})

	t.Run("plain transfer stays twelve bytes", func(t *testing.T) {
		td := setup(t)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithNumber(1), testsuite.AccountWithBalance(5))
		receiver := td.RandAccAddress()
		td.execute(t, tx.NewTransferTx(td.sbx.CurrentHeight(), addr, receiver, 3, 0))
		raw, err := hex.DecodeString("010000000200000000000000")
		require.NoError(t, err)
		require.Equal(t, raw, accountBytes(t, acc))
		expected, err := hash.FromString("c3b75f08e64a66cb980fdc03c3a0b78635a7b1db049096e8bbbd9a2873f3071a")
		require.NoError(t, err)
		require.Equal(t, expected, acc.Hash())
		require.False(t, acc.HasAnchor())
	})

	t.Run("bond and withdraw cannot use the lock", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		td.sbx.EXPECT().PowerDelta().Return(int64(0)).AnyTimes()
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(20*MinAnchorDeposit))
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		locked := acc.LockedDeposit()
		validators := len(td.sbx.FakeValidators)

		valPub, _ := td.RandBLSKeyPair()
		bondOver := tx.NewBondTx(td.sbx.CurrentHeight(), addr, valPub.ValidatorAddress(), valPub, acc.Balance()+1, 0)
		td.check(t, bondOver, true, ErrInsufficientFunds)
		require.Equal(t, locked, acc.LockedDeposit())
		require.Len(t, td.sbx.FakeValidators, validators)

		val := td.addTestValidator(t, testsuite.ValidatorWithStake(1))
		withdrawOver := tx.NewWithdrawTx(td.sbx.CurrentHeight(), val.Address(), addr, acc.Balance()+1, 0)
		td.check(t, withdrawOver, true, ErrInsufficientFunds)
		require.Equal(t, locked, acc.LockedDeposit())

		bondFee := amount.Amount(1)
		stake := acc.Balance() - bondFee
		require.GreaterOrEqual(t, stake, td.sbx.Params().MinimumStake)
		valAddr := valPub.ValidatorAddress()
		td.sbx.FakeCommittee.EXPECT().Contains(valAddr).Return(false).Times(1)
		td.sbx.EXPECT().IsJoinedCommittee(valAddr).Return(false).Times(1)
		td.sbx.EXPECT().UpdatePowerDelta(stake.ToNanoPAC()).Times(1)
		bondTx := tx.NewBondTx(td.sbx.CurrentHeight(), addr, valAddr, valPub, stake, bondFee)
		td.check(t, bondTx, true, nil)
		td.execute(t, bondTx)
		require.Equal(t, locked, acc.LockedDeposit())
		require.Equal(t, amount.Amount(0), acc.Balance())
	})

	t.Run("set and delete leave power alone", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		td.sbx.EXPECT().PowerDelta().Return(int64(0)).AnyTimes()
		_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(2*MinAnchorDeposit))
		validators := len(td.sbx.FakeValidators)
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		td.execute(t, td.deleteAnchorTx(addr, 0))
		require.Equal(t, int64(0), td.sbx.PowerDelta())
		require.Len(t, td.sbx.FakeValidators, validators)
	})
}

func TestAnchorAbsurdPayloads(t *testing.T) {
	root := bytes.Repeat([]byte{0x11}, 32)

	t.Run("rejected payloads do not debit", func(t *testing.T) {
		cases := []struct {
			name string
			tx   func(td *testData, addr crypto.Address) *tx.Tx
		}{
			{
				name: "short hash",
				tx: func(td *testData, addr crypto.Address) *tx.Tx {
					return td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 31), "", 0, MinAnchorDeposit, 1)
				},
			},
			{
				name: "long hash",
				tx: func(td *testData, addr crypto.Address) *tx.Tx {
					return td.setAnchorTx(addr, bytes.Repeat([]byte{0x11}, 65), "", 0, MinAnchorDeposit, 1)
				},
			},
			{
				name: "empty hash",
				tx: func(td *testData, addr crypto.Address) *tx.Tx {
					return td.setAnchorTx(addr, nil, "", 0, MinAnchorDeposit, 1)
				},
			},
			{
				name: "long uri",
				tx: func(td *testData, addr crypto.Address) *tx.Tx {
					return td.setAnchorTx(addr, root, strings.Repeat("a", 129), 0, MinAnchorDeposit, 1)
				},
			},
			{
				name: "bad utf-8",
				tx: func(td *testData, addr crypto.Address) *tx.Tx {
					return td.setAnchorTx(addr, root, string([]byte{0xC3}), 0, MinAnchorDeposit, 1)
				},
			},
			{
				name: "negative deposit",
				tx: func(td *testData, addr crypto.Address) *tx.Tx {
					return td.setAnchorTx(addr, root, "", 0, -1, 1)
				},
			},
		}
		for _, testCase := range cases {
			td := setup(t)
			td.useAnchor(1)
			acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(10*MinAnchorDeposit))
			before := accountBytes(t, acc)
			trx := testCase.tx(td, addr)
			exe, err := MakeExecutor(trx, td.sbx)
			require.NoError(t, err)
			require.Error(t, exe.Check(td.sbx, true))
			exe.Execute(td.sbx)
			requireUnchanged(t, acc, before)
		}
	})

	t.Run("treasury and validator", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		treasury, _ := td.addTestAccount(t, testsuite.AccountWithAddress(crypto.TreasuryAddress),
			testsuite.AccountWithBalance(10*MinAnchorDeposit))
		before := accountBytes(t, treasury)
		trx := td.setAnchorTx(crypto.TreasuryAddress, root, "", 0, MinAnchorDeposit, 1)
		exe, err := MakeExecutor(trx, td.sbx)
		require.NoError(t, err)
		require.Error(t, exe.Check(td.sbx, true))
		exe.Execute(td.sbx)
		requireUnchanged(t, treasury, before)

		valAddr := td.RandValAddress()
		valAcc, _ := td.addTestAccount(t, testsuite.AccountWithAddress(valAddr),
			testsuite.AccountWithBalance(10*MinAnchorDeposit))
		before = accountBytes(t, valAcc)
		trx = td.setAnchorTx(valAddr, root, "", 0, MinAnchorDeposit, 1)
		exe, err = MakeExecutor(trx, td.sbx)
		require.NoError(t, err)
		require.Error(t, exe.Check(td.sbx, true))
		exe.Execute(td.sbx)
		requireUnchanged(t, valAcc, before)
	})

	t.Run("action flipped after construction", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(10*MinAnchorDeposit))
		trx := td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 1)
		exe, err := MakeExecutor(trx, td.sbx)
		require.NoError(t, err)
		trx.Payload().(*payload.AnchorPayload).Action = 2
		before := accountBytes(t, acc)
		require.Error(t, exe.Check(td.sbx, true))
		require.Error(t, exe.Check(td.sbx, false))
		exe.Execute(td.sbx)
		requireUnchanged(t, acc, before)

		trx.Payload().(*payload.AnchorPayload).Action = 0xFF
		exe.Execute(td.sbx)
		requireUnchanged(t, acc, before)
	})

	t.Run("negative fee", func(t *testing.T) {
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(10*MinAnchorDeposit))
		before := accountBytes(t, acc)
		td.check(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, -1), true, ErrAmountOverflow)
		td.check(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit,
			amount.Amount(amount.MaxNanoPAC)+1), true, ErrAmountOverflow)
		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, -1))
		requireUnchanged(t, acc, before)

		td.execute(t, td.setAnchorTx(addr, root, "", 0, MinAnchorDeposit, 0))
		before = accountBytes(t, acc)
		balance := acc.Balance()
		td.check(t, td.deleteAnchorTx(addr, -1), true, ErrAmountOverflow)
		exe, err := MakeExecutor(td.deleteAnchorTx(addr, -1), td.sbx)
		require.NoError(t, err)
		exe.Execute(td.sbx)
		require.Equal(t, balance, acc.Balance())
		requireUnchanged(t, acc, before)
	})

	t.Run("execute alone does not wrap", func(t *testing.T) {
		maxAmt := amount.Amount(amount.MaxNanoPAC)
		td := setup(t)
		td.useAnchor(1)
		acc, addr := td.addTestAccount(t, testsuite.AccountWithBalance(1))
		before := accountBytes(t, acc)
		exe, err := MakeExecutor(td.setAnchorTx(addr, root, "", 0, maxAmt, 1), td.sbx)
		require.NoError(t, err)
		exe.Execute(td.sbx)
		requireUnchanged(t, acc, before)
		require.Equal(t, amount.Amount(1), acc.Balance())

		exe, err = MakeExecutor(td.deleteAnchorTx(addr, 1), td.sbx)
		require.NoError(t, err)
		exe.Execute(td.sbx)
		requireUnchanged(t, acc, before)
	})
}

func TestAnchorOnRealSandbox(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	fakeStore := store.NewFakeStore(ts)
	params := param.FromGenesis(genesis.MainnetGenesis())
	reader := committee.NewFakeCommittee(ts)
	addr, acc := ts.GenerateTestAccount(testsuite.AccountWithBalance(10 * MinAnchorDeposit))
	fakeStore.FakeAccounts[addr] = acc

	sbx := sandbox.NewSandbox(10, fakeStore, params, reader, 0)
	require.Equal(t, params.BlockVersion, sbx.BlockVersion())
	require.Equal(t, uint32(0), sbx.CurrentUnixTime())

	const unixTime = uint32(1_700_000_000)
	sbx.SetBlockContext(protocol.Version(5), unixTime)
	require.Equal(t, protocol.Version(5), sbx.BlockVersion())
	require.Equal(t, unixTime, sbx.CurrentUnixTime())

	root := bytes.Repeat([]byte{0x11}, 32)
	trx := tx.NewAnchorTx(sbx.CurrentHeight(), addr, payload.AnchorActionSet, root, "abc", 4, MinAnchorDeposit, 1)
	exe, err := MakeExecutor(trx, sbx)
	require.NoError(t, err)
	require.NoError(t, exe.Check(sbx, true))
	untouched := sbx.Account(addr)
	require.False(t, untouched.HasAnchor())
	require.Equal(t, 10*MinAnchorDeposit, untouched.Balance())

	exe.Execute(sbx)

	got := sbx.Account(addr)
	require.True(t, got.HasAnchor())
	require.Equal(t, root, got.RootHash())
	require.Equal(t, MinAnchorDeposit, got.LockedDeposit())
	require.Equal(t, 10*MinAnchorDeposit-MinAnchorDeposit-1, got.Balance())
	require.Equal(t, sbx.CurrentHeight(), got.CreatedAtHeight())
	require.Equal(t, unixTime, got.CreatedAtTime())
	require.False(t, fakeStore.FakeAccounts[addr].HasAnchor())
	require.Len(t, accountBytes(t, fakeStore.FakeAccounts[addr]), 12)

	del := tx.NewAnchorTx(sbx.CurrentHeight(), addr, payload.AnchorActionDelete, nil, "", 0, 0, 1)
	delExe, err := MakeExecutor(del, sbx)
	require.NoError(t, err)
	delExe.Execute(sbx)
	after := sbx.Account(addr)
	require.False(t, after.HasAnchor())
	require.Len(t, accountBytes(t, after), 12)
	require.Equal(t, 10*MinAnchorDeposit-2, after.Balance())
	require.Equal(t, int64(0), sbx.PowerDelta())
}
