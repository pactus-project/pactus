package bls

import (
	"bytes"
	"encoding/hex"
	"io"

	cbor "github.com/fxamacker/cbor/v2"
	bls12381 "github.com/kilic/bls12-381"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/encoding"
	"github.com/pactus-project/pactus/util/errors"
)

var _ crypto.Signature = &Signature{}

const SignatureSize = 48

type Signature struct {
	pointG1 bls12381.PointG1
}

func SignatureFromString(text string) (*Signature, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidSignature, err.Error())
	}

	return SignatureFromBytes(data)
}

// SignatureFromBytes constructs a BLS signature from the raw bytes.
func SignatureFromBytes(data []byte) (*Signature, error) {
	if len(data) != SignatureSize {
		return nil, errors.Errorf(errors.ErrInvalidSignature,
			"signature should be %d bytes, but it is %v bytes", SignatureSize, len(data))
	}
	g1 := bls12381.NewG1()

	pointG1, err := g1.FromCompressed(data)
	if err != nil {
		return nil, errors.Errorf(errors.ErrInvalidSignature, err.Error())
	}
	if g1.IsZero(pointG1) {
		return nil, errors.Errorf(errors.ErrInvalidSignature,
			"signature is zero")
	}

	return &Signature{pointG1: *pointG1}, nil
}

func (sig *Signature) Bytes() []byte {
	g1 := bls12381.NewG1()

	return g1.ToCompressed(sig.point())
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

func (sig *Signature) EqualsTo(right crypto.Signature) bool {
	g1 := bls12381.NewG1()

	return g1.Equal(sig.point(), right.(*Signature).point())
}

// clonePoint clones the pointG1 to make sure it remains intact.
func (sig *Signature) point() *bls12381.PointG1 {
	return bls12381.NewG1().New().Set(&sig.pointG1)
}
