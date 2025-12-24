package addresspath

import (
	"strconv"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/stretchr/testify/assert"
)

func TestHardenedKey(t *testing.T) {
	num := uint32(12345)
	hardenedNum := Harden(num)
	unhardenedNum := UnHarden(hardenedNum)

	assert.Equal(t, num, unhardenedNum)
}

func TestPathToString(t *testing.T) {
	h := hardenedKeyStart
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
	for no, tt := range tests {
		assert.Equal(t, tt.wantStr, tt.path.String(), "case %d failed", no)
	}
}

func TestStringToPath(t *testing.T) {
	h := hardenedKeyStart
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
	for no, tt := range tests {
		path, err := FromString(tt.str)
		assert.Equal(t, tt.wantPath, path, "case %d failed", no)
		assert.ErrorIsf(t, err, tt.wantErr, "case %d failed", no)
	}
}

func TestPathHelpers(t *testing.T) {
	path := Path{
		12381 + hardenedKeyStart,
		21888 + hardenedKeyStart,
		2 + hardenedKeyStart,
		1,
	}

	assert.Equal(t, PurposeBLS12381, path.Purpose())
	assert.Equal(t, CoinTypePactusMainnet, path.CoinType())
	assert.Equal(t, crypto.AddressTypeBLSAccount, path.AddressType())
	assert.Equal(t, uint32(1), path.AddressIndex())
}
