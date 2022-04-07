package bls

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/herumi/bls-go-binary/bls"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/encoding"
)

const SignatureSize = 48

type Signature struct {
	signature bls.Sign
}

func SignatureFromString(text string) (*Signature, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}

	return SignatureFromBytes(data)
}

func SignatureFromBytes(data []byte) (*Signature, error) {
	if len(data) != SignatureSize {
		return nil, fmt.Errorf("invalid signature")
	}
	s := new(bls.Sign)
	if err := s.Deserialize(data); err != nil {
		return nil, err
	}

	var sig Signature
	sig.signature = *s

	return &sig, nil
}

func (sig Signature) Bytes() []byte {
	return sig.signature.Serialize()
}

func (sig Signature) String() string {
	return sig.signature.SerializeToHexStr()
}

func (sig *Signature) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(sig.Bytes())
}

func (sig *Signature) UnmarshalCBOR(bs []byte) error {
	var data []byte
	if err := cbor.Unmarshal(bs, &data); err != nil {
		return err
	}

	return sig.Decode(bytes.NewReader(data))
}

func (sig *Signature) Encode(w io.Writer) error {
	return encoding.WriteElements(w, sig.Bytes())
}

func (sig *Signature) Decode(r io.Reader) error {
	data := make([]byte, SignatureSize)
	err := encoding.ReadElements(r, data)
	if err != nil {
		return err
	}

	p, err := SignatureFromBytes(data)
	if err != nil {
		return err
	}
	*sig = *p
	return nil
}

func (sig *Signature) SanityCheck() error {
	if sig.signature.IsZero() {
		return fmt.Errorf("signature is zero")
	}

	return nil
}

func (sig Signature) EqualsTo(right crypto.Signature) bool {
	return sig.signature.IsEqual(&right.(*Signature).signature)
}
