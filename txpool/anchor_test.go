package txpool

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAnchorRejectedByPool(t *testing.T) {
	for _, version := range []protocol.Version{
		protocol.ProtocolVersion4,
		protocol.ProtocolVersion5,
	} {
		t.Run(version.String(), func(t *testing.T) {
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

			acc := account.NewAccount(1)
			acc.AddToBalance(10 * executor.MinAnchorDeposit)
			pub, prv := td.RandBLSKeyPair()
			from := pub.AccountAddress()

			td.sbx.EXPECT().BlockVersion().Return(version).AnyTimes()
			td.sbx.EXPECT().Account(gomock.Any()).DoAndReturn(func(addr crypto.Address) *account.Account {
				if addr == from {
					return acc
				}

				return nil
			}).AnyTimes()

			root := bytes.Repeat([]byte{0x11}, 32)
			setTx := tx.NewAnchorTx(td.sbx.CurrentHeight(), from, payload.AnchorActionSet,
				root, "abc", 4, executor.MinAnchorDeposit, 0)
			td.HelperSignTransaction(prv, setTx)

			balance := acc.Balance()
			size := td.pool.Size()
			err := td.pool.AppendTx(setTx)
			require.ErrorIs(t, err, ErrPayloadTypeNotSupported)
			require.Equal(t, size, td.pool.Size())
			require.Equal(t, balance, acc.Balance())
			require.False(t, acc.HasAnchor())
			require.Equal(t, amount.Amount(0), acc.LockedDeposit())

			err = td.pool.AppendTxAndBroadcast(setTx)
			require.ErrorIs(t, err, ErrPayloadTypeNotSupported)
			require.Equal(t, size, td.pool.Size())
			select {
			case <-td.broadcastPipe.UnsafeGetChannel():
				require.Fail(t, "anchor was broadcast")
			default:
			}

			delTx := tx.NewAnchorTx(td.sbx.CurrentHeight(), from, payload.AnchorActionDelete, nil, "", 0, 0, 0)
			td.HelperSignTransaction(prv, delTx)
			err = td.pool.AppendTx(delTx)
			require.ErrorIs(t, err, ErrPayloadTypeNotSupported)
			require.Equal(t, size, td.pool.Size())
			require.Equal(t, balance, acc.Balance())

			priced := tx.NewAnchorTx(td.sbx.CurrentHeight(), from, payload.AnchorActionSet,
				root, "abc", 4, executor.MinAnchorDeposit, td.pool.config.fixedFee())
			td.HelperSignTransaction(prv, priced)
			err = td.pool.AppendTx(priced)
			require.ErrorIs(t, err, ErrPayloadTypeNotSupported)
			require.Equal(t, balance, acc.Balance())
			require.False(t, acc.HasAnchor())

			unset := tx.NewAnchorTx(0, from, payload.AnchorActionSet,
				root, "abc", 4, executor.MinAnchorDeposit, td.pool.config.fixedFee())
			err = td.pool.AppendTx(unset)
			require.Error(t, err)
			require.NotErrorIs(t, err, ErrPayloadTypeNotSupported)
			require.Equal(t, balance, acc.Balance())

			transfer := td.makeTransferTx()
			td.mockExecution(transfer, nil)
			require.NoError(t, td.pool.AppendTx(transfer))
			require.Equal(t, size+1, td.pool.Size())

			require.Equal(t, amount.Amount(0), td.pool.EstimatedFee(1, payload.TypeAnchor))
			require.Equal(t, amount.Amount(0), td.pool.EstimatedFee(1, payload.Type(8)))
			require.Equal(t, td.pool.config.fixedFee(), td.pool.EstimatedFee(1, payload.TypeTransfer))
		})
	}
}
