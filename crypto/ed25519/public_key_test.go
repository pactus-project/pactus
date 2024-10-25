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

func TestPublicKeyCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandEd25519KeyPair()
	pub2 := new(ed25519.PublicKey)

	bs, err := pub1.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, pub2.UnmarshalCBOR(bs))
	assert.True(t, pub1.EqualsTo(pub2))

	assert.Error(t, pub2.UnmarshalCBOR([]byte("abcd")))

	inv, _ := hex.DecodeString(strings.Repeat("ff", ed25519.PublicKeySize))
	data, _ := cbor.Marshal(inv)
	assert.NoError(t, pub2.UnmarshalCBOR(data))
}

func TestPublicKeyEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandEd25519KeyPair()
	pub2, _ := ts.RandEd25519KeyPair()
	pub3, _ := ts.RandBLSKeyPair()

	assert.True(t, pub1.EqualsTo(pub1))
	assert.False(t, pub1.EqualsTo(pub2))
	assert.False(t, pub1.EqualsTo(pub3))
}

func TestPublicKeyEncoding(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandEd25519KeyPair()
	fw1 := util.NewFixedWriter(20)
	assert.Error(t, pub.Encode(fw1))

	fw2 := util.NewFixedWriter(ed25519.PublicKeySize)
	assert.NoError(t, pub.Encode(fw2))

	fr1 := util.NewFixedReader(20, fw2.Bytes())
	assert.Error(t, pub.Decode(fr1))

	fr2 := util.NewFixedReader(ed25519.PublicKeySize, fw2.Bytes())
	assert.NoError(t, pub.Decode(fr2))
	assert.Equal(t, ed25519.PublicKeySize, pub.SerializeSize())
}

func TestPublicKeyVerifyAddress(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandEd25519KeyPair()
	pub2, _ := ts.RandEd25519KeyPair()

	err := pub1.VerifyAddress(pub1.AccountAddress())
	assert.NoError(t, err)

	err = pub1.VerifyAddress(pub2.AccountAddress())
	assert.Equal(t, crypto.AddressMismatchError{
		Expected: pub1.AccountAddress(),
		Got:      pub2.AccountAddress(),
	}, err)
}

func TestNilPublicKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub := &ed25519.PublicKey{}
	randSig := ts.RandEd25519Signature()
	assert.Error(t, pub.VerifyAddress(ts.RandAccAddress()))
	assert.Error(t, pub.VerifyAddress(ts.RandValAddress()))
	assert.Error(t, pub.Verify(nil, nil))
	assert.Panics(t, func() { _ = pub.Verify(nil, randSig) })
}

func TestNilSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandEd25519KeyPair()
	randSig := ts.RandEd25519Signature()
	assert.Error(t, pub.Verify(nil, nil))
	assert.Error(t, pub.Verify(nil, randSig))
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
			"invalid checksum (expected ztd56p got ztd5p6)",
			"public1rafnl324uwngqdq455ax4e52fedmfcvskkwas6wsau0u0nwj4g96qztd5p6",
			false, nil,
		},
		{
			"invalid HRP: xxx",
			"xxx1rafnl324uwngqdq455ax4e52fedmfcvskkwas6wsau0u0nwj4g96qvguamu",
			false, nil,
		},
		{
			"invalid length: 31",
			"public1ruwz86xyvhyehy8g7wg98jsmy07cfkjp6dy8zwxa8hqtdj99hquk7xyus",
			false, nil,
		},
		{
			"invalid signature type: 4",
			"public1yafnl324uwngqdq455ax4e52fedmfcvskkwas6wsau0u0nwj4g96qdnx0mf",
			false, nil,
		},
		{
			"",
			"public1rafnl324uwngqdq455ax4e52fedmfcvskkwas6wsau0u0nwj4g96qztd56p",
			true,
			[]byte{
				0xea, 0x67, 0xf8, 0xaa, 0xbc, 0x74, 0xd0, 0x06, 0x82, 0xb4, 0xa7, 0x4d, 0x5c, 0xd1, 0x49, 0xcb,
				0x76, 0x9c, 0x32, 0x16, 0xb3, 0xbb, 0x0d, 0x3a, 0x1d, 0xe3, 0xf8, 0xf9, 0xba, 0x55, 0x41, 0x74,
			},
		},
	}

	for no, tt := range tests {
		pub, err := ed25519.PublicKeyFromString(tt.encoded)
		if tt.valid {
			assert.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, tt.result, pub.Bytes(), "test %v: invalid bytes", no)
			assert.Equal(t, tt.encoded, pub.String(), "test %v: invalid encoded", no)
		} else {
			assert.Contains(t, err.Error(), tt.errMsg, "test %v: error not matched", no)
		}
	}
}
