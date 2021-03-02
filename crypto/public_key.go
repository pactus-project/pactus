package crypto

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/herumi/bls-go-binary/bls"
)

const PublicKeySize = 96

type PublicKey struct {
	data publicKeyData
}

type publicKeyData struct {
	PublicKey *bls.PublicKey
}

func PublicKeyFromString(text string) (PublicKey, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return PublicKey{}, err
	}

	return PublicKeyFromRawBytes(data)
}

func PublicKeyFromRawBytes(data []byte) (PublicKey, error) {
	if len(data) != PublicKeySize {
		return PublicKey{}, fmt.Errorf("Invalid public key")
	}
	pk := new(bls.PublicKey)
	if err := pk.Deserialize(data); err != nil {
		return PublicKey{}, err
	}

	var pb PublicKey
	pb.data.PublicKey = pk

	if err := pb.SanityCheck(); err != nil {
		return PublicKey{}, err
	}

	return pb, nil
}

/// -------
/// CASTING

func (pb PublicKey) RawBytes() []byte {
	if pb.data.PublicKey == nil {
		return nil
	}
	return pb.data.PublicKey.Serialize()
}

func (pb PublicKey) String() string {
	if pb.data.PublicKey == nil {
		return ""
	}
	return pb.data.PublicKey.SerializeToHexStr()
}

/// ----------
/// MARSHALING

func (pb PublicKey) MarshalText() ([]byte, error) {
	return []byte(pb.String()), nil
}

func (pb *PublicKey) UnmarshalText(text []byte) error {
	p, err := PublicKeyFromString(string(text))
	if err != nil {
		return err
	}

	*pb = p
	return nil
}

func (pb PublicKey) MarshalJSON() ([]byte, error) {
	bz, err := pb.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(bz))
}

func (pb *PublicKey) UnmarshalJSON(bz []byte) error {
	var text string
	if err := json.Unmarshal(bz, &text); err != nil {
		return err
	}
	return pb.UnmarshalText([]byte(text))
}

func (pb PublicKey) MarshalCBOR() ([]byte, error) {
	if pb.data.PublicKey == nil {
		return nil, fmt.Errorf("Invalid public key")
	}
	return cbor.Marshal(pb.RawBytes())
}

func (pb *PublicKey) UnmarshalCBOR(bs []byte) error {
	var data []byte
	if err := cbor.Unmarshal(bs, &data); err != nil {
		return err
	}

	p, err := PublicKeyFromRawBytes(data)
	if err != nil {
		return err
	}

	*pb = p
	return nil
}

/// ----------
/// ATTRIBUTES

func (pb *PublicKey) SanityCheck() error {
	bs := pb.RawBytes()
	if len(bs) != PublicKeySize {
		return fmt.Errorf("Public key should be %v bytes but it is %v bytes", PublicKeySize, len(bs))
	}
	return nil
}

func (pb *PublicKey) Verify(msg []byte, sig Signature) bool {
	return sig.data.Signature.VerifyByte(pb.data.PublicKey, Hash256(msg))
}

func (pb *PublicKey) EqualsTo(right PublicKey) bool {
	return pb.data.PublicKey.IsEqual(right.data.PublicKey)
}

func (pb PublicKey) Address() Address {
	addr := new(Address)
	copy(addr.data.Address[:], Hash160(Hash256(pb.RawBytes())))
	return *addr
}
