package sortition

import (
	"math/big"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

type VRF struct {
	max int64
}

func NewVRF() *VRF {
	return &VRF{}
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
func (vrf *VRF) Evaluate(seed Seed, signer crypto.Signer) (index int64, proof Proof) {
	sig := signer.SignData(seed[:])

	proof, _ = ProofFromRawBytes(sig.RawBytes())
	index = vrf.getIndex(proof)

	return index, proof
}

// Verify ensures the proof is valid
func (vrf *VRF) Verify(seed Seed, public crypto.PublicKey, proof Proof) (index int64, result bool) {
	proofSig, err := crypto.SignatureFromRawBytes(proof[:])
	if err != nil {
		return 0, false
	}

	// Verify signature (proof)
	if !public.Verify(seed[:], proofSig) {
		return 0, false
	}

	index = vrf.getIndex(proof)

	return index, true
}

func (vrf *VRF) getIndex(proof Proof) int64 {
	h := crypto.HashH(proof[:])

	rnd64 := util.SliceToInt64(h.RawBytes())
	rnd64 = rnd64 & 0x7fffffffffffffff

	// construct the numerator and denominator for normalizing the proof uint between [0, 1]
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
