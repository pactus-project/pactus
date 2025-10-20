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
// signature. It returns nil if no signatures are provided or if the first
// signature fails to decode to a valid point.
func SignatureAggregate(sigs ...*Signature) *Signature {
	if len(sigs) == 0 {
		return nil
	}
	grp1 := new(bls12381.G1Affine)
	aggPointG1, err := sigs[0].PointG1()
	if err != nil {
		return nil
	}
	for i := 1; i < len(sigs); i++ {
		pointG1, _ := sigs[i].PointG1()
		aggPointG1 = grp1.Add(aggPointG1, pointG1)
	}

	data := aggPointG1.Bytes()

	return &Signature{
		data:    data[:],
		pointG1: aggPointG1,
	}
}

// PublicKeyAggregate aggregates one or more BLS public keys into a single
// public key. It returns nil if no public keys are provided or if the first
// public key fails to decode to a valid point.
func PublicKeyAggregate(pubs ...*PublicKey) *PublicKey {
	if len(pubs) == 0 {
		return nil
	}
	grp2 := new(bls12381.G2Affine)
	aggPointG2, err := pubs[0].PointG2()
	if err != nil {
		return nil
	}
	for i := 1; i < len(pubs); i++ {
		pointG2, _ := pubs[i].PointG2()
		aggPointG2 = grp2.Add(aggPointG2, pointG2)
	}

	data := aggPointG2.Bytes()

	return &PublicKey{
		data:    data[:],
		pointG2: aggPointG2,
	}
}
