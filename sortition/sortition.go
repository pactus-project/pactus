package sortition

import (
	"github.com/pactus-project/pactus/crypto/bls"
)

func EvaluateSortition(seed VerifiableSeed, prv *bls.PrivateKey, total, threshold int64) (bool, Proof) {
	index, proof := Evaluate(seed, prv, uint64(total))
	if int64(index) < threshold {
		return true, proof
	}

	return false, Proof{}
}

func VerifyProof(seed VerifiableSeed, proof Proof, pub *bls.PublicKey, total, threshold int64) bool {
	index, result := Verify(seed, pub, proof, uint64(total))
	if !result {
		return false
	}

	return int64(index) < threshold
}
