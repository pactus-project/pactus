package bls

import (
	"github.com/herumi/bls-go-binary/bls"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

func init() {
	if err := bls.Init(bls.BLS12_381); err != nil {
		panic(err)
	}

	// Check subgroup order for pubkeys and signatures.
	bls.VerifyPublicKeyOrder(true)
	bls.VerifySignatureOrder(true)
}

func Aggregate(sigs []*Signature) *Signature {
	aggregated := new(bls.Sign)
	signatures := make([]bls.Sign, len(sigs))

	for i, s := range sigs {
		signatures[i] = *s.data.Signature
	}

	aggregated.Aggregate(signatures)

	return &Signature{
		data: signatureData{
			Signature: aggregated,
		},
	}
}

func VerifyAggregated(aggregated *Signature, pubs []*PublicKey, msg []byte) bool {
	pubVec := make([]bls.PublicKey, len(pubs))
	for i, p := range pubs {
		pubVec[i] = *p.data.PublicKey
	}
	return aggregated.data.Signature.FastAggregateVerify(pubVec, hash.Hash256(msg))
}

func RandomKeyPair() (*PublicKey, *PrivateKey) {
	prv := new(PrivateKey)
	prv.data.SecretKey = new(bls.SecretKey)
	prv.data.SecretKey.SetByCSPRNG()

	pub := new(PublicKey)
	pub.data.PublicKey = prv.data.SecretKey.GetPublicKey()

	return pub, prv
}

// ---------
// For tests
func GenerateTestSigner() crypto.Signer {
	_, prv := RandomKeyPair()
	return crypto.NewSigner(prv)
}

func GenerateTestKeyPair() (*PublicKey, *PrivateKey) {
	return RandomKeyPair()
}
