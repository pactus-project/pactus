package bls

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicKeyMarshaling(t *testing.T) {
	pub1, _ := GenerateTestKeyPair()
	pub2 := new(PublicKey)
	pub3 := new(PublicKey)
	pub4 := new(PublicKey)

	js, err := json.Marshal(pub1)
	assert.NoError(t, err)
	require.Error(t, pub2.UnmarshalJSON([]byte("bad")))
	require.NoError(t, json.Unmarshal(js, pub2))

	bs, err := pub2.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, pub3.UnmarshalCBOR(bs))

	txt, err := pub2.MarshalText()
	assert.NoError(t, err)
	assert.NoError(t, pub4.UnmarshalText(txt))

	require.True(t, pub1.EqualsTo(pub4))
	require.NoError(t, pub1.SanityCheck())
}

func TestPublicKeyFromBytes(t *testing.T) {
	_, err := PublicKeyFromRawBytes(nil)
	assert.Error(t, err)
	pub1, _ := GenerateTestKeyPair()
	pub2, err := PublicKeyFromRawBytes(pub1.RawBytes())
	assert.NoError(t, err)
	require.True(t, pub1.EqualsTo(pub2))

	inv, _ := hex.DecodeString(strings.Repeat("ff", PublicKeySize))
	_, err = PublicKeyFromRawBytes(inv)
	assert.Error(t, err)
}

func TestPublicKeyFromString(t *testing.T) {
	pub1, _ := GenerateTestKeyPair()
	pub2, err := PublicKeyFromString(pub1.String())
	assert.NoError(t, err)
	require.True(t, pub1.EqualsTo(pub2))

	_, err = PublicKeyFromString("inv")
	assert.Error(t, err)
}

func TestMarshalingEmptyPublicKey(t *testing.T) {
	pub1 := PublicKey{}

	js, err := json.Marshal(pub1)
	assert.NoError(t, err)
	assert.Equal(t, js, []byte{0x22, 0x22}) // ""
	var pub2 PublicKey
	err = json.Unmarshal(js, &pub2)
	assert.Error(t, err)

	bs, err := pub1.MarshalCBOR()
	assert.Error(t, err)

	var pub3 PublicKey
	err = pub3.UnmarshalCBOR(bs)
	assert.Error(t, err)

	assert.Equal(t, pub1.Address().String(), "zc15y7u67dfcrsgsvtmrzwgseqlf2m8r2c4qs0v3g")
}
