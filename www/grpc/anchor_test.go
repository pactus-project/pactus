package grpc_test

import (
	"encoding/hex"
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetAccountWithoutAnchor(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	addr, acc := td.GenerateTestAccount(testsuite.AccountWithBalance(5 * 1e9))
	td.FakeState.AddTestAccount(addr, acc)

	res, err := client.GetAccount(t.Context(), &pactus.GetAccountRequest{Address: addr.String()})
	require.NoError(t, err)
	require.Nil(t, res.Account.Anchor)
	data, err := hex.DecodeString(res.Account.Data)
	require.NoError(t, err)
	require.Len(t, data, 12)
	require.Equal(t, acc.Balance().ToNanoPAC(), res.Account.Balance)
	require.Equal(t, acc.Hash().String(), res.Account.Hash)
	require.Equal(t, hash.CalcHash(data).String(), res.Account.Hash)

	got, err := client.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: addr.String()})
	require.NoError(t, err)
	require.False(t, got.Found)
	require.Nil(t, got.Anchor)
	require.Equal(t, addr.String(), got.Address)
}

func TestGetAccountAnchorFields(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	root := bytesRepeat(0x11, 32)
	addr, acc := td.GenerateTestAccount(testsuite.AccountWithBalance(5 * 1e9))
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:        root,
		ManifestURI:     "line1\nline2",
		AnchorType:      0x00,
		LockedDeposit:   7 * 1e9,
		CreatedAtHeight: 3,
		CreatedAtTime:   111,
		UpdatedAtHeight: 4,
		UpdatedAtTime:   222,
	}))
	root[0] = 0
	td.FakeState.AddTestAccount(addr, acc)

	res, err := client.GetAccount(t.Context(), &pactus.GetAccountRequest{Address: addr.String()})
	require.NoError(t, err)
	require.Equal(t, int64(5*1e9), res.Account.Balance)
	require.Equal(t, int64(7*1e9), res.Account.Anchor.LockedDeposit)
	require.Equal(t, bytesRepeat(0x11, 32), res.Account.Anchor.RootHash)
	require.Equal(t, "line1\nline2", res.Account.Anchor.ManifestUri)
	require.Equal(t, uint32(0), res.Account.Anchor.AnchorType)
	require.Equal(t, uint32(3), res.Account.Anchor.CreatedAtHeight)
	require.Equal(t, uint32(111), res.Account.Anchor.CreatedAtTime)
	require.Equal(t, uint32(4), res.Account.Anchor.UpdatedAtHeight)
	require.Equal(t, uint32(222), res.Account.Anchor.UpdatedAtTime)

	data, err := hex.DecodeString(res.Account.Data)
	require.NoError(t, err)
	require.Greater(t, len(data), 12)
	require.Equal(t, hash.CalcHash(data).String(), res.Account.Hash)

	plain := account.NewAccount(acc.Number())
	plain.AddToBalance(acc.Balance())
	require.NotEqual(t, plain.Hash().String(), res.Account.Hash)

	again, err := client.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: addr.String()})
	require.NoError(t, err)
	require.True(t, again.Found)
	require.Equal(t, uint32(111), again.Anchor.CreatedAtTime)
	require.Equal(t, uint32(222), again.Anchor.UpdatedAtTime)

	wide := bytesRepeat(0xFF, 64)
	wideAddr, wideAcc := td.GenerateTestAccount(testsuite.AccountWithBalance(1))
	require.NoError(t, wideAcc.SetAnchor(account.AnchorData{
		RootHash:      wide,
		AnchorType:    0xFF,
		LockedDeposit: amount.MaxNanoPAC,
	}))
	td.FakeState.AddTestAccount(wideAddr, wideAcc)
	wideRes, err := client.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: wideAddr.String()})
	require.NoError(t, err)
	require.True(t, wideRes.Found)
	require.Equal(t, wide, wideRes.Anchor.RootHash)
	require.Equal(t, uint32(0xFF), wideRes.Anchor.AnchorType)
	require.Equal(t, int64(amount.MaxNanoPAC), wideRes.Anchor.LockedDeposit)
	require.Positive(t, wideRes.Anchor.LockedDeposit)
}

func TestGetAnchorSecondSetKeepsCreated(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)
	addr, acc := td.GenerateTestAccount(testsuite.AccountWithBalance(20 * 1e9))
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:        bytesRepeat(0x22, 32),
		LockedDeposit:   9 * 1e9,
		CreatedAtHeight: 3,
		CreatedAtTime:   111,
		UpdatedAtHeight: 8,
		UpdatedAtTime:   333,
	}))
	td.FakeState.AddTestAccount(addr, acc)

	res, err := client.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: addr.String()})
	require.NoError(t, err)
	require.Equal(t, uint32(3), res.Anchor.CreatedAtHeight)
	require.Equal(t, uint32(111), res.Anchor.CreatedAtTime)
	require.Equal(t, uint32(8), res.Anchor.UpdatedAtHeight)
	require.Equal(t, uint32(333), res.Anchor.UpdatedAtTime)
	require.Equal(t, int64(9*1e9), res.Anchor.LockedDeposit)
	require.Equal(t, bytesRepeat(0x22, 32), res.Anchor.RootHash)
}

