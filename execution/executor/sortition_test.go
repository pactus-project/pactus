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
	stamp29 := crypto.GenerateTestHash()
	stamp24 := crypto.GenerateTestHash()
	tSandbox.AppendStampAndUpdateHeight(29, stamp29)

	proof1 := tVal1Signer.Sign(stamp29.RawBytes()).RawBytes()
	trx1 := tx.NewSortitionTx(stamp29, 1, addr1, proof1, "invalid-address", &pub1, nil)
	trx1.SetSignature(priv1.Sign(trx1.SignBytes()))
	assert.Error(t, exe.Execute(trx1))

	trx2 := tx.NewSortitionTx(stamp29, 2, tVal1Signer.Address(), proof1, "invalid-sequence", &tVal1Pub, nil)
	tVal1Signer.SignMsg(trx2)
	assert.Error(t, exe.Execute(trx2))

	trx3 := tx.NewSortitionTx(stamp29, 1, tVal1Signer.Address(), proof1, "too early", &tVal1Pub, nil)
	tVal1Signer.SignMsg(trx3)
	assert.Error(t, exe.Execute(trx3))

	tSandbox.AppendStampAndUpdateHeight(30, stamp24)

	trx4 := tx.NewSortitionTx(stamp24, 1, tVal1Signer.Address(), proof1, "too early", &tVal1Pub, nil)
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

	assert.Equal(t, tSandbox.Validator(tVal1Signer.Address()).Sequence(), 1)
	assert.NoError(t, exe.Execute(trx2))
	assert.Error(t, exe.Execute(trx2)) // invalid sequence
	assert.Zero(t, exe.Fee())
}
