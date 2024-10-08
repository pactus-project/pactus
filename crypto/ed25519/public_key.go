package ed25519

import (
	"bytes"
	"crypto/ed25519"
	"io"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/encoding"
)

var _ crypto.PublicKey = &PublicKey{}

const PublicKeySize = 32

type PublicKey struct {
	inner ed25519.PublicKey
}

// PublicKeyFromString decodes the input string and returns the PublicKey
// if the string is a valid bech32m encoding of a Ed25519 public key.
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

	if typ != crypto.SignatureTypeEd25519 {
		return nil, crypto.InvalidSignatureTypeError(typ)
	}

	return PublicKeyFromBytes(data)
}

// PublicKeyFromBytes constructs a Ed25519 public key from the raw bytes.
func PublicKeyFromBytes(data []byte) (*PublicKey, error) {
	if len(data) != PublicKeySize {
		return nil, crypto.InvalidLengthError(len(data))
	}

	return &PublicKey{data}, nil
}

// Bytes returns the raw byte representation of the public key.
func (pub *PublicKey) Bytes() []byte {
	return pub.inner
}

// String returns a human-readable string for the Ed25519 public key.
func (pub *PublicKey) String() string {
	str, _ := bech32m.EncodeFromBase256WithType(
		crypto.PublicKeyHRP,
		crypto.SignatureTypeEd25519,
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

func (*PublicKey) SerializeSize() int {
	return PublicKeySize
}

// Verify checks that a signature is valid for the given message and public key.
// It's defined in section 2.6 of the spec: CoreVerify.
func (pub *PublicKey) Verify(msg []byte, sig crypto.Signature) error {
	if sig == nil {
		return crypto.ErrInvalidPublicKey
	}

	if !ed25519.Verify(pub.inner, msg, sig.Bytes()) {
		return crypto.ErrInvalidSignature
	}

	return nil
}

// EqualsTo checks if the current public key is equal to another public key.
func (pub *PublicKey) EqualsTo(x crypto.PublicKey) bool {
	xEd25519, ok := x.(*PublicKey)
	if !ok {
		return false
	}

	return pub.inner.Equal(xEd25519.inner)
}

// AccountAddress returns the account address derived from the public key.
func (pub *PublicKey) AccountAddress() crypto.Address {
	data := hash.Hash160(hash.Hash256(pub.Bytes()))
	addr := crypto.NewAddress(crypto.AddressTypeEd25519Account, data)

	return addr
}

// VerifyAddress checks if the provided address matches the derived address from the public key.
func (pub *PublicKey) VerifyAddress(addr crypto.Address) error {
	if addr != pub.AccountAddress() {
		return crypto.AddressMismatchError{
			Expected: pub.AccountAddress(),
			Got:      addr,
		}
	}

	return nil
}