func TestGetAnchorZeroRecordIsPresent(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)
	addr, acc := td.GenerateTestAccount()
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      make([]byte, 32),
		LockedDeposit: 1,
	}))
	td.FakeState.AddTestAccount(addr, acc)

	res, err := client.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: addr.String()})
	require.NoError(t, err)
	require.True(t, res.Found)
	require.NotNil(t, res.Anchor)
	require.Equal(t, make([]byte, 32), res.Anchor.RootHash)
	require.Empty(t, res.Anchor.ManifestUri)
	require.Equal(t, uint32(0), res.Anchor.AnchorType)
	require.Equal(t, uint32(0), res.Anchor.CreatedAtTime)
	require.Equal(t, uint32(0), res.Anchor.UpdatedAtTime)

	emptyAddr, emptyAcc := td.GenerateTestAccount()
	td.FakeState.AddTestAccount(emptyAddr, emptyAcc)
	missing, err := client.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: emptyAddr.String()})
	require.NoError(t, err)
	require.False(t, missing.Found)
	require.Nil(t, missing.Anchor)
}

func TestGetAnchorAfterDelete(t *testing.T) {
	td := setup(t, nil)
	chain := td.blockchainClient(t)
	txs := td.transactionClient(t)

	pub, prv := td.RandBLSKeyPair()
	addr := pub.AccountAddress()
	acc := account.NewAccount(7)
	lock := amount.Amount(4 * 1e9)
	fee := amount.Amount(1)
	root := bytesRepeat(0x31, 32)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      root,
		ManifestURI:   "keep",
		AnchorType:    4,
		LockedDeposit: lock,
	}))
	acc.ClearAnchor()
	acc.AddToBalance(lock - fee)
	td.FakeState.AddTestAccount(addr, acc)

	setTx := tx.NewAnchorTx(1, addr, payload.AnchorActionSet, root, "keep", 4,
		lock, fee, tx.WithMemo("set"))
	signTx(prv, setTx)
	delTx := tx.NewAnchorTx(2, addr, payload.AnchorActionDelete, nil, "", 0, 0, fee, tx.WithMemo("del"))
	signTx(prv, delTx)
	blk, cert := td.GenerateTestBlock(5, testsuite.BlockWithTransactions([]*tx.Tx{setTx, delTx}))
	td.FakeState.AddTestBlock(blk, cert)

	accRes, err := chain.GetAccount(t.Context(), &pactus.GetAccountRequest{Address: addr.String()})
	require.NoError(t, err)
	require.Nil(t, accRes.Account.Anchor)
	data, err := hex.DecodeString(accRes.Account.Data)
	require.NoError(t, err)
	require.Len(t, data, 12)
	require.Equal(t, (lock - fee).ToNanoPAC(), accRes.Account.Balance)

	plain := account.NewAccount(acc.Number())
	plain.AddToBalance(lock - fee)
	require.Equal(t, plain.Hash().String(), accRes.Account.Hash)

	got, err := chain.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: addr.String()})
	require.NoError(t, err)
	require.False(t, got.Found)
	require.Nil(t, got.Anchor)
	require.Equal(t, addr.String(), got.Address)

	setRes, err := txs.GetTransaction(t.Context(), &pactus.GetTransactionRequest{
		Id:        setTx.ID().String(),
		Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_INFO,
	})
	require.NoError(t, err)
	require.Equal(t, pactus.PayloadType_PAYLOAD_TYPE_ANCHOR, setRes.Transaction.PayloadType)
	require.Equal(t, uint32(payload.AnchorActionSet), setRes.Transaction.GetAnchor().Action)
	require.Equal(t, root, setRes.Transaction.GetAnchor().RootHash)
	require.Equal(t, lock.ToNanoPAC(), setRes.Transaction.GetAnchor().Deposit)

	delRes, err := txs.GetTransaction(t.Context(), &pactus.GetTransactionRequest{
		Id:        delTx.ID().String(),
		Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_INFO,
	})
	require.NoError(t, err)
	require.NotEqual(t, setTx.ID(), delTx.ID())
	require.Equal(t, pactus.PayloadType_PAYLOAD_TYPE_ANCHOR, delRes.Transaction.PayloadType)
	require.Equal(t, uint32(payload.AnchorActionDelete), delRes.Transaction.GetAnchor().Action)
	require.Empty(t, delRes.Transaction.GetAnchor().RootHash)
	require.Equal(t, int64(0), delRes.Transaction.Value)
	require.Equal(t, int64(0), delRes.Transaction.GetAnchor().Deposit)
	require.Equal(t, fee.ToNanoPAC(), delRes.Transaction.Fee)
	require.NotEqual(t, delRes.Transaction.Fee, delRes.Transaction.GetAnchor().Deposit)
}

