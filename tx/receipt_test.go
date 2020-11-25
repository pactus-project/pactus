package tx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestEncoding(t *testing.T) {
	r1 := Receipt{
		data: receiptData{
			TxHash:    crypto.GenerateTestHash(),
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
