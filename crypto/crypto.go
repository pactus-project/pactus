package crypto

import (
	"github.com/herumi/bls-go-binary/bls"
	"github.com/btcsuite/btcutil/bech32"
)

var TreasuryAddress = Address{data: addressData{Address: [AddressSize]byte{0}}}

func init() {
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

// EncodeFromBase256 converts a base256-encoded byte slice into a base32-encoded
// byte slice and then encodes it into a bech32 string with the given
// human-readable part (HRP).  The HRP will be converted to lowercase if needed
// since mixed cased encodings are not permitted and lowercase is used for
// checksum purposes.
func EncodeFromBase256(hrp string, data []byte) (string, error) {
	converted, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", err
	}
	return bech32.Encode(hrp, converted)
}

// DecodeToBase256 decodes a bech32-encoded string into its associated
// human-readable part (HRP) and base32-encoded data, converts that data to a
// base256-encoded byte slice and returns it along with the lowercase HRP.
func DecodeToBase256(bech string) (string, []byte, error) {
	hrp, data, err := bech32.Decode(bech)
	if err != nil {
		return "", nil, err
	}
	converted, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return "", nil, err
	}
	return hrp, converted, nil
}
