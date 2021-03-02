package crypto

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/herumi/bls-go-binary/bls"
)

const PrivateKeySize = 32

type PrivateKey struct {
	data privateKeyData
}

type privateKeyData struct {
	SecretKey *bls.SecretKey
}

func PrivateKeyFromString(text string) (PrivateKey, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return PrivateKey{}, err
	}

	return PrivateKeyFromRawBytes(data)
}

func PrivateKeyFromSeed(seed []byte) (PrivateKey, error) {
	sc := new(bls.SecretKey)
	err := sc.SetLittleEndianMod(seed)
	if err != nil {
		return PrivateKey{}, err
	}

	var pv PrivateKey
	pv.data.SecretKey = sc

	return pv, nil
}

func PrivateKeyFromRawBytes(data []byte) (PrivateKey, error) {
	if len(data) != PrivateKeySize {
		return PrivateKey{}, fmt.Errorf("Invalid private key")
	}
	sc := new(bls.SecretKey)
	if err := sc.Deserialize(data); err != nil {
		return PrivateKey{}, err
	}

	var pv PrivateKey
	pv.data.SecretKey = sc

	return pv, nil
}

/// -------
/// CASTING
func (pv PrivateKey) RawBytes() []byte {
	if pv.data.SecretKey == nil {
		return nil
	}
	return pv.data.SecretKey.Serialize()
}

func (pv PrivateKey) String() string {
	if pv.data.SecretKey == nil {
		return ""
	}
	return pv.data.SecretKey.SerializeToHexStr()
}

/// ----------
/// MARSHALING

func (pv PrivateKey) MarshalText() ([]byte, error) {
	return []byte(pv.String()), nil
}

func (pv *PrivateKey) UnmarshalText(text []byte) error {
	p, err := PrivateKeyFromString(string(text))
	if err != nil {
		return err
	}

	*pv = p
	return nil
}

func (pv PrivateKey) MarshalJSON() ([]byte, error) {
	bz, err := pv.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(bz))
}

func (pv *PrivateKey) UnmarshalJSON(bz []byte) error {
	var text string
	if err := json.Unmarshal(bz, &text); err != nil {
		return err
	}
	return pv.UnmarshalText([]byte(text))
}

func (pv PrivateKey) MarshalCBOR() ([]byte, error) {
	if pv.data.SecretKey == nil {
		return nil, fmt.Errorf("Invalid private key")
	}
	return cbor.Marshal(pv.RawBytes())
}

func (pv *PrivateKey) UnmarshalCBOR(bs []byte) error {
	var data []byte
	if err := cbor.Unmarshal(bs, &data); err != nil {
		return err
	}

	p, err := PrivateKeyFromRawBytes(data)
	if err != nil {
		return err
	}

	*pv = p
	return nil
}

func (pv *PrivateKey) SanityCheck() error {
	bs := pv.RawBytes()
	if len(bs) != PrivateKeySize {
		return fmt.Errorf("Private key should be %v bytes but it is %v bytes", PrivateKeySize, len(bs))
	}
	return nil
}

func (pv *PrivateKey) Sign(msg []byte) Signature {
	sig := new(Signature)
	sig.data.Signature = pv.data.SecretKey.SignByte(Hash256(msg))

	return *sig
}

func (pv *PrivateKey) PublicKey() PublicKey {
	pb := new(PublicKey)
	pb.data.PublicKey = pv.data.SecretKey.GetPublicKey()

	return *pb
}

func (pv *PrivateKey) EqualsTo(right PrivateKey) bool {
	return pv.data.SecretKey.IsEqual(right.data.SecretKey)
}

func RandomKeyPair() (Address, PublicKey, PrivateKey) {
	pv := new(PrivateKey)
	pv.data.SecretKey = new(bls.SecretKey)
	pv.data.SecretKey.SetByCSPRNG()

	return pv.PublicKey().Address(), pv.PublicKey(), *pv
}

// ---------
// For tests
func GenerateTestKeyPair() (Address, PublicKey, PrivateKey) {
	return RandomKeyPair()
}
