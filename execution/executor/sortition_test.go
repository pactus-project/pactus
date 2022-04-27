package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/types/validator"
	"github.com/zarbchain/zarb-go/util/errors"
)

func TestExecuteSortitionTx(t *testing.T) {
	setup(t)
	exe := NewSortitionExecutor(true)

	existingVal := tSandbox.TestStore.RandomTestVal()
	val := tSandbox.TestStore.AddTestValidator()
	proof := sortition.GenerateRandomProof()

	val.UpdateLastBondingHeight(tSandbox.CurrentHeight() - tSandbox.BondInterval())
	tSandbox.UpdateValidator(val)
	t.Run("Should fail, Invalid address", func(t *testing.T) {
		trx := tx.NewSortitionTx(tStamp500000, 1, crypto.GenerateTestAddress(), proof)
		tSandbox.AcceptTestSortition = true
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidAddress)
	})

	val.UpdateLastBondingHeight(tSandbox.CurrentHeight() - tSandbox.BondInterval() + 1)
	tSandbox.UpdateValidator(val)
	t.Run("Should fail, Bonding period", func(t *testing.T) {

		trx := tx.NewSortitionTx(tStamp500000, val.Sequence()+1, val.Address(), proof)
		tSandbox.AcceptTestSortition = true
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidHeight)
	})

	tSandbox.TestStore.AddTestBlock(500001)

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewSortitionTx(tStamp500000, val.Sequence()+2, val.Address(), proof)
		tSandbox.AcceptTestSortition = true
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidSequence)
	})

	t.Run("Should fail, Invalid proof", func(t *testing.T) {
		trx := tx.NewSortitionTx(tStamp500000, val.Sequence()+1, val.Address(), proof)
		tSandbox.AcceptTestSortition = false
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidProof)
	})

	t.Run("Should fail, Committee has free seats and validator is in the committee", func(t *testing.T) {
		trx := tx.NewSortitionTx(tStamp500000, existingVal.Sequence()+1, existingVal.Address(), proof)
		tSandbox.AcceptTestSortition = true
		assert.Equal(t, errors.Code(exe.Execute(trx, tSandbox)), errors.ErrInvalidTx)
	})

	t.Run("Should be ok", func(t *testing.T) {
		trx := tx.NewSortitionTx(tStamp500000, val.Sequence()+1, val.Address(), proof)
		tSandbox.AcceptTestSortition = true
		assert.NoError(t, exe.Execute(trx, tSandbox))

		// Execute again, should fail
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	assert.Equal(t, tSandbox.Validator(val.Address()).LastJoinedHeight(), tSandbox.CurrentHeight())
	assert.Zero(t, exe.Fee())

	checkTotalCoin(t, 0)
}

func TestSortitionNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewSortitionExecutor(true)
	exe2 := NewSortitionExecutor(false)

	val := tSandbox.TestStore.RandomTestVal()
	proof := sortition.GenerateRandomProof()

	tSandbox.AcceptTestSortition = true
	trx := tx.NewSortitionTx(tStamp500000, val.Sequence(), val.Address(), proof)
	assert.Error(t, exe1.Execute(trx, tSandbox))
	assert.NoError(t, exe2.Execute(trx, tSandbox))
}

func TestChangePower1(t *testing.T) {
	setup(t)

	exe := NewSortitionExecutor(true)

	// Let's create validators first
	pub1, _ := bls.GenerateTestKeyPair()
	amt1 := tSandbox.Committee().TotalPower() / 3
	val1 := tSandbox.MakeNewValidator(pub1)
	val1.AddToStake(amt1 - 1)
	val1.UpdateLastBondingHeight(tSandbox.CurrentHeight() - tSandbox.BondInterval())
	tSandbox.UpdateValidator(val1)
	proof1 := sortition.GenerateRandomProof()

	pub2, _ := bls.GenerateTestKeyPair()
	val2 := tSandbox.MakeNewValidator(pub2)
	val2.AddToStake(2)
	val2.UpdateLastBondingHeight(tSandbox.CurrentHeight() - tSandbox.BondInterval())
	tSandbox.UpdateValidator(val2)
	proof2 := sortition.GenerateRandomProof()

	val3 := tSandbox.Committee().Validators()[0]
	proof3 := sortition.GenerateRandomProof()

	tSandbox.Params.CommitteeSize = 4
	tSandbox.AcceptTestSortition = true
	trx1 := tx.NewSortitionTx(tStamp500000, val1.Sequence()+1, val1.Address(), proof1)
	assert.NoError(t, exe.Execute(trx1, tSandbox))

	trx2 := tx.NewSortitionTx(tStamp500000, val2.Sequence()+1, val2.Address(), proof2)
	assert.Error(t, exe.Execute(trx2, tSandbox), "More than 1/3 of power is joining at the same height")

	// Committee member
	trx3 := tx.NewSortitionTx(tStamp500000, val3.Sequence()+1, val3.Address(), proof3)
	assert.NoError(t, exe.Execute(trx3, tSandbox))
}

