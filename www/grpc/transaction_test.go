package grpc

import (
	"context"
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestGetTransaction(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.transactionClient(t)

	testBlock := td.mockState.TestStore.AddTestBlock(1)
	trx1 := testBlock.Transactions()[0]

	t.Run("Should return transaction (verbosity: 0)", func(t *testing.T) {
		res, err := client.GetTransaction(context.Background(),
			&pactus.GetTransactionRequest{
				Id:        trx1.ID().Bytes(),
				Verbosity: pactus.TransactionVerbosity_TRANSACTION_DATA,
			})
		data, _ := trx1.Bytes()

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Transaction)
		assert.Equal(t, uint32(0x1), res.BlockHeight)
		assert.Equal(t, trx1.ID().Bytes(), res.Transaction.Id)
		assert.Equal(t, data, res.Transaction.Data)
		assert.Nil(t, res.Transaction.Payload)
	})

	t.Run("Should return transaction (verbosity: 1)", func(t *testing.T) {
		res, err := client.GetTransaction(context.Background(),
			&pactus.GetTransactionRequest{
				Id:        trx1.ID().Bytes(),
				Verbosity: pactus.TransactionVerbosity_TRANSACTION_INFO,
			})
		pld := res.Transaction.Payload.(*pactus.TransactionInfo_Transfer)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Transaction)
		assert.Empty(t, res.Transaction.Data)
		assert.Equal(t, uint32(0x1), res.BlockHeight)
		assert.Equal(t, testBlock.Header().UnixTime(), res.BlockTime)
		assert.Equal(t, trx1.ID().Bytes(), res.Transaction.Id)
		assert.Equal(t, trx1.Fee().ToNanoPAC(), res.Transaction.Fee)
		assert.Equal(t, trx1.Memo(), res.Transaction.Memo)
		assert.Equal(t, trx1.Payload().Type(), payload.Type(res.Transaction.PayloadType))
		assert.Equal(t, trx1.LockTime(), res.Transaction.LockTime)
		assert.Equal(t, trx1.Signature().Bytes(), res.Transaction.Signature)
		assert.Equal(t, trx1.PublicKey().String(), res.Transaction.PublicKey)
		assert.Equal(t, trx1.Payload().(*payload.TransferPayload).Amount.ToNanoPAC(), pld.Transfer.Amount)
		assert.Equal(t, trx1.Payload().(*payload.TransferPayload).From.String(), pld.Transfer.Sender)
		assert.Equal(t, trx1.Payload().(*payload.TransferPayload).To.String(), pld.Transfer.Receiver)
	})

	t.Run("Should return nil value because transaction id is invalid", func(t *testing.T) {
		res, err := client.GetTransaction(context.Background(),
			&pactus.GetTransactionRequest{Id: []byte("invalid_id")})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil value because transaction doesn't exist", func(t *testing.T) {
		id := td.RandHash()
		res, err := client.GetTransaction(context.Background(),
			&pactus.GetTransactionRequest{Id: id.Bytes()})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestSendRawTransaction(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.transactionClient(t)

	t.Run("Should fail, invalid cbor", func(t *testing.T) {
		res, err := client.BroadcastTransaction(context.Background(),
			&pactus.BroadcastTransactionRequest{SignedRawTransaction: []byte("00000000")})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	t.Run("Should fail, transaction with invalid signature", func(t *testing.T) {
		trx, _ := td.GenerateTestTransferTx()
		_, pValKey := td.GenerateTestTransferTx()
		trx.SetSignature(pValKey.Sign(trx.SignBytes()))
		data, _ := trx.Bytes()
		res, err := client.BroadcastTransaction(context.Background(),
			&pactus.BroadcastTransactionRequest{SignedRawTransaction: data})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
	trx, _ := td.GenerateTestTransferTx()
	data, _ := trx.Bytes()
	t.Run("Should pass", func(t *testing.T) {
		res, err := client.BroadcastTransaction(context.Background(),
			&pactus.BroadcastTransactionRequest{SignedRawTransaction: data})
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("Should fail, Not Broadcasted", func(t *testing.T) {
		res, err := client.BroadcastTransaction(context.Background(),
			&pactus.BroadcastTransactionRequest{SignedRawTransaction: data})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestGetRawTransaction(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.transactionClient(t)

	t.Run("Transfer", func(t *testing.T) {
		amt := td.RandAmount()
		res, err := client.GetRawTransferTransaction(context.Background(),
			&pactus.GetRawTransferTransactionRequest{
				Sender:   td.RandAccAddress().String(),
				Receiver: td.RandAccAddress().String(),
				Amount:   amt.ToNanoPAC(),
				Memo:     td.RandString(32),
			})
		assert.NoError(t, err)
		assert.NotEmpty(t, res.RawTransaction)

		decodedTrx, _ := tx.FromBytes(res.RawTransaction)
		expectedLockTime := td.mockState.LastBlockHeight()
		expectedFee := td.mockState.CalculateFee(amt, payload.TypeTransfer)

		assert.Equal(t, amt, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Equal(t, expectedFee, decodedTrx.Fee())
	})

	t.Run("Bond", func(t *testing.T) {
		amt := td.RandAmount()
		pub, _ := td.RandBLSKeyPair()

		res, err := client.GetRawBondTransaction(context.Background(),
			&pactus.GetRawBondTransactionRequest{
				Sender:    td.RandAccAddress().String(),
				Receiver:  td.RandValAddress().String(),
				Stake:     amt.ToNanoPAC(),
				PublicKey: pub.String(),
				Memo:      td.RandString(32),
			})
		assert.NoError(t, err)
		assert.NotEmpty(t, res.RawTransaction)

		decodedTrx, _ := tx.FromBytes(res.RawTransaction)
		expectedLockTime := td.mockState.LastBlockHeight()
		expectedFee := td.mockState.CalculateFee(amt, payload.TypeBond)

		assert.Equal(t, amt, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Equal(t, expectedFee, decodedTrx.Fee())
	})

	t.Run("Unbond", func(t *testing.T) {
		res, err := client.GetRawUnbondTransaction(context.Background(),
			&pactus.GetRawUnbondTransactionRequest{
				ValidatorAddress: td.RandValAddress().String(),
				Memo:             td.RandString(32),
			})
		assert.NoError(t, err)
		assert.NotEmpty(t, res.RawTransaction)

		decodedTrx, _ := tx.FromBytes(res.RawTransaction)
		expectedLockTime := td.mockState.LastBlockHeight()

		assert.Zero(t, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Zero(t, decodedTrx.Fee())
	})

	t.Run("Withdraw", func(t *testing.T) {
		amt := td.RandAmount()
		res, err := client.GetRawWithdrawTransaction(context.Background(),
			&pactus.GetRawWithdrawTransactionRequest{
				ValidatorAddress: td.RandValAddress().String(),
				AccountAddress:   td.RandAccAddress().String(),
				Amount:           amt.ToNanoPAC(),
				Memo:             td.RandString(32),
			})

		assert.NoError(t, err)
		assert.NotEmpty(t, res.RawTransaction)

		decodedTrx, _ := tx.FromBytes(res.RawTransaction)
		expectedLockTime := td.mockState.LastBlockHeight()
		expectedFee := td.mockState.CalculateFee(amt, payload.TypeWithdraw)

		assert.Equal(t, amt, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Equal(t, expectedFee, decodedTrx.Fee())
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}

func TestCalculateFee(t *testing.T) {
	td := setup(t, nil)
	conn, client := td.transactionClient(t)

	t.Run("Not fixed amount", func(t *testing.T) {
		amt := amount.Amount(100e9)
		expectedFee := td.mockState.TestParams.MaximumFee
		res, err := client.CalculateFee(context.Background(),
			&pactus.CalculateFeeRequest{
				Amount:      amt.ToNanoPAC(),
				PayloadType: pactus.PayloadType_TRANSFER_PAYLOAD,
				FixedAmount: false,
			})
		assert.NoError(t, err)
		assert.Equal(t, res.Amount, amt.ToNanoPAC())
		assert.Equal(t, res.Fee, expectedFee.ToNanoPAC())
	})

	t.Run("Fixed amount", func(t *testing.T) {
		amt := amount.Amount(100e9)
		expectedFee := td.mockState.TestParams.MaximumFee
		res, err := client.CalculateFee(context.Background(),
			&pactus.CalculateFeeRequest{
				Amount:      100e9,
				PayloadType: pactus.PayloadType_TRANSFER_PAYLOAD,
				FixedAmount: true,
			})
		assert.NoError(t, err)
		assert.Equal(t, res.Amount, (amt - expectedFee).ToNanoPAC())
		assert.Equal(t, res.Fee, expectedFee.ToNanoPAC())
	})

	t.Run("Insufficient amount to pay fee", func(t *testing.T) {
		amt := amount.Amount(1)
		expectedFee := td.mockState.TestParams.MinimumFee
		res, err := client.CalculateFee(context.Background(),
			&pactus.CalculateFeeRequest{
				Amount:      amt.ToNanoPAC(),
				PayloadType: pactus.PayloadType_TRANSFER_PAYLOAD,
				FixedAmount: true,
			})
		assert.NoError(t, err)
		assert.Negative(t, res.Amount)
		assert.Equal(t, res.Fee, expectedFee.ToNanoPAC())
	})

	assert.Nil(t, conn.Close(), "Error closing connection")
	td.StopServer()
}
