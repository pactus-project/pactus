package hdkeychain

import (
	"encoding/hex"
	"io"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/secp256k1"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBIP32Vector1(t *testing.T) {
	testSeed, err := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	require.NoError(t, err)

	h := hardenedKeyStart
	tests := []struct {
		name     string
		path     []uint32
		wantPriv string
		wantPub  string
	}{
		{
			name:     "m",
			path:     []uint32{},
			wantPriv: "e8f32e723decf4051aefac8e2c93c9c5b214313817cdb01a1494b917c8436b35",
			wantPub:  "0339a36013301597daef41fbe593a02cc513d0b55527ec2df1050e2e8ff49c85c2",
		},
		{
			name:     "m/0H",
			path:     []uint32{h},
			wantPriv: "edb2e14f9ee77d26dd93b4ecede8d16ed408ce149b6cd80b0715a2d911a0afea",
			wantPub:  "035a784662a4a20a65bf6aab9ae98a6c068a81c52e4b032c0fb5400c706cfccc56",
		},
		{
			name:     "m/0H/1",
			path:     []uint32{h, 1},
			wantPriv: "3c6cb8d0f6a264c91ea8b5030fadaa8e538b020f0a387421a12de9319dc93368",
			wantPub:  "03501e454bf00751f24b1b489aa925215d66af2234e3891c3b21a52bedb3cd711c",
		},
		{
			name:     "m/0H/1/2H",
			path:     []uint32{h, 1, 2 + h},
			wantPriv: "cbce0d719ecf7431d88e6a89fa1483e02e35092af60c042b1df2ff59fa424dca",
			wantPub:  "0357bfe1e341d01c69fe5654309956cbea516822fba8a601743a012a7896ee8dc2",
		},
		{
			name:     "m/0H/1/2H/2",
			path:     []uint32{h, 1, 2 + h, 2},
			wantPriv: "0f479245fb19a38a1954c5c7c0ebab2f9bdfd96a17563ef28a6a4b1a2a764ef4",
			wantPub:  "02e8445082a72f29b75ca48748a914df60622a609cacfce8ed0e35804560741d29",
		},
		{
			name:     "m/0H/1/2H/2/1000000000",
			path:     []uint32{h, 1, 2 + h, 2, 1000000000},
			wantPriv: "471b76e389e528d6de6d816857e012c5455051cad6660850e58372a6c3e6e7c8",
			wantPub:  "022a471424da5e657499d1ff51cb43c47481a03b1e77f951fe64cec9f5a48f7011",
		},
	}

	masterKey, err := NewMaster(testSeed)
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extKey, err := masterKey.DerivePath(tt.path)
			require.NoError(t, err)

			require.Equal(t, tt.path, extKey.Path())

			privKey, err := extKey.RawPrivateKey()
			require.NoError(t, err)
			require.Equal(t, tt.wantPriv, hex.EncodeToString(privKey))

			pubKey := extKey.RawPublicKey()
			require.Equal(t, tt.wantPub, hex.EncodeToString(pubKey))

			secpPrivKey, err := secp256k1.PrivateKeyFromBytes(privKey)
			require.NoError(t, err)
			require.Equal(t, pubKey, secpPrivKey.PublicKeyNative().Bytes())

			neuterKey := extKey.Neuter()
			_, err = neuterKey.RawPrivateKey()
			require.ErrorIs(t, err, ErrNotPrivExtKey)
			require.Equal(t, pubKey, neuterKey.RawPublicKey())
			require.False(t, neuterKey.IsPrivate())
		})
	}
}

func TestNonHardenedDerivation(t *testing.T) {
	testSeed, err := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	require.NoError(t, err)

	masterKey, err := NewMaster(testSeed)
	require.NoError(t, err)

	parent, err := masterKey.Derive(hardenedKeyStart)
	require.NoError(t, err)

	extKey1, err := parent.Derive(1)
	require.NoError(t, err)

	extKey2, err := parent.Neuter().Derive(1)
	require.NoError(t, err)

	require.Equal(t, extKey1.RawPublicKey(), extKey2.RawPublicKey())
}

