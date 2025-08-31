package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected Version
		hasError bool
	}{
		{"1", ProtocolVersion1, false},
		{"2", ProtocolVersion2, false},
		{"2", ProtocolVersionLatest, false},
		{"invalid", 0, true},
		{"0", 0, false},
		{"-1", Version(255), false},
		{"127", Version(127), false},
		{"128", 0, true}, // out of int8 range
	}

	for _, test := range tests {
		result, err := ParseVersion(test.input)
		if test.hasError {
			assert.Error(t, err, "ParseVersion(%q) should return error", test.input)
		} else {
			assert.NoError(t, err, "ParseVersion(%q) should not return error", test.input)
			assert.Equal(t, test.expected, result, "ParseVersion(%q)", test.input)
		}
	}
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		version  Version
		expected string
	}{
		{ProtocolVersion1, "1"},
		{ProtocolVersion2, "2"},
		{ProtocolVersionLatest, "2"},
		{0, "0"},
		{127, "127"},
		{255, "255"},
	}

	for _, test := range tests {
		result := test.version.String()
		assert.Equal(t, test.expected, result, "Version(%d).String()", test.version)
	}
}
