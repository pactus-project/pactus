package crypto

import (
	"github.com/herumi/bls-go-binary/bls"
)

var MintbaseAddress = Address{data: addressData{Address: [AddressSize]byte{0}}}

func init() {
	bls.Init(bls.BLS12_381)

	if err := bls.Init(bls.BLS12_381); err != nil {
		panic(err)
	}

	// Check subgroup order for pubkeys and signatures.
	bls.VerifyPublicKeyOrder(true)
	bls.VerifySignatureOrder(true)
}

func Aggregate(sigs []Signature) Signature {
	aggregated := new(bls.Sign)
	signatures := make([]bls.Sign, len(sigs))

	for i, s := range sigs {
		signatures[i] = *s.data.Signature
	}

	aggregated.Aggregate(signatures)

	return Signature{
		data: signatureData{
			Signature: aggregated,
		},
	}
}

func VerifyAggregated(aggregated Signature, pubs []PublicKey, msg []byte) bool {
	pubVec := make([]bls.PublicKey, len(pubs))
	for i, p := range pubs {
		pubVec[i] = *p.data.PublicKey
	}
	return aggregated.data.Signature.FastAggregateVerify(pubVec, Hash256(msg))
}
