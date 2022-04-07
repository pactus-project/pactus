package bls

import (
	"encoding/hex"
	"strings"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

func TestSignatureCBORMarshaling(t *testing.T) {
	_, prv := GenerateTestKeyPair()
	sig1 := prv.Sign(util.Uint64ToSlice(util.RandUint64(0)))
	sig2 := new(Signature)

	bs, err := sig1.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, sig2.UnmarshalCBOR(bs))
	assert.True(t, sig1.EqualsTo(sig2))
	assert.NoError(t, sig1.SanityCheck())

	assert.Error(t, sig2.UnmarshalCBOR([]byte("abcd")))

	inv, _ := hex.DecodeString(strings.Repeat("ff", SignatureSize))
	data, _ := cbor.Marshal(inv)
	assert.Error(t, sig2.UnmarshalCBOR(data))
}

func TestSignatureEncoding(t *testing.T) {
	_, prv := GenerateTestKeyPair()
	sig := prv.Sign(util.Uint64ToSlice(util.RandUint64(0)))
	w1 := util.NewFixedWriter(20)
	assert.Error(t, sig.Encode(w1))

	w2 := util.NewFixedWriter(SignatureSize)
	assert.NoError(t, sig.Encode(w2))

	r1 := util.NewFixedReader(20, w2.Bytes())
	assert.Error(t, sig.Decode(r1))

	r2 := util.NewFixedReader(SignatureSize, w2.Bytes())
	assert.NoError(t, sig.Decode(r2))
}

func TestVerifyingSignature(t *testing.T) {
	msg := []byte("zarb")

	pb1, pv1 := GenerateTestKeyPair()
	pb2, pv2 := GenerateTestKeyPair()
	sig1 := pv1.Sign(msg)
	sig2 := pv2.Sign(msg)

	assert.NotEqual(t, sig1, sig2)
	assert.NoError(t, pb1.Verify(msg, sig1))
	assert.NoError(t, pb2.Verify(msg, sig2))
	assert.Equal(t, errors.Code(pb1.Verify(msg, sig2)), errors.ErrInvalidSignature)
	assert.Equal(t, errors.Code(pb2.Verify(msg, sig1)), errors.ErrInvalidSignature)
	assert.Equal(t, errors.Code(pb1.Verify(msg[1:], sig1)), errors.ErrInvalidSignature)
}

func TestSignatureBytes(t *testing.T) {
	tests := []struct {
		name      string
		encoded   string
		bytes     []byte
		decodable bool
		valid     bool
	}{
		{"invalid input",
			"inv",
			nil,
			false, false},
		{"nil signature",
			"",
			nil,
			false, false},
		{"invalid length",
			"00",
			nil,
			false, false},
		{"zero signature",
			"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			nil,
			false, false},
		{"invalid signature",
			"fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			nil,
			false, false},
		{"infinite signature",
			"c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			[]byte{0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			true, false,
		},
		{"valid signature",
			"ad0f88cec815e9b8af3f0136297cb242ed8b6369af723fbdac077fa927f5780db7df47c77fb53f3a22324673f000c792",
			[]byte{0xad, 0x0f, 0x88, 0xce, 0xc8, 0x15, 0xe9, 0xb8, 0xaf, 0x3f, 0x01, 0x36, 0x29, 0x7c, 0xb2, 0x42,
				0xed, 0x8b, 0x63, 0x69, 0xaf, 0x72, 0x3f, 0xbd, 0xac, 0x07, 0x7f, 0xa9, 0x27, 0xf5, 0x78, 0x0d,
				0xb7, 0xdf, 0x47, 0xc7, 0x7f, 0xb5, 0x3f, 0x3a, 0x22, 0x32, 0x46, 0x73, 0xf0, 0x00, 0xc7, 0x92},
			true, true,
		},
	}

	for _, test := range tests {
		sig, err := SignatureFromString(test.encoded)
		if test.decodable {
			assert.NoError(t, err, "test %v. unexpected error", test.name)
			assert.Equal(t, sig.SanityCheck() == nil, test.valid, "test %v. sanity check failed", test.name)
			assert.Equal(t, sig.Bytes(), test.bytes, "test %v. invalid bytes", test.name)
			assert.Equal(t, sig.String(), test.encoded, "test %v. invalid encoded", test.name)

		} else {
			assert.Error(t, err, "test %v. should failed", test.name)
			assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature)
		}
	}
}
