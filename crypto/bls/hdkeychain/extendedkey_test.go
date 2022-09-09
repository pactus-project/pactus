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
			wantPrivG1: "33c25c9358f2b3ad3ddcf7060e163c5b27e3af8d3297b568e0b527e274f7c19a",
			wantPrivG2: "33c25c9358f2b3ad3ddcf7060e163c5b27e3af8d3297b568e0b527e274f7c19a",
			wantPubG1:  "a46249287d576289d2e99ddbf8ebb4115fa91badacf25c1d16b3923d85b2611ac73db5a3c15e78cc789b3e83e4c5b731",
			wantPubG2:  "8dfc988fa761d3c26949f3e0b7d21442e416251321192ff0619ba42452152e84699f48199e53c7e435ff912f949e846f1876a4ca81c9e484f5e8ed591696e8ac05d9a3a944a9bb0f47aec5a7dbb5e58ac32cb8b66282bd2b2b51655648625115",
		},
		{
			name:       "derivation path: m/0",
			path:       Path{0},
			wantPrivG1: "1466ddb94f10d14f88f70b772caa9612961ef8e860781ff71cef77cb7ba20711",
			wantPrivG2: "1ed5ab89dc552b15b6eb986094e80b16544215111c5bec0ed85546b458165fb3",
			wantPubG1:  "a09c95fc6f5c47f16bd4eda8728472b3595f92ecfd6676a6b38de3f603a527796241d7818c84444aa0a076c07a1accc5",
			wantPubG2:  "a36bbc51a5e0a81b2916319a5a8d8cc46933f65f222dc333b5b2f7542edf24c778c2a91d144667c9a1d37e57e7d679ba0225c741377e6cce27374339e4062d1a8f66c5469098be5f6ede6778a0f1949d2f80eee0ac01baa0504e895f49060c95",
		},
		{
			name:       "derivation path: m/0/1",
			path:       Path{0, 1},
			wantPrivG1: "2cdde39e222b0591db3e243056637227df517bead054c878b76daf8701f94156",
			wantPrivG2: "31f03abb81e0796c4451215cf7c6b0ee56996220fa1a89b6cc877b5397fc729c",
			wantPubG1:  "8320667eb5b15f92b5847f8e3c9c95005e5594f72fb754fcfdbc9f02ba923eb9fc86ec2287ed844ca362118b7a0b7966",
			wantPubG2:  "8d17b968427674b9928ea0530c5cd76ccb82d24093d8240439128e1e63658b60e500d1d6c0258212a02ef1e02c0039d50cebb8893882c30d8d59643f198bc52279cb2b8baec0856d88fada4ae5feb6417babffc11022bc51284bf4ec0762ae13",
		},
		{
			name:       "derivation path: m/0/1/2",
			path:       Path{0, 1, 2},
			wantPrivG1: "05e173c528d5672d2732e2c2a17c0cd8798f60edb563ababa5ba2b3a813751d3",
			wantPrivG2: "674a1aba003388cb92091c2cf10951e15d280fd0605395d60df5bb8c3564849a",
			wantPubG1:  "af5f8ab633806adefcba56fc9a069eae63ea3e3ed5c79beabc89247d4a4ccd56ebcacdf871cf668d66276301c5d3f6b5",
			wantPubG2:  "b15828396b16ddde2f92994bf15ea950f87a4302db0bde0fcec764d1941adeb0e2c9b1af52b42fc5778c2d2ba1865de30d1aeaa330f60ae1780061f553d7530c463634a5580e4d3afcf642a08d3b12c4e277a75ecc99a934ae72b515d1e40588",
		},
		{
			name:       "derivation path: m/0/1/2/2",
			path:       Path{0, 1, 2, 2},
			wantPrivG1: "1e9c469649399870039e18c833ecbf7c51e8ef719af749fe18bef9a1de855193",
			wantPrivG2: "72e5579288262af7db277e9e2bf26b61bd01fcf9c00d87a93b513ccbe3bc6ec3",
			wantPubG1:  "90a3cc3fbf1b64cec7153e87b4c48adecbfa82df78331bab300b26da27d45d76217bd75bd3e40f5743c4cb082075bba3",
			wantPubG2:  "84dbcfb19b0f31563ad80ae9f937749fabeb700981e21fb037ec240bcab08782a12849773ce6722bee7b05df648f6d3516bbb8cbcf6ea9ade53bcf4c056f64a5a4622906e72b97e14dba3d350383dfeee8c9ab79b4caa08a1448bc53c93a4e5d",
		},
		{
			name:       "derivation path: m/0/1/2/2/1000000000",
			path:       Path{0, 1, 2, 2, 1000000000},
			wantPrivG1: "1329676c0a944575c4c7002c1a14eda0ab06df2b95a703d04139b8e25bc4bafd",
			wantPrivG2: "108516c07a91668c58480ad7626b30c096df80c285d76aa9ff124e35a5dc3ea3",
			wantPubG1:  "81bac3a42514cd503eeaf601a5ef71a42915ff5f81b7c91a3b228391791d4ed5c2bd1bb17ac9e0b106cd6ba972fbdca9",
			wantPubG2:  "b4c58961b71370f342f60ff99ffc06666f7efbfccdcc17aa2638132b407be7d658fa3a394c76dd0a1d46fcd18fcf08a90811e15aa8de73786a3b50f0dc105fd7e2fb8a906645d804b093f8a9a820a219336d91a7d3e9a654ab0c2cfed3b6ee0c",
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
			wantPriv:  "5610023c035ea3ccaf7397fc72abc0f4f5e26f14bb87c403584e525fdd4e88a1",
			wantPubG1: "a93c0b2121a59c4227f03358025a512578f512274d4256f3e41f43ecc038cb56917a4f7e2ce09243d958eb78640bf80e",
			wantPubG2: "b17bfe58f5ac0e3fd065b35fbd4c67355bc215673b11480f0a1afc25b979dc768d7f4ac0685f42bcb2a53109f9a50b8118c6a8ba208f3aa7632be0001cef3d8fbe0f72390733ad990b0c097ad6bb4634308fe1ebea0b0e9bb9f65b1b269995e8",
		},
		{
			name:      "derivation path: m/0H/1H",
			path:      Path{0 + h, 1 + h},
			wantPriv:  "634f5e1dbd6da00d3f6fdef805d9ba5b18c452d5c5e6b221e88704fa2f65aabf",
			wantPubG1: "a25843ef7b4af69a498b3f641622aff00bf6a4a5d7b51ac92c3b9c69c57d74db3010cd2b12d6be6cfff0b9fe14d62a14",
			wantPubG2: "a163ad8792fc74fa44d782a13446e083b6abd36526b8f6ac3cce02ee226f04f304a17e4e041bc6f83ef52788700b389105cfa385165be6c322350c7589ebbd0e6f2573019cc04447eff6d29069da25a7749a6ac162718aad4ffb856b4c41ae50",
		},
		{
			name:      "derivation path: m/0H/1H/2H",
			path:      Path{0 + h, 1 + h, 2 + h},
			wantPriv:  "3ea153e94166ef78e678bfc2e91a90a3a14352a53996e0ab8aca79386dfb9643",
			wantPubG1: "96617b1781317fbfa9449c86c029b1e197a8ddf74bfb03301e6488467bfa61753049a440b2cd5c5c1607dd4b0209360c",
			wantPubG2: "b0cf3fd5dff4946cc9045e5b638620c11627d862604c5157e5fbc1d6dba09cbaeea54ed9db515d7fffc3f39f0559be490fb5ed8987b0ede8f646f7baf669c795108a04568082e8b03d3a0b986066642c0c6ba19520b04fb6d3e5f8bda765298c",
		},
		{
			name:      "derivation path: m/0H/1H/2H/2H",
			path:      Path{0 + h, 1 + h, 2 + h, 2 + h},
			wantPriv:  "494f993bb8543e0f083d01a5fc895544a627f368579995551872b493c57620d8",
			wantPubG1: "b7733ce3508d8be6e40f4e19bf19b532fc9b97b79a3d1ed3d468912123ae0421e5344427a29e07b906479da644811566",
			wantPubG2: "a6bcd91222cc42776c31c80fa6d7f4c5cbc6f4f15db64a10f71d5eb2d062ff62d0482a7ecd8768a3d2ea98da5e25c81209ab986c3a2917aeb411c95041bbf14e92f8d0825e379727f0e0263259c87f8451eba350bed9623272fe0e01285df203",
		},
		{
			name:      "derivation path: m/0H/1H/2H/2H/1000000000H",
			path:      Path{0 + h, 1 + h, 2 + h, 2 + h, 1000000000 + h},
			wantPriv:  "7307dbc1593dd53fade74ee51baaf87c1539c7a52cacad2e9b280d41c3196397",
			wantPubG1: "ac8f9f67bc9cbd52385f759050c29b0d7d233582119171e96cabad6f7bb7619607844f84d32bba852513fa6cd3b71b83",
			wantPubG2: "8df1c40893676593535d245af19f525468013eb2c3a900192ff901307a8d8aaae3bce8f047c9f990459b67c2cc6da32e0edc034853596a6d6403543565481aacde9328cf6e331777b9a66fa5d62b1061752a8e6d84001f212d07caed70d75cc1",
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
			wantPrivG1: "33c25c9358f2b3ad3ddcf7060e163c5b27e3af8d3297b568e0b527e274f7c19a",
			wantPrivG2: "33c25c9358f2b3ad3ddcf7060e163c5b27e3af8d3297b568e0b527e274f7c19a",
			wantPubG1:  "a46249287d576289d2e99ddbf8ebb4115fa91badacf25c1d16b3923d85b2611ac73db5a3c15e78cc789b3e83e4c5b731",
			wantPubG2:  "8dfc988fa761d3c26949f3e0b7d21442e416251321192ff0619ba42452152e84699f48199e53c7e435ff912f949e846f1876a4ca81c9e484f5e8ed591696e8ac05d9a3a944a9bb0f47aec5a7dbb5e58ac32cb8b66282bd2b2b51655648625115",
		},
		{
			name:       "derivation path: m/0H",
			path:       Path{h},
			wantPrivG1: "5610023c035ea3ccaf7397fc72abc0f4f5e26f14bb87c403584e525fdd4e88a1",
			wantPrivG2: "5610023c035ea3ccaf7397fc72abc0f4f5e26f14bb87c403584e525fdd4e88a1",
			wantPubG1:  "a93c0b2121a59c4227f03358025a512578f512274d4256f3e41f43ecc038cb56917a4f7e2ce09243d958eb78640bf80e",
			wantPubG2:  "b17bfe58f5ac0e3fd065b35fbd4c67355bc215673b11480f0a1afc25b979dc768d7f4ac0685f42bcb2a53109f9a50b8118c6a8ba208f3aa7632be0001cef3d8fbe0f72390733ad990b0c097ad6bb4634308fe1ebea0b0e9bb9f65b1b269995e8",
		},
		{
			name:       "derivation path: m/0H/1",
			path:       Path{h, 1},
			wantPrivG1: "674d241570a26d26488b23c808a5495763d86ac6b945b6228b2ac77364420777",
			wantPrivG2: "133f2a1948d7932dfb13b4496e8973e5a7f1a3104f031532ab60a9aa60ba7796",
			wantPubG1:  "9601169d048fe0e552af7b20b865dd48d5b17121657cd6cc03ad4da2471379ea9162552e98a6722c603d9a24cd295c3a",
			wantPubG2:  "838e08ca5a7fe2f7068fb3e978e84d2b77e8421d4ac29d0da8dde28bf97517821f0f4950a1355d1441e52fef3cb3ff16178495a413ae03ab5a08fd82f0749be2ff07f0e293a04cbf77d12ac41bcb16695efa92fcc4a7512177adf5a3f4b1a239",
		},
		{
			name:       "derivation path: m/0H/1/2H",
			path:       Path{h, 1, 2 + h},
			wantPrivG1: "0e611260644e31302061fc6fd06f88c85c44de134b2f68def387c5756f9fa247",
			wantPrivG2: "07bc7452e7f645573d2f0d575a9fa52a43ed45caeca470a5b7b804ef4ebf1e69",
			wantPubG1:  "90bf3cfef1aae20b77d33ecc83ddb8b4abd4bbb4693a72c5c9a186bfff3604dd1954852500cb79d9d9e500c2f6783a52",
			wantPubG2:  "87b68c3845403e293b03d82e36173b28d370b0b626bc4941e42ab6957b4e8d4376131caf304b652c86b7679d457b62960fa6097289d86dab4a12af4f68d1ba874e3f76a22e17a888c6bb080a775bb23ca2341ef91ff5bf801860bdb5da41a615",
		},
		{
			name:       "derivation path: m/0H/1/2H/2",
			path:       Path{h, 1, 2 + h, 2},
			wantPrivG1: "6ba9b834987c30182d7722f584ddaa26a454df5d352ae06ec5440396de6f5d76",
			wantPrivG2: "1d481654ce0b95ad1d437c3079b1d548a889295d719e963e59864f48893038ab",
			wantPubG1:  "a364a765ed027a8ac3bc423f30f5ac7cd413eb36341faa0d36baa51f3535c810b6e5d4196a0e7ad22b213987712ee36f",
			wantPubG2:  "a6cf23189b00c60330bb5017cb63afa362c9f9e6412f4cd4c132f7efd69fd261b8a0ad1318210231dfc9d19f8d2f38e104498b7163521d69f6a1339b1cf15fdf942c0366c0c97d70a7207166af98774e7bb0a3ae833e908528401e9869dacae6",
		},
		{
			name:       "derivation path: m/0H/1/2H/2/1000000000",
			path:       Path{h, 1, 2 + h, 2, 1000000000},
			wantPrivG1: "3410023a0681a3bf76c73b4fa375f61a2bf3b819fce3bc152642ba1f7ccda8d4",
			wantPrivG2: "3f62b68e92a73b3b2d7f422815eab0d1736e0cc4de7be62e455fac38ef511b8a",
			wantPubG1:  "94ada0b7874a39d1e18da0ea05e73a30f1aa96a11140a50530429aee849620ac39f1cd96634feb98c8ea5e997cafb600",
			wantPubG2:  "a07cb22e78477d493f21948e237db54a3ebb5eab78ec2d0c28a4b40fddc7be7aa2f38eb5b4d204cdd73b5a1eaa0b592f1062b0f0db96cd58a1ca2dfbcbc66f76b374b8ae5bcd5a0ac6ea6967efaa587e443628a0912d3f962d526fc3f0a52bea",
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
			wantXPrivG1: "XSECRET1PQP0R4SGK8Y84J2G9LQD2E4W5RYXQRPWKYSG8TT4KUMZD0QF7TT8PSVPNCFWFXK8JKWKNMH8HQC8PV0ZMYL36LRFJJ76K3C94YL38FA7PNGF4LRNQ",
			wantXPrivG2: "XSECRET1PQP0R4SGK8Y84J2G9LQD2E4W5RYXQRPWKYSG8TT4KUMZD0QF7TT8PSCPNCFWFXK8JKWKNMH8HQC8PV0ZMYL36LRFJJ76K3C94YL38FA7PNGJYH4F0",
			wantXPubG1:  "xpublic1pqp0r4sgk8y84j2g9lqd2e4w5ryxqrpwkysg8tt4kumzd0qf7tt8psv9yvfyjsl2hv2ya96vam0uwhdq3t753htdv7fwp694njg7ctvnprtrnmddrc9083nrcnvlg8ex9kucs39tdxg",
			wantXPubG2:  "xpublic1pqp0r4sgk8y84j2g9lqd2e4w5ryxqrpwkysg8tt4kumzd0qf7tt8pscydljvglfmp60pxjj0nuzmay9zzustz2yepryhlqcvm5sj9y9fws35e7jqenefu0ep4l7gjl9y7s3h3sa4ye2quneyy7h5w6kgkjm52cpwe5w55f2dmpar6a3d8mw67tzkr9jutvc5zh54jk5t92eyxy5g44q3law",
		},
		{
			name:        "derivation path: m/0H",
			path:        Path{h},
			wantXPrivG1: "XSECRET1PQXQGPQYQPP5SQ4JYQ75NUCW4VWGHZET3907HQPNWDRZ3VK2E57GYLGNMRHJP2VZKZQPRCQ6750X27UUHL3E2HS857H3X799MSLZQXKZW2F0A6N5G5YX7849T",
			wantXPrivG2: "XSECRET1PQXQGPQYQPP5SQ4JYQ75NUCW4VWGHZET3907HQPNWDRZ3VK2E57GYLGNMRHJP2CZKZQPRCQ6750X27UUHL3E2HS857H3X799MSLZQXKZW2F0A6N5G5YA00RLY",
			wantXPubG1:  "xpublic1pqxqgpqyqpp5sq4jyq75nucw4vwghzet3907hqpnwdrz3vk2e57gylgnmrhjp2v9f8s9jzgd9n3pz0upntqp955f90r63yf6dgft08eqlg0kvqwxt26gh5nm79nsfys7etr4hseqtlq8q0fuwhm",
			wantXPubG2:  "xpublic1pqxqgpqyqpp5sq4jyq75nucw4vwghzet3907hqpnwdrz3vk2e57gylgnmrhjp2c9300l93advpclaqednt775cee4t0pp2eemz9yq7zs6lsjmj7wuw6xh7jkqdp05909j55csn7d9pwq3334ghgsg7w48vv47qqquau7cl0s0wguswvadny9sczt666a5vdps3ls7h6stp6dmnajmrvnfn90ggfvk63",
		},
		{
			name:        "derivation path: m/0H/1",
			path:        Path{h, 1},
			wantXPrivG1: "XSECRET1PQ2QGPQYQPQQJ5T5660PXXYXX2ZSJ2MHZKCWAPX30WVMJXHLU660AC5N89KS0G83SVAXJG9TS5FKJVJYTY0YQ3F2F2A3AS6KXH9ZMVG5T9TRHXEZZQAMSYAE2NA",
			wantXPrivG2: "XSECRET1PQ2QGPQYQPQQMK0SSJRU6REVAXFC8M5SAZ3DFRD7600W2VX4D0WGNVSZ5NAS00CMQZVLJ5X2G67FJM7CNK3YKAZTNUKNLRGCSFUP32V4TVZ565C96W7TQHLVN7U",
			wantXPubG1:  "xpublic1pq2qgpqyqpqqj5t5660pxxyxx2zsj2mhzkcwapx30wvmjxhlu660ac5n89ks0g83sjcq3d8gy3lsw25400vstsewafr2mzufpv47ddnqr44x6y3cn084fzcj496v2vu3vvq7e5fxd99wr55mtvqw",
			wantXPubG2:  "xpublic1pq2qgpqyqpqqmk0ssjru6revaxfc8m5saz3dfrd7600w2vx4d0wgnvsz5nas00cmqsw8q3jj60l30wp50k05h36zd9dm7sssaftpf6rdgmh3gh7t4z7pp7r6f2zsn2hg5g8jjlmeuk0l3v9uyjkjp8tsr4ddq3lvz7p6fhchlqlcw9yaqfjlh05f2csduk9nftmaf9lxy5agjzaad7k3lfvdz8ymw20j9",
		},
		{
			name:        "derivation path: m/0H/1/2H",
			path:        Path{h, 1, 2 + h},
			wantXPrivG1: "XSECRET1PQWQGPQYQPQQC9QYQSQYYHAGUT2VR02X6W39ZTJZSTT2RNNRQY035MK5G37DLKH9SA883GH3SPES3YCRYFCCNQGRPL3HAQMUGEPWYFHSNFVHK3HHNSLZH2MUL5FRSP5ZLK0",
			wantXPrivG2: "XSECRET1PQWQGPQYQPQQC9QYQSQYDKAEFYJMMS69PPPRPA7CKMK55R7W3AVQX6NNGHKZTHAET95G4MMRQQ778G5H87EZ4W0F0P4T448A99FP763W2AJJ8PFDHHQZW7N4LRE5SW5TVFF",
			wantXPubG1:  "xpublic1pqwqgpqyqpqqc9qyqsqyyhagut2vr02x6w39ztjzstt2rnnrqy035mk5g37dlkh9sa883gh3sjzlnelh34t3qka7n8mxg8hdckj4afwa5dya893wf5xrtllekqnw3j4y9y5qvk7wem8jspshk0qa9yqkzurf",
			wantXPubG2:  "xpublic1pqwqgpqyqpqqc9qyqsqydkaefyjmms69ppprpa7ckmk55r7w3avqx6nnghkzthaet95g4mmrqs7mgcwz9gqlzjwcrmqhrv9em9rfhpv9ky67yjs0y92mf276w34phvycu4ucykefvs6mk08290d3fvraxp9egnkrd4d9p9t60drgm4p6w8am2ytsh4zyvdwcgpfm4hv3u5g6pa7gl7klcqxrqhk6a5sdxz5d6n5gl",
		},
		{
			name:        "derivation path: m/0H/1/2H/2",
			path:        Path{h, 1, 2 + h, 2},
			wantXPrivG1: "XSECRET1PQJQGPQYQPQQC9QYQSQYQYR4SYD8NRLYH4AG43LT0JGYAUG3ALW6WTYSF2WZ9MJC9DK2UG4ETXP46NWP5NP7RQXPDWU30TPXA4GN2G4XLT56J4CRWC4ZQ89K7DAWHVA0AKNL",
			wantXPrivG2: "XSECRET1PQJQGPQYQPQQC9QYQSQYQ90WDS78JMZDFTD4ADJDXJEEXDNVW3FCMZQSQGVXKZA0P3FTV3PA3VQW5S9J5EC9ETTGAGD7RQ7D364Y23ZFFT4CEA937TXRY7JYFXQU2KPNHS3E",
			wantXPubG1:  "xpublic1pqjqgpqyqpqqc9qyqsqyqyr4syd8nrlyh4ag43lt0jgyaug3alw6wtysf2wz9mjc9dk2ug4etxz3kffm9a5p84zkrh3pr7v84437dgyltxc6pl2sdx6a228e4xhyppdh96svk5rn66g4jzwv8wyhwxmc3zk50v",
			wantXPubG2:  "xpublic1pqjqgpqyqpqqc9qyqsqyq90wds78jmzdftd4adjdxjeexdnvw3fcmzqsqgvxkza0p3ftv3pa3vznv7gccnvqvvqeshdgp0jmr473k9j0eueqj7nx5cye00m7knlfxrw9q45f3sggzx80un5vl35hn3cgyfx9hzc6jr45ldgfnnvw0zh7ljskqxekqe97hpfeqw9n2lxrhfeampgawsvlfppfggq0fs6w6etnqqvxyc7",
		},
		{
			name:        "derivation path: m/0H/1/2H/2/1000000000",
			path:        Path{h, 1, 2 + h, 2, 1000000000},
			wantXPrivG1: "XSECRET1PQKQGPQYQPQQC9QYQSQYQ9QY5A0WQXLSAHT8CK43XPC2D97TS8C5VV06VS07Z6APNFRGJ4QKAA00SM9T9XQ6PQQ36Q6Q680MKCUA5LGM47CDZHUACR87W80Q4YEPT58MUEK5DGA82UND",
			wantXPrivG2: "XSECRET1PQKQGPQYQPQQC9QYQSQYQ9QY5A0WQX3VZ7AAQ85SQYHYH5WC2ZUU6064LMF6HKJ78ZWW96HXSJVMS9477VQLK9D5WJ2NNKWED0APZS902KRGHXMSVCN08HE3WG406CW802YDC5ZJCWKU",
			wantXPubG1:  "xpublic1pqkqgpqyqpqqc9qyqsqyq9qy5a0wqxlsaht8ck43xpc2d97ts8c5vv06vs07z6apnfrgj4qkaa00sm9t9xz22mg9hsa9rn50p3ksw5p088gc0r25k5yg5pfg9xppf4m5yjcs2cw03ektxxnltnryw5h5e0jhmvqqs9q7tc",
			wantXPubG2:  "xpublic1pqkqgpqyqpqqc9qyqsqyq9qy5a0wqx3vz7aaq85sqyhyh5wc2zuu6064lmf6hkj78zww96hxsjvms9477vzs8ev3w0prh6jflyx2gugmak49raw674duwctgv9zjtgr7ac7l84ghn366mf5syehtnkks74g94jtcsv2c0pkuke4v2rj3dl09uvmmkkd6t3tjme4dq43h2d9n7l2jc0ezrv29qjyknl93d2fhu8u99904qdlqv7t",
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
