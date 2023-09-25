package tests

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sendRawTx(t *testing.T, raw []byte) error {
	t.Helper()

	_, err := tTransaction.SendRawTransaction(tCtx,
		&pactus.SendRawTransactionRequest{Data: raw})
	return err
}

func broadcastSendTransaction(t *testing.T, sender *bls.ValidatorKey, receiver crypto.Address, amt, fee int64) error {
	t.Helper()

	stamp := lastHash().Stamp()
	lockTime := lastHeight() + 1
	trx := tx.NewTransferTx(stamp, lockTime, sender.PublicKey().AccountAddress(), receiver, amt, fee, "")
	sig := sender.Sign(trx.SignBytes())

	trx.SetPublicKey(sender.PublicKey())
	trx.SetSignature(sig)

	d, _ := trx.Bytes()
	return sendRawTx(t, d)
}

func broadcastBondTransaction(t *testing.T, sender *bls.ValidatorKey, pub *bls.PublicKey, stake, fee int64) error {
	t.Helper()

	stamp := lastHash().Stamp()
	lockTime := lastHeight() + 1
	trx := tx.NewBondTx(stamp, lockTime, sender.PublicKey().AccountAddress(), pub.ValidatorAddress(), pub, stake, fee, "")
	sig := sender.Sign(trx.SignBytes())

	trx.SetPublicKey(sender.PublicKey())
	trx.SetSignature(sig)

	d, _ := trx.Bytes()
	return sendRawTx(t, d)
}

func TestTransactions(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pubAlice, prvAlice := ts.RandBLSKeyPair()
	pubBob, prvBob := ts.RandBLSKeyPair()
	pubCarol, _ := ts.RandBLSKeyPair()
	pubDave, _ := ts.RandBLSKeyPair()

	signerAlice := bls.NewValidatorKey(prvAlice)
	signerBob := bls.NewValidatorKey(prvBob)

	t.Run("Sending normal transaction", func(t *testing.T) {
		require.NoError(t, broadcastSendTransaction(t, tValKeys[tNodeIdx2][0], pubAlice.AccountAddress(), 80000000, 8000))
	})

	t.Run("Invalid fee", func(t *testing.T) {
		require.Error(t, broadcastSendTransaction(t, signerAlice, pubBob.AccountAddress(), 500000, 0))
	})

	t.Run("Alice tries double spending", func(t *testing.T) {
		require.NoError(t, broadcastSendTransaction(t, signerAlice, pubBob.AccountAddress(), 50000000, 5000))

		require.Error(t, broadcastSendTransaction(t, signerAlice, pubCarol.AccountAddress(), 50000000, 5000))
	})

	t.Run("Bob sends two transaction at once", func(t *testing.T) {
		require.NoError(t, broadcastSendTransaction(t, signerBob, pubCarol.AccountAddress(), 10, 1000))

		require.NoError(t, broadcastSendTransaction(t, signerBob, pubDave.AccountAddress(), 1, 1000))
	})

	t.Run("Bonding transactions", func(t *testing.T) {
		// These validators are not in the committee now.
		// Bond transactions are valid and they can enter the committee soon
		for i := 0; i < tTotalNodes; i++ {
			amt := int64(1000000)
			fee := int64(1000)
			signer := tValKeys[tNodeIdx1][0]

			require.NoError(t, broadcastBondTransaction(t, signer, tValKeys[i][1].PublicKey(), amt, fee))
			fmt.Printf("Staking %v to %v\n", amt, tValKeys[i][1].Address())

			require.NoError(t, broadcastBondTransaction(t, signer, tValKeys[i][2].PublicKey(), amt, fee))
			fmt.Printf("Staking %v to %v\n", amt, tValKeys[i][2].Address())
		}
	})

	// Make sure all transactions are confirmed
	waitForNewBlocks(8)

	accAlice := getAccount(t, pubAlice.AccountAddress())
	accBob := getAccount(t, pubBob.AccountAddress())
	accCarol := getAccount(t, pubCarol.AccountAddress())
	accDave := getAccount(t, pubDave.AccountAddress())
	require.NotNil(t, accAlice)
	require.NotNil(t, accBob)
	require.NotNil(t, accCarol)
	require.NotNil(t, accDave)

	assert.Equal(t, accAlice.Balance, int64(80000000-50005000))
	assert.Equal(t, accBob.Balance, int64(50000000-2011))
	assert.Equal(t, accCarol.Balance, int64(10))
	assert.Equal(t, accDave.Balance, int64(1))
}
