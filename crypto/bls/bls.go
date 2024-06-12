package bls

import (
	bls12381 "github.com/kilic/bls12-381"
)

// set Ciphersuite for Basic mode
// https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-bls-signature-04#section-4.2.1
var dst = []byte("BLS_SIG_BLS12381G1_XMD:SHA-256_SSWU_RO_NUL_")

func SignatureAggregate(sigs ...*Signature) *Signature {
	if len(sigs) == 0 {
		return nil
	}
	g1 := bls12381.NewG1()
	aggPointG1, err := sigs[0].PointG1()
	if err != nil {
		return nil
	}
	for i := 1; i < len(sigs); i++ {
		s, err := sigs[i].PointG1()
		if err != nil {
			return nil
		}
		g1.Add(
			&aggPointG1,
			&aggPointG1,
			&s)
	}

	data := g1.ToCompressed(&aggPointG1)

	return &Signature{
		data:    data,
		pointG1: &aggPointG1,
	}
}

func PublicKeyAggregate(pubs ...*PublicKey) *PublicKey {
	if len(pubs) == 0 {
		return nil
	}
	g2 := bls12381.NewG2()
	aggPointG2, err := pubs[0].PointG2()
	if err != nil {
		return nil
	}
	for i := 1; i < len(pubs); i++ {
		pointG2, _ := pubs[i].PointG2()
		g2.Add(
			&aggPointG2,
			&aggPointG2,
			&pointG2)
	}

	data := g2.ToCompressed(&aggPointG2)

	return &PublicKey{
		data:    data,
		pointG2: &aggPointG2,
	}
}