func TestListAnchorsDropsDeletedAccount(t *testing.T) {
	td := setup(t, nil)
	chain := td.blockchainClient(t)
	txs := td.transactionClient(t)

	pub, prv := td.RandBLSKeyPair()
	addr := pub.AccountAddress()
	acc := account.NewAccount(1)
	root := bytesRepeat(0x11, 32)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      root,
		LockedDeposit: 1,
	}))
	td.FakeState.AddTestAccount(addr, acc)
	neighbor := putAnchor(t, td, 2, 0x22)

	setTx := tx.NewAnchorTx(1, addr, payload.AnchorActionSet, root, "", 0, 1, 1)
	signTx(prv, setTx)
	blk, cert := td.GenerateTestBlock(4, testsuite.BlockWithTransactions([]*tx.Tx{setTx}))
	td.FakeState.AddTestBlock(blk, cert)

	acc.ClearAnchor()

	page, err := chain.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Count: 20})
	require.NoError(t, err)
	require.Equal(t, uint32(1), page.Total)
	require.Equal(t, []string{neighbor}, itemAddresses(page))

	setRes, err := txs.GetTransaction(t.Context(), &pactus.GetTransactionRequest{
		Id:        setTx.ID().String(),
		Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_INFO,
	})
	require.NoError(t, err)
	require.Equal(t, root, setRes.Transaction.GetAnchor().RootHash)
}

func TestGetAnchorRecreateAfterDelete(t *testing.T) {
	td := setup(t, nil)
	chain := td.blockchainClient(t)
	txs := td.transactionClient(t)

	pub, prv := td.RandBLSKeyPair()
	addr := pub.AccountAddress()
	acc := account.NewAccount(4)
	oldRoot := bytesRepeat(0x11, 32)
	newRoot := bytesRepeat(0x22, 32)
	oldLock := amount.Amount(4 * 1e9)
	newLock := amount.Amount(1 * 1e9)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:        newRoot,
		ManifestURI:     "new",
		AnchorType:      2,
		LockedDeposit:   newLock,
		CreatedAtHeight: 9,
		CreatedAtTime:   900,
		UpdatedAtHeight: 9,
		UpdatedAtTime:   900,
	}))
	td.FakeState.AddTestAccount(addr, acc)

	oldSet := tx.NewAnchorTx(1, addr, payload.AnchorActionSet, oldRoot, "old", 1, oldLock, 1)
	signTx(prv, oldSet)
	delTx := tx.NewAnchorTx(2, addr, payload.AnchorActionDelete, nil, "", 0, 0, 1)
	signTx(prv, delTx)
	newSet := tx.NewAnchorTx(3, addr, payload.AnchorActionSet, newRoot, "new", 2, newLock, 1)
	signTx(prv, newSet)
	blk, cert := td.GenerateTestBlock(9, testsuite.BlockWithTransactions([]*tx.Tx{oldSet, delTx, newSet}))
	td.FakeState.AddTestBlock(blk, cert)

	got, err := chain.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: addr.String()})
	require.NoError(t, err)
	require.True(t, got.Found)
	require.Equal(t, newRoot, got.Anchor.RootHash)
	require.Equal(t, uint32(9), got.Anchor.CreatedAtHeight)
	require.Equal(t, uint32(900), got.Anchor.CreatedAtTime)
	require.Equal(t, got.Anchor.CreatedAtHeight, got.Anchor.UpdatedAtHeight)
	require.Equal(t, got.Anchor.CreatedAtTime, got.Anchor.UpdatedAtTime)
	require.Equal(t, newLock.ToNanoPAC(), got.Anchor.LockedDeposit)

	oldRes, err := txs.GetTransaction(t.Context(), &pactus.GetTransactionRequest{
		Id:        oldSet.ID().String(),
		Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_INFO,
	})
	require.NoError(t, err)
	require.Equal(t, oldRoot, oldRes.Transaction.GetAnchor().RootHash)
	require.Equal(t, oldLock.ToNanoPAC(), oldRes.Transaction.GetAnchor().Deposit)

	page, err := chain.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Count: 20})
	require.NoError(t, err)
	require.Equal(t, uint32(1), page.Total)
	require.Equal(t, []string{addr.String()}, itemAddresses(page))
	require.Equal(t, newRoot, page.Items[0].Anchor.RootHash)
}

func TestGetAnchorWithoutTxBody(t *testing.T) {
	td := setup(t, nil)
	chain := td.blockchainClient(t)
	txs := td.transactionClient(t)

	pub, prv := td.RandBLSKeyPair()
	addr := pub.AccountAddress()
	acc := account.NewAccount(3)
	root := bytesRepeat(0x41, 32)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:        root,
		ManifestURI:     "kept",
		LockedDeposit:   9,
		CreatedAtHeight: 3,
		CreatedAtTime:   111,
		UpdatedAtHeight: 8,
		UpdatedAtTime:   222,
	}))
	td.FakeState.AddTestAccount(addr, acc)

	setTx := tx.NewAnchorTx(1, addr, payload.AnchorActionSet, root, "kept", 0, 9, 1)
	signTx(prv, setTx)

	got, err := chain.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: addr.String()})
	require.NoError(t, err)
	require.True(t, got.Found)
	require.Equal(t, root, got.Anchor.RootHash)
	require.Equal(t, uint32(3), got.Anchor.CreatedAtHeight)
	require.Equal(t, uint32(111), got.Anchor.CreatedAtTime)
	require.Equal(t, uint32(8), got.Anchor.UpdatedAtHeight)
	require.Equal(t, uint32(222), got.Anchor.UpdatedAtTime)

	_, err = txs.GetTransaction(t.Context(), &pactus.GetTransactionRequest{
		Id:        setTx.ID().String(),
		Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_INFO,
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestListAnchorsWithoutTxBody(t *testing.T) {
	td := setup(t, nil)
	chain := td.blockchainClient(t)

	pub, _ := td.RandBLSKeyPair()
	addr := pub.AccountAddress()
	acc := account.NewAccount(6)
	root := bytesRepeat(0x42, 32)
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      root,
		LockedDeposit: 1,
	}))
	td.FakeState.AddTestAccount(addr, acc)

	page, err := chain.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Count: 20})
	require.NoError(t, err)
	require.Equal(t, uint32(1), page.Total)
	require.Equal(t, []string{addr.String()}, itemAddresses(page))
	require.Equal(t, root, page.Items[0].Anchor.RootHash)
}

