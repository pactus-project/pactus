package store

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestPruneKeepsAnchorAccount(t *testing.T) {
	td := setup(t, nil)
	addr, _ := saveAnchoredAccount(t, td, 3, 3*1e9, keptAnchor())
	setTx := anchorSetTx(addr, 1, bytesFill(0x11), 7*1e9)
	saveTxBlock(t, td, 11, setTx)
	pruneHeight(t, td, 11)

	got := mustStoreAccount(t, td.store, addr)
	assertKeptAnchor(t, got)
	require.Equal(t, amount.Amount(3*1e9), got.Balance())
	requireMissingBlock(t, td.store, 11)
	requireMissingTx(t, td.store, setTx)
}

func TestPruneCreateKeepsLaterUpdate(t *testing.T) {
	td := setup(t, nil)
	addr, _ := td.GenerateTestAccount(testsuite.AccountWithNumber(4), testsuite.AccountWithBalance(1))
	createRoot := bytesFill(0x11)
	updateRoot := bytesFill(0x22)
	addrAccount(t, td, addr, account.AnchorData{
		RootHash:        updateRoot,
		ManifestURI:     "second",
		LockedDeposit:   8,
		CreatedAtHeight: 3,
		CreatedAtTime:   111,
		UpdatedAtHeight: 9,
		UpdatedAtTime:   222,
	})

	createTx := anchorSetTx(addr, 1, createRoot, 7)
	updateTx := anchorSetTx(addr, 2, updateRoot, 1)
	saveTxBlock(t, td, 11, createTx)
	saveTxBlock(t, td, 12, updateTx)
	pruneHeight(t, td, 11)

	got := mustStoreAccount(t, td.store, addr)
	require.Equal(t, updateRoot, got.RootHash())
	require.Equal(t, types.Height(3), got.CreatedAtHeight())
	require.Equal(t, uint32(111), got.CreatedAtTime())
	require.Equal(t, types.Height(9), got.UpdatedAtHeight())
	require.Equal(t, uint32(222), got.UpdatedAtTime())
	requireMissingTx(t, td.store, createTx)

	kept, err := td.store.Transaction(updateTx.ID())
	require.NoError(t, err)
	require.NotNil(t, kept)
	require.Equal(t, types.Height(12), kept.Height)
}

func TestPruneDeleteLeavesPlainAccount(t *testing.T) {
	td := setup(t, nil)
	addr, plain := td.GenerateTestAccount(testsuite.AccountWithNumber(5), testsuite.AccountWithBalance(2*1e9))
	td.store.UpdateAccount(addr, plain)
	require.NoError(t, td.store.WriteBatch())

	setTx := anchorSetTx(addr, 1, bytesFill(0x31), 4*1e9)
	delTx := tx.NewAnchorTx(2, addr, payload.AnchorActionDelete, nil, "", 0, 0, 1)
	saveTxBlock(t, td, 11, setTx, delTx)

	other, _ := saveAnchoredAccount(t, td, 6, 1, account.AnchorData{
		RootHash:        bytesFill(0x44),
		ManifestURI:     "stays",
		LockedDeposit:   1,
		CreatedAtHeight: 12,
		CreatedAtTime:   300,
		UpdatedAtHeight: 12,
		UpdatedAtTime:   300,
	})
	otherTx := anchorSetTx(other, 1, bytesFill(0x44), 1)
	saveTxBlock(t, td, 12, otherTx)
	pruneHeight(t, td, 11)

	got := mustStoreAccount(t, td.store, addr)
	require.False(t, got.HasAnchor())
	require.Len(t, accountBytes(t, got), 12)
	require.Equal(t, amount.Amount(2*1e9), got.Balance())
	requireMissingTx(t, td.store, setTx)
	requireMissingTx(t, td.store, delTx)

	kept := mustStoreAccount(t, td.store, other)
	require.True(t, kept.HasAnchor())
	require.Equal(t, bytesFill(0x44), kept.RootHash())
	require.Greater(t, len(accountBytes(t, kept)), 12)
	otherStored, err := td.store.Transaction(otherTx.ID())
	require.NoError(t, err)
	require.NotNil(t, otherStored)
}

