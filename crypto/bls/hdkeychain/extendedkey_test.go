package hdkeychain

import (
	"encoding/hex"
	"io"
	"testing"

	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNonHardenedDerivation tests derive private key and public key in
// non hardened mode.
func TestNonHardenedDerivation(t *testing.T) {
	testSeed, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	tests := []struct {
		name       string
		path       Path
		wantPrivG1 string
		wantPrivG2 string
		wantPubG1  string
		wantPubG2  string
	}{
		{
			name:       "derivation path: m",
			path:       Path{},
			wantPrivG1: "38167a7b6fcea7929deb6af40123e37b2ef21e488f4871f16d411914490657f5",
			wantPrivG2: "38167a7b6fcea7929deb6af40123e37b2ef21e488f4871f16d411914490657f5",
			wantPubG1:  "b99b512321d9dbae271f4d418b10a2345fa84c1c883d0f9a82163b84c405948ea123f01141258cdbed2d033eae4a551a",
			wantPubG2:  "97167b36223b32c3b0b8c8234eefeeee12a37b3aafa9ed4e156facd2d5a4206d72f814bd69270ae0a95d6de13f1a1e6618330e1d02506ed8ec2c184660912213bf2cd54f66b38fe3e5dd241320cf2684f6c86affd93a07689870dedeb6cf133b",
		},
		{
			name:       "derivation path: m/0",
			path:       Path{0},
			wantPrivG1: "32833bcae69d2296e71356dc444d08c9d78671a6dbc1234fd22ee78ddd21359a",
			wantPrivG2: "3914227ccdfedae713b1289c14a4125de77dc8ae3dcb9c4ede1912819d03f807",
			wantPubG1:  "a33ee09eee91c2d7b345f8ca28daf9fcbc92d2fd9b0be0f582ff3e5e71b9b13999fc4918d6894ba6b5ab6acdc556f665",
			wantPubG2:  "adc3abc6d60ab9cf5dcea68fce523959f0564013536c7328f420fe83c94de483a99230a620d9620bb822c74ee112d4f40dd3ec63b20f731b459d875efb0c9f79f55567ec3aa4a361b8b68134ce11f99fa9b9a15dea38edd610d74d0b4ac9b31a",
		},
		{
			name:       "derivation path: m/0/1",
			path:       Path{0, 1},
			wantPrivG1: "67a6f33ab789e703bc33364591230662b4a067fd217c0bf94781a0438447f7e6",
			wantPrivG2: "1571b4b6c4233eb6132b98682e430919b32ab4e0b18b27edf4e7505a78814c9c",
			wantPubG1:  "82de21510f071ff70f6eba967e6876487a7e63d8fa43e9d37d0e89af8ccad9d177d3fbdaac566fe7a9a19e8288b9df7c",
			wantPubG2:  "95ac09809e30db3a6cab7fee871e3fbc0b670a25ceebc2a1146d72c67c03f03a8766692773609c0898887095b03db71b08470fc21905ee737511bbb976a1d77a9fc523c98aea3439fbf22dc0e8bff9878054116a95d1607bf9a2225a8f96c331",
		},
		{
			name:       "derivation path: m/0/1/2",
			path:       Path{0, 1, 2},
			wantPrivG1: "083e105a2722e56a473c8ce703d264789e9f4642e0163d6d7f2e91133a4f6c14",
			wantPrivG2: "0ab855d785a73ac1d434993238eac9465398848c940c569304a6cf8104432323",
			wantPubG1:  "a931438adb1639201f1d5ea7445a5cfeb319c0143f83964bd2650ea7b835ec9360e55eb5f48345f8cff6297bde663773",
			wantPubG2:  "a04d61ace6bb9f147feaa4f655bf61d3915de1274ae3c455d1cd61260b4c4fb154b53b86cbf06d7f035e8f3046f46c8a0ee825aeb0d871b34993964f43e587dde51f26609b8bc3c58db93ec975d7b21f947d5ef13fe611a6c624584604bd9745",
		},
		{
			name:       "derivation path: m/0/1/2/2",
			path:       Path{0, 1, 2, 2},
			wantPrivG1: "45dc329488337cce410446a6735abb1f2885ef3edcc717e4a1541ff5e0a302e5",
			wantPrivG2: "1af4a741d566a744bad2ea4d632138af1084b1a85c9deb486920421ceacee200",
			wantPubG1:  "8dca90a2ef835d22b6a1a31a0f35a42afa9f995d7ff5050391d63246e9d870f9da7d30201006d60b2b2520edd49123a5",
			wantPubG2:  "95994434c768d7273ba943453b60912b9094e267e4611816d2f83905931c433ba9db841bd1aef53ac469cc759718f4641985be8d5b1e99d3b216daf61660e37f90ef77c61fa27c3b462199ab04a8a5091a15f85f5f166ee281ea8415b357f07d",
		},
		{
			name:       "derivation path: m/0/1/2/2/1000000000",
			path:       Path{0, 1, 2, 2, 1000000000},
			wantPrivG1: "3255d4c3b4839073d8bd8e58a49cb35579b23c3ab393844efb6f382901642ffe",
			wantPrivG2: "59fef95b26def9f7c068d821025f87b7b7d05b5a4a983ccdff3c375faa1cdaf5",
			wantPubG1:  "9660f7ef3d08328dd166063a47ab6a9858b56833ac5d1f991195e8a75698b07a45b8cb3e620835062862693f80e1b012",
			wantPubG2:  "a96976500c782dd268fd7703813b7710f5896ade7a7823518fec005b8bc5c39c3c00df82073658fabb47339ce83435500b5f64b3159e0c0acfe48852bc42bb09e83795783f69b1f25a41f02a431aad50b7d0eac09f33f672db95a092af186754",
		},
	}

	masterKeyG1, _ := NewMaster(testSeed, true)
	masterKeyG2, _ := NewMaster(testSeed, false)
	neuteredMasterKeyG1 := masterKeyG1.Neuter()
	neuteredMasterKeyG2 := masterKeyG2.Neuter()
	for i, test := range tests {
		extKeyG1, err := masterKeyG1.DerivePath(test.path)
		require.NoError(t, err)

		extKeyG2, err := masterKeyG2.DerivePath(test.path)
		require.NoError(t, err)

		neuterKeyG1, err := neuteredMasterKeyG1.DerivePath(test.path)
		require.NoError(t, err)

		neuterKeyG2, err := neuteredMasterKeyG2.DerivePath(test.path)
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
	}
}

// TestHardenedDerivation tests derive private key and public key in
// hardened mode.
func TestHardenedDerivationG2(t *testing.T) {
	testSeed, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	h := HardenedKeyStart
	tests := []struct {
		name      string
		path      Path
		wantPriv  string
		wantPubG1 string
		wantPubG2 string
	}{
		{
			name:      "derivation path: m/0H",
			path:      Path{0 + h},
			wantPriv:  "0da9a671ff3cc10514fa9cf8368af0756ad01e242bb3ae4d06a4061e6c3ec6be",
			wantPubG1: "ad880d92cb1276c8d5768d26dcc3ebf23de4f0a1d12c9b99fdb6b745ff931d2c01925b4fec4147d7452faddc18564b98",
			wantPubG2: "b5bfa5aa0ee3c8c4f9222f4f244c3fd6fddcdd0accf484e56d09c70a4e2eada059073c890bbfc30d2693a8a68262c3ac06da94222f5d9de2b8ea4fb7cf29664196becd7a36f49a55f9364afe2536de909b87a8c2900cd38590f8dd4b51c8642a",
		},
		{
			name:      "derivation path: m/0H/1H",
			path:      Path{0 + h, 1 + h},
			wantPriv:  "319ebd81870da6bc3608b80fa22418c3dbb854f8719c19be687533266cb0f291",
			wantPubG1: "95dae848c45f54076c720be4d1c9f543993ccdb8850d47b9d21450b0e861de12473f5a68193d2f60e1d44ef37d57dd32",
			wantPubG2: "94c2ef2238d79c7148719a9faade73c7c27f33ce6526b3d374950f76f074d1ef3a4798f8d4dbd97b6c800640421bbb4d0c39943d28c1bcf62acd85f957f29fd293d375c10bc3642e66f282d129fbddf81c9ddc0aba453f3b6e0f404b0b4a12a5",
		},
		{
			name:      "derivation path: m/0H/1H/2H",
			path:      Path{0 + h, 1 + h, 2 + h},
			wantPriv:  "028a08fe12f776e2538749510644c9ad78e020c36b9cee216518ba6d52fd7549",
			wantPubG1: "829d3fa8890ddad5289c9429af9a3a6fa284ecdc3ca6a91dc7af466b1edd3e03952869def6baec191bdbfd76de79d754",
			wantPubG2: "923c884ce3f91e9cc4595114c5fbe3aee03b8698c7d35ae358f76816d132c6972c6612ad18c2178cf1c9e47d8e0e5822164c8bee956ddccf7103d4a41c3fcf092f746493c333cd33e95efa60b62bb5ef6beaecd3b3c338bc4ed56d6000e59748",
		},
		{
			name:      "derivation path: m/0H/1H/2H/2H",
			path:      Path{0 + h, 1 + h, 2 + h, 2 + h},
			wantPriv:  "67a7b55c4cf0620cd99e55c77b8075f3ec5ce2ac9f840678dada4047849d4e91",
			wantPubG1: "91a562e2f28d775560426db45dda4b9b12cd9c55bff4501c8a1253ffee26e2299f18c13c81d586b40955c120a7bd0ea4",
			wantPubG2: "a800a88b89dcbf614e498e4dacfcaccb9735fad06b5ca722e75debf4b22a0ffdfb761898f880de07220d584e3ea4e28e13b7ab01b56e0883587ee73ce0bdcfbbc2f2948ce3ed14aa88d69cbaefab38bf3f374ba1be72001a44bd183184053ee6",
		},
		{
			name:      "derivation path: m/0H/1H/2H/2H/1000000000H",
			path:      Path{0 + h, 1 + h, 2 + h, 2 + h, 1000000000 + h},
			wantPriv:  "5fa9eb585d673c081b07d56639a980b69d7a1f31f4725d23884aaeb2bf0e6fdc",
			wantPubG1: "949a2b9d20bd8c260f1c9f0b4631dc3abff6bcf3b723652e65c06447b1f0b7ca1f2f9b5fecd3692222d1e7280e8c0a2d",
			wantPubG2: "952ca984734510e531507cdb653a3bce2ed4b9689d917718663ab04c84bb78596d7f46252a1ea0e3dd7a380a5a5bf8f208d9571004508dbc01464e95712a2f5df3be43c1f09cd06f0c0944a4eaa5d5a1fcde5cd2ec2661c6ac3307b78b694cea",
		},
	}

	masterKeyG1, _ := NewMaster(testSeed, true)
	masterKeyG2, _ := NewMaster(testSeed, false)
	neuteredMasterKeyG1 := masterKeyG1.Neuter()
	neuteredMasterKeyG2 := masterKeyG2.Neuter()
	for i, test := range tests {
		extKeyG1, err := masterKeyG1.DerivePath(test.path)
		require.NoError(t, err)

		extKeyG2, err := masterKeyG2.DerivePath(test.path)
		require.NoError(t, err)

		_, err = neuteredMasterKeyG1.DerivePath(test.path)
		require.ErrorIs(t, err, ErrDeriveHardFromPublic)

		_, err = neuteredMasterKeyG2.DerivePath(test.path)
		require.ErrorIs(t, err, ErrDeriveHardFromPublic)

		privKeyG1, err := extKeyG1.RawPrivateKey()
		require.NoError(t, err)
		require.Equal(t, hex.EncodeToString(privKeyG1), test.wantPriv,
			"mismatched serialized private key for test #%v", i+1)

		privKeyG2, err := extKeyG2.RawPrivateKey()
		require.NoError(t, err)
		require.Equal(t, privKeyG1, privKeyG2)

		pubKeyG1 := extKeyG1.RawPublicKey()
		require.Equal(t, hex.EncodeToString(pubKeyG1), test.wantPubG1,
			"mismatched serialized public key for test #%v", i+1)

		pubKeyG2 := extKeyG2.RawPublicKey()
		require.Equal(t, hex.EncodeToString(pubKeyG2), test.wantPubG2,
			"mismatched serialized public key for test #%v", i+1)

		require.True(t, extKeyG1.IsPrivate())
		require.True(t, extKeyG2.IsPrivate())
		require.Equal(t, extKeyG1.Path(), test.path)
		require.Equal(t, extKeyG2.Path(), test.path)
	}
}

// TestDerivation tests derive private keys in hardened and non hardened modes.
func TestDerivation(t *testing.T) {
	testSeed, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	h := HardenedKeyStart
	tests := []struct {
		name       string
		path       Path
		wantPrivG1 string
		wantPrivG2 string
		wantPubG1  string
		wantPubG2  string
	}{
		{
			name:       "derivation path: m",
			path:       Path{},
			wantPrivG1: "38167a7b6fcea7929deb6af40123e37b2ef21e488f4871f16d411914490657f5",
			wantPrivG2: "38167a7b6fcea7929deb6af40123e37b2ef21e488f4871f16d411914490657f5",
			wantPubG1:  "b99b512321d9dbae271f4d418b10a2345fa84c1c883d0f9a82163b84c405948ea123f01141258cdbed2d033eae4a551a",
			wantPubG2:  "97167b36223b32c3b0b8c8234eefeeee12a37b3aafa9ed4e156facd2d5a4206d72f814bd69270ae0a95d6de13f1a1e6618330e1d02506ed8ec2c184660912213bf2cd54f66b38fe3e5dd241320cf2684f6c86affd93a07689870dedeb6cf133b",
		},
		{
			name:       "derivation path: m/0H",
			path:       Path{h},
			wantPrivG1: "0da9a671ff3cc10514fa9cf8368af0756ad01e242bb3ae4d06a4061e6c3ec6be",
			wantPrivG2: "0da9a671ff3cc10514fa9cf8368af0756ad01e242bb3ae4d06a4061e6c3ec6be",
			wantPubG1:  "ad880d92cb1276c8d5768d26dcc3ebf23de4f0a1d12c9b99fdb6b745ff931d2c01925b4fec4147d7452faddc18564b98",
			wantPubG2:  "b5bfa5aa0ee3c8c4f9222f4f244c3fd6fddcdd0accf484e56d09c70a4e2eada059073c890bbfc30d2693a8a68262c3ac06da94222f5d9de2b8ea4fb7cf29664196becd7a36f49a55f9364afe2536de909b87a8c2900cd38590f8dd4b51c8642a",
		},
		{
			name:       "derivation path: m/0H/1",
			path:       Path{h, 1},
			wantPrivG1: "0f356ce1bc5fd0d8c6ac274c7b206c2cb868b502af7a7065e836ca28f0d865a5",
			wantPrivG2: "647e9be74b8648a198267e18c40d988b74f1ff8a33132fb012187ffd22c96fc4",
			wantPubG1:  "92c9d3908b636a383c174bd3d27a8ae5b116240b90aa4d4b2425405b11eded72419284f04c66438578adffbc62e3698d",
			wantPubG2:  "ae08105f86d96658083e0799bb2fa282f6431a7d23044ef6aaf3083b027aaed32798005e27668a7327e36a7db07fcab40253921fe842ba2b7411a44bfd292c7cf124970b70da554fceff6eef4ba2edddeb713d40bebbd79c24b9888541da63a3",
		},
		{
			name:       "derivation path: m/0H/1/2H",
			path:       Path{h, 1, 2 + h},
			wantPrivG1: "1f8dd29cf00257ffd4bc4001580d2011ce23a0727ce20ee4156c7a6fb61964ce",
			wantPrivG2: "578d64b2246b6b8c8906073b2de5cad50313997206ba872d38886d3d2d5eb508",
			wantPubG1:  "93981a7ebf6ec8ffe6a5be1a4dfda93c8e9b192f79417ba5d4c4590c034306666ae3408146c33325ed8bab3fc0b01926",
			wantPubG2:  "a89393efab309f3be02e501ebb6f472ef6392d48f832289247992497f4da816dfefc89fccf8b423b6e3005f7f29cdbf305c4aa0bd6163ece7d816dd34d83193512f17a752fb1a64b7e6455b7e301dcddf4eab0170304a1398d53abaabc63bdce",
		},
		{
			name:       "derivation path: m/0H/1/2H/2",
			path:       Path{h, 1, 2 + h, 2},
			wantPrivG1: "1ff1c89c8ac0387cb6b66d79c3da6b4bdc400d991d71ae78aa569e8d005d84ae",
			wantPrivG2: "22d2d08809f4fe0103b1b4750565da52147d8694c6c382320a988368ea0fe14c",
			wantPubG1:  "8eec88f3d3dd95b6d234ca83a9862989c3b960196a482d7c81eafd207faf60acac5624101184394a9dd4a7848115320e",
			wantPubG2:  "aabdc2a56f89fce40ee071dc3193aa63784391930af044258d3e104bdc9b2bef9aacb310466bada59df534c46cb6cfa2073e2c7c1e92905cba9d55bc32316ccfc2b71ac00fa75ed26098f301475a3dcf8fa5793b49aa6f538dd74873a5ca8baa",
		},
		{
			name:       "derivation path: m/0H/1/2H/2/1000000000",
			path:       Path{h, 1, 2 + h, 2, 1000000000},
			wantPrivG1: "3475af3ba6a3dbb5c181d9be4eaf0bea3ba6d1493a7a028978b563acd0799a67",
			wantPrivG2: "4a19cfe432a510604ca9e6cdfda424e021f477ebbdd7e4130c01d32d567a92ee",
			wantPubG1:  "860c0f97c1e8ed8d8c0265898e872428a4d6f33a10700e26a2536a333a6a5579029e9dd870c3801867667dba6b750fff",
			wantPubG2:  "b12f350aaa1d66451a8a6fedc1bbcd3d3895808b7788d09e9c5d14d91a05e4c983e0fdbffa3a2afa393528b9cbbb49090880ef251db20a6d671f00a82afa8cba5ba5964a1cf96c8592751dc65c38c268d2bf1bcdfd88aab585c21a96dfa4ba6f",
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

		require.True(t, extKeyG1.IsPrivate())
		require.True(t, extKeyG2.IsPrivate())

		neuterKeyG1 := extKeyG1.Neuter()
		neuterKeyG2 := extKeyG2.Neuter()
		neuterPubKeyG1 := neuterKeyG1.RawPublicKey()
		neuterPubKeyG2 := neuterKeyG2.RawPublicKey()

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
	}
}

// TestInvalidDerivation tests Derive function for invalid data
func TestInvalidDerivation(t *testing.T) {
	t.Run("Private key is 31 bytes. It should be 32 bytes", func(t *testing.T) {
		key := [31]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], Path{}, true, false)
		_, err := ext.Derive(HardenedKeyStart)
		assert.ErrorIs(t, err, ErrInvalidKeyData)
	})

	t.Run("Public key on G1 is 96 bytes. It should be 48 bytes", func(t *testing.T) {
		key := [96]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], Path{}, false, true)
		_, err := ext.Derive(0)
		assert.ErrorIs(t, err, ErrInvalidKeyData)
	})

	t.Run("Public key on G2 is 42 bytes. It should be 96 bytes", func(t *testing.T) {
		key := [95]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], Path{}, false, false)
		_, err := ext.Derive(0)
		assert.ErrorIs(t, err, ErrInvalidKeyData)
	})

	t.Run("Invalid key", func(t *testing.T) {
		key := [95]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], Path{}, false, false)
		_, err := ext.Derive(0)
		assert.ErrorIs(t, err, ErrInvalidKeyData)
	})

	t.Run("Derive public key from hardened key", func(t *testing.T) {
		key := [32]byte{0}
		chainCode := [32]byte{0}
		ext := newExtendedKey(key[:], chainCode[:], Path{}, false, false)
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
	testSeed, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f")
	h := HardenedKeyStart
	tests := []struct {
		name        string
		path        Path
		wantXPrivG1 string
		wantXPrivG2 string
		wantXPubG1  string
		wantXPubG2  string
	}{
		{
			name:        "derivation path: m",
			path:        Path{},
			wantXPrivG1: "XSECRET1PQP3PH9J2P809LAD5S8G2T5MCD32Z7KS9QA2LEY8WRR2ASUDKT4R5JVPCZEA8KM7W57FFM6M27SQJ8CMM9MEPUJY0FPCLZM2PRY2YJPJH753GQ87P",
			wantXPrivG2: "XSECRET1PQP3PH9J2P809LAD5S8G2T5MCD32Z7KS9QA2LEY8WRR2ASUDKT4R5JCPCZEA8KM7W57FFM6M27SQJ8CMM9MEPUJY0FPCLZM2PRY2YJPJH752EG3YW",
			wantXPubG1:  "xpublic1pqp3ph9j2p809lad5s8g2t5mcd32z7ks9qa2ley8wrr2asudkt4r5jv9endgjxgwemwhzw86dgx93pg35t75yc8yg858e4qsk8wzvgpv536sj8uq3gyjcekld95pnatj225dqkj87e7",
			wantXPubG2:  "xpublic1pqp3ph9j2p809lad5s8g2t5mcd32z7ks9qa2ley8wrr2asudkt4r5jcyhzeanvg3mxtpmpwxgyd8wlmhwz23hkw4048k5u9t04nfdtfpqd4e0s99adyns4c9ft4k7z0c6renpsvcwr5p9qmkcaskps3nqjy3p80ev648kdvu0u0ja6fqnyr8jdp8kep40lkf6qa5fsux7m6mv7yem46g8ve",
		},
		{
			name:        "derivation path: m/0H",
			path:        Path{h},
			wantXPrivG1: "XSECRET1PQXQGPQYQPZLR8Y9H6D4SCKKR707WUTDW22J2K5K0X77J45PLNFU8UTS3DDSJ2VQD4XN8RLEUCYZ3F75ULQMG4UR4DTGPUFPTKWHY6P4YQC0XC0KXHCC8M7LE",
			wantXPrivG2: "XSECRET1PQXQGPQYQPZLR8Y9H6D4SCKKR707WUTDW22J2K5K0X77J45PLNFU8UTS3DDSJ2CQD4XN8RLEUCYZ3F75ULQMG4UR4DTGPUFPTKWHY6P4YQC0XC0KXHCRKNG9K",
			wantXPubG1:  "xpublic1pqxqgpqyqpzlr8y9h6d4sckkr707wutdw22j2k5k0x77j45plnfu8uts3ddsj2v9d3qxe9jcjwmyd2a5dymwv86lj8hj0pgw39jdenldkkazllyca9sqeyk60a3q5046997kacxzkfwvq90zc4g",
			wantXPubG2:  "xpublic1pqxqgpqyqpzlr8y9h6d4sckkr707wutdw22j2k5k0x77j45plnfu8uts3ddsj2c94h7j65rhrerz0jg30fujyc07klhwd6zkv7jzw2mgfcu9yut4d5pvsw0yfpwluxrfxjw52dqnzcwkqdk55ygh4m80zhr4yld7099nyr947e4arday62hunvjh7y5mdayyms75v9yqv6wzep7xafdgusep2my7jwj",
		},
		{
			name:        "derivation path: m/0H/1",
			path:        Path{h, 1},
			wantXPrivG1: "XSECRET1PQ2QGPQYQPQQ65UA6MN3X4QULRNLCTZY30LX8U38KJ0NX3W88K8RVSP3KS06NA6PSPU6KECDUTLGD334VYAX8KGRV9JUX3DGZ4AA8QE0GXM9Z3UXCVKJSU8APV0",
			wantXPrivG2: "XSECRET1PQ2QGPQYQPQQJCR83YKMDJ39MCEN24DDMK2ESK935LZT6J4WDUN47GXWDPCGUEVNQV3LFHE6TSEY2RXPX0CVVGRVC3D60RLU2XVFJLVQJRPLL6GKFDLZQT4UW7A",
			wantXPubG1:  "xpublic1pq2qgpqyqpqq65ua6mn3x4qulrnlctzy30lx8u38kj0nx3w88k8rvsp3ks06na6psjtya8yytvd4rs0qhf0fay752ukc3vfqtjz4y6jeyy4q9ky0da4eyry5y7pxxvsu90zkll0rzud5c6p4824p",
			wantXPubG2:  "xpublic1pq2qgpqyqpqqjcr83ykmdj39mcen24ddmk2esk935lzt6j4wdun47gxwdpcguevnq4cypqhuxm9n9szp7q7vmktazstmyxxnayvzyaa427vyrkqn64mfj0xqqtcnkdznnyl3k5lds0l9tgqjnjg07ss469d6prfztl55jcl83yjtskux6248ualmwaa969mwaadcn6s97h0tecf9e3zz5rknr5vu6qds4",
		},
		{
			name:        "derivation path: m/0H/1/2H",
			path:        Path{h, 1, 2 + h},
			wantXPrivG1: "XSECRET1PQWQGPQYQPQQC9QYQSQY2WE42YAR2XHFJKE9PG5734DUFWY84SPPNQU28A2NXTEW7RQX3KFESR7XA988SQFTLL49UGQQ4SRFQZ88Z8GRJ0N3QAEQ4D3AXLDSEVN8QSYHQXG",
			wantXPrivG2: "XSECRET1PQWQGPQYQPQQC9QYQSQYXKJF465QLR30NMALS94SHVQ53KQ76J3DRVRURLZM7ZHYM8W279YMQ27XKFV3YDD4CEZGXQUAJMEW265P38XTJQ6AGWTFC3PKN6T27K5YQFFDQ0M",
			wantXPubG1:  "xpublic1pqwqgpqyqpqqc9qyqsqy2we42yar2xhfjke9pg5734dufwy84sppnqu28a2nxtew7rqx3kfesjwvp5l4ldmy0le49hcdymldf8j8fkxf009qhhfw5c3vscq6rqenx4c6qs9rvxve9ak96k07qkqvjvvplzpf",
			wantXPubG2:  "xpublic1pqwqgpqyqpqqc9qyqsqyxkjf465qlr30nmals94shvq53kq76j3drvrurlzm7zhym8w279ymq4zfe8matxz0nhcpw2q0tkm689mmrjt2glqez3yj8nyjf0ax6s9klalyfln8cks3mdccqtaljnndlxpwy4g9av937ee7czmwnfkp3jdgj79a82ta35e9huez4kl3srhxa7n4tq9crqjsnnr2n4w4tccaaecqrq7q9",
		},
		{
			name:        "derivation path: m/0H/1/2H/2",
			path:        Path{h, 1, 2 + h, 2},
			wantXPrivG1: "XSECRET1PQJQGPQYQPQQC9QYQSQYQ9P9NU8EW0SLNHN9YY0KWD9CR4VR30CLFD3Y085WAXSREUC94J8JAXQ0LRJYU3TQRSL9KKEKHNS76DD9ACSQDNYWHRTNC4FTFARGQTKZ2UX2Z56E",
			wantXPrivG2: "XSECRET1PQJQGPQYQPQQC9QYQSQYQ9LLEUG0UDJJ8WZVAXSA3HW9VFPM3D2DJJXNK9S0K6R74CF9ZTE7RVQ3D95YGP860UQGRKX682PT9MFFPGLVXJNRV8Q3JP2VGX682PLS5CA364WS",
			wantXPubG1:  "xpublic1pqjqgpqyqpqqc9qyqsqyq9p9nu8ew0slnhn9yy0kwd9cr4vr30clfd3y085waxsreuc94j8jaxz8wez8n60wetdkjxn9g82vx9xyu8wtqr94ysttus8406grl4as2etzkysgprppef2waffuysy2nyrslzejaa",
			wantXPubG2:  "xpublic1pqjqgpqyqpqqc9qyqsqyq9lleug0udjj8wzvaxsa3hw9vfpm3d2djjxnk9s0k6r74cf9zte7rvz4tms49d7yleeqwupcacvvn4f3hssu3jv90q3p935lpqj7unv47lx4vkvgyv6ad5kwl2dxydjmvlgs88ck8c85jjpwt4824hserzmx0c2m34sq05a0dycyc7vq5wk3ae78627fmfx4x75ud6ay88fw23w4q7x3faw",
		},
		{
			name:        "derivation path: m/0H/1/2H/2/1000000000",
			path:        Path{h, 1, 2 + h, 2, 1000000000},
			wantXPrivG1: "XSECRET1PQKQGPQYQPQQC9QYQSQYQ9QY5A0WQ8PVGMPS3CW8PX6JX68GKNJ3FJU9GHDQK9ZY8S4UKMFDVX9MW7HFDXQ68TTEM563AHDWPS8VMUN40P04RHFK3FYA85Q5F0Z6K8TXS0XDXWJ5TMS2",
			wantXPrivG2: "XSECRET1PQKQGPQYQPQQC9QYQSQYQ9QY5A0WQX58Q3Z4FS7DY3Z5QZ82LHG056HEL6RMTG0JTVC5E050DWFCZKDPWVP9PNNLYX2J3QCZV48NVMLDYYNSZRARHAW7A0EQNPSQAXT2K02FWUY38N3L",
			wantXPubG1:  "xpublic1pqkqgpqyqpqqc9qyqsqyq9qy5a0wq8pvgmps3cw8px6jx68gknj3fju9ghdqk9zy8s4ukmfdvx9mw7hfdxzrqcruhc85wmrvvqfjcnr58ys52f4hn8gg8qr3x5ffk5ve6df2hjq57nhv8psuqrpnkvld6dd6sllcvmc8y5",
			wantXPubG2:  "xpublic1pqkqgpqyqpqqc9qyqsqyq9qy5a0wqx58q3z4fs7dy3z5qz82lhg056hel6rmtg0jtvc5e050dwfczkdpwvzcj7dg24gwkv3g63fh7msdme57n39vq3dmc35y7n3w3fkg6qhjvnqlqlkll5w32lgun229eewa5jzggsrhj28djpfkkw8cq4q404r96twjevjsul9kgtyn4rhr9cwxzdrft7x7dlky24dv9cgdfdhayhfhs8yp58j",
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
func TestInvalidString(t *testing.T) {
	tests := []struct {
		desc          string
		str           string
		expectedError error
	}{
		{
			desc:          "invalid checksum",
			str:           "XSECRET1PQP0R4SGK8Y84J2G9LQD2E4W5RYXQRPWKYSG8TT4KUMZD0QF7TT8PSVPNCFWFXK8JKWKNMH8HQC8PV0ZMYL36LRFJJ76K3C94YL38FA7PNGF4LRNP",
			expectedError: bech32m.ErrInvalidChecksum{Expected: "f4lrnq", Actual: "f4lrnp"},
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
	extKey, _ := NewKeyFromString("XSECRET1PQP0R4SGK8Y84J2G9LQD2E4W5RYXQRPWKYSG8TT4KUMZD0QF7TT8PSVPNCFWFXK8JKWKNMH8HQC8PV0ZMYL36LRFJJ76K3C94YL38FA7PNGF4LRNQ")
	neuterKey := extKey.Neuter()
	assert.Equal(t, neuterKey.String(), "xpublic1pqp0r4sgk8y84j2g9lqd2e4w5ryxqrpwkysg8tt4kumzd0qf7tt8psv9yvfyjsl2hv2ya96vam0uwhdq3t753htdv7fwp694njg7ctvnprtrnmddrc9083nrcnvlg8ex9kucs39tdxg")
	assert.Equal(t, neuterKey, neuterKey.Neuter())
}
