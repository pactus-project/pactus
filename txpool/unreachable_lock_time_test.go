package txpool

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// execution.CheckLockTime only enforces an upper bound on
// LockTime in strict mode (used when a block is proposed/validated); the
// pool's admission path (txPool.checkTx -> execution.CheckAndExecute(...,
// strict=false)) accepts any future LockTime, however far. Combined with
// pool.prune() always evicting the HEAD of the sub-pool (the oldest entry,
// regardless of whether it is close to becoming eligible), a handful of paid
// transactions with a lock time far beyond any realistic chain lifetime can
// evict an honest, soon-mineable transaction that was already sitting in the
// pool — at essentially no cost to the attacker, since these transactions can
// never be mined and their sender is therefore never actually charged.
//
// This test expresses the desired behavior (an honest, near-term transaction
// must not be evicted by transactions that can never become eligible) and
// failed before the fix.
func TestUnreachableLockTimeSpamEvictsHonestTx(t *testing.T) {
	conf := testDefaultConfig()
	conf.MaxSize = 10 // transferPoolSize() == 5
	td := setup(t, conf)

	honestTx := td.makeTransferTx() // lock time == current height: mineable as soon as it is included.
	td.mockExecution(honestTx, nil)
	require.NoError(t, td.pool.AppendTx(honestTx))
	require.True(t, td.pool.HasTx(honestTx.ID()), "sanity check: the honest transaction was admitted")

	// Far beyond any realistic chain lifetime (current height + 1 billion
	// blocks). Non-strict admission accepts it unconditionally: there is no
	// upper bound anywhere in the admission path.
	unreachableLockTime := types.Height(uint32(td.sbx.CurrentHeight())) + 1_000_000_000

	for i := 0; i < td.pool.config.transferPoolSize(); i++ {
		spam := td.makeTransferTx(testsuite.TransactionWithLockTime(unreachableLockTime))
		td.mockExecution(spam, nil)
		require.NoError(t, td.pool.AppendTx(spam), "unreachable lock time transactions are still admitted")
	}

	assert.True(t, td.pool.HasTx(honestTx.ID()),
		"an honest, soon-mineable transaction must not be evicted by spam whose lock time will never be reached")
}
