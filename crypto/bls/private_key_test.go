package bls

import (
	"encoding/hex"
	"testing"

	"github.com/herumi/bls-go-binary/bls"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/util/errors"
)

func TestPrivateKeyToString(t *testing.T) {
	_, randomPrv := GenerateTestKeyPair()
	tests := []struct {
		name      string
		encoded   string
		decodable bool
		valid     bool
		result    []byte
	}{
		{
			"empty private key",
			"",
			false, false,
			nil,
		},
		{
			"invalid character",
			"INV",
			false, false,
			nil,
		},
		{
			"no type",
			"SECRET1HPZZU9",
			false, false,
			nil,
		},
		{
			"invalid private key",
			"SECRET1PLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLSJGV7U3",
			false, false,
			nil,
		},
		{
			"invalid hrp",
			"XXX1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQ7E8ACD",
			false, false,
			nil,
		},
		{
			"invalid checksum",
			"SECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQQJVK68",
			false, false,
			nil,
		},
		{
			"invalid length",
			"SECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CCZ0EU7Z",
			false, false,
			nil,
		},
		{
			"invalid type",
			"SECRET1ZDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQG04E54",
			false, false,
			nil,
		},
		{
			"valid private key in lowercase format",
			"secret1pdrwtlp5px0fahdx39gxzjp7fkzfalml0d5u9tt9kvqhduc99cmgqqjvk67",
			true, true,
			[]byte{0x68, 0xdc, 0xbf, 0x86, 0x81, 0x33, 0xd3, 0xdb, 0xb4, 0xd1, 0x2a, 0xc, 0x29, 0x7, 0xc9, 0xb0,
				0x93, 0xdf, 0xef, 0xef, 0x6d, 0x38, 0x55, 0xac, 0xb6, 0x60, 0x2e, 0xde, 0x60, 0xa5, 0xc6, 0xd0},
		},
		{
			"valid private key",
			"SECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQQJVK67",
			true, true,
			[]byte{0x68, 0xdc, 0xbf, 0x86, 0x81, 0x33, 0xd3, 0xdb, 0xb4, 0xd1, 0x2a, 0xc, 0x29, 0x7, 0xc9, 0xb0,
				0x93, 0xdf, 0xef, 0xef, 0x6d, 0x38, 0x55, 0xac, 0xb6, 0x60, 0x2e, 0xde, 0x60, 0xa5, 0xc6, 0xd0},
		},
		{
			"random private key",
			randomPrv.String(),
			true, true,
			randomPrv.secretKey.Serialize(),
		},
	}
	for _, test := range tests {
		prv, err := PrivateKeyFromString(test.encoded)
		if test.decodable {
			assert.NoError(t, err, "test %v. unexpected error", test.name)
			assert.Equal(t, prv.SanityCheck() == nil, test.valid, "test %v. sanity check failed", test.name)

			prv2, _ := PrivateKeyFromBytes(test.result)
			assert.True(t, prv.EqualsTo(prv2))
		} else {
			assert.Error(t, err, "test %v. should failed", test.name)
			assert.Equal(t, errors.Code(err), errors.ErrInvalidPrivateKey)
		}
	}
}

// TestKeyGen ensures the KeyGen function works as intended.
func TestKeyGen(t *testing.T) {
	tests := []struct {
		ikm string
		sk  string
	}{
		{
			"",
			"Err",
		},
		{
			"00000000000000000000000000000000000000000000000000000000000000",
			"Err",
		},
		{
			"0000000000000000000000000000000000000000000000000000000000000000",
			"4d129a19df86a0f5345bad4cc6f249ec2a819ccc3386895beb4f7d98b3db6235",
		},
		{
			"2b1eb88002e83a622792d0b96d4f0695e328f49fdd32480ec0cf39c2c76463af",
			"0000f678e80740072a4a7fe8c7344db88a00ccc7db36aa51fa51f9c68e561584",
		},
		// The test vectors from EIP-2333
		// https://github.com/ethereum/EIPs/blob/784107449bd83a9327b54f82aba96de28d72b89a/EIPS/eip-2333.md#test-cases
		{
			"c55257c360c07c72029aebc1b53c05ed0362ada38ead3e3e9efa3708e53495531f09a6987599d18264c1e1c92f2cf141630c7a3c4ab7c81b2f001698e7463b04",
			"0d7359d57963ab8fbbde1852dcf553fedbc31f464d80ee7d40ae683122b45070",
		},
		{
			"3141592653589793238462643383279502884197169399375105820974944592",
			"41c9e07822b092a93fd6797396338c3ada4170cc81829fdfce6b5d34bd5e7ec7",
		},
		{
			"0099FF991111002299DD7744EE3355BBDD8844115566CC55663355668888CC00",
			"3cfa341ab3910a7d00d933d8f7c4fe87c91798a0397421d6b19fd5b815132e80",
		},
		{
			"d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3",
			"2a0e28ffa5fbbe2f8e7aad4ed94f745d6bf755c51182e119bb1694fe61d3afca",
		},
	}

	for i, test := range tests {
		ikm, _ := hex.DecodeString(test.ikm)
		prv, err := KeyGen(ikm, nil)
		if test.sk == "Err" {
			assert.Error(t, err, "test %v failed", i)
		} else {
			assert.NoError(t, err, "test %v failed", i)
			assert.Equal(t, hex.EncodeToString(prv.Bytes()), test.sk,
				"test %v failed", i)
		}
	}
}

func TestPrivateKeySanityCheck(t *testing.T) {
	sc := new(bls.SecretKey)
	err := sc.DeserializeHexStr("0000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	prv := PrivateKey{
		secretKey: *sc,
	}
	assert.NoError(t, err)
	assert.Error(t, prv.SanityCheck())
}
