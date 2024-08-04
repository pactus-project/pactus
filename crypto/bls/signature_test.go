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
	assert.NoError(t, err)
	assert.NoError(t, sig2.UnmarshalCBOR(bs))
	assert.True(t, sig1.EqualsTo(sig2))

	assert.Error(t, sig2.UnmarshalCBOR([]byte("abcd")))

	inv, _ := hex.DecodeString(strings.Repeat("ff", bls.SignatureSize))
	data, _ := cbor.Marshal(inv)
	assert.NoError(t, sig2.UnmarshalCBOR(data))

	_, err = sig2.PointG1()
	assert.Error(t, err)
}

func TestSignatureEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv := ts.RandBLSKeyPair()
	sig1 := prv.Sign([]byte("foo"))
	sig2 := prv.Sign([]byte("bar"))

	assert.True(t, sig1.EqualsTo(sig1))
	assert.False(t, sig1.EqualsTo(sig2))
	assert.Equal(t, sig1, sig1)
	assert.NotEqual(t, sig1, sig2)
}

func TestSignatureEncoding(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv := ts.RandBLSKeyPair()
	sig := prv.Sign(ts.RandBytes(16))
	w1 := util.NewFixedWriter(20)
	assert.Error(t, sig.Encode(w1))

	w2 := util.NewFixedWriter(bls.SignatureSize)
	assert.NoError(t, sig.Encode(w2))

	r1 := util.NewFixedReader(20, w2.Bytes())
	assert.Error(t, sig.Decode(r1))

	r2 := util.NewFixedReader(bls.SignatureSize, w2.Bytes())
	assert.NoError(t, sig.Decode(r2))
}

func TestVerifyingSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	msg := []byte("zarb")

	pb1, pv1 := ts.RandBLSKeyPair()
	pb2, pv2 := ts.RandBLSKeyPair()
	sig1 := pv1.Sign(msg)
	sig2 := pv2.Sign(msg)

	assert.False(t, sig1.EqualsTo(sig2))
	assert.NoError(t, pb1.Verify(msg, sig1))
	assert.NoError(t, pb2.Verify(msg, sig2))
	assert.ErrorIs(t, pb1.Verify(msg, sig2), crypto.ErrInvalidSignature)
	assert.ErrorIs(t, pb2.Verify(msg, sig1), crypto.ErrInvalidSignature)
	assert.ErrorIs(t, pb1.Verify(msg[1:], sig1), crypto.ErrInvalidSignature)
}

func TestSignatureBytes(t *testing.T) {
	tests := []struct {
		errMsg  string
		encoded string
		valid   bool
		bytes   []byte
	}{
		{
			"encoding/hex: invalid byte: U+006E 'n'",
			"not_proper_encoded",
			false, nil,
		},
		{
			"signature should be 48 bytes, but it is 0 bytes",
			"",
			false, nil,
		},
		{
			"encoding/hex: odd length hex string",
			"0",
			false, nil,
		},
		{
			"signature should be 48 bytes, but it is 1 bytes",
			"00",
			false, nil,
		},
		{
			"",
			"ad0f88cec815e9b8af3f0136297cb242ed8b6369af723fbdac077fa927f5780db7df47c77fb53f3a22324673f000c792",
			true,
			[]byte{
				0xad, 0x0f, 0x88, 0xce, 0xc8, 0x15, 0xe9, 0xb8, 0xaf, 0x3f, 0x01, 0x36, 0x29, 0x7c, 0xb2, 0x42,
				0xed, 0x8b, 0x63, 0x69, 0xaf, 0x72, 0x3f, 0xbd, 0xac, 0x07, 0x7f, 0xa9, 0x27, 0xf5, 0x78, 0x0d,
				0xb7, 0xdf, 0x47, 0xc7, 0x7f, 0xb5, 0x3f, 0x3a, 0x22, 0x32, 0x46, 0x73, 0xf0, 0x00, 0xc7, 0x92,
			},
		},
	}

	for no, test := range tests {
		sig, err := bls.SignatureFromString(test.encoded)
		if test.valid {
			assert.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, test.bytes, sig.Bytes(), "test %v: invalid bytes", no)
			assert.Equal(t, test.encoded, sig.String(), "test %v: invalid encode", no)
		} else {
			assert.Contains(t, err.Error(), test.errMsg, "test %v: error not matched", no)
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
			"compression flag must be set",
			"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			false,
		},
		{
			"input string must be zero when infinity flag is set",
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			false,
		},
		{
			"signature is zero",
			"c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			false,
		},
		{
			"",
			"ad0f88cec815e9b8af3f0136297cb242ed8b6369af723fbdac077fa927f5780db7df47c77fb53f3a22324673f000c792",
			true,
		},
	}

	for no, test := range tests {
		sig, err := bls.SignatureFromString(test.encoded)
		require.NoError(t, err)

		_, err = sig.PointG1()
		if test.valid {
			assert.NoError(t, err, "test %v: unexpected error", no)
		} else {
			assert.Contains(t, err.Error(), test.errMsg, "test %v: error not matched", no)
		}
	}
}
