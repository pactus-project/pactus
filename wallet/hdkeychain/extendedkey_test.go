package hdkeychain

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNonHardenedDerivation tests several vectors which derive private keys in
// non hardened mode
func TestNonHardenedDerivation(t *testing.T) {
	testSeed, _ := hex.DecodeString("c7195bf60f19a6ceee5bc4920ff6acf66d64126d743baa27bb9ead164c9b393dd3fc087ae9eafd83e32cb57f24dd653f5d4045dc5f45a4b461fb520cd8a10ca3")

	tests := []struct {
		name     string
		path     []uint32
		wantPriv string
		wantPub  string
	}{
		{
			name:     "test vector 1 chain m",
			path:     []uint32{},
			wantPriv: "SECRET1PF06E3APS3Y5YG0RHC7JTTT4SYJHMT6NGD70QD9C5JFKQR9A0HDTQ44N79F",
			wantPub:  "public1pjxufj2ts7yry49lnjdx2ex3jqtvry5f9sng6pfus2y0qqkx83r24tlc3aq08qkxu0yk4q45e30y3zzpjw87v3z6xrvty45duwq4elm9ac45y7dndy873cvw4wf7htq2dkhu637fqm9f3a7znsmvvlwfmaswt68r8",
		},
		{
			name:     "test vector 1 chain m/0",
			path:     []uint32{0},
			wantPriv: "SECRET1PTSXQ5NP0AY8E8PWNMLDQPJE7PR5THENTKN406EV2NEL6GK9NSUFSFKRZL6",
			wantPub:  "public1pkeaquf7hgnvgv6f59k2s4was58yk2thwe4qs2q82l62yuxg60cuqwe397dl3emkqqksyl25fx4sn79yvnzw46ptpp7hgwvcwqjqlpucc9jg35sn3r4fmfdaq02ezuqx9yh8rr7u6xavrsjphw95n8qvvgyp8mcx6",
		},
		{
			name:     "test vector 1 chain m/0/1",
			path:     []uint32{0, 1},
			wantPriv: "SECRET1PPRMQ9CXCVSWU492U0XYL9S463KWZJ8TE2ULYQ85VJC9EY3GPTS6S70KNGK",
			wantPub:  "public1pjxu9e7nf6nxn66c27utr744m7r4zga9q4g6l3zjdxcx8s36yfwwkmtxxhjw2v90enk7csheu378mqryku0v45xqe39c7mrz9j6g2ldfrzhzy2kf7yc3490syk83lm3g9rwjr5qr4rqsl9rep985xywve2sxq52ak",
		},
		{
			name:     "test vector 1 chain m/0/1/2",
			path:     []uint32{0, 1, 2},
			wantPriv: "SECRET1PR07C9KG87J0T5W0JRK73GCWPY4CLHXGLLJFMM0MU082PR9DG3Q4SPZM5U6",
			wantPub:  "public1p56763hprv80dhr5e6ywvrcpcucn6gtc98deymd4nlrp9m06873ltgrhwjmtqlau9tphhvneq65vacxt7ph9jc7sd6epx0cd05jhh0j3dkmg5n44etzrdm2h43um2x6qs3e5r54ux7saggkx4v7swlyh8753lfuqw",
		},
		{
			name:     "test vector 1 chain m/0/1/2/2",
			path:     []uint32{0, 1, 2, 2},
			wantPriv: "SECRET1PWY9JXXERMUL0PG9SKADSCZA6KSZ2DP7AA6CNY5ZXN8ZW8CKJSFYS85ZM9Q",
			wantPub:  "public1p3uqr94hsm8veulpz3wvy83yh58xvffs0exzzcnshlv26mt7c89hescg3t6nqc42qnqprlflghxjm2qe7clx7x29qcl3td0h6au42g5x204dwv8ak7aj05a46z6tyw47wkal7qv0p4en49sj8usqh0ky90cz89qme",
		},
		{
			name:     "test vector 1 chain m/0/1/2/2/1000000000",
			path:     []uint32{0, 1, 2, 2, 1000000000},
			wantPriv: "SECRET1PDWA39FF0HKFJZLGRMLKRQJA4QZQCR5DCK9J65FN8FK46DFFDWRRSXFSF8C",
			wantPub:  "public1pj8wttg09tm8zvkxvvuekd3872e42wvsj29dh8zwpsx2tvy2seeyt5cgxxm0mxkpcn4kwrtk8g3jp2re3shpxr4x9zznuumn4qxecnfpckyj3evyj8thmlu2pdfzh7kqv7ffgxqwg8d50hhmm93z3t9n92qx402zl",
		},
	}

	for i, test := range tests {
		extKey, _ := NewMaster(testSeed)
		neuterKey, _ := extKey.Neuter()

		for _, childNum := range test.path {
			var err error
			extKey, err = extKey.Derive(childNum)
			require.Equal(t, childNum, extKey.ChildIndex())
			require.NoError(t, err)

			neuterKey, err = neuterKey.Derive(childNum)
			require.Equal(t, childNum, extKey.ChildIndex())
			require.NoError(t, err)
		}

		privKey, err := extKey.BLSPrivateKey()
		require.NoError(t, err)
		require.Equal(t, privKey.String(), test.wantPriv,
			"mismatched serialized private key for test #%v", i+1)

		pubKey, err := extKey.BLSPublicKey()
		require.NoError(t, err)
		require.Equal(t, pubKey.String(), test.wantPub,
			"mismatched serialized public key for test #%v", i+1)

		require.True(t, extKey.IsPrivate())
		require.False(t, neuterKey.IsPrivate())
		neuterPubKey, _ := neuterKey.BLSPublicKey()
		require.True(t, neuterPubKey.EqualsTo(pubKey))
		require.True(t, neuterKey.Address().EqualsTo(pubKey.Address()))
	}
}

