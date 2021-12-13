package sortition

import (
	"encoding/hex"
	"fmt"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

type Seed [48]byte

func SeedFromString(text string) (Seed, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return Seed{}, err
	}

	return SeedFromRawBytes(data)
}

func SeedFromRawBytes(data []byte) (Seed, error) {
	if len(data) != 48 {
		return Seed{}, fmt.Errorf("invalid seed length")
	}

	s := Seed{}
	copy(s[:], data)

	return s, nil
}

func (s *Seed) Generate(signer crypto.Signer) Seed {
	hash := hash.CalcHash(s[:])
	sig := signer.SignData(hash.RawBytes())
	newSeed, _ := SeedFromRawBytes(sig.RawBytes())
	return newSeed
}

func (s *Seed) Validate(public crypto.PublicKey, prevSeed Seed) bool {
	sig, _ := bls.SignatureFromRawBytes(s[:])
	hash := hash.CalcHash(prevSeed[:])
	return public.Verify(hash.RawBytes(), sig)
}

func (s Seed) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(s[:])), nil
}

func (s *Seed) UnmarshalText(text []byte) error {
	seed, err := SeedFromString(string(text))
	if err != nil {
		return err
	}
	*s = seed
	return nil
}

func GenerateRandomSeed() Seed {
	h := hash.GenerateTestHash()
	signer := bls.GenerateTestSigner()
	sig := signer.SignData(h.RawBytes())
	seed, _ := SeedFromRawBytes(sig.RawBytes())
	return seed
}
