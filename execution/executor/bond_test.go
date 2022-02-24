package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecuteBondTx(t *testing.T) {
	setup(t)
	exe := NewBondExecutor(true)

	bonder := tAcc1.Address()
	pub, _ := bls.GenerateTestKeyPair()
	addr := pub.Address()
	hash100 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(100, hash100)

	t.Run("Should fail, Invalid bonder", func(t *testing.T) {
		trx := tx.NewBondTx(hash100.Stamp(), 1, pub.Address(), pub, 1000, 1000, "invalid bonder")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewBondTx(hash100.Stamp(), tSandbox.AccSeq(bonder)+2, bonder, pub, 1000, 1000, "invalid sequence")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Insufficient balance", func(t *testing.T) {
		trx := tx.NewBondTx(hash100.Stamp(), tSandbox.AccSeq(bonder)+1, bonder, pub, 10000000000, 10000000, "insufficient balance")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Inside committee", func(t *testing.T) {
		tSandbox.InCommittee = true
		trx := tx.NewBondTx(hash100.Stamp(), tSandbox.AccSeq(bonder)+1, bonder, pub, 1000, 1000, "inside committee")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Ok", func(t *testing.T) {
		tSandbox.InCommittee = false
		trx := tx.NewBondTx(hash100.Stamp(), tSandbox.AccSeq(bonder)+1, bonder, pub, 1000, 1000, "ok")

		assert.NoError(t, exe.Execute(trx, tSandbox))

		// Replay
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	assert.Equal(t, tSandbox.Account(bonder).Balance(), int64(10000000000-2000))
	assert.Equal(t, tSandbox.Validator(addr).Stake(), int64(1000))
	assert.Equal(t, tSandbox.Validator(addr).LastBondingHeight(), 101)
	tSandbox.AppendNewBlock(101, hash.GenerateTestHash())

	t.Run("Should be able to rebond if hasn't ever unbonded", func(t *testing.T) {
		tSandbox.InCommittee = false
		trx := tx.NewBondTx(hash100.Stamp(), tSandbox.AccSeq(bonder)+1, bonder, pub, 1000, 1000, "Rebond")

		assert.NoError(t, exe.Execute(trx, tSandbox))
	})

	assert.Equal(t, tSandbox.Validator(addr).Power(), int64(2000))
	assert.Equal(t, tSandbox.Validator(addr).Stake(), int64(2000))
	assert.Equal(t, tSandbox.Validator(addr).LastBondingHeight(), 102)
	tSandbox.AppendNewBlock(102, hash.GenerateTestHash())

	t.Run("Shouldn't be able to rebond after unbonding", func(t *testing.T) {
		tSandbox.InCommittee = false
		uexe := NewUnbondExecutor(true)

		unbondTrx := tx.NewUnbondTx(hash100.Stamp(), 1, addr, "Unbond")
		assert.NoError(t, uexe.Execute(unbondTrx, tSandbox))

		trx := tx.NewBondTx(hash100.Stamp(), tSandbox.AccSeq(bonder)+1, bonder, pub, 1000, 1000, "Rebond")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	assert.Equal(t, int64(10000000000-4000), tSandbox.Account(bonder).Balance())
	assert.Equal(t, int64(2000), tSandbox.Validator(addr).Stake())
	assert.Equal(t, int64(0), tSandbox.Validator(addr).Power())
	assert.Equal(t, 102, tSandbox.Validator(addr).LastBondingHeight())
	assert.Equal(t, 103, tSandbox.Validator(addr).UnbondingHeight())

	assert.Equal(t, exe.Fee(), int64(1000))

	checkTotalCoin(t, 2000)
}

func TestBondNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewBondExecutor(true)
	exe2 := NewBondExecutor(false)

	tSandbox.InCommittee = true
	hash100 := hash.GenerateTestHash()
	tSandbox.AppendNewBlock(100, hash100)
	pub, _ := bls.GenerateTestKeyPair()

	trx := tx.NewBondTx(hash100.Stamp(), tSandbox.AccSeq(tAcc1.Address())+1, tAcc1.Address(), pub, 1000, 1000, "")

	assert.Error(t, exe1.Execute(trx, tSandbox))
	assert.NoError(t, exe2.Execute(trx, tSandbox))
}
