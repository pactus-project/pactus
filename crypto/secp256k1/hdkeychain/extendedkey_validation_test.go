package hdkeychain

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func serializeKey(t *testing.T, isPrivate bool, chainCode []byte, reserved uint8, key, trailing []byte) string {
	t.Helper()

	buf := bytes.NewBuffer(nil)
	require.NoError(t, encoding.WriteElement(buf, uint8(0)))
	require.NoError(t, encoding.WriteVarBytes(buf, chainCode))
	require.NoError(t, encoding.WriteElement(buf, reserved))
	require.NoError(t, encoding.WriteVarBytes(buf, key))
	buf.Write(trailing)

	hrp := crypto.XPublicKeyHRP()
	if isPrivate {
		hrp = crypto.XPrivateKeyHRP()
	}
	str, err := bech32m.EncodeFromBase256WithType(hrp, byte(crypto.SignatureTypeSecp256k1), buf.Bytes())
	require.NoError(t, err)

	return str
}

func TestNewKeyFromStringValidatesLayout(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	master, err := NewMaster(ts.RandBytes(32))
	require.NoError(t, err)

	chainCode := ts.RandBytes(32)
	prv, err := master.RawPrivateKey()
	require.NoError(t, err)
	pub := master.Neuter().RawPublicKey()

	tests := []struct {
		name      string
		isPrivate bool
		chainCode []byte
		reserved  uint8
		key       []byte
		trailing  []byte
		valid     bool
	}{
		{"private", true, chainCode, 0, prv, nil, true},
		{"public", false, chainCode, 0, pub, nil, true},
		{"short chain code", true, chainCode[:31], 0, prv, nil, false},
		{"long chain code", true, append(append([]byte{}, chainCode...), 0), 0, prv, nil, false},
		{"reserved byte set", true, chainCode, 1, prv, nil, false},
		{"short private key", true, chainCode, 0, prv[:31], nil, false},
		{"public key as private", true, chainCode, 0, pub, nil, false},
		{"private key as public", false, chainCode, 0, prv, nil, false},
		{"trailing bytes", true, chainCode, 0, prv, []byte{1}, false},
	}

	for _, tt := range tests {
		str := serializeKey(t, tt.isPrivate, tt.chainCode, tt.reserved, tt.key, tt.trailing)
		_, err := NewKeyFromString(str)
		if tt.valid {
			assert.NoError(t, err, tt.name)
		} else {
			assert.ErrorIs(t, err, ErrInvalidKeyData, tt.name)
		}
	}
}

func TestDeriveMaxDepth(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	master, err := NewMaster(ts.RandBytes(32))
	require.NoError(t, err)

	master.path = make([]uint32, maxDepth-1)
	child, err := master.Derive(hardenedKeyStart)
	require.NoError(t, err)
	assert.Len(t, child.Path(), maxDepth)

	_, err = child.Derive(hardenedKeyStart)
	assert.ErrorIs(t, err, ErrMaxDepthExceeded)

	parsed, err := NewKeyFromString(child.String())
	require.NoError(t, err)
	assert.Equal(t, child.Path(), parsed.Path())
}
