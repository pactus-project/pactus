package grpc

import (
	"encoding/hex"
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	pactus "github.com/pactus-project/pactus/www/grpc/proto"
	"github.com/stretchr/testify/assert"
)

func TestGetTransaction(t *testing.T) {
	conn, client := callTransactionServer(t)

	tx1 := tMockState.TestStore.AddTestTransaction()

	t.Run("Should return transaction", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &pactus.TransactionRequest{Id: tx1.ID().String()})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Transaction)
		assert.Equal(t, tx1.ID().Bytes(), res.Transaction.Id)
		assert.Equal(t, tx1.Stamp().Bytes(), res.Transaction.Stamp)
		assert.Equal(t, tx1.Fee(), res.Transaction.Fee)
		assert.Equal(t, tx1.Memo(), res.Transaction.Memo)
		assert.Equal(t, tx1.Sequence(), res.Transaction.Sequence)
		assert.Equal(t, tx1.Signature().Bytes(), res.Transaction.Signature)
		assert.Equal(t, tx1.PublicKey().String(), res.Transaction.PublicKey)
		assert.Equal(t, tx1.Payload().(*payload.SendPayload).Amount, res.Transaction.Payload.(*pactus.TransactionInfo_Send).Send.Amount)
		assert.Equal(t, tx1.Payload().(*payload.SendPayload).Sender.String(), res.Transaction.Payload.(*pactus.TransactionInfo_Send).Send.Sender)
		assert.Equal(t, tx1.Payload().(*payload.SendPayload).Receiver.String(), res.Transaction.Payload.(*pactus.TransactionInfo_Send).Send.Receiver)
	})

	t.Run("Should return nil value because transaction id is invalid", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &pactus.TransactionRequest{Id: "invalid_id"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil value because transaction doesn't exist", func(t *testing.T) {
		id := hash.GenerateTestHash()
		res, err := client.GetTransaction(tCtx, &pactus.TransactionRequest{Id: id.String()})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestSendRawTransaction(t *testing.T) {
	conn, client := callTransactionServer(t)

	t.Run("Should fail, invalid raw data", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &pactus.SendRawTransactionRequest{Data: "invalid raw data"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should fail, invalid cbor", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &pactus.SendRawTransactionRequest{Data: hex.EncodeToString([]byte("00000000"))})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should fail, transaction with invalid signature", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		_, signer := tx.GenerateTestSendTx()
		trx.SetSignature(signer.SignData(trx.SignBytes()))
		data, _ := trx.Bytes()
		res, err := client.SendRawTransaction(tCtx, &pactus.SendRawTransactionRequest{Data: hex.EncodeToString(data)})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	trx, _ := tx.GenerateTestSendTx()
	data, _ := trx.Bytes()
	t.Run("Should pass", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &pactus.SendRawTransactionRequest{Data: hex.EncodeToString(data)})
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("Should fail, Not Broadcasted", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &pactus.SendRawTransactionRequest{Data: hex.EncodeToString(data)})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	assert.Nil(t, conn.Close(), "Error closing connection")
}
