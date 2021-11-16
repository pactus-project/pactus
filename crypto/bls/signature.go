package bls

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/herumi/bls-go-binary/bls"
	"github.com/zarbchain/zarb-go/crypto"
)

const SignatureSize = 48

type BLSSignature struct {
	data signatureData
}

type signatureData struct {
	Signature *bls.Sign
}

func SignatureFromString(text string) (crypto.Signature, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}

	return SignatureFromRawBytes(data)
}

func SignatureFromRawBytes(data []byte) (crypto.Signature, error) {
	if len(data) != SignatureSize {
		return nil, fmt.Errorf("invalid signature")
	}
	s := new(bls.Sign)
	if err := s.Deserialize(data); err != nil {
		return nil, err
	}

	var sig BLSSignature
	sig.data.Signature = s

	if err := sig.SanityCheck(); err != nil {
		return nil, err
	}
	return &sig, nil
}

func (sig BLSSignature) RawBytes() []byte {
	if sig.data.Signature == nil {
		return nil
	}

	return sig.data.Signature.Serialize()
}

func (sig BLSSignature) String() string {
	if sig.data.Signature == nil {
		return ""
	}
	return sig.data.Signature.SerializeToHexStr()
}

func (sig BLSSignature) Fingerprint() string {
	return hex.EncodeToString(sig.RawBytes()[:6])
}

func (sig BLSSignature) MarshalText() ([]byte, error) {
	return []byte(sig.String()), nil
}

func (sig *BLSSignature) UnmarshalText(text []byte) error {
	s, err := SignatureFromString(string(text))
	if err != nil {
		return err
	}

	*sig = *s.(*BLSSignature)
	return nil
}

func (sig BLSSignature) MarshalJSON() ([]byte, error) {
	return json.Marshal(sig.String())
}

func (sig *BLSSignature) UnmarshalJSON(bz []byte) error {
	var text string
	if err := json.Unmarshal(bz, &text); err != nil {
		return err
	}
	return sig.UnmarshalText([]byte(text))
}

func (sig BLSSignature) MarshalCBOR() ([]byte, error) {
	if sig.data.Signature == nil {
		return nil, fmt.Errorf("invalid signature")
	}
	return cbor.Marshal(sig.RawBytes())
}

func (sig *BLSSignature) UnmarshalCBOR(bs []byte) error {
	var data []byte
	if err := cbor.Unmarshal(bs, &data); err != nil {
		return err
	}

	s, err := SignatureFromRawBytes(data)
	if err != nil {
		return err
	}

	*sig = *s.(*BLSSignature)
	return nil
}

func (sig BLSSignature) SanityCheck() error {
	bs := sig.RawBytes()
	if len(bs) != SignatureSize {
		return fmt.Errorf("signature should be %v bytes but it is %v bytes", SignatureSize, len(bs))
	}

	return nil
}

func (sig BLSSignature) EqualsTo(right crypto.Signature) bool {
	return sig.data.Signature.IsEqual(right.(*BLSSignature).data.Signature)
}
