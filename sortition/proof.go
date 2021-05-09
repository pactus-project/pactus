package sortition

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type Proof [48]byte

func ProofFromString(text string) (Proof, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return Proof{}, err
	}

	return ProofFromRawBytes(data)
}

func ProofFromRawBytes(data []byte) (Proof, error) {
	if len(data) != 48 {
		return Proof{}, fmt.Errorf("invalid proof length")
	}

	p := Proof{}
	copy(p[:], data)

	return p, nil
}

func (p Proof) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(p[:])), nil
}

func (p *Proof) UnmarshalText(text []byte) error {
	proof, err := ProofFromString(string(text))
	if err != nil {
		return err
	}
	*p = proof
	return nil
}

func GenerateRandomProof() Proof {
	p := Proof{}
	_, err := rand.Read(p[:])
	if err != nil {
		panic(err)
	}
	return p
}
