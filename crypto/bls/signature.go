package bls

import (
	"bytes"
	"crypto/subtle"
	"encoding/hex"
	"io"
	"sync/atomic"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	cbor "github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/encoding"
)

var _ crypto.Signature = &Signature{}

const SignatureSize = 48

type Signature struct {
	pointG1 atomic.Pointer[bls12381.G1Affine] // Lazily initialized point on G1.
	data    []byte                            // Raw signature data.
}

// SignatureFromString decodes the input string and returns the Signature
// if the string is a valid hexadecimal encoding of a BLS signature.
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
		return nil, crypto.InvalidLengthError(len(data))
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

	// data is always exactly SignatureSize bytes here, so this can only
	// fail on length, which never happens; pointG1 is lazy and never
	// populated by SignatureFromBytes, so there's nothing to carry over.
	sig.data = data
	sig.pointG1.Store(nil)

	return nil
}

func (*Signature) SerializeSize() int {
	return SignatureSize
}

// EqualsTo checks if the current signature is equal to another signature.
func (sig *Signature) EqualsTo(x crypto.Signature) bool {
	xBLS, ok := x.(*Signature)
	if !ok || xBLS == nil {
		return false
	}

	return subtle.ConstantTimeCompare(sig.data, xBLS.data) == 1
}

// PointG1 returns the point on G1 for the signature.
func (sig *Signature) PointG1() (*bls12381.G1Affine, error) {
	if g1Aff := sig.pointG1.Load(); g1Aff != nil {
		return g1Aff, nil
	}

	g1Aff := new(bls12381.G1Affine)
	err := g1Aff.Unmarshal(sig.data)
	if err != nil {
		return nil, err
	}
	if g1Aff.IsInfinity() || !g1Aff.IsInSubGroup() {
		return nil, crypto.ErrInvalidSignature
	}

	// Concurrent calls may each compute g1Aff (the result is deterministic,
	// so a lost race just means a harmless recomputation); atomic.Pointer
	// makes the eventual read/write itself race-free.
	sig.pointG1.Store(g1Aff)

	return g1Aff, nil
}
