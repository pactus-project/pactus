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
		p.SetRawTx(raw)
		return nil
	}).Result()

	_, err := res.Struct()
	if err != nil {
		return err
	}

	return nil
}
func TestSendTransaction(t *testing.T) {
	b := getBlockAt(t, 1)
	sender := getAccount(t, tSigners["node_1"].Address())
	receiverAddr, _, _ := crypto.GenerateTestKeyPair()
	pub := tSigners["node_1"].PublicKey()
	trx1 := tx.NewSendTx(b.Hash(), sender.Sequence()+1, sender.Address(), receiverAddr, 10000, 1000, "", &pub, nil)
	tSigners["node_1"].SignMsg(trx1)

	d, _ := trx1.Encode()
	require.NoError(t, sendRawTx(t, d))

	receiver := getAccount(t, receiverAddr)
	assert.Equal(t, receiver.Balance(), int64(10000))
}
