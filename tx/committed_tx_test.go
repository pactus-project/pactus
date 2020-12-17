package tx

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestSerialization(t *testing.T) {
	tx, _ := GenerateTestSendTx()
	r := tx.GenerateReceipt(Ok, crypto.GenerateTestHash())
	ctrx := CommittedTx{tx, r}
	d, err := cbor.Marshal(ctrx)
	assert.NoError(t, err)
	ctrx2 := new(CommittedTx)
	err = cbor.Unmarshal(d, ctrx2)
	assert.NoError(t, err)
	assert.Equal(t, tx.ID(), ctrx2.Tx.ID())
	assert.Equal(t, r.Hash(), ctrx2.Receipt.Hash())
	assert.Equal(t, r, ctrx2.Receipt)
}

func TestCommittedTxSanityCheck(t *testing.T) {
	tx, _ := GenerateTestSendTx()
	r := tx.GenerateReceipt(Ok, crypto.GenerateTestHash())
	ctrx := CommittedTx{tx, r}
	assert.NoError(t, ctrx.SanityCheck())
	ctrx.Receipt.data.TxID = crypto.GenerateTestHash()
	assert.Error(t, ctrx.SanityCheck())
}
