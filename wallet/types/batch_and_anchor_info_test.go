package types_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A batch transfer gives one entry per recipient, and an anchor gives one entry without a receiver.
func TestMakeTransactionInfosOfBatchAndAnchor(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("batch transfer has one entry per recipient, with its own amount", func(t *testing.T) {
		recipients := []payload.BatchRecipient{
			{To: ts.RandAccAddress(), Amount: amount.Amount(1e9)},
			{To: ts.RandAccAddress(), Amount: amount.Amount(2e9)},
			{To: ts.RandAccAddress(), Amount: amount.Amount(3e9)},
		}
		trx := tx.NewBatchTransferTx(ts.RandHeight(), ts.RandAccAddress(), recipients, ts.RandFee())

		infos, err := types.MakeTransactionInfos(trx, types.TransactionStatusPending, 0)
		require.NoError(t, err)
		require.Len(t, infos, len(recipients))
		for i, recipient := range recipients {
			assert.Equal(t, recipient.To.String(), infos[i].Receiver)
			assert.Equal(t, recipient.Amount, infos[i].Amount)
		}
	})

	t.Run("anchor has one entry without receiver", func(t *testing.T) {
		from := ts.RandAccAddress()
		deposit := amount.Amount(5e9)
		trx := tx.NewAnchorTx(ts.RandHeight(), from, payload.AnchorActionSet,
			ts.RandBytes(32), "ipfs://manifest", 0, deposit, ts.RandFee())

		height := ts.RandHeight()
		infos, err := types.MakeTransactionInfos(trx, types.TransactionStatusConfirmed, height)
		require.NoError(t, err)
		require.Len(t, infos, 1)
		assert.Equal(t, from.String(), infos[0].Sender)
		assert.Empty(t, infos[0].Receiver)
		assert.Equal(t, deposit, infos[0].Amount)
		assert.Equal(t, payload.TypeAnchor, infos[0].PayloadType)
		assert.Equal(t, height, infos[0].BlockHeight)
	})
}
