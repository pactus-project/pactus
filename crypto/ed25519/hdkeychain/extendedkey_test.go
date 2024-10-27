package hdkeychain

import (
	"encoding/hex"
	"io"
	"testing"

	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNonHardenedDerivation tests deriving a new key in non-hardened mode.
// It should return an error.
func TestNonHardenedDerivation(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	testSeed := ts.RandBytes(32)
	path := []uint32{
		ts.RandUint32(HardenedKeyStart),
	}

	masterKey, _ := NewMaster(testSeed)
	_, err := masterKey.DerivePath(path)
	assert.ErrorIs(t, err, ErrNonHardenedPath)
}

// TestHardenedDerivation tests derive key in hardened mode.
func TestHardenedDerivation(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	testSeed := ts.RandBytes(32)
	path := []uint32{
		ts.RandUint32(HardenedKeyStart) + HardenedKeyStart,
	}

	masterKey, err := NewMaster(testSeed)
	require.NoError(t, err)

	extKey, err := masterKey.DerivePath(path)
	require.NoError(t, err)

	assert.Equal(t, path, extKey.Path())
}

// TestDerivation verifies the derivation of new keys in hardened mode.
// The test cases are based on the SLIP-0010 standard.
func TestDerivation(t *testing.T) {
	testSeed, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	h := HardenedKeyStart
	tests := []struct {
		name    string
		path    []uint32
		wantPrv string
		wantPub string
	}{
		{
			name:    "derivation path: m",
			path:    []uint32{},
			wantPrv: "2b4be7f19ee27bbf30c667b642d5f4aa69fd169872f8fc3059c08ebae2eb19e7",
			wantPub: "a4b2856bfec510abab89753fac1ac0e1112364e7d250545963f135f2a33188ed",
		},
		{
			name:    "derivation path: m/0H",
			path:    []uint32{h},
			wantPrv: "68e0fe46dfb67e368c75379acec591dad19df3cde26e63b93a8e704f1dade7a3",
			wantPub: "8c8a13df77a28f3445213a0f432fde644acaa215fc72dcdf300d5efaa85d350c",
		},
		{
			name:    "derivation path: m/0H/1H",
			path:    []uint32{h, 1 + h},
			wantPrv: "b1d0bad404bf35da785a64ca1ac54b2617211d2777696fbffaf208f746ae84f2",
			wantPub: "1932a5270f335bed617d5b935c80aedb1a35bd9fc1e31acafd5372c30f5c1187",
		},
		{
			name:    "derivation path: m/0H/1H/2H",
			path:    []uint32{h, 1 + h, 2 + h},
			wantPrv: "92a5b23c0b8a99e37d07df3fb9966917f5d06e02ddbd909c7e184371463e9fc9",
			wantPub: "ae98736566d30ed0e9d2f4486a64bc95740d89c7db33f52121f8ea8f76ff0fc1",
		},
		{
			name:    "derivation path: m/0H/1H/2H/2H",
			path:    []uint32{h, 1 + h, 2 + h, 2 + h},
			wantPrv: "30d1dc7e5fc04c31219ab25a27ae00b50f6fd66622f6e9c913253d6511d1e662",
			wantPub: "8abae2d66361c879b900d204ad2cc4984fa2aa344dd7ddc46007329ac76c429c",
		},
		{
			name:    "derivation path: m/0H/1H/2H/2H/1000000000H",
			path:    []uint32{h, 1 + h, 2 + h, 2 + h, 1000000000 + h},
			wantPrv: "8f94d394a8e8fd6b1bc2f3f49f5c47e385281d5c17e65324b0f62483e37e8793",
			wantPub: "3c24da049451555d51a7014a37337aa4e12d41e485abccfa46b47dfb2af54b7a",
		},
	}

	masterKey, _ := NewMaster(testSeed)
	for no, tt := range tests {
		extKey, err := masterKey.DerivePath(tt.path)
		require.NoError(t, err)

		privKey := extKey.RawPrivateKey()
		require.Equal(t, tt.wantPrv, hex.EncodeToString(privKey),
			"mismatched serialized private key for test #%v", no+1)

		pubKey := extKey.RawPublicKey()
		require.Equal(t, tt.wantPub, hex.EncodeToString(pubKey),
			"mismatched serialized public key for test #%v", no+1)

		require.Equal(t, tt.path, extKey.Path())
	}
}

// TestGenerateSeed ensures the GenerateSeed function works as intended.
func TestGenerateSeed(t *testing.T) {
	tests := []struct {
		name   string
		length uint8
		err    error
	}{
		// Test various valid lengths.
		{name: "16 bytes", length: 16},
		{name: "17 bytes", length: 17},
		{name: "20 bytes", length: 20},
		{name: "32 bytes", length: 32},
		{name: "64 bytes", length: 64},

		// Test invalid lengths.
		{name: "15 bytes", length: 15, err: ErrInvalidSeedLen},
		{name: "65 bytes", length: 65, err: ErrInvalidSeedLen},
	}

	for no, tt := range tests {
		seed, err := GenerateSeed(tt.length)
		assert.ErrorIs(t, err, tt.err)

		if tt.err == nil {
			assert.Len(t, seed, int(tt.length),
				"GenerateSeed #%d (%s): length mismatch -- got %d, want %d",
				no+1, tt.name, len(seed), tt.length)
		}
	}
}

// TestNewMaster ensures the NewMaster function works as intended.
func TestNewMaster(t *testing.T) {
	tests := []struct {
		name string
		seed string
		key  string
		err  error
	}{
		// Test various valid seeds.
		{
			name: "16 bytes",
			seed: "000102030405060708090a0b0c0d0e0f",
			key:  "2b4be7f19ee27bbf30c667b642d5f4aa69fd169872f8fc3059c08ebae2eb19e7",
		},
		{
			name: "32 bytes",
			seed: "3ddd5602285899a946114506157c7997e5444528f3003f6134712147db19b678",
			key:  "4b36bc63a15797f4d506074f36f2f3904bc0f10179b5ab91183c167e9c2dcf0e",
		},
		{
			name: "64 bytes",
			seed: "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784" +
				"817e7b7875726f6c696663605d5a5754514e4b484542",
			key: "171cb88b1b3c1db25add599712e36245d75bc65a1a5c9e18d76f9f2b1eab4012",
		},

		// Test invalid seeds.
		{
			name: "empty seed",
			seed: "",
			err:  ErrInvalidSeedLen,
		},
		{
			name: "15 bytes",
			seed: "000000000000000000000000000000",
			err:  ErrInvalidSeedLen,
		},
		{
			name: "65 bytes",
			seed: "000000000000000000000000000000000000000000000000000000000000000000000000000" +
				"0000000000000000000000000000000000000000000000000000000",
			err: ErrInvalidSeedLen,
		},
	}

	for no, tt := range tests {
		seed, _ := hex.DecodeString(tt.seed)
		extKey, err := NewMaster(seed)
		assert.ErrorIs(t, err, tt.err)

		if tt.err == nil {
			privKey := extKey.RawPrivateKey()
			assert.Equal(t, tt.key, hex.EncodeToString(privKey),
				"NewMaster #%d (%s): key mismatch -- got %x, want %s",
				no+1, tt.name, privKey, tt.key)
		}
	}
}

// TestKeyToString ensures the String function works as intended.
//
//nolint:lll // long extended keys
func TestKeyToString(t *testing.T) {
	testSeed, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	h := HardenedKeyStart
	tests := []struct {
		name      string
		path      []uint32
		wantXPriv string
	}{
		{
			name:      "derivation path: m",
			path:      []uint32{},
			wantXPriv: "XSECRET1RQQSFQPR2J0098Q989D0Y2QG8FPT86H4Q9WLK2GHE08S9CRVD3J5LL7CQYQ45HEL3NM38H0ESCENMVSK47J4XNLGKNPE03LPST8QGAWHZAVV7WQ3TTQ3",
		},
		{
			name:      "derivation path: m/0H",
			path:      []uint32{h},
			wantXPriv: "XSECRET1RQYQQQQYQYZ94N2S38Q9KYN5P2PAZ0LKA5K075MGTW7D80ZGC5T7NTY8PD6WXJQPQDRS0U3KLKELRDRR4X7DVA3V3MTGEMU7DUFHX8WF63ECY78DDU73SSP8VGW",
		},
		{
			name:      "derivation path: m/0H/1H",
			path:      []uint32{h, 1 + h},
			wantXPriv: "XSECRET1RQGQQQQYQQYQQPQPQ5VSYYHMH6X6UY5Z6DVDJWWPTXUMGAEJQUD2HCV25Z6QPYS649U2QQG936ZADGP9LXHD8SKNYEGDV2JEXZUS36FMHD9HML7HJPRM5DT5Y7GG0AYKR",
		},
		{
			name:      "derivation path: m/0H/1H/2H",
			path:      []uint32{h, 1 + h, 2 + h},
			wantXPriv: "XSECRET1RQVQQQQYQQYQQPQQZQQQGQGPWDXFFUQ944VJS7JWRLVWP9UJJME876TQAHZPCWZ22P7XYE8XDDSQZPY49KG7QHZ5EUD7S0HELHXTXJ9L46PHQ9HDAJZW8UXZRW9RRA87FNLUP3Y",
		},
		{
			name:      "derivation path: m/0H/1H/2H/2H",
			path:      []uint32{h, 1 + h, 2 + h, 2 + h},
			wantXPriv: "XSECRET1RQSQQQQYQQYQQPQQZQQQGQQSQQZQZPRMDSLUN6AGWPM7VMGQH6E32RVC6YEHY5M6EJWC47HQLQLM5M4WVQQSRP5WU0E0UQNP3YXDTYK384CQT2RM06ENZ9AHFEYFJ20T9Z8G7VCSW77AAC",
		},
		{
			name:      "derivation path: m/0H/1H/2H/2H/1000000000H",
			path:      []uint32{h, 1 + h, 2 + h, 2 + h, 1000000000 + h},
			wantXPriv: "XSECRET1RQ5QQQQYQQYQQPQQZQQQGQQSQQZQQPJ56HVSXS7YEYWSV4SKDTG53W2J8TL57P7C5E44DKKKE3GL6WQENU7H6YVQQYZ8EF5U54R5066CMCTELF86UGL3C22QATST7V5EYKRMZFQLR06REXGPVGYE",
		},
	}

	masterKey, _ := NewMaster(testSeed)
	for no, tt := range tests {
		extKey, _ := masterKey.DerivePath(tt.path)
		require.Equal(t, tt.wantXPriv, extKey.String(), "test %d failed", no)

		recoveredExtKey, err := NewKeyFromString(tt.wantXPriv)
		require.NoError(t, err)

		require.Equal(t, extKey, recoveredExtKey)
		require.Equal(t, tt.path, recoveredExtKey.path)
	}
}

// TestInvalidString checks errors corresponding to the invalid strings
//
//nolint:lll // long extended private keys
func TestInvalidString(t *testing.T) {
	tests := []struct {
		desc          string
		str           string
		expectedError error
	}{
		{
			desc:          "invalid checksum",
			str:           "XSECRET1RQGQQQQYQQYQQPQPQ5VSYYHMH6X6UY5Z6DVDJWWPTXUMGAEJQUD2HCV25Z6QPYS649U2QQG936ZADGP9LXHD8SKNYEGDV2JEXZUS36FMHD9HML7HJPRM5DT5Y7GG0AYRK",
			expectedError: bech32m.InvalidChecksumError{Expected: "g0aykr", Actual: "g0ayrk"},
		},
		{
			desc:          "no depth",
			str:           "XSECRET1RFK28CY",
			expectedError: io.EOF,
		},
		{
			desc:          "wrong path",
			str:           "XSECRET1RQGQQQQYQ6EJ6DE",
			expectedError: io.EOF,
		},
		{
			desc:          "no chain code",
			str:           "XSECRET1RQGQQQQYQQYQQPQQ98TS98",
			expectedError: io.EOF,
		},
		{
			desc:          "no reserved",
			str:           "XSECRET1RQGQQQQYQQYQQPQPQ5VSYYHMH6X6UY5Z6DVDJWWPTXUMGAEJQUD2HCV25Z6QPYS649U2Q8GUZZJ",
			expectedError: io.EOF,
		},
		{
			desc:          "no key",
			str:           "XSECRET1RQGQQQQYQQYQQPQPQ5VSYYHMH6X6UY5Z6DVDJWWPTXUMGAEJQUD2HCV25Z6QPYS649U2QQ85HSJA",
			expectedError: io.EOF,
		},
		{
			desc:          "invalid type",
			str:           "XSECRET1YQGQQQQYQQYQQPQPQ5VSYYHMH6X6UY5Z6DVDJWWPTXUMGAEJQUD2HCV25Z6QPYS649U2QQG936ZADGP9LXHD8SKNYEGDV2JEXZUS36FMHD9HML7HJPRM5DT5Y7GTKSQQT",
			expectedError: ErrInvalidKeyData,
		},
		{
			desc:          "invalid hrp",
			str:           "SECRET1RQGQQQQYQQYQQPQPQ5VSYYHMH6X6UY5Z6DVDJWWPTXUMGAEJQUD2HCV25Z6QPYS649U2QQG936ZADGP9LXHD8SKNYEGDV2JEXZUS36FMHD9HML7HJPRM5DT5Y7GYQ7VAT",
			expectedError: ErrInvalidHRP,
		},
	}

	for no, tt := range tests {
		_, err := NewKeyFromString(tt.str)
		assert.ErrorIs(t, err, tt.expectedError, "test %d error is not matched", no)
	}
}
