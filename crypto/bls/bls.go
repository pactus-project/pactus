package bls

import (
	bls12381 "github.com/kilic/bls12-381"
)

// set Ciphersuite for Basic mode
// https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-bls-signature-04#section-4.2.1
var dst = []byte("BLS_SIG_BLS12381G1_XMD:SHA-256_SSWU_RO_NUL_")

func Aggregate(sigs []*Signature) *Signature {
	if len(sigs) == 0 {
		return nil
	}
	g1 := bls12381.NewG1()
	aggPointG1 := sigs[0].pointG1
	for i := 1; i < len(sigs); i++ {
		g1.Add(
			&aggPointG1,
			&aggPointG1,
			&sigs[i].pointG1)
	}

	return &Signature{
		pointG1: aggPointG1,
	}
}

func VerifyAggregated(sig *Signature, pubs []*PublicKey, msg []byte) bool {
	if len(pubs) == 0 {
		return false
	}
	g2 := bls12381.NewG2()
	aggPointG2 := pubs[0].pointG2
	for i := 1; i < len(pubs); i++ {
		if g2.IsZero(&pubs[i].pointG2) {
			return false
		}
		g2.Add(
			&aggPointG2,
			&aggPointG2,
			&pubs[i].pointG2)
	}

	aggPub := PublicKey{pointG2: aggPointG2}
	return aggPub.Verify(msg, sig) == nil
}
