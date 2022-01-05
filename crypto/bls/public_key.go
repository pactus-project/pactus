package bls

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/herumi/bls-go-binary/bls"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
)

const PublicKeySize = 96

type PublicKey struct {
	data publicKeyData
}

type publicKeyData struct {
	PublicKey *bls.PublicKey
}

func PublicKeyFromString(text string) (*PublicKey, error) {
	data, err := hex.DecodeString(text) // from bech32 string
	if err != nil {
		return nil, err
	}

	return PublicKeyFromRawBytes(data)
}

func PublicKeyFromRawBytes(data []byte) (*PublicKey, error) {
	if len(data) != PublicKeySize {
		return nil, fmt.Errorf("invalid public key")
	}
	pk := new(bls.PublicKey)
	if err := pk.Deserialize(data); err != nil {
		return nil, err
	}

	var pub PublicKey
	pub.data.PublicKey = pk

	if err := pub.SanityCheck(); err != nil {
		return nil, err
	}

	return &pub, nil
}

func (pub PublicKey) RawBytes() []byte {
	if pub.data.PublicKey == nil {
		return nil
	}
	return pub.data.PublicKey.Serialize()
}

func (pub PublicKey) String() string {
	if pub.data.PublicKey == nil {
		return ""
	}
	return pub.data.PublicKey.SerializeToHexStr()
}

func (pub PublicKey) MarshalText() ([]byte, error) {
	return []byte(pub.String()), nil
}

func (pub *PublicKey) UnmarshalText(text []byte) error {
	p, err := PublicKeyFromString(string(text))
	if err != nil {
		return err
	}

	*pub = *p
	return nil
}

func (pub *PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pub.String())
}

func (pub *PublicKey) UnmarshalJSON(bz []byte) error {
	var text string
	if err := json.Unmarshal(bz, &text); err != nil {
		return err
	}
	return pub.UnmarshalText([]byte(text))
}

func (pub *PublicKey) MarshalCBOR() ([]byte, error) {
	if pub.data.PublicKey == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return cbor.Marshal(pub.RawBytes())
}

func (pub *PublicKey) UnmarshalCBOR(bs []byte) error {
	var data []byte
	if err := cbor.Unmarshal(bs, &data); err != nil {
		return err
	}

	p, err := PublicKeyFromRawBytes(data)
	if err != nil {
		return err
	}

	*pub = *p
	return nil
}

func (pub *PublicKey) SanityCheck() error {
	bs := pub.RawBytes()
	if len(bs) != PublicKeySize {
		return fmt.Errorf("public key should be %v bytes but it is %v bytes", PublicKeySize, len(bs))
	}
	return nil
}

func (pub *PublicKey) Verify(msg []byte, sig crypto.Signature) bool {
	return sig.(*Signature).data.Signature.VerifyByte(pub.data.PublicKey, msg)
}

func (pub *PublicKey) EqualsTo(right crypto.PublicKey) bool {
	return pub.data.PublicKey.IsEqual(right.(*PublicKey).data.PublicKey)
}

func (pub *PublicKey) Address() crypto.Address {
	data := hash.Hash160(hash.Hash256(pub.RawBytes()))
	data = append([]byte{crypto.AddressTypeBLS}, data...)
	addr, _ := crypto.AddressFromRawBytes(data)
	return addr
}
