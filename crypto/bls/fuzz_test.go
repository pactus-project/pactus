package bls_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/stretchr/testify/require"
)

// A public key comes from any peer: decoding any bytes never panics,
// and a key that is decoded is encoded back to the same bytes.
func FuzzPublicKeyFromBytes(f *testing.F) {
	prv, err := bls.KeyGen(make([]byte, 32), nil)
	require.NoError(f, err)
	f.Add(prv.PublicKeyNative().Bytes())
	f.Add([]byte{})
	f.Add(make([]byte, bls.PublicKeySize))

	f.Fuzz(func(t *testing.T, data []byte) {
		pub, err := bls.PublicKeyFromBytes(data)
		if err != nil {
			return
		}
		require.Equal(t, data, pub.Bytes())

		again, err := bls.PublicKeyFromString(pub.String())
		require.NoError(t, err)
		require.True(t, pub.EqualsTo(again))

		_ = pub.AccountAddress()
		_ = pub.ValidatorAddress()
	})
}

// A signature comes from any peer: decoding any bytes never panics,
// and a signature that is decoded is encoded back to the same bytes.
func FuzzSignatureFromBytes(f *testing.F) {
	prv, err := bls.KeyGen(make([]byte, 32), nil)
	require.NoError(f, err)
	f.Add(prv.SignNative([]byte("fuzz")).Bytes())
	f.Add([]byte{})
	f.Add(make([]byte, bls.SignatureSize))

	f.Fuzz(func(t *testing.T, data []byte) {
		sig, err := bls.SignatureFromBytes(data)
		if err != nil {
			return
		}
		require.Equal(t, data, sig.Bytes())

		again, err := bls.SignatureFromString(sig.String())
		require.NoError(t, err)
		require.True(t, sig.EqualsTo(again))
	})
}
