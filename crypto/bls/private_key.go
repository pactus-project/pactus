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

const PrivateKeySize = 32

type BLSPrivateKey struct {
	data privateKeyData
}

type privateKeyData struct {
	SecretKey *bls.SecretKey
}

func PrivateKeyFromString(text string) (*BLSPrivateKey, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}

	return PrivateKeyFromRawBytes(data)
}

func PrivateKeyFromSeed(seed []byte) (*BLSPrivateKey, error) {
	sc := new(bls.SecretKey)
	err := sc.SetLittleEndianMod(seed)
	if err != nil {
		return nil, err
	}

	var pv BLSPrivateKey
	pv.data.SecretKey = sc

	return &pv, nil
}

func PrivateKeyFromRawBytes(data []byte) (*BLSPrivateKey, error) {
	if len(data) != PrivateKeySize {
		return nil, fmt.Errorf("invalid private key")
	}
	sc := new(bls.SecretKey)
	if err := sc.Deserialize(data); err != nil {
		return nil, err
	}

	var pv BLSPrivateKey
	pv.data.SecretKey = sc

	return &pv, nil
}

func (pv BLSPrivateKey) RawBytes() []byte {
	if pv.data.SecretKey == nil {
		return nil
	}
	return pv.data.SecretKey.Serialize()
}

func (pv BLSPrivateKey) String() string {
	if pv.data.SecretKey == nil {
		return ""
	}
	return pv.data.SecretKey.SerializeToHexStr()
}

func (pv *BLSPrivateKey) MarshalText() ([]byte, error) {
	return []byte(pv.String()), nil
}

func (pv *BLSPrivateKey) UnmarshalText(text []byte) error {
	p, err := PrivateKeyFromString(string(text))
	if err != nil {
		return err
	}

	*pv = *p
	return nil
}

func (pv *BLSPrivateKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pv.String())
}

func (pv *BLSPrivateKey) UnmarshalJSON(bz []byte) error {
	var text string
	if err := json.Unmarshal(bz, &text); err != nil {
		return err
	}
	return pv.UnmarshalText([]byte(text))
}

func (pv *BLSPrivateKey) MarshalCBOR() ([]byte, error) {
	if pv.data.SecretKey == nil {
		return nil, fmt.Errorf("invalid private key")
	}
	return cbor.Marshal(pv.RawBytes())
}

func (pv *BLSPrivateKey) UnmarshalCBOR(bs []byte) error {
	var data []byte
	if err := cbor.Unmarshal(bs, &data); err != nil {
		return err
	}

	p, err := PrivateKeyFromRawBytes(data)
	if err != nil {
		return err
	}

	*pv = *p
	return nil
}

func (pv *BLSPrivateKey) SanityCheck() error {
	bs := pv.RawBytes()
	if len(bs) != PrivateKeySize {
		return fmt.Errorf("private key should be %v bytes but it is %v bytes", PrivateKeySize, len(bs))
	}
	return nil
}

func (pv *BLSPrivateKey) Sign(msg []byte) crypto.Signature {
	sig := new(BLSSignature)
	sig.data.Signature = pv.data.SecretKey.SignByte(hash.Hash256(msg))

	return sig
}

func (pv *BLSPrivateKey) PublicKey() crypto.PublicKey {
	pb := new(BLSPublicKey)
	pb.data.PublicKey = pv.data.SecretKey.GetPublicKey()

	return pb
}

func (pv *BLSPrivateKey) EqualsTo(right crypto.PrivateKey) bool {
	return pv.data.SecretKey.IsEqual(right.(*BLSPrivateKey).data.SecretKey)
}
