package executor

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	sbx       *sandbox.FakeSandbox
	totalCoin amount.Amount
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	sbx := sandbox.NewFakeSandbox(ts)

	return &testData{
		TestSuite: ts,
		sbx:       sbx,
	}
}

func (td *testData) randStake() amount.Amount {
	return td.RandAmountRange(td.sbx.Params().MinimumStake, td.sbx.Params().MaximumStake)
}

func (td *testData) randAmountFee(max amount.Amount) (amt, fee amount.Amount) {
	for {
		amt = td.RandAmountRange(0, max)
		fee = td.RandFee()

		if amt+fee < max {
			return amt, fee
		}
	}
}

func (td *testData) addTestValidator(t *testing.T, opts ...testsuite.ValidatorMakerOption) *validator.Validator {
	t.Helper()

	val := td.GenerateTestValidator(opts...)
	td.sbx.AddValidator(val)
	td.totalCoin += val.Stake()

	return val
}

func (td *testData) addTestAccount(t *testing.T, opts ...testsuite.AccountMakerOption) (
	*account.Account, crypto.Address,
) {
	t.Helper()

	acc, addr := td.GenerateTestAccount(opts...)
	td.sbx.AddAccount(addr, acc)
	td.totalCoin += acc.Balance()

	return acc, addr
}

func (td *testData) checkTotalCoin(t *testing.T, fee amount.Amount) {
	t.Helper()

	total := amount.Amount(0)
	for _, acc := range td.sbx.FakeAccounts {
		total += acc.Balance()
	}

	for _, val := range td.sbx.FakeValidators {
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

	err = exe.Check(td.sbx, strict)
	require.ErrorIs(t, err, expectedErr)
}

func (td *testData) execute(t *testing.T, trx *tx.Tx) {
	t.Helper()

	exe, err := MakeExecutor(trx, td.sbx)
	require.NoError(t, err)

	exe.Execute(td.sbx)
}

func TestBondAndUnbond(t *testing.T) {
	td := setup(t)
	_, senderAddr := td.addTestAccount(t,
		testsuite.AccountWithBalance(10_000e9))
	valPub, _ := td.RandBLSKeyPair()
	valAddr := valPub.ValidatorAddress()

	stake := td.RandAmountRange(td.sbx.Params().MinimumStake, td.sbx.Params().MaximumStake)
	fee := td.RandFee()
	lockTime := td.sbx.CurrentHeight()

	require.Nil(t, td.sbx.Validator(valAddr))

	t.Run("First bond", func(t *testing.T) {
		trxBond := tx.NewBondTx(lockTime, senderAddr, valAddr, valPub, stake, fee)

		td.sbx.FakeCommittee.EXPECT().Contains(valAddr).Return(false).Times(1)
		td.sbx.EXPECT().IsJoinedCommittee(valAddr).Return(false).Times(1)
		td.sbx.EXPECT().UpdatePowerDelta(stake.ToNanoPAC()).Times(1)

		td.check(t, trxBond, true, nil)
		td.check(t, trxBond, false, nil)
		td.execute(t, trxBond)

		updatedVal := td.sbx.Validator(valAddr)
		require.NotNil(t, updatedVal)

		assert.Equal(t, stake, updatedVal.Stake())
		assert.Equal(t, stake.ToNanoPAC(), updatedVal.Power())
		assert.Equal(t, td.sbx.CurrentHeight(), updatedVal.LastBondingHeight())
		assert.Equal(t, types.Height(0), updatedVal.UnbondingHeight())
	})

	t.Run("First bond in same block", func(t *testing.T) {
		trxUnbond := tx.NewUnbondTx(lockTime, valAddr)

		td.sbx.FakeCommittee.EXPECT().Contains(valAddr).Return(false).Times(1)
		td.sbx.EXPECT().IsJoinedCommittee(valAddr).Return(false).Times(1)
		td.sbx.EXPECT().UpdatePowerDelta(-1 * stake.ToNanoPAC()).Return().Times(1)

		td.check(t, trxUnbond, true, nil)
		td.check(t, trxUnbond, false, nil)
		td.execute(t, trxUnbond)

		updatedVal := td.sbx.Validator(valAddr)
		require.NotNil(t, updatedVal)

		assert.Equal(t, stake, updatedVal.Stake())
		assert.Equal(t, int64(0), updatedVal.Power())
		assert.Equal(t, td.sbx.CurrentHeight(), updatedVal.LastBondingHeight())
		assert.Equal(t, td.sbx.CurrentHeight(), updatedVal.UnbondingHeight())
	})

	td.checkTotalCoin(t, fee)
}

func TestTransferBetweenAccounts(t *testing.T) {
	td := setup(t)

	senderAcc, addr1 := td.addTestAccount(t)
	senderBalance := senderAcc.Balance()
	addr2 := td.RandAccAddress()
	addr3 := td.RandAccAddress()

	amt1, fee1 := td.randAmountFee(senderBalance)
	amt2, fee2 := td.randAmountFee(amt1)

	lockTime := td.sbx.CurrentHeight()

	require.NotNil(t, td.sbx.Account(addr1))
	require.Nil(t, td.sbx.Account(addr2))
	require.Nil(t, td.sbx.Account(addr3))

	t.Run("First transfer", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, addr1, addr2, amt1, fee1)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)

		require.NotNil(t, td.sbx.Account(addr1))
		require.NotNil(t, td.sbx.Account(addr2))
		require.Nil(t, td.sbx.Account(addr3))
	})

	t.Run("Second transfer", func(t *testing.T) {
		trx := tx.NewTransferTx(lockTime, addr2, addr3, amt2, fee2)

		td.check(t, trx, true, nil)
		td.check(t, trx, false, nil)
		td.execute(t, trx)

		require.NotNil(t, td.sbx.Account(addr1))
		require.NotNil(t, td.sbx.Account(addr2))
		require.NotNil(t, td.sbx.Account(addr3))
	})

	td.checkTotalCoin(t, fee1+fee2)
}
