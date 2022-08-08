package hdkeychain

import (
	"encoding/hex"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/util/bech32m"
)

// TestNonHardenedDerivation tests derive private keys in non hardened mode
func TestNonHardenedDerivation(t *testing.T) {
	testSeed, _ := hex.DecodeString("c7195bf60f19a6ceee5bc4920ff6acf66d64126d743baa27bb9ead164c9b393dd3fc087ae9eafd83e32cb57f24dd653f5d4045dc5f45a4b461fb520cd8a10ca3")
	tests := []struct {
		name     string
		path     Path
		wantPriv string
		wantPub  string
	}{
		{
			name:     "derivation path: m",
			path:     Path{},
			wantPriv: "SECRET1PF06E3APS3Y5YG0RHC7JTTT4SYJHMT6NGD70QD9C5JFKQR9A0HDTQ44N79F",
			wantPub:  "public1pjxufj2ts7yry49lnjdx2ex3jqtvry5f9sng6pfus2y0qqkx83r24tlc3aq08qkxu0yk4q45e30y3zzpjw87v3z6xrvty45duwq4elm9ac45y7dndy873cvw4wf7htq2dkhu637fqm9f3a7znsmvvlwfmaswt68r8",
		},
		{
			name:     "derivation path: m/0",
			path:     Path{0},
			wantPriv: "SECRET1P9EWLCAGDLKVDV9XTZJMH5YTMQSUD7SDH7QKEKJQ6P6THHLMNVPPS70CSS0",
			wantPub:  "public1ps3wcrjhqk8mhp8rhlunvyze8wa509tl5mkz8s8fa86pspvzkzcahcvga6q9dtar0mswc9x3wkqpt2907ycj54mlq57vpwat5hcltqsnnjywm9h7x8he59wshe4dnvqqyxzuj9e2s2h56mrywcxm67hsu858epzkd",
		},
		{
			name:     "derivation path: m/0/1",
			path:     Path{0, 1},
			wantPriv: "SECRET1PD7TUGN5Y0WQU3PDWWFXELW03YJ4JHD3MLT594NP4CFN09TUSPAAS3GU38U",
			wantPub:  "public1pjlq2nfuz6v23ry589ev59spjvucmyarvhj5kvfs7xyq26djazrfefyz04npxseamtfvrvkknnxemqq672ft9rf2vfljpxr5tzhexvpg5ptmsrg2kexk869nhe02666p60ssck4up402the5uv4jrdll55syw49a7",
		},
		{
			name:     "derivation path: m/0/1/2",
			path:     Path{0, 1, 2},
			wantPriv: "SECRET1P9MZEFVAN89Y7TLHYN5D7J487VYPM40KYU9C8YLVTYG3VSN8EPQFSK3W9ZJ",
			wantPub:  "public1psjhd8k6wv8rpx4s2ep82rlst93h9j3jrv576n9f3dwqhl275279ahxg2yvwj8cuzhnzk9qpxtf20sxxerwtmvkxdty3pyah5p4f960l25v3npk5533e54utejw3g6t9yxec5u48ny96s7thea6kaye38xv8hex68",
		},
		{
			name:     "derivation path: m/0/1/2/2",
			path:     Path{0, 1, 2, 2},
			wantPriv: "SECRET1PWTRVTPYRZ3NALLR6QF04V5Q9Q5GDURA4XFMWJE6Y49VKFPKWSPKQU4MTQV",
			wantPub:  "public1pjapnwas5tlh0az66ugjnwe79wuuwkksgz9fa5hwx47lnh20eaq269mscxg2xxwaza02w7kzye25t7xpa2qdf45v2cjc0dswk2ezp09ff2xlhv7m2rh02dfmrl8g9jl9n58629gpmxa3s0tl47u84jp9ecuqhwrtm",
		},
		{
			name:     "derivation path: m/0/1/2/2/1000000000",
			path:     Path{0, 1, 2, 2, 1000000000},
			wantPriv: "SECRET1PXP68X5L58APYQX0V9DAUYHA3WELALF80R6SC4PDFM358VHD6RXTQ2WCUWT",
			wantPub:  "public1p5edx5vcuzjfv9wrwyh6rp9wa7asyj32hvww7r5383ry4q5t4npmtu475se82xwmlj5jx2kzkge67s9q6mhqzj56unf5pujj702ahjh6gxdzn29nnxgxa6z9l32sndau42mug4xcwdf8p45jx5r3jx2se35l5de60",
		},
	}

	masterKey, _ := NewMaster(testSeed)
	neuteredMasterKey := masterKey.Neuter()
	for i, test := range tests {
		extKey, err := masterKey.DerivePath(test.path)
		require.NoError(t, err)

		neuterKey, err := neuteredMasterKey.DerivePath(test.path)
		require.NoError(t, err)

		privKey, err := extKey.BLSPrivateKey()
		require.NoError(t, err)
		require.Equal(t, privKey.String(), test.wantPriv,
			"mismatched serialized private key for test #%v", i+1)

		pubKey := extKey.BLSPublicKey()
		require.Equal(t, pubKey.String(), test.wantPub,
			"mismatched serialized public key for test #%v", i+1)

		require.True(t, extKey.IsPrivate())
		require.False(t, neuterKey.IsPrivate())
		neuterPubKey := neuterKey.BLSPublicKey()
		require.True(t, neuterPubKey.EqualsTo(pubKey))
		require.True(t, neuterKey.Address().EqualsTo(pubKey.Address()))
		require.Equal(t, extKey.Path(), test.path)
		require.Equal(t, neuterKey.Path(), test.path)

		_, err = neuterKey.BLSPrivateKey()
		assert.ErrorIs(t, err, ErrNotPrivExtKey)
	}
}

