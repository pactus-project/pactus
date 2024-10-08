package crypto

import "io"

type PublicKey interface {
	Bytes() []byte
	String() string
	MarshalCBOR() ([]byte, error)
	UnmarshalCBOR([]byte) error
	Encode(io.Writer) error
	Decode(io.Reader) error
	SerializeSize() int
	Verify(msg []byte, sig Signature) error
	VerifyAddress(addr Address) error
	EqualsTo(right PublicKey) bool
}