func TestGetAccountAfterPrunedDelete(t *testing.T) {
	td := setup(t, nil)
	chain := td.blockchainClient(t)
	txs := td.transactionClient(t)

	pub, prv := td.RandBLSKeyPair()
	addr := pub.AccountAddress()
	acc := account.NewAccount(5)
	acc.AddToBalance(3)
	td.FakeState.AddTestAccount(addr, acc)

	setTx := tx.NewAnchorTx(1, addr, payload.AnchorActionSet, bytesRepeat(0x43, 32), "gone", 0, 1, 1)
	signTx(prv, setTx)

	accRes, err := chain.GetAccount(t.Context(), &pactus.GetAccountRequest{Address: addr.String()})
	require.NoError(t, err)
	require.Nil(t, accRes.Account.Anchor)
	data, err := hex.DecodeString(accRes.Account.Data)
	require.NoError(t, err)
	require.Len(t, data, 12)

	got, err := chain.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: addr.String()})
	require.NoError(t, err)
	require.False(t, got.Found)
	require.Nil(t, got.Anchor)

	_, err = txs.GetTransaction(t.Context(), &pactus.GetTransactionRequest{
		Id:        setTx.ID().String(),
		Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_INFO,
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetAnchorMissingAccount(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)
	addr := td.RandAccAddress()
	before := len(td.FakeState.FakeStore.FakeAccounts)

	_, err := client.GetAccount(t.Context(), &pactus.GetAccountRequest{Address: addr.String()})
	require.Equal(t, codes.NotFound, status.Code(err))

	res, err := client.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: addr.String()})
	require.NoError(t, err)
	require.False(t, res.Found)
	require.Nil(t, res.Anchor)
	require.Len(t, td.FakeState.FakeStore.FakeAccounts, before)
}

func TestGetAnchorRejectsBadAddress(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)
	before := len(td.FakeState.FakeStore.FakeAccounts)

	for _, address := range []string{"", "not-an-address"} {
		_, err := client.GetAccount(t.Context(), &pactus.GetAccountRequest{Address: address})
		require.Equal(t, codes.InvalidArgument, status.Code(err), address)

		_, err = client.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: address})
		require.Equal(t, codes.InvalidArgument, status.Code(err), address)
	}

	val := td.RandValAddress().String()
	_, err := client.GetAccount(t.Context(), &pactus.GetAccountRequest{Address: val})
	require.Equal(t, codes.NotFound, status.Code(err))
	_, err = client.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: val})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	treasury := account.NewAccount(0)
	treasury.AddToBalance(5)
	td.FakeState.AddTestAccount(crypto.TreasuryAddress, treasury)
	acc, err := client.GetAccount(t.Context(), &pactus.GetAccountRequest{
		Address: crypto.TreasuryAddress.String(),
	})
	require.NoError(t, err)
	require.Nil(t, acc.Account.Anchor)
	require.Equal(t, int64(5), acc.Account.Balance)
	_, err = client.GetAnchor(t.Context(), &pactus.GetAnchorRequest{
		Address: crypto.TreasuryAddress.String(),
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	require.Len(t, td.FakeState.FakeStore.FakeAccounts, before+1)
}

func TestGetAnchorFindsBothAccountCurves(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	edPub, _ := td.RandEd25519KeyPair()
	secPub, _ := td.RandSecp256k1KeyPair()
	for _, addr := range []crypto.Address{edPub.AccountAddress(), secPub.AccountAddress()} {
		acc := account.NewAccount(1)
		require.NoError(t, acc.SetAnchor(account.AnchorData{
			RootHash:      bytesRepeat(0x44, 32),
			LockedDeposit: 1,
		}))
		td.FakeState.AddTestAccount(addr, acc)

		res, err := client.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: addr.String()})
		require.NoError(t, err)
		require.True(t, res.Found)
		require.Equal(t, addr.String(), res.Address)
	}
}

