package sortition

import (
	"math/big"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
)

// evaluate returns a random number between 0 and max with the proof
func evaluate(seed VerifiableSeed, signer crypto.Signer, max int64) (index int64, proof Proof) {
	signData := append(seed[:], signer.PublicKey().RawBytes()...)
	sig := signer.SignData(signData)

	proof, _ = ProofFromRawBytes(sig.RawBytes())
	index = getIndex(proof, max)

	return index, proof
}

// verify ensures the proof is valid
func verify(seed VerifiableSeed, publicKey crypto.PublicKey, proof Proof, max int64) (index int64, result bool) {
	proofSig, err := bls.SignatureFromRawBytes(proof[:])
	if err != nil {
		return 0, false
	}

	// Verify signature (proof)
	signData := append(seed[:], publicKey.RawBytes()...)
	if !publicKey.Verify(signData, proofSig) {
		return 0, false
	}

	index = getIndex(proof, max)

	return index, true
}

func getIndex(proof Proof, max int64) int64 {
	h := hash.CalcHash(proof[:])

	rnd64 := util.SliceToInt64(h.RawBytes())
	rnd64 = rnd64 & 0x7fffffffffffffff

	// construct the numerator and denominator for normalizing the proof uint
	bigRnd := big.NewInt(rnd64)
	bigMax := big.NewInt(max)

	numerator := big.NewInt(0)
	numerator = numerator.Mul(bigRnd, bigMax)

	denominator := big.NewInt(util.MaxInt64) // 0x7FFFFFFFFFFFFFFF

	// divide numerator and denominator to get the election ratio for this block height
	index := big.NewInt(0)
	index = index.Div(numerator, denominator)

	return index.Int64()
}
