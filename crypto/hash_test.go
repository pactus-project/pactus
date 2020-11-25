package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash256(t *testing.T) {
	var data = []byte("zarb")
	h := Hash256(data)
	expected, _ := hex.DecodeString("d68b7866ebb7412bb398517c7a990663b00d8cca4159ab7e5620a19f0fb2fb8e")
	assert.Equal(t, h, expected)
}

func TestHash160(t *testing.T) {
	var data = []byte("zarb")
	h := Hash160(data)
	expected, _ := hex.DecodeString("e93efc0c83176034cb828e39435eeecc07a29298")
	assert.Equal(t, h, expected)
}
