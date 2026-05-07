package executor

import (
	"testing"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type testData struct {
	*testsuite.TestSuite

	params     *param.Params
	accounts   map[crypto.Address]*account.Account
	validators map[crypto.Address]*validator.Validator
	totalCoin  amount.Amount
	committee  *committee.MockCommittee
	sbx        *sandbox.MockSandbox
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	params := param.FromGenesis(genesis.MainnetGenesis())
	params.BlockVersion = protocol.ProtocolVersion3
	params.CommitteeSize = 7
	committee := committee.NewMockCommittee(ts.Ctrl)

	accounts := make(map[crypto.Address]*account.Account)
	validators := make(map[crypto.Address]*validator.Validator)

	sbx := sandbox.NewMockSandbox(ts.Ctrl)

	sbx.EXPECT().Account(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) *account.Account {
			return accounts[addr]
		},
	).AnyTimes()

	sbx.EXPECT().MakeNewAccount(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) *account.Account {
			return account.NewAccount(ts.RandInt32())
		},
	).AnyTimes()

	sbx.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).DoAndReturn(
		func(addr crypto.Address, acc *account.Account) {
			accounts[addr] = acc
		},
	).AnyTimes()

	sbx.EXPECT().Validator(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) *validator.Validator {
			return validators[addr]
		},
	).AnyTimes()

	sbx.EXPECT().CurrentHeight().Return(ts.RandHeight()).AnyTimes()

	sbx.EXPECT().Validator(gomock.Any()).DoAndReturn(
		func(addr crypto.Address) *validator.Validator {
			return validators[addr]
		},
	).AnyTimes()

	sbx.EXPECT().MakeNewValidator(gomock.Any()).DoAndReturn(
		func(pub *bls.PublicKey) *validator.Validator {
			return validator.NewValidator(pub, ts.RandInt32())
		},
	).AnyTimes()

	sbx.EXPECT().UpdateValidator(gomock.Any()).DoAndReturn(
		func(val *validator.Validator) {
			validators[val.Address()] = val
		},
	).AnyTimes()

	sbx.EXPECT().CurrentHeight().Return(ts.RandHeight()).AnyTimes()
	sbx.EXPECT().Params().Return(params).AnyTimes()
	sbx.EXPECT().Committee().Return(committee).AnyTimes()

	return &testData{
		TestSuite:  ts,
		params:     params,
		committee:  committee,
		accounts:   accounts,
		validators: validators,
		sbx:        sbx,
	}
}
func (td *testData) addTestValidator(t *testing.T, opts ...testsuite.ValidatorMakerOption) *validator.Validator {
	t.Helper()

	val := td.GenerateTestValidator(opts...)
	td.validators[val.Address()] = val
	td.totalCoin += val.Stake()

	return val
}

func (td *testData) addTestAccount(t *testing.T, opts ...testsuite.AccountMakerOption) (
	*account.Account, crypto.Address) {
	t.Helper()

	acc, addr := td.GenerateTestAccount(opts...)
	td.accounts[addr] = acc
	td.totalCoin += acc.Balance()

	return acc, addr
}

func (td *testData) checkTotalCoin(t *testing.T, fee amount.Amount) {
	t.Helper()

	total := amount.Amount(0)
	for _, acc := range td.accounts {
		total += acc.Balance()
	}

	for _, val := range td.validators {
		total += val.Stake()
	}
	assert.Equal(t, total+fee, td.totalCoin)
}

func (td *testData) check(t *testing.T, trx *tx.Tx, strict bool, expectedErr error) {
	t.Helper()

	exe, err := MakeExecutor(trx, td.sbx)
	if err != nil {
		require.ErrorIs(t, err, expectedErr)

		return
	}

	err = exe.Check(strict)
	require.ErrorIs(t, err, expectedErr)
}

func (td *testData) execute(t *testing.T, trx *tx.Tx) {
	t.Helper()

	exe, err := MakeExecutor(trx, td.sbx)
	require.NoError(t, err)

	exe.Execute()
}
