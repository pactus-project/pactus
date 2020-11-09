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
