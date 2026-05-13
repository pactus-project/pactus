package secp256k1_test

import (
	"encoding/hex"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/secp256k1"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicKeyCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandSecp256k1KeyPair()
	pub2 := new(secp256k1.PublicKey)

	bs, err := pub1.MarshalCBOR()
	require.NoError(t, err)
	require.NoError(t, pub2.UnmarshalCBOR(bs))
	assert.True(t, pub1.EqualsTo(pub2))

	require.Error(t, pub2.UnmarshalCBOR([]byte("abcd")))

	inv, _ := ts.RandSecp256k1KeyPair()
	invBytes := inv.Bytes()
	data, _ := cbor.Marshal(invBytes)
	require.NoError(t, pub2.UnmarshalCBOR(data))
}

func TestPublicKeyEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandSecp256k1KeyPair()
	pub2, _ := ts.RandSecp256k1KeyPair()
	pub3, _ := ts.RandEd25519KeyPair()

	assert.True(t, pub1.EqualsTo(pub1))
	assert.False(t, pub1.EqualsTo(pub2))
	assert.False(t, pub1.EqualsTo(pub3))
}

func TestPublicKeyEncoding(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandSecp256k1KeyPair()
	fw1 := util.NewFixedWriter(20)
	require.Error(t, pub.Encode(fw1))

	fw2 := util.NewFixedWriter(secp256k1.PublicKeySize)
	require.NoError(t, pub.Encode(fw2))

	fr1 := util.NewFixedReader(20, fw2.Bytes())
	require.Error(t, pub.Decode(fr1))

	fr2 := util.NewFixedReader(secp256k1.PublicKeySize, fw2.Bytes())
	require.NoError(t, pub.Decode(fr2))
	assert.Equal(t, secp256k1.PublicKeySize, pub.SerializeSize())
}

func TestPublicKeyVerifyAddress(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub1, _ := ts.RandSecp256k1KeyPair()
	pub2, _ := ts.RandSecp256k1KeyPair()

	err := pub1.VerifyAddress(pub1.AccountAddress())
	require.NoError(t, err)

	err = pub1.VerifyAddress(pub2.AccountAddress())
	assert.Equal(t, crypto.AddressMismatchError{
		Expected: pub1.AccountAddress(),
		Got:      pub2.AccountAddress(),
	}, err)
}

func TestNilPublicKey(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub := &secp256k1.PublicKey{}
	randSig := ts.RandSecp256k1Signature()
	require.Error(t, pub.VerifyAddress(ts.RandAccAddress()))
	require.Error(t, pub.VerifyAddress(ts.RandValAddress()))
	require.Error(t, pub.Verify(nil, nil))
	assert.Panics(t, func() { _ = pub.Verify(nil, randSig) })
}

func TestNilSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandSecp256k1KeyPair()
	randSig := ts.RandSecp256k1Signature()
	require.Error(t, pub.Verify(nil, nil))
	require.Error(t, pub.Verify(nil, randSig))
}

func TestPublicKeyFromString(t *testing.T) {
	pipPubBytes, err := hex.DecodeString("036d6caac248af96f6afa7f904f550253a0f3ef3f5aa2fe6838a95b216691468e2")
	require.NoError(t, err)

	wrongHRP, err := bech32m.EncodeFromBase256WithType("XXX", crypto.SignatureTypeSecp256k1, pipPubBytes)
	require.NoError(t, err)

	shortPayload := pipPubBytes[:32]
	shortPublic, err := bech32m.EncodeFromBase256WithType(crypto.PublicKeyHRP, crypto.SignatureTypeSecp256k1, shortPayload)
	require.NoError(t, err)

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
			"invalid checksum",
			"public1yqdkke2kzfzheda405lusfa2sy5aq70hn7k4zle5r322my9nfz35wyfamrft",
			false, nil,
		},
		{
			"invalid HRP",
			wrongHRP,
			false, nil,
		},
		{
			"invalid length: 32",
			shortPublic,
			false, nil,
		},
		{
			"invalid signature type: 3",
			"public1rafnl324uwngqdq455ax4e52fedmfcvskkwas6wsau0u0nwj4g96qztd56p",
			false, nil,
		},
		{
			"",
			"public1yqdkke2kzfzheda405lusfa2sy5aq70hn7k4zle5r322my9nfz35wyfamrfs",
			true,
			[]byte{
				0x03, 0x6d, 0x6c, 0xaa, 0xc2, 0x48, 0xaf, 0x96, 0xf6, 0xaf, 0xa7, 0xf9, 0x04, 0xf5, 0x50, 0x25,
				0x3a, 0x0f, 0x3e, 0xf3, 0xf5, 0xaa, 0x2f, 0xe6, 0x83, 0x8a, 0x95, 0xb2, 0x16, 0x69, 0x14, 0x68, 0xe2,
			},
		},
	}

	for no, tt := range tests {
		pub, err := secp256k1.PublicKeyFromString(tt.encoded)
		if tt.valid {
			require.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, tt.result, pub.Bytes(), "test %v: invalid bytes", no)
			assert.Equal(t, tt.encoded, pub.String(), "test %v: invalid encoded", no)
		} else {
			require.Error(t, err, "test %v", no)
			assert.Contains(t, err.Error(), tt.errMsg, "test %v: error not matched", no)
		}
	}
}

func TestPublicKeyFromStringInvalidPoint(t *testing.T) {
	bad := "public1yqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqglnhh9"
	_, err := secp256k1.PublicKeyFromString(bad)
	require.Error(t, err)
}

func TestPublicKeyFromBytesInvalidLength(t *testing.T) {
	_, err := secp256k1.PublicKeyFromBytes([]byte{0x02, 0x01})
	require.Error(t, err)
}

func TestPublicKeyFromBytesInvalidHexPattern(t *testing.T) {
	_, err := secp256k1.PublicKeyFromBytes(make([]byte, secp256k1.PublicKeySize))
	require.Error(t, err)
}
