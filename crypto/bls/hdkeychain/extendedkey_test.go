package hdkeychain

import (
	"encoding/hex"
	"io"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNonHardenedDerivation tests derive private key and public key in
// non hardened mode.
func TestNonHardenedDerivation(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	testSeed := ts.RandBytes(32)
	path := []uint32{
		ts.RandUint32(hardenedKeyStart),
		ts.RandUint32(hardenedKeyStart),
		ts.RandUint32(hardenedKeyStart),
		ts.RandUint32(hardenedKeyStart),
	}

	checkPublicKeyDerivation := func(masterKey *ExtendedKey, path []uint32) {
		neuterKey := masterKey.Neuter()

		extKey1, _ := masterKey.DerivePath(path)
		extKey2, _ := neuterKey.DerivePath(path)
		pubKey1 := extKey1.RawPublicKey()
		pubKey2 := extKey2.RawPublicKey()

		require.Equal(t, path, extKey1.Path())
		require.Equal(t, pubKey1, pubKey2)
	}

	masterKeyG1, _ := NewMaster(testSeed, true)
	masterKeyG2, _ := NewMaster(testSeed, false)

	checkPublicKeyDerivation(masterKeyG1, path)
	checkPublicKeyDerivation(masterKeyG2, path)
}

// TestHardenedDerivation tests derive private key and public key in
// hardened mode.
func TestHardenedDerivation(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	testSeed := ts.RandBytes(32)
	path := []uint32{
		ts.RandUint32(hardenedKeyStart) + hardenedKeyStart,
		ts.RandUint32(hardenedKeyStart) + hardenedKeyStart,
		ts.RandUint32(hardenedKeyStart) + hardenedKeyStart,
		ts.RandUint32(hardenedKeyStart) + hardenedKeyStart,
	}

	masterKey, _ := NewMaster(testSeed, false)
	extKey, _ := masterKey.DerivePath(path)
	privKey, _ := extKey.RawPrivateKey()
	blsPrivKey, _ := bls.PrivateKeyFromBytes(privKey)
	pubKey := extKey.RawPublicKey()

	assert.Equal(t, path, extKey.Path())
	assert.Equal(t, pubKey, blsPrivKey.PublicKey().Bytes())
}

// TestDerivation tests derive private keys in hardened and non hardened modes.
func TestDerivation(t *testing.T) {
	testSeed, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	h := hardenedKeyStart
	tests := []struct {
		name       string
		path       []uint32
		wantPrivG1 string
		wantPrivG2 string
		wantPubG1  string
		wantPubG2  string
	}{
		{
			name:       "derivation path: m",
			path:       []uint32{},
			wantPrivG1: "4f55e31ee1c4f58af0840fd3f5e635fd6c07eacd14283c45d7d43729003abb84",
			wantPrivG2: "4f55e31ee1c4f58af0840fd3f5e635fd6c07eacd14283c45d7d43729003abb84",
			wantPubG1:  "8fbed8842588b629377c0a0d0d9547a9ee17527d5fd6d2c609034a8c3c074dda031e0dfe886b454499bfe0f40a7c4b18",
			wantPubG2: "b1bad3bf4a4ae87c89dec2c32512603ca08e2db62cfd2254c96bfe75068f5a98e7c4cd7d37cf0496dd6e79703e7c88e5" +
				"046bdec9c896ef2ad030096bbcf73c6cff17add3da9530f22491901fdf7fd2076c0f08ea35a4fdaa00e7ac6d0a5442e3",
		},
		{
			name:       "derivation path: m/0H",
			path:       []uint32{h},
			wantPrivG1: "5f5d7bfae7eabf2cc3faebc12449e1c7116c2777d7e384ead79df299667b8d9a",
			wantPrivG2: "5695ba5087a27f8c0d7270455104658b2367b8e90ab6f7f57ac7ce22d4a6836c",
			wantPubG1:  "b2826a89a22fec3349d64f4379a1eb5632b0b345b985b738324a5b8db640307421201efe36ae6c8c639d32d4124496ae",
			wantPubG2: "b37da3080662ceeb7f07289801a56e5c555d413434ad096079c084caa162c8d224891f68816921f5bd1453af7d085bc4" +
				"00341d61ce496ffb11cd10f8e90522447fada1a5f646c45797e00460925876f0b63f4023bf27e828688f7b4dd833e641",
		},
		{
			name:       "derivation path: m/0H/1",
			path:       []uint32{h, 1},
			wantPrivG1: "3bea739c9a2695ba4af566bc3f28e5c62da8e721b977709f9d492f7129b83521",
			wantPrivG2: "555422bcbffd1d55eea6f87a924ba5d046bb60e2bffe2182daf78bab6a6e179f",
			wantPubG1:  "af5980f4172797c07174a4040eb0b1859b357b05f0a29ac65c35d957730fd722ffd520d861e8fbe3126d26ceb08dbe52",
			wantPubG2: "b5f783bb1f1173feebb083f146c5a83470e84f26177862c5ab5b8be34ae6e3955d1b324f501a0d2751d971805f0612bc" +
				"0b5e966c9060eeb08cf38a7e71037863ffb2f6433694e69db59f731dbe55125f995d2d6ccd139d56d5b481d3bce76baa",
		},
		{
			name:       "derivation path: m/0H/1/2H",
			path:       []uint32{h, 1, 2 + h},
			wantPrivG1: "221e1f998e9599aecdab1c9671162bea925ee50d5f1c5bca2ed19908ac0f2ddd",
			wantPrivG2: "39e4906c49c05f5daeed89ced104a32cda82782654dcc116346144424746f871",
			wantPubG1:  "b06503dda77e1408478fc4b2d044a0ce2ab73691e8497a37f99d00e1076782698aacceb8e68fb9c3db6deccb0b8375fe",
			wantPubG2: "81461b89b446d055ac3bc38b9384363cbabc47cc0a16c97a7c7ea24eeffd70f213daacdfd736a49c45befececcd81832" +
				"12f04e186bcc9fbf67bfa5de862c57298cff4d36d5409380a166b9e37348b665186019b15498608309936e7ff36a87b5",
		},
		{
			name:       "derivation path: m/0H/1/2H/2",
			path:       []uint32{h, 1, 2 + h, 2},
			wantPrivG1: "26a19ca5ff2f6b32871de71aabd87a30ce79cdde3b0556cbb46692295f0aee15",
			wantPrivG2: "3aa1e19a9bf2bf631d95b401e29d5f042160edd76ced9696e42a98be80b41faa",
			wantPubG1:  "afd589792ba6bcb1866598a673a96fdaef9bf94026ef875a1a3e8d4fd839360f4659c9495afaf24c52577c0aa1fb5d45",
			wantPubG2: "92b20565b4a02bf82229f32e0ccc6f23446ded5ca2d67067afc70931b5a934f9469651e67e1105b5601cb585a1f44538" +
				"124fe3529f5b1edb27ab44f0900e59a27f57df87aa03395a70825d02433c2498d8396c90986dad79d5ba9e0fc438bea8",
		},
		{
			name:       "derivation path: m/0H/1/2H/2/1000000000",
			path:       []uint32{h, 1, 2 + h, 2, 1000000000},
			wantPrivG1: "44b743b059c2e4cb720378f4f0eda9369a1f02294e140e6a2e444bfdd36b1ad9",
			wantPrivG2: "2b01ef29730eb62c7114621d9d28ad77cf33f2434572a2bf9b73f1e502fea770",
			wantPubG1:  "99b404130a1ae6b6dd90ddf2a25c692f405536fee11046257ed6ba11629f101ad80658c61c039f0523de4c6e9f58a5c8",
			wantPubG2: "b05a01a80c3fe465227c23df7e36be1adcf557111f4cc50bf0f00c66c2b084d1e1d96e2f1c754496cb1f83dd1123456e" +
				"17697e77a9b99ea557a63c9bf29668a966732882e7baebf079a4afad212910deb10e5151e18ae98ee4a57d0e622332aa",
		},
	}

	masterKeyG1, _ := NewMaster(testSeed, true)
	masterKeyG2, _ := NewMaster(testSeed, false)
	for no, tt := range tests {
		extKeyG1, err := masterKeyG1.DerivePath(tt.path)
		require.NoError(t, err)

		extKeyG2, err := masterKeyG2.DerivePath(tt.path)
		require.NoError(t, err)

		privKeyG1, err := extKeyG1.RawPrivateKey()
		require.NoError(t, err)
		require.Equal(t, tt.wantPrivG1, hex.EncodeToString(privKeyG1),
			"mismatched serialized private key for test #%v", no+1)

		privKeyG2, err := extKeyG2.RawPrivateKey()
		require.NoError(t, err)
		require.Equal(t, tt.wantPrivG2, hex.EncodeToString(privKeyG2),
			"mismatched serialized private key for test #%v", no+1)

		pubKeyG1 := extKeyG1.RawPublicKey()
		require.Equal(t, tt.wantPubG1, hex.EncodeToString(pubKeyG1),
			"mismatched serialized public key for test #%v", no+1)

		pubKeyG2 := extKeyG2.RawPublicKey()
		require.Equal(t, tt.wantPubG2, hex.EncodeToString(pubKeyG2),
			"mismatched serialized public key for test #%v", no+1)

		neuterKeyG1 := extKeyG1.Neuter()
		neuterKeyG2 := extKeyG2.Neuter()
		neuterPubKeyG1 := neuterKeyG1.RawPublicKey()
		neuterPubKeyG2 := neuterKeyG2.RawPublicKey()

		require.True(t, extKeyG1.IsPrivate())
		require.True(t, extKeyG2.IsPrivate())
		require.False(t, neuterKeyG1.IsPrivate())
		require.False(t, neuterKeyG2.IsPrivate())
		require.Equal(t, pubKeyG1, neuterPubKeyG1)
		require.Equal(t, pubKeyG2, neuterPubKeyG2)
		require.Equal(t, tt.path, extKeyG1.Path())
		require.Equal(t, tt.path, extKeyG2.Path())
		require.Equal(t, tt.path, neuterKeyG1.Path())
		require.Equal(t, tt.path, neuterKeyG2.Path())

		_, err = neuterKeyG1.RawPrivateKey()
		assert.ErrorIs(t, err, ErrNotPrivExtKey)

		_, err = neuterKeyG2.RawPrivateKey()
		assert.ErrorIs(t, err, ErrNotPrivExtKey)

		blsPrivKey, _ := bls.PrivateKeyFromBytes(privKeyG2)
		require.Equal(t, pubKeyG2, blsPrivKey.PublicKey().Bytes())
	}
}

// TestInvalidDerivation tests Derive function for invalid data.
func TestInvalidDerivation(t *testing.T) {
	t.Run("Private key is 31 bytes. It should be 32 bytes", func(t *testing.T) {
		key := [31]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], []uint32{}, true, false)
		_, err := ext.Derive(hardenedKeyStart)
		assert.ErrorIs(t, err, ErrInvalidKeyData)
	})

	t.Run("Public key on G1 is 96 bytes. It should be 48 bytes", func(t *testing.T) {
		key := [96]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], []uint32{}, false, true)
		_, err := ext.Derive(0)
		assert.ErrorIs(t, err, ErrInvalidKeyData)
	})

	t.Run("Public key on G2 is 42 bytes. It should be 96 bytes", func(t *testing.T) {
		key := [95]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], []uint32{}, false, false)
		_, err := ext.Derive(0)
		assert.ErrorIs(t, err, ErrInvalidKeyData)
	})

	t.Run("Invalid key", func(t *testing.T) {
		key := [95]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], []uint32{}, false, false)
		_, err := ext.Derive(0)
		assert.ErrorIs(t, err, ErrInvalidKeyData)
	})

	t.Run("Derive public key from hardened key", func(t *testing.T) {
		key := [32]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], []uint32{}, false, false)
		_, err := ext.Derive(hardenedKeyStart)
		assert.ErrorIs(t, err, ErrDeriveHardFromPublic)
	})
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
				no, tt.name, len(seed), tt.length)
		}
	}
}

