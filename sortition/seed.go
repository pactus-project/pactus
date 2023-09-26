package sortition

import (
	"encoding/hex"
	"fmt"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
)

type VerifiableSeed [48]byte

var UndefVerifiableSeed = VerifiableSeed{}

func VerifiableSeedFromString(text string) (VerifiableSeed, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return UndefVerifiableSeed, err
	}

	return VerifiableSeedFromBytes(data)
}

func VerifiableSeedFromBytes(data []byte) (VerifiableSeed, error) {
	if len(data) != 48 {
		return UndefVerifiableSeed, fmt.Errorf("invalid seed length")
	}

	s := UndefVerifiableSeed
	copy(s[:], data)

	return s, nil
}

func (s *VerifiableSeed) GenerateNext(prv *bls.PrivateKey) VerifiableSeed {
	hash := hash.CalcHash(s[:])
	sig := prv.Sign(hash.Bytes())
	newSeed, _ := VerifiableSeedFromBytes(sig.Bytes())
	return newSeed
}

func (s *VerifiableSeed) Verify(public *bls.PublicKey, prevSeed VerifiableSeed) bool {
	sig, err := bls.SignatureFromBytes(s[:])
	if err != nil {
		return false
	}
	hash := hash.CalcHash(prevSeed[:])
	return public.Verify(hash.Bytes(), sig) == nil
}
