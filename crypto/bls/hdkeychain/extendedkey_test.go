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
	for i := 0; i < 10000; i++ {
		ts := testsuite.NewTestSuite(t)

		testSeed := ts.RandBytes(32)
		path := []uint32{
			ts.RandUint32(HardenedKeyStart),
			ts.RandUint32(HardenedKeyStart),
			ts.RandUint32(HardenedKeyStart),
			ts.RandUint32(HardenedKeyStart),
		}

		checkPublicKeyDerivation := func(masterKey *ExtendedKey, path []uint32) {
			neuterKey := masterKey.Neuter()

			extKey1, _ := masterKey.DerivePath(path)
			extKey2, _ := neuterKey.DerivePath(path)
			pubKey1 := extKey1.RawPublicKey()
			pubKey2 := extKey2.RawPublicKey()

			require.Equal(t, extKey1.Path(), path)
			require.Equal(t, pubKey1, pubKey2)
		}

		masterKeyG1, _ := NewMaster(testSeed, true)
		masterKeyG2, _ := NewMaster(testSeed, false)

		checkPublicKeyDerivation(masterKeyG1, path)
		checkPublicKeyDerivation(masterKeyG2, path)
	}
}

// TestHardenedDerivation tests derive private key and public key in
// hardened mode.
func TestHardenedDerivation(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	testSeed := ts.RandBytes(32)
	path := []uint32{
		ts.RandUint32(HardenedKeyStart) + HardenedKeyStart,
		ts.RandUint32(HardenedKeyStart) + HardenedKeyStart,
		ts.RandUint32(HardenedKeyStart) + HardenedKeyStart,
		ts.RandUint32(HardenedKeyStart) + HardenedKeyStart,
	}

	masterKey, _ := NewMaster(testSeed, false)
	extKey, _ := masterKey.DerivePath(path)
	privKey, _ := extKey.RawPrivateKey()
	blsPrivKey, _ := bls.PrivateKeyFromBytes(privKey)
	pubKey := extKey.RawPublicKey()

	assert.Equal(t, extKey.Path(), path)
	assert.Equal(t, blsPrivKey.PublicKey().Bytes(), pubKey)
}

// TestDerivation tests derive private keys in hardened and non hardened modes.
func TestDerivation(t *testing.T) {
	testSeed, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	h := HardenedKeyStart
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
	for i, test := range tests {
		extKeyG1, err := masterKeyG1.DerivePath(test.path)
		require.NoError(t, err)

		extKeyG2, err := masterKeyG2.DerivePath(test.path)
		require.NoError(t, err)

		privKeyG1, err := extKeyG1.RawPrivateKey()
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(privKeyG1), test.wantPrivG1,
			"mismatched serialized private key for test #%v", i+1)

		privKeyG2, err := extKeyG2.RawPrivateKey()
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(privKeyG2), test.wantPrivG2,
			"mismatched serialized private key for test #%v", i+1)

		pubKeyG1 := extKeyG1.RawPublicKey()
		require.Equal(t, hex.EncodeToString(pubKeyG1), test.wantPubG1,
			"mismatched serialized public key for test #%v", i+1)

		pubKeyG2 := extKeyG2.RawPublicKey()
		require.Equal(t, hex.EncodeToString(pubKeyG2), test.wantPubG2,
			"mismatched serialized public key for test #%v", i+1)

		neuterKeyG1 := extKeyG1.Neuter()
		neuterKeyG2 := extKeyG2.Neuter()
		neuterPubKeyG1 := neuterKeyG1.RawPublicKey()
		neuterPubKeyG2 := neuterKeyG2.RawPublicKey()

		require.True(t, extKeyG1.IsPrivate())
		require.True(t, extKeyG2.IsPrivate())
		require.False(t, neuterKeyG1.IsPrivate())
		require.False(t, neuterKeyG2.IsPrivate())
		require.Equal(t, neuterPubKeyG1, pubKeyG1)
		require.Equal(t, neuterPubKeyG2, pubKeyG2)
		require.Equal(t, extKeyG1.Path(), test.path)
		require.Equal(t, extKeyG2.Path(), test.path)
		require.Equal(t, neuterKeyG1.Path(), test.path)
		require.Equal(t, neuterKeyG2.Path(), test.path)

		_, err = neuterKeyG1.RawPrivateKey()
		assert.ErrorIs(t, err, ErrNotPrivExtKey)

		_, err = neuterKeyG2.RawPrivateKey()
		assert.ErrorIs(t, err, ErrNotPrivExtKey)

		blsPrivKey, _ := bls.PrivateKeyFromBytes(privKeyG2)
		require.Equal(t, blsPrivKey.PublicKey().Bytes(), pubKeyG2)
	}
}

