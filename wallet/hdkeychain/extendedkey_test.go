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
		// Test vector 1
		{
			name:     "test vector 1 chain m",
			path:     []uint32{},
			wantPriv: "SECRET1PPHL58434CM687WHZ0CDL2A36KDFRQ9GZS8SXJHX90E25W6DH6V3Q8SJRYL",
			wantPub:  "public1pjhtucuyxzsmy526l5yjffj9ddt0lj05rhgk4rt7hyxrxhj27sdtc6qd8gxkseexlpfnjyd9jp2r7zqsy87mh72dtl9s2egmftkxgsc63rfj2yntedhz9xvx5smc8sqvhswzdv4a5zzmknzre4qz4cw2urupv25wf",
		},
		{
			name:     "test vector 1 chain m/0",
			path:     []uint32{0},
			wantPriv: "SECRET1PR83YEN7X5NQMUTQFZJYDXKELCUL4R48GCTC3W5L9JDHEDU2HVCWS5JA5H7",
			wantPub:  "public1pny43ff27stg5uvj8n6wd59aw54su9ngddg5lyn8g280kcyxlp3vjq7w4a5tz75ekmf8pslgkvrnryyryfwld7a7eq6rxrenw06gf5a69xkg8rkmew7lk3e5445dfcnf3v06fg4zcnqju6cum8m7fej47xc3pke58",
		},
		{
			name:     "test vector 1 chain m/0/1",
			path:     []uint32{0, 1},
			wantPriv: "SECRET1PF9W7HJ85Y2LXJ3U2EN86TFY6YPA0KPV28F7UKCTHMT87J5C9F86Q9826VN",
			wantPub:  "public1p5e9desagtctlq72f5sncjkf3a8wm3g9w6jz6grjq337rwx2zfqfg79kx4df2e292z0457ap0pjegz9mar87acpg3alr2tlr3hh8pjvuyelhhmhzn4qlyyf85p6eljlxf7fvzy0uf60e2nefajyqum0tpjurxvvs4",
		},
		{
			name:     "test vector 1 chain m/0/1/2",
			path:     []uint32{0, 1, 2},
			wantPriv: "SECRET1PDU6VZEEDUW6KLATRTEK6CTEZCVDRHTJKZC4NM5GC08VSDFDR93XS7DJE4Z",
			wantPub:  "public1p46264agws63z8rc9678vkytpjtyq4vhd05qp6f4njdgkg3k2uwrm99m75mmrep56me3acq390c2569lx7krghct8fnyldkcjacplnk4t7zw7h3r2ftadkdt0wyjdzjk6pdkerewqn6snygvsknkfjcajcvzluaxk",
		},
		{
			name:     "test vector 1 chain m/0/1/2/2",
			path:     []uint32{0, 1, 2, 2},
			wantPriv: "SECRET1P88ZLGQ3QQUNKRA32302VC7ZPQYXYXRZAVY7HQCWL53YLCCG52EZS29HPCK",
			wantPub:  "public1pnx2jce5y68sl7rxjndquy5w324dklc6a2f5rs73l5slwcfhp3zklr9aywnt3f4zyl7lgllpcv7e0qq7wlhrkanlxxgt3xhprf6mfke63lqx45mrznkl7yakd46hkqce0vx27afng93x7ddvl3jq82yv9cuadqcgl",
		},
		{
			name:     "test vector 1 chain m/0/1/2/2/1000000000",
			path:     []uint32{0, 1, 2, 2, 1000000000},
			wantPriv: "SECRET1P87P43XJQVD4ANFD9QFDD0J34A8ZSSTFYJKMKEW53C8NQEMLTVH0SRY2G0P",
			wantPub:  "public1pswx2u5xqxw3p223wxpcql2p5082kvdc879n7f9f26xwq2eqxc2v8pgwazazzqk09lpcrw2r030dxkyswc3ggmulkjzlzxaaaglc8yex822q39d5tfjerz5769nqlwp0qwaks2qgnv6g8dft55f0t9w38uscadmsk",
		},
	}

	extKey, _ := NewMaster(testSeed)
	extKey2, _ := NewMaster(testSeed)
	neuterKey, _ := extKey2.Neuter()
	for i, test := range tests {
		for _, childNum := range test.path {
			var err error
			extKey, err = extKey.Derive(childNum)
			require.NoError(t, err)

			neuterKey, err = neuterKey.Derive(childNum)
			require.NoError(t, err)
		}

		privStr, err := extKey.BLSPrivateKey()
		require.NoError(t, err)
		require.Equal(t, privStr.String(), test.wantPriv,
			"mismatched serialized private key for test #%v", i+1)

		pubStr, err := extKey.BLSPublicKey()
		require.NoError(t, err)
		require.Equal(t, pubStr.String(), test.wantPub,
			"mismatched serialized public key for test #%v", i+1)

		require.True(t, extKey.IsPrivate())
		require.False(t, neuterKey.IsPrivate())
		neuterPubStr, _ := extKey.BLSPublicKey()
		require.Equal(t, neuterPubStr, pubStr)
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
		wantPub  string
	}{
		// Test vector 1
		{
			name:     "test vector 1 chain m",
			path:     []uint32{},
			wantPriv: "SECRET1PPHL58434CM687WHZ0CDL2A36KDFRQ9GZS8SXJHX90E25W6DH6V3Q8SJRYL",
			wantPub:  "public1pjhtucuyxzsmy526l5yjffj9ddt0lj05rhgk4rt7hyxrxhj27sdtc6qd8gxkseexlpfnjyd9jp2r7zqsy87mh72dtl9s2egmftkxgsc63rfj2yntedhz9xvx5smc8sqvhswzdv4a5zzmknzre4qz4cw2urupv25wf",
		},
		{
			name:     "test vector 1 chain m/0H",
			path:     []uint32{hkStart},
			wantPriv: "SECRET1PGQ2F6J2DF25W2G5C7UKJKC9FKADFQT6R4DL7PJ4YLTFHVMEQFCCQCAMQVK",
			wantPub:  "public1p3urj0ue8wty7yqsv6u8xl5yl2muznywdzxj8sfyga2e48r7hl6kvd28l35um8y23jn4vh6cl24mujz80clm39w57gu528lzk33wpdwpk4ftfgy5vyr5vwxxe506qepsdtv2refd253wwc80m5zrqjv0fyvrhvh7f",
		},
		{
			name:     "test vector 1 chain m/0H/1",
			path:     []uint32{hkStart, 1},
			wantPriv: "SECRET1P9SLS00AZEVKQPYAEWJQG0X3KGA2WAEV83APCKUTMSDZFJ9MSZZUQ07VRGU",
			wantPub:  "public1pjgdrhs9dtqsa8v0u8aqg7g3xmwnduxfe3cwfaqkmu7jygqeduuae2u2l9say5xszga2z3zvvuul0qxprwze5qlzf3ec5de5707e0sth94dn409t8y7rw66vjefv6aqnf5my5l0jjzzngqlw5yqc964j25y9w3fws",
		},
		{
			name:     "test vector 1 chain m/0H/1/2H",
			path:     []uint32{hkStart, 1, hkStart + 2},
			wantPriv: "SECRET1PVS595W2ATQT8S3PTRNU5JSXWU0DKMYW8A9LHYK6CUDAUC7CPNJPSDYY75P",
			wantPub:  "public1pspvarul3exq4v2pk54ryrx4pmqwg55avz85jhh63r0mnmtsmpk3r5a70jy0v99hth4j4ppjr9mg56rk4l9pnetpkyts7e6fv0vxmccpxrp9m3vjgwxhau7psrzrqv27klpgluqs9smztllhvmhykd39g7sus2yqy",
		},
		{
			name:     "test vector 1 chain m/0H/1/2H/2",
			path:     []uint32{hkStart, 1, hkStart + 2, 2},
			wantPriv: "SECRET1PTGS4GXEDKNG82RQFTF5CTVWTRJH5RCNDV60JEWWW74EQR87QYNES2PRFNP",
			wantPub:  "public1p3hn0zrackxprg00h9sgylw40j5a0qw45kv5z905rfn488f2wsq8r8yerargse3sy5uwkskz8t5wkvz23wxeescxny70nzp5g86ku9kn576gv8dvphz0gvm0puup9tzqce86tw740vchu6glxta46vdxz8ctz6k96",
		},
		{
			name:     "test vector 1 chain m/0H/1/2H/2/1000000000",
			path:     []uint32{hkStart, 1, hkStart + 2, 2, 1000000000},
			wantPriv: "SECRET1PTTX080V6FMFWCECSCJNX366S8MG6USLH0UFCQVPUHSLW0WXLDJQQ9E2KYK",
			wantPub:  "public1p3qsmjpzf3w3sfqs2wscxkqjzp7l3pnl7d7yhk032s52d039rw63tye55d5s3kchslefnh0jm09g2gq293ywnseqpdskzsqy8xslp9ulu9pqeh60h02594rl9g4uzl7qfjf4gdcmcl5r6nmvl2w5w5a6ur5x9ellk",
		},
	}

	extKey, _ := NewMaster(testSeed)
	for i, test := range tests {
		for _, childNum := range test.path {
			var err error
			extKey, err = extKey.Derive(childNum)
			require.NoError(t, err)
		}

		privStr, err := extKey.BLSPrivateKey()
		require.NoError(t, err)
		require.Equal(t, privStr.String(), test.wantPriv,
			"mismatched serialized private key for test #%v", i+1)

		pubStr, err := extKey.BLSPublicKey()
		require.NoError(t, err)
		require.Equal(t, pubStr.String(), test.wantPub,
			"mismatched serialized public key for test #%v", i+1)

		require.True(t, extKey.IsPrivate())
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
