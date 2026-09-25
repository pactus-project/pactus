package integration

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/require"
)

// waitFor polls check until it succeeds or about 40 seconds (20 blocks) pass.
func waitFor(t *testing.T, what string, check func() bool) {
	t.Helper()

	for i := 0; i < 40; i++ {
		if check() {
			return
		}
		time.Sleep(time.Second)
	}
	require.Fail(t, "timeout waiting for "+what)
}

// broadcastAnchor builds an anchor transaction with GetRawAnchorTransaction, as a
// client would, signs it with the owner key and broadcasts it.
func broadcastAnchor(t *testing.T, owner *bls.ValidatorKey, req *pactus.GetRawAnchorTransactionRequest) *tx.Tx {
	t.Helper()

	res, err := tTransactionClient.GetRawAnchorTransaction(tCtx, req)
	require.NoError(t, err)
	raw, err := hex.DecodeString(res.RawTransaction)
	require.NoError(t, err)
	trx, err := tx.FromBytes(raw)
	require.NoError(t, err)

	trx.SetPublicKey(owner.PublicKey())
	trx.SetSignature(owner.Sign(trx.SignBytes()))
	signed, err := trx.Bytes()
	require.NoError(t, err)
	require.NoError(t, sendRawTx(t, signed))

	return trx
}

func getAnchor(t *testing.T, addr crypto.Address) *pactus.GetAnchorResponse {
	t.Helper()

	res, err := tBlockchainClient.GetAnchor(tCtx, &pactus.GetAnchorRequest{Address: addr.String()})
	require.NoError(t, err)

	return res
}

// requireSameOnAllNodes waits until every node's state passes check for the
// account, then requires the account bytes to be identical on all nodes.
func requireSameOnAllNodes(t *testing.T, addr crypto.Address, check func(*account.Account) bool) {
	t.Helper()

	var first []byte
	for i, node := range tNodes {
		waitFor(t, fmt.Sprintf("node %d", i+1), func() bool {
			acc, err := node.State().AccountByAddress(addr)

			return err == nil && check(acc)
		})

		acc, err := node.State().AccountByAddress(addr)
		require.NoError(t, err)
		raw, err := acc.Bytes()
		require.NoError(t, err)
		if first == nil {
			first = raw

			continue
		}
		require.Equal(t, first, raw, "node %d disagrees with node 1", i+1)
	}
}

// An anchor goes the whole way a client sees: built by GetRawAnchorTransaction,
// signed, broadcast over gRPC, gossiped, included by consensus, then read back
// identically from every node.
func TestAnchors(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	_, prv := ts.RandBLSKeyPair()
	owner := bls.NewValidatorKey(prv)
	ownerAddr := owner.PublicKey().AccountAddress()
	root := ts.RandBytes(32)
	deposit := amount.Amount(1e9)

	// The genesis account belongs to the last key of the last node (see TestMain).
	funder := tValKeys[tNodeIdx4][2]
	funderAcc := getAccount(t, funder.PublicKey().AccountAddress())
	require.NotNil(t, funderAcc)
	require.GreaterOrEqual(t, funderAcc.Balance, int64(5e9))
	require.NoError(t, broadcastSendTransaction(t, funder, ownerAddr, 5e9, 1000))
	waitFor(t, "funding", func() bool {
		acc := getAccount(t, ownerAddr)

		return acc != nil && acc.Balance == 5e9
	})

	var created *pactus.AnchorInfo
	t.Run("create", func(t *testing.T) {
		trx := broadcastAnchor(t, owner, &pactus.GetRawAnchorTransactionRequest{
			From:        ownerAddr.String(),
			Action:      uint32(payload.AnchorActionSet),
			RootHash:    root,
			ManifestUri: "ipfs://integration",
			AnchorType:  2,
			Deposit:     deposit.ToNanoPAC(),
		})
		waitFor(t, "anchor creation", func() bool {
			return getAnchor(t, ownerAddr).Found
		})

		created = getAnchor(t, ownerAddr).Anchor
		require.Equal(t, root, created.RootHash)
		require.Equal(t, "ipfs://integration", created.ManifestUri)
		require.Equal(t, uint32(2), created.AnchorType)
		require.Equal(t, deposit.ToNanoPAC(), created.LockedDeposit)
		require.Equal(t, created.CreatedAtHeight, created.UpdatedAtHeight)
		require.Equal(t, created.CreatedAtTime, created.UpdatedAtTime)

		// The clocks are those of the block that includes the transaction.
		included, err := tTransactionClient.GetTransaction(tCtx, &pactus.GetTransactionRequest{
			Id:        trx.ID().String(),
			Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_INFO,
		})
		require.NoError(t, err)
		require.Equal(t, included.BlockHeight, created.CreatedAtHeight)
		require.Equal(t, included.BlockTime, created.CreatedAtTime)

		requireSameOnAllNodes(t, ownerAddr, func(acc *account.Account) bool {
			return acc.HasAnchor()
		})
	})

	t.Run("a top-up keeps the content date", func(t *testing.T) {
		broadcastAnchor(t, owner, &pactus.GetRawAnchorTransactionRequest{
			From:        ownerAddr.String(),
			Action:      uint32(payload.AnchorActionSet),
			RootHash:    root,
			ManifestUri: "ipfs://integration",
			AnchorType:  2,
			Deposit:     deposit.ToNanoPAC(),
		})
		waitFor(t, "top-up", func() bool {
			res := getAnchor(t, ownerAddr)

			return res.Found && res.Anchor.LockedDeposit == 2*deposit.ToNanoPAC()
		})

		toppedUp := getAnchor(t, ownerAddr).Anchor
		require.Equal(t, created.CreatedAtHeight, toppedUp.CreatedAtHeight)
		require.Equal(t, created.UpdatedAtHeight, toppedUp.UpdatedAtHeight)
		require.Equal(t, created.UpdatedAtTime, toppedUp.UpdatedAtTime)

		requireSameOnAllNodes(t, ownerAddr, func(acc *account.Account) bool {
			return acc.LockedDeposit() == 2*deposit
		})
	})

	t.Run("delete refunds the deposit", func(t *testing.T) {
		before := getAccount(t, ownerAddr).Balance
		trx := broadcastAnchor(t, owner, &pactus.GetRawAnchorTransactionRequest{
			From:   ownerAddr.String(),
			Action: uint32(payload.AnchorActionDelete),
		})
		waitFor(t, "delete", func() bool {
			return !getAnchor(t, ownerAddr).Found
		})

		after := getAccount(t, ownerAddr)
		require.Equal(t, before+2*deposit.ToNanoPAC()-trx.Fee().ToNanoPAC(), after.Balance)
		require.Nil(t, after.Anchor)

		requireSameOnAllNodes(t, ownerAddr, func(acc *account.Account) bool {
			raw, err := acc.Bytes()

			return err == nil && !acc.HasAnchor() && len(raw) == 12
		})
	})
}
