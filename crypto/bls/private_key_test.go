package bls_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv1 := ts.RandBLSKeyPair()
	_, prv2 := ts.RandBLSKeyPair()
	_, prv3 := ts.RandEd25519KeyPair()

	assert.True(t, prv1.EqualsTo(prv1))
	assert.False(t, prv1.EqualsTo(prv2))
	assert.False(t, prv1.EqualsTo(prv3))
}

func TestPrivateKeyFromString(t *testing.T) {
	tests := []struct {
		errMsg  string
		encoded string
		valid   bool
		result  []byte
	}{
		{
			"invalid separator index -1",
			"not_proper_encoded",
			false, nil,
		},
		{
			"invalid bech32 string length 0",
			"",
			false, nil,
		},
		{
			"invalid character not part of charset: 105",
			"SECRET1IOIOOI",
			false, nil,
		},
		{
			"invalid bech32 string length 0",
			"SECRET1HPZZU9",
			false, nil,
		},
		{
			"invalid HRP: xxx",
			"XXX1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQMUUMJT",
			false, nil,
		},
		{
			"invalid checksum (expected qjvk67 got qjvk68)",
			"SECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQQJVK68",
			false, nil,
		},
		{
			"invalid length: 31",
			"SECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CCZ0EU7Z",
			false, nil,
		},
		{
			"invalid signature type: 2",
			"SECRET1ZDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQG04E54",
			false, nil,
		},
		{
			"",
			"secret1pdrwtlp5px0fahdx39gxzjp7fkzfalml0d5u9tt9kvqhduc99cmgqqjvk67", // lowercase
			true,
			[]byte{
				0x68, 0xdc, 0xbf, 0x86, 0x81, 0x33, 0xd3, 0xdb, 0xb4, 0xd1, 0x2a, 0xc, 0x29, 0x7, 0xc9, 0xb0,
				0x93, 0xdf, 0xef, 0xef, 0x6d, 0x38, 0x55, 0xac, 0xb6, 0x60, 0x2e, 0xde, 0x60, 0xa5, 0xc6, 0xd0,
			},
		},
		{
			"",
			"SECRET1PDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CMGQQJVK67",
			true,
			[]byte{
				0x68, 0xdc, 0xbf, 0x86, 0x81, 0x33, 0xd3, 0xdb, 0xb4, 0xd1, 0x2a, 0xc, 0x29, 0x7, 0xc9, 0xb0,
				0x93, 0xdf, 0xef, 0xef, 0x6d, 0x38, 0x55, 0xac, 0xb6, 0x60, 0x2e, 0xde, 0x60, 0xa5, 0xc6, 0xd0,
			},
		},
	}
	for no, tt := range tests {
		prv, err := bls.PrivateKeyFromString(tt.encoded)
		if tt.valid {
			assert.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, tt.result, prv.Bytes(), "test %v: invalid bytes", no)
			assert.Equal(t, strings.ToUpper(tt.encoded), prv.String(), "test %v: invalid encoded", no)
		} else {
			assert.Contains(t, err.Error(), tt.errMsg, "test %v: error not matched", no)
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
			"c55257c360c07c72029aebc1b53c05ed0362ada38ead3e3e9efa3708e5349553" +
				"1f09a6987599d18264c1e1c92f2cf141630c7a3c4ab7c81b2f001698e7463b04",
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

	for no, tt := range tests {
		ikm, _ := hex.DecodeString(tt.ikm)
		prv, err := bls.KeyGen(ikm, nil)
		if tt.sk == "Err" {
			assert.Error(t, err,
				"test '%v' failed. no error", no)
		} else {
			assert.NoError(t, err,
				"test'%v' failed. has error", no)
			assert.Equal(t, tt.sk, hex.EncodeToString(prv.Bytes()),
				"test '%v' failed. not equal", no)
		}
	}
}