func TestReopenKeepsPrunedAnchor(t *testing.T) {
	td := setup(t, nil)
	addr, _ := saveAnchoredAccount(t, td, 3, 3*1e9, keptAnchor())
	setTx := anchorSetTx(addr, 1, bytesFill(0x11), 7*1e9)
	saveTxBlock(t, td, 11, setTx)
	pruneHeight(t, td, 11)
	requireMissingTx(t, td.store, setTx)

	saveTxBlock(t, td, 12)
	saveTxBlock(t, td, 13)
	pruneHeight(t, td, 1)
	td.store.config.TxCacheWindow = 1

	td.store.Close()
	reopened, err := NewStore(td.store.config)
	require.NoError(t, err)
	t.Cleanup(func() {
		reopened.Close()
	})

	require.True(t, reopened.IsPruned())
	got := mustStoreAccount(t, reopened, addr)
	assertKeptAnchor(t, got)
	requireMissingTx(t, reopened, setTx)
}

func keptAnchor() account.AnchorData {
	return account.AnchorData{
		RootHash:        bytesFill(0x11),
		ManifestURI:     "keep",
		AnchorType:      4,
		LockedDeposit:   7 * 1e9,
		CreatedAtHeight: 3,
		CreatedAtTime:   111,
		UpdatedAtHeight: 4,
		UpdatedAtTime:   222,
	}
}

func assertKeptAnchor(t *testing.T, got *account.Account) {
	t.Helper()

	require.True(t, got.HasAnchor())
	require.Equal(t, bytesFill(0x11), got.RootHash())
	require.Equal(t, "keep", got.ManifestURI())
	require.Equal(t, amount.Amount(7*1e9), got.LockedDeposit())
	require.Equal(t, types.Height(3), got.CreatedAtHeight())
	require.Equal(t, uint32(111), got.CreatedAtTime())
	require.Equal(t, types.Height(4), got.UpdatedAtHeight())
	require.Equal(t, uint32(222), got.UpdatedAtTime())
	require.Greater(t, len(accountBytes(t, got)), 12)
}

func saveAnchoredAccount(t *testing.T, td *testData, number int32, balance amount.Amount,
	data account.AnchorData,
) (crypto.Address, *account.Account) {
	t.Helper()

	addr, acc := td.GenerateTestAccount(
		testsuite.AccountWithNumber(number),
		testsuite.AccountWithBalance(balance),
	)
	require.NoError(t, acc.SetAnchor(data))
	td.store.UpdateAccount(addr, acc)
	require.NoError(t, td.store.WriteBatch())

	return addr, acc
}

func addrAccount(t *testing.T, td *testData, addr crypto.Address, data account.AnchorData) *account.Account {
	t.Helper()

	acc := account.NewAccount(4)
	acc.AddToBalance(1)
	require.NoError(t, acc.SetAnchor(data))
	td.store.UpdateAccount(addr, acc)
	require.NoError(t, td.store.WriteBatch())

	return acc
}

func saveTxBlock(t *testing.T, td *testData, height types.Height, txs ...*tx.Tx) {
	t.Helper()

	var blk *block.Block
	var cert *certificate.Certificate
	if len(txs) == 0 {
		blk, cert = td.GenerateTestBlock(height)
	} else {
		blk, cert = td.GenerateTestBlock(height, testsuite.BlockWithTransactions(txs))
	}
	td.store.SaveBlock(blk, cert)
	require.NoError(t, td.store.WriteBatch())
}

func pruneHeight(t *testing.T, td *testData, height types.Height) {
	t.Helper()

	pruned, err := td.store.pruneBlock(height)
	require.NoError(t, err)
	require.True(t, pruned)
	require.NoError(t, td.store.WriteBatch())
}

func anchorSetTx(from crypto.Address, stamp types.Height, root []byte, deposit amount.Amount) *tx.Tx {
	return tx.NewAnchorTx(stamp, from, payload.AnchorActionSet, root, "manifest", 1, deposit, 1)
}

func mustStoreAccount(t *testing.T, reader Reader, addr crypto.Address) *account.Account {
	t.Helper()

	got, err := reader.Account(addr)
	require.NoError(t, err)

	return got
}

func requireMissingBlock(t *testing.T, reader Reader, height types.Height) {
	t.Helper()

	blk, err := reader.Block(height)
	require.Error(t, err)
	require.Nil(t, blk)
}

func requireMissingTx(t *testing.T, reader Reader, trx *tx.Tx) {
	t.Helper()

	got, err := reader.Transaction(trx.ID())
	require.Error(t, err)
	require.Nil(t, got)
}

func accountBytes(t *testing.T, acc *account.Account) []byte {
	t.Helper()

	bs, err := acc.Bytes()
	require.NoError(t, err)

	return bs
}

func bytesFill(fill byte) []byte {
	return bytes.Repeat([]byte{fill}, 32)
}
