package tx

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestEncodingReceipt(t *testing.T) {
	r1 := Receipt{
		data: receiptData{
			TxID:    crypto.GenerateTestHash(),
			BlockHash: crypto.GenerateTestHash(),
			Status:    Ok,
		},
	}

	bz, err := r1.MarshalCBOR()
	require.NoError(t, err)
	var r2 Receipt
	err = r2.UnmarshalCBOR(bz)
	fmt.Printf("%x\n", bz)

	require.NoError(t, err)
	require.Equal(t, r1, r2)
}

func TestSanityCheck(t *testing.T) {
	r := Receipt{
		data: receiptData{
			TxID:    crypto.GenerateTestHash(),
			BlockHash: crypto.GenerateTestHash(),
			Status:    Ok,
		},
	}

	assert.NoError(t, r.SanityCheck())
	r.data.Status = 1
	assert.Error(t, r.SanityCheck())
	r.data.Status = 0
	r.data.BlockHash = crypto.UndefHash
	assert.Error(t, r.SanityCheck())
}
func TestReceiptDecodingAndHash(t *testing.T) {
	d, _ := hex.DecodeString("a30100025820fa62c80a6e5a929d89acc2d5b169c47e2f12dd79b8ee9ccb209f38abaacc510f035820b3f91e81559252054698b20e658c25b2dd7b4f6e4cb928641921e7cef19de346")
	h, _ := crypto.HashFromString("5aef9dfba6969624095dd4eb593cd0212a1500f82d48a82c77f622941de5692b")
	var r Receipt
	err := r.UnmarshalCBOR(d)
	assert.NoError(t, err)
	d2, _ := r.MarshalCBOR()
	assert.Equal(t, d, d2)
	assert.Equal(t, r.Hash(), h)
}
