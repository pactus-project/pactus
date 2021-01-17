package executor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

func TestExecuteSortitionTx(t *testing.T) {
	setup(t)
	exe := NewSortitionExecutor(tSandbox)

	addr1, pub1, priv1 := crypto.GenerateTestKeyPair()
	stamp4 := crypto.GenerateTestHash()
	stamp14 := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(4, stamp4)

	proof1 := tVal1Signer.Sign(stamp4.RawBytes()).RawBytes()
	trx1 := tx.NewSortitionTx(stamp4, 1, addr1, proof1, "invalid-address", &pub1, nil)
	trx1.SetSignature(priv1.Sign(trx1.SignBytes()))
	assert.Error(t, exe.Execute(trx1))

	trx2 := tx.NewSortitionTx(stamp4, 1, tVal1Signer.Address(), proof1, "invalid-address", &tVal1Pub, nil)
	tVal1Signer.SignMsg(trx2)
	assert.Error(t, exe.Execute(trx2))

	tSandbox.AppendStampAndUpdateHeight(14, stamp14)

	tSandbox.AcceptSortition = false
	assert.Error(t, exe.Execute(trx2))

	tSandbox.ErrorAddToSet = true
	assert.Error(t, exe.Execute(trx2))

	tSandbox.AcceptSortition = true
	tSandbox.ErrorAddToSet = false

	assert.NoError(t, exe.Execute(trx2))
	assert.Error(t, exe.Execute(trx2)) // invalid sequence
	assert.Zero(t, exe.Fee())
}
