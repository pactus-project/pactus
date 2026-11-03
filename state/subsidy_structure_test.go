package state

import (
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestSubsidyStructureIsChecked(t *testing.T) {
	td := setup(t)

	proposerAddr := td.state.Proposer(0).Address()

	t.Run("A subsidy with less than two recipients is refused without a panic", func(t *testing.T) {
		recipients := []payload.BatchRecipient{
			{To: td.RandAccAddress(), Amount: 1},
			{To: td.RandAccAddress(), Amount: 1},
		}
		trx := td.GenerateTestSubsidyTx(testsuite.TransactionWithRecipients(recipients))

		require.NotPanics(t, func() {
			// Even a malformed payload is refused gracefully.
			trx.Payload().(*payload.BatchTransferPayload).Recipients = recipients[:1]
			err := td.state.checkSubsidy(trx, proposerAddr, true)
			require.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
		})
	})

	t.Run("A block without transaction is refused without a panic", func(t *testing.T) {
		blk, _ := td.GenerateTestBlock(td.state.LastBlockHeight()+1, testsuite.BlockWithTransactions(block.Txs{}))

		sb := td.state.concreteSandbox()
		require.NotPanics(t, func() {
			err := td.state.executeBlock(blk, sb, false)
			require.ErrorIs(t, err, ErrInvalidSubsidyTransaction)
		})
	})
}
