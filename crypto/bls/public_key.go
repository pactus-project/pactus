package bls

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/herumi/bls-go-binary/bls"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/encoding"
)

const PublicKeySize = 96

type PublicKey struct {
	publicKey *bls.PublicKey
}

func PublicKeyFromString(text string) (*PublicKey, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}

	return PublicKeyFromBytes(data)
}

func PublicKeyFromBytes(data []byte) (*PublicKey, error) {
	if len(data) != PublicKeySize {
		return nil, fmt.Errorf("invalid public key")
	}
	pk := new(bls.PublicKey)
	if err := pk.Deserialize(data); err != nil {
		return nil, err
	}

	var pub PublicKey
	pub.publicKey = pk

	return &pub, nil
}

func (pub PublicKey) Bytes() []byte {
	if pub.publicKey == nil {
		return nil
	}
	return pub.publicKey.Serialize()
}

func (pub PublicKey) String() string {
	if pub.publicKey == nil {
		return ""
	}
	return pub.publicKey.SerializeToHexStr()
}

func (pub *PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pub.String())
}

func (pub *PublicKey) MarshalCBOR() ([]byte, error) {
	if pub.publicKey == nil {
		return nil, fmt.Errorf("invalid public key")
	}
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
		return fmt.Errorf("public key is zero")
	}

	return nil
}

func (pub *PublicKey) Verify(msg []byte, sig crypto.Signature) bool {
	return sig.(*Signature).signature.VerifyByte(pub.publicKey, msg)
}

func (pub *PublicKey) EqualsTo(right crypto.PublicKey) bool {
	return pub.publicKey.IsEqual(right.(*PublicKey).publicKey)
}

func (pub *PublicKey) Address() crypto.Address {
	data := hash.Hash160(hash.Hash256(pub.Bytes()))
	data = append([]byte{crypto.AddressTypeBLS}, data...)
	addr, _ := crypto.AddressFromBytes(data)
	return addr
}

func (pub *PublicKey) VerifyAddress(addr crypto.Address) bool {
	return addr.EqualsTo(pub.Address())
}
