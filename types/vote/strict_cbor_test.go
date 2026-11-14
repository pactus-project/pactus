package vote_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

// selfDescribeTag is the CBOR tag 55799, which can prefix any item. The encoder never writes it,
// so the decoder of the data that comes from the network rejects it, like any other tag.
var selfDescribeTag = []byte{0xd9, 0xd9, 0xf7}

func TestUnmarshalCBORRejectsTags(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	vte, _ := ts.GenerateTestPrepareVote(ts.RandHeight(), ts.RandRound())
	data, err := vte.MarshalCBOR()
	require.NoError(t, err)

	decoded := new(vote.Vote)
	require.NoError(t, decoded.UnmarshalCBOR(data))
	require.Error(t, new(vote.Vote).UnmarshalCBOR(append(selfDescribeTag, data...)))
}
