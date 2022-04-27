package crypto

import "io"

type Signature interface {
	Bytes() []byte
	String() string
	MarshalCBOR() ([]byte, error)
	UnmarshalCBOR([]byte) error
	Encode(io.Writer) error
	Decode(io.Reader) error
	SanityCheck() error
	EqualsTo(right Signature) bool
}
