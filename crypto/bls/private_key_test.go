package bls

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyMarshaling(t *testing.T) {
	_, prv1 := GenerateTestKeyPair()
	prv2 := new(PrivateKey)

	bs, err := prv1.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, prv2.UnmarshalCBOR(bs))
	assert.True(t, prv1.EqualsTo(prv2))
	assert.NoError(t, prv1.SanityCheck())

	js, err := prv1.MarshalJSON()
	assert.NoError(t, err)
	assert.Contains(t, string(js), prv1.String())

	inv, _ := hex.DecodeString(strings.Repeat("ff", PrivateKeySize))
	data, _ := cbor.Marshal(inv)
	assert.Error(t, prv2.UnmarshalCBOR(data))
}

func TestPrivateKeyFromString(t *testing.T) {
	_, prv1 := GenerateTestKeyPair()
	prv2, err := PrivateKeyFromString(prv1.String())
	assert.NoError(t, err)
	assert.True(t, prv1.EqualsTo(prv2))

	_, err = PrivateKeyFromString("")
	assert.Error(t, err)

	_, err = PrivateKeyFromString("inv")
	assert.Error(t, err)

	_, err = PrivateKeyFromString("00")
	assert.Error(t, err)
}

func TestPrivateKeyEmpty(t *testing.T) {
	prv1 := &PrivateKey{}

	bs, err := prv1.MarshalCBOR()
	assert.Error(t, err)
	assert.Empty(t, prv1.String())
	assert.Empty(t, prv1.Bytes())

	var prv2 PrivateKey
	err = prv2.UnmarshalCBOR(bs)
	assert.Error(t, err)
}

func TestPrivateKeyFromSeed(t *testing.T) {
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
		/// The test vectors from EIP-2333
		/// https://github.com/ethereum/EIPs/blob/784107449bd83a9327b54f82aba96de28d72b89a/EIPS/eip-2333.md#test-cases
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
		prv, err := PrivateKeyFromSeed(ikm)
		if test.sk == "Err" {
			assert.Error(t, err, "test #i failed", i)
		} else {
			assert.NoError(t, err, "test #i failed", i)
			assert.Equal(t, prv.String(), test.sk, "test #i failed", i)
		}
	}
}

func TestPrivateKeySanityCheck(t *testing.T) {
	prv, err := PrivateKeyFromString("0000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	assert.Error(t, prv.SanityCheck())
}