// TestDerivation tests derive private keys in hardened and non hardened modes.
func TestDerivation(t *testing.T) {
	testSeed, _ := hex.DecodeString("c7195bf60f19a6ceee5bc4920ff6acf66d64126d743baa27bb9ead164c9b393dd3fc087ae9eafd83e32cb57f24dd653f5d4045dc5f45a4b461fb520cd8a10ca3")
	h := HardenedKeyStart
	tests := []struct {
		name     string
		path     Path
		wantPriv string
		wantPub  string
	}{
		{
			name:     "derivation path: m",
			path:     Path{},
			wantPriv: "SECRET1PF06E3APS3Y5YG0RHC7JTTT4SYJHMT6NGD70QD9C5JFKQR9A0HDTQ44N79F",
			wantPub:  "public1pjxufj2ts7yry49lnjdx2ex3jqtvry5f9sng6pfus2y0qqkx83r24tlc3aq08qkxu0yk4q45e30y3zzpjw87v3z6xrvty45duwq4elm9ac45y7dndy873cvw4wf7htq2dkhu637fqm9f3a7znsmvvlwfmaswt68r8",
		},
		{
			name:     "derivation path: m/0H",
			path:     Path{h},
			wantPriv: "SECRET1PQTLW8H0D3RN77DXN8XS8XQH4N2AT8SHUYGH507VXD8CXP387J3ASXMYVUM",
			wantPub:  "public1pkgp40pcrj9rv7lz52d88smuqh2kwu4sfxjr40vl0jrmgpza5lszfcvhfmqfxetdgdgtnvdqva3mljqkkdgtkchtrltn5sp466k9kvz4mjhunufyulpuan2w03h4qscctdycvw3tr782n8rum724s6y2x5qddkqq7",
		},
		{
			name:     "derivation path: m/0H/1",
			path:     Path{h, 1},
			wantPriv: "SECRET1P22Z2NPSWG09Z5K4LY72RD58JXA330MNHMKJ0ETKGS4VWDSGF32SQSJKTSA",
			wantPub:  "public1ps5hvvw75jtpga5m72vspywjdav7jf6548nm63qxdh8f8kght6qdfdgma5lnr4ppwu66880c9d94ec9j2z324u23qxdqywtkp6ztjyuvw8vsqyff3rjqy9lpsd8m7c65w4l3jjrvwc3fj3gqsxe7w2f8x5ynv3ju2",
		},
		{
			name:     "derivation path: m/0H/1/2H",
			path:     Path{h, 1, 2 + h},
			wantPriv: "SECRET1PRSRDLLE030ATWX4GAEY0U000JX4QTX42PZQNMPEV26NECY2N2DYSTC23RR",
			wantPub:  "public1p5ffwfn22lec0hwdz4sr8wjnyasu0ujeq7h6upfny4h6ltxun4m006k9mnegcpjjj0qlne46ud0cmsre5lnu48p84kcrt3vh8p7hqckhg2d09g5mf93pt2yx376sa7upjlwcvw45xtk5zv539lk7nq88thq70t390",
		},
		{
			name:     "derivation path: m/0H/1/2H/2",
			path:     Path{h, 1, 2 + h, 2},
			wantPriv: "SECRET1PPACZ03VNYU4FQ63GYTMJ6TRCGLYQJ77SF9EMSF29Q4PKHQAMS3CQ7TCV37",
			wantPub:  "public1pshjth962p7atqjtswgawhmfwpmpn2ltsy7f5xs8e88d0aftqc252ra829xswlpq3hqtf8vaq8jzr59pkvzcfwja4k530gh0w5qcz3eqn9egtpns2pynfldt09l6y8k3zd4p3dnuvy5l4syac6mw4lx2zgsljpcq5",
		},
		{
			name:     "derivation path: m/0H/1/2H/2/1000000000",
			path:     Path{h, 1, 2 + h, 2, 1000000000},
			wantPriv: "SECRET1PY8YA4MUZWE3UEWNH3L6N9TXF4S4PFVMU90NX08FS4MAWS8XEHGVSKC8MUP",
			wantPub:  "public1p30dx22dx20ffvpn4mv5wqgxf59x5kt6udt9jceawqjnl6s0h5d00er6avj5j9lue3etp5aw45dq8kqpe08g556fgva87f9x0mq55vn0nx2zse5vjqv3zjwr6yv25rhmtf6mvx24ne886dl45h9jxyxwnhuxhkh7p",
		},
	}

	masterKey, _ := NewMaster(testSeed)
	for i, test := range tests {
		extKey, err := masterKey.DerivePath(test.path)
		require.NoError(t, err)

		privKey, err := extKey.BLSPrivateKey()
		require.NoError(t, err)
		require.Equal(t, privKey.String(), test.wantPriv,
			"mismatched serialized private key for test #%v", i+1)

		pubKey := extKey.BLSPublicKey()
		require.Equal(t, pubKey.String(), test.wantPub,
			"mismatched serialized public key for test #%v", i+1)

		require.True(t, extKey.IsPrivate())

		neuterKey := extKey.Neuter()
		require.False(t, neuterKey.IsPrivate())
		neuterPubKey := neuterKey.BLSPublicKey()
		require.True(t, neuterPubKey.EqualsTo(pubKey))
		require.True(t, neuterKey.Address().EqualsTo(pubKey.Address()))
		require.ElementsMatch(t, extKey.Path(), test.path)
		require.ElementsMatch(t, neuterKey.Path(), test.path)

		_, err = neuterKey.BLSPrivateKey()
		assert.ErrorIs(t, err, ErrNotPrivExtKey)
	}
}

