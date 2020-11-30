package sortition

import (
	"math/big"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

type VRF struct {
	signer crypto.PrivateKey
	max    int64
}

func NewVRF(signer crypto.PrivateKey) VRF {
	return VRF{
		signer: signer,
	}
}

func (vrf *VRF) SetMax(max int64) {
	vrf.max = max
}

// Evaluate returns a random number between 0 and max with the proof
func (vrf *VRF) Evaluate(m []byte) (index int64, proof []byte) {
	// sign the hashed block height
	sig := vrf.signer.Sign(m)

	proof = make([]byte, 0)
	addrBytes := vrf.signer.PublicKey().Address().RawBytes()
	sigBytes := sig.RawBytes()
	proof = append(proof, addrBytes...)
	proof = append(proof, sigBytes...)

	index = vrf.getIndex(sigBytes)

	return index, proof
}

// Verify ensure the proof is valid
func (vrf *VRF) Verify(msg []byte, publicKey crypto.PublicKey, proof []byte) (index int64, result bool) {
	address, err := crypto.AddressFromRawBytes(proof[0:crypto.AddressSize])
	if err != nil {
		return 0, false
	}

	sig, err := crypto.SignatureFromRawBytes(proof[crypto.AddressSize:])
	if err != nil {
		return 0, false
	}

	// Verify address
	if !address.Verify(publicKey) {
		return 0, false
	}

	// Verify signature (proof)
	if !publicKey.Verify(msg, &sig) {
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
