package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParamsUint8(t *testing.T) {
	tests := []struct {
		key string
		val uint8
	}{
		{"k1", uint8(0)},
		{"k2", uint8(0xFF)},
	}

	p := params{}
	for _, test := range tests {
		p.setUint8(test.key, test.val)
		assert.Equal(t, test.val, p.getUint8(test.key))
	}
}

func TestParamsUint32(t *testing.T) {
	tests := []struct {
		key string
		val uint32
	}{
		{"k1", uint32(0)},
		{"k2", uint32(0xFFFFFFFF)},
	}

	p := params{}
	for _, test := range tests {
		p.setUint32(test.key, test.val)
		assert.Equal(t, test.val, p.getUint32(test.key))
	}
}

func TestParamsBytes(t *testing.T) {
	tests := []struct {
		key string
		val []byte
	}{
		{"k1", []byte{0, 0}},
		{"k2", []byte{0xff, 0xff}},
	}

	p := params{}
	for _, test := range tests {
		p.setBytes(test.key, test.val)
		assert.Equal(t, test.val, p.getBytes(test.key))
	}
}
