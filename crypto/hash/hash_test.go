package hash_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestHashFromString(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	hash1 := ts.RandHash()
	hash2, err := hash.FromString(hash1.String())
	assert.Contains(t, strings.ToUpper(hash1.String()), hash1.ShortString())
	assert.NoError(t, err)
	assert.Equal(t, hash1, hash2)

	_, err = hash.FromString("")
	assert.Error(t, err)

	_, err = hash.FromString("inv")
	assert.Error(t, err)

	_, err = hash.FromString("00")
	assert.Error(t, err)
}

func TestHashEmpty(t *testing.T) {
	_, err := hash.FromBytes(nil)
	assert.Error(t, err)

	_, err = hash.FromBytes([]byte{1})
	assert.Error(t, err)
}

func TestHash256(t *testing.T) {
	data := []byte("zarb")
	h1 := hash.Hash256(data)
	expected, _ := hex.DecodeString("12b38977f2d67f06f0c0cd54aaf7324cf4fee184398ea33d295e8d1543c2ee1a")
	assert.Equal(t, h1, expected)
}

func TestHash160(t *testing.T) {
	data := []byte("zarb")
	h := hash.Hash160(data)
	expected, _ := hex.DecodeString("e93efc0c83176034cb828e39435eeecc07a29298")
	assert.Equal(t, h, expected)
}

func TestHashBasicCheck(t *testing.T) {
	h, err := hash.FromString("0000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	assert.True(t, h.IsUndef())
	assert.Equal(t, hash.UndefHash.Bytes(), h.Bytes())
}
