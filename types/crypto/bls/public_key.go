package bls

import (
	"bytes"
	"encoding/hex"
	"io"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/herumi/bls-go-binary/bls"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/util/encoding"
	"github.com/zarbchain/zarb-go/util/errors"
)

const PublicKeySize = 96

type PublicKey struct {
	publicKey bls.PublicKey
}

func PublicKeyFromString(text string) (*PublicKey, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey, err.Error())
	}

	return PublicKeyFromBytes(data)
}

func PublicKeyFromBytes(data []byte) (*PublicKey, error) {
	if len(data) != PublicKeySize {
		return nil, errors.Errorf(errors.ErrInvalidPublicKey, "public key should be %d bytes, but it is %v bytes", PublicKeySize, len(data))
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

func (pub PublicKey) String() string {
	return pub.publicKey.SerializeToHexStr()
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