func TestInvalidDerivation(t *testing.T) {
	t.Run("private key is 31 bytes", func(t *testing.T) {
		key := [31]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], []uint32{}, true)
		_, err := ext.Derive(hardenedKeyStart)
		require.ErrorIs(t, err, ErrInvalidKeyData)
	})

	t.Run("public key is invalid length", func(t *testing.T) {
		key := [32]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], []uint32{}, false)
		_, err := ext.Derive(0)
		require.ErrorIs(t, err, ErrInvalidKeyData)
	})

	t.Run("derive hardened from public key", func(t *testing.T) {
		pub, err := hex.DecodeString("0339a36013301597daef41fbe593a02cc513d0b55527ec2df1050e2e8ff49c85c2")
		require.NoError(t, err)
		chainCode := [32]byte{0}
		ext := newExtendedKey(pub, chainCode[:], []uint32{}, false)
		_, err = ext.Derive(hardenedKeyStart)
		require.ErrorIs(t, err, ErrDeriveHardFromPublic)
	})
}

func TestNewMasterInvalidSeed(t *testing.T) {
	_, err := NewMaster([]byte{0, 1, 2})
	require.ErrorIs(t, err, ErrInvalidSeedLen)
}

func TestNeuter(t *testing.T) {
	extKey, err := NewKeyFromString(
		"XSECRET1YQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VSV0QLX0WFP62LT3A4N38XXGMTJVW6PG4CPMTWJZC25329S0RN9A8" +
			"SQYPR3KAHR38JJ34K7DKQKS4LQZTZ525Z3ETTXVZZSUKPH9FKRUMNUSGMAYA2",
	)
	require.NoError(t, err)

	neuterKey := extKey.Neuter()
	assert.Equal(
		t,
		"xpublic1yq5qqqqyqqyqqqqqzqqqgqqsqqqqqpj568vsv0qlx0wfp62lt3a4n38xxgmtjvw6pg4cpmtwjzc25329s0rn9a8"+
			"sqyypz53c5ynd9uet5n8gl75wtg0z8fqdq8v080723lejvaj045j8hqygefdyx2",
		neuterKey.String(),
	)
	assert.Equal(t, neuterKey, neuterKey.Neuter())

	_, err = neuterKey.RawPrivateKey()
	require.ErrorIs(t, err, ErrNotPrivExtKey)
}

