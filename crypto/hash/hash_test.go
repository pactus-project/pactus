package hash

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

func TestHashMarshaling(t *testing.T) {
	hash1 := GenerateTestHash()
	hash2 := new(Hash)

	bs, err := cbor.Marshal(hash1)
	assert.NoError(t, err)
	assert.NoError(t, cbor.Unmarshal(bs, hash2))
	assert.True(t, hash1.EqualsTo(*hash2))
	assert.NoError(t, hash1.SanityCheck())

	js, err := hash1.MarshalJSON()
	assert.NoError(t, err)
	assert.Contains(t, string(js), hash1.String())
	assert.Contains(t, strings.ToUpper(hash1.String()), hash1.Fingerprint())
}

func TestHashFromString(t *testing.T) {
	hash1 := GenerateTestHash()
	hash2, err := FromString(hash1.String())
	assert.NoError(t, err)
	assert.True(t, hash1.EqualsTo(hash2))

	_, err = FromString("")
	assert.Error(t, err)

	_, err = FromString("inv")
	assert.Error(t, err)

	_, err = FromString("00")
	assert.Error(t, err)
}

func TestHashEmpty(t *testing.T) {
	hash1 := Hash{}
	assert.Error(t, hash1.SanityCheck())
}

func TestHash256(t *testing.T) {
	var data = []byte("zarb")
	h := Hash256(data)
	expected, _ := hex.DecodeString("12b38977f2d67f06f0c0cd54aaf7324cf4fee184398ea33d295e8d1543c2ee1a")
	assert.Equal(t, h, expected)

	hash, _ := FromBytes(h)
	stamp, _ := StampFromString("12b38977")
	assert.Equal(t, hash.Stamp(), stamp)
}

func TestHash160(t *testing.T) {
	var data = []byte("zarb")
	h := Hash160(data)
	expected, _ := hex.DecodeString("e93efc0c83176034cb828e39435eeecc07a29298")
	assert.Equal(t, h, expected)
}

func TestHashSanityCheck(t *testing.T) {
	h, err := FromString("0000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	assert.True(t, h.IsUndef())
	assert.Error(t, h.SanityCheck())
	assert.Equal(t, UndefHash.RawBytes(), h.RawBytes())
}
