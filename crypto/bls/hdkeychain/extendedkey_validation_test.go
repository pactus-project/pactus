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

func serializeKey(t *testing.T, isPrivate, pubOnG1 bool, chainCode, key, trailing []byte) string {
	t.Helper()

	buf := bytes.NewBuffer(nil)
	require.NoError(t, encoding.WriteElement(buf, uint8(0)))
	require.NoError(t, encoding.WriteVarBytes(buf, chainCode))
	require.NoError(t, encoding.WriteElement(buf, pubOnG1))
	require.NoError(t, encoding.WriteVarBytes(buf, key))
	buf.Write(trailing)

	hrp := crypto.XPublicKeyHRP()
	if isPrivate {
		hrp = crypto.XPrivateKeyHRP()
	}
	str, err := bech32m.EncodeFromBase256WithType(hrp, byte(crypto.SignatureTypeBLS), buf.Bytes())
	require.NoError(t, err)

	return str
}

func TestNewKeyFromStringValidatesLayout(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	chainCode := ts.RandBytes(32)

	for _, pubOnG1 := range []bool{true, false} {
		master, err := NewMaster(ts.RandBytes(32), pubOnG1)
		require.NoError(t, err)

		prv, err := master.RawPrivateKey()
		require.NoError(t, err)
		pub := master.Neuter().RawPublicKey()

		tests := []struct {
			name      string
			isPrivate bool
			chainCode []byte
			key       []byte
			trailing  []byte
			valid     bool
		}{
			{"private", true, chainCode, prv, nil, true},
			{"public", false, chainCode, pub, nil, true},
			{"short chain code", true, chainCode[:31], prv, nil, false},
			{"long chain code", true, append(append([]byte{}, chainCode...), 0), prv, nil, false},
			{"short private key", true, chainCode, prv[:31], nil, false},
			{"public key as private", true, chainCode, pub, nil, false},
			{"private key as public", false, chainCode, prv, nil, false},
			{"short public key", false, chainCode, pub[:len(pub)-1], nil, false},
			{"trailing bytes", true, chainCode, prv, []byte{1}, false},
		}

		for _, tt := range tests {
			str := serializeKey(t, tt.isPrivate, pubOnG1, tt.chainCode, tt.key, tt.trailing)
			_, err := NewKeyFromString(str)
			if tt.valid {
				assert.NoError(t, err, tt.name)
			} else {
				assert.ErrorIs(t, err, ErrInvalidKeyData, tt.name)
			}
		}
	}
}

func TestDeriveMaxDepth(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	master, err := NewMaster(ts.RandBytes(32), false)
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
