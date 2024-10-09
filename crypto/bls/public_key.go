package bls

import (
	"bytes"
	"crypto/subtle"
	"io"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	cbor "github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/encoding"
)

var _ crypto.PublicKey = &PublicKey{}

const PublicKeySize = 96

type PublicKey struct {
	pointG2 *bls12381.G2Affine // Lazily initialized point on G2.
	data    []byte             // Raw public key data.
}

// PublicKeyFromString decodes the input string and returns the PublicKey
// if the string is a valid bech32m encoding of a BLS public key.
func PublicKeyFromString(text string) (*PublicKey, error) {
	// Decode the bech32m encoded public key.
	hrp, typ, data, err := bech32m.DecodeToBase256WithTypeNoLimit(text)
	if err != nil {
		return nil, err
	}

	// Check if hrp is valid
	if hrp != crypto.PublicKeyHRP {
		return nil, crypto.InvalidHRPError(hrp)
	}

	if typ != crypto.SignatureTypeBLS {
		return nil, crypto.InvalidSignatureTypeError(typ)
	}

	return PublicKeyFromBytes(data)
}

// PublicKeyFromBytes constructs a BLS public key from the raw bytes.
func PublicKeyFromBytes(data []byte) (*PublicKey, error) {
	if len(data) != PublicKeySize {
		return nil, crypto.InvalidLengthError(len(data))
	}

	return &PublicKey{data: data}, nil
}

// Bytes returns the raw byte representation of the public key.
func (pub *PublicKey) Bytes() []byte {
	return pub.data
}

// String returns a human-readable string for the BLS public key.
func (pub *PublicKey) String() string {
	str, _ := bech32m.EncodeFromBase256WithType(
		crypto.PublicKeyHRP,
		crypto.SignatureTypeBLS,
		pub.Bytes())

	return str
}

// MarshalCBOR encodes the public key into CBOR format.
func (pub *PublicKey) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(pub.Bytes())
}

// UnmarshalCBOR decodes the public key from CBOR format.
func (pub *PublicKey) UnmarshalCBOR(bs []byte) error {
	var data []byte
	if err := cbor.Unmarshal(bs, &data); err != nil {
		return err
	}

	return pub.Decode(bytes.NewReader(data))
}

// Encode writes the raw bytes of the public key to the provided writer.
func (pub *PublicKey) Encode(w io.Writer) error {
	return encoding.WriteElements(w, pub.Bytes())
}

// Decode reads the raw bytes of the public key from the provided reader and initializes the public key.
func (pub *PublicKey) Decode(r io.Reader) error {
	data := make([]byte, PublicKeySize)
	err := encoding.ReadElements(r, data)
	if err != nil {
		return err
	}

	p, _ := PublicKeyFromBytes(data)
	*pub = *p

	return nil
}

// Verify checks that a signature is valid for the given message and public key.
// It's defined in section 2.6 of the spec: CoreVerify.
func (pub *PublicKey) Verify(msg []byte, sig crypto.Signature) error {
	if sig == nil {
		return crypto.ErrInvalidPublicKey
	}

	r := sig.(*Signature)
	pointG1, err := r.PointG1()
	if err != nil {
		return err
	}
	pointG2, err := pub.PointG2()
	if err != nil {
		return err
	}
	qAffine, err := bls12381.HashToG1(msg, dst)
	if err != nil {
		panic(err)
	}

	var negP bls12381.G2Affine
	negP.Neg(&gen2Aff)

	check, _ := bls12381.PairingCheck(
		[]bls12381.G1Affine{qAffine, *pointG1},
		[]bls12381.G2Affine{*pointG2, negP})

	if !check {
		return crypto.ErrInvalidSignature
	}

	return nil
}

func (*PublicKey) SerializeSize() int {
	return PublicKeySize
}

// EqualsTo checks if the current public key is equal to another public key.
func (pub *PublicKey) EqualsTo(x crypto.PublicKey) bool {
	xBLS, ok := x.(*PublicKey)
	if !ok {
		return false
	}

	return subtle.ConstantTimeCompare(pub.data, xBLS.data) == 1
}

// AccountAddress returns the account address derived from the public key.
func (pub *PublicKey) AccountAddress() crypto.Address {
	data := hash.Hash160(hash.Hash256(pub.Bytes()))
	addr := crypto.NewAddress(crypto.AddressTypeBLSAccount, data)

	return addr
}

// ValidatorAddress returns the validator address derived from the public key.
func (pub *PublicKey) ValidatorAddress() crypto.Address {
	data := hash.Hash160(hash.Hash256(pub.Bytes()))
	addr := crypto.NewAddress(crypto.AddressTypeValidator, data)

	return addr
}

// VerifyAddress checks if the provided address matches the derived address from the public key.
func (pub *PublicKey) VerifyAddress(addr crypto.Address) error {
	if addr.IsValidatorAddress() {
		if addr != pub.ValidatorAddress() {
			return crypto.AddressMismatchError{
				Expected: pub.ValidatorAddress(),
				Got:      addr,
			}
		}
	} else {
		if addr != pub.AccountAddress() {
			return crypto.AddressMismatchError{
				Expected: pub.AccountAddress(),
				Got:      addr,
			}
		}
	}

	return nil
}

// PointG2 returns the point on G2 for the public key.
func (pub *PublicKey) PointG2() (*bls12381.G2Affine, error) {
	if pub.pointG2 != nil {
		return pub.pointG2, nil
	}

	g2Aff := new(bls12381.G2Affine)
	err := g2Aff.Unmarshal(pub.data)
	if err != nil {
		return nil, err
	}
	if g2Aff.IsInfinity() || !g2Aff.IsInSubGroup() {
		return nil, crypto.ErrInvalidPublicKey
	}

	pub.pointG2 = g2Aff

	return g2Aff, nil
}
