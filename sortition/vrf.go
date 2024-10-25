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

// Evaluate returns a provable random number between [0, max) along with a proof.
// It returns the random number and the proof that can regenerate the random number using
// the public key of the signer, without revealing the private key.
func Evaluate(seed VerifiableSeed, prv *bls.PrivateKey, max uint64) (uint64, Proof) {
	signData := make([]byte, 0, bls.SignatureSize+bls.PublicKeySize)
	signData = append(signData, seed[:]...)
	signData = append(signData, prv.PublicKey().Bytes()...)

	sig := prv.Sign(signData)

	proof, _ := ProofFromBytes(sig.Bytes())
	index := GetIndex(proof, max)

	return index, proof
}

// Verify checks if the provided proof, based on the seed and public key, is valid.
// If the proof is valid, it calculates the random number that
// can be generated based on the given proof.
func Verify(seed VerifiableSeed, pub *bls.PublicKey, proof Proof, max uint64) (uint64, bool) {
	proofSig, err := bls.SignatureFromBytes(proof[:])
	if err != nil {
		return 0, false
	}

	// Verify signature (proof)
	signData := make([]byte, 0, bls.SignatureSize+bls.PublicKeySize)
	signData = append(signData, seed[:]...)
	signData = append(signData, pub.Bytes()...)
	if err := pub.Verify(signData, proofSig); err != nil {
		return 0, false
	}

	index := GetIndex(proof, max)

	return index, true
}

func GetIndex(proof Proof, max uint64) uint64 {
	hash := hash.CalcHash(proof[:])

	// construct the numerator and denominator for normalizing the proof uint
	bigRnd := &big.Int{}
	bigMax := &big.Int{}
	numerator := &big.Int{}

	bigRnd.SetBytes(hash.Bytes())
	bigMax.SetUint64(max)

	numerator = numerator.Mul(bigRnd, bigMax)

	// divide numerator and denominator to get the election ratio for this block height
	index := big.NewInt(0)
	index = index.Div(numerator, denominator)

	return index.Uint64()
}