func TestListAnchorsPage(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	empty, err := client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{})
	require.NoError(t, err)
	require.Equal(t, uint32(0), empty.Total)
	require.Empty(t, empty.Items)

	_, err = client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Count: 101})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	third := putAnchor(t, td, 30, 0x33)
	plainAddr, plain := td.GenerateTestAccount(testsuite.AccountWithNumber(25))
	td.FakeState.AddTestAccount(plainAddr, plain)
	first := putAnchor(t, td, 10, 0x11)
	second := putAnchor(t, td, 20, 0x22)

	page, err := client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{})
	require.NoError(t, err)
	require.Equal(t, uint32(3), page.Total)
	require.Equal(t, []string{first, second, third}, itemAddresses(page))

	zero, err := client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Count: 0})
	require.NoError(t, err)
	require.Equal(t, itemAddresses(page), itemAddresses(zero))

	one, err := client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Count: 1})
	require.NoError(t, err)
	require.Equal(t, []string{first}, itemAddresses(one))
	require.Equal(t, uint32(3), one.Total)

	next, err := client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Skip: 1, Count: 1})
	require.NoError(t, err)
	require.Equal(t, []string{second}, itemAddresses(next))
	require.Equal(t, uint32(3), next.Total)

	full, err := client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Count: 100})
	require.NoError(t, err)
	require.Len(t, full.Items, 3)

	_, err = client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Count: 101})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	for _, skip := range []uint32{3, 4, math.MaxUint32} {
		tail, err := client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Skip: skip, Count: 20})
		require.NoError(t, err)
		require.Empty(t, tail.Items)
		require.Equal(t, uint32(3), tail.Total)
	}

	again, err := client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Count: 20})
	require.NoError(t, err)
	require.Equal(t, itemAddresses(page), itemAddresses(again))

	for _, item := range page.Items {
		one, err := client.GetAnchor(t.Context(), &pactus.GetAnchorRequest{Address: item.Address})
		require.NoError(t, err)
		require.Equal(t, item.Anchor.RootHash, one.Anchor.RootHash)
	}
	require.NotEqual(t, page.Items[0].Address, page.Items[1].Address)
}

func TestListAnchorsSliceAndMutation(t *testing.T) {
	td := setup(t, nil)
	client := td.blockchainClient(t)

	addresses := make([]string, 5)
	for i := range addresses {
		addresses[i] = putAnchor(t, td, int32(i*10), byte(i+1))
	}
	for number := int32(1); number < 10; number += 2 {
		addr, acc := td.GenerateTestAccount(testsuite.AccountWithNumber(number))
		td.FakeState.AddTestAccount(addr, acc)
	}

	page, err := client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Skip: 2, Count: 2})
	require.NoError(t, err)
	require.Equal(t, uint32(5), page.Total)
	require.Equal(t, addresses[2:4], itemAddresses(page))

	short, err := client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Skip: 2, Count: 50})
	require.NoError(t, err)
	require.Equal(t, uint32(5), short.Total)
	require.Equal(t, addresses[2:], itemAddresses(short))

	gone := addresses[1]
	for addr, acc := range td.FakeState.FakeStore.FakeAccounts {
		if addr.String() == gone {
			acc.ClearAnchor()
		}
	}
	after, err := client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Count: 20})
	require.NoError(t, err)
	require.Equal(t, uint32(4), after.Total)
	require.Equal(t, []string{addresses[0], addresses[2], addresses[3], addresses[4]}, itemAddresses(after))

	for addr, acc := range td.FakeState.FakeStore.FakeAccounts {
		if addr.String() == addresses[0] {
			require.NoError(t, acc.SetAnchor(account.AnchorData{
				RootHash:      bytesRepeat(0x99, 32),
				LockedDeposit: 1,
			}))
		}
	}
	updated, err := client.ListAnchors(t.Context(), &pactus.ListAnchorsRequest{Count: 1})
	require.NoError(t, err)
	require.Equal(t, addresses[0], updated.Items[0].Address)
	require.Equal(t, bytesRepeat(0x99, 32), updated.Items[0].Anchor.RootHash)
	require.Equal(t, uint32(4), updated.Total)
}

func TestGetRawAnchorTransaction(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)
	pub, prv := td.RandBLSKeyPair()
	from := pub.AccountAddress().String()
	root := bytesRepeat(0x21, 32)

	td.FakeState.FakeHeight = 42
	res, err := client.GetRawAnchorTransaction(t.Context(), &pactus.GetRawAnchorTransactionRequest{
		From:        from,
		Action:      uint32(payload.AnchorActionSet),
		RootHash:    root,
		ManifestUri: "line\nuri",
		AnchorType:  0xFF,
		Deposit:     5,
		Fee:         7,
		Memo:        "memo\"q",
		LockTime:    9,
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.RawTransaction)
	require.Empty(t, res.Id)

	decoded := decodeAnchor(t, client, res.RawTransaction)
	require.Equal(t, pactus.PayloadType_PAYLOAD_TYPE_ANCHOR, decoded.PayloadType)
	require.Equal(t, uint32(payload.AnchorActionSet), decoded.GetAnchor().Action)
	require.Equal(t, root, decoded.GetAnchor().RootHash)
	require.Equal(t, "line\nuri", decoded.GetAnchor().ManifestUri)
	require.Equal(t, uint32(0xFF), decoded.GetAnchor().AnchorType)
	require.Equal(t, int64(5), decoded.GetAnchor().Deposit)
	require.Equal(t, int64(7), decoded.Fee)
	require.Equal(t, "memo\"q", decoded.Memo)
	require.Equal(t, uint32(9), decoded.LockTime)
	require.Empty(t, decoded.PublicKey)
	require.Empty(t, decoded.Signature)

	raw, err := tx.FromString(res.RawTransaction)
	require.NoError(t, err)
	signTx(prv, raw)
	signed, err := raw.Bytes()
	require.NoError(t, err)
	signedInfo := decodeAnchor(t, client, hex.EncodeToString(signed))
	require.NotEmpty(t, signedInfo.PublicKey)
	require.NotEmpty(t, signedInfo.Signature)
	require.Equal(t, root, signedInfo.GetAnchor().RootHash)
	require.Equal(t, int64(5), signedInfo.GetAnchor().Deposit)
	require.Equal(t, int64(7), signedInfo.Fee)
	require.NoError(t, raw.BasicCheck())

	signed[20] ^= 0x01
	flipped, flipErr := tx.FromBytes(signed)
	if flipErr == nil {
		require.Error(t, flipped.BasicCheck())
	}

	other, otherPrv := td.RandBLSKeyPair()
	mismatch, err := tx.FromString(res.RawTransaction)
	require.NoError(t, err)
	signTx(otherPrv, mismatch)
	require.NotEqual(t, other.AccountAddress(), mismatch.Payload().Signer())
	require.Error(t, mismatch.BasicCheck())

	again, err := tx.FromString(res.RawTransaction)
	require.NoError(t, err)
	signTx(prv, again)
	require.Equal(t, raw.Signature().Bytes(), again.Signature().Bytes())
}

