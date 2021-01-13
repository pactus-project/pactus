package crypto

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashMarshaling(t *testing.T) {
	hash1 := GenerateTestHash()
	hash2 := new(Hash)
	hash3 := new(Hash)
	hash4 := new(Hash)

	js, err := json.Marshal(hash1)
	assert.NoError(t, err)
	require.Error(t, hash2.UnmarshalJSON([]byte("bad")))
	require.NoError(t, json.Unmarshal(js, hash2))

	bs, err := hash2.MarshalCBOR()
	assert.NoError(t, err)
	assert.NoError(t, hash3.UnmarshalCBOR(bs))

	txt, err := hash2.MarshalText()
	assert.NoError(t, err)
	assert.NoError(t, hash4.UnmarshalText(txt))

	require.True(t, hash1.EqualsTo(*hash4))
}

func TestHashFromBytes(t *testing.T) {
	_, err := HashFromRawBytes(nil)
	assert.Error(t, err)
	hash1 := GenerateTestHash()
	hash2, err := HashFromRawBytes(hash1.RawBytes())
	assert.NoError(t, err)
	require.True(t, hash1.EqualsTo(hash2))

	inv, _ := hex.DecodeString("0102")
	_, err = HashFromRawBytes(inv)
	assert.Error(t, err)
}

func TestHashFromString(t *testing.T) {
	hash1 := GenerateTestHash()
	hash2, err := HashFromString(hash1.String())
	assert.NoError(t, err)
	require.True(t, hash1.EqualsTo(hash2))

	_, err = HashFromString("inv")
	assert.Error(t, err)
}

func TestUndefHash(t *testing.T) {
	h, err := HashFromString("0000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	assert.True(t, h.IsUndef())
	assert.Error(t, h.SanityCheck())
	assert.Equal(t, UndefHash.RawBytes(), h.RawBytes())
}

func TestEmptyHash(t *testing.T) {
	expected, err := HashFromString("0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8")
	assert.NoError(t, err)
	var data = []byte{}
	h := Hash256(data)
	assert.Equal(t, h, expected.RawBytes())
}

func TestHash256(t *testing.T) {
	var data = []byte("zarb")
	h := Hash256(data)
	expected, _ := hex.DecodeString("12b38977f2d67f06f0c0cd54aaf7324cf4fee184398ea33d295e8d1543c2ee1a")
	assert.Equal(t, h, expected)
}

func TestHash160(t *testing.T) {
	var data = []byte("zarb")
	h := Hash160(data)
	expected, _ := hex.DecodeString("e93efc0c83176034cb828e39435eeecc07a29298")
	assert.Equal(t, h, expected)
}
