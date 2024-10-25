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
	assert.NoError(t, pub2.UnmarshalCBOR(data))

	_, err = pub2.PointG2()
	assert.Error(t, err)
}

func TestPublicKeyEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandBLSKeyPair()
	pub2, _ := ts.RandBLSKeyPair()
	pub3, _ := ts.RandEd25519KeyPair()

	assert.True(t, pub1.EqualsTo(pub1))
	assert.False(t, pub1.EqualsTo(pub2))
	assert.False(t, pub1.EqualsTo(pub3))
}

func TestPublicKeyEncoding(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandBLSKeyPair()
	fw1 := util.NewFixedWriter(20)
	assert.Error(t, pub.Encode(fw1))

	fw2 := util.NewFixedWriter(bls.PublicKeySize)
	assert.NoError(t, pub.Encode(fw2))

	fr1 := util.NewFixedReader(20, fw2.Bytes())
	assert.Error(t, pub.Decode(fr1))

	fr2 := util.NewFixedReader(bls.PublicKeySize, fw2.Bytes())
	assert.NoError(t, pub.Decode(fr2))
	assert.Equal(t, bls.PublicKeySize, pub.SerializeSize())
}

func TestPublicKeyVerifyAddress(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandBLSKeyPair()
	pub2, _ := ts.RandBLSKeyPair()

	err := pub1.VerifyAddress(pub1.AccountAddress())
	assert.NoError(t, err)
	err = pub1.VerifyAddress(pub1.ValidatorAddress())
	assert.NoError(t, err)

	err = pub1.VerifyAddress(pub2.AccountAddress())
	assert.Equal(t, crypto.AddressMismatchError{
		Expected: pub1.AccountAddress(),
		Got:      pub2.AccountAddress(),
	}, err)

	err = pub1.VerifyAddress(pub2.ValidatorAddress())
	assert.Equal(t, crypto.AddressMismatchError{
		Expected: pub1.ValidatorAddress(),
		Got:      pub2.ValidatorAddress(),
	}, err)
}

func TestNilPublicKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub := &bls.PublicKey{}
	randSig := ts.RandBLSSignature()
	assert.Error(t, pub.VerifyAddress(ts.RandAccAddress()))
	assert.Error(t, pub.VerifyAddress(ts.RandValAddress()))
	assert.Error(t, pub.Verify(nil, nil))
	assert.Error(t, pub.Verify(nil, &bls.Signature{}))
	assert.Error(t, pub.Verify(nil, randSig))
}

func TestNilSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandBLSKeyPair()
	assert.Error(t, pub.Verify(nil, nil))
	assert.Error(t, pub.Verify(nil, &bls.Signature{}))
}

func TestPublicKeyFromString(t *testing.T) {
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
			"invalid HRP: xxx",
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
			"invalid length: 95",
			"public1p4u8hfytl2pj6l9rj0t54gxcdmna4hq52ncqkkqjf3arha5mlk3x4mzpyjkhmdl20jae7f65aamjr" +
				"vqcvf4sudcapz52ctcwc8r9wz3z2gwxs38880cgvfy49ta5ssyjut05myd4zgmjqstggmetyuyg73y98kl",
			false, nil,
		},
		{
			"invalid signature type: 2",
			"public1z372l5frmm5e7cn7ewfjdkx5t7y62kztqr82rtatat70cl8p8ng3rdzr02mzpwcfl6s2v26kry6mwg" +
				"xpqy92ywx9wtff80mc9p3kr4cmhgekj048gavx2zdh78tsnh7eg5jzdw6s3et6c0dqyp22vslcgkukxh4l4",
			false, nil,
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

	for no, tt := range tests {
		pub, err := bls.PublicKeyFromString(tt.encoded)
		if tt.valid {
			assert.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, tt.result, pub.Bytes(), "test %v: invalid bytes", no)
			assert.Equal(t, tt.encoded, pub.String(), "test %v: invalid encoded", no)
		} else {
			assert.Contains(t, err.Error(), tt.errMsg, "test %v: error not matched", no)
		}
	}
}

func TestPointG2(t *testing.T) {
	tests := []struct {
		errMsg  string
		encoded string
		valid   bool
	}{
		{
			"short buffer",
			"public1pqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq" +
				"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqjzu9w8",
			false,
		},
		{
			"invalid point encoding",
			"public1pllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll" +
				"llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllluhpuzyf",
			false,
		},
		{
			"invalid public key",
			"public1pcqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq" +
				"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqglnhh9",
			false,
		},
		{
			"",
			"public1p4u8hfytl2pj6l9rj0t54gxcdmna4hq52ncqkkqjf3arha5mlk3x4mzpyjkhmdl20jae7f65aamjr" +
				"vqcvf4sudcapz52ctcwc8r9wz3z2gwxs38880cgvfy49ta5ssyjut05myd4zgmjqstggmetyuyg7v5jhx47a",
			true,
		},
	}

	for no, tt := range tests {
		pub, err := bls.PublicKeyFromString(tt.encoded)
		require.NoError(t, err)

		_, err = pub.PointG2()
		if tt.valid {
			assert.NoError(t, err, "test %v: unexpected error", no)
		} else {
			assert.Contains(t, err.Error(), tt.errMsg, "test %v: error not matched", no)
		}
	}
}