func TestGetRawAnchorFeeAndLockTime(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)
	from := td.RandAccAddress().String()
	root := bytesRepeat(0x11, 32)
	td.FakeState.FakeHeight = 77
	td.FakeState.EXPECT().CalculateFee(amount.Amount(5), payload.TypeAnchor).Return(amount.Amount(4242))

	res, err := client.GetRawAnchorTransaction(t.Context(), &pactus.GetRawAnchorTransactionRequest{
		From:     from,
		RootHash: root,
		Deposit:  5,
	})
	require.NoError(t, err)
	decoded := decodeAnchor(t, client, res.RawTransaction)
	require.Equal(t, int64(4242), decoded.Fee)
	require.Equal(t, uint32(77), decoded.LockTime)

	kept, err := client.GetRawAnchorTransaction(t.Context(), &pactus.GetRawAnchorTransactionRequest{
		From:     from,
		RootHash: root,
		Deposit:  5,
		Fee:      1,
		LockTime: 1,
	})
	require.NoError(t, err)
	decoded = decodeAnchor(t, client, kept.RawTransaction)
	require.Equal(t, int64(1), decoded.Fee)
	require.Equal(t, uint32(1), decoded.LockTime)

	high, err := client.GetRawAnchorTransaction(t.Context(), &pactus.GetRawAnchorTransactionRequest{
		From:     from,
		RootHash: root,
		Deposit:  5,
		Fee:      amount.MaxNanoPAC,
		LockTime: math.MaxUint32,
	})
	require.NoError(t, err)
	decoded = decodeAnchor(t, client, high.RawTransaction)
	require.Equal(t, int64(amount.MaxNanoPAC), decoded.Fee)
	require.Equal(t, uint32(math.MaxUint32), decoded.LockTime)
}

func TestGetRawAnchorBounds(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)
	from := td.RandAccAddress().String()

	ok := func(req *pactus.GetRawAnchorTransactionRequest) *pactus.TransactionInfo {
		t.Helper()
		res, err := client.GetRawAnchorTransaction(t.Context(), req)
		require.NoError(t, err)

		return decodeAnchor(t, client, res.RawTransaction)
	}
	bad := func(req *pactus.GetRawAnchorTransactionRequest) {
		t.Helper()
		_, err := client.GetRawAnchorTransaction(t.Context(), req)
		require.Equal(t, codes.InvalidArgument, status.Code(err))
	}

	base := &pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), Deposit: 1, Fee: 1, LockTime: 1,
	}
	require.Equal(t, bytesRepeat(0x11, 32), ok(base).GetAnchor().RootHash)

	memo64 := strings.Repeat("m", 64)
	got := ok(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), Deposit: 1, Fee: 1, LockTime: 1, Memo: memo64,
	})
	require.Equal(t, memo64, got.Memo)
	bad(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), Deposit: 1, Fee: 1, LockTime: 1,
		Memo: strings.Repeat("m", 65),
	})

	bad(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 31), Deposit: 1, Fee: 1, LockTime: 1,
	})
	require.Len(t, ok(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), Deposit: 1, Fee: 1, LockTime: 1,
	}).GetAnchor().RootHash, 32)
	require.Len(t, ok(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 64), Deposit: 1, Fee: 1, LockTime: 1,
	}).GetAnchor().RootHash, 64)
	bad(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 65), Deposit: 1, Fee: 1, LockTime: 1,
	})

	uri128 := strings.Repeat("u", 128)
	require.Equal(t, uri128, ok(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), ManifestUri: uri128, Deposit: 1, Fee: 1, LockTime: 1,
	}).GetAnchor().ManifestUri)
	bad(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), ManifestUri: strings.Repeat("u", 129),
		Deposit: 1, Fee: 1, LockTime: 1,
	})
	res, err := client.GetRawAnchorTransaction(t.Context(), &pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), ManifestUri: "\xC3", Deposit: 1, Fee: 1, LockTime: 1,
	})
	require.Error(t, err)
	require.Nil(t, res)

	bad(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), Deposit: -1, Fee: 1, LockTime: 1,
	})
	bad(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), Deposit: 1, Fee: -1, LockTime: 1,
	})
	bad(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), Deposit: 1, Fee: amount.MaxNanoPAC + 1, LockTime: 1,
	})

	require.Equal(t, int64(0), ok(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), Fee: 1, LockTime: 1,
	}).GetAnchor().Deposit)
	require.Equal(t, int64(executor.MinAnchorDeposit-1), ok(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), Deposit: int64(executor.MinAnchorDeposit - 1),
		Fee: 1, LockTime: 1,
	}).GetAnchor().Deposit)
	require.Equal(t, int64(amount.MaxNanoPAC), ok(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), Deposit: amount.MaxNanoPAC, Fee: 1, LockTime: 1,
	}).GetAnchor().Deposit)
	bad(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), Deposit: amount.MaxNanoPAC + 1, Fee: 1, LockTime: 1,
	})

	zeros := ok(&pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: make([]byte, 32), Fee: 1, LockTime: 1,
	})
	require.Equal(t, make([]byte, 32), zeros.GetAnchor().RootHash)
	require.NotEmpty(t, zeros.GetAnchor().RootHash)

	bad(&pactus.GetRawAnchorTransactionRequest{From: from, Fee: 1, LockTime: 1})
	bad(&pactus.GetRawAnchorTransactionRequest{From: from, Action: 2, Fee: 1, LockTime: 1})
	bad(&pactus.GetRawAnchorTransactionRequest{
		From: crypto.TreasuryAddress.String(), RootHash: bytesRepeat(0x11, 32), Deposit: 1, Fee: 1, LockTime: 1,
	})
	bad(&pactus.GetRawAnchorTransactionRequest{
		From: td.RandValAddress().String(), RootHash: bytesRepeat(0x11, 32), Deposit: 1, Fee: 1, LockTime: 1,
	})

	edPub, _ := td.RandEd25519KeyPair()
	secPub, _ := td.RandSecp256k1KeyPair()
	for _, addr := range []crypto.Address{edPub.AccountAddress(), secPub.AccountAddress()} {
		info := ok(&pactus.GetRawAnchorTransactionRequest{
			From: addr.String(), RootHash: bytesRepeat(0x11, 32), Deposit: 1, Fee: 1, LockTime: 1,
		})
		require.Equal(t, addr.String(), info.GetAnchor().From)
	}
}

