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
	_, pv0 := GenerateRandomKey()
	_, pv1 := GenerateRandomKey()
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

	defer func() { recover() }()

	pv3.Sign([]byte{})
	t.Errorf("did not panic")
}

/*
func TestPrivateKeyValidity(t *testing.T) {
	var err error
	_, err = PrivateKeyFromString("skZfztcE4vkJLYNQ3TcvAkgH24TV1hQfuojiwReVto9JknsoWNZPJVmd6agFiCyGx1px45HJjgRQvRNRrc4oeqZgaPXhQHM")
	assert.NoError(t, err)

	_, err = PrivateKeyFromString("skzfztcE4vkJLYNQ3TcvAkgH24TV1hQfuojiwReVto9JknsoWNZPJVmd6agFiCyGx1px45HJjgRQvRNRrc4oeqZgaPXhQHM")
	assert.Error(t, err)

	_, err = PrivateKeyFromString("SKZfztcE4vkJLYNQ3TcvAkgH24TV1hQfuojiwReVto9JknsoWNZPJVmd6agFiCyGx1px45HJjgRQvRNRrc4oeqZgaPXhQHM")
	assert.Error(t, err)

	_, err = PrivateKeyFromString("invalid_private_key")
	assert.Error(t, err)

	_, err = PrivateKeyFromRawBytes([]byte{0, 1, 2, 3, 4, 5, 6})
	assert.Error(t, err)
}
*/
