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
	assert.Contains(t, strings.ToUpper(hash1.String()), hash1.LogString())
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
	data := []byte("pactus")
	h1 := hash.Hash256(data)
	expected, _ := hex.DecodeString("ea020ace5c968f755dfc1b5921e574191cd9ff438639badae8a69f667e0d5970")
	assert.Equal(t, expected, h1)
}

func TestHash160(t *testing.T) {
	data := []byte("pactus")
	h := hash.Hash160(data)
	expected, _ := hex.DecodeString("1e4633f850c9590a97633631eae007e18648e30e")
	assert.Equal(t, expected, h)
}

func TestHashBasicCheck(t *testing.T) {
	h, err := hash.FromString("0000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	assert.True(t, h.IsUndef())
	assert.Equal(t, hash.UndefHash.Bytes(), h.Bytes())
}
