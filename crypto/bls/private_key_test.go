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
	assert.Empty(t, prv1.RawBytes())

	var prv2 PrivateKey
	err = prv2.UnmarshalCBOR(bs)
	assert.Error(t, err)
}

func TestPrivateKeyFromSeed(t *testing.T) {
	_, err := PrivateKeyFromSeed([]byte{0})
	assert.Error(t, err)
	seed := [32]byte{}
	prv, err := PrivateKeyFromSeed(seed[:])
	assert.NoError(t, err)
	assert.Equal(t, prv.RawBytes(), []byte{74, 53, 59, 227, 218, 192, 145, 160, 167, 230, 64, 98, 3, 114, 245, 225, 226, 228, 64, 23, 23, 193, 231, 156, 172, 111, 251, 168, 246, 144, 86, 4})
}

func TestPrivateKeySanityCheck(t *testing.T) {
	prv, err := PrivateKeyFromString("0000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	assert.Error(t, prv.SanityCheck())
}
