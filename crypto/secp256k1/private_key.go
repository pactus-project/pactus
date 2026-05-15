package secp256k1

import (
	"strings"

	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/bech32m"
)

var _ crypto.PrivateKey = &PrivateKey{}

const PrivateKeySize = 32

type PrivateKey struct {
	inner *secp.PrivateKey
}

// PrivateKeyFromString decodes the input string and returns the PrivateKey
// if the string is a valid bech32m encoding of a secp256k1 private key.
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

	if typ != crypto.SignatureTypeSecp256k1 {
		return nil, crypto.InvalidSignatureTypeError(typ)
	}

	return PrivateKeyFromBytes(data)
}

// PrivateKeyFromBytes constructs a secp256k1 private key from the raw bytes.
func PrivateKeyFromBytes(data []byte) (*PrivateKey, error) {
	if len(data) != PrivateKeySize {
		return nil, crypto.InvalidLengthError(len(data))
	}

	var scalar secp.ModNScalar
	overflow := scalar.SetByteSlice(data)
	if overflow || scalar.IsZero() {
		return nil, crypto.ErrInvalidPrivateKey
	}

	inner := secp.NewPrivateKey(&scalar)

	return &PrivateKey{inner}, nil
}

// String returns a human-readable string for the secp256k1 private key.
func (prv *PrivateKey) String() string {
	str, _ := bech32m.EncodeFromBase256WithType(
		crypto.PrivateKeyHRP,
		crypto.SignatureTypeSecp256k1,
		prv.Bytes())

	return strings.ToUpper(str)
}

// Bytes return the raw bytes of the private key.
func (prv *PrivateKey) Bytes() []byte {
	return prv.inner.Serialize()
}

// Sign calculates the signature from the private key and given message.
// It's defined in section 2.6 of the spec: CoreSign.
func (prv *PrivateKey) Sign(msg []byte) crypto.Signature {
	return prv.SignNative(msg)
}

func (prv *PrivateKey) SignNative(msg []byte) *Signature {
	// RFC6979 deterministic ECDSA with low-s canonicalization.
	sig := ecdsa.Sign(prv.inner, hash.Hash256(msg))

	var bs [SignatureSize]byte
	rScalar := sig.R()
	sScalar := sig.S()
	r := rScalar.Bytes()
	s := sScalar.Bytes()
	copy(bs[:32], r[:])
	copy(bs[32:], s[:])

	return &Signature{
		data: bs[:],
	}
}

func (prv *PrivateKey) PublicKeyNative() *PublicKey {
	return &PublicKey{inner: prv.inner.PubKey()}
}

func (prv *PrivateKey) PublicKey() crypto.PublicKey {
	return prv.PublicKeyNative()
}

func (prv *PrivateKey) EqualsTo(x crypto.PrivateKey) bool {
	xSecp, ok := x.(*PrivateKey)
	if !ok {
		return false
	}

	return prv.inner.Key.Equals(&xSecp.inner.Key)
}
