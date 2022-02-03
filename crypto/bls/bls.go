package bls

import (
	"github.com/herumi/bls-go-binary/bls"
	"github.com/zarbchain/zarb-go/crypto"
)

func init() {
	err := bls.Init(bls.BLS12_381)
	if err != nil {
		panic(err)
	}

	// use serialization mode compatible with ETH
	bls.SetETHserialization(true)

	err = bls.SetMapToMode(bls.IRTF)
	if err != nil {
		panic(err)
	}

	// set G2 generator
	// https://docs.rs/bls12_381_plus/0.6.0/bls12_381_plus/notes/design/index.html#fixed-generators
	var gen bls.PublicKey
	err = gen.SetHexString("1 24aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8 13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801 606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be")
	if err != nil {
		panic(err)
	}
	err = bls.SetGeneratorOfPublicKey(&gen)
	if err != nil {
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
	return aggregated.data.Signature.FastAggregateVerify(pubVec, msg)
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
