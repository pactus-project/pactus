package secp256k1_test

import (
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/secp256k1"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrivateKeyEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv1 := ts.RandSecp256k1KeyPair()
	_, prv2 := ts.RandSecp256k1KeyPair()
	_, prv3 := ts.RandEd25519KeyPair()

	assert.True(t, prv1.EqualsTo(prv1))
	assert.False(t, prv1.EqualsTo(prv2))
	assert.False(t, prv1.EqualsTo(prv3))
}

func TestPrivateKeyFromString(t *testing.T) {
	pipPrvBytes := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}
	wrongHRP, err := bech32m.EncodeFromBase256WithType("XXX", crypto.SignatureTypeSecp256k1, pipPrvBytes)
	require.NoError(t, err)

	shortPayload := pipPrvBytes[:31]
	shortSecret, err := bech32m.EncodeFromBase256WithType(crypto.PrivateKeyHRP, crypto.SignatureTypeSecp256k1, shortPayload)
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
			"SECRET1YQQQSYQCYQ5RQWZQFPG9SCRGWPUGPZYSNZS23V9CCRYDPK8QARC0SPVXU8Y",
			false, nil,
		},
		{
			"invalid HRP",
			strings.ToUpper(wrongHRP),
			false, nil,
		},
		{
			"invalid signature type: 3",
			"SECRET1RJ6STNTA7Y3P2QLQF8A6QCX05F2H5TFNE5RSH066KZME4WVFXKE7QW097LG",
			false, nil,
		},
		{
			"invalid length: 31",
			strings.ToUpper(shortSecret),
			false, nil,
		},
		{
			"",
			"SECRET1YQQQSYQCYQ5RQWZQFPG9SCRGWPUGPZYSNZS23V9CCRYDPK8QARC0SPVXU8Z",
			true,
			[]byte{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			},
		},
	}

	for no, tt := range tests {
		prv, err := secp256k1.PrivateKeyFromString(tt.encoded)
		if tt.valid {
			require.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, tt.result, prv.Bytes(), "test %v: invalid bytes", no)
			assert.Equal(t, strings.ToUpper(tt.encoded), prv.String(), "test %v: invalid encoded", no)
		} else {
			require.Error(t, err, "test %v", no)
			assert.Contains(t, err.Error(), tt.errMsg, "test %v: error not matched", no)
		}
	}
}

func TestPrivateKeyFromBytes(t *testing.T) {
	_, err := secp256k1.PrivateKeyFromBytes(make([]byte, 31))
	require.Error(t, err)

	_, err = secp256k1.PrivateKeyFromBytes(make([]byte, 32))
	require.ErrorIs(t, err, crypto.ErrInvalidPrivateKey)
}
