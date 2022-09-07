package bls

import (
	"github.com/herumi/bls-go-binary/bls"
	bls12381 "github.com/kilic/bls12-381"
	"github.com/zarbchain/zarb-go/crypto"
)

var g1 = bls12381.NewG1()
var g2 = bls12381.NewG2()

// set Ciphersuite for Basic mode
// https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-bls-signature-04#section-4.2.1
var dst = []byte("BLS_SIG_BLS12381G1_XMD:SHA-256_SSWU_RO_NUL_")

func Aggregate(sigs []*Signature) *Signature {
	if len(sigs) == 0 {
		return nil
	}
	aggPointG1 := sigs[0].pointG1
	for i := 1; i < len(sigs); i++ {
		g1.Add(
			&aggPointG1,
			&aggPointG1,
			&sigs[i].pointG1)
	}

	err = bls.SetMapToMode(bls.IRTF)
	if err != nil {
		panic(err)
	}

	// set G2 generator
	// https://docs.rs/bls12_381_plus/0.6.0/bls12_381_plus/notes/design/index.html#fixed-generators
	var gen bls.PublicKey
	err = gen.SetHexString(`
		1
		24aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8
		13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e
		ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801
		606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be`)
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
	// aggregated := new(bls.Sign)
	// signatures := make([]bls.Sign, len(sigs))

	// for i, s := range sigs {
	// 	signatures[i] = s.signature
	// }

	// aggregated.Aggregate(signatures)

	// return &Signature{
	// 	signature: *aggregated,
	// }

	return nil
}

	aggPub := PublicKey{pointG2: aggPointG2}
	return aggPub.Verify(msg, sig) == nil
}

// GenerateTestSigner generates a signer for testing.
func GenerateTestSigner() crypto.Signer {
	_, prv := GenerateTestKeyPair()
	return crypto.NewSigner(prv)
}

// GenerateTestKeyPair generates a key pair for testing.
func GenerateTestKeyPair() (*PublicKey, *PrivateKey) {
	buf := make([]byte, PrivateKeySize)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	prv, _ := PrivateKeyFromBytes(buf)
	pub := prv.PublicKey().(*PublicKey)

	return pub, prv
}
