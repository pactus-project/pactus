package txpool

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/stretchr/testify/require"
)

func TestAnchorRejectedByPool(t *testing.T) {
	td := setup(t, nil)

	previous := executor.DefaultFactory
	t.Cleanup(func() {
		executor.DefaultFactory = previous
	})
	mockExe := td.exe
	executor.DefaultFactory = func(trx *tx.Tx, sbx sandbox.Sandbox) (executor.Executor, error) {
		if trx.Payload().Type() == payload.TypeAnchor {
			return executor.MakeExecutorImpl(trx, sbx)
		}

		return mockExe, nil
	}

	pub, prv := td.RandBLSKeyPair()
	from := pub.AccountAddress()
	root := bytes.Repeat([]byte{0x11}, 32)
	setTx := tx.NewAnchorTx(td.sbx.CurrentHeight(), from, payload.AnchorActionSet, root, "abc", 4, 1, 0)
	td.HelperSignTransaction(prv, setTx)

	size := td.pool.Size()
	err := td.pool.AppendTx(setTx)
	require.ErrorIs(t, err, executor.InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})
	require.Equal(t, size, td.pool.Size())

	err = td.pool.AppendTxAndBroadcast(setTx)
	require.ErrorIs(t, err, executor.InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})
	require.Equal(t, size, td.pool.Size())
	select {
	case <-td.broadcastPipe.UnsafeGetChannel():
		require.Fail(t, "anchor was broadcast")
	default:
	}

	delTx := tx.NewAnchorTx(td.sbx.CurrentHeight(), from, payload.AnchorActionDelete, nil, "", 0, 0, 1)
	td.HelperSignTransaction(prv, delTx)
	err = td.pool.AppendTx(delTx)
	require.ErrorIs(t, err, executor.InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})
	require.Equal(t, size, td.pool.Size())

	transfer := td.makeTransferTx()
	td.mockExecution(transfer, nil)
	require.NoError(t, td.pool.AppendTx(transfer))
	require.Equal(t, size+1, td.pool.Size())

	require.Equal(t, amount.Amount(0), td.pool.EstimatedFee(1, payload.TypeAnchor))
	require.Equal(t, amount.Amount(0), td.pool.EstimatedFee(1, payload.Type(8)))
	require.Equal(t, td.pool.config.fixedFee(), td.pool.EstimatedFee(1, payload.TypeTransfer))
}
