package secp256k1_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/secp256k1"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

// selfDescribeTag is the CBOR tag 55799, which can prefix any item. The encoder never writes it,
// so the decoder of the data that comes from the network rejects it, like any other tag.
var selfDescribeTag = []byte{0xd9, 0xd9, 0xf7}

func TestUnmarshalCBORRejectsTags(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, prv := ts.RandSecp256k1KeyPair()
	sig := prv.Sign([]byte("message"))

	pubData, err := pub.MarshalCBOR()
	require.NoError(t, err)
	require.NoError(t, new(secp256k1.PublicKey).UnmarshalCBOR(pubData))
	require.Error(t, new(secp256k1.PublicKey).UnmarshalCBOR(append(selfDescribeTag, pubData...)))

	sigData, err := sig.MarshalCBOR()
	require.NoError(t, err)
	require.NoError(t, new(secp256k1.Signature).UnmarshalCBOR(sigData))
	require.Error(t, new(secp256k1.Signature).UnmarshalCBOR(append(selfDescribeTag, sigData...)))
}
