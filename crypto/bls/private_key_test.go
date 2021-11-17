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
	_, priv1 := GenerateTestKeyPair()
	priv2 := new(BLSPrivateKey)
	priv3 := new(BLSPrivateKey)
	priv4 := new(BLSPrivateKey)

	js, err := json.Marshal(priv1)
	assert.NoError(t, err)
	require.Error(t, priv2.UnmarshalJSON([]byte("bad")))
	require.NoError(t, json.Unmarshal(js, priv2))

	bs, err := priv2.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, priv3.UnmarshalCBOR(bs))

	txt, err := priv2.MarshalText()
	assert.NoError(t, err)
	assert.NoError(t, priv4.UnmarshalText(txt))

	require.True(t, priv1.EqualsTo(priv4))
	require.NoError(t, priv1.SanityCheck())
}

func TestPrivateKeyFromBytes(t *testing.T) {
	_, err := PrivateKeyFromRawBytes(nil)
	assert.Error(t, err)
	_, priv1 := GenerateTestKeyPair()
	priv2, err := PrivateKeyFromRawBytes(priv1.RawBytes())
	assert.NoError(t, err)
	require.True(t, priv1.EqualsTo(priv2))

	inv, _ := hex.DecodeString(strings.Repeat("ff", PrivateKeySize))
	_, err = PrivateKeyFromRawBytes(inv)
	assert.Error(t, err)
}

func TestPrivateKeyFromString(t *testing.T) {
	_, priv1 := GenerateTestKeyPair()
	priv2, err := PrivateKeyFromString(priv1.String())
	assert.NoError(t, err)
	require.True(t, priv1.EqualsTo(priv2))

	_, err = PrivateKeyFromString("inv")
	assert.Error(t, err)
}

func TestMarshalingEmptyPrivateKey(t *testing.T) {
	pv1 := &BLSPrivateKey{}

	js, err := json.Marshal(pv1)
	assert.NoError(t, err)
	assert.Equal(t, js, []byte{0x22, 0x22}) // ""
	var pv2 BLSPrivateKey
	err = json.Unmarshal(js, &pv2)
	assert.Error(t, err)

	bs, err := pv1.MarshalCBOR()
	assert.Error(t, err)
	var pv3 BLSPrivateKey
	err = pv3.UnmarshalCBOR(bs)
	assert.Error(t, err)
}

func TestPrivateKeyFromSeed(t *testing.T) {
	priv, err := PrivateKeyFromSeed([]byte{0})
	assert.NoError(t, err)
	assert.NoError(t, priv.SanityCheck())
}
