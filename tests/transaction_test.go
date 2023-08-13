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

func sendRawTx(_ *testing.T, raw []byte) error {
	_, err := tTransaction.SendRawTransaction(tCtx,
		&pactus.SendRawTransactionRequest{Data: raw})
	return err
}

func broadcastSendTransaction(t *testing.T, sender crypto.Signer, receiver crypto.Address, amt, fee int64) error {
	stamp := lastHash().Stamp()
	seq := getSequence(sender.Address())
	trx := tx.NewTransferTx(stamp, seq+1, sender.Address(), receiver, amt, fee, "")
	sender.SignMsg(trx)

	d, _ := trx.Bytes()
	return sendRawTx(t, d)
}

func broadcastBondTransaction(t *testing.T, sender crypto.Signer, pub crypto.PublicKey, stake, fee int64) error {
	stamp := lastHash().Stamp()
	seq := getSequence(sender.Address())
	trx := tx.NewBondTx(stamp, seq+1, sender.Address(), pub.Address(), pub.(*bls.PublicKey), stake, fee, "")
	sender.SignMsg(trx)

	d, _ := trx.Bytes()
	return sendRawTx(t, d)
}

func TestTransactions(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pubAlice, prvAlice := ts.RandomBLSKeyPair()
	pubBob, prvBob := ts.RandomBLSKeyPair()
	pubCarol, _ := ts.RandomBLSKeyPair()
	pubDave, _ := ts.RandomBLSKeyPair()

	signerAlice := crypto.NewSigner(prvAlice)
	signerBob := crypto.NewSigner(prvBob)

	t.Run("Sending normal transaction", func(t *testing.T) {
		require.NoError(t, broadcastSendTransaction(t, tSigners[tNodeIdx2][0], pubAlice.Address(), 80000000, 8000))
		incSequence(tSigners[tNodeIdx2][0].Address())
	})

	t.Run("Invalid fee", func(t *testing.T) {
		require.Error(t, broadcastSendTransaction(t, signerAlice, pubBob.Address(), 500000, 0))
	})

	t.Run("Alice tries double spending", func(t *testing.T) {
		require.NoError(t, broadcastSendTransaction(t, signerAlice, pubBob.Address(), 50000000, 5000))
		incSequence(signerAlice.Address())

		require.Error(t, broadcastSendTransaction(t, signerAlice, pubCarol.Address(), 50000000, 5000))
	})

	t.Run("Bob sends two transaction at once", func(t *testing.T) {
		require.NoError(t, broadcastSendTransaction(t, signerBob, pubCarol.Address(), 10, 1000))
		incSequence(signerBob.Address())

		require.NoError(t, broadcastSendTransaction(t, signerBob, pubDave.Address(), 1, 1000))
		incSequence(signerBob.Address())
	})

	t.Run("Bonding transactions", func(t *testing.T) {
		// These validators are not in the committee now.
		// Bond transactions are valid and they can enter the committee soon
		for i := 0; i < tTotalNodes; i++ {
			amt := int64(1000000)
			fee := int64(1000)
			signer := tSigners[tNodeIdx1][0]

			require.NoError(t, broadcastBondTransaction(t, signer, tSigners[i][1].PublicKey(), amt, fee))
			fmt.Printf("Staking %v to %v\n", amt, tSigners[i][1].Address())
			incSequence(signer.Address())

			require.NoError(t, broadcastBondTransaction(t, signer, tSigners[i][2].PublicKey(), amt, fee))
			fmt.Printf("Staking %v to %v\n", amt, tSigners[i][2].Address())
			incSequence(signer.Address())
		}
	})

	// Make sure all transactions are confirmed
	waitForNewBlocks(8)

	accAlice := getAccount(t, pubAlice.Address())
	accBob := getAccount(t, pubBob.Address())
	accCarol := getAccount(t, pubCarol.Address())
	accDave := getAccount(t, pubDave.Address())
	require.NotNil(t, accAlice)
	require.NotNil(t, accBob)
	require.NotNil(t, accCarol)
	require.NotNil(t, accDave)

	assert.Equal(t, accAlice.Balance, int64(80000000-50005000))
	assert.Equal(t, accBob.Balance, int64(50000000-2011))
	assert.Equal(t, accCarol.Balance, int64(10))
	assert.Equal(t, accDave.Balance, int64(1))
}
