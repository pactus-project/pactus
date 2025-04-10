package version_test

import (
	"fmt"
	"testing"

	"github.com/pactus-project/pactus/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionString(t *testing.T) {
	tests := []struct {
		name     string
		version  version.Version
		expected string
	}{
		{
			name:     "Version with Meta",
			version:  version.Version{Major: 1, Minor: 2, Patch: 3, Meta: "beta"},
			expected: "1.2.3-beta",
		},
		{
			name:     "Version without Meta",
			version:  version.Version{Major: 2, Minor: 0, Patch: 0, Meta: ""},
			expected: "2.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    version.Version
		expectedErr bool
	}{
		{
			name:        "Valid Version String starts with 'v'",
			input:       "v1.2.3",
			expected:    version.Version{Major: 1, Minor: 2, Patch: 3, Meta: ""},
			expectedErr: false,
		},
		{
			name:        "Valid Version String with Meta",
			input:       "1.2.3-beta",
			expected:    version.Version{Major: 1, Minor: 2, Patch: 3, Meta: "beta"},
			expectedErr: false,
		},
		{
			name:        "Valid Version String without Meta",
			input:       "2.0.0",
			expected:    version.Version{Major: 2, Minor: 0, Patch: 0, Meta: ""},
			expectedErr: false,
		},
		{
			name:        "Invalid Version String",
			input:       "1.2",
			expected:    version.Version{},
			expectedErr: true,
		},
		{
			name:        "Invalid Major number",
			input:       "one.2.3",
			expected:    version.Version{},
			expectedErr: true,
		},
		{
			name:        "Invalid Minor number",
			input:       "1.two.3",
			expected:    version.Version{},
			expectedErr: true,
		},
		{
			name:        "Invalid Patch number",
			input:       "1.2.three",
			expected:    version.Version{},
			expectedErr: true,
		},
		{
			name:        "Invalid Patch-Meta Format",
			input:       "1.2.3-rc1-dev",
			expected:    version.Version{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ver, err := version.ParseVersion(tt.input)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, ver)
			}
		})
	}
}

func TestVersionComparison(t *testing.T) {
	tests := []struct {
		name         string
		v1Input      string
		v2Input      string
		expectedSign int
	}{
		{
			name:         "Equal Versions",
			v1Input:      "1.2.3",
			v2Input:      "1.2.3",
			expectedSign: 0,
		},
		{
			name:         "Equal Versions",
			v1Input:      "1.2.3-beta",
			v2Input:      "1.2.3",
			expectedSign: 0,
		},
		{
			name:         "Equal Versions",
			v1Input:      "1.2.3-beta",
			v2Input:      "1.2.3-rc1",
			expectedSign: 0,
		},
		{
			name:         "Lesser Patch Version",
			v1Input:      "1.2.2",
			v2Input:      "1.2.3",
			expectedSign: -1,
		},
		{
			name:         "Lesser Minor Version",
			v1Input:      "2.1.0",
			v2Input:      "2.2.0",
			expectedSign: -1,
		},
		{
			name:         "Greater Major Version",
			v1Input:      "2.1.0",
			v2Input:      "1.2.3",
			expectedSign: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ver1, err := version.ParseVersion(tt.v1Input)
			assert.NoError(t, err)

			ver2, err := version.ParseVersion(tt.v2Input)
			assert.NoError(t, err)

			expectedSign := tt.expectedSign
			actualSign := ver1.Compare(ver2)

			assert.Equal(t, expectedSign, actualSign,
				fmt.Sprintf("Comparison result mismatch for %s vs %s", tt.v1Input, tt.v2Input))
		})
	}
}

// TestCheckVersionString checks if the current version string is valid and parsable.
func TestCheckVersionString(t *testing.T) {
	curVer := version.NodeVersion()
	parsedVer, err := version.ParseVersion(curVer.String())
	require.NoError(t, err)
	assert.Equal(t, curVer.Major, parsedVer.Major)
	assert.Equal(t, curVer.Minor, parsedVer.Minor)
	assert.Equal(t, curVer.Patch, parsedVer.Patch)
	assert.Equal(t, curVer.Meta, parsedVer.Meta)
}
