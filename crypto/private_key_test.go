package crypto

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalingEmptyPrivateKey(t *testing.T) {
	pv1 := PrivateKey{}

	js, err := json.Marshal(pv1)
	assert.Error(t, err)
	var pv2 PrivateKey
	err = json.Unmarshal(js, &pv2)
	assert.Error(t, err)

	bs, err := pv1.MarshalCBOR()
	assert.Error(t, err)
	var pv3 PrivateKey
	err = pv3.UnmarshalCBOR(bs)
	assert.Error(t, err)
}

func TestMarshalingPrivateKey(t *testing.T) {
	_, _, pv0 := GenerateTestKeyPair()
	_, _, pv1 := GenerateTestKeyPair()
	js, err := json.Marshal(pv1)
	//fmt.Println(string(js))
	assert.NoError(t, err)

	var pv2 PrivateKey
	assert.NoError(t, json.Unmarshal(js, &pv2))
	require.Equal(t, pv1, pv2)

	bs, err := pv1.MarshalCBOR()
	assert.NoError(t, err)

	var pv3 PrivateKey
	assert.NoError(t, pv3.UnmarshalCBOR(bs))
	require.False(t, pv2.EqualsTo(pv0))
	require.True(t, pv2.EqualsTo(pv3))
	require.Equal(t, pv3, pv1)

	require.Nil(t, pv3.Sign([]byte{}))
}
