package html_test

import (
	"bytes"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestAccountAnchorRows(t *testing.T) {
	td := setup(t)

	balance := amount.Amount(2 * 1e9)
	lock := amount.Amount(7 * 1e9)
	addr, acc := td.GenerateTestAccount(testsuite.AccountWithBalance(balance))
	root := bytesRepeat(0xAB)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:        root,
		ManifestURI:     "first\nsecond",
		AnchorType:      4,
		LockedDeposit:   lock,
		CreatedAtHeight: 3,
		UpdatedAtHeight: 4,
	}))
	td.gRPCServer.FakeState.AddTestAccount(addr, acc)

	body := accountPage(t, td, addr.String())
	require.Contains(t, body, "Anchor Type")
	require.Contains(t, body, "first")
	require.Contains(t, body, "second")
	require.Greater(t, indexOf(body, "Locked Deposit"), indexOf(body, "second"))
	require.Contains(t, body, balance.String())
	require.Contains(t, body, lock.String())
	require.NotEqual(t, balance.String(), lock.String())
	require.Contains(t, body, hex.EncodeToString(root))

	next := bytesRepeat(0xCD)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:        next,
		ManifestURI:     "first\nsecond",
		AnchorType:      4,
		LockedDeposit:   lock + 1e9,
		CreatedAtHeight: 3,
		UpdatedAtHeight: 9,
	}))
	updated := accountPage(t, td, addr.String())
	require.Contains(t, updated, hex.EncodeToString(next))
	require.NotContains(t, updated, hex.EncodeToString(root))
	require.Contains(t, updated, (lock + 1e9).String())

	acc.ClearAnchor()
	cleared := accountPage(t, td, addr.String())
	require.NotContains(t, cleared, "Anchor Type")
	require.NotContains(t, cleared, "Locked Deposit")
}

func TestTransactionAnchorRows(t *testing.T) {
	td := setup(t)

	sender, prv := td.RandBLSKeyPair()
	root := bytesRepeat(0x21)
	trx := tx.NewAnchorTx(1, sender.AccountAddress(), payload.AnchorActionSet, root, "manifest", 4,
		amount.Amount(9), amount.Amount(1), tx.WithMemo("note"))
	trx.SetSignature(prv.Sign(trx.SignBytes()))
	trx.SetPublicKey(prv.PublicKey())
	blk, cert := td.GenerateTestBlock(3, testsuite.BlockWithTransactions([]*tx.Tx{trx}))
	td.gRPCServer.FakeState.AddTestBlock(blk, cert)

	w := httptest.NewRecorder()
	r := mux.SetURLVars(new(http.Request), map[string]string{"id": trx.ID().String()})
	td.httpServer.GetTransactionHandler(w, r)
	require.Equal(t, 200, w.Code)
	body := w.Body.String()
	require.Contains(t, body, "Action")
	require.Contains(t, body, hex.EncodeToString(root))
	require.Contains(t, body, "manifest")
	require.Contains(t, body, amount.Amount(9).String())
}

func TestTransactionAnchorDelete(t *testing.T) {
	td := setup(t)

	sender, prv := td.RandBLSKeyPair()
	trx := tx.NewAnchorTx(1, sender.AccountAddress(), payload.AnchorActionDelete, nil, "", 0, 0, 1)
	signAnchor(prv, trx)
	blk, cert := td.GenerateTestBlock(3, testsuite.BlockWithTransactions([]*tx.Tx{trx}))
	td.gRPCServer.FakeState.AddTestBlock(blk, cert)

	body := transactionPage(t, td, trx.ID().String())
	require.Contains(t, body, "<td>Action</td><td>1</td>")
	require.Contains(t, body, amount.Amount(0).String())
	require.NotContains(t, body, "Root Hash")
	require.NotContains(t, body, "Manifest URI")
}

