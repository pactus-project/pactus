package crypto

import "io"

type Signature interface {
	RawBytes() []byte
	String() string
	MarshalJSON() ([]byte, error) // TODO: remove me
	MarshalCBOR() ([]byte, error)
	UnmarshalCBOR([]byte) error
	Encode(io.Writer) error
	Decode(io.Reader) error
	SanityCheck() error
	EqualsTo(right Signature) bool
}
