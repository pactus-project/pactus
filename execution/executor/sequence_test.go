package executor

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// totalCoins returns the coins of the accounts and the stakes of the validators.
func (td *testData) totalCoins() amount.Amount {
	total := amount.Amount(0)
	for _, acc := range td.sbx.FakeAccounts {
		total += acc.Balance()
		total += acc.LockedDeposit()
	}
	for _, val := range td.sbx.FakeValidators {
		total += val.Stake()
	}

	return total
}

// A long random sequence of transfers, bonds, unbonds and withdraws, valid or not, is run
// on the same state. After each transaction:
//   - a transaction that is refused changes no coin,
//   - a transaction that is executed removes its fee from the coins, and nothing else,
//   - no balance and no stake is negative, and no stake is above the maximum.
func TestRandomSequenceKeepsTheCoins(t *testing.T) {
	td := setup(t)
	sbx := td.sbx

	sbx.FakeCommittee.EXPECT().Contains(gomock.Any()).Return(false).AnyTimes()
	sbx.FakeCommittee.EXPECT().Size().Return(51).AnyTimes()
	sbx.EXPECT().IsJoinedCommittee(gomock.Any()).Return(false).AnyTimes()
	sbx.EXPECT().UpdatePowerDelta(gomock.Any()).AnyTimes()

	accounts := make([]crypto.Address, 0, 8)
	for range 5 {
		_, addr := td.addTestAccount(t, testsuite.AccountWithBalance(td.RandAmountRange(100e9, 5000e9)))
		accounts = append(accounts, addr)
	}
	// Some addresses have no account yet.
	accounts = append(accounts, td.RandAccAddress(), td.RandAccAddress(), td.RandAccAddress())

	valKeys := make([]*bls.PublicKey, 0, 4)
	for range 4 {
		pub, _ := td.RandBLSKeyPair()
		valKeys = append(valKeys, pub)
	}

	randAccount := func() crypto.Address { return accounts[td.RandIntMax(len(accounts))] }
	randValKey := func() *bls.PublicKey { return valKeys[td.RandIntMax(len(valKeys))] }
	// The amounts go from small to larger than what an account holds.
	randAmount := func() amount.Amount { return td.RandAmountRange(0, 2000e9) }

	executed := map[payload.Type]int{}
	refused := map[payload.Type]int{}

	for step := range 3000 {
		// The height moves on, so that an unbonded validator can withdraw later.
		if step%10 == 0 {
			sbx.FakeHeight += 2000
		}

		lockTime := sbx.CurrentHeight()
		fee := td.RandFee()

		var trx *tx.Tx
		switch td.RandIntMax(5) {
		case 0:
			trx = tx.NewTransferTx(lockTime, randAccount(), randAccount(), randAmount(), fee)
		case 1:
			recipients := []payload.BatchRecipient{
				{To: randAccount(), Amount: randAmount()},
				{To: randAccount(), Amount: randAmount()},
				{To: td.RandAccAddress(), Amount: td.RandAmountRange(1, 10e9)},
			}
			trx = tx.NewBatchTransferTx(lockTime, randAccount(), recipients, fee)
		case 2:
			pub := randValKey()
			trx = tx.NewBondTx(lockTime, randAccount(), pub.ValidatorAddress(), pub, randAmount(), fee)
		case 3:
			trx = tx.NewUnbondTx(lockTime, randValKey().ValidatorAddress())
			fee = 0
		default:
			trx = tx.NewWithdrawTx(lockTime, randValKey().ValidatorAddress(), randAccount(), randAmount(), fee)
		}
		typ := trx.Payload().Type()

		// The executors run on the transactions that passed the basic checks:
		// a batch with the same recipient twice, for example, never reaches them.
		if err := trx.Payload().BasicCheck(); err != nil {
			continue
		}

		before := td.totalCoins()

		exe, err := MakeExecutor(trx, sbx)
		if err == nil {
			err = exe.Check(sbx, true)
		}
		if err != nil {
			refused[typ]++
			require.Equal(t, before, td.totalCoins(), "step %d: a refused %s changed the coins: %v", step, typ, err)

			continue
		}

		exe.Execute(sbx)
		executed[typ]++

		require.Equal(t, before-fee, td.totalCoins(), "step %d: %s", step, typ)
		for addr, acc := range sbx.FakeAccounts {
			require.GreaterOrEqual(t, acc.Balance(), amount.Amount(0), "step %d: account %s", step, addr)
		}
		for addr, val := range sbx.FakeValidators {
			require.GreaterOrEqual(t, val.Stake(), amount.Amount(0), "step %d: validator %s", step, addr)
			require.LessOrEqual(t, val.Stake(), sbx.Params().MaximumStake, "step %d: validator %s", step, addr)
		}
	}

	// The sequence exercised each type of transaction, in both ways.
	for _, typ := range []payload.Type{
		payload.TypeTransfer, payload.TypeBatchTransfer, payload.TypeBond, payload.TypeUnbond, payload.TypeWithdraw,
	} {
		require.Positive(t, executed[typ], "no %s was executed: %v", typ, executed)
		require.Positive(t, refused[typ], "no %s was refused: %v", typ, refused)
	}
}
