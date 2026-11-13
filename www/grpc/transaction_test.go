package grpc_test

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetTransaction(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)

	blockHeight := td.RandHeight()
	valPubKey, _ := td.RandBLSKeyPair()
	textTrx := td.GenerateTestBondTx(
		testsuite.TransactionWithValidatorPublicKey(valPubKey),
	)

	testBlock, testCert := td.GenerateTestBlock(blockHeight, testsuite.BlockWithTransactions([]*tx.Tx{textTrx}))
	td.FakeState.AddTestBlock(testBlock, testCert)

	t.Run("Should return transaction (verbosity: 0)", func(t *testing.T) {
		res, err := client.GetTransaction(t.Context(),
			&pactus.GetTransactionRequest{
				Id:        textTrx.ID().String(),
				Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_DATA,
			})

		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Transaction)
		assert.Equal(t, uint32(blockHeight), res.BlockHeight)
		assert.Equal(t, textTrx.ID().String(), res.Transaction.Id)

		b, err := textTrx.Bytes()
		require.NoError(t, err)

		assert.Equal(t, hex.EncodeToString(b), res.Transaction.Data)
		assert.Nil(t, res.Transaction.Payload)
	})

	t.Run("Should return transaction (verbosity: 1)", func(t *testing.T) {
		res, err := client.GetTransaction(t.Context(),
			&pactus.GetTransactionRequest{
				Id:        textTrx.ID().String(),
				Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_INFO,
			})
		pld := res.Transaction.Payload.(*pactus.TransactionInfo_Bond)

		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.Transaction)
		assert.Empty(t, res.Transaction.Data)
		assert.Equal(t, uint32(blockHeight), res.BlockHeight)
		assert.Equal(t, testBlock.Header().UnixTime(), res.BlockTime)
		assert.Equal(t, textTrx.ID().String(), res.Transaction.Id)
		assert.Equal(t, textTrx.Fee().ToNanoPAC(), res.Transaction.Fee)
		assert.Equal(t, textTrx.Memo(), res.Transaction.Memo)
		assert.Equal(t, textTrx.Payload().Type(), payload.Type(res.Transaction.PayloadType))
		assert.Equal(t, uint32(textTrx.LockTime()), res.Transaction.LockTime)
		assert.Equal(t, textTrx.Signature().String(), res.Transaction.Signature)
		assert.Equal(t, textTrx.PublicKey().String(), res.Transaction.PublicKey)
		assert.Equal(t, textTrx.Payload().(*payload.BondPayload).Stake.ToNanoPAC(), pld.Bond.Stake)
		assert.Equal(t, textTrx.Payload().(*payload.BondPayload).From.String(), pld.Bond.Sender)
		assert.Equal(t, textTrx.Payload().(*payload.BondPayload).To.String(), pld.Bond.Receiver)
		assert.Equal(t, textTrx.Payload().(*payload.BondPayload).PublicKey.String(), pld.Bond.PublicKey)
	})

	t.Run("Should return nil value because transaction id is invalid", func(t *testing.T) {
		res, err := client.GetTransaction(t.Context(),
			&pactus.GetTransactionRequest{Id: "invalid_id"})
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Should return nil value because transaction doesn't exist", func(t *testing.T) {
		id := td.RandHash()
		res, err := client.GetTransaction(t.Context(),
			&pactus.GetTransactionRequest{Id: id.String()})
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestSendRawTransaction(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)

	t.Run("Should fail, invalid cbor", func(t *testing.T) {
		res, err := client.BroadcastTransaction(t.Context(),
			&pactus.BroadcastTransactionRequest{SignedRawTransaction: "00000000"})
		require.Error(t, err)
		assert.Nil(t, res)
	})

	trx := td.GenerateTestTransferTx()
	data, _ := trx.Bytes()
	t.Run("Should pass", func(t *testing.T) {
		td.FakeState.EXPECT().AddPendingTxAndBroadcast(gomock.Any()).Return(nil).Times(1)

		res, err := client.BroadcastTransaction(t.Context(),
			&pactus.BroadcastTransactionRequest{SignedRawTransaction: hex.EncodeToString(data)})
		require.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("Should fail and not broadcast", func(t *testing.T) {
		td.FakeState.EXPECT().AddPendingTxAndBroadcast(gomock.Any()).Return(errors.New("some error")).Times(1)

		res, err := client.BroadcastTransaction(t.Context(),
			&pactus.BroadcastTransactionRequest{SignedRawTransaction: hex.EncodeToString(data)})
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestGetRawTransaction(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)
	// A transaction without a lock time is invalid: the chain must have a block.
	td.FakeState.FakeHeight = td.RandHeight()

	td.FakeState.EXPECT().CalculateFee(gomock.Any(), gomock.Any()).Return(td.RandFee()).AnyTimes()
	td.FakeState.EXPECT().CalculateFeeForTx(gomock.Any()).DoAndReturn(
		func(trx *tx.Tx) (amount.Amount, error) { return trx.Fee(), nil }).AnyTimes()

	t.Run("Transfer", func(t *testing.T) {
		amt := td.RandAmount()
		res, err := client.GetRawTransferTransaction(t.Context(),
			&pactus.GetRawTransferTransactionRequest{
				Sender:   td.RandAccAddress().String(),
				Receiver: td.RandAccAddress().String(),
				Amount:   amt.ToNanoPAC(),
				Memo:     td.RandString(32),
			})
		require.NoError(t, err)
		assert.NotEmpty(t, res.RawTransaction)

		decodedTrx, err := tx.FromString(res.RawTransaction)
		require.NoError(t, err)
		expectedLockTime := td.FakeState.LastBlockHeight()
		expectedFee := td.FakeState.CalculateFee(amt, payload.TypeTransfer)

		assert.Equal(t, amt, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Equal(t, expectedFee, decodedTrx.Fee())
	})

	t.Run("Batch Transfer", func(t *testing.T) {
		amt1 := td.RandAmount()
		amt2 := td.RandAmount()
		totalAmt := amt1 + amt2

		res, err := client.GetRawBatchTransferTransaction(
			t.Context(),
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
		require.NoError(t, err)
		assert.NotEmpty(t, res.RawTransaction)

		decodedTrx, err := tx.FromString(res.RawTransaction)
		require.NoError(t, err)
		expectedLockTime := td.FakeState.LastBlockHeight()
		expectedFee := td.FakeState.CalculateFee(totalAmt, payload.TypeBatchTransfer)

		assert.Equal(t, totalAmt, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Equal(t, expectedFee, decodedTrx.Fee())
	})

	t.Run("Bond with the Public Key", func(t *testing.T) {
		amt := td.RandAmount()
		pub, _ := td.RandBLSKeyPair()

		res, err := client.GetRawBondTransaction(t.Context(),
			&pactus.GetRawBondTransactionRequest{
				Sender:    td.RandAccAddress().String(),
				Receiver:  pub.ValidatorAddress().String(),
				Stake:     amt.ToNanoPAC(),
				PublicKey: pub.String(),
				Memo:      td.RandString(32),
			})
		require.NoError(t, err)
		assert.NotEmpty(t, res.RawTransaction)

		decodedTrx, err := tx.FromString(res.RawTransaction)
		require.NoError(t, err)
		expectedLockTime := td.FakeState.LastBlockHeight()
		expectedFee := td.FakeState.CalculateFee(amt, payload.TypeBond)

		assert.Equal(t, amt, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Equal(t, expectedFee, decodedTrx.Fee())
	})

	t.Run("Bond without the Public Key", func(t *testing.T) {
		amt := td.RandAmount()

		res, err := client.GetRawBondTransaction(t.Context(),
			&pactus.GetRawBondTransactionRequest{
				Sender:   td.RandAccAddress().String(),
				Receiver: td.RandValAddress().String(),
				Stake:    amt.ToNanoPAC(),
				Memo:     td.RandString(32),
			})
		require.NoError(t, err)
		assert.NotEmpty(t, res.RawTransaction)

		decodedTrx, err := tx.FromString(res.RawTransaction)
		require.NoError(t, err)
		expectedLockTime := td.FakeState.LastBlockHeight()
		expectedFee := td.FakeState.CalculateFee(amt, payload.TypeBond)

		assert.Equal(t, amt, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Equal(t, expectedFee, decodedTrx.Fee())
	})

	t.Run("Unbond", func(t *testing.T) {
		res, err := client.GetRawUnbondTransaction(t.Context(),
			&pactus.GetRawUnbondTransactionRequest{
				ValidatorAddress: td.RandValAddress().String(),
				Memo:             td.RandString(32),
			})
		require.NoError(t, err)
		assert.NotEmpty(t, res.RawTransaction)

		decodedTrx, err := tx.FromString(res.RawTransaction)
		require.NoError(t, err)
		expectedLockTime := td.FakeState.LastBlockHeight()

		assert.Zero(t, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Zero(t, decodedTrx.Fee())
	})

	t.Run("Withdraw", func(t *testing.T) {
		amt := td.RandAmount()
		res, err := client.GetRawWithdrawTransaction(t.Context(),
			&pactus.GetRawWithdrawTransactionRequest{
				ValidatorAddress: td.RandValAddress().String(),
				AccountAddress:   td.RandAccAddress().String(),
				Amount:           amt.ToNanoPAC(),
				Memo:             td.RandString(32),
			})

		require.NoError(t, err)
		assert.NotEmpty(t, res.RawTransaction)

		decodedTrx, err := tx.FromString(res.RawTransaction)
		require.NoError(t, err)
		expectedLockTime := td.FakeState.LastBlockHeight()
		expectedFee := td.FakeState.CalculateFee(amt, payload.TypeWithdraw)

		assert.Equal(t, amt, decodedTrx.Payload().Value())
		assert.Equal(t, expectedLockTime, decodedTrx.LockTime())
		assert.Equal(t, expectedFee, decodedTrx.Fee())
	})
}

// TestGetRawTransactionFollowsMinimumFee checks that a raw transaction built with the
// estimated fee is rebuilt when the pool requires more for this particular transaction,
// and that an explicit fee is used as it is.
func TestGetRawTransactionFollowsMinimumFee(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)
	td.FakeState.FakeHeight = td.RandHeight()

	const estimated = amount.Amount(1_000)
	const required = amount.Amount(5_000)

	td.FakeState.EXPECT().CalculateFee(gomock.Any(), gomock.Any()).Return(estimated).AnyTimes()
	td.FakeState.EXPECT().CalculateFeeForTx(gomock.Any()).DoAndReturn(
		func(*tx.Tx) (amount.Amount, error) { return required, nil }).AnyTimes()

	sender := td.RandAccAddress().String()
	root := make([]byte, 32)

	tests := []struct {
		name  string
		build func(fee int64) (string, error)
	}{
		{"Transfer", func(fee int64) (string, error) {
			res, err := client.GetRawTransferTransaction(t.Context(), &pactus.GetRawTransferTransactionRequest{
				Sender: sender, Receiver: td.RandAccAddress().String(), Amount: 10, Fee: fee,
			})

			return rawOf(res, err)
		}},
		{"Bond", func(fee int64) (string, error) {
			res, err := client.GetRawBondTransaction(t.Context(), &pactus.GetRawBondTransactionRequest{
				Sender: sender, Receiver: td.RandValAddress().String(), Stake: 10, Fee: fee,
			})

			return rawOf(res, err)
		}},
		{"Withdraw", func(fee int64) (string, error) {
			res, err := client.GetRawWithdrawTransaction(t.Context(), &pactus.GetRawWithdrawTransactionRequest{
				ValidatorAddress: td.RandValAddress().String(), AccountAddress: sender, Amount: 10, Fee: fee,
			})

			return rawOf(res, err)
		}},
		{"Batch transfer", func(fee int64) (string, error) {
			res, err := client.GetRawBatchTransferTransaction(t.Context(), &pactus.GetRawBatchTransferTransactionRequest{
				Sender: sender, Fee: fee,
				// A batch transfer has between 2 and 8 recipients.
				Recipients: []*pactus.Recipient{
					{Receiver: td.RandAccAddress().String(), Amount: 10},
					{Receiver: td.RandAccAddress().String(), Amount: 10},
				},
			})

			return rawOf(res, err)
		}},
		{"Anchor", func(fee int64) (string, error) {
			res, err := client.GetRawAnchorTransaction(t.Context(), &pactus.GetRawAnchorTransactionRequest{
				Sender: sender, Action: uint32(payload.AnchorActionSet), RootHash: root, Deposit: 5, Fee: fee,
			})

			return rawOf(res, err)
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name+": estimated fee follows the minimum", func(t *testing.T) {
			raw, err := tt.build(0)
			require.NoError(t, err)

			trx, err := tx.FromString(raw)
			require.NoError(t, err)
			assert.Equal(t, required, trx.Fee())
		})

		t.Run(tt.name+": explicit fee is kept", func(t *testing.T) {
			raw, err := tt.build(7)
			require.NoError(t, err)

			trx, err := tx.FromString(raw)
			require.NoError(t, err)
			assert.Equal(t, amount.Amount(7), trx.Fee())
		})
	}
}

func rawOf(res *pactus.GetRawTransactionResponse, err error) (string, error) {
	if err != nil {
		return "", err
	}

	return res.RawTransaction, nil
}

func TestCalculateFee(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)

	expectedFee := td.RandFee()
	td.FakeState.EXPECT().CalculateFee(gomock.Any(), gomock.Any()).Return(expectedFee).AnyTimes()

	t.Run("Not fixed amount", func(t *testing.T) {
		amt := td.RandAmount()
		res, err := client.CalculateFee(t.Context(),
			&pactus.CalculateFeeRequest{
				Amount:      amt.ToNanoPAC(),
				PayloadType: pactus.PayloadType_PAYLOAD_TYPE_TRANSFER,
				FixedAmount: false,
			})
		require.NoError(t, err)
		assert.Equal(t, amt.ToNanoPAC(), res.Amount)
		assert.Equal(t, expectedFee.ToNanoPAC(), res.Fee)
	})

	t.Run("Fixed amount", func(t *testing.T) {
		amt := td.RandAmount()
		res, err := client.CalculateFee(t.Context(),
			&pactus.CalculateFeeRequest{
				Amount:      amt.ToNanoPAC(),
				PayloadType: pactus.PayloadType_PAYLOAD_TYPE_TRANSFER,
				FixedAmount: true,
			})
		require.NoError(t, err)
		assert.Equal(t, (amt - expectedFee).ToNanoPAC(), res.Amount)
		assert.Equal(t, expectedFee.ToNanoPAC(), res.Fee)
	})

	t.Run("Insufficient amount to pay fee", func(t *testing.T) {
		amt := amount.Amount(0)
		res, err := client.CalculateFee(t.Context(),
			&pactus.CalculateFeeRequest{
				Amount:      amt.ToNanoPAC(),
				PayloadType: pactus.PayloadType_PAYLOAD_TYPE_TRANSFER,
				FixedAmount: true,
			})
		require.NoError(t, err)
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
		res, err := client.DecodeRawTransaction(t.Context(),
			&pactus.DecodeRawTransactionRequest{RawTransaction: hex.EncodeToString(data)})
		require.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, trx.ID().String(), res.Transaction.Id)
		assert.Equal(t, trx.Fee().ToNanoPAC(), res.Transaction.Fee)
		assert.Equal(t, trx.Memo(), res.Transaction.Memo)
		assert.Equal(t, trx.Payload().Type(), payload.Type(res.Transaction.PayloadType))
		assert.Equal(t, uint32(trx.LockTime()), res.Transaction.LockTime)
		assert.Equal(t, trx.Signature().String(), res.Transaction.Signature)
		assert.Equal(t, trx.PublicKey().String(), res.Transaction.PublicKey)
	})

	t.Run("Should fail to decode invalid raw transaction", func(t *testing.T) {
		res, err := client.DecodeRawTransaction(t.Context(),
			&pactus.DecodeRawTransactionRequest{RawTransaction: "invalid_raw_transaction"})
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestCheckTransaction(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)

	trx := td.GenerateTestTransferTx()
	data, _ := trx.Bytes()
	req := &pactus.CheckTransactionRequest{RawTransaction: hex.EncodeToString(data)}

	t.Run("a transaction that can't be decoded is an error", func(t *testing.T) {
		_, err := client.CheckTransaction(t.Context(), &pactus.CheckTransactionRequest{RawTransaction: "00"})
		require.Error(t, err)
	})

	t.Run("a valid transaction", func(t *testing.T) {
		td.FakeState.EXPECT().CheckTransaction(gomock.Any()).Return(nil).Times(1)

		res, err := client.CheckTransaction(t.Context(), req)
		require.NoError(t, err)
		assert.True(t, res.IsValid)
	})

	t.Run("an invalid transaction is in the answer, not an error", func(t *testing.T) {
		td.FakeState.EXPECT().CheckTransaction(gomock.Any()).Return(errors.New("insufficient funds")).Times(1)

		res, err := client.CheckTransaction(t.Context(), req)
		require.NoError(t, err)
		assert.False(t, res.IsValid)
		assert.Equal(t, "insufficient funds", res.ErrorMessage)
	})
}

// A transaction that the node would refuse once it is signed is reported when it is built.
func TestGetRawTransactionRejectsInvalid(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)
	td.FakeState.FakeHeight = td.RandHeight()
	td.FakeState.EXPECT().CalculateFee(gomock.Any(), gomock.Any()).Return(td.RandFee()).AnyTimes()
	td.FakeState.EXPECT().CalculateFeeForTx(gomock.Any()).DoAndReturn(
		func(trx *tx.Tx) (amount.Amount, error) { return trx.Fee(), nil }).AnyTimes()

	tests := map[string]*pactus.GetRawTransferTransactionRequest{
		"negative amount": {Sender: td.RandAccAddress().String(), Receiver: td.RandAccAddress().String(), Amount: -1},
		"negative fee": {
			Sender: td.RandAccAddress().String(), Receiver: td.RandAccAddress().String(), Amount: 1, Fee: -1,
		},
		"memo too long": {
			Sender: td.RandAccAddress().String(), Receiver: td.RandAccAddress().String(), Amount: 1,
			Memo: strings.Repeat("m", 65),
		},
	}

	for name, req := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := client.GetRawTransferTransaction(t.Context(), req)
			require.Equal(t, codes.InvalidArgument, status.Code(err))
		})
	}
}
