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

type BLSPublicKey struct {
	data publicKeyData
}

type publicKeyData struct {
	PublicKey *bls.PublicKey
}

func PublicKeyFromString(text string) (*BLSPublicKey, error) {
	data, err := hex.DecodeString(text) // from bech32 string
	if err != nil {
		return nil, err
	}

	return PublicKeyFromRawBytes(data)
}

func PublicKeyFromRawBytes(data []byte) (*BLSPublicKey, error) {
	if len(data) != PublicKeySize {
		return nil, fmt.Errorf("invalid public key")
	}
	pk := new(bls.PublicKey)
	if err := pk.Deserialize(data); err != nil {
		return nil, err
	}

	var pub BLSPublicKey
	pub.data.PublicKey = pk

	if err := pub.SanityCheck(); err != nil {
		return nil, err
	}

	return &pub, nil
}

func (pub BLSPublicKey) RawBytes() []byte {
	if pub.data.PublicKey == nil {
		return nil
	}
	return pub.data.PublicKey.Serialize()
}

func (pub BLSPublicKey) String() string {
	if pub.data.PublicKey == nil {
		return ""
	}
	return pub.data.PublicKey.SerializeToHexStr()
}

func (pub BLSPublicKey) MarshalText() ([]byte, error) {
	return []byte(pub.String()), nil
}

func (pub *BLSPublicKey) UnmarshalText(text []byte) error {
	p, err := PublicKeyFromString(string(text))
	if err != nil {
		return err
	}

	*pub = *p
	return nil
}

func (pub *BLSPublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pub.String())
}

func (pub *BLSPublicKey) UnmarshalJSON(bz []byte) error {
	var text string
	if err := json.Unmarshal(bz, &text); err != nil {
		return err
	}
	return pub.UnmarshalText([]byte(text))
}

func (pub *BLSPublicKey) MarshalCBOR() ([]byte, error) {
	if pub.data.PublicKey == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return cbor.Marshal(pub.RawBytes())
}

func (pub *BLSPublicKey) UnmarshalCBOR(bs []byte) error {
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

func (pub *BLSPublicKey) SanityCheck() error {
	bs := pub.RawBytes()
	if len(bs) != PublicKeySize {
		return fmt.Errorf("public key should be %v bytes but it is %v bytes", PublicKeySize, len(bs))
	}
	return nil
}

func (pub *BLSPublicKey) Verify(msg []byte, sig crypto.Signature) bool {
	return sig.(*BLSSignature).data.Signature.VerifyByte(pub.data.PublicKey, hash.Hash256(msg))
}

func (pub *BLSPublicKey) EqualsTo(right crypto.PublicKey) bool {
	return pub.data.PublicKey.IsEqual(right.(*BLSPublicKey).data.PublicKey)
}

func (pub *BLSPublicKey) Address() crypto.Address {
	addr, _ := crypto.AddressFromRawBytes(hash.Hash160(hash.Hash256(pub.RawBytes())))
	return addr
}
