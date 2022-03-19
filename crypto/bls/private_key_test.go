package bls

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrivateKeyMarshaling(t *testing.T) {
	_, prv1 := GenerateTestKeyPair()
	prv2 := new(PrivateKey)
	prv3 := new(PrivateKey)
	prv4 := new(PrivateKey)

	js, err := json.Marshal(prv1)
	assert.NoError(t, err)
	require.Error(t, prv2.UnmarshalJSON([]byte("bad")))
	require.NoError(t, json.Unmarshal(js, prv2))

	bs, err := prv2.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, prv3.UnmarshalCBOR(bs))

	txt, err := prv2.MarshalText()
	assert.NoError(t, err)
	assert.NoError(t, prv4.UnmarshalText(txt))

	require.True(t, prv1.EqualsTo(prv4))
	require.NoError(t, prv1.SanityCheck())
}

func TestPrivateKeyFromBytes(t *testing.T) {
	_, err := PrivateKeyFromRawBytes(nil)
	assert.Error(t, err)
	_, prv1 := GenerateTestKeyPair()
	prv2, err := PrivateKeyFromRawBytes(prv1.RawBytes())
	assert.NoError(t, err)
	require.True(t, prv1.EqualsTo(prv2))

	inv, _ := hex.DecodeString(strings.Repeat("ff", PrivateKeySize))
	_, err = PrivateKeyFromRawBytes(inv)
	assert.Error(t, err)
}

func TestPrivateKeyFromString(t *testing.T) {
	_, prv1 := GenerateTestKeyPair()
	prv2, err := PrivateKeyFromString(prv1.String())
	assert.NoError(t, err)
	require.True(t, prv1.EqualsTo(prv2))

	_, err = PrivateKeyFromString("inv")
	assert.Error(t, err)
}

func TestMarshalingEmptyPrivateKey(t *testing.T) {
	prv1 := &PrivateKey{}

	js, err := json.Marshal(prv1)
	assert.NoError(t, err)
	assert.Equal(t, js, []byte{0x22, 0x22}) // ""
	var prv2 PrivateKey
	err = json.Unmarshal(js, &prv2)
	assert.Error(t, err)

	bs, err := prv1.MarshalCBOR()
	assert.Error(t, err)
	var pv3 PrivateKey
	err = pv3.UnmarshalCBOR(bs)
	assert.Error(t, err)
}

func TestPrivateKeyFromSeed(t *testing.T) {
	_, err := PrivateKeyFromSeed([]byte{0})
	assert.Error(t, err)
	seed := [32]byte{}
	prv, err := PrivateKeyFromSeed(seed[:])
	assert.NoError(t, err)
	assert.NoError(t, prv.SanityCheck())
	assert.Equal(t, prv.RawBytes(), []byte{74, 53, 59, 227, 218, 192, 145, 160, 167, 230, 64, 98, 3, 114, 245, 225, 226, 228, 64, 23, 23, 193, 231, 156, 172, 111, 251, 168, 246, 144, 86, 4})
}
