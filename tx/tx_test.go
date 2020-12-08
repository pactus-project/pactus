package tx

import (
	"fmt"
	"testing"

	"github.com/zarbchain/zarb-go/crypto"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestEncodingTx(t *testing.T) {
	tx, _ := GenerateTestSendTx()

	bz, err := tx.MarshalCBOR()
	require.NoError(t, err)
	var tx2 Tx
	err = tx2.UnmarshalCBOR(bz)

	bz2, _ := tx2.MarshalCBOR()

	fmt.Printf("%x\n", bz)
	fmt.Printf("%x\n", bz2)

	require.NoError(t, err)
	require.Equal(t, bz, bz2)
	require.Equal(t, tx.Hash(), tx2.Hash())
}

func TestEncodingTxNoMemo(t *testing.T) {
	tx, _ := GenerateTestSendTx()
	tx.data.Memo = ""

	bz, err := tx.MarshalCBOR()
	require.NoError(t, err)
	var tx2 Tx
	err = tx2.UnmarshalCBOR(bz)

	bz2, _ := tx2.MarshalCBOR()

	fmt.Printf("%x\n", bz)
	fmt.Printf("%x\n", bz2)

	require.NoError(t, err)
	require.Equal(t, bz, bz2)
	require.Equal(t, tx.Hash(), tx2.Hash())
}

func TestEncodingTxNoSig(t *testing.T) {
	tx, _ := GenerateTestSendTx()
	tx.SetPublicKey(nil)
	tx.SetSignature(nil)

	bz, err := tx.MarshalCBOR()
	require.NoError(t, err)
	var tx2 Tx
	err = tx2.UnmarshalCBOR(bz)

	bz2, _ := tx2.MarshalCBOR()

	fmt.Printf("%x\n", bz)
	fmt.Printf("%x\n", bz2)

	require.NoError(t, err)
	require.Equal(t, bz, bz2)
	require.Equal(t, tx.Hash(), tx2.Hash())
}

func TestTxSanityCheck(t *testing.T) {
	tx, _ := GenerateTestSendTx()
	tx.data.Version = 2
	assert.Error(t, tx.SanityCheck())
	tx.data.Version = 1
	tx.data.Sequence = -1
	assert.Error(t, tx.SanityCheck())
}

func TestInvalidSignature(t *testing.T) {
	tx, pv := GenerateTestSendTx()
	assert.NoError(t, tx.SanityCheck())

	fmt.Printf("%x\n", tx.SignBytes())

	tx.SetPublicKey(nil)
	assert.Error(t, tx.SanityCheck())

	_, pbInv, pvInv := crypto.GenerateTestKeyPair()
	tx.SetPublicKey(&pbInv)
	assert.Error(t, tx.SanityCheck())

	sig := pvInv.Sign(tx.SignBytes())
	tx.SetSignature(sig)
	assert.Error(t, tx.SanityCheck())

	// Invalid sign Bytes
	var tx2 Tx
	tx2 = *tx
	tx2.data.Memo = "Hack me"
	sig = pv.Sign(tx2.SignBytes())
	pb := pv.PublicKey()
	tx.SetPublicKey(&pb)
	tx.SetSignature(sig)
	assert.Error(t, tx.SanityCheck())
}
