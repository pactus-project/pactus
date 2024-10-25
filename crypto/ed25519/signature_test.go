package ed25519_test

import (
	"encoding/hex"
	"strings"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestSignatureCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv := ts.RandEd25519KeyPair()
	sig1 := prv.Sign(ts.RandBytes(16))
	sig2 := new(ed25519.Signature)

	bs, err := sig1.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, sig2.UnmarshalCBOR(bs))
	assert.True(t, sig1.EqualsTo(sig2))

	assert.Error(t, sig2.UnmarshalCBOR([]byte("abcd")))

	inv, _ := hex.DecodeString(strings.Repeat("ff", ed25519.SignatureSize))
	data, _ := cbor.Marshal(inv)
	assert.NoError(t, sig2.UnmarshalCBOR(data))
}

func TestSignatureEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv1 := ts.RandEd25519KeyPair()
	_, prv2 := ts.RandEd25519KeyPair()
	_, prv3 := ts.RandBLSKeyPair()

	sig1 := prv1.Sign([]byte("foo"))
	sig2 := prv2.Sign([]byte("foo"))
	sig3 := prv3.Sign([]byte("foo"))

	assert.True(t, sig1.EqualsTo(sig1))
	assert.False(t, sig1.EqualsTo(sig2))
	assert.False(t, sig1.EqualsTo(sig3))
}

func TestSignatureEncoding(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv := ts.RandEd25519KeyPair()
	sig := prv.Sign(ts.RandBytes(16))
	fw1 := util.NewFixedWriter(20)
	assert.Error(t, sig.Encode(fw1))

	fw2 := util.NewFixedWriter(ed25519.SignatureSize)
	assert.NoError(t, sig.Encode(fw2))

	fr1 := util.NewFixedReader(20, fw2.Bytes())
	assert.Error(t, sig.Decode(fr1))

	fr2 := util.NewFixedReader(ed25519.SignatureSize, fw2.Bytes())
	assert.NoError(t, sig.Decode(fr2))
	assert.Equal(t, ed25519.SignatureSize, sig.SerializeSize())
}

func TestVerifyingSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	msg := []byte("pactus")

	pb1, pv1 := ts.RandEd25519KeyPair()
	pb2, pv2 := ts.RandEd25519KeyPair()
	sig1 := pv1.Sign(msg)
	sig2 := pv2.Sign(msg)

	assert.False(t, sig1.EqualsTo(sig2))
	assert.NoError(t, pb1.Verify(msg, sig1))
	assert.NoError(t, pb2.Verify(msg, sig2))
	assert.ErrorIs(t, pb1.Verify(msg, sig2), crypto.ErrInvalidSignature)
	assert.ErrorIs(t, pb2.Verify(msg, sig1), crypto.ErrInvalidSignature)
	assert.ErrorIs(t, pb1.Verify(msg[1:], sig1), crypto.ErrInvalidSignature)
}

func TestSignatureFromString(t *testing.T) {
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
			"invalid length: 0",
			"",
			false, nil,
		},
		{
			"",
			"7d6af02f788422319781b03d7f4ed575b78c4c4dc8060ce145624fc8dc9ad92b" +
				"ae2d28c70242f03a644f313009ad9cc88b5dc37d501e43279c8fbc40b973af04",
			true,
			[]byte{
				0x7d, 0x6a, 0xf0, 0x2f, 0x78, 0x84, 0x22, 0x31, 0x97, 0x81, 0xb0, 0x3d, 0x7f, 0x4e, 0xd5, 0x75,
				0xb7, 0x8c, 0x4c, 0x4d, 0xc8, 0x06, 0x0c, 0xe1, 0x45, 0x62, 0x4f, 0xc8, 0xdc, 0x9a, 0xd9, 0x2b,
				0xae, 0x2d, 0x28, 0xc7, 0x02, 0x42, 0xf0, 0x3a, 0x64, 0x4f, 0x31, 0x30, 0x09, 0xad, 0x9c, 0xc8,
				0x8b, 0x5d, 0xc3, 0x7d, 0x50, 0x1e, 0x43, 0x27, 0x9c, 0x8f, 0xbc, 0x40, 0xb9, 0x73, 0xaf, 0x04,
			},
		},
	}

	for no, tt := range tests {
		sig, err := ed25519.SignatureFromString(tt.encoded)
		if tt.valid {
			assert.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, tt.bytes, sig.Bytes(), "test %v: invalid bytes", no)
			assert.Equal(t, tt.encoded, sig.String(), "test %v: invalid encode", no)
		} else {
			assert.Contains(t, err.Error(), tt.errMsg, "test %v: error not matched", no)
		}
	}
}