// TestNewMaster ensures the NewMaster function works as intended.
func TestNewMaster(t *testing.T) {
	tests := []struct {
		name    string
		seed    string
		privKey string
		err     error
	}{
		// Test various valid seeds.
		{
			name:    "16 bytes",
			seed:    "000102030405060708090a0b0c0d0e0f",
			privKey: "4f55e31ee1c4f58af0840fd3f5e635fd6c07eacd14283c45d7d43729003abb84",
		},
		{
			name:    "32 bytes",
			seed:    "3ddd5602285899a946114506157c7997e5444528f3003f6134712147db19b678",
			privKey: "4c101174339ffca5cc0afca5d2d8e2538834781318e5e1c8afdabf7e6fb77444",
		},
		{
			name: "64 bytes",
			seed: "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7" +
				"b7875726f6c696663605d5a5754514e4b484542",
			privKey: "47b660cc8dc2d4dc2cdf8893048bda9d5dc6318eb31f301b272b291b26cb20a1",
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
			seed: "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
				"000000000000000000000000000000000000000000000",
			err: ErrInvalidSeedLen,
		},
	}

	for no, tt := range tests {
		seed, _ := hex.DecodeString(tt.seed)
		extKeyG1, err := NewMaster(seed, true)
		assert.ErrorIs(t, err, tt.err)

		extKeyG2, err := NewMaster(seed, true)
		assert.ErrorIs(t, err, tt.err)

		if tt.err == nil {
			privKeyG1, _ := extKeyG1.RawPrivateKey()
			assert.Equal(t, tt.privKey, hex.EncodeToString(privKeyG1),
				"NewMaster #%d (%s): privKeyG1 mismatch -- got %x, want %s",
				no+1, tt.name, privKeyG1, tt.privKey)

			privKeyG2, _ := extKeyG2.RawPrivateKey()
			assert.Equal(t, tt.privKey, hex.EncodeToString(privKeyG2),
				"NewMaster #%d (%s): privKeyG2 mismatch -- got %x, want %s",
				no+1, tt.name, privKeyG2, tt.privKey)
		}
	}
}

