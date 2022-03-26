package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/www/capnp"
)

func sendRawTx(t *testing.T, raw []byte) error {
	res := tCapnpServer.SendRawTransaction(tCtx, func(p capnp.ZarbServer_sendRawTransaction_Params) error {
		assert.NoError(t, p.SetRawTx(raw))
		return nil
	}).Result()

	_, err := res.Struct()
	if err != nil {
		return err
	}

	return nil
}

func broadcastSendTransaction(t *testing.T, sender crypto.Signer, receiver crypto.Address, amt, fee int64) error {
	stamp := lastHash().Stamp()
	seq := getSequence(sender.Address())
	trx := tx.NewSendTx(stamp, seq+1, sender.Address(), receiver, amt, fee, "")
	sender.SignMsg(trx)

	d, _ := trx.Bytes()
	return sendRawTx(t, d)
}

func broadcastBondTransaction(t *testing.T, sender crypto.Signer, val crypto.PublicKey, stake, fee int64) error {
	stamp := lastHash().Stamp()
	seq := getSequence(sender.Address())
	trx := tx.NewBondTx(stamp, seq+1, sender.Address(), val.(*bls.PublicKey), stake, fee, "")
	sender.SignMsg(trx)

	d, _ := trx.Bytes()
	return sendRawTx(t, d)
}

func TestBondingTransactions(t *testing.T) {
	t.Run("Bonding transactions", func(t *testing.T) {
		// These validators are not in the committee now.
		// Bond transactions are valid and they can enter the committee soon
		for i := 4; i < tTotalNodes; i++ {
			amt := util.RandInt64(1000000 - 1) // fee is always 1000
			require.NoError(t, broadcastBondTransaction(t, tSigners[tNodeIdx1], tSigners[i].PublicKey(), amt, 1000))

			fmt.Printf("Staking %v to %v\n", amt, tSigners[i].Address())
			incSequence(tSigners[tNodeIdx1].Address())
		}
	})
}

func TestSendingTransactions(t *testing.T) {
	pubAlice, prvAlice := bls.GenerateTestKeyPair()
	pubBob, prvBob := bls.GenerateTestKeyPair()
	pubCarol, _ := bls.GenerateTestKeyPair()
	pubDave, _ := bls.GenerateTestKeyPair()

	signerAlice := crypto.NewSigner(prvAlice)
	signerBob := crypto.NewSigner(prvBob)

	t.Run("Sending normal transaction", func(t *testing.T) {
		require.NoError(t, broadcastSendTransaction(t, tSigners[tNodeIdx2], pubAlice.Address(), 80000000, 80000))
		incSequence(tSigners[tNodeIdx1].Address())
	})

	t.Run("Invalid fee", func(t *testing.T) {
		require.Error(t, broadcastSendTransaction(t, signerAlice, pubBob.Address(), 500000, 1))
	})

	t.Run("Alice tries double spending", func(t *testing.T) {
		require.NoError(t, broadcastSendTransaction(t, signerAlice, pubBob.Address(), 50000000, 50000))
		incSequence(signerAlice.Address())

		require.Error(t, broadcastSendTransaction(t, signerAlice, pubCarol.Address(), 50000000, 50000))
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
		for i := tTotalNodes; i < tTotalNodes; i++ {
			amt := util.RandInt64(1000000 - 1) // fee is always 1000
			require.NoError(t, broadcastBondTransaction(t, tSigners[tNodeIdx2], tSigners[i].PublicKey(), amt, 1000))

			fmt.Printf("Staking %v to %v\n", amt, tSigners[i].Address())
			incSequence(tSigners[tNodeIdx2].Address())
		}
	})

	// Make sure all transactions are confirmed
	waitForNewBlock(8)

	accAlice := getAccount(t, pubAlice.Address())
	accBob := getAccount(t, pubBob.Address())
	accCarol := getAccount(t, pubCarol.Address())
	accDave := getAccount(t, pubDave.Address())
	require.NotNil(t, accAlice)
	require.NotNil(t, accBob)
	require.NotNil(t, accCarol)
	require.NotNil(t, accDave)

	assert.Equal(t, accAlice.Balance(), int64(80000000-50050000))
	assert.Equal(t, accBob.Balance(), int64(50000000-2011))
	assert.Equal(t, accCarol.Balance(), int64(10))
	assert.Equal(t, accDave.Balance(), int64(1))
}
