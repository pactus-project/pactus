package block_test

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

// A block in a CBOR message has no byte after it: the extra bytes would be a second encoding.
func TestUnmarshalCBORRejectsTrailingBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blk, _ := ts.GenerateTestBlock(ts.RandHeight())
	data, err := blk.Bytes()
	require.NoError(t, err)

	valid, err := cbor.Marshal(data)
	require.NoError(t, err)
	require.NoError(t, new(block.Block).UnmarshalCBOR(valid))

	trailing, err := cbor.Marshal(append(data, 0))
	require.NoError(t, err)
	require.ErrorIs(t, new(block.Block).UnmarshalCBOR(trailing), block.ErrTrailingBytes)
}