// TestInvalidDerivation tests Derive function for invalid data
func TestInvalidDerivation(t *testing.T) {
	t.Run("Invalid key", func(t *testing.T) {
		key := [31]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], Path{}, true)
		_, err := ext.Derive(HardenedKeyStart)
		assert.ErrorIs(t, err, ErrInvalidKeyData)
	})

	t.Run("Invalid key", func(t *testing.T) {
		key := [95]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], Path{}, false)
		_, err := ext.Derive(0)
		assert.ErrorIs(t, err, ErrInvalidKeyData)
	})

	t.Run("Derive public key from hardened key", func(t *testing.T) {
		key := [32]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], Path{}, false)
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
func TestKeyToString(t *testing.T) {
	testSeed, _ := hex.DecodeString("c7195bf60f19a6ceee5bc4920ff6acf66d64126d743baa27bb9ead164c9b393dd3fc087ae9eafd83e32cb57f24dd653f5d4045dc5f45a4b461fb520cd8a10ca3")
	h := HardenedKeyStart
	tests := []struct {
		name      string
		path      Path
		wantXPriv string
		wantXPub  string
	}{
		{
			name:      "derivation path: m",
			path:      Path{},
			wantXPriv: "XSECRET1PQRLH34S0AV6M6DQ273T5G0CZJ9P597K8MHUSLK7EMSV5QT0X34FGWJL4NR6RPZFGGS7803AYKKHTQF90KH4XSMU7Q6T3FYNVQXT6LW6K2QT5GC",
			wantXPub:  "xpublic1pqrlh34s0av6m6dq273t5g0czj9p597k8mhuslk7emsv5qt0x34fg0ydcny5hpugxf2tl8y6v4jdryqkcxfgjtpx35zneq5g7qpvv0zx42hl3r6q7wpvdc7fd2ptfnz7fzyyryu0uez95vxckftgmcuptnlktm3tgfumx6g0arsca2unawkq5md0e4rujpk2nrmu98pkce7unhmqtxyp6w",
		},
		{
			name:      "derivation path: m/0H",
			path:      Path{h},
			wantXPriv: "XSECRET1PQXQGPQYQPRM5HCPEAF8A7L2SZ869UKDVJMFD0DLS2256RN6VTG3DLQK2HV3RQQH7U0W7MZ88AU6DXWDQWVP0TX46K0P0CG30GLUCV60SVRZ0A9RMMC4JYH",
			wantXPub:  "xpublic1pqxqgpqyqprm5hcpeaf8a7l2sz869ukdvjmfd0dls2256rn6vtg3dlqk2hv3rpvsr27rs8y2xea79g56w0phcpw4vaetqjdy827e7ly8ksz9mflqynsewnkqjdjk6s6shxc6qemrhlypdv6shd3wk87h8fqrt44vtvc9th90e8cjfe7remx5ulr02pp3sk6fscazk8uw4xw8ehu4tp5g5dgqm5xfh7",
		},
		{
			name:      "derivation path: m/0H/1",
			path:      Path{h, 1},
			wantXPriv: "XSECRET1PQ2QGPQYQPQQLSPCJLQWPR645ALQJ7F6297W6CDEJRSWW5F2MVV4DMY903MVCJS6JSJ5CVRJREG4940E8JSMDPU3HVVT7UA7A5N72AJY9TRNVZZV25QHE8RCZ",
			wantXPub:  "xpublic1pq2qgpqyqpqqlspcjlqwpr645alqj7f6297w6cdejrsww5f2mvv4dmy903mvcjsu99mrrh4yjc28dxljnyqfr5n0t85jw49fu775gpnde6faj967sr2t2xld8ucaggthxk3em7ptfdwwpvjs5240z5gpngprjaswsju38rr3myqpz2vguspp0cvrf7lkx4r40uv5smrky2v52qypk0njjfe4pqjrgak",
		},
		{
			name:      "derivation path: m/0H/1/2H",
			path:      Path{h, 1, 2 + h},
			wantXPriv: "XSECRET1PQWQGPQYQPQQC9QYQSQY0H8GERGSUYVJC3NMUS793F94FYLN4LHKG2LKR2NMEVTFCYNS5LTGUQM0L7TUTL2M3428WFRLRMMU34GZE42SGSY7CWTZK57WPZ56NFY5T96P7",
			wantXPub:  "xpublic1pqwqgpqyqpqqc9qyqsqy0h8gergsuyvjc3nmus793f94fyln4lhkg2lkr2nmevtfcyns5ltdz2tjv6jh7wramng4vqem55e8v8rlykg847hq2ve9d7h6ehyawmm743wu72xqv55nc8u7dwhrt7xuq7d8ul9fcfadkq6ut9ec04cx946znte29x6fvg263p50k580hqvhmkrr4dpja4qn9yf0ah5cpe6acc2vyqg",
		},
		{
			name:      "derivation path: m/0H/1/2H/2",
			path:      Path{h, 1, 2 + h, 2},
			wantXPriv: "XSECRET1PQJQGPQYQPQQC9QYQSQYQ9J7EM47JQN338TXZ8UTRHKUMR45WZ2HXLZWL9NEXMPDWUFQ7PP7QPACZ03VNYU4FQ63GYTMJ6TRCGLYQJ77SF9EMSF29Q4PKHQAMS3CQL5R976",
			wantXPub:  "xpublic1pqjqgpqyqpqqc9qyqsqyq9j7em47jqn338txz8utrhkumr45wz2hxlzwl9nexmpdwufq7pp7qshjth962p7atqjtswgawhmfwpmpn2ltsy7f5xs8e88d0aftqc252ra829xswlpq3hqtf8vaq8jzr59pkvzcfwja4k530gh0w5qcz3eqn9egtpns2pynfldt09l6y8k3zd4p3dnuvy5l4syac6mw4lx2zgshnh8ag",
		},
		{
			name:      "derivation path: m/0H/1/2H/2/1000000000",
			path:      Path{h, 1, 2 + h, 2, 1000000000},
			wantXPriv: "XSECRET1PQKQGPQYQPQQC9QYQSQYQ9QY5A0WQXZPJWCN5JTGU8CHX8QWJEHSDAYUYPXL9MR6AFFANTQSEMLD4DXSQY8YA4MUZWE3UEWNH3L6N9TXF4S4PFVMU90NX08FS4MAWS8XEHGVSDM6D0R",
			wantXPub:  "xpublic1pqkqgpqyqpqqc9qyqsqyq9qy5a0wqxzpjwcn5jtgu8chx8qwjehsdayuypxl9mr6affantqsemld4dxsq30dx22dx20ffvpn4mv5wqgxf59x5kt6udt9jceawqjnl6s0h5d00er6avj5j9lue3etp5aw45dq8kqpe08g556fgva87f9x0mq55vn0nx2zse5vjqv3zjwr6yv25rhmtf6mvx24ne886dl45h9jxyxwnhucftqpz",
		},
	}

	masterKey, _ := NewMaster(testSeed)
	for i, test := range tests {
		extKey, _ := masterKey.DerivePath(test.path)
		neuterKey := extKey.Neuter()

		assert.Equal(t, extKey.String(), test.wantXPriv, "test %d failed", i)
		assert.Equal(t, neuterKey.String(), test.wantXPub, "test %d failed", i)

		extKey2, err := NewKeyFromString(test.wantXPriv)
		assert.NoError(t, err)

		neuterKey2, err := NewKeyFromString(test.wantXPub)
		assert.NoError(t, err)

		assert.Equal(t, extKey, extKey2)
		assert.Equal(t, neuterKey, neuterKey2)
	}
}

