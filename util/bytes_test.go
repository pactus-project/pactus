package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyBytes(t *testing.T) {
	tests := []struct {
		bytes []byte
	}{
		{nil},
		{[]byte{}},
		{[]byte{1}},
		{[]byte{1, 2, 3, 4}},
	}

	for _, test := range tests {
		copied := CopyBytes(test.bytes)
		assert.Equal(t, copied, test.bytes)
	}
}

func TestHex2Bytes(t *testing.T) {
	tests := []struct {
		hex   string
		bytes []byte
	}{
		{"", []byte{}},
		{"0x01", []byte{1}},
		{"01", []byte{1}},
		{"1", []byte{1}},
		{"x01", []byte{}},
	}

	for _, test := range tests {
		bytes := Hex2Bytes(test.hex)
		assert.Equal(t, bytes, test.bytes)
	}
}
func TestBytesToHex(t *testing.T) {
	tests := []struct {
		bytes []byte
		hex   string
	}{
		{nil, ""},
		{[]byte{}, ""},
		{[]byte{1}, "01"},
		{[]byte{01}, "01"},
	}

	for _, test := range tests {
		hex := Bytes2Hex(test.bytes)
		assert.Equal(t, hex, test.hex)
	}
}
