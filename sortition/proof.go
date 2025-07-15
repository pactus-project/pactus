package sortition

import (
	"encoding/hex"
	"errors"
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
		return Proof{}, errors.New("invalid proof length")
	}

	p := Proof{}
	copy(p[:], data)

	return p, nil
}
