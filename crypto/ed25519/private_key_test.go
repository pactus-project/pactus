package ed25519_test

import (
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	_, prv1 := ts.RandEd25519KeyPair()
	_, prv2 := ts.RandEd25519KeyPair()
	_, prv3 := ts.RandBLSKeyPair()

	assert.True(t, prv1.EqualsTo(prv1))
	assert.False(t, prv1.EqualsTo(prv2))
	assert.False(t, prv1.EqualsTo(prv3))
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
			"not_proper_encoded",
			false, nil,
		},
		{
			"invalid checksum (expected uuk3y0 got uuk30y)",
			"SECRET1RYY62A96X25ZAL4DPL5Z63G83GCSFCCQ7K0CMQD3MFNLYK3A6R26QUUK30Y",
			false, nil,
		},
		{
			"invalid HRP: xxx",
			"XXX1RYY62A96X25ZAL4DPL5Z63G83GCSFCCQ7K0CMQD3MFNLYK3A6R26Q8JXUV6",
			false, nil,
		},
		{
			"invalid signature type: 4",
			"SECRET1YVKPE43FDU9TC4C8LPFD4JY9METET3GEKQE7E7ECK4EJYV20WVAPQZCU0KL",
			false, nil,
		},
		{
			"invalid length: 31",
			"SECRET1RDRWTLP5PX0FAHDX39GXZJP7FKZFALML0D5U9TT9KVQHDUC99CCPV3HNE",
			false, nil,
		},
		{
			"",
			"SECRET1RJ6STNTA7Y3P2QLQF8A6QCX05F2H5TFNE5RSH066KZME4WVFXKE7QW097LG",
			true,
			[]byte{
				0x96, 0xa0, 0xb9, 0xaf, 0xbe, 0x24, 0x42, 0xa0, 0x7c, 0x09, 0x3f, 0x74, 0x0c, 0x19, 0xf4, 0x4a,
				0xaf, 0x45, 0xa6, 0x79, 0xa0, 0xe1, 0x77, 0xeb, 0x56, 0x16, 0xf3, 0x57, 0x31, 0x26, 0xb6, 0x7c,
			},
		},
	}

	for no, tt := range tests {
		prv, err := ed25519.PrivateKeyFromString(tt.encoded)
		if tt.valid {
			assert.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, tt.result, prv.Bytes(), "test %v: invalid bytes", no)
			assert.Equal(t, strings.ToUpper(tt.encoded), prv.String(), "test %v: invalid encoded", no)
		} else {
			assert.Contains(t, err.Error(), tt.errMsg, "test %v: error not matched", no)
		}
	}
}
