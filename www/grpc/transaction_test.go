package grpc

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestGetTransaction(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conn, client := testTransactionClient(t)

	testBlock := tMockState.TestStore.AddTestBlock(1)
	trx1 := testBlock.Transactions()[0]

	t.Run("Should return transaction", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &pactus.GetTransactionRequest{Id: trx1.ID().Bytes(),
			Verbosity: pactus.TransactionVerbosity_TRANSACTION_INFO})
		pld := res.Transaction.Payload.(*pactus.TransactionInfo_Transfer)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Transaction)
		assert.Equal(t, uint32(0x1), res.BlockHeight)
		assert.Equal(t, testBlock.Header().UnixTime(), res.BlockTime)
		assert.Equal(t, trx1.ID().Bytes(), res.Transaction.Id)
		assert.Equal(t, trx1.Stamp().Bytes(), res.Transaction.Stamp)
		assert.Equal(t, trx1.Fee(), res.Transaction.Fee)
		assert.Equal(t, trx1.Memo(), res.Transaction.Memo)
		assert.Equal(t, trx1.Payload().Type(), payload.Type(res.Transaction.PayloadType))
		assert.Equal(t, trx1.Sequence(), res.Transaction.Sequence)
		assert.Equal(t, trx1.Signature().Bytes(), res.Transaction.Signature)
		assert.Equal(t, trx1.PublicKey().String(), res.Transaction.PublicKey)
		assert.Equal(t, trx1.Payload().(*payload.TransferPayload).Amount, pld.Transfer.Amount)
		assert.Equal(t, trx1.Payload().(*payload.TransferPayload).Sender.String(), pld.Transfer.Sender)
		assert.Equal(t, trx1.Payload().(*payload.TransferPayload).Receiver.String(), pld.Transfer.Receiver)
	})

	t.Run("Should return nil value because transaction id is invalid", func(t *testing.T) {
		res, err := client.GetTransaction(tCtx, &pactus.GetTransactionRequest{Id: []byte("invalid_id")})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil value because transaction doesn't exist", func(t *testing.T) {
		id := ts.RandomHash()
		res, err := client.GetTransaction(tCtx, &pactus.GetTransactionRequest{Id: id.Bytes()})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	assert.Nil(t, conn.Close(), "Error closing connection")
}

func TestSendRawTransaction(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	conn, client := testTransactionClient(t)

	t.Run("Should fail, invalid cbor", func(t *testing.T) {
		res, err := client.SendRawTransaction(tCtx, &pactus.SendRawTransactionRequest{Data: []byte("00000000")})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should fail, transaction with invalid signature", func(t *testing.T) {
		trx, _ := ts.GenerateTestSendTx()
		_, signer := ts.GenerateTestSendTx()
		trx.SetSignature(signer.SignData(trx.SignBytes()))
		data, _ := trx.Bytes()
		res, err := client.SendRawTransaction(tCtx, &pactus.SendRawTransactionRequest{Data: data})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	trx, _ := ts.GenerateTestSendTx()
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
