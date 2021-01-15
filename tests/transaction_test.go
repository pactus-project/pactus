package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
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

func broadcastSendTransaction(t *testing.T, sender crypto.Signer, receiver crypto.Address, amt, fee int64, expectError bool) {
	pub := sender.PublicKey()
	stamp := lastBlock(t).Hash()
	seq := getSequence(t, pub.Address())
	trx := tx.NewSendTx(stamp, seq+1, pub.Address(), receiver, amt, fee, "", &pub, nil)
	sender.SignMsg(trx)

	d, _ := trx.Encode()
	if expectError {
		require.Error(t, sendRawTx(t, d))
	} else {
		require.NoError(t, sendRawTx(t, d))
		incSequence(t, pub.Address())
	}
}

func TestSendingTransactions(t *testing.T) {
	aliceAddr, _, alicePriv := crypto.GenerateTestKeyPair()
	bobAddr, _, bobPriv := crypto.GenerateTestKeyPair()
	carolAddr, _, _ := crypto.GenerateTestKeyPair()
	daveAddr, _, _ := crypto.GenerateTestKeyPair()

	aliceSigner := crypto.NewSigner(alicePriv)
	bobSigner := crypto.NewSigner(bobPriv)

	t.Run("Sending normal transaction", func(t *testing.T) {
		broadcastSendTransaction(t, tSigners[tNodeIdx2], aliceAddr, 80000000, 80000, false)
	})

	t.Run("Invalid fee", func(t *testing.T) {
		broadcastSendTransaction(t, aliceSigner, bobAddr, 500000, 1, true)
	})

	t.Run("Alice tries double spending", func(t *testing.T) {
		broadcastSendTransaction(t, aliceSigner, bobAddr, 50000000, 50000, false)
		broadcastSendTransaction(t, aliceSigner, carolAddr, 50000000, 50000, true)
	})

	t.Run("Bob sends two transaction at once", func(t *testing.T) {
		broadcastSendTransaction(t, bobSigner, carolAddr, 10, 1000, false)
		broadcastSendTransaction(t, bobSigner, daveAddr, 1, 1000, false)
	})

	waitForNewBlock(t)
	waitForNewBlock(t)

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
