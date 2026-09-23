package protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected Version
		hasError bool
	}{
		{"1", ProtocolVersion1, false},
		{"2", ProtocolVersion2, false},
		{"3", ProtocolVersion3, false},
		{"4", ProtocolVersion4, false},
		{"5", ProtocolVersion5, false},
		{"5", ProtocolVersionLatest, false},
		{"6", Version(6), false},
		{"invalid", 0, true},
		{"0", 0, false},
		{"-1", Version(255), false},
		{"127", Version(127), false},
		{"128", 0, true}, // out of int8 range
	}

	for _, test := range tests {
		result, err := ParseVersion(test.input)
		if test.hasError {
			require.Error(t, err, "ParseVersion(%q) should return error", test.input)
		} else {
			require.NoError(t, err, "ParseVersion(%q) should not return error", test.input)
			assert.Equal(t, test.expected, result, "ParseVersion(%q)", test.input)
		}
	}

	four, err := ParseVersion("4")
	require.NoError(t, err)
	assert.Equal(t, ProtocolVersion4, four)
	assert.NotEqual(t, ProtocolVersionLatest, four)

	six, err := ParseVersion("6")
	require.NoError(t, err)
	assert.NotEqual(t, ProtocolVersionLatest, six)
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		version  Version
		expected string
	}{
		{ProtocolVersion1, "1"},
		{ProtocolVersion2, "2"},
		{ProtocolVersion3, "3"},
		{ProtocolVersion4, "4"},
		{ProtocolVersion5, "5"},
		{ProtocolVersionLatest, "5"},
		{0, "0"},
		{127, "127"},
		{255, "255"},
	}

	for _, test := range tests {
		result := test.version.String()
		assert.Equal(t, test.expected, result, "Version(%d).String()", test.version)
	}
}
