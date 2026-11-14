package bls_test

import (
	"testing"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fp"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/stretchr/testify/require"
)

// The points of the signatures and the public keys must be in the prime-order subgroup.
// The decoding of the points is the only place that checks it, so these tests guard it:
// they fail if a new version of the curve library stops checking the subgroup by default.

func TestSignaturePointNotInSubgroup(t *testing.T) {
	// A point of the curve is almost never in the subgroup: the cofactor is large.
	// Find the first x that has a point: y^2 = x^3 + 4.
	var point bls12381.G1Affine
	var four fp.Element
	four.SetUint64(4)
	for i := uint64(1); ; i++ {
		var ySquare fp.Element
		point.X.SetUint64(i)
		ySquare.Square(&point.X).Mul(&ySquare, &point.X).Add(&ySquare, &four)
		if point.Y.Sqrt(&ySquare) != nil {
			break
		}
	}
	require.True(t, point.IsOnCurve())
	require.False(t, point.IsInSubGroup())

	data := point.Bytes()
	sig, err := bls.SignatureFromBytes(data[:])
	require.NoError(t, err)

	_, err = sig.PointG1()
	require.Error(t, err)

	prv, err := bls.KeyGen([]byte("0123456789abcdef0123456789abcdef"), nil)
	require.NoError(t, err)
	require.Error(t, prv.PublicKeyNative().Verify([]byte("msg"), sig))
}

func TestPublicKeyPointNotInSubgroup(t *testing.T) {
	// Find the first x that has a point on the twist: y^2 = x^3 + 4(1+u).
	var point bls12381.G2Affine
	var coeff bls12381.E2
	coeff.A0.SetUint64(4)
	coeff.A1.SetUint64(4)
	for i := uint64(1); ; i++ {
		var ySquare bls12381.E2
		point.X.A0.SetUint64(i)
		point.X.A1.SetUint64(1)
		ySquare.Square(&point.X).Mul(&ySquare, &point.X).Add(&ySquare, &coeff)
		if ySquare.Legendre() == 1 {
			point.Y.Sqrt(&ySquare)

			break
		}
	}
	require.True(t, point.IsOnCurve())
	require.False(t, point.IsInSubGroup())

	data := point.Bytes()
	pub, err := bls.PublicKeyFromBytes(data[:])
	if err == nil {
		_, err = pub.PointG2()
	}
	require.Error(t, err)
}

func TestPointAtInfinityIsRejected(t *testing.T) {
	var g1 bls12381.G1Affine
	g1Data := g1.Bytes()
	sig, err := bls.SignatureFromBytes(g1Data[:])
	require.NoError(t, err)
	_, err = sig.PointG1()
	require.ErrorIs(t, err, crypto.ErrInvalidSignature)

	var g2 bls12381.G2Affine
	g2Data := g2.Bytes()
	pub, err := bls.PublicKeyFromBytes(g2Data[:])
	if err == nil {
		_, err = pub.PointG2()
	}
	require.ErrorIs(t, err, crypto.ErrInvalidPublicKey)
}