func TestGetRawAnchorDelete(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)
	from := td.RandAccAddress().String()

	_, err := client.GetRawAnchorTransaction(t.Context(), &pactus.GetRawAnchorTransactionRequest{
		From: from, Action: uint32(payload.AnchorActionDelete), Deposit: 1, Fee: 1, LockTime: 1,
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = client.GetRawAnchorTransaction(t.Context(), &pactus.GetRawAnchorTransactionRequest{
		From: from, Action: uint32(payload.AnchorActionDelete), RootHash: bytesRepeat(0x11, 32),
		Fee: 1, LockTime: 1,
	})
	require.Equal(t, codes.InvalidArgument, status.Code(err))

	res, err := client.GetRawAnchorTransaction(t.Context(), &pactus.GetRawAnchorTransactionRequest{
		From: from, Action: uint32(payload.AnchorActionDelete), ManifestUri: "leftover",
		AnchorType: 0xFF, Memo: "gone", Fee: 1, LockTime: 1,
	})
	require.NoError(t, err)
	decoded := decodeAnchor(t, client, res.RawTransaction)
	require.Equal(t, uint32(payload.AnchorActionDelete), decoded.GetAnchor().Action)
	require.Empty(t, decoded.GetAnchor().RootHash)
	require.Empty(t, decoded.GetAnchor().ManifestUri)
	require.Equal(t, uint32(0), decoded.GetAnchor().AnchorType)
	require.Equal(t, int64(0), decoded.GetAnchor().Deposit)
	require.Equal(t, int64(0), decoded.Value)
	require.Equal(t, "gone", decoded.Memo)
	require.NotContains(t, res.RawTransaction, hex.EncodeToString([]byte("leftover")))
}

func TestDecodeRawAnchorAndTransfer(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)
	from := td.RandAccAddress().String()
	res, err := client.GetRawAnchorTransaction(t.Context(), &pactus.GetRawAnchorTransactionRequest{
		From: from, RootHash: bytesRepeat(0x11, 32), Deposit: 1, Fee: 1, LockTime: 1,
	})
	require.NoError(t, err)
	require.Equal(t, pactus.PayloadType_PAYLOAD_TYPE_ANCHOR,
		decodeAnchor(t, client, res.RawTransaction).PayloadType)

	transfer := td.GenerateTestTransferTx()
	raw, err := transfer.Bytes()
	require.NoError(t, err)
	decoded, err := client.DecodeRawTransaction(t.Context(), &pactus.DecodeRawTransactionRequest{
		RawTransaction: hex.EncodeToString(raw),
	})
	require.NoError(t, err)
	require.Equal(t, pactus.PayloadType_PAYLOAD_TYPE_TRANSFER, decoded.Transaction.PayloadType)

	for _, badHex := range []string{"", "abc", res.RawTransaction + "ff"} {
		_, err = client.DecodeRawTransaction(t.Context(), &pactus.DecodeRawTransactionRequest{RawTransaction: badHex})
		require.Error(t, err)
	}
}

