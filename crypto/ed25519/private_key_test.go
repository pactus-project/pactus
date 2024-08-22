package ed25519_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv1 := ts.RandED25519KeyPair()
	_, prv2 := ts.RandED25519KeyPair()

	fmt.Println(prv1.String())

	assert.True(t, prv1.EqualsTo(prv1))
	assert.False(t, prv1.EqualsTo(prv2))
	assert.Equal(t, prv1, prv1)
	assert.NotEqual(t, prv1, prv2)
}

func TestPrivateKeyFromString(t *testing.T) {
	tests := []struct {
		errMsg  string
		encoded string
		valid   bool
		result  []byte
	}{
		{
			"invalid separator index -1",
			"XXXXXXR2SYCC5TDQKMJ73J64J8GJTMTKREEQNQAS0M5SLZ9LVJV7Y940NVZQD9JUS" +
				"GV2N44C9H5PVGRXARNGZ7QF3PSKH7805E5SZXPE7ZHHAGX0NFQR",
			false, nil,
		},
		{
			"invalid checksum (expected s9c56g got czlgh0)",
			"SECRET1RAC7048K666DCCYG7FJW68ZE2G6P32UAPLRLWDV3RTAR4PWZUX2CDSFAL55VM" +
				"YS06CY35LA72AWZN5DY5NZA078S4S4K654UFJ0YCCZLGH0",
			false, nil,
		},
		{
			"invalid bech32 string length 0",
			"",
			false, nil,
		},
		{
			"invalid character not part of charset: 105",
			"SECRET1IOIOOI",
			false, nil,
		},
		{
			"invalid bech32 string length 0",
			"SECRET1HPZZU9",
			false, nil,
		},
		{
			"",
			"SECRET1RJ6STNTA7Y3P2QLQF8A6QCX05F2H5TFNE5RSH066KZME4WVFXKE7QW097LG",
			true,
			[]byte{
				0x96, 0xa0, 0xb9, 0xaf, 0xbe, 0x24, 0x42, 0xa0, 0x7c, 0x9, 0x3f, 0x74, 0xc, 0x19, 0xf4,
				0x4a, 0xaf, 0x45, 0xa6, 0x79, 0xa0, 0xe1, 0x77, 0xeb, 0x56, 0x16, 0xf3, 0x57, 0x31, 0x26,
				0xb6, 0x7c,
			},
		},
	}

	for no, test := range tests {
		prv, err := ed25519.PrivateKeyFromString(test.encoded)
		if test.valid {
			assert.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, test.result, prv.Bytes(), "test %v: invalid bytes", no)
			assert.Equal(t, strings.ToUpper(test.encoded), prv.String(), "test %v: invalid encoded", no)
		} else {
			assert.Contains(t, err.Error(), test.errMsg, "test %v: error not matched", no)
		}
	}
}

// TestKeyGen ensures the KeyGen function works as intended.
func TestKeyGen(t *testing.T) {
	tests := []struct {
		seed []byte
		sk   string
	}{}

	for i, test := range tests {
		prv, err := ed25519.KeyGen(test.seed)
		if test.sk == "Err" {
			assert.Error(t, err,
				"test '%v' failed. no error", i)
		} else {
			assert.NoError(t, err,
				"test'%v' failed. has error", i)
			assert.Equal(t, test.sk, hex.EncodeToString(prv.Bytes()),
				"test '%v' failed. not equal", i)
		}
	}
}