// TestHardenedDerivation tests several vectors which derive private keys in
// hardened mode
func TestHardenedDerivation(t *testing.T) {
	testSeed, _ := hex.DecodeString("c7195bf60f19a6ceee5bc4920ff6acf66d64126d743baa27bb9ead164c9b393dd3fc087ae9eafd83e32cb57f24dd653f5d4045dc5f45a4b461fb520cd8a10ca3")

	tests := []struct {
		name     string
		path     []uint32
		wantPriv string
	}{
		// Test vector 1
		{
			name:     "test vector 1 chain m",
			path:     []uint32{},
			wantPriv: "SECRET1PF06E3APS3Y5YG0RHC7JTTT4SYJHMT6NGD70QD9C5JFKQR9A0HDTQ44N79F",
		},
		{
			name:     "test vector 1 chain m/0H",
			path:     []uint32{HardenedKeyStart},
			wantPriv: "SECRET1PW0X8NMEHKZTR79G95ESXL57ZWAH2JYQY0JPARUWHG93U0YQC97TS8CTQ9M",
		},
		{
			name:     "test vector 1 chain m/0H/2H",
			path:     []uint32{HardenedKeyStart, HardenedKeyStart + 2},
			wantPriv: "SECRET1PVT2EKG3US008TD9V8AAGLTV7K70QP2GMXXHFJ64DC0RNULSL2FRSQGGZWF",
		},
		{
			name:     "test vector 1 chain m/0H/2H/2H",
			path:     []uint32{HardenedKeyStart, HardenedKeyStart + 2, HardenedKeyStart + 2},
			wantPriv: "SECRET1PVXKCPT0D662R95GJM4ERQL42J9VXZ0W4PW00FPJZ63APR9RX89SQWCWZPM",
		},
	}

	for i, test := range tests {
		extKey, _ := NewMaster(testSeed)

		var err error
		for _, childNum := range test.path {
			extKey, err = extKey.Derive(childNum)
			require.NoError(t, err)
			require.Equal(t, childNum, extKey.ChildIndex())
		}

		privStr, err := extKey.BLSPrivateKey()
		require.NoError(t, err)
		require.Equal(t, privStr.String(), test.wantPriv,
			"mismatched serialized private key for test #%v", i+1)
		require.True(t, extKey.IsPrivate())
	}
}

// TestInvalidDerivation tests Derive function for invalid data
func TestInvalidDerivation(t *testing.T) {
	t.Run("Deriving hardened from non-hardened", func(t *testing.T) {
		seed, _ := GenerateSeed(32)
		ext, _ := NewMaster(seed)

		derived1, _ := ext.Derive(1)
		_, err := derived1.Derive(1 + HardenedKeyStart)
		assert.ErrorIs(t, err, ErrInvalidChild)
	})

	t.Run("Deriving non-hardened from hardened", func(t *testing.T) {
		seed, _ := GenerateSeed(32)
		ext, _ := NewMaster(seed)

		derived1, _ := ext.Derive(1 + HardenedKeyStart)
		_, err := derived1.Derive(1)
		assert.ErrorIs(t, err, ErrInvalidChild)
	})

	t.Run("Invalid key", func(t *testing.T) {
		key := [31]byte{0}
		chainCode := [32]byte{0}
		ext := NewExtendedKey(key[:], chainCode[:], 0, 0, true)
		_, err := ext.Derive(HardenedKeyStart)
		assert.ErrorIs(t, err, ErrInvalidKey)
	})

	t.Run("Invalid key", func(t *testing.T) {
		key := [95]byte{0}
		chainCode := [32]byte{0}
		ext := NewExtendedKey(key[:], chainCode[:], 0, 0, false)
		_, err := ext.Derive(0)
		assert.ErrorIs(t, err, ErrInvalidKey)
	})

	t.Run("Derive public key from hardened key", func(t *testing.T) {
		key := [32]byte{0}
		chainCode := [32]byte{0}
		ext := NewExtendedKey(key[:], chainCode[:], 0, 0, false)
		_, err := ext.Derive(HardenedKeyStart)
		assert.ErrorIs(t, err, ErrDeriveHardFromPublic)
	})

	t.Run("Derive beyond maximum depth", func(t *testing.T) {
		key := [32]byte{0}
		chainCode := [32]byte{0}
		ext := NewExtendedKey(key[:], chainCode[:], maxUint8, 0, true)
		_, err := ext.Derive(HardenedKeyStart)
		assert.ErrorIs(t, err, ErrDeriveBeyondMaxDepth)
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
