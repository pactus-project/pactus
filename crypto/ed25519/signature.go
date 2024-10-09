package ed25519

import (
	"bytes"
	"crypto/subtle"
	"encoding/hex"
	"io"

	cbor "github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/encoding"
)

var _ crypto.Signature = &Signature{}

const SignatureSize = 64

type Signature struct {
	data []byte
}

// SignatureFromString decodes the input string and returns the Signature
// if the string is a valid hexadecimal encoding of a Ed25519 signature.
func SignatureFromString(text string) (*Signature, error) {
	data, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}

	return SignatureFromBytes(data)
}

// SignatureFromBytes constructs a Ed25519 signature from the raw bytes.
func SignatureFromBytes(data []byte) (*Signature, error) {
	if len(data) != SignatureSize {
		return nil, crypto.InvalidLengthError(len(data))
	}

	return &Signature{data: data}, nil
}

// Bytes returns the raw byte representation of the signature.
func (sig *Signature) Bytes() []byte {
	return sig.data[:SignatureSize]
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

func (*Signature) SerializeSize() int {
	return SignatureSize
}

// EqualsTo checks if the current signature is equal to another signature.
func (sig *Signature) EqualsTo(x crypto.Signature) bool {
	xEd25519, ok := x.(*Signature)
	if !ok {
		return false
	}

	return subtle.ConstantTimeCompare(sig.data, xEd25519.data) == 1
}
