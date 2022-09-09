package bls

import (
	"bytes"
	"io"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/herumi/bls-go-binary/bls"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/errors"
)

const PublicKeySize = 96

type PublicKey struct {
	publicKey bls.PublicKey
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

func PublicKeyFromBytes(data []byte) (*PublicKey, error) {
	if len(data) != PublicKeySize {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey,
			"public key should be %d bytes, but it is %v bytes", PublicKeySize, len(data))
	}
	pk := new(bls.PublicKey)
	if err := pk.Deserialize(data); err != nil {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey, err.Error())
	}

	var pub PublicKey
	pub.publicKey = *pk

	return &pub, nil
}

func (pub PublicKey) Bytes() []byte {
	return pub.publicKey.Serialize()
}

// String returns a human-readable string for the BLS public key.
func (pub *PublicKey) String() string {
	data := pub.publicKey.Serialize()

	str, err := bech32m.EncodeFromBase256WithType(
		crypto.PublicKeyHRP,
		crypto.SignatureTypeBLS,
		data)
	if err != nil {
		panic(err.Error())
	}

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

func (pub *PublicKey) SanityCheck() error {
	if pub.publicKey.IsZero() {
		return errors.Errorf(errors.ErrInvalidPublicKey, "public key is zero")
	}

	return nil
}

func (pub *PublicKey) Verify(msg []byte, sig crypto.Signature) error {
	if sig == nil {
		return errors.Error(errors.ErrInvalidSignature)
	}
	if !sig.(*Signature).signature.VerifyByte(&pub.publicKey, msg) {
		return errors.Error(errors.ErrInvalidSignature)
	}
	return nil
}

func (pub *PublicKey) EqualsTo(right crypto.PublicKey) bool {
	return pub.publicKey.IsEqual(&right.(*PublicKey).publicKey)
}

func (pub *PublicKey) Address() crypto.Address {
	data := hash.Hash160(hash.Hash256(pub.Bytes()))
	data = append([]byte{crypto.SignatureTypeBLS}, data...)
	var addr crypto.Address
	copy(addr[:], data[:])
	return addr
}

func (pub *PublicKey) VerifyAddress(addr crypto.Address) error {
	if !addr.EqualsTo(pub.Address()) {
		return errors.Error(errors.ErrInvalidAddress)
	}
	return nil
}
