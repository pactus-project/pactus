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

func serializeKey(t *testing.T, chainCode []byte, reserved uint8, key, trailing []byte) string {
	t.Helper()

	buf := bytes.NewBuffer(nil)
	require.NoError(t, encoding.WriteElement(buf, uint8(0)))
	require.NoError(t, encoding.WriteVarBytes(buf, chainCode))
	require.NoError(t, encoding.WriteElement(buf, reserved))
	require.NoError(t, encoding.WriteVarBytes(buf, key))
	buf.Write(trailing)

	str, err := bech32m.EncodeFromBase256WithType(crypto.XPrivateKeyHRP(), byte(crypto.SignatureTypeEd25519), buf.Bytes())
	require.NoError(t, err)

	return str
}

func TestNewKeyFromStringValidatesLayout(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	chainCode := ts.RandBytes(32)
	key := ts.RandBytes(32)

	tests := []struct {
		name      string
		chainCode []byte
		reserved  uint8
		key       []byte
		trailing  []byte
		valid     bool
	}{
		{"valid", chainCode, 0, key, nil, true},
		{"short chain code", chainCode[:31], 0, key, nil, false},
		{"long chain code", append(append([]byte{}, chainCode...), 0), 0, key, nil, false},
		{"reserved byte set", chainCode, 1, key, nil, false},
		{"short key", chainCode, 0, key[:31], nil, false},
		{"long key", chainCode, 0, append(append([]byte{}, key...), 0), nil, false},
		{"trailing bytes", chainCode, 0, key, []byte{1}, false},
	}

	for _, tt := range tests {
		str := serializeKey(t, tt.chainCode, tt.reserved, tt.key, tt.trailing)
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
