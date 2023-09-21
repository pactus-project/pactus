package sortition

import (
	"math/big"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
)

var denominator *big.Int

func init() {
	denominator = &big.Int{}
	denominator.SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
}

// Evaluate returns a random number between 0 and max with the proof.
func Evaluate(seed VerifiableSeed, prv *bls.PrivateKey, max uint64) (index uint64, proof Proof) {
	signData := append(seed[:], prv.PublicKey().Bytes()...)
	sig := prv.Sign(signData)

	proof, _ = ProofFromBytes(sig.Bytes())
	index = GetIndex(proof, max)

	return index, proof
}

// Verify ensures the proof is valid.
func Verify(seed VerifiableSeed, pub *bls.PublicKey, proof Proof, max uint64) (index uint64, result bool) {
	proofSig, err := bls.SignatureFromBytes(proof[:])
	if err != nil {
		return 0, false
	}

	// Verify signature (proof)
	signData := append(seed[:], pub.Bytes()...)
	if err := pub.Verify(signData, proofSig); err != nil {
		return 0, false
	}

	index = GetIndex(proof, max)

	return index, true
}

func GetIndex(proof Proof, max uint64) uint64 {
	h := hash.CalcHash(proof[:])

	// construct the numerator and denominator for normalizing the proof uint
	bigRnd := &big.Int{}
	bigMax := &big.Int{}
	numerator := &big.Int{}

	bigRnd.SetBytes(h.Bytes())
	bigMax.SetUint64(max)

	numerator = numerator.Mul(bigRnd, bigMax)

	// divide numerator and denominator to get the election ratio for this block height
	index := big.NewInt(0)
	index = index.Div(numerator, denominator)

	return index.Uint64()
}
