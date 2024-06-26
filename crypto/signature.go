package crypto

import "io"

type Signature interface {
	Bytes() []byte
	String() string
	Hex() string
	MarshalCBOR() ([]byte, error)
	UnmarshalCBOR([]byte) error
	Encode(io.Writer) error
	Decode(io.Reader) error
	EqualsTo(right Signature) bool
}