// TestKeyToString ensures the String function works as intended.
//
//nolint:lll // long extended keys
func TestKeyToString(t *testing.T) {
	testSeed, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	h := hardenedKeyStart
	tests := []struct {
		name      string
		path      []uint32
		wantXPriv string
		wantXPub  string
	}{
		{
			name:      "derivation path: m",
			path:      []uint32{},
			wantXPriv: "XSECRET1YQQSGW00LS8QZ75JKY073LEGK06KR54DQF80R6V2TKSHWYFLLA5MA2ZQQYR50XTNJ8HK0GPG6A7KGUTYNE8ZMY9P38QTUMVQ6ZJ2TJ97GGD4N2DCRPEC",
			wantXPub:  "xpublic1yqqsgw00ls8qz75jky073legk06kr54dqf80r6v2tkshwyflla5ma2zqqyypnngmqzvcpt976aaqlhevn5qkv2y7sk42j0mpd7yzsut507jwgtssw0ur97",
		},
		{
			name:      "derivation path: m/0H",
			path:      []uint32{h},
			wantXPriv: "XSECRET1YQYQQQQYQYPRLMT9APUGFWPPM0RRRCGXRFM6WMXS3RKVQQ3ADZC5ZC7HXYDS5ZQPQAKEWZNU7UA7JDHVNKNKWM6X3DM2Q3NS5NDKDSZC8ZK3DJYDQ4L4QH8692L",
			wantXPub:  "xpublic1yqyqqqqyqyprlmt9apugfwppm0rrrcgxrfm6wmxs3rkvqq3adzc5zc7hxyds5zqppqdd8s3nz5j3q5edld24e46v2dsrg4qw99e9sxtq0k4qqcurvlnx9vjud6tv",
		},
		{
			name:      "derivation path: m/0H/1",
			path:      []uint32{h, 1},
			wantXPriv: "XSECRET1YQGQQQQYQQYQQQQPQ9FU9WCCNS6AZ8KK2CDQCPHGESDE5U3Z0M0MHGPQ40R5MDTDN0SVSQGPUDJUDPA4ZVNY3A294QV86M25W2W9SYRC28P6ZRGFDAYCEMJFNDQKP4GEW",
			wantXPub:  "xpublic1yqgqqqqyqqyqqqqpq9fu9wccns6az8kk2cdqcphgesde5u3z0m0mhgpq40r5mdtdn0svsqggr2q0y2jlsqaglyjcmfzd2jffpt4n27g35uwy3cwep5547mv7dwywqlzx79l",
		},
		{
			name:      "derivation path: m/0H/1/2H",
			path:      []uint32{h, 1, 2 + h},
			wantXPriv: "XSECRET1YQVQQQQYQQYQQQQQZQQQGQGQYGE4EEJ8PV85KVSYU555CD3VY7PLFMJQLWDWMDQ7RLAHV0V2S8UQZPJ7WP4CEANM5X8VGU65FLG2G8CPWX5YJ4ASVQS43MUHLT8AYYNW277UQXD",
			wantXPub:  "xpublic1yqvqqqqyqqyqqqqqzqqqgqgqyge4eej8pv85kvsyu555cd3vy7plfmjqlwdwmdq7rlahv0v2s8uqzzq6hhls7xswsr35lu4j5xzv4djl2295z97ag5cqhgwsp9fufdm5dcg3jm7ef",
		},
		{
			name:      "derivation path: m/0H/1/2H/2",
			path:      []uint32{h, 1, 2 + h, 2},
			wantXPriv: "XSECRET1YQSQQQQYQQYQQQQQZQQQGQQSQQQQZPNAHRZPLQ9NK7KRAQG7V2W34H3LC3AEYK8UV9ZF2CYN44JPZ50KAQQSQ73UJGHA3NGU2R92VT37QAW4JLX7LM94PW437729X5JC69FMYAAQF53EN6",
			wantXPub:  "xpublic1yqsqqqqyqqyqqqqqzqqqgqqsqqqqzpnahrzplq9nk7kraqg7v2w34h3lc3aeyk8uv9zf2cyn44jpz50kaqqss96zy2zp2wtefkaw2fp6g4y2d7crz9fsfet8uarksudvqg4s8g8ffwz5e44",
		},
		{
			name:      "derivation path: m/0H/1/2H/2/1000000000",
			path:      []uint32{h, 1, 2 + h, 2, 1000000000},
			wantXPriv: "XSECRET1YQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VSV0QLX0WFP62LT3A4N38XXGMTJVW6PG4CPMTWJZC25329S0RN9A8SQYPR3KAHR38JJ34K7DKQKS4LQZTZ525Z3ETTXVZZSUKPH9FKRUMNUSGMAYA2",
			wantXPub:  "xpublic1yq5qqqqyqqyqqqqqzqqqgqqsqqqqqpj568vsv0qlx0wfp62lt3a4n38xxgmtjvw6pg4cpmtwjzc25329s0rn9a8sqyypz53c5ynd9uet5n8gl75wtg0z8fqdq8v080723lejvaj045j8hqygefdyx2",
		},
	}

	masterKey, err := NewMaster(testSeed)
	require.NoError(t, err)

	for no, tt := range tests {
		extKey, err := masterKey.DerivePath(tt.path)
		require.NoError(t, err)
		neuterKey := extKey.Neuter()

		require.Equal(t, tt.wantXPriv, extKey.String(), "test %d failed", no)
		require.Equal(t, tt.wantXPub, neuterKey.String(), "test %d failed", no)

		recoveredExtKey, err := NewKeyFromString(tt.wantXPriv)
		require.NoError(t, err)

		recoveredNeuterKey, err := NewKeyFromString(tt.wantXPub)
		require.NoError(t, err)

		privKey, err := extKey.RawPrivateKey()
		require.NoError(t, err)
		recoveredPrivKey, err := recoveredExtKey.RawPrivateKey()
		require.NoError(t, err)

		require.Equal(t, privKey, recoveredPrivKey)
		require.Equal(t, extKey.ChainCode(), recoveredExtKey.ChainCode())
		require.Equal(t, extKey.RawPublicKey(), recoveredNeuterKey.RawPublicKey())
		require.Equal(t, neuterKey.ChainCode(), recoveredNeuterKey.ChainCode())
		require.Equal(t, tt.path, recoveredExtKey.Path())
		require.Equal(t, tt.path, recoveredNeuterKey.Path())
		require.True(t, recoveredExtKey.IsPrivate())
		require.False(t, recoveredNeuterKey.IsPrivate())
	}
}

