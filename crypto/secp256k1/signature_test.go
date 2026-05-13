package secp256k1_test

import (
	"encoding/hex"
	"strings"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/secp256k1"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignatureCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv := ts.RandSecp256k1KeyPair()
	sig1 := prv.Sign(ts.RandBytes(16))
	sig2 := new(secp256k1.Signature)

	bs, err := sig1.MarshalCBOR()
	require.NoError(t, err)
	require.NoError(t, sig2.UnmarshalCBOR(bs))
	assert.True(t, sig1.EqualsTo(sig2))

	require.Error(t, sig2.UnmarshalCBOR([]byte("abcd")))

	inv, _ := hex.DecodeString(strings.Repeat("ff", secp256k1.SignatureSize))
	data, _ := cbor.Marshal(inv)
	require.NoError(t, sig2.UnmarshalCBOR(data))
}

func TestSignatureEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv1 := ts.RandSecp256k1KeyPair()
	_, prv2 := ts.RandSecp256k1KeyPair()
	_, prv3 := ts.RandEd25519KeyPair()

	sig1 := prv1.Sign([]byte("foo"))
	sig2 := prv2.Sign([]byte("foo"))
	sig3 := prv3.Sign([]byte("foo"))

	assert.True(t, sig1.EqualsTo(sig1))
	assert.False(t, sig1.EqualsTo(sig2))
	assert.False(t, sig1.EqualsTo(sig3))
}

func TestSignatureEncoding(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv := ts.RandSecp256k1KeyPair()
	sig := prv.Sign(ts.RandBytes(16))
	fw1 := util.NewFixedWriter(20)
	require.Error(t, sig.Encode(fw1))

	fw2 := util.NewFixedWriter(secp256k1.SignatureSize)
	require.NoError(t, sig.Encode(fw2))

	fr1 := util.NewFixedReader(20, fw2.Bytes())
	require.Error(t, sig.Decode(fr1))

	fr2 := util.NewFixedReader(secp256k1.SignatureSize, fw2.Bytes())
	require.NoError(t, sig.Decode(fr2))
	assert.Equal(t, secp256k1.SignatureSize, sig.SerializeSize())
}

func TestVerifyingSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	msg := []byte("pactus")

	pb1, pv1 := ts.RandSecp256k1KeyPair()
	pb2, pv2 := ts.RandSecp256k1KeyPair()
	sig1 := pv1.Sign(msg)
	sig2 := pv2.Sign(msg)

	assert.False(t, sig1.EqualsTo(sig2))
	require.NoError(t, pb1.Verify(msg, sig1))
	require.NoError(t, pb2.Verify(msg, sig2))
	require.ErrorIs(t, pb1.Verify(msg, sig2), crypto.ErrInvalidSignature)
	require.ErrorIs(t, pb2.Verify(msg, sig1), crypto.ErrInvalidSignature)
	require.ErrorIs(t, pb1.Verify(msg[1:], sig1), crypto.ErrInvalidSignature)
}

func TestSignatureFromString(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	_, prv := ts.RandSecp256k1KeyPair()
	validHex := hex.EncodeToString(prv.Sign([]byte("test-vector")).Bytes())

	tests := []struct {
		errMsg  string
		encoded string
		valid   bool
		bytes   []byte
	}{
		{
			"encoding/hex: invalid byte",
			"not_proper_encoded",
			false, nil,
		},
		{
			"invalid length: 0",
			"",
			false, nil,
		},
		{
			"invalid length: 32",
			strings.Repeat("00", 32),
			false, nil,
		},
		{
			"",
			validHex,
			true,
			prv.Sign([]byte("test-vector")).Bytes(),
		},
	}

	for no, tt := range tests {
		sig, err := secp256k1.SignatureFromString(tt.encoded)
		if tt.valid {
			require.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, tt.bytes, sig.Bytes(), "test %v: invalid bytes", no)
			assert.Equal(t, tt.encoded, sig.String(), "test %v: invalid encode", no)
		} else {
			require.Error(t, err, "test %v", no)
			assert.Contains(t, err.Error(), tt.errMsg, "test %v: error not matched", no)
		}
	}
}