func TestChangePower2(t *testing.T) {
	setup(t)

	exe := NewSortitionExecutor(true)

	// Let's create validators first
	pub1, _ := bls.GenerateTestKeyPair()
	val1 := tSandbox.MakeNewValidator(pub1)
	val1.AddToStake(1)
	val1.UpdateLastBondingHeight(tSandbox.CurrentHeight() - tSandbox.BondInterval())
	tSandbox.UpdateValidator(val1)
	proof1 := sortition.GenerateRandomProof()

	pub2, _ := bls.GenerateTestKeyPair()
	val2 := tSandbox.MakeNewValidator(pub2)
	val2.AddToStake(1)
	val2.UpdateLastBondingHeight(tSandbox.CurrentHeight() - tSandbox.BondInterval())
	tSandbox.UpdateValidator(val2)
	proof2 := sortition.GenerateRandomProof()

	pub3, _ := bls.GenerateTestKeyPair()
	val3 := tSandbox.MakeNewValidator(pub3)
	val3.AddToStake(1)
	val3.UpdateLastBondingHeight(tSandbox.CurrentHeight() - tSandbox.BondInterval())
	tSandbox.UpdateValidator(val3)
	proof3 := sortition.GenerateRandomProof()

	val4 := tSandbox.Committee().Validators()[0]
	proof4 := sortition.GenerateRandomProof()

	tSandbox.Params.CommitteeSize = 7
	tSandbox.AcceptTestSortition = true
	trx1 := tx.NewSortitionTx(tStamp500000, val1.Sequence()+1, val1.Address(), proof1)
	assert.NoError(t, exe.Execute(trx1, tSandbox))

	trx2 := tx.NewSortitionTx(tStamp500000, val2.Sequence()+1, val2.Address(), proof2)
	assert.NoError(t, exe.Execute(trx2, tSandbox))

	trx3 := tx.NewSortitionTx(tStamp500000, val3.Sequence()+1, val3.Address(), proof3)
	assert.Error(t, exe.Execute(trx3, tSandbox), "More than 1/3 of power is leaving at the same height")

	// Committee member
	trx4 := tx.NewSortitionTx(tStamp500000, val4.Sequence()+1, val4.Address(), proof4)
	assert.NoError(t, exe.Execute(trx4, tSandbox))
}

// TestOldestDidNotPropose tests if the oldest validator in the committee had chance to propose a block or not.
func TestOldestDidNotPropose(t *testing.T) {
	setup(t)

	exe := NewSortitionExecutor(true)

	// Let's create validators first
	vals := make([]*validator.Validator, 9)
	for i := 0; i < 9; i++ {
		pub, _ := bls.GenerateTestKeyPair()
		val := tSandbox.MakeNewValidator(pub)
		val.AddToStake(1 * 10e8)
		val.UpdateLastBondingHeight(tSandbox.CurrentHeight() - tSandbox.BondInterval())
		tSandbox.UpdateValidator(val)
		vals[i] = val
	}

	tSandbox.Params.CommitteeSize = 7
	tSandbox.AcceptTestSortition = true

	stamp := tStamp500000
	for i := 0; i < 8; i = i + 2 {
		b := tSandbox.TestStore.AddTestBlock(int32(500001 + (i / 2)))
		stamp = b.Stamp()

		trx1 := tx.NewSortitionTx(stamp, vals[i].Sequence()+1, vals[i].Address(), sortition.GenerateRandomProof())
		assert.NoError(t, exe.Execute(trx1, tSandbox))

		trx2 := tx.NewSortitionTx(stamp, vals[i+1].Sequence()+1, vals[i+1].Address(), sortition.GenerateRandomProof())
		assert.NoError(t, exe.Execute(trx2, tSandbox))

		joined := make([]*validator.Validator, 0)
		tSandbox.IterateValidators(func(vs *sandbox.ValidatorStatus) {
			if vs.Validator.LastJoinedHeight() == tSandbox.CurrentHeight() {
				joined = append(joined, &vs.Validator)
			}
		})
		tSandbox.TestCommittee.Update(0, joined)
	}

	trx := tx.NewSortitionTx(stamp, vals[8].Sequence()+1, vals[8].Address(), sortition.GenerateRandomProof())
	assert.Error(t, exe.Execute(trx, tSandbox))
}
