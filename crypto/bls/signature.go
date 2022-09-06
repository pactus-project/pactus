package bls

import (
	"bytes"
	"encoding/hex"
	"io"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/herumi/bls-go-binary/bls"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util/encoding"
	"github.com/zarbchain/zarb-go/util/errors"
)

const SignatureSize = 48

type Signature struct {
	pointG1 *bls12381.PointG1
}

func SignatureFromString(text string) (*Signature, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidSignature, err.Error())
	}

	return SignatureFromBytes(data)
}

func SignatureFromBytes(data []byte) (*Signature, error) {
	if len(data) != SignatureSize {
		return nil, errors.Errorf(errors.ErrInvalidSignature,
			"signature should be %d bytes, but it is %v bytes", SignatureSize, len(data))
	}
	pointG1, err := g1.FromCompressed(data)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidSignature, err.Error())
	}

	return &Signature{pointG1: pointG1}, nil
}

func (sig *Signature) Bytes() []byte {
	return g1.ToCompressed(sig.pointG1)
}

func (sig *Signature) String() string {
	return hex.EncodeToString(sig.Bytes())
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
	if g1.IsZero(sig.pointG1) {
		return errors.Errorf(errors.ErrInvalidSignature, "signature is zero")
	}

	return nil
}

func (sig Signature) EqualsTo(right crypto.Signature) bool {
	return g1.Equal(sig.pointG1, right.(*Signature).pointG1)
}