// TestNewKeyFromString ensures the NewKeyFromString function works as intended.
func TestNewKeyFromString(t *testing.T) {
	extKey, _ := NewKeyFromString("XSECRET1PQXQGPQYQPRM5HCPEAF8A7L2SZ869UKDVJMFD0DLS2256RN6VTG3DLQK2HV3RQQH7U0W7MZ88AU6DXWDQWVP0TX46K0P0CG30GLUCV60SVRZ0A9RMMC4JYH")
	neuterKey, _ := NewKeyFromString("xpublic1pqxqgpqyqprm5hcpeaf8a7l2sz869ukdvjmfd0dls2256rn6vtg3dlqk2hv3rpvsr27rs8y2xea79g56w0phcpw4vaetqjdy827e7ly8ksz9mflqynsewnkqjdjk6s6shxc6qemrhlypdv6shd3wk87h8fqrt44vtvc9th90e8cjfe7remx5ulr02pp3sk6fscazk8uw4xw8ehu4tp5g5dgqm5xfh7")
	h := HardenedKeyStart

	// Case 1
	extKey1, _ := extKey.Derive(1)
	prv1, _ := extKey1.BLSPrivateKey()
	assert.Equal(t, prv1.String(), "SECRET1P22Z2NPSWG09Z5K4LY72RD58JXA330MNHMKJ0ETKGS4VWDSGF32SQSJKTSA")
	assert.Equal(t, extKey1.Path(), Path{0 + h, 1})

	neuterKey1, _ := neuterKey.Derive(1)
	_, err := neuterKey1.BLSPrivateKey()
	assert.ErrorIs(t, err, ErrNotPrivExtKey)
	assert.Equal(t, neuterKey1.BLSPublicKey().String(), "public1ps5hvvw75jtpga5m72vspywjdav7jf6548nm63qxdh8f8kght6qdfdgma5lnr4ppwu66880c9d94ec9j2z324u23qxdqywtkp6ztjyuvw8vsqyff3rjqy9lpsd8m7c65w4l3jjrvwc3fj3gqsxe7w2f8x5ynv3ju2")
	assert.Equal(t, neuterKey1.Path(), Path{0 + h, 1})

	// Case 2
	extKey2, _ := extKey1.Derive(2 + h)
	prv2, _ := extKey2.BLSPrivateKey()
	assert.Equal(t, prv2.String(), "SECRET1PRSRDLLE030ATWX4GAEY0U000JX4QTX42PZQNMPEV26NECY2N2DYSTC23RR")
	assert.Equal(t, extKey2.Path(), Path{0 + h, 1, 2 + h})

	_, err = neuterKey1.Derive(2 + h)
	assert.ErrorIs(t, err, ErrDeriveHardFromPublic)
}