func TestCalculateFeeAnchor(t *testing.T) {
	td := setup(t, nil)
	const fixed amount.Amount = 10_000_000
	td.FakeState.EXPECT().CalculateFee(gomock.Any(), payload.TypeAnchor).Return(fixed).AnyTimes()
	client := td.transactionClient(t)

	for _, amt := range []int64{-1, 0, 1, amount.MaxNanoPAC} {
		res, err := client.CalculateFee(t.Context(), &pactus.CalculateFeeRequest{
			Amount:      amt,
			PayloadType: pactus.PayloadType_PAYLOAD_TYPE_ANCHOR,
		})
		require.NoError(t, err)
		require.Equal(t, int64(fixed), res.Fee)
	}
}

func TestBroadcastAnchor(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)
	pub, prv := td.RandBLSKeyPair()
	res, err := client.GetRawAnchorTransaction(t.Context(), &pactus.GetRawAnchorTransactionRequest{
		From: pub.AccountAddress().String(), RootHash: bytesRepeat(0x11, 32), Deposit: 1, Fee: 1, LockTime: 1,
	})
	require.NoError(t, err)

	td.FakeState.EXPECT().AddPendingTxAndBroadcast(gomock.Any()).Return(errors.New("unsigned"))
	_, err = client.BroadcastTransaction(t.Context(), &pactus.BroadcastTransactionRequest{
		SignedRawTransaction: res.RawTransaction,
	})
	require.ErrorContains(t, err, "unsigned")

	raw, err := tx.FromString(res.RawTransaction)
	require.NoError(t, err)
	signTx(prv, raw)
	signed, err := raw.Bytes()
	require.NoError(t, err)
	td.FakeState.EXPECT().AddPendingTxAndBroadcast(gomock.Any()).Return(errors.New("pool says no"))
	_, err = client.BroadcastTransaction(t.Context(), &pactus.BroadcastTransactionRequest{
		SignedRawTransaction: hex.EncodeToString(signed),
	})
	require.ErrorContains(t, err, "pool says no")
}

func TestSignRawAnchor(t *testing.T) {
	conf := testConfig()
	conf.EnableWallet = true
	td := setup(t, conf)
	wallets := td.walletClient(t)
	txs := td.transactionClient(t)
	pub, prv := td.RandBLSKeyPair()

	res, err := txs.GetRawAnchorTransaction(t.Context(), &pactus.GetRawAnchorTransactionRequest{
		From: pub.AccountAddress().String(), RootHash: bytesRepeat(0x11, 32), Deposit: 3, Fee: 4, LockTime: 1,
	})
	require.NoError(t, err)

	td.FakeWalletMgr.EXPECT().SignRawTransaction("test", "pw", gomock.Any()).DoAndReturn(
		func(_, _ string, raw []byte) ([]byte, []byte, error) {
			trx, err := tx.FromBytes(raw)
			if err != nil {
				return nil, nil, err
			}
			signTx(prv, trx)
			data, err := trx.Bytes()
			if err != nil {
				return nil, nil, err
			}

			return trx.ID().Bytes(), data, nil
		},
	).Times(2)

	first, err := wallets.SignRawTransaction(t.Context(), &pactus.SignRawTransactionRequest{
		WalletName: "test", Password: "pw", RawTransaction: res.RawTransaction,
	})
	require.NoError(t, err)
	second, err := wallets.SignRawTransaction(t.Context(), &pactus.SignRawTransactionRequest{
		WalletName: "test", Password: "pw", RawTransaction: res.RawTransaction,
	})
	require.NoError(t, err)
	require.Equal(t, first.SignedRawTransaction, second.SignedRawTransaction)

	signed := decodeAnchor(t, txs, first.SignedRawTransaction)
	require.NotEmpty(t, signed.PublicKey)
	require.NotEmpty(t, signed.Signature)
	require.Equal(t, bytesRepeat(0x11, 32), signed.GetAnchor().RootHash)
	require.Equal(t, int64(3), signed.GetAnchor().Deposit)
	require.Equal(t, int64(4), signed.Fee)
	parsed, err := tx.FromString(first.SignedRawTransaction)
	require.NoError(t, err)
	require.NoError(t, parsed.BasicCheck())
}

func putAnchor(t *testing.T, td *testData, number int32, mark byte) string {
	t.Helper()

	addr, acc := td.GenerateTestAccount(testsuite.AccountWithNumber(number))
	require.NoError(t, acc.SetAnchor(account.AnchorData{
		RootHash:      bytesRepeat(mark, 32),
		LockedDeposit: 1,
	}))
	td.FakeState.AddTestAccount(addr, acc)

	return addr.String()
}

func itemAddresses(res *pactus.ListAnchorsResponse) []string {
	out := make([]string, 0, len(res.Items))
	for _, item := range res.Items {
		out = append(out, item.Address)
	}

	return out
}

func decodeAnchor(t *testing.T, client pactus.TransactionClient, raw string) *pactus.TransactionInfo {
	t.Helper()

	res, err := client.DecodeRawTransaction(t.Context(), &pactus.DecodeRawTransactionRequest{RawTransaction: raw})
	require.NoError(t, err)

	return res.Transaction
}

func signTx(prv crypto.PrivateKey, trx *tx.Tx) {
	trx.SetSignature(prv.Sign(trx.SignBytes()))
	trx.SetPublicKey(prv.PublicKey())
}

func bytesRepeat(fill byte, n int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = fill
	}

	return out
}
