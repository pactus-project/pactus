package vault

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
		p.SetUint8(test.key, test.val)
		assert.Equal(t, test.val, p.GetUint8(test.key))
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
		p.SetUint32(test.key, test.val)
		assert.Equal(t, test.val, p.GetUint32(test.key))
	}
}

func TestParamsBytes(t *testing.T) {
	tests := []struct {
		key string
		val []byte
	}{
		{"k1", []byte{0, 0}},
		{"k2", []byte{0xff, 0xff}},
		{"k2", []byte{}},
	}

	p := params{}
	for _, test := range tests {
		p.SetBytes(test.key, test.val)
		assert.Equal(t, test.val, p.GetBytes(test.key))
	}
}

func TestParamsString(t *testing.T) {
	tests := []struct {
		key string
		val string
	}{
		{"k1", "foo"},
		{"k2", "bar"},
		{"k3", "bar"},
	}

	p := params{}
	for _, test := range tests {
		p.SetString(test.key, test.val)
		assert.Equal(t, test.val, p.GetString(test.key))
	}
}
