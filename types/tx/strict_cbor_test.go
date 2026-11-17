package tx_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

// selfDescribeTag is the CBOR tag 55799, which can prefix any item. The encoder never writes it,
// so the decoder of the data that comes from the network rejects it, like any other tag.
var selfDescribeTag = []byte{0xd9, 0xd9, 0xf7}

func TestUnmarshalCBORRejectsTags(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx := ts.GenerateTestTransferTx()
	data, err := trx.MarshalCBOR()
	require.NoError(t, err)

	require.NoError(t, new(tx.Tx).UnmarshalCBOR(data))
	require.Error(t, new(tx.Tx).UnmarshalCBOR(append(selfDescribeTag, data...)))
}
