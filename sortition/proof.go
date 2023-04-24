package sortition

import (
	"encoding/hex"
	"fmt"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util"
)

type Proof [48]byte

func ProofFromString(text string) (Proof, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return Proof{}, err
	}

	return ProofFromBytes(data)
}

func ProofFromBytes(data []byte) (Proof, error) {
	if len(data) != 48 {
		return Proof{}, fmt.Errorf("invalid proof length")
	}

	p := Proof{}
	copy(p[:], data)

	return p, nil
}

func GenerateRandomProof() Proof {
	sig := bls.GenerateTestSigner().SignData(
		util.Int64ToSlice(util.RandInt64(0)))
	proof, _ := ProofFromBytes(sig.Bytes())
	return proof
}
