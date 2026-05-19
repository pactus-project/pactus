package bls_test

import (
	"encoding/hex"
	"strings"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignatureCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv := ts.RandBLSKeyPair()
	sig1 := prv.Sign(ts.RandBytes(16))
	sig2 := new(bls.Signature)

	bs, err := sig1.MarshalCBOR()
	require.NoError(t, err)
	require.NoError(t, sig2.UnmarshalCBOR(bs))
	assert.True(t, sig1.EqualsTo(sig2))

	require.Error(t, sig2.UnmarshalCBOR([]byte("abcd")))

	inv, _ := hex.DecodeString(strings.Repeat("ff", bls.SignatureSize))
	data, _ := cbor.Marshal(inv)
	require.NoError(t, sig2.UnmarshalCBOR(data))

	_, err = sig2.PointG1()
	require.Error(t, err)
}

func TestSignatureEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv1 := ts.RandBLSKeyPair()
	_, prv2 := ts.RandBLSKeyPair()
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

	_, prv := ts.RandBLSKeyPair()
	sig := prv.Sign(ts.RandBytes(16))
	fw1 := util.NewFixedWriter(20)
	require.Error(t, sig.Encode(fw1))

	fw2 := util.NewFixedWriter(bls.SignatureSize)
	require.NoError(t, sig.Encode(fw2))

	fr1 := util.NewFixedReader(20, fw2.Bytes())
	require.Error(t, sig.Decode(fr1))

	fr2 := util.NewFixedReader(bls.SignatureSize, fw2.Bytes())
	require.NoError(t, sig.Decode(fr2))
	assert.Equal(t, bls.SignatureSize, sig.SerializeSize())
}

func TestVerifyingSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	msg := []byte("pactus")

	pb1, pv1 := ts.RandBLSKeyPair()
	pb2, pv2 := ts.RandBLSKeyPair()
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
	tests := []struct {
		errMsg  string
		encoded string
		valid   bool
	}{
		{
			"encoding/hex: invalid byte: U+006E 'n'",
			"not_proper_encoded",
			false,
		},
		{
			"invalid length: 0",
			"",
			false,
		},
		{
			"encoding/hex: odd length hex string",
			"0",
			false,
		},
		{
			"invalid length: 1",
			"00",
			false,
		},
		{
			"",
			"ad0f88cec815e9b8af3f0136297cb242ed8b6369af723fbdac077fa927f5780db7df47c77fb53f3a22324673f000c792",
			true,
		},
	}

	for no, tt := range tests {
		sig, err := bls.SignatureFromString(tt.encoded)
		if tt.valid {
			require.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, tt.encoded, sig.String(), "test %v: invalid encode", no)
		} else {
			assert.Contains(t, err.Error(), tt.errMsg, "test %v: error not matched", no)
		}
	}
}

func TestPointG1(t *testing.T) {
	tests := []struct {
		errMsg  string
		encoded string
		valid   bool
	}{
		{
			"short buffer",
			"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			false,
		},
		{
			"invalid point encoding",
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			false,
		},
		{
			"invalid point",
			"c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			false,
		},
		{
			"",
			"ad0f88cec815e9b8af3f0136297cb242ed8b6369af723fbdac077fa927f5780db7df47c77fb53f3a22324673f000c792",
			true,
		},
	}

	for no, tt := range tests {
		sig, err := bls.SignatureFromString(tt.encoded)
		require.NoError(t, err)

		_, err = sig.PointG1()
		if tt.valid {
			require.NoError(t, err, "test %v: unexpected error", no)
		} else {
			require.Error(t, err, "test %v", no)
			assert.Contains(t, err.Error(), tt.errMsg, "test %v: error not matched", no)
		}
	}
}
