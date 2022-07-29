package hdkeychain

import (
	"encoding/hex"
	"fmt"
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
			wantPriv: "SECRET1PZ5JAWEKSLAJHEM04V9RN0AMLLZ7K54GHSFK25TJNS067UXUS8UDS5L3PDC",
			wantPub:  "public1p3q7xtdhplzefrsgmr34mpdf2nfvulv9qqu8hsfn30tn5phqc85u9g9rqueavr4vx6xwjz2e9a8l5j9e4avyyfy2ugxnk2npwek6xx7j393cly0j8lsja8un3g9wawhmkgvkm3c9g7dcaaxfllfn8yucsqslkvcn0",
		},
		{
			name:     "test vector 1 chain m/0/1",
			path:     []uint32{0, 1},
			wantPriv: "SECRET1P936XUE9V7UC3KAZ4E2AUTSYZE7PPW4SE8V4EG0HW34C06P0ZW3WSGMJC9C",
			wantPub:  "public1phx9c6xke9nggxk6dtd08p27cg2re0ax549rjg40286juqu9eu8kky5wukazpa4rmkudm5uds8gds693jtjajldhkf6vl60m7thtva0n3pm0n52afeeu82rmc43z835ny6pyyrherqywz5trptsegadvllg434xuw",
		},
		{
			name:     "test vector 1 chain m/0/1/2",
			path:     []uint32{0, 1, 2},
			wantPriv: "SECRET1P8SASES40AQ4DJU35LDUKFL8JZ2ZLFSL9D5FT3HKL9KLM2TKX2KCS6ELVP2",
			wantPub:  "public1p3agp8dutddjnsdjxvfrm5q8lhhp802pgyd9cq7eydlty3lc28gxqum3e7xkavv6dv32pugtsupch7pdl63vaddxhq9s87ftu4zcdjyj8zwhvct73098rs6c7ks7n96ujeswl5ygsr3meq39ypup4l6ylygw4ljwy",
		},
		{
			name:     "test vector 1 chain m/0/1/2/2",
			path:     []uint32{0, 1, 2, 2},
			wantPriv: "SECRET1PPPASCFMDW9MQNAVJL74NA90JZ4HCKGNR56A6F20SUF4VSWU0M8XQZUGT2X",
			wantPub:  "public1p5yk990j0l95vzutfrwc4sgzeycressw0mw6h0aswqsxl4lhk4aaeewpsm5ndyt5refl4wdepm4htzxzfrcgqz36dqrpsrklvm2e8uymw5v866tyq4skwnfgsunvpqwfyp5dhfcsz722mjc52lzy53nprqvm89zdq",
		},
		{
			name:     "test vector 1 chain m/0/1/2/2/1000000000",
			path:     []uint32{0, 1, 2, 2, 1000000000},
			wantPriv: "SECRET1PVJK70HLGMUTRK0ZM7DL8AHPYJ7HWQ2JNCE4KXKF7VAYT8PPTZYCQ4LEY4G",
			wantPub:  "public1p3534gvup2jsy932x2k8r6nyl6ynsp6us3840cv6vleqakhrpk5v67u2xyuu9huxt64rpq8gnnrtecz3p2pyasnm5e76fznjdhl2ss4vk2d3m9jmeadml6adqd6q5wjrwmas8g23dk8gu0sw7paerv77aju8e5mhn",
		},
	}

	for i, test := range tests {
		extKey, _ := NewMaster(testSeed)
		neuterKey, _ := extKey.Neuter()

		for _, childNum := range test.path {
			var err error
			extKey, err = extKey.Derive(childNum)
			require.NoError(t, err)

			neuterKey, err = neuterKey.Derive(childNum)
			require.NoError(t, err)
		}

		privKey, err := extKey.BLSPrivateKey()
		fmt.Printf("privKey: %x\n", privKey.Bytes())
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
		fmt.Println(pubKey.String())
		fmt.Println(neuterPubKey.String())
		require.True(t, neuterPubKey.EqualsTo(pubKey))
	}
}

// TestHardenedDerivation tests several vectors which derive private keys in
// hardened mode
func TestHardenedDerivation(t *testing.T) {
	testSeed, _ := hex.DecodeString("c7195bf60f19a6ceee5bc4920ff6acf66d64126d743baa27bb9ead164c9b393dd3fc087ae9eafd83e32cb57f24dd653f5d4045dc5f45a4b461fb520cd8a10ca3")
	hkStart := uint32(0x80000000)

	tests := []struct {
		name     string
		path     []uint32
		wantPriv string
		wantErr  error
	}{
		// Test vector 1
		{
			name:     "test vector 1 chain m",
			path:     []uint32{},
			wantPriv: "SECRET1PF06E3APS3Y5YG0RHC7JTTT4SYJHMT6NGD70QD9C5JFKQR9A0HDTQ44N79F",
			wantErr:  nil,
		},
		{
			name:     "test vector 1 chain m/0H",
			path:     []uint32{hkStart},
			wantPriv: "SECRET1PW0X8NMEHKZTR79G95ESXL57ZWAH2JYQY0JPARUWHG93U0YQC97TS8CTQ9M",
			wantErr:  nil,
		},
		{
			name:     "test vector 1 chain m/0H/1",
			path:     []uint32{hkStart, 1},
			wantPriv: "",
			wantErr:  ErrInvalidChild,
		},
		{
			name:     "test vector 1 chain m/0H/2H",
			path:     []uint32{hkStart, hkStart + 2},
			wantPriv: "SECRET1PVT2EKG3US008TD9V8AAGLTV7K70QP2GMXXHFJ64DC0RNULSL2FRSQGGZWF",
			wantErr:  nil,
		},
		{
			name:     "test vector 1 chain m/0H/2H/2H",
			path:     []uint32{hkStart, hkStart + 2, hkStart + 2},
			wantPriv: "SECRET1PVXKCPT0D662R95GJM4ERQL42J9VXZ0W4PW00FPJZ63APR9RX89SQWCWZPM",
			wantErr:  nil,
		},
	}

	for i, test := range tests {
		extKey, _ := NewMaster(testSeed)

		var err error
		for _, childNum := range test.path {
			extKey, err = extKey.Derive(childNum)
			if err != nil {
				break
			}
		}

		if test.wantErr != nil {
			require.ErrorIs(t, err, test.wantErr)
		} else {
			privStr, err := extKey.BLSPrivateKey()
			require.NoError(t, err)
			require.Equal(t, privStr.String(), test.wantPriv,
				"mismatched serialized private key for test #%v", i+1)
			require.True(t, extKey.IsPrivate())
		}
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
