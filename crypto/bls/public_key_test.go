package bls_test

import (
	"encoding/hex"
	"strings"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestPublicKeyCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandBLSKeyPair()
	pub2 := new(bls.PublicKey)

	bs, err := pub1.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, pub2.UnmarshalCBOR(bs))
	assert.True(t, pub1.EqualsTo(pub2))

	assert.Error(t, pub2.UnmarshalCBOR([]byte("abcd")))

	inv, _ := hex.DecodeString(strings.Repeat("ff", bls.PublicKeySize))
	data, _ := cbor.Marshal(inv)
	assert.Error(t, pub2.UnmarshalCBOR(data))
}

func TestPublicKeyEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandBLSKeyPair()
	pub2, _ := ts.RandBLSKeyPair()

	assert.True(t, pub1.EqualsTo(pub1))
	assert.False(t, pub1.EqualsTo(pub2))
	assert.Equal(t, pub1, pub1)
	assert.NotEqual(t, pub1, pub2)
}

func TestPublicKeyEncoding(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandBLSKeyPair()
	w1 := util.NewFixedWriter(20)
	assert.Error(t, pub.Encode(w1))

	w2 := util.NewFixedWriter(bls.PublicKeySize)
	assert.NoError(t, pub.Encode(w2))

	r1 := util.NewFixedReader(20, w2.Bytes())
	assert.Error(t, pub.Decode(r1))

	r2 := util.NewFixedReader(bls.PublicKeySize, w2.Bytes())
	assert.NoError(t, pub.Decode(r2))
}

func TestPublicKeyVerifyAddress(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandBLSKeyPair()
	pub2, _ := ts.RandBLSKeyPair()

	err := pub1.VerifyAddress(pub1.AccountAddress())
	assert.NoError(t, err)
	err = pub1.VerifyAddress(pub2.AccountAddress())
	assert.Equal(t, err, crypto.AddressMismatchError{
		Expected: pub1.AccountAddress(),
		Got:      pub2.AccountAddress(),
	})
}

func TestNilPublicKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub := &bls.PublicKey{}
	assert.Error(t, pub.VerifyAddress(ts.RandAccAddress()))
	assert.Error(t, pub.Verify(nil, nil))
	assert.Error(t, pub.Verify(nil, &bls.Signature{}))
}

func TestNilSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandBLSKeyPair()
	assert.Error(t, pub.Verify(nil, nil))
	assert.Error(t, pub.Verify(nil, &bls.Signature{}))
}

func TestPublicKeyBytes(t *testing.T) {
	tests := []struct {
		errMsg  string
		encoded string
		valid   bool
		result  []byte
	}{
		{
			"invalid separator index -1",
			"not_proper_encoded",
			false, nil,
		},
		{
			"invalid bech32 string length 0",
			"",
			false, nil,
		},
		{
			"invalid character not part of charset: 105",
			"public1ioiooi",
			false, nil,
		},
		{
			"invalid bech32 string length 0",
			"public134jkgz",
			false, nil,
		},
		{
			"compression flag must be set",
			"public1pqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq" +
				"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqjzu9w8",
			false, nil,
		},
		{
			"input string must be zero when infinity flag is set",
			"public1pllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll" +
				"llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllluhpuzyf",
			false, nil,
		},
		{
			"public key is zero",
			"public1pcqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq" +
				"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqglnhh9",
			false, nil,
		},
		{
			"invalid hrp: xxx",
			"xxx1p4u8hfytl2pj6l9rj0t54gxcdmna4hq52ncqkkqjf3arha5mlk3x4mzpyjkhmdl20jae7f65aamjrvqc" +
				"vf4sudcapz52ctcwc8r9wz3z2gwxs38880cgvfy49ta5ssyjut05myd4zgmjqstggmetyuyg7v5evslaq",
			false, nil,
		},
		{
			"invalid checksum (expected jhx47a got jhx470)",
			"public1p4u8hfytl2pj6l9rj0t54gxcdmna4hq52ncqkkqjf3arha5mlk3x4mzpyjkhmdl20jae7f65aamjr" +
				"vqcvf4sudcapz52ctcwc8r9wz3z2gwxs38880cgvfy49ta5ssyjut05myd4zgmjqstggmetyuyg7v5jhx470",
			false, nil,
		},
		{
			"public key should be 96 bytes, but it is 95 bytes",
			"public1p4u8hfytl2pj6l9rj0t54gxcdmna4hq52ncqkkqjf3arha5mlk3x4mzpyjkhmdl20jae7f65aamjr" +
				"vqcvf4sudcapz52ctcwc8r9wz3z2gwxs38880cgvfy49ta5ssyjut05myd4zgmjqstggmetyuyg73y98kl",
			false,
			nil,
		},
		{
			"",
			"public1p4u8hfytl2pj6l9rj0t54gxcdmna4hq52ncqkkqjf3arha5mlk3x4mzpyjkhmdl20jae7f65aamjr" +
				"vqcvf4sudcapz52ctcwc8r9wz3z2gwxs38880cgvfy49ta5ssyjut05myd4zgmjqstggmetyuyg7v5jhx47a",
			true,
			[]byte{
				0xaf, 0x0f, 0x74, 0x91, 0x7f, 0x50, 0x65, 0xaf, 0x94, 0x72, 0x7a, 0xe9, 0x54, 0x1b, 0x0d, 0xdc,
				0xfb, 0x5b, 0x82, 0x8a, 0x9e, 0x01, 0x6b, 0x02, 0x49, 0x8f, 0x47, 0x7e, 0xd3, 0x7f, 0xb4, 0x4d, 0x5d,
				0x88, 0x24, 0x95, 0xaf, 0xb6, 0xfd, 0x4f, 0x97, 0x73, 0xe4, 0xea, 0x9d, 0xee, 0xe4, 0x36, 0x03, 0x0c,
				0x4d, 0x61, 0xc6, 0xe3, 0xa1, 0x15, 0x15, 0x85, 0xe1, 0xd8, 0x38, 0xca, 0xe1, 0x44, 0x4a, 0x43, 0x8d,
				0x08, 0x9c, 0xe7, 0x7e, 0x10, 0xc4, 0x92, 0xa5, 0x5f, 0x69, 0x08, 0x12, 0x5c, 0x5b, 0xe9, 0xb2, 0x36,
				0xa2, 0x46, 0xe4, 0x08, 0x2d, 0x08, 0xde, 0x56, 0x4e, 0x11, 0x1e, 0x65,
			},
		},
	}

	for no, test := range tests {
		pub, err := bls.PublicKeyFromString(test.encoded)
		if test.valid {
			assert.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, pub.Bytes(), test.result, "test %v: invalid bytes", no)
			assert.Equal(t, pub.String(), test.encoded, "test %v: invalid encoded", no)
		} else {
			assert.Equal(t, errors.Code(err), errors.ErrInvalidPublicKey)
			assert.Contains(t, err.Error(), test.errMsg, "test %v: error not matched", no)
		}
	}
}
