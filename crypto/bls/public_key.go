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

func PublicKeyFromString(text string) (crypto.PublicKey, error) {
	data, err := hex.DecodeString(text) // from bech32 string
	if err != nil {
		return nil, err
	}

	return PublicKeyFromRawBytes(data)
}

func PublicKeyFromRawBytes(data []byte) (crypto.PublicKey, error) {
	if len(data) != PublicKeySize {
		return nil, fmt.Errorf("invalid public key")
	}
	pk := new(bls.PublicKey)
	if err := pk.Deserialize(data); err != nil {
		return nil, err
	}

	var pb BLSPublicKey
	pb.data.PublicKey = pk

	if err := pb.SanityCheck(); err != nil {
		return nil, err
	}

	return &pb, nil
}

func (pb BLSPublicKey) RawBytes() []byte {
	if pb.data.PublicKey == nil {
		return nil
	}
	return pb.data.PublicKey.Serialize()
}

func (pb BLSPublicKey) String() string {
	if pb.data.PublicKey == nil {
		return ""
	}
	return pb.data.PublicKey.SerializeToHexStr()
}

func (pb BLSPublicKey) MarshalText() ([]byte, error) {
	return []byte(pb.String()), nil
}

func (pb *BLSPublicKey) UnmarshalText(text []byte) error {
	p, err := PublicKeyFromString(string(text))
	if err != nil {
		return err
	}

	*pb = *p.(*BLSPublicKey)
	return nil
}

func (pb *BLSPublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pb.String())
}

func (pb *BLSPublicKey) UnmarshalJSON(bz []byte) error {
	var text string
	if err := json.Unmarshal(bz, &text); err != nil {
		return err
	}
	return pb.UnmarshalText([]byte(text))
}

func (pb *BLSPublicKey) MarshalCBOR() ([]byte, error) {
	if pb.data.PublicKey == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return cbor.Marshal(pb.RawBytes())
}

func (pb *BLSPublicKey) UnmarshalCBOR(bs []byte) error {
	var data []byte
	if err := cbor.Unmarshal(bs, &data); err != nil {
		return err
	}

	p, err := PublicKeyFromRawBytes(data)
	if err != nil {
		return err
	}

	*pb = *p.(*BLSPublicKey)
	return nil
}

func (pb *BLSPublicKey) SanityCheck() error {
	bs := pb.RawBytes()
	if len(bs) != PublicKeySize {
		return fmt.Errorf("public key should be %v bytes but it is %v bytes", PublicKeySize, len(bs))
	}
	return nil
}

func (pb *BLSPublicKey) Verify(msg []byte, sig crypto.Signature) bool {
	return sig.(*BLSSignature).data.Signature.VerifyByte(pb.data.PublicKey, hash.Hash256(msg))
}

func (pb *BLSPublicKey) EqualsTo(right crypto.PublicKey) bool {
	return pb.data.PublicKey.IsEqual(right.(*BLSPublicKey).data.PublicKey)
}

func (pb *BLSPublicKey) Address() crypto.Address {
	addr, _ := crypto.AddressFromRawBytes(hash.Hash160(hash.Hash256(pb.RawBytes())))
	return addr
}
