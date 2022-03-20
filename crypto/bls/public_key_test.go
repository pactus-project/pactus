package bls

import (
	"encoding/hex"
	"strings"
	"testing"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

func TestPublicKeyMarshaling(t *testing.T) {
	pub1, _ := GenerateTestKeyPair()
	pub2 := new(PublicKey)

	bs, err := pub1.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, pub2.UnmarshalCBOR(bs))
	assert.True(t, pub1.EqualsTo(pub2))
	assert.NoError(t, pub1.SanityCheck())

	js, err := pub1.MarshalJSON()
	assert.NoError(t, err)
	assert.Contains(t, string(js), pub1.String())

	inv, _ := hex.DecodeString(strings.Repeat("ff", PublicKeySize))
	data, _ := cbor.Marshal(inv)
	assert.Error(t, pub2.UnmarshalCBOR(data))
}

func TestPublicKeyFromString(t *testing.T) {
	pub1, _ := GenerateTestKeyPair()
	pub2, err := PublicKeyFromString(pub1.String())
	assert.NoError(t, err)
	assert.True(t, pub1.EqualsTo(pub2))

	_, err = PublicKeyFromString("")
	assert.Error(t, err)

	_, err = PublicKeyFromString("inv")
	assert.Error(t, err)

	_, err = PublicKeyFromString("00")
	assert.Error(t, err)
}

func TestPublicKeyEmpty(t *testing.T) {
	pub1 := PublicKey{}

	bs, err := pub1.MarshalCBOR()
	assert.Error(t, err)
	assert.Empty(t, pub1.String())
	assert.Empty(t, pub1.RawBytes())

	var pub2 PublicKey
	err = pub2.UnmarshalCBOR(bs)
	assert.Error(t, err)
}

func TestPublicKeySanityCheck(t *testing.T) {
	pub, err := PublicKeyFromString("C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	assert.Error(t, pub.SanityCheck())
}

func TestPublicKeyVerifyAddress(t *testing.T) {
	pub1, _ := GenerateTestKeyPair()
	pub2, _ := GenerateTestKeyPair()

	assert.True(t, pub1.VerifyAddress(pub1.Address()))
	assert.False(t, pub1.VerifyAddress(pub2.Address()))
}