// TestInvalidDerivation tests Derive function for invalid data.
func TestInvalidDerivation(t *testing.T) {
	t.Run("Private key is 31 bytes. It should be 32 bytes", func(t *testing.T) {
		key := [31]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], []uint32{}, true, false)
		_, err := ext.Derive(HardenedKeyStart)
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
		_, err := ext.Derive(HardenedKeyStart)
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

	for i, test := range tests {
		seed, err := GenerateSeed(test.length)
		assert.ErrorIs(t, err, test.err)

		if test.err == nil && len(seed) != int(test.length) {
			t.Errorf("GenerateSeed #%d (%s): length mismatch -- "+
				"got %d, want %d", i, test.name, len(seed),
				test.length)
			continue
		}
	}
}

// TestKeyToString ensures the String function works as intended.
//
//nolint:lll
func TestKeyToString(t *testing.T) {
	testSeed, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	h := HardenedKeyStart
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
			wantXPrivG1: "XSECRET1PQZU8NVYHHG5E99FQ4YW7U2W7RK2RNRY3QA4YY3D7V9CYYEWJXRYHYVZ02H33ACWY7K90PPQ06067VD0ADSR74NG59Q7YT475XU5SQW4MSSW6Q78L",
			wantXPrivG2: "XSECRET1PQZU8NVYHHG5E99FQ4YW7U2W7RK2RNRY3QA4YY3D7V9CYYEWJXRYHYCZ02H33ACWY7K90PPQ06067VD0ADSR74NG59Q7YT475XU5SQW4MSS4TGGAS",
			wantXPubG1:  "xpublic1pqzu8nvyhhg5e99fq4yw7u2w7rk2rnry3qa4yy3d7v9cyyewjxryhyvy0hmvggfvgkc5nwlq2p5xe23afact4yl2l6mfvvzgrf2xrcp6dmgp3ur073p4523yehls0gznufvvq3m96zp",
			wantXPubG2:  "xpublic1pqzu8nvyhhg5e99fq4yw7u2w7rk2rnry3qa4yy3d7v9cyyewjxryhyc93htfm7jj2ap7gnhkzcvj3ycpu5z8zmd3vl539fjttle6sdr66nrnufntaxl8sf9kadeuhq0nu3rjsg677e8yfdme26qcqj6au7u7xelch4hfa49fs7gjfryqlmalaypmvpuyw5ddylk4qpeavd599gshr2mw2pg",
		},
		{
			name:        "derivation path: m/0H",
			path:        []uint32{h},
			wantXPrivG1: "XSECRET1PQXQGPQYQPQDNX9T02WPS2RZ5SYUKE3JPHE8RGDHJMTNU7684679WEQWRN8STWVZLT4AL4EL2HUKV87HTCYJYNCW8Z9KZWA7HUWZW44UA72VKV7UDNGKV5QYH",
			wantXPrivG2: "XSECRET1PQXQGPQYQPR38R7SGQNLMC6H96CANRN8XE3WVFVLF0V5XW2LE0FDSPYT52FUNSCZKJKA9PPAZ07XQ6UNSG4GSGEVTYDNM36G2KMML27K8EC3DFF5RDS8L8040",
			wantXPubG1:  "xpublic1pqxqgpqyqpqdnx9t02wps2rz5syuke3jphe8rgdhjmtnu7684679weqwrn8stwv9jsf4gng30ase5n4j0gdu6r66kx2ctx3deskmnsvj2twxmvspswssjq8h7x6hxerrrn5edgyjyj6hqt7m8st",
			wantXPubG2:  "xpublic1pqxqgpqyqpr38r7sgqnlmc6h96canrn8xe3wvfvlf0v5xw2le0fdspyt52funsc9n0k3sspnzem4h7pegnqq62mju24w5zdp545ykq7wqsn92zckg6gjgj8mgs95jradaz3f67lggt0zqqdqav88yjmlmz8x3p78fq53yglad5xjlv3ky27t7qprqjfv8du9k8aqz80e8aq5x3rmmfhvr8ejp8hhlk0",
		},
		{
			name:        "derivation path: m/0H/1",
			path:        []uint32{h, 1},
			wantXPrivG1: "XSECRET1PQ2QGPQYQPQQAWNF96GJ6GZ3NJAUCU420ERWS4QXW0AN0GG7YES9X6J383M3CNJPS804888Y6Y62M5JH4V67R7289CCK63EEPH9MHP8UAFYHHZ2DCX5SSQVC0PJ",
			wantXPrivG2: "XSECRET1PQ2QGPQYQPQQ5UZAW3QE20CFTVGCV52TWY5JS0WJ4UN9RTLJPXD30V5JKH59DHSRQ242Z909LL5W4TM4XLPAFYJA96PRTKC8ZHLLZRQK67796K6NWZ70SJFA5Q2",
			wantXPubG1:  "xpublic1pq2qgpqyqpqqawnf96gj6gz3njaucu420erws4qxw0an0gg7yes9x6j383m3cnjps4avcpaqhy7tuqut55szqav93skdn27c97z3f43juxhv4wuc06u30l4fqmps737lrzfkjdn4s3kl9ygk9k4n",
			wantXPubG2:  "xpublic1pq2qgpqyqpqq5uzaw3qe20cftvgcv52twy5js0wj4un9rtljpxd30v5jkh59dhsrqkhmc8wclz9ela6ass0c5d3dgx3cwsnexzaux93dttw97xjhxuw246xejfagp5rf828vhrqzlqcftcz67jekfqc8wkzx08zn7wyphscllktmyxd55u6wmt8mnrkl92yjln9wj6mxdzww4d4d5s8fmeemt4gnl9wl9",
		},
		{
			name:        "derivation path: m/0H/1/2H",
			path:        []uint32{h, 1, 2 + h},
			wantXPrivG1: "XSECRET1PQWQGPQYQPQQC9QYQSQYVS75S2U3CMRR437PA74GDTXR83NAFM24TUX4MAPZUTPRUVPQPUJPSYG0PLXVWJKV6ANDTRJT8Z93TA2F9AEGDTUW9HJ3W6XVS3TQ09HWS94JH9W",
			wantXPrivG2: "XSECRET1PQWQGPQYQPQQC9QYQSQYWZYEV9LDU5X7EQ3LPMDLTNKVTCA2E72U985SWSDS420AU3W36NVMQ88JFQMZFCP04MTHD388DZP9R9NDGY7PX2NWVZ935V9ZYY36XLPCS6JPKQ5",
			wantXPubG1:  "xpublic1pqwqgpqyqpqqc9qyqsqyvs75s2u3cmrr437pa74gdtxr83nafm24tux4mapzutpruvpqpujpskpjs8hd80c2qs3u0cjedq39qec4twd53apyh5dlen5qwzpm8sf5c4txwhrnglwwrmdk7ejctsd6lup8j3ul",
			wantXPubG2:  "xpublic1pqwqgpqyqpqqc9qyqsqywzyev9ldu5x7eq3lpmdltnkvtca2e72u985swsds420au3w36nvmqs9rphzd5gmg9ttpmcw9e8ppk8jatc37vpgtvj7nu063yamlawrep8k4vmltndfyugkl0ankvmqvryyhsfcvxhnylhanmlfw7sck9w2vvlaxnd42qjwq2ze4eude53dn9rpspnv25npsgxzvndellx658k52wvpgk",
		},
		{
			name:        "derivation path: m/0H/1/2H/2",
			path:        []uint32{h, 1, 2 + h, 2},
			wantXPrivG1: "XSECRET1PQJQGPQYQPQQC9QYQSQYQ9ZWYN98T9Y4TWRN08T5M0ZPVEFVXQCKLYSK269XY7U90VNPXEJJZXQN2R899LUHKKV58RHN3427C0GCVU7WDMCAS24KTK3NFY22LPTHP2MHG35K",
			wantXPrivG2: "XSECRET1PQJQGPQYQPQQC9QYQSQYQYTZRQ5QNVZ5MDELTWXSKWAXCS7JGA6SNUM44Z06Q5TRL5WE8W9EQVQA2RCV6N0ET7CCAJK6QRC5ATUZZZC8D6AKWM95KUS4F305QKS065XR6F35",
			wantXPubG1:  "xpublic1pqjqgpqyqpqqc9qyqsqyq9zwyn98t9y4twrn08t5m0zpvefvxqcklysk269xy7u90vnpxejjzxzhatzte9wntevvxvkv2vuafdldwlxlegqnwlp66rglg6n7c8ymq73jee9y447hjf3f9wlq258a463gz9wywh",
			wantXPubG2:  "xpublic1pqjqgpqyqpqqc9qyqsqyqytzrq5qnvz5mdeltwxskwaxcs7jga6snum44z06q5trl5we8w9eqvzftypt9kjszh7pz98ejurxvdu35gm0dtj3dvur84lrsjvd44y60j35k28n8uyg9k4spedv9586y2wqjfl34986mrmdj026y7zgqukdz0atalpa2qvu45uyzt5pyx0pynrvrjmysnpk667w4h20ql3pch65qv8476a",
		},
		{
			name:        "derivation path: m/0H/1/2H/2/1000000000",
			path:        []uint32{h, 1, 2 + h, 2, 1000000000},
			wantXPrivG1: "XSECRET1PQKQGPQYQPQQC9QYQSQYQ9QY5A0WQXHUTT62EEEU8FVQSKVJSLA3U8PSVQP0H8WEPNTN72WQ55NG72LP3XPZTWSAST8PWFJMJQDU0FU8D4YMF58CZ998PGRN29EZYHLWNDVDDJCEH9HC",
			wantXPrivG2: "XSECRET1PQKQGPQYQPQQC9QYQSQYQ9QY5A0WQ8TCVN9Y0WAPHDU9GH4KL29D3W0VDDKPTDXSD2YME0MNL72PL694VVQ4SRMEFWV8TVTR3Z33PM8FG44MU7VLJGDZH9G4LNDELREGZL6NHQSD2RKK",
			wantXPubG1:  "xpublic1pqkqgpqyqpqqc9qyqsqyq9qy5a0wqxhutt62eeeu8fvqskvjsla3u8psvqp0h8wepntn72wq55ng72lp3xzvmgpqnpgdwddkajrwl9gjudyh5q4fklms3q3390mtt5ytznugp4kqxtrrpcqulq53aunrwnav2tjqt3xrwy",
			wantXPubG2:  "xpublic1pqkqgpqyqpqqc9qyqsqyq9qy5a0wq8tcvn9y0waphdu9gh4kl29d3w0vddkptdxsd2yme0mnl72pl694vvzc95qdgpsl7gefz0s3a7l3khcddea2hzy05e3gt7rcqcekzkzzdrcwedch3ca2yjm93lq7azy352mshd9l802den6j40f3un0efv69fveej3qh8ht4lq7dy47kjz2gsm6csu523ux9wnrhy547suc3rx24qg487k2",
		},
	}

	masterKeyG1, _ := NewMaster(testSeed, true)
	masterKeyG2, _ := NewMaster(testSeed, false)
	for i, test := range tests {
		extKeyG1, _ := masterKeyG1.DerivePath(test.path)
		neuterKeyG1 := extKeyG1.Neuter()

		extKeyG2, _ := masterKeyG2.DerivePath(test.path)
		neuterKeyG2 := extKeyG2.Neuter()

		require.Equal(t, extKeyG1.String(), test.wantXPrivG1, "test %d failed", i)
		require.Equal(t, neuterKeyG1.String(), test.wantXPubG1, "test %d failed", i)
		require.Equal(t, extKeyG2.String(), test.wantXPrivG2, "test %d failed", i)
		require.Equal(t, neuterKeyG2.String(), test.wantXPubG2, "test %d failed", i)

		recoveredExtKeyG1, err := NewKeyFromString(test.wantXPrivG1)
		require.NoError(t, err)

		recoveredExtKeyG2, err := NewKeyFromString(test.wantXPrivG2)
		require.NoError(t, err)

		recoveredNeuterKeyG1, err := NewKeyFromString(test.wantXPubG1)
		require.NoError(t, err)

		recoveredNeuterKeyG2, err := NewKeyFromString(test.wantXPubG2)
		require.NoError(t, err)

		require.Equal(t, extKeyG1, recoveredExtKeyG1)
		require.Equal(t, extKeyG2, recoveredExtKeyG2)
		require.Equal(t, neuterKeyG1, recoveredNeuterKeyG1)
		require.Equal(t, neuterKeyG2, recoveredNeuterKeyG2)
	}
}

