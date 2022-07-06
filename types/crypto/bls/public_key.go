package bls

import (
	"bytes"
	"io"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/herumi/bls-go-binary/bls"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/util/bech32m"
	"github.com/zarbchain/zarb-go/util/encoding"
	"github.com/zarbchain/zarb-go/util/errors"
)

const (
	PublicKeySize = 96
	hrpPublicKey  = "public"
)

type PublicKey struct {
	publicKey bls.PublicKey
}

// PublicKeyFromString decodes the string encoding of a BLS public key
// and returns the public key if text is a valid encoding for BLS public key.
func PublicKeyFromString(text string) (*PublicKey, error) {
	// Decode the bech32m encoded public key.
	hrp, data, err := bech32m.DecodeNoLimit(text)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey, err.Error())
	}

	// Check if hrp is valid
	if hrp != hrpPublicKey {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey, "invalid hrp: %v", hrp)
	}

	// The first byte of the decoded public key is the signature type, it must
	// exist.
	if len(data) < 1 {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey, "no public key type")
	}

	// ...and should be 1 for BLS signature.
	sigType := data[0]
	if sigType != crypto.SignatureTypeBLS {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey, "invalid public key type: %v", sigType)
	}

	// The remaining characters of the public key returned are grouped into
	// words of 5 bits. In order to restore the original program
	// bytes, we'll need to regroup into 8 bit words.
	regrouped, err := bech32m.ConvertBits(data[1:], 5, 8, false)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey, err.Error())
	}

	return PublicKeyFromBytes(regrouped)
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

	// Group the public key bytes into 5 bit groups, as this is what is used to
	// encode each character in the public key string.
	converted, err := bech32m.ConvertBits(data, 8, 5, true)
	if err != nil {
		panic(err.Error())
	}

	// Concatenate the type of the public key which is 1 for BLS and program,
	// and encode the resulting bytes using bech32m encoding.
	combined := make([]byte, len(converted)+1)
	combined[0] = crypto.SignatureTypeBLS
	copy(combined[1:], converted)
	str, err := bech32m.Encode(hrpPublicKey, combined)
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
