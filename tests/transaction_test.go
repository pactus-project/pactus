package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
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
	stamp := lastBlock().Hash()
	seq := getSequence(t, sender.Address())
	trx := tx.NewSendTx(stamp, seq+1, sender.Address(), receiver, amt, fee, "")
	sender.SignMsg(trx)

	d, _ := trx.Encode()
	return sendRawTx(t, d)
}

func broadcastBondTransaction(t *testing.T, sender crypto.Signer, val crypto.PublicKey, stake, fee int64) error {
	stamp := lastBlock().Hash()
	seq := getSequence(t, sender.Address())
	trx := tx.NewBondTx(stamp, seq+1, sender.Address(), val, stake, fee, "")
	sender.SignMsg(trx)

	d, _ := trx.Encode()
	return sendRawTx(t, d)
}

func TestBondingTransactions(t *testing.T) {
	t.Run("Bonding transactions", func(t *testing.T) {
		// These validators are not in the committee now.
		// Bond transactions are valid and they can enter the committee soon
		for i := tCommitteeSize; i < tTotalNodes; i++ {
			amt := util.RandInt64(1000000 - 1) // fee is always 1000
			require.NoError(t, broadcastBondTransaction(t, tSigners[tNodeIdx1], tSigners[i].PublicKey(), amt, 1000))

			fmt.Printf("Staking %v to %v\n", amt, tSigners[i].Address())
			incSequence(t, tSigners[tNodeIdx1].Address())
		}
	})
}

func TestSendingTransactions(t *testing.T) {
	aliceAddr, _, alicePriv := crypto.GenerateTestKeyPair()
	bobAddr, _, bobPriv := crypto.GenerateTestKeyPair()
	carolAddr, _, _ := crypto.GenerateTestKeyPair()
	daveAddr, _, _ := crypto.GenerateTestKeyPair()

	aliceSigner := crypto.NewSigner(alicePriv)
	bobSigner := crypto.NewSigner(bobPriv)

	t.Run("Sending normal transaction", func(t *testing.T) {
		require.NoError(t, broadcastSendTransaction(t, tSigners[tNodeIdx2], aliceAddr, 80000000, 80000))
		incSequence(t, tSigners[tNodeIdx1].Address())
	})

	t.Run("Invalid fee", func(t *testing.T) {
		require.Error(t, broadcastSendTransaction(t, aliceSigner, bobAddr, 500000, 1))
	})

	t.Run("Alice tries double spending", func(t *testing.T) {
		require.NoError(t, broadcastSendTransaction(t, aliceSigner, bobAddr, 50000000, 50000))
		incSequence(t, aliceSigner.Address())

		require.Error(t, broadcastSendTransaction(t, aliceSigner, carolAddr, 50000000, 50000))
	})

	t.Run("Bob sends two transaction at once", func(t *testing.T) {
		require.NoError(t, broadcastSendTransaction(t, bobSigner, carolAddr, 10, 1000))
		incSequence(t, bobSigner.Address())

		require.NoError(t, broadcastSendTransaction(t, bobSigner, daveAddr, 1, 1000))
		incSequence(t, bobSigner.Address())
	})

	t.Run("Bonding transactions", func(t *testing.T) {
		// These validators are not in the committee now.
		// Bond transactions are valid and they can enter the committee soon
		for i := tTotalNodes; i < tTotalNodes; i++ {
			amt := util.RandInt64(1000000 - 1) // fee is always 1000
			require.NoError(t, broadcastBondTransaction(t, tSigners[tNodeIdx2], tSigners[i].PublicKey(), amt, 1000))

			fmt.Printf("Staking %v to %v\n", amt, tSigners[i].Address())
			incSequence(t, tSigners[tNodeIdx2].Address())
		}
	})

	// Make sure all transaction confirmed
	for i := 0; i < 10; i++ {
		waitForNewBlock()
	}

	aliceAcc := getAccount(t, aliceAddr)
	bobAcc := getAccount(t, bobAddr)
	carolAcc := getAccount(t, carolAddr)
	daveAcc := getAccount(t, daveAddr)
	require.NotNil(t, aliceAcc)
	require.NotNil(t, bobAcc)
	require.NotNil(t, carolAcc)
	require.NotNil(t, daveAcc)

	assert.Equal(t, aliceAcc.Balance(), int64(80000000-50050000))
	assert.Equal(t, bobAcc.Balance(), int64(50000000-2011))
	assert.Equal(t, carolAcc.Balance(), int64(10))
	assert.Equal(t, daveAcc.Balance(), int64(1))
}
