package sortition

import (
	"encoding/hex"
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

type VerifiableSeed [48]byte

func VerifiableSeedFromString(text string) (VerifiableSeed, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return VerifiableSeed{}, err
	}

	return VerifiableSeedFromRawBytes(data)
}

func VerifiableSeedFromRawBytes(data []byte) (VerifiableSeed, error) {
	if len(data) != 48 {
		return VerifiableSeed{}, fmt.Errorf("invalid seed length")
	}

	s := VerifiableSeed{}
	copy(s[:], data)

	return s, nil
}

func (s *VerifiableSeed) Generate(signer crypto.Signer) VerifiableSeed {
	hash := hash.CalcHash(s[:])
	sig := signer.SignData(hash.RawBytes())
	newSeed, _ := VerifiableSeedFromRawBytes(sig.RawBytes())
	return newSeed
}

func (s *VerifiableSeed) Verify(public crypto.PublicKey, prevSeed VerifiableSeed) bool {
	sig, _ := bls.SignatureFromRawBytes(s[:])
	hash := hash.CalcHash(prevSeed[:])
	return public.Verify(hash.RawBytes(), sig)
}

func (s VerifiableSeed) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(s[:])), nil
}

func (s *VerifiableSeed) UnmarshalText(text []byte) error {
	seed, err := VerifiableSeedFromString(string(text))
	if err != nil {
		return err
	}
	*s = seed
	return nil
}

func GenerateRandomSeed() VerifiableSeed {
	h := hash.GenerateTestHash()
	signer := bls.GenerateTestSigner()
	sig := signer.SignData(h.RawBytes())
	seed, _ := VerifiableSeedFromRawBytes(sig.RawBytes())
	return seed
}
