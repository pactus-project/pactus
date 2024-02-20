package addresspath

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathToString(t *testing.T) {
	h := HardenedKeyStart
	tests := []struct {
		path    Path
		wantStr string
	}{
		{NewPath(), "m"},
		{NewPath(0), "m/0"},
		{NewPath(0, 1), "m/0/1"},
		{NewPath(0, 1, 1000000000), "m/0/1/1000000000"},
		{NewPath(h), "m/0'"},
		{NewPath(h, h+1), "m/0'/1'"},
		{NewPath(h, h+1, h+1000000000), "m/0'/1'/1000000000'"},
	}
	for i, test := range tests {
		assert.Equal(t, test.path.String(), test.wantStr, "case %d failed", i)
	}
}

func TestStringToPath(t *testing.T) {
	h := HardenedKeyStart
	tests := []struct {
		str      string
		wantPath Path
		wantErr  error
	}{
		{"m", nil, nil},
		{"m/0", Path{0}, nil},
		{"m/0/1", Path{0, 1}, nil},
		{"m/0/1/1000000000", Path{0, 1, 1000000000}, nil},
		{"m/0'", Path{h}, nil},
		{"m/0'/1'", Path{h, h + 1}, nil},
		{"m/0'/1'/1000000000'", Path{h, h + 1, h + 1000000000}, nil},
		{"i", nil, ErrInvalidPath},
		{"m/'", nil, strconv.ErrSyntax},
		{"m/abc'", nil, strconv.ErrSyntax},
	}
	for i, test := range tests {
		path, err := FromString(test.str)
		assert.Equal(t, path, test.wantPath, "case %d failed", i)
		assert.ErrorIsf(t, err, test.wantErr, "case %d failed", i)
	}
}

func TestPathHelpers(t *testing.T) {
	purpose, coinType, addressType, addressIndex := 12381, 21888, 2, 0
	path := Path{12381, 21888, 2, 0}

	assert.Equal(t, uint32(purpose), path.Purpose())
	assert.Equal(t, uint32(coinType), path.CoinType())
	assert.Equal(t, uint32(addressType), path.AddressType())
	assert.Equal(t, uint32(addressIndex), path.AddressIndex())
}
