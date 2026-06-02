package vault

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateMnemonic(t *testing.T) {
	_, err := GenerateMnemonic(127)
	require.Error(t, err, "low entropy")

	_, err = GenerateMnemonic(128)
	require.NoError(t, err)

	_, err = GenerateMnemonic(257)
	require.Error(t, err, "high entropy")

	_, err = GenerateMnemonic(256)
	require.NoError(t, err)
}

func TestValidateMnemonic(t *testing.T) {
	tests := []struct {
		mnenomic string
		errStr   string
	}{
		{
			"",
			"invalid mnenomic",
		},
		{
			"abandon ability able about above absent absorb abstract absurd abuse access",
			"invalid mnenomic",
		},
		{
			"bandon ability able about above absent absorb abstract absurd abuse access ability",
			"word `bandon` not found in reverse map",
		},
		{
			"abandon ability able about above absent absorb abstract absurd abuse access accident",
			"checksum incorrect",
		},
		{
			"abandon ability able about above absent absorb abstract absurd abuse access ability",
			"",
		},
	}
	for no, tt := range tests {
		err := CheckMnemonic(tt.mnenomic)
		if err != nil {
			assert.Equal(t, tt.errStr, err.Error(), "test %v failed", no)
		}
	}
}

func TestPrivateKeyFromString(t *testing.T) {
	tests := []struct {
		str string
		err error
	}{
		{str: "SECRET1PQQQSYQCYQ5RQWZQFPG9SCRGWPUGPZYSNZS23V9CCRYDPK8QARC0SEZYD4L", err: nil},
		{str: "SECRET1RQQQSYQCYQ5RQWZQFPG9SCRGWPUGPZYSNZS23V9CCRYDPK8QARC0SW5D8X2", err: nil},
		{str: "SECRET1YQQQSYQCYQ5RQWZQFPG9SCRGWPUGPZYSNZS23V9CCRYDPK8QARC0SPVXU8Z", err: nil},
		{str: "SECRET19QQQSYQCYQ5RQWZQFPG9SCRGWPUGPZYSNZS23V9CCRYDPK8QARC0S78KE6U", err: ErrInvalidPrivateKey},
	}

	for _, tt := range tests {
		_, err := PrivateKeyFromString(tt.str)
		require.ErrorIs(t, err, tt.err, "unexpected error for input: %s", tt.str)
	}
}

func TestPublicKeyFromString(t *testing.T) {
	tests := []struct {
		str string
		err error
	}{
		{str: "public1p5u5sljqq6tg57tw9uh95z6lt7vn8mlk3cmp608rwm38t68n90k2km2sx5t724l2ze99ktvedf4p7" +
			"5ymglpssq6pfc36m0428vwjs9h7hzl5a28zucl02u2vpu4sfp2ppe8zmet7p9xu9nysr4wvsx86vuujrva2z", err: nil},
		{str: "public1rqwss00lnecgtu8tsm5vwwj7qn9n7f43snwjs6hcamjrxgyj4xxuq5agu5g", err: nil},
		{str: "public1yqdkke2kzfzheda405lusfa2sy5aq70hn7k4zle5r322my9nfz35wyfamrfs", err: nil},
		{str: "public19qdkke2kzfzheda405lusfa2sy5aq70hn7k4zle5r322my9nfz35wygp96am", err: ErrInvalidPublicKey},
	}

	for _, tt := range tests {
		_, err := PublicKeyFromString(tt.str)
		require.ErrorIs(t, err, tt.err, "unexpected error for input: %s", tt.str)
	}
}

func TestSignatureFromString(t *testing.T) {
	tests := []struct {
		str string
		typ crypto.SignatureType
		err error
	}{
		{
			str: "8bdda74336efdf43b428a3811d3d6867a19e20889c91261b02a6b950b130f5bb22621394667c27660bfed2a8719d9c52",
			typ: crypto.SignatureTypeBLS, err: nil,
		},
		{
			str: "1fc2c800499342d08242db9c3eb654027cb7b821e6af9ede56dfdb67e824f15b" +
				"ddb419d2db3fd5aaf3ef1a9ebb9a9deb749380f0d6a110cbe95319fe9f794305",
			typ: crypto.SignatureTypeEd25519, err: nil,
		},
		{
			str: "c86779676d217b04979434e5bd37eddd02b671e9a54b48d3a812c7862dcb5396" +
				"31bb5e8459fec007608f50ea5661e0a5215aac976705404cb4f36ee623e63199",
			typ: crypto.SignatureTypeSecp256k1, err: nil,
		},
	}

	for _, tt := range tests {
		_, err := SignatureFromString(tt.str, tt.typ)
		require.ErrorIs(t, err, tt.err, "unexpected error for input: %s", tt.str)
	}
}
