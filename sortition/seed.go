package sortition

import (
	"encoding/hex"
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
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

func (s *VerifiableSeed) Generate(signer crypto.Signer) VerifiableSeed {
	hash := hash.CalcHash(s[:])
	sig := signer.SignData(hash.Bytes())
	newSeed, _ := VerifiableSeedFromBytes(sig.Bytes())
	return newSeed
}

func (s *VerifiableSeed) Verify(public crypto.PublicKey, prevSeed VerifiableSeed) bool {
	sig, err := bls.SignatureFromBytes(s[:])
	if err != nil {
		return false
	}
	if err := sig.SanityCheck(); err != nil {
		return false
	}
	hash := hash.CalcHash(prevSeed[:])
	return public.Verify(hash.Bytes(), sig)
}

func GenerateRandomSeed() VerifiableSeed {
	h := hash.GenerateTestHash()
	signer := bls.GenerateTestSigner()
	sig := signer.SignData(h.Bytes())
	seed, _ := VerifiableSeedFromBytes(sig.Bytes())
	return seed
}
