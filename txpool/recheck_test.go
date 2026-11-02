package txpool

import (
	"testing"

	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/execution/executor"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// realExecutionData is a transaction pool backed by a real sandbox and the
// real executors, instead of the mocked executor used by the other tests.
type realExecutionData struct {
	*testsuite.TestSuite

	pool      *txPool
	fakeStore *store.FakeStore
	params    *param.Params
}

func setupRealExecution(t *testing.T) *realExecutionData {
	t.Helper()

	// Other tests of this package replace the global executor factory by a mock.
	oldFactory := executor.DefaultFactory
	executor.DefaultFactory = executor.MakeExecutorImpl
	t.Cleanup(func() { executor.DefaultFactory = oldFactory })

	ts := testsuite.NewTestSuite(t)

	fakeStore := store.NewFakeStore(ts)
	fakeStore.EXPECT().RecentTransaction(gomock.Any()).Return(false).AnyTimes()
	fakeStore.EXPECT().IsBanned(gomock.Any()).Return(false).AnyTimes()

	broadcastPipe := pipeline.New[message.Message](t.Context())
	eventPipe := pipeline.New[any](t.Context())

	data := &realExecutionData{
		TestSuite: ts,
		pool:      NewTxPool(t.Context(), DefaultConfig(), fakeStore, broadcastPipe, eventPipe).(*txPool),
		fakeStore: fakeStore,
		params:    param.FromGenesis(genesis.MainnetGenesis()),
	}
	data.pool.SetNewSandboxAndRecheck(data.newSandbox())

	return data
}

func (data *realExecutionData) newSandbox() sandbox.Sandbox {
	return sandbox.NewSandbox(100, data.params.BlockVersion, 0,
		data.fakeStore, data.params, committee.NewFakeCommittee(data.TestSuite), 1_000_000_000_000)
}

func (data *realExecutionData) addAccount(balance amount.Amount) (crypto.Address, *bls.PrivateKey) {
	pub, prv := data.RandBLSKeyPair()
	acc := account.NewAccount(data.RandInt32())
	acc.AddToBalance(balance)
	addr := pub.AccountAddress()
	data.fakeStore.FakeAccounts[addr] = acc

	return addr, prv
}

func (data *realExecutionData) setBalance(addr crypto.Address, balance amount.Amount) {
	acc := account.NewAccount(data.fakeStore.FakeAccounts[addr].Number())
	acc.AddToBalance(balance)
	data.fakeStore.FakeAccounts[addr] = acc
}

// TestRecheckIsDeterministic checks that rechecking the pool after a sandbox
// update drops the same transaction every time, instead of depending on the
// iteration order of the map that holds the sub-pools.
func TestRecheckIsDeterministic(t *testing.T) {
	retainedBond := 0
	retainedTransfer := 0
	const rounds = 40

	for i := 0; i < rounds; i++ {
		data := setupRealExecution(t)

		sender, prv := data.addAccount(101e9)
		lockTime := data.pool.sbx.CurrentHeight()
		fee := data.pool.fixedFee()

		valPub, _ := data.RandBLSKeyPair()
		bondTrx := tx.NewBondTx(lockTime, sender, valPub.ValidatorAddress(), valPub, 50e9, fee)
		data.HelperSignTransaction(prv, bondTrx)
		transferTrx := tx.NewTransferTx(lockTime, sender, data.RandAccAddress(), 50e9, fee)
		data.HelperSignTransaction(prv, transferTrx)

		require.NoError(t, data.pool.AppendTx(bondTrx))
		require.NoError(t, data.pool.AppendTx(transferTrx))

		// The balance of the sender drops: only one of the two pooled
		// transactions can still be executed.
		data.setBalance(sender, 61e9)
		data.pool.SetNewSandboxAndRecheck(data.newSandbox())

		hasBond := data.pool.HasTx(bondTrx.ID())
		hasTransfer := data.pool.HasTx(transferTrx.ID())
		require.NotEqual(t, hasBond, hasTransfer, "exactly one transaction should stay in the pool")

		if hasBond {
			retainedBond++
		} else {
			retainedTransfer++
		}
	}

	assert.True(t, retainedBond == rounds || retainedTransfer == rounds,
		"the transaction kept after a recheck must not vary: bond kept %d times, transfer kept %d times out of %d",
		retainedBond, retainedTransfer, rounds)
}
