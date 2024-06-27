package bls

import (
	"bytes"
	"encoding/hex"
	"fmt"
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
	pointG1 *bls12381.PointG1 // Lazily initialized point on G1.
	data    []byte            // Raw signature data.
}

// SignatureFromString decodes the input string and returns the Signature
// if the string is a valid hex encoding of a BLS signature.
func SignatureFromString(text string) (*Signature, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}

	return SignatureFromBytes(data)
}

// SignatureFromBytes constructs a BLS signature from the raw bytes.
func SignatureFromBytes(data []byte) (*Signature, error) {
	if len(data) != SignatureSize {
		return nil, errors.Errorf(errors.ErrInvalidSignature,
			"signature should be %d bytes, but it is %v bytes", SignatureSize, len(data))
	}

	return &Signature{data: data}, nil
}

// Bytes returns the raw byte representation of the signature.
func (sig *Signature) Bytes() []byte {
	return sig.data
}

// String returns the hex-encoded string representation of the signature.
func (sig *Signature) String() string {
	return hex.EncodeToString(sig.Bytes())
}

// MarshalCBOR encodes the signature into CBOR format.
func (sig *Signature) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(sig.Bytes())
}

// UnmarshalCBOR decodes the signature from CBOR format.
func (sig *Signature) UnmarshalCBOR(bs []byte) error {
	var data []byte
	if err := cbor.Unmarshal(bs, &data); err != nil {
		return err
	}

	return sig.Decode(bytes.NewReader(data))
}

// Encode writes the raw bytes of the signature to the provided writer.
func (sig *Signature) Encode(w io.Writer) error {
	return encoding.WriteElements(w, sig.Bytes())
}

// Decode reads the raw bytes of the signature from the provided reader and initializes the signature.
func (sig *Signature) Decode(r io.Reader) error {
	data := make([]byte, SignatureSize)
	err := encoding.ReadElements(r, data)
	if err != nil {
		return err
	}

	s, _ := SignatureFromBytes(data)
	*sig = *s

	return nil
}

// EqualsTo checks if the current signature is equal to another signature.
func (sig *Signature) EqualsTo(right crypto.Signature) bool {
	return bytes.Equal(sig.data, right.(*Signature).data)
}

// PointG1 returns the point on G1 for the signature.
func (sig *Signature) PointG1() (bls12381.PointG1, error) {
	if sig.pointG1 != nil {
		return *sig.pointG1, nil
	}

	g1 := bls12381.NewG1()
	pointG1, err := g1.FromCompressed(sig.data)
	if err != nil {
		return bls12381.PointG1{}, err
	}
	if g1.IsZero(pointG1) {
		return bls12381.PointG1{}, fmt.Errorf("signature is zero")
	}

	sig.pointG1 = pointG1

	return *pointG1, nil
}
