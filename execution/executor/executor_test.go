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

	sbx *sandbox.MockSandbox
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	sbx := sandbox.MockingSandbox(ts)
	randHeight := ts.RandHeight()
	_ = sbx.TestStore.AddTestBlock(randHeight)

	return &testData{
		TestSuite: ts,
		sbx:       sbx,
	}
}

func (td *testData) checkTotalCoin(t *testing.T, fee amount.Amount) {
	t.Helper()

	total := amount.Amount(0)
	for _, acc := range td.sbx.TestStore.Accounts {
		total += acc.Balance()
	}

	for _, val := range td.sbx.TestStore.Validators {
		total += val.Stake()
	}
	assert.Equal(t, total+fee, amount.Amount(21_000_000*1e9))
}

func (td *testData) check(t *testing.T, trx *tx.Tx, strict bool, expectedErr error) {
	t.Helper()

	exe, err := MakeExecutor(trx, td.sbx)
	if err != nil {
		assert.ErrorIs(t, err, expectedErr)

		return
	}

	err = exe.Check(strict)
	assert.ErrorIs(t, err, expectedErr)
}

func (td *testData) execute(t *testing.T, trx *tx.Tx) {
	t.Helper()

	exe, err := MakeExecutor(trx, td.sbx)
	require.NoError(t, err)

	exe.Execute()
}
