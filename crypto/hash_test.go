package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUndefHash(t *testing.T) {
	expected, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	assert.Equal(t, UndefHash.RawBytes(), expected)
}

func TestEmptyHash(t *testing.T) {
	expected, _ := hex.DecodeString("0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8")
	var data = []byte{}
	h := Hash256(data)
	assert.Equal(t, h, expected)
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
