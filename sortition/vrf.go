package sortition

import (
	"math/big"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

type VRF struct {
	signer crypto.Signer
	max    int64
}

func NewVRF(signer crypto.Signer) *VRF {
	return &VRF{
		signer: signer,
	}
}

func (vrf *VRF) SetMax(max int64) {
	vrf.max = max
}

func (vrf *VRF) AddToMax(num int64) {
	vrf.max += num
}

func (vrf *VRF) Max() int64 {
	return vrf.max
}

// Evaluate returns a random number between 0 and max with the proof
func (vrf *VRF) Evaluate(hash crypto.Hash) (index int64, proof []byte) {
	sig := vrf.signer.Sign(hash.RawBytes())

	proof = sig.RawBytes()
	index = vrf.getIndex(proof)

	return index, proof
}

// Verify ensures the proof is valid
func (vrf *VRF) Verify(hash crypto.Hash, publicKey crypto.PublicKey, proof []byte) (index int64, result bool) {
	sig, err := crypto.SignatureFromRawBytes(proof)
	if err != nil {
		return 0, false
	}

	// Verify signature (proof)
	if !publicKey.Verify(hash.RawBytes(), &sig) {
		return 0, false
	}

	index = vrf.getIndex(sig.RawBytes())

	return index, true
}

func (vrf *VRF) getIndex(sig []byte) int64 {
	h := crypto.HashH(sig)

	rnd64 := util.SliceToInt64(h.RawBytes())
	rnd64 = rnd64 & 0x7fffffffffffffff

	// construct the numerator and denominator for normalizing the signature uint between [0, 1]
	index := big.NewInt(0)
	numerator := big.NewInt(0)
	rnd := big.NewInt(rnd64)
	max := big.NewInt(vrf.max)
	denominator := big.NewInt(util.MaxInt64)

	numerator = numerator.Mul(rnd, max)

	// divide numerator and denominator to get the election ratio for this block height
	index = index.Div(numerator, denominator)

	return index.Int64()
}
