package sortition

import (
	"github.com/pactus-project/pactus/crypto"
)

func EvaluateSortition(seed VerifiableSeed, signer crypto.Signer, total, threshold int64) (bool, Proof) {
	index, proof := Evaluate(seed, signer, uint64(total))
	if int64(index) < threshold {
		return true, proof
	}

	return false, Proof{}
}

func VerifyProof(seed VerifiableSeed, proof Proof, public crypto.PublicKey, total, threshold int64) bool {
	index, result := Verify(seed, public, proof, uint64(total))
	if !result {
		return false
	}
	return int64(index) < threshold
}
