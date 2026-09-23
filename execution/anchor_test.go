package execution

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAnchorReplayAndLockTime(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	sbx := sandbox.NewFakeSandbox(ts)
	sbx.FakeBlockVersion = protocol.Version(5)
	sbx.FakeUnixTime = 40
	sbx.FakeHeight = 20

	addr, acc := ts.GenerateTestAccount(testsuite.AccountWithBalance(10 * executor.MinAnchorDeposit))
	sbx.AddAccount(addr, acc)

	committed := map[tx.ID]bool{}
	var fee amount.Amount
	sbx.EXPECT().IsBanned(gomock.Any()).Return(false).AnyTimes()
	sbx.EXPECT().RecentTransaction(gomock.Any()).DoAndReturn(func(id tx.ID) bool {
		return committed[id]
	}).AnyTimes()
	sbx.EXPECT().CommitTransaction(gomock.Any()).DoAndReturn(func(trx *tx.Tx) {
		committed[trx.ID()] = true
		fee += trx.Fee()
	}).AnyTimes()
	sbx.EXPECT().AccumulatedFee().DoAndReturn(func() amount.Amount {
		return fee
	}).AnyTimes()

	root := bytes.Repeat([]byte{0x11}, 32)
	future := tx.NewAnchorTx(sbx.CurrentHeight()+1, addr, payload.AnchorActionSet,
		root, "", 0, executor.MinAnchorDeposit, 1)
	before, err := acc.Bytes()
	require.NoError(t, err)
	err = CheckAndExecute(future, sbx, true)
	require.ErrorIs(t, err, LockTimeInFutureError{LockTime: future.LockTime()})
	got, err := acc.Bytes()
	require.NoError(t, err)
	require.Equal(t, before, got)

	trx := tx.NewAnchorTx(sbx.CurrentHeight(), addr, payload.AnchorActionSet,
		root, "ok", 1, executor.MinAnchorDeposit, 7)
	require.NoError(t, CheckAndExecute(trx, sbx, true))
	require.Equal(t, amount.Amount(7), sbx.AccumulatedFee())
	require.Equal(t, executor.MinAnchorDeposit, acc.LockedDeposit())
	balance := acc.Balance()

	err = CheckAndExecute(trx, sbx, false)
	require.ErrorIs(t, err, TransactionCommittedError{ID: trx.ID()})
	require.Equal(t, balance, acc.Balance())
	require.Equal(t, executor.MinAnchorDeposit, acc.LockedDeposit())
	require.Equal(t, amount.Amount(7), sbx.AccumulatedFee())
}