// TestInvalidString checks errors corresponding to the invalid strings.
//
//nolint:lll // long extended keys
func TestInvalidString(t *testing.T) {
	tests := []struct {
		desc          string
		str           string
		expectedError error
	}{
		{
			desc:          "invalid checksum",
			str:           "XSECRET1YQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VSV0QLX0WFP62LT3A4N38XXGMTJVW6PG4CPMTWJZC25329S0RN9A8SQYPR3KAHR38JJ34K7DKQKS4LQZTZ525Z3ETTXVZZSUKPH9FKRUMNUSGMAYAX",
			expectedError: bech32m.InvalidChecksumError{Expected: "gmaya2", Actual: "gmayax"},
		},
		{
			desc:          "no depth",
			str:           "XSECRET1YG0AHEP",
			expectedError: io.EOF,
		},
		{
			desc:          "wrong path",
			str:           "XSECRET1YQYNXV7PF",
			expectedError: io.EOF,
		},
		{
			desc:          "no chain code",
			str:           "XSECRET1YQYQQQQYQJXNNTG",
			expectedError: io.EOF,
		},
		{
			desc:          "no reserved",
			str:           "XSECRET1YQYQQQQYQYQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQK8637R",
			expectedError: io.EOF,
		},
		{
			desc:          "no key",
			str:           "XSECRET1YQYQQQQYQYQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQTFRFAC",
			expectedError: io.EOF,
		},
		{
			desc:          "invalid type",
			str:           "XSECRET1PQGQQQQYQQYQQQQPQ9FU9WCCNS6AZ8KK2CDQCPHGESDE5U3Z0M0MHGPQ40R5MDTDN0SVSQGPUDJUDPA4ZVNY3A294QV86M25W2W9SYRC28P6ZRGFDAYCEMJFNDQFAQJAN",
			expectedError: ErrInvalidKeyData,
		},
		{
			desc:          "invalid type",
			str:           "xpublic1pqgqqqqyqqyqqqqpq9fu9wccns6az8kk2cdqcphgesde5u3z0m0mhgpq40r5mdtdn0svsqgpudjudpa4zvny3a294qv86m25w2w9syrc28p6zrgfdaycemjfndq7d8ma2",
			expectedError: ErrInvalidKeyData,
		},
		{
			desc:          "invalid hrp",
			str:           "SECRET1YQGQQQQYQQYQQQQPQ9FU9WCCNS6AZ8KK2CDQCPHGESDE5U3Z0M0MHGPQ40R5MDTDN0SVSQGPUDJUDPA4ZVNY3A294QV86M25W2W9SYRC28P6ZRGFDAYCEMJFNDQ6WKQJX",
			expectedError: crypto.InvalidHRPError("secret"),
		},
		{
			desc:          "invalid hrp",
			str:           "public1yqgqqqqyqqyqqqqpq9fu9wccns6az8kk2cdqcphgesde5u3z0m0mhgpq40r5mdtdn0svsqgpudjudpa4zvny3a294qv86m25w2w9syrc28p6zrgfdaycemjfndqd73fjl",
			expectedError: crypto.InvalidHRPError("public"),
		},
	}

	for no, tt := range tests {
		_, err := NewKeyFromString(tt.str)
		require.ErrorIs(t, err, tt.expectedError, "test %d error is not matched", no)
	}
}
