package bls

import (
	"encoding/hex"
	"strings"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/util/errors"
)

func TestPublicKeyCBORMarshaling(t *testing.T) {
	pub1, _ := GenerateTestKeyPair()
	pub2 := new(PublicKey)

	bs, err := pub1.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, pub2.UnmarshalCBOR(bs))
	assert.True(t, pub1.EqualsTo(pub2))
	assert.NoError(t, pub1.SanityCheck())

	assert.Error(t, pub2.UnmarshalCBOR([]byte("abcd")))

	inv, _ := hex.DecodeString(strings.Repeat("ff", PublicKeySize))
	data, _ := cbor.Marshal(inv)
	assert.Error(t, pub2.UnmarshalCBOR(data))
}

func TestPublicKeyEncoding(t *testing.T) {
	pub, _ := GenerateTestKeyPair()
	w1 := util.NewFixedWriter(20)
	assert.Error(t, pub.Encode(w1))

	w2 := util.NewFixedWriter(PublicKeySize)
	assert.NoError(t, pub.Encode(w2))

	r1 := util.NewFixedReader(20, w2.Bytes())
	assert.Error(t, pub.Decode(r1))

	r2 := util.NewFixedReader(PublicKeySize, w2.Bytes())
	assert.NoError(t, pub.Decode(r2))
}

func TestPublicKeyVerifyAddress(t *testing.T) {
	pub1, _ := GenerateTestKeyPair()
	pub2, _ := GenerateTestKeyPair()

	assert.NoError(t, pub1.VerifyAddress(pub1.Address()))
	assert.Equal(t, errors.Code(pub1.VerifyAddress(pub2.Address())), errors.ErrInvalidAddress)
}

func TestNilPublicKey(t *testing.T) {
	pub := &PublicKey{}
	assert.Error(t, pub.VerifyAddress(crypto.GenerateTestAddress()))
	assert.Error(t, pub.Verify(nil, nil))
	assert.Error(t, pub.Verify(nil, &Signature{}))
}

func TestNilSignature(t *testing.T) {
	pub, _ := GenerateTestKeyPair()
	assert.Error(t, pub.Verify(nil, nil))
	assert.Error(t, pub.Verify(nil, &Signature{}))
}

func TestPublicKeyBytes(t *testing.T) {
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
		{"nil public key",
			"",
			nil,
			false, false},
		{"invalid length",
			"00",
			nil,
			false, false},
		{"zero public key",
			"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			nil,
			false, false},
		{"invalid public key",
			"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			nil,
			false, false},
		{"infinite public key",
			"c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			[]byte{0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			true, false,
		},
		{"valid public key",
			"af0f74917f5065af94727ae9541b0ddcfb5b828a9e016b02498f477ed37fb44d5d882495afb6fd4f9773e4ea9deee436030c4d61c6e3a1151585e1d838cae1444a438d089ce77e10c492a55f6908125c5be9b236a246e4082d08de564e111e65",
			[]byte{0xaf, 0x0f, 0x74, 0x91, 0x7f, 0x50, 0x65, 0xaf, 0x94, 0x72, 0x7a, 0xe9, 0x54, 0x1b, 0x0d, 0xdc,
				0xfb, 0x5b, 0x82, 0x8a, 0x9e, 0x01, 0x6b, 0x02, 0x49, 0x8f, 0x47, 0x7e, 0xd3, 0x7f, 0xb4, 0x4d, 0x5d,
				0x88, 0x24, 0x95, 0xaf, 0xb6, 0xfd, 0x4f, 0x97, 0x73, 0xe4, 0xea, 0x9d, 0xee, 0xe4, 0x36, 0x03, 0x0c,
				0x4d, 0x61, 0xc6, 0xe3, 0xa1, 0x15, 0x15, 0x85, 0xe1, 0xd8, 0x38, 0xca, 0xe1, 0x44, 0x4a, 0x43, 0x8d,
				0x08, 0x9c, 0xe7, 0x7e, 0x10, 0xc4, 0x92, 0xa5, 0x5f, 0x69, 0x08, 0x12, 0x5c, 0x5b, 0xe9, 0xb2, 0x36,
				0xa2, 0x46, 0xe4, 0x08, 0x2d, 0x08, 0xde, 0x56, 0x4e, 0x11, 0x1e, 0x65},
			true, true,
		},
	}

	for _, test := range tests {
		pub, err := PublicKeyFromString(test.encoded)
		if test.decodable {
			assert.NoError(t, err, "test %v. unexpected error", test.name)
			assert.Equal(t, pub.SanityCheck() == nil, test.valid, "test %v. sanity check failed", test.name)
			assert.Equal(t, pub.Bytes(), test.bytes, "test %v. invalid bytes", test.name)
			assert.Equal(t, pub.String(), test.encoded, "test %v. invalid encoded", test.name)
		} else {
			assert.Error(t, err, "test %v. should failed", test.name)
			assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
		}
	}
}
