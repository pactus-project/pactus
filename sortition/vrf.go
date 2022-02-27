package sortition

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

// evaluate returns a proof if random number is less than 1
func evaluate(seed VerifiableSeed, signer crypto.Signer, coins, max int) (bool, Proof) {
	sig := signer.SignData(seed).RawBytes()

	for coin := 0; coin < coins; coin++ {
		index, proof := getIndex(sig, coin+1, max)
		if index == 0 {
			return true, proof
		}
	}

	return false, Proof{}
}

// getIndex returns a random number between 0 and totalCoins with the proof
func getIndex(base []byte, coin, totalCoins int) (int, Proof) {
	proof := NewProof(base, coin)
	data, _ := cbor.Marshal(proof)
	fmt.Printf("data: %x\n", data)
	h := hash.Hash256(data)
	fmt.Printf("hash: %x\n", h)

	point := uint64(h[17]) | uint64(h[7])<<8 | uint64(h[8])<<16 | uint64(h[21])<<24
	fmt.Printf("point: %x\n", point)
	index := (uint64(totalCoins) * point) / uint64(0xffffffff)
	fmt.Printf("index: %x\n", index)

	return int(index), proof
}

// verifyProof ensures the proof is valid
func verifyProof(seed VerifiableSeed, public crypto.PublicKey, proof Proof, totalCoin int) bool {
	proofSig, err := bls.SignatureFromRawBytes(proof.Base)
	if err != nil {
		return false
	}

	// Verify signature (base)
	if !public.Verify(seed[:], proofSig) {
		return false
	}

	index, _ := getIndex(proof.Base, proof.Coin, totalCoin)

	return index == 0
}
