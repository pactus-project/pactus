package grpc

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestGetTransaction(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)

	blockHeight := td.RandHeight()
	valPubKey, _ := td.RandBLSKeyPair()
	textTrx := td.GenerateTestBondTx(
		testsuite.TransactionWithValidatorPublicKey(valPubKey))
	testBlock, testCert := td.GenerateTestBlock(blockHeight,
		testsuite.BlockWithTransactions([]*tx.Tx{textTrx}))
	td.mockState.TestStore.SaveBlock(testBlock, testCert)

	t.Run("Should return transaction (verbosity: 0)", func(t *testing.T) {
		res, err := client.GetTransaction(context.Background(),
			&pactus.GetTransactionRequest{
				Id:        textTrx.ID().String(),
				Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_DATA,
			})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Transaction)
		assert.Equal(t, blockHeight, res.BlockHeight)
		assert.Equal(t, textTrx.ID().String(), res.Transaction.Id)

		b, err := textTrx.Bytes()
		assert.NoError(t, err)

		assert.Equal(t, hex.EncodeToString(b), res.Transaction.Data)
		assert.Nil(t, res.Transaction.Payload)
	})

	t.Run("Should return transaction (verbosity: 1)", func(t *testing.T) {
		res, err := client.GetTransaction(context.Background(),
			&pactus.GetTransactionRequest{
				Id:        textTrx.ID().String(),
				Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_INFO,
			})
		pld := res.Transaction.Payload.(*pactus.TransactionInfo_Bond)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Transaction)
		assert.Empty(t, res.Transaction.Data)
		assert.Equal(t, blockHeight, res.BlockHeight)
		assert.Equal(t, testBlock.Header().UnixTime(), res.BlockTime)
		assert.Equal(t, textTrx.ID().String(), res.Transaction.Id)
		assert.Equal(t, textTrx.Fee().ToNanoPAC(), res.Transaction.Fee)
		assert.Equal(t, textTrx.Memo(), res.Transaction.Memo)
		assert.Equal(t, textTrx.Payload().Type(), payload.Type(res.Transaction.PayloadType))
		assert.Equal(t, textTrx.LockTime(), res.Transaction.LockTime)
		assert.Equal(t, textTrx.Signature().String(), res.Transaction.Signature)
		assert.Equal(t, textTrx.PublicKey().String(), res.Transaction.PublicKey)
		assert.Equal(t, textTrx.Payload().(*payload.BondPayload).Stake.ToNanoPAC(), pld.Bond.Stake)
		assert.Equal(t, textTrx.Payload().(*payload.BondPayload).From.String(), pld.Bond.Sender)
		assert.Equal(t, textTrx.Payload().(*payload.BondPayload).To.String(), pld.Bond.Receiver)
		assert.Equal(t, textTrx.Payload().(*payload.BondPayload).PublicKey.String(), pld.Bond.PublicKey)
	})

	t.Run("Should return nil value because transaction id is invalid", func(t *testing.T) {
		res, err := client.GetTransaction(context.Background(),
			&pactus.GetTransactionRequest{Id: "invalid_id"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil value because transaction doesn't exist", func(t *testing.T) {
		id := td.RandHash()
		res, err := client.GetTransaction(context.Background(),
			&pactus.GetTransactionRequest{Id: id.String()})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestSendRawTransaction(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)

	t.Run("Should fail, invalid cbor", func(t *testing.T) {
		res, err := client.BroadcastTransaction(context.Background(),
			&pactus.BroadcastTransactionRequest{SignedRawTransaction: "00000000"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	trx := td.GenerateTestTransferTx()
	data, _ := trx.Bytes()
	t.Run("Should pass", func(t *testing.T) {
		res, err := client.BroadcastTransaction(context.Background(),
			&pactus.BroadcastTransactionRequest{SignedRawTransaction: hex.EncodeToString(data)})
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
	t.Run("Should fail and not broadcast", func(t *testing.T) {
		td.mockState.TestPool.AppendError = errors.New("some error")
		res, err := client.BroadcastTransaction(context.Background(),
			&pactus.BroadcastTransactionRequest{SignedRawTransaction: hex.EncodeToString(data)})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestGetRawTransaction(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)

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

		decodedTrx, err := tx.FromString(res.RawTransaction)
		assert.NoError(t, err)
		expectedLockTime := td.mockState.LastBlockHeight()
		expectedFee := td.mockState.CalculateFee(amt, payload.TypeTransfer)

		assert.Equal(t, amt, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Equal(t, expectedFee, decodedTrx.Fee())
	})

	t.Run("Batch Transfer", func(t *testing.T) {
		amt1 := td.RandAmount()
		amt2 := td.RandAmount()
		totalAmt := amt1 + amt2

		res, err := client.GetRawBatchTransferTransaction(context.Background(),
			&pactus.GetRawBatchTransferTransactionRequest{
				Sender: td.RandAccAddress().String(),
				Recipients: []*pactus.Recipient{
					{
						Receiver: td.RandAccAddress().String(),
						Amount:   amt1.ToNanoPAC(),
					},
					{
						Receiver: td.RandAccAddress().String(),
						Amount:   amt2.ToNanoPAC(),
					},
				},
				Memo: td.RandString(32),
			},
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.RawTransaction)

		decodedTrx, err := tx.FromString(res.RawTransaction)
		assert.NoError(t, err)
		expectedLockTime := td.mockState.LastBlockHeight()
		expectedFee := td.mockState.CalculateFee(totalAmt, payload.TypeBatchTransfer)

		assert.Equal(t, totalAmt, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Equal(t, expectedFee, decodedTrx.Fee())
	})

	t.Run("Bond with the Public Key", func(t *testing.T) {
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

		decodedTrx, err := tx.FromString(res.RawTransaction)
		assert.NoError(t, err)
		expectedLockTime := td.mockState.LastBlockHeight()
		expectedFee := td.mockState.CalculateFee(amt, payload.TypeBond)

		assert.Equal(t, amt, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Equal(t, expectedFee, decodedTrx.Fee())
	})

	t.Run("Bond without the Public Key", func(t *testing.T) {
		amt := td.RandAmount()

		res, err := client.GetRawBondTransaction(context.Background(),
			&pactus.GetRawBondTransactionRequest{
				Sender:   td.RandAccAddress().String(),
				Receiver: td.RandValAddress().String(),
				Stake:    amt.ToNanoPAC(),
				Memo:     td.RandString(32),
			})
		assert.NoError(t, err)
		assert.NotEmpty(t, res.RawTransaction)

		decodedTrx, err := tx.FromString(res.RawTransaction)
		assert.NoError(t, err)
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

		decodedTrx, err := tx.FromString(res.RawTransaction)
		assert.NoError(t, err)
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

		decodedTrx, err := tx.FromString(res.RawTransaction)
		assert.NoError(t, err)
		expectedLockTime := td.mockState.LastBlockHeight()
		expectedFee := td.mockState.CalculateFee(amt, payload.TypeWithdraw)

		assert.Equal(t, amt, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Equal(t, expectedFee, decodedTrx.Fee())
	})
}

func TestCalculateFee(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)

	t.Run("Not fixed amount", func(t *testing.T) {
		amt := amount.Amount(100e9)
		expectedFee := amount.Amount(0.1e9)
		res, err := client.CalculateFee(context.Background(),
			&pactus.CalculateFeeRequest{
				Amount:      amt.ToNanoPAC(),
				PayloadType: pactus.PayloadType_PAYLOAD_TYPE_TRANSFER,
				FixedAmount: false,
			})
		assert.NoError(t, err)
		assert.Equal(t, amt.ToNanoPAC(), res.Amount)
		assert.Equal(t, expectedFee.ToNanoPAC(), res.Fee)
	})

	t.Run("Fixed amount", func(t *testing.T) {
		amt := amount.Amount(100e9)
		expectedFee := amount.Amount(0.1e9)
		res, err := client.CalculateFee(context.Background(),
			&pactus.CalculateFeeRequest{
				Amount:      100e9,
				PayloadType: pactus.PayloadType_PAYLOAD_TYPE_TRANSFER,
				FixedAmount: true,
			})
		assert.NoError(t, err)
		assert.Equal(t, (amt - expectedFee).ToNanoPAC(), res.Amount)
		assert.Equal(t, expectedFee.ToNanoPAC(), res.Fee)
	})

	t.Run("Insufficient amount to pay fee", func(t *testing.T) {
		amt := amount.Amount(1)
		expectedFee := amount.Amount(0.1e9)
		res, err := client.CalculateFee(context.Background(),
			&pactus.CalculateFeeRequest{
				Amount:      amt.ToNanoPAC(),
				PayloadType: pactus.PayloadType_PAYLOAD_TYPE_TRANSFER,
				FixedAmount: true,
			})
		assert.NoError(t, err)
		assert.Negative(t, res.Amount)
		assert.Equal(t, expectedFee.ToNanoPAC(), res.Fee)
	})
}

func TestDecodeRawTransaction(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)

	t.Run("Should decode valid raw transaction", func(t *testing.T) {
		trx := td.GenerateTestTransferTx()
		data, _ := trx.Bytes()
		res, err := client.DecodeRawTransaction(context.Background(),
			&pactus.DecodeRawTransactionRequest{RawTransaction: hex.EncodeToString(data)})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, trx.ID().String(), res.Transaction.Id)
		assert.Equal(t, trx.Fee().ToNanoPAC(), res.Transaction.Fee)
		assert.Equal(t, trx.Memo(), res.Transaction.Memo)
		assert.Equal(t, trx.Payload().Type(), payload.Type(res.Transaction.PayloadType))
		assert.Equal(t, trx.LockTime(), res.Transaction.LockTime)
		assert.Equal(t, trx.Signature().String(), res.Transaction.Signature)
		assert.Equal(t, trx.PublicKey().String(), res.Transaction.PublicKey)
	})

	t.Run("Should fail to decode invalid raw transaction", func(t *testing.T) {
		res, err := client.DecodeRawTransaction(context.Background(),
			&pactus.DecodeRawTransactionRequest{RawTransaction: "invalid_raw_transaction"})
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}
