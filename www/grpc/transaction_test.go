package grpc

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	pactus "github.com/pactus-project/pactus/www/grpc/proto"
	"github.com/stretchr/testify/assert"
)

func TestGetTransaction(t *testing.T) {
	conn, client := callTransactionServer(t)

	testBlock := tMockState.TestStore.AddTestBlock(1)
	trx1 := testBlock.Transactions()[1]

	t.Run("Should return transaction", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &pactus.TransactionRequest{Id: trx1.ID().Bytes(),
			Verbosity: pactus.TransactionVerbosity_TRANSACTION_INFO})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Transaction)
		assert.Equal(t, trx1.ID().Bytes(), res.Transaction.Id)
		assert.Equal(t, trx1.Stamp().Bytes(), res.Transaction.Stamp)
		assert.Equal(t, trx1.Fee(), res.Transaction.Fee)
		assert.Equal(t, trx1.Memo(), res.Transaction.Memo)
		assert.Equal(t, trx1.Sequence(), res.Transaction.Sequence)
		assert.Equal(t, trx1.Signature().Bytes(), res.Transaction.Signature)
		assert.Equal(t, trx1.PublicKey().String(), res.Transaction.PublicKey)
		assert.Equal(t, trx1.Payload().(*payload.SendPayload).Amount, res.Transaction.Payload.(*pactus.TransactionInfo_Send).Send.Amount)
		assert.Equal(t, trx1.Payload().(*payload.SendPayload).Sender.String(), res.Transaction.Payload.(*pactus.TransactionInfo_Send).Send.Sender)
		assert.Equal(t, trx1.Payload().(*payload.SendPayload).Receiver.String(), res.Transaction.Payload.(*pactus.TransactionInfo_Send).Send.Receiver)
	})

	t.Run("Should return nil value because transaction id is invalid", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &pactus.TransactionRequest{Id: []byte("invalid_id")})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil value because transaction doesn't exist", func(t *testing.T) {
		id := hash.GenerateTestHash()
		res, err := client.GetTransaction(tCtx, &pactus.TransactionRequest{Id: id.Bytes()})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestSendRawTransaction(t *testing.T) {
	conn, client := callTransactionServer(t)

	t.Run("Should fail, invalid cbor", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &pactus.SendRawTransactionRequest{Data: []byte("00000000")})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should fail, transaction with invalid signature", func(t *testing.T) {
		trx, _ := tx.GenerateTestSendTx()
		_, signer := tx.GenerateTestSendTx()
		trx.SetSignature(signer.SignData(trx.SignBytes()))
		data, _ := trx.Bytes()
		res, err := client.SendRawTransaction(tCtx, &pactus.SendRawTransactionRequest{Data: data})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	trx, _ := tx.GenerateTestSendTx()
	data, _ := trx.Bytes()
	t.Run("Should pass", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &pactus.SendRawTransactionRequest{Data: data})
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("Should fail, Not Broadcasted", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &pactus.SendRawTransactionRequest{Data: data})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	assert.Nil(t, conn.Close(), "Error closing connection")
}