// TestInvalidString checks errors corresponding to the invalid strings
func TestInvalidString(t *testing.T) {
	tests := []struct {
		desc          string
		str           string
		expectedError error
	}{
		{
			desc:          "invalid checksum",
			str:           "XSECRET1PQ2QGPQYQPQQLSPCJLQWPR645ALQJ7F6297W6CDEJRSWW5F2MVV4DMY903MVCJS6JSJ5CVRJREG4940E8JSMDPU3HVVT7UA7A5N72AJY9TRNVZZV25QHE8RC0",
			expectedError: bech32m.ErrInvalidChecksum{Expected: "he8rcz", Actual: "he8rc0"},
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
func TestNeuter(t *testing.T) {
	extKey, _ := NewKeyFromString("XSECRET1PQ2QGPQYQPQQLSPCJLQWPR645ALQJ7F6297W6CDEJRSWW5F2MVV4DMY903MVCJS6JSJ5CVRJREG4940E8JSMDPU3HVVT7UA7A5N72AJY9TRNVZZV25QHE8RCZ")
	neuterKey := extKey.Neuter()
	assert.Equal(t, neuterKey.String(), "xpublic1pq2qgpqyqpqqlspcjlqwpr645alqj7f6297w6cdejrsww5f2mvv4dmy903mvcjsu99mrrh4yjc28dxljnyqfr5n0t85jw49fu775gpnde6faj967sr2t2xld8ucaggthxk3em7ptfdwwpvjs5240z5gpngprjaswsju38rr3myqpz2vguspp0cvrf7lkx4r40uv5smrky2v52qypk0njjfe4pqjrgak")
	assert.Equal(t, neuterKey, neuterKey.Neuter())
}
