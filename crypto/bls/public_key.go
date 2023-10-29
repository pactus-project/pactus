package bls

import (
	"bytes"
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
	pointG2 bls12381.PointG2
}

// PublicKeyFromString decodes the string encoding of a BLS public key
// and returns the public key if text is a valid encoding for BLS public key.
func PublicKeyFromString(text string) (*PublicKey, error) {
	// Decode the bech32m encoded public key.
	hrp, typ, data, err := bech32m.DecodeToBase256WithTypeNoLimit(text)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey, err.Error())
	}

	// Check if hrp is valid
	if hrp != crypto.PublicKeyHRP {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey, "invalid hrp: %v", hrp)
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

	g2 := bls12381.NewG2()

	pointG2, err := g2.FromCompressed(data)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey, err.Error())
	}

	if g2.IsZero(pointG2) {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey,
			"public key is zero")
	}

	return &PublicKey{pointG2: *pointG2}, nil
}

func (pub *PublicKey) Bytes() []byte {
	g2 := bls12381.NewG2()

	return g2.ToCompressed(pub.point())
}

// String returns a human-readable string for the BLS public key.
func (pub *PublicKey) String() string {
	str, _ := bech32m.EncodeFromBase256WithType(
		crypto.PublicKeyHRP,
		crypto.SignatureTypeBLS,
		pub.Bytes())

	return str
}

func (pub *PublicKey) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(pub.Bytes())
}

func (pub *PublicKey) UnmarshalCBOR(bs []byte) error {
	var data []byte
	if err := cbor.Unmarshal(bs, &data); err != nil {
		return err
	}

	return pub.Decode(bytes.NewReader(data))
}

func (pub *PublicKey) Encode(w io.Writer) error {
	return encoding.WriteElements(w, pub.Bytes())
}

func (pub *PublicKey) Decode(r io.Reader) error {
	data := make([]byte, PublicKeySize)

	err := encoding.ReadElements(r, data)
	if err != nil {
		return err
	}

	p, err := PublicKeyFromBytes(data)
	if err != nil {
		return err
	}

	*pub = *p

	return nil
}

// The Verify checks that a signature is valid for the given message and public key.
// It's defined in section 2.6 of the spec: CoreVerify.
func (pub *PublicKey) Verify(msg []byte, sig crypto.Signature) error {
	if sig == nil {
		return errors.Error(errors.ErrInvalidSignature)
	}

	g1 := bls12381.NewG1()

	r := sig.(*Signature)
	if g1.IsZero(&r.pointG1) {
		return errors.Errorf(errors.ErrInvalidSignature,
			"signature is zero")
	}

	q, err := g1.HashToCurve(msg, dst)
	if err != nil {
		panic(err)
	}

	g2one := bls12381.NewG2().New().Set(&bls12381.G2One)

	eng := bls12381.NewEngine()
	eng.AddPair(q, pub.point())
	eng.AddPairInv(&r.pointG1, g2one)

	if !eng.Check() {
		return crypto.ErrInvalidSignature
	}

	return nil
}

func (pub *PublicKey) EqualsTo(right crypto.PublicKey) bool {
	g2 := bls12381.NewG2()

	return g2.Equal(pub.point(), right.(*PublicKey).point())
}

func (pub *PublicKey) AccountAddress() crypto.Address {
	data := hash.Hash160(hash.Hash256(pub.Bytes()))
	addr := crypto.NewAddress(crypto.AddressTypeBLSAccount, data)

	return addr
}

func (pub *PublicKey) ValidatorAddress() crypto.Address {
	data := hash.Hash160(hash.Hash256(pub.Bytes()))
	addr := crypto.NewAddress(crypto.AddressTypeValidator, data)

	return addr
}

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

// clonePoint clones the pointG2 to make sure it remains intact.
func (pub *PublicKey) point() *bls12381.PointG2 {
	return bls12381.NewG2().New().Set(&pub.pointG2)
}