// TestInvalidString checks errors corresponding to the invalid strings
//
//nolint:lll
func TestInvalidString(t *testing.T) {
	tests := []struct {
		desc          string
		str           string
		expectedError error
	}{
		{
			desc:          "invalid checksum",
			str:           "XSECRET1PQP0R4SGK8Y84J2G9LQD2E4W5RYXQRPWKYSG8TT4KUMZD0QF7TT8PSVPNCFWFXK8JKWKNMH8HQC8PV0ZMYL36LRFJJ76K3C94YL38FA7PNGF4LRNP",
			expectedError: bech32m.InvalidChecksumError{Expected: "f4lrnq", Actual: "f4lrnp"},
		},
		{
			desc:          "no path len",
			str:           "XSECRET1P6NTYTF",
			expectedError: io.EOF,
		},
		{
			desc:          "wrong path",
			str:           "XSECRET1PQ2QGPQYQPQZ0DRED",
			expectedError: io.EOF,
		},
		{
			desc:          "no key",
			str:           "XSECRET1PQ2QGPQYQPQQLSPCJLQWPR645ALQJ7F6297W6CDEJRSWW5F2MVV4DMY903MVCJSCAN908R",
			expectedError: ErrInvalidKeyData,
		},
		{
			desc:          "invalid type",
			str:           "XSECRET1LQ2QGPQYQPQQLSPCJLQWPR645ALQJ7F6297W6CDEJRSWW5F2MVV4DMY903MVCJS6JSJ5CVRJREG4940E8JSMDPU3HVVT7UA7A5N72AJY9TRNVZZV25QL9W9UJ",
			expectedError: ErrInvalidKeyData,
		},
		{
			desc:          "no key",
			str:           "xpublic1pq2qgpqyqpqqlspcjlqwpr645alqj7f6297w6cdejrsww5f2mvv4dmy903mvcjscn8ga0y",
			expectedError: ErrInvalidKeyData,
		},
		{
			desc:          "invalid type",
			str:           "xpublic1lq2qgpqyqpqqlspcjlqwpr645alqj7f6297w6cdejrsww5f2mvv4dmy903mvcjsu99mrrh4yjc28dxljnyqfr5n0t85jw49fu775gpnde6faj967sr2t2xld8ucaggthxk3em7ptfdwwpvjs5240z5gpngprjaswsju38rr3myqpz2vguspp0cvrf7lkx4r40uv5smrky2v52qypk0njjfe4phxhhu2",
			expectedError: ErrInvalidKeyData,
		},
		{
			str:           "SECRET1PQ2QGPQYQPQQLSPCJLQWPR645ALQJ7F6297W6CDEJRSWW5F2MVV4DMY903MVCJS6JSJ5CVRJREG4940E8JSMDPU3HVVT7UA7A5N72AJY9TRNVZZV25QMJLMS4",
			expectedError: ErrInvalidKeyData,
		},
	}

	for i, test := range tests {
		_, err := NewKeyFromString(test.str)
		assert.ErrorIs(t, err, test.expectedError, "test %d error is not matched", i)
	}
}

// TestNeuter ensures the Neuter function works as intended.
//
//nolint:lll
func TestNeuter(t *testing.T) {
	extKey, _ := NewKeyFromString("XSECRET1PQP0R4SGK8Y84J2G9LQD2E4W5RYXQRPWKYSG8TT4KUMZD0QF7TT8PSVPNCFWFXK8JKWKNMH8HQC8PV0ZMYL36LRFJJ76K3C94YL38FA7PNGF4LRNQ")
	neuterKey := extKey.Neuter()
	assert.Equal(t, neuterKey.String(), "xpublic1pqp0r4sgk8y84j2g9lqd2e4w5ryxqrpwkysg8tt4kumzd0qf7tt8psv9yvfyjsl2hv2ya96vam0uwhdq3t753htdv7fwp694njg7ctvnprtrnmddrc9083nrcnvlg8ex9kucs39tdxg")
	assert.Equal(t, neuterKey, neuterKey.Neuter())
}
