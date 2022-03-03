package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecuteBondTx(t *testing.T) {
	setup(t)
	exe := NewBondExecutor(true)

	bonder := tAcc1.Address()
	pub, _ := bls.GenerateTestKeyPair()
	addr := pub.Address()

	t.Run("Should fail, Invalid bonder", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, 1, pub.Address(), pub, 100000, 1000, "invalid bonder")
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, tSandbox.AccSeq(bonder)+2, bonder, pub, 100000, 1000, "invalid sequence")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Insufficient balance", func(t *testing.T) {
		trx := tx.NewBondTx(tStamp500000, tSandbox.AccSeq(bonder)+1, bonder, pub, tAcc1Balance+1, 0, "insufficient balance")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Should fail, Inside committee", func(t *testing.T) {
		tSandbox.InCommittee = true
		trx := tx.NewBondTx(tStamp500000, tSandbox.AccSeq(bonder)+1, bonder, pub, 100000, 1000, "inside committee")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Ok", func(t *testing.T) {
		tSandbox.InCommittee = false
		trx := tx.NewBondTx(tStamp500000, tSandbox.AccSeq(bonder)+1, bonder, pub, 100000, 1000, "ok")

		assert.NoError(t, exe.Execute(trx, tSandbox))

		// Replay
		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	t.Run("Unbonded before", func(t *testing.T) {
		tVal1.UpdateUnbondingHeight(tSandbox.CurHeight)
		trx := tx.NewBondTx(tStamp500000, tSandbox.AccSeq(bonder)+1, bonder, tVal1.PublicKey(), 100000, 1000, "ok")

		assert.Error(t, exe.Execute(trx, tSandbox))
	})

	assert.Equal(t, tSandbox.Account(bonder).Balance(), tAcc1Balance-(100000+1000))
	assert.Equal(t, tSandbox.Validator(addr).Stake(), int64(100000))
	assert.Equal(t, tSandbox.Validator(addr).LastBondingHeight(), tSandbox.CurHeight)
	assert.Equal(t, exe.Fee(), int64(1000))

	checkTotalCoin(t, 1000)
}

func TestBondNonStrictMode(t *testing.T) {
	setup(t)
	exe1 := NewBondExecutor(true)
	exe2 := NewBondExecutor(false)

	tSandbox.InCommittee = true
	pub, _ := bls.GenerateTestKeyPair()

	trx := tx.NewBondTx(tStamp500001, tSandbox.AccSeq(tAcc1.Address())+1, tAcc1.Address(), pub, 1000, 1000, "")

	assert.Error(t, exe1.Execute(trx, tSandbox))
	assert.NoError(t, exe2.Execute(trx, tSandbox))
}
