package secp256k1_test

import (
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/secp256k1"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

// The CBOR bytes of a key or a signature must have exactly the expected length:
// the bytes in excess would be silently dropped.
func TestUnmarshalCBORExactLength(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandSecp256k1KeyPair()
	sig := ts.RandSecp256k1Signature()

	unmarshalPub := func(data []byte) error {
		bs, err := cbor.Marshal(data)
		require.NoError(t, err)

		return new(secp256k1.PublicKey).UnmarshalCBOR(bs)
	}
	unmarshalSig := func(data []byte) error {
		bs, err := cbor.Marshal(data)
		require.NoError(t, err)

		return new(secp256k1.Signature).UnmarshalCBOR(bs)
	}

	pubData, sigData := pub.Bytes(), sig.Bytes()

	require.NoError(t, unmarshalPub(pubData))
	require.NoError(t, unmarshalSig(sigData))

	require.ErrorIs(t, unmarshalPub(append(append([]byte{}, pubData...), 0)), crypto.InvalidLengthError(len(pubData)+1))
	require.ErrorIs(t, unmarshalSig(append(append([]byte{}, sigData...), 0)), crypto.InvalidLengthError(len(sigData)+1))
	require.Error(t, unmarshalPub(pubData[:len(pubData)-1]))
	require.Error(t, unmarshalSig(sigData[:len(sigData)-1]))
}
