package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecuteBondTx(t *testing.T) {
	setup(t)
	exe := NewBondExecutor(tSandbox)

	bonder := tAcc1.Address()
	addr, pub, _ := crypto.GenerateTestKeyPair()
	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(100, stamp)

	t.Run("Should fail, Invalid bonder", func(t *testing.T) {
		trx := tx.NewBondTx(stamp, 1, addr, pub, 1000, 1000, "invalid bonder")
		assert.Error(t, exe.Execute(trx))
	})

	t.Run("Should fail, Invalid sequence", func(t *testing.T) {
		trx := tx.NewBondTx(stamp, tSandbox.AccSeq(bonder)+2, bonder, pub, 1000, 1000, "invalid sequence")

		assert.Error(t, exe.Execute(trx))
	})

	t.Run("Should fail, Insufficient balance", func(t *testing.T) {
		trx := tx.NewBondTx(stamp, tSandbox.AccSeq(bonder)+1, bonder, pub, 10000000000, 10000000, "insufficient balance")

		assert.Error(t, exe.Execute(trx))
	})

	t.Run("Ok", func(t *testing.T) {
		trx := tx.NewBondTx(stamp, tSandbox.AccSeq(bonder)+1, bonder, pub, 1000, 1000, "ok")

		assert.NoError(t, exe.Execute(trx))

		// Replay
		assert.Error(t, exe.Execute(trx))
	})
	assert.Equal(t, tSandbox.Account(bonder).Balance(), int64(10000000000-2000))
	assert.Equal(t, tSandbox.Validator(addr).Stake(), int64(1000))
	assert.Equal(t, tSandbox.Validator(addr).BondingHeight(), 102)
	assert.Equal(t, exe.Fee(), int64(1000))

	checkTotalCoin(t, 1000)
}
