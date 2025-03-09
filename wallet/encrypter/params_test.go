package encrypter

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
	for _, tt := range tests {
		p.SetUint8(tt.key, tt.val)
		assert.Equal(t, tt.val, p.GetUint8(tt.key, 0))
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
	for _, tt := range tests {
		p.SetUint32(tt.key, tt.val)
		assert.Equal(t, tt.val, p.GetUint32(tt.key, 0))
	}
}

func TestParamsUint64(t *testing.T) {
	tests := []struct {
		key string
		val uint64
	}{
		{"k1", uint64(0)},
		{"k2", uint64(0xFFFFFFFFFFFFFFFF)},
	}

	p := params{}
	for _, tt := range tests {
		p.SetUint64(tt.key, tt.val)
		assert.Equal(t, tt.val, p.GetUint64(tt.key, 0))
	}
}

func TestParamsDefaultValue(t *testing.T) {
	p := params{}
	assert.Equal(t, uint64(24), p.GetUint64("not-exist", 24))
	assert.Equal(t, uint32(24), p.GetUint32("not-exist", 24))
	assert.Equal(t, uint8(24), p.GetUint8("not-exist", 24))

}

func TestParamsBytes(t *testing.T) {
	tests := []struct {
		key    string
		val    []byte
		base64 string
	}{
		{"k1", []byte{0, 0}, "AAA="},
		{"k2", []byte{0xff, 0xff}, "//8="},
		{"k2", []byte{}, ""},
	}

	p := params{}
	for _, tt := range tests {
		p.SetBytes(tt.key, tt.val)
		assert.Equal(t, tt.val, p.GetBytes(tt.key))
		assert.Equal(t, tt.base64, p.GetString(tt.key))
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
	for _, tt := range tests {
		p.SetString(tt.key, tt.val)
		assert.Equal(t, tt.val, p.GetString(tt.key))
	}
}
