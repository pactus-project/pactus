// Package bls implements BLS signatures over the BLS12-381 pairing-friendly curve.
//
// This package uses the gnark-crypto BLS12-381 implementation and provides
// the main primitives used across Pactus: PrivateKey, PublicKey, and
// Signature. It also exposes helpers for aggregating signatures and public
// keys when multiple participants are involved.
//
// The ciphersuite/domain separation follows `BLS_SIG_BLS12381G1_XMD:SHA-256_SSWU_RO_NUL_`.
package bls

import (
	"errors"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
)

// Ciphersuite for basic mode as defined in:
// https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-bls-signature-04#section-4.2.1
var (
	dst     = []byte("BLS_SIG_BLS12381G1_XMD:SHA-256_SSWU_RO_NUL_")
	gen2Aff bls12381.G2Affine
	gen2Jac bls12381.G2Jac
)

func init() {
	_, gen2Jac, _, gen2Aff = bls12381.Generators()
}

// SignatureAggregate aggregates one or more BLS signatures into a single
// signature. It returns an error if no signatures are provided or if any
// signature fails to decode to a valid point.
func SignatureAggregate(sigs ...*Signature) (*Signature, error) {
	if len(sigs) == 0 {
		return nil, errors.New("no signatures provided")
	}
	firstPointG1, err := sigs[0].PointG1()
	if err != nil {
		return nil, err
	}

	// The points are added in Jacobian coordinates and converted back once:
	// an addition in affine coordinates needs a field inversion each time.
	var aggJac bls12381.G1Jac
	aggJac.FromAffine(firstPointG1)
	for i := 1; i < len(sigs); i++ {
		pointG1, err := sigs[i].PointG1()
		if err != nil {
			return nil, err
		}
		aggJac.AddMixed(pointG1)
	}
	aggPointG1 := new(bls12381.G1Affine).FromJacobian(&aggJac)

	data := aggPointG1.Bytes()

	sig := &Signature{data: data[:]}
	sig.pointG1.Store(aggPointG1)

	return sig, nil
}

// PublicKeyAggregate aggregates one or more BLS public keys into a single
// public key. It returns an error if no public keys are provided or if any
// public key fails to decode to a valid point.
func PublicKeyAggregate(pubs ...*PublicKey) (*PublicKey, error) {
	if len(pubs) == 0 {
		return nil, errors.New("no public keys provided")
	}
	firstPointG2, err := pubs[0].PointG2()
	if err != nil {
		return nil, err
	}

	// The points are added in Jacobian coordinates and converted back once:
	// an addition in affine coordinates needs a field inversion each time.
	var aggJac bls12381.G2Jac
	aggJac.FromAffine(firstPointG2)
	for i := 1; i < len(pubs); i++ {
		pointG2, err := pubs[i].PointG2()
		if err != nil {
			return nil, err
		}
		aggJac.AddMixed(pointG2)
	}
	aggPointG2 := new(bls12381.G2Affine).FromJacobian(&aggJac)

	data := aggPointG2.Bytes()

	pub := &PublicKey{data: data[:]}
	pub.pointG2.Store(aggPointG2)

	return pub, nil
}
