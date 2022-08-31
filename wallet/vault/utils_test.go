package vault

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls/hdkeychain"
)

func TestPathToString(t *testing.T) {
	h := hdkeychain.HardenedKeyStart
	tests := []struct {
		path    []uint32
		wantStr string
	}{
		{[]uint32{}, "m"},
		{[]uint32{0}, "m/0"},
		{[]uint32{0, 1}, "m/0/1"},
		{[]uint32{0, 1, 1000000000}, "m/0/1/1000000000"},
		{[]uint32{h}, "m/0'"},
		{[]uint32{h, h + 1}, "m/0'/1'"},
		{[]uint32{h, h + 1, h + 1000000000}, "m/0'/1'/1000000000'"},
	}
	for i, test := range tests {
		assert.Equal(t, derivePathToString(test.path), test.wantStr, "case %d failed", i)
	}
}

func TestStringToPath(t *testing.T) {
	h := hdkeychain.HardenedKeyStart
	tests := []struct {
		str      string
		wantPath []uint32
		wantErr  error
	}{
		{"m", []uint32{}, nil},
		{"m/0", []uint32{0}, nil},
		{"m/0/1", []uint32{0, 1}, nil},
		{"m/0/1/1000000000", []uint32{0, 1, 1000000000}, nil},
		{"m/0'", []uint32{h}, nil},
		{"m/0'/1'", []uint32{h, h + 1}, nil},
		{"m/0'/1'/1000000000'", []uint32{h, h + 1, h + 1000000000}, nil},
		{"i", nil, ErrInvalidPath},
		{"m/'", nil, strconv.ErrSyntax},
		{"m/abc'", nil, strconv.ErrSyntax},
	}
	for i, test := range tests {
		path, err := stringToDerivePath(test.str)
		assert.Equal(t, path, test.wantPath, "case %d failed", i)
		assert.ErrorIsf(t, err, test.wantErr, "case %d failed", i)
	}
}
