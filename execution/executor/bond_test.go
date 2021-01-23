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

	bonder := tAcc1Signer.Address()
	addr1, pub1, priv1 := crypto.GenerateTestKeyPair()
	signer1 := crypto.NewSigner(priv1)
	stamp := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(100, stamp)

	trx1 := tx.NewBondTx(stamp, 1, addr1, pub1, 1000, 1000, "invalid boner")
	tAcc1Signer.SignMsg(trx1)
	assert.Error(t, exe.Execute(trx1))

	trx2 := tx.NewBondTx(stamp, tSandbox.AccSeq(bonder)+2, bonder, pub1, 1000, 1000, "invalid sequence")
	signer1.SignMsg(trx2)
	assert.Error(t, exe.Execute(trx2))

	trx3 := tx.NewBondTx(stamp, tSandbox.AccSeq(bonder)+1, bonder, pub1, 10000000000, 10000000, "insufficient balance")
	signer1.SignMsg(trx3)
	assert.Error(t, exe.Execute(trx3))

	trx5 := tx.NewBondTx(stamp, tSandbox.AccSeq(bonder)+1, bonder, pub1, 1000, 1000, "ok")
	signer1.SignMsg(trx5)
	assert.NoError(t, exe.Execute(trx5))

	// Duplicated. Invalid sequence
	assert.Error(t, exe.Execute(trx5))

	assert.Equal(t, tSandbox.Account(bonder).Balance(), int64(10000000000-2000))
	assert.Equal(t, tSandbox.Validator(addr1).Stake(), int64(1000))
	assert.Equal(t, tSandbox.Validator(addr1).BondingHeight(), 102)
	assert.Equal(t, exe.Fee(), int64(1000))

	checkTotalCoin(t, 1000)
}
