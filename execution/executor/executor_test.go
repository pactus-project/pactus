package executor

import (
	"testing"

	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	sandbox *sandbox.MockSandbox
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	sb := sandbox.MockingSandbox(ts)
	randHeight := ts.RandHeight()
	_ = sb.TestStore.AddTestBlock(randHeight)

	return &testData{
		TestSuite: ts,
		sandbox:   sb,
	}
}

func (td *testData) checkTotalCoin(t *testing.T, fee amount.Amount) {
	t.Helper()

	total := amount.Amount(0)
	for _, acc := range td.sandbox.TestStore.Accounts {
		total += acc.Balance()
	}

	for _, val := range td.sandbox.TestStore.Validators {
		total += val.Stake()
	}
	assert.Equal(t, total+fee, amount.Amount(21_000_000*1e9))
}

func (td *testData) check(t *testing.T, trx *tx.Tx, strict bool, expectedErr error) {
	exe, err := MakeExecutor(trx, td.sandbox)
	if err != nil {
		assert.ErrorIs(t, err, expectedErr)

		return
	}

	err = exe.Check(strict)
	assert.ErrorIs(t, err, expectedErr)
}

func (td *testData) execute(t *testing.T, trx *tx.Tx) {
	exe, err := MakeExecutor(trx, td.sandbox)
	require.NoError(t, err)

	exe.Execute()
}