// TestKeyToString ensures the String function works as intended.
//
//nolint:lll // long extended keys
func TestKeyToString(t *testing.T) {
	testSeed, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	h := hardenedKeyStart
	tests := []struct {
		name        string
		path        []uint32
		wantXPrivG1 string
		wantXPrivG2 string
		wantXPubG1  string
		wantXPubG2  string
	}{
		{
			name:        "derivation path: m",
			path:        []uint32{},
			wantXPrivG1: "XSECRET1PQQSTS7DSJ7AZNY54YZ53MM3FMCWEGWVVJYRK5SJ9HESHQSN96GCVJUSPYP84TCC7U8Z0TZHSSS8A8A0XXH7KCPL2E52ZS0Z96L2RW2GQ82ACG058EDE",
			wantXPrivG2: "XSECRET1PQQSTS7DSJ7AZNY54YZ53MM3FMCWEGWVVJYRK5SJ9HESHQSN96GCVJUSQYP84TCC7U8Z0TZHSSS8A8A0XXH7KCPL2E52ZS0Z96L2RW2GQ82ACGWGEQEJ",
			wantXPubG1:  "xpublic1pqqsts7dsj7azny54yz53mm3fmcwegwvvjyrk5sj9heshqsn96gcvjuspxz8makyyykytv2fh0s9q6rv4g757u96j040ad5kxpyp54rpuqaxa5qc7phlgs669gjvmlc85pf7ykxqqale8c",
			wantXPubG2:  "xpublic1pqqsts7dsj7azny54yz53mm3fmcwegwvvjyrk5sj9heshqsn96gcvjusqvzcm45alff9wslyfmmpvxfgjvq72pr3dkck06gj5e94luagx3adf3e7ye47n0ncyjmwku7ts8e7g3egyd00vnjykau4dqvqfdw70w0rvlut6m576j5c0yfy3jq0a7l7jqakq7z82xkj0m2squ7kx6zj5gt3snpvg7k",
		},
		{
			name:        "derivation path: m/0H",
			path:        []uint32{h},
			wantXPrivG1: "XSECRET1PQYQQQQYQYQDNX9T02WPS2RZ5SYUKE3JPHE8RGDHJMTNU7684679WEQWRN8STWQFQTAWHH7H8A2LJESL6A0QJGJ0PCUGKCFMH6L3CF6KHNHEFJENM3KDQRD8L6L",
			wantXPrivG2: "XSECRET1PQYQQQQYQYR38R7SGQNLMC6H96CANRN8XE3WVFVLF0V5XW2LE0FDSPYT52FUNSQPQ262M55Y85FLCCRTJWPZ4ZPR93V3K0W8FP2M00AT6CL8Z949XSDKQ5PJCQ7",
			wantXPubG1:  "xpublic1pqyqqqqyqyqdnx9t02wps2rz5syuke3jphe8rgdhjmtnu7684679weqwrn8stwqfsk2px4zdz9lkrxjwkfaphng0t2cetpv69hxzmwwpjffdcmdjqxp6zzgq7lcm2umyvvwwn94qjgjt2u0x45lj",
			wantXPubG2:  "xpublic1pqyqqqqyqyr38r7sgqnlmc6h96canrn8xe3wvfvlf0v5xw2le0fdspyt52funsqrqkd76xzqxvt8wklc89zvqrftwt3246sf5xjksjcreczzv4gtzerfzfzgldzqkjg04h5298tmappdugqp5r4suujt0lvgu6y8cayzjy3rl4ks6tajxc3te0cqyvzf9sahskcl5qgalyl5zs6y00dxasvlxgyepz53e",
		},
		{
			name:        "derivation path: m/0H/1",
			path:        []uint32{h, 1},
			wantXPrivG1: "XSECRET1PQGQQQQYQQYQQQQPQ6AXJT5395S9R89ME3E25LJXAP2QVULMX7S3UFNQ2D49Z0RHR38YQZGPMAFEEEX3XJKAY4ATXHSLJ3EWX9K5WWGDEWACFL82F9ACJNWP4YYSJCVL6",
			wantXPrivG2: "XSECRET1PQGQQQQYQQYQQQQPQFC96AZPJ5LSJKC3SEG5KUFF9Q7A9TEX2XHLYZVMZ7EF9D0G2M0QQQGZ42S3TE0LAR427AFHC02FYHFWSG6AKPC4LLCSC9KHH3W4K5MSHNUXN9SDC",
			wantXPubG1:  "xpublic1pqgqqqqyqqyqqqqpq6axjt5395s9r89me3e25ljxap2qvulmx7s3ufnq2d49z0rhr38yqzv90txq0g9e8jlq8za9yqs8tpvv9nv6hkp0s52dvvhp4m9thxr7hytla2gxcv850hccjd5nvavydhefqvzaww5",
			wantXPubG2:  "xpublic1pqgqqqqyqqyqqqqpqfc96azpj5lsjkc3seg5kuff9q7a9tex2xhlyzvmz7ef9d0g2m0qqqc9477pmk8c3w0lwhvyr79rvt2p5wr5y7fsh0p3vt26m30354ehrj4w3kvj02qdq6f63m9ccqhcxz27qkh5kdjgxpm4s3nec5ln3qdux8laj7epnd98xnk6e7ucahe23yhuet5kkengnn4tdtdyp6w7ww6a2nwzn8e",
		},
		{
			name:        "derivation path: m/0H/1/2H",
			path:        []uint32{h, 1, 2 + h},
			wantXPrivG1: "XSECRET1PQVQQQQYQQYQQQQQZQQQGQGXG02G9WGUD336CLQ7L25X4NPNCE75A4247R2A7S3W9S37XQSQ7FQQJQGS7R7VCA9VE4MX6K8YKWYTZH65JTMJS6HCUT09ZA5VEPZKQ7TWAD83EJ6",
			wantXPrivG2: "XSECRET1PQVQQQQYQQYQQQQQZQQQGQG8PZVKZLK72R0VSGLSAKL4EMX9UW4VL9WZN6G8GXC24877GHGAFKVQZQW0YJPKYNSZLTKHWMZWW6YZ2XTX6SFUZV4XUCYTRGC2YGFR5D7R3N8ZL8R",
			wantXPubG1:  "xpublic1pqvqqqqyqqyqqqqqzqqqgqgxg02g9wgud336clq7l25x4npnce75a4247r2a7s3w9s37xqsq7fqqnpvr9q0w6wls5pprcl39j6pz2pn32kumfr6zf0gmln8gquyrk0qnf32kvaw8x37uu8kmdan9shqm4lcszuz63",
			wantXPubG2:  "xpublic1pqvqqqqyqqyqqqqqzqqqgqg8pzvkzlk72r0vsglsakl4emx9uw4vl9wzn6g8gxc24877ghgafkvqxpq2xrwymg3ks2kkrhsutjwzrv096h3ruczske9a8cl4zfmhl6u8jz0d2eh7hx6jfc3d7lm8vekqcxgf0qnscd0xfl0m8h7jaap3v2u5cel6dxm25pyuq59ntncmnfzmx2xrqrxc4fxrqsvyexmnl7d4g0dgd7v822",
		},
		{
			name:        "derivation path: m/0H/1/2H/2",
			path:        []uint32{h, 1, 2 + h, 2},
			wantXPrivG1: "XSECRET1PQSQQQQYQQYQQQQQZQQQGQQSQQQQZPZWYN98T9Y4TWRN08T5M0ZPVEFVXQCKLYSK269XY7U90VNPXEJJZQYSZDGVU5HLJ76EJSUW7WX4TMPARPNNEEH0RKP2KEW6XDY3FTU9WU9GXSKX25",
			wantXPrivG2: "XSECRET1PQSQQQQYQQYQQQQQZQQQGQQSQQQQZQTZRQ5QNVZ5MDELTWXSKWAXCS7JGA6SNUM44Z06Q5TRL5WE8W9EQQQSR4G0PN2DL90MRRK2MGQ0ZN40SGGTQAHTKEMVKJMJZ4X97SZ6PL2SP66YRM",
			wantXPubG1:  "xpublic1pqsqqqqyqqyqqqqqzqqqgqqsqqqqzpzwyn98t9y4twrn08t5m0zpvefvxqcklysk269xy7u90vnpxejjzqyc2l4vf0y46d093seje3fnn49ha4muml9qzdmu8tgdrar20mqunvr6xt8y5jkh67fx9y4mup2slkh29x3erlu",
			wantXPubG2:  "xpublic1pqsqqqqyqqyqqqqqzqqqgqqsqqqqzqtzrq5qnvz5mdeltwxskwaxcs7jga6snum44z06q5trl5we8w9eqqpsf9vs9vk62q2lcyg5lxtsve3hjx3rda4w294nsv7huwzf3kk5nf72xjeg7vls3qk6kq894skslg3fczf87x55ltv0dkfatgncfqrje5fl40hu84gpnjknssfwsyseuyjvdswtvjzvxmtte6kafur7y8zl2s9y5g8j",
		},
		{
			name:        "derivation path: m/0H/1/2H/2/1000000000",
			path:        []uint32{h, 1, 2 + h, 2, 1000000000},
			wantXPrivG1: "XSECRET1PQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VS9LZ67JKWW0P6TQY9NY58LV0PCVRQQTAEMKGV6ULJNS99Y68JHCVGPYPZTWSAST8PWFJMJQDU0FU8D4YMF58CZ998PGRN29EZYHLWNDVDDJAT3F4D",
			wantXPrivG2: "XSECRET1PQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VS27RYEFRMHGDM0P29ADH63TVTNMRTDS2MF5R23X7T7ULLJS073DTQQYQ4SRMEFWV8TVTR3Z33PM8FG44MU7VLJGDZH9G4LNDELREGZL6NHQCCQPXC",
			wantXPubG1:  "xpublic1pq5qqqqyqqyqqqqqzqqqgqqsqqqqqpj568vs9lz67jkww0p6tqy9ny58lv0pcvrqqtaemkgv6uljns99y68jhcvgpxzvmgpqnpgdwddkajrwl9gjudyh5q4fklms3q3390mtt5ytznugp4kqxtrrpcqulq53aunrwnav2tjqzv47re",
			wantXPubG2:  "xpublic1pq5qqqqyqqyqqqqqzqqqgqqsqqqqqpj568vs27ryefrmhgdm0p29adh63tvtnmrtds2mf5r23x7t7ulljs073dtqqvzc95qdgpsl7gefz0s3a7l3khcddea2hzy05e3gt7rcqcekzkzzdrcwedch3ca2yjm93lq7azy352mshd9l802den6j40f3un0efv69fveej3qh8ht4lq7dy47kjz2gsm6csu523ux9wnrhy547suc3rx24qeppfv3",
		},
	}

	masterKeyG1, _ := NewMaster(testSeed, true)
	masterKeyG2, _ := NewMaster(testSeed, false)
	for no, tt := range tests {
		extKeyG1, _ := masterKeyG1.DerivePath(tt.path)
		neuterKeyG1 := extKeyG1.Neuter()

		extKeyG2, _ := masterKeyG2.DerivePath(tt.path)
		neuterKeyG2 := extKeyG2.Neuter()

		require.Equal(t, tt.wantXPrivG1, extKeyG1.String(), "test %d failed", no)
		require.Equal(t, tt.wantXPubG1, neuterKeyG1.String(), "test %d failed", no)
		require.Equal(t, tt.wantXPrivG2, extKeyG2.String(), "test %d failed", no)
		require.Equal(t, tt.wantXPubG2, neuterKeyG2.String(), "test %d failed", no)

		recoveredExtKeyG1, err := NewKeyFromString(tt.wantXPrivG1)
		require.NoError(t, err)

		recoveredExtKeyG2, err := NewKeyFromString(tt.wantXPrivG2)
		require.NoError(t, err)

		recoveredNeuterKeyG1, err := NewKeyFromString(tt.wantXPubG1)
		require.NoError(t, err)

		recoveredNeuterKeyG2, err := NewKeyFromString(tt.wantXPubG2)
		require.NoError(t, err)

		require.Equal(t, extKeyG1, recoveredExtKeyG1)
		require.Equal(t, extKeyG2, recoveredExtKeyG2)
		require.Equal(t, neuterKeyG1, recoveredNeuterKeyG1)
		require.Equal(t, neuterKeyG2, recoveredNeuterKeyG2)
		require.Equal(t, tt.path, recoveredExtKeyG1.path)
		require.Equal(t, tt.path, recoveredExtKeyG2.path)
		require.Equal(t, tt.path, recoveredNeuterKeyG1.path)
		require.Equal(t, tt.path, recoveredNeuterKeyG2.path)
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
			str:           "XSECRET1PQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VS9LZ67JKWW0P6TQY9NY58LV0PCVRQQTAEMKGV6ULJNS99Y68JHCVGPYPZTWSAST8PWFJMJQDU0FU8D4YMF58CZ998PGRN29EZYHLWNDVDDJAT3FD4",
			expectedError: bech32m.InvalidChecksumError{Expected: "at3f4d", Actual: "at3fd4"},
		},
		{
			desc:          "no depth",
			str:           "XSECRET1P6NTYTF",
			expectedError: io.EOF,
		},
		{
			desc:          "wrong path",
			str:           "XSECRET1PQ5QQQQYQQYQQQQQZQQQGQESEG08",
			expectedError: io.EOF,
		},
		{
			desc:          "no chain code",
			str:           "XSECRET1PQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568V5G6A4P",
			expectedError: io.EOF,
		},
		{
			desc:          "no group",
			str:           "XSECRET1PQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VS9LZ67JKWW0P6TQY9NY58LV0PCVRQQTAEMKGV6ULJNS99Y68JHCVGS579U6",
			expectedError: io.EOF,
		},
		{
			desc:          "no key",
			str:           "XSECRET1PQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VS9LZ67JKWW0P6TQY9NY58LV0PCVRQQTAEMKGV6ULJNS99Y68JHCVGPJCMNNE",
			expectedError: io.EOF,
		},
		{
			desc:          "invalid type",
			str:           "XSECRET1ZQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VS9LZ67JKWW0P6TQY9NY58LV0PCVRQQTAEMKGV6ULJNS99Y68JHCVGPYPZTWSAST8PWFJMJQDU0FU8D4YMF58CZ998PGRN29EZYHLWNDVDDJPJKVRX",
			expectedError: ErrInvalidKeyData,
		},
		{
			desc:          "invalid type",
			str:           "XPUBLIC1ZQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VS9LZ67JKWW0P6TQY9NY58LV0PCVRQQTAEMKGV6ULJNS99Y68JHCVGPYPZTWSAST8PWFJMJQDU0FU8D4YMF58CZ998PGRN29EZYHLWNDVDDJ3HALEC",
			expectedError: ErrInvalidKeyData,
		},
		{
			desc:          "invalid hrp",
			str:           "SECRET1PQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VS9LZ67JKWW0P6TQY9NY58LV0PCVRQQTAEMKGV6ULJNS99Y68JHCVGPYPZTWSAST8PWFJMJQDU0FU8D4YMF58CZ998PGRN29EZYHLWNDVDDJ98PYV5",
			expectedError: ErrInvalidHRP,
		},
		{
			desc:          "invalid hrp",
			str:           "PUBLIC1PQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VS9LZ67JKWW0P6TQY9NY58LV0PCVRQQTAEMKGV6ULJNS99Y68JHCVGPYPZTWSAST8PWFJMJQDU0FU8D4YMF58CZ998PGRN29EZYHLWNDVDDJ4Z2HK2",
			expectedError: ErrInvalidHRP,
		},
	}

	for no, tt := range tests {
		_, err := NewKeyFromString(tt.str)
		assert.ErrorIs(t, err, tt.expectedError, "test %d error is not matched", no)
	}
}

// TestNeuter ensures the Neuter function works as intended.
//
//nolint:lll // Long extended private key
func TestNeuter(t *testing.T) {
	extKey, _ := NewKeyFromString("XSECRET1PQ5QQQQYQQYQQQQQZQQQGQQSQQQQQPJ568VS9LZ67JKWW0P6TQY9NY58LV0PCVRQQTAEMKGV6ULJNS99Y68JHCVGPYPZTWSAST8PWFJMJQDU0FU8D4YMF58CZ998PGRN29EZYHLWNDVDDJAT3F4D")
	neuterKey := extKey.Neuter()
	assert.Equal(t,
		"xpublic1pq5qqqqyqqyqqqqqzqqqgqqsqqqqqpj568vs9lz67jkww0p6tqy9ny58lv0pcvrqqtaemkgv6uljns99y68jhcvgpxzvmgpqnpgdwddkajrwl9gjudyh5q4fklms3q3390mtt5ytznugp4kqxtrrpcqulq53aunrwnav2tjqzv47re",
		neuterKey.String())
	assert.Equal(t, neuterKey, neuterKey.Neuter())
}
