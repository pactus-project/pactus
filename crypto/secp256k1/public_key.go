package secp256k1

import (
	"bytes"
	"crypto/subtle"
	"io"

	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	cbor "github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/encoding"
)

var _ crypto.PublicKey = &PublicKey{}

const PublicKeySize = 33

type PublicKey struct {
	inner *secp.PublicKey
}

// PublicKeyFromString decodes the input string and returns the PublicKey
// if the string is a valid bech32m encoding of a Secp256k1 public key.
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

	if typ != crypto.SignatureTypeSecp256k1 {
		return nil, crypto.InvalidSignatureTypeError(typ)
	}

	return PublicKeyFromBytes(data)
}

// PublicKeyFromBytes constructs a secp256k1 public key from compressed bytes.
func PublicKeyFromBytes(data []byte) (*PublicKey, error) {
	if len(data) != PublicKeySize {
		return nil, crypto.InvalidLengthError(len(data))
	}
	inner, err := secp.ParsePubKey(data)
	if err != nil {
		return nil, err
	}

	return &PublicKey{inner: inner}, nil
}

// Bytes returns the raw byte representation of the public key.
func (pub *PublicKey) Bytes() []byte {
	if pub == nil || pub.inner == nil {
		return nil
	}

	return pub.inner.SerializeCompressed()
}

// String returns a human-readable string for the secp256k1 public key.
func (pub *PublicKey) String() string {
	str, _ := bech32m.EncodeFromBase256WithType(
		crypto.PublicKeyHRP,
		crypto.SignatureTypeSecp256k1,
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

	p, err := PublicKeyFromBytes(data)
	if err != nil {
		return err
	}
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

	sigSecp, ok := sig.(*Signature)
	if !ok {
		return crypto.ErrInvalidSignature
	}
	rBytes := sigSecp.data[:32]
	sBytes := sigSecp.data[32:]
	var r, s secp.ModNScalar
	if r.SetByteSlice(rBytes) || s.SetByteSlice(sBytes) {
		return crypto.ErrInvalidSignature
	}
	if r.IsZero() || s.IsZero() || s.IsOverHalfOrder() {
		return crypto.ErrInvalidSignature
	}

	if !ecdsa.NewSignature(&r, &s).Verify(msg, pub.inner) {
		return crypto.ErrInvalidSignature
	}

	return nil
}

// EqualsTo checks if the current public key is equal to another public key.
func (pub *PublicKey) EqualsTo(x crypto.PublicKey) bool {
	xSecp, ok := x.(*PublicKey)
	if !ok {
		return false
	}

	return subtle.ConstantTimeCompare(pub.Bytes(), xSecp.Bytes()) == 1
}

// AccountAddress returns the account address derived from the public key.
func (pub *PublicKey) AccountAddress() crypto.Address {
	data := hash.Hash160(hash.Hash256(pub.Bytes()))
	addr := crypto.NewAddress(crypto.AddressTypeSecp256k1Account, data)

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
