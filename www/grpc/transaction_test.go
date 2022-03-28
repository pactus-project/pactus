package grpc

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetTransaction(t *testing.T) {
	conn, client := callServer(t)

	tx1 := tMockState.TestStore.AddTestTransaction()

	t.Run("Should return transaction", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &zarb.TransactionRequest{Id: tx1.ID().String()})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Tranaction)
		assert.Equal(t, tx1.ID().Bytes(), res.Tranaction.Id)
		assert.Equal(t, tx1.Stamp().Bytes(), res.Tranaction.Stamp)
		assert.Equal(t, tx1.Fee(), res.Tranaction.Fee)
		assert.Equal(t, tx1.Memo(), res.Tranaction.Memo)
		assert.Equal(t, tx1.Sequence(), res.Tranaction.Sequence)
		assert.Equal(t, tx1.Signature().Bytes(), res.Tranaction.Signature)
		assert.Equal(t, tx1.Payload().(*payload.SendPayload).Amount, res.Tranaction.Payload.(*zarb.TransactionInfo_Send).Send.Amount)
		assert.Equal(t, tx1.Payload().(*payload.SendPayload).Sender.String(), res.Tranaction.Payload.(*zarb.TransactionInfo_Send).Send.Sender)
		assert.Equal(t, tx1.Payload().(*payload.SendPayload).Receiver.String(), res.Tranaction.Payload.(*zarb.TransactionInfo_Send).Send.Receiver)
	})

	t.Run("Should return nil value because transcation id is invalid", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &zarb.TransactionRequest{Id: "invalid_id"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil value because transcation doesn't exist", func(t *testing.T) {
		id := hash.GenerateTestHash()
		res, err := client.GetTransaction(tCtx, &zarb.TransactionRequest{Id: id.String()})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	assert.Nil(t, conn.Close(), "Error closing connection")

}

func TestSendRawTransaction(t *testing.T) {
	conn, client := callServer(t)

	t.Run("Should fail, invalid raw data", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &zarb.SendRawTransactionRequest{Data: "invalid raw data"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should fail, invalid cbor", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &zarb.SendRawTransactionRequest{Data: hex.EncodeToString([]byte("00000000"))})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should fail, transaction with invalid signature", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		_, signer := tx.GenerateTestSendTx()
		trx.SetSignature(signer.SignData(trx.SignBytes()))
		data, _ := trx.Bytes()
		res, err := client.SendRawTransaction(tCtx, &zarb.SendRawTransactionRequest{Data: hex.EncodeToString(data)})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	trx, _ := tx.GenerateTestSendTx()
	data, _ := trx.Bytes()
	t.Run("Should pass", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &zarb.SendRawTransactionRequest{Data: hex.EncodeToString(data)})
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("Should fail, Not Broadcasted", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &zarb.SendRawTransactionRequest{Data: hex.EncodeToString(data)})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	assert.Nil(t, conn.Close(), "Error closing connection")
}
