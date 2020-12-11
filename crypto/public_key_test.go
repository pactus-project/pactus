package crypto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalingEmptyPublicKey(t *testing.T) {
	pb1 := PublicKey{}

	js, err := json.Marshal(pb1)
	assert.Error(t, err)
	var pb2 PublicKey
	err = json.Unmarshal(js, &pb2)
	assert.Error(t, err)

	bs, err := pb1.MarshalCBOR()
	assert.Error(t, err)

	var pb3 PublicKey
	err = pb3.UnmarshalCBOR(bs)
	assert.Error(t, err)
}

func TestMarshalingPublicKey(t *testing.T) {
	_, pb0, _ := RandomKeyPair()
	_, pb1, _ := RandomKeyPair()
	js, err := json.Marshal(&pb1)
	assert.NoError(t, err)

	var pb2 PublicKey
	require.NoError(t, json.Unmarshal(js, &pb2))
	require.False(t, pb1.EqualsTo(pb0))
	require.True(t, pb1.EqualsTo(pb2))

	bs, err := pb1.MarshalCBOR()
	assert.NoError(t, err)

	pb3 := new(PublicKey)
	assert.NoError(t, pb3.UnmarshalCBOR(bs))

	require.True(t, pb3.EqualsTo(pb2))
}

func TestPublicKeyToAddress(t *testing.T) {
	addr, err := AddressFromString("ea378dac42172c7836ed17753e105e66e024508a")
	assert.NoError(t, err)
	pub, err := PublicKeyFromString("ccd169a31a7bfa611480072137f77efd8c1cfb0f811957972d15bab4e8c8998ade29d99b03815d3873e57d21e67ce210480270ca0b77698de0623ab1e6a241bd05a00a2e3a5b319c99fa1b9ecb6f53564e4c53dbb8a2b6b46315bf258208f614")
	assert.NoError(t, err)

	assert.Equal(t, pub.Address(), addr)
}
