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

func Aggregate(sigs []crypto.Signature) crypto.Signature {
	aggregated := new(bls.Sign)
	signatures := make([]bls.Sign, len(sigs))

	for i, s := range sigs {
		signatures[i] = *s.(*BLSSignature).data.Signature
	}

	aggregated.Aggregate(signatures)

	return &BLSSignature{
		data: signatureData{
			Signature: aggregated,
		},
	}
}

func VerifyAggregated(aggregated crypto.Signature, pubs []crypto.PublicKey, msg []byte) bool {
	pubVec := make([]bls.PublicKey, len(pubs))
	for i, p := range pubs {
		pubVec[i] = *p.(*BLSPublicKey).data.PublicKey
	}
	return aggregated.(*BLSSignature).data.Signature.FastAggregateVerify(pubVec, hash.Hash256(msg))
}

func RandomKeyPair() (*BLSPublicKey, *BLSPrivateKey) {
	pv := new(BLSPrivateKey)
	pv.data.SecretKey = new(bls.SecretKey)
	pv.data.SecretKey.SetByCSPRNG()

	pb := new(BLSPublicKey)
	pb.data.PublicKey = pv.data.SecretKey.GetPublicKey()

	return pb, pv
}

// ---------
// For tests
func GenerateTestSigner() crypto.Signer {
	_, priv := RandomKeyPair()
	return crypto.NewSigner(priv)
}

func GenerateTestKeyPair() (*BLSPublicKey, *BLSPrivateKey) {
	return RandomKeyPair()
}
