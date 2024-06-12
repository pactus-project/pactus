package bls

import (
	"bytes"
	"fmt"
	"io"

	cbor "github.com/fxamacker/cbor/v2"
	bls12381 "github.com/kilic/bls12-381"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/errors"
)

var _ crypto.PublicKey = &PublicKey{}

const PublicKeySize = 96

type PublicKey struct {
	pointG2 *bls12381.PointG2 // Lazily initialized point on G2.
	data    []byte            // Raw public key data.
}

// PublicKeyFromString decodes the string encoding of a BLS public key
// and returns the public key if text is a valid encoding for BLS public key.
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
		return nil, errors.Errorf(errors.ErrInvalidPublicKey, "invalid public key type: %v", typ)
	}

	return PublicKeyFromBytes(data)
}

// PublicKeyFromBytes constructs a BLS public key from the raw bytes.
func PublicKeyFromBytes(data []byte) (*PublicKey, error) {
	if len(data) != PublicKeySize {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey,
			"public key should be %d bytes, but it is %v bytes", PublicKeySize, len(data))
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
		return errors.Error(errors.ErrInvalidSignature)
	}
	g1 := bls12381.NewG1()
	g2 := bls12381.NewG2()

	r := sig.(*Signature)
	pointG1, err := r.PointG1()
	if err != nil {
		return err
	}
	pointG2, err := pub.PointG2()
	if err != nil {
		return err
	}
	q, err := g1.HashToCurve(msg, dst)
	if err != nil {
		return err
	}
	g2one := g2.New().Set(&bls12381.G2One)

	eng := bls12381.NewEngine()
	eng.AddPair(q, &pointG2)
	eng.AddPairInv(&pointG1, g2one)

	if !eng.Check() {
		return crypto.ErrInvalidSignature
	}

	return nil
}

// EqualsTo checks if the current public key is equal to another public key.
func (pub *PublicKey) EqualsTo(right crypto.PublicKey) bool {
	return bytes.Equal(pub.data, right.(*PublicKey).data)
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
func (pub *PublicKey) PointG2() (bls12381.PointG2, error) {
	if pub.pointG2 != nil {
		return *pub.pointG2, nil
	}

	g2 := bls12381.NewG2()
	pointG2, err := g2.FromCompressed(pub.data)
	if err != nil {
		return bls12381.PointG2{}, err
	}
	if g2.IsZero(pointG2) {
		return bls12381.PointG2{}, fmt.Errorf("public key is zero")
	}

	pub.pointG2 = pointG2

	return *pointG2, nil
}
