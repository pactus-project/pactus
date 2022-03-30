package crypto

import "io"

type PublicKey interface {
	Bytes() []byte
	String() string
	MarshalJSON() ([]byte, error) // TODO: remove me
	MarshalCBOR() ([]byte, error)
	UnmarshalCBOR([]byte) error
	Encode(io.Writer) error
	Decode(io.Reader) error
	SanityCheck() error
	Address() Address
	Verify(msg []byte, sig Signature) error
	VerifyAddress(addr Address) error
	EqualsTo(right PublicKey) bool
}
