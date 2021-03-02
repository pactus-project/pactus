package crypto

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicKeyMarshaling(t *testing.T) {
	_, pub1, _ := GenerateTestKeyPair()
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

	require.True(t, pub1.EqualsTo(*pub4))
	require.NoError(t, pub1.SanityCheck())
}

func TestPublicKeyFromBytes(t *testing.T) {
	_, err := PublicKeyFromRawBytes(nil)
	assert.Error(t, err)
	_, pub1, _ := GenerateTestKeyPair()
	pub2, err := PublicKeyFromRawBytes(pub1.RawBytes())
	assert.NoError(t, err)
	require.True(t, pub1.EqualsTo(pub2))

	inv, _ := hex.DecodeString(strings.Repeat("ff", PublicKeySize))
	_, err = PublicKeyFromRawBytes(inv)
	assert.Error(t, err)
}

func TestPublicKeyFromString(t *testing.T) {
	_, pub1, _ := GenerateTestKeyPair()
	pub2, err := PublicKeyFromString(pub1.String())
	assert.NoError(t, err)
	require.True(t, pub1.EqualsTo(pub2))

	_, err = PublicKeyFromString("inv")
	assert.Error(t, err)
}

func TestMarshalingEmptyPublicKey(t *testing.T) {
	pb1 := PublicKey{}

	js, err := json.Marshal(pb1)
	assert.NoError(t, err)
	assert.Equal(t, js, []byte{0x22, 0x22}) // ""
	var pb2 PublicKey
	err = json.Unmarshal(js, &pb2)
	assert.Error(t, err)

	bs, err := pb1.MarshalCBOR()
	assert.Error(t, err)

	var pb3 PublicKey
	err = pb3.UnmarshalCBOR(bs)
	assert.Error(t, err)

	assert.Equal(t, pb1.Address().String(), "a13dcd79a9c0e088317b189c88641f4ab671ab15") // Invalid address
}

func TestPublicKeyToAddress(t *testing.T) {
	addr, err := AddressFromString("fc617c7591b3a1db27281497ce9d57f76faef212")
	assert.NoError(t, err)
	pub, err := PublicKeyFromString("ccd169a31a7bfa611480072137f77efd8c1cfb0f811957972d15bab4e8c8998ade29d99b03815d3873e57d21e67ce210480270ca0b77698de0623ab1e6a241bd05a00a2e3a5b319c99fa1b9ecb6f53564e4c53dbb8a2b6b46315bf258208f614")
	assert.NoError(t, err)
	assert.Equal(t, pub.Address(), addr)
}