func TestTransactionAnchorSetBesideDelete(t *testing.T) {
	td := setup(t)

	sender, prv := td.RandBLSKeyPair()
	addr := sender.AccountAddress()
	root := bytesRepeat(0x21)
	deposit := amount.Amount(9)
	setTx := tx.NewAnchorTx(1, addr, payload.AnchorActionSet, root, "manifest", 4, deposit, 1)
	delTx := tx.NewAnchorTx(2, addr, payload.AnchorActionDelete, nil, "", 0, 0, 1)
	signAnchor(prv, setTx)
	signAnchor(prv, delTx)
	blk, cert := td.GenerateTestBlock(3, testsuite.BlockWithTransactions([]*tx.Tx{setTx, delTx}))
	td.gRPCServer.FakeState.AddTestBlock(blk, cert)

	setPage := transactionPage(t, td, setTx.ID().String())
	require.Contains(t, setPage, hex.EncodeToString(root))
	require.Contains(t, setPage, deposit.String())

	delPage := transactionPage(t, td, delTx.ID().String())
	require.Contains(t, delPage, "<td>Action</td><td>1</td>")
	require.Contains(t, delPage, amount.Amount(0).String())
	require.NotContains(t, delPage, "Root Hash")
}

func TestAccountClearedKeepsSetPage(t *testing.T) {
	td := setup(t)

	pub, prv := td.RandBLSKeyPair()
	addr := pub.AccountAddress()
	lock := amount.Amount(7 * 1e9)
	fee := amount.Amount(1)
	acc := account.NewAccount(3)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytesRepeat(0x21),
		LockedDeposit: lock,
	}))
	acc.ClearAnchor()
	acc.AddToBalance(lock - fee)
	td.gRPCServer.FakeState.AddTestAccount(addr, acc)

	root := bytesRepeat(0x21)
	setTx := tx.NewAnchorTx(1, addr, payload.AnchorActionSet, root, "manifest", 4, lock, fee)
	signAnchor(prv, setTx)
	blk, cert := td.GenerateTestBlock(3, testsuite.BlockWithTransactions([]*tx.Tx{setTx}))
	td.gRPCServer.FakeState.AddTestBlock(blk, cert)

	page := accountPage(t, td, addr.String())
	require.NotContains(t, page, "Anchor Type")
	require.NotContains(t, page, "Locked Deposit")
	require.Contains(t, page, (lock - fee).String())

	setPage := transactionPage(t, td, setTx.ID().String())
	require.Contains(t, setPage, hex.EncodeToString(root))
	require.Contains(t, setPage, lock.String())
}

func TestTransactionAnchorEscapesURI(t *testing.T) {
	td := setup(t)

	sender, prv := td.RandBLSKeyPair()
	root := bytesRepeat(0x21)
	deposit := amount.Amount(9)
	trx := tx.NewAnchorTx(1, sender.AccountAddress(), payload.AnchorActionSet, root, "a<b\nsecond", 4,
		deposit, 1)
	signAnchor(prv, trx)
	blk, cert := td.GenerateTestBlock(3, testsuite.BlockWithTransactions([]*tx.Tx{trx}))
	td.gRPCServer.FakeState.AddTestBlock(blk, cert)

	body := transactionPage(t, td, trx.ID().String())
	require.Contains(t, body, "a&lt;b")
	require.NotContains(t, body, "a<b")
	require.Greater(t, indexOf(body, "Deposit"), indexOf(body, "a&lt;b"))
	require.Greater(t, indexOf(body, "Deposit"), indexOf(body, "second"))
	require.Contains(t, body, deposit.String())
}

func transactionPage(t *testing.T, td *testData, id string) string {
	t.Helper()

	w := httptest.NewRecorder()
	r := mux.SetURLVars(new(http.Request), map[string]string{"id": id})
	td.httpServer.GetTransactionHandler(w, r)
	require.Equal(t, 200, w.Code)

	return w.Body.String()
}

func signAnchor(prv crypto.PrivateKey, trx *tx.Tx) {
	trx.SetSignature(prv.Sign(trx.SignBytes()))
	trx.SetPublicKey(prv.PublicKey())
}

func accountPage(t *testing.T, td *testData, address string) string {
	t.Helper()

	w := httptest.NewRecorder()
	r := mux.SetURLVars(new(http.Request), map[string]string{"address": address})
	td.httpServer.GetAccountHandler(w, r)
	require.Equal(t, 200, w.Code)

	return w.Body.String()
}

func indexOf(body, part string) int {
	for i := 0; i+len(part) <= len(body); i++ {
		if body[i:i+len(part)] == part {
			return i
		}
	}

	return -1
}

func bytesRepeat(fill byte) []byte {
	return bytes.Repeat([]byte{fill}, 32)
}
