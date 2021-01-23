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

	addr1, _, priv1 := crypto.GenerateTestKeyPair()
	stamp40 := crypto.GenerateTestHash()
	stamp41 := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(40, stamp40)

	proof1 := tVal1Signer.SignData(stamp40.RawBytes()).RawBytes()
	trx1 := tx.NewSortitionTx(stamp40, 1, addr1, proof1, "invalid-address")
	trx1.SetSignature(priv1.Sign(trx1.SignBytes()))
	assert.Error(t, exe.Execute(trx1))

	trx2 := tx.NewSortitionTx(stamp40, 2, tVal1Signer.Address(), proof1, "invalid-sequence")
	tVal1Signer.SignMsg(trx2)
	assert.Error(t, exe.Execute(trx2))

	trx3 := tx.NewSortitionTx(stamp40, 1, tVal1Signer.Address(), proof1, "too early")
	tVal1Signer.SignMsg(trx3)
	assert.Error(t, exe.Execute(trx3))

	tSandbox.AppendStampAndUpdateHeight(41, stamp41)

	trx4 := tx.NewSortitionTx(stamp41, 1, tVal1Signer.Address(), proof1, "ok")
	tVal1Signer.SignMsg(trx4)

	tSandbox.AcceptSortition = true
	tSandbox.ErrorAddToSet = true
	assert.Error(t, exe.Execute(trx4))

	tSandbox.AcceptSortition = false
	tSandbox.ErrorAddToSet = true
	assert.Error(t, exe.Execute(trx4))

	tSandbox.AcceptSortition = true
	tSandbox.ErrorAddToSet = false
	assert.NoError(t, exe.Execute(trx4))

	val := tSandbox.Validator(tVal1Signer.Address())
	assert.Equal(t, val.Sequence(), 1)
	assert.Equal(t, val.LastJoinedHeight(), 42)
	assert.NoError(t, exe.Execute(trx2))
	assert.Error(t, exe.Execute(trx2)) // invalid sequence
	assert.Zero(t, exe.Fee())
}
