package ed25519

import (
	"crypto/ed25519"
	"strings"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/bech32m"
)

var _ crypto.PrivateKey = &PrivateKey{}

const PrivateKeySize = 32

type PrivateKey struct {
	inner ed25519.PrivateKey
}

// PrivateKeyFromString decodes the input string and returns the PrivateKey
// if the string is a valid bech32m encoding of a Ed25519 public key.
func PrivateKeyFromString(text string) (*PrivateKey, error) {
	// Decode the bech32m encoded private key.
	hrp, typ, data, err := bech32m.DecodeToBase256WithTypeNoLimit(text)
	if err != nil {
		return nil, err
	}

	// Check if hrp is valid
	if hrp != crypto.PrivateKeyHRP {
		return nil, crypto.InvalidHRPError(hrp)
	}

	if typ != crypto.SignatureTypeEd25519 {
		return nil, crypto.InvalidSignatureTypeError(typ)
	}

	return PrivateKeyFromBytes(data)
}

// PrivateKeyFromBytes constructs a ED25519 private key from the raw bytes.
func PrivateKeyFromBytes(data []byte) (*PrivateKey, error) {
	if len(data) != PrivateKeySize {
		return nil, crypto.InvalidLengthError(len(data))
	}
	inner := ed25519.NewKeyFromSeed(data)

	return &PrivateKey{inner}, nil
}

// String returns a human-readable string for the ED25519 private key.
func (prv *PrivateKey) String() string {
	str, _ := bech32m.EncodeFromBase256WithType(
		crypto.PrivateKeyHRP,
		crypto.SignatureTypeEd25519,
		prv.Bytes())

	return strings.ToUpper(str)
}

// Bytes return the raw bytes of the private key.
func (prv *PrivateKey) Bytes() []byte {
	return prv.inner[:PrivateKeySize]
}

// Sign calculates the signature from the private key and given message.
// It's defined in section 2.6 of the spec: CoreSign.
func (prv *PrivateKey) Sign(msg []byte) crypto.Signature {
	return prv.SignNative(msg)
}

func (prv *PrivateKey) SignNative(msg []byte) *Signature {
	sig := ed25519.Sign(prv.inner, msg)

	return &Signature{
		data: sig,
	}
}

func (prv *PrivateKey) PublicKeyNative() *PublicKey {
	pub := prv.inner.Public()

	// TODO: fix me, should get from scalar multiplication.
	return &PublicKey{
		inner: pub.(ed25519.PublicKey),
	}
}

func (prv *PrivateKey) PublicKey() crypto.PublicKey {
	return prv.PublicKeyNative()
}

func (prv *PrivateKey) EqualsTo(x crypto.PrivateKey) bool {
	xEd25519, ok := x.(*PrivateKey)
	if !ok {
		return false
	}

	return prv.inner.Equal(xEd25519.inner)
}
