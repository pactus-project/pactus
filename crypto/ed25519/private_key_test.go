package ed25519_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv1 := ts.RandED25519KeyPair()
	_, prv2 := ts.RandED25519KeyPair()

	fmt.Println(prv1.String())

	assert.True(t, prv1.EqualsTo(prv1))
	assert.False(t, prv1.EqualsTo(prv2))
	assert.Equal(t, prv1, prv1)
	assert.NotEqual(t, prv1, prv2)
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
			"XXXXXXR2SYCC5TDQKMJ73J64J8GJTMTKREEQNQAS0M5SLZ9LVJV7Y940NVZQD9JUS" +
				"GV2N44C9H5PVGRXARNGZ7QF3PSKH7805E5SZXPE7ZHHAGX0NFQR",
			false, nil,
		},
		{
			"invalid checksum (expected s9c56g got czlgh0)",
			"SECRET1RAC7048K666DCCYG7FJW68ZE2G6P32UAPLRLWDV3RTAR4PWZUX2CDSFAL55VM" +
				"YS06CY35LA72AWZN5DY5NZA078S4S4K654UFJ0YCCZLGH0",
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
			"",
			"SECRET1RSW0NCUTJDFMWMSEYQVMXLW9ZP3WT2920S24XT55R6YKHL2G3ZG74C39" +
				"GHXVZRY4F7F9ZZ5VJNWFYWSSN4MCFSSCP34DKY0ML0EF26SQ3J5575",
			true,
			[]byte{},
		},
	}

	for no, test := range tests {
		prv, err := ed25519.PrivateKeyFromString(test.encoded)
		if test.valid {
			assert.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, test.result, prv.Bytes(), "test %v: invalid bytes", no)
			assert.Equal(t, strings.ToUpper(test.encoded), prv.String(), "test %v: invalid encoded", no)
		} else {
			assert.Contains(t, err.Error(), test.errMsg, "test %v: error not matched", no)
		}
	}
}

// // TestKeyGen ensures the KeyGen function works as intended.
// func TestKeyGen(t *testing.T) {
// 	tests := []struct {
// 		ikm string
// 		sk  string
// 	}{
// 		{
// 			"",
// 			"Err",
// 		},
// 		{
// 			"00000000000000000000000000000000000000000000000000000000000000",
// 			"Err",
// 		},
// 		{
// 			"0000000000000000000000000000000000000000000000000000000000000000",
// 			"4d129a19df86a0f5345bad4cc6f249ec2a819ccc3386895beb4f7d98b3db6235",
// 		},
// 		{
// 			"2b1eb88002e83a622792d0b96d4f0695e328f49fdd32480ec0cf39c2c76463af",
// 			"0000f678e80740072a4a7fe8c7344db88a00ccc7db36aa51fa51f9c68e561584",
// 		},
// 		// The test vectors from EIP-2333
// 		// https://github.com/ethereum/EIPs/blob/784107449bd83a9327b54f82aba96de28d72b89a/EIPS/eip-2333.md#test-cases
// 		{
// 			"c55257c360c07c72029aebc1b53c05ed0362ada38ead3e3e9efa3708e5349553" +
// 				"1f09a6987599d18264c1e1c92f2cf141630c7a3c4ab7c81b2f001698e7463b04",
// 			"0d7359d57963ab8fbbde1852dcf553fedbc31f464d80ee7d40ae683122b45070",
// 		},
// 		{
// 			"3141592653589793238462643383279502884197169399375105820974944592",
// 			"41c9e07822b092a93fd6797396338c3ada4170cc81829fdfce6b5d34bd5e7ec7",
// 		},
// 		{
// 			"0099FF991111002299DD7744EE3355BBDD8844115566CC55663355668888CC00",
// 			"3cfa341ab3910a7d00d933d8f7c4fe87c91798a0397421d6b19fd5b815132e80",
// 		},
// 		{
// 			"d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3",
// 			"2a0e28ffa5fbbe2f8e7aad4ed94f745d6bf755c51182e119bb1694fe61d3afca",
// 		},
// 	}

// 	for i, test := range tests {
// 		ikm, _ := hex.DecodeString(test.ikm)
// 		prv, err := bls.KeyGen(ikm, nil)
// 		if test.sk == "Err" {
// 			assert.Error(t, err,
// 				"test '%v' failed. no error", i)
// 		} else {
// 			assert.NoError(t, err,
// 				"test'%v' failed. has error", i)
// 			assert.Equal(t, test.sk, hex.EncodeToString(prv.Bytes()),
// 				"test '%v' failed. not equal", i)
// 		}
// 	}
// }
