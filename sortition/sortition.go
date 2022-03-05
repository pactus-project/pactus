package sortition

import (
	"github.com/zarbchain/zarb-go/crypto"
)

func EvaluateSortition(seed VerifiableSeed, signer crypto.Signer, total, threshold int64) (bool, Proof) {
	index, proof := evaluate(seed, signer, total)
	if index < threshold {
		return true, proof
	}

	return false, Proof{}
}

func VerifyProof(seed VerifiableSeed, proof Proof, public crypto.PublicKey, total, threshold int64) bool {
	index, result := verify(seed, public, proof, total)
	if !result {
		return false
	}
	return index < threshold
}
