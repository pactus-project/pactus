package types_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusAndDirectionStrings(t *testing.T) {
	assert.Equal(t, "failed", types.TransactionStatusFailed.String())
	assert.Equal(t, "pending", types.TransactionStatusPending.String())
	assert.Equal(t, "confirmed", types.TransactionStatusConfirmed.String())
	assert.Equal(t, "unknown", types.TransactionStatus(9).String())

	assert.Equal(t, "any", types.TxDirectionAny.String())
	assert.Equal(t, "incoming", types.TxDirection(types.TxDirectionIncoming).String())
	assert.Equal(t, "outgoing", types.TxDirection(types.TxDirectionOutgoing).String())
	assert.Equal(t, "unknown", types.TxDirection(9).String())
}

// A transaction with one receiver gives one entry of the history, which carries
// the fields of the transaction and its serialized bytes.
func TestMakeTransactionInfos(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	transfer := ts.GenerateTestTransferTx()
	bond := ts.GenerateTestBondTx()
	unbond := ts.GenerateTestUnbondTx()
	withdraw := ts.GenerateTestWithdrawTx()

	tests := []struct {
		name     string
		trx      *tx.Tx
		receiver string
	}{
		{"transfer", transfer, transfer.Payload().(*payload.TransferPayload).To.String()},
		{"bond", bond, bond.Payload().(*payload.BondPayload).To.String()},
		{"unbond", unbond, unbond.Payload().(*payload.UnbondPayload).Validator.String()},
		{"withdraw", withdraw, withdraw.Payload().(*payload.WithdrawPayload).To.String()},
		{"sortition", ts.GenerateTestSortitionTx(), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infos, err := types.MakeTransactionInfos(tt.trx, types.TransactionStatusConfirmed, 42)
			require.NoError(t, err)
			require.Len(t, infos, 1)

			info := infos[0]
			assert.Equal(t, tt.trx.ID().String(), info.TxID)
			assert.Equal(t, tt.trx.Payload().Signer().String(), info.Sender)
			assert.Equal(t, tt.receiver, info.Receiver)
			assert.Equal(t, tt.trx.Payload().Value(), info.Amount)
			assert.Equal(t, tt.trx.Fee(), info.Fee)
			assert.Equal(t, tt.trx.Memo(), info.Memo)
			assert.Equal(t, types.TransactionStatusConfirmed, info.Status)
			assert.EqualValues(t, 42, info.BlockHeight)
			assert.Equal(t, tt.trx.Payload().Type(), info.PayloadType)
			assert.Equal(t, types.TxDirectionAny, info.Direction)

			decoded, err := tx.FromBytes(info.Data)
			require.NoError(t, err)
			assert.Equal(t, tt.trx.ID(), decoded.ID())
		})
	}
}
