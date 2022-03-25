package crypto

import "io"

type PublicKey interface {
	RawBytes() []byte
	String() string
	MarshalJSON() ([]byte, error) // TODO: remove me
	MarshalCBOR() ([]byte, error)
	UnmarshalCBOR([]byte) error
	Encode(io.Writer) error
	Decode(io.Reader) error
	SanityCheck() error
	Verify(msg []byte, sig Signature) bool
	Address() Address
	VerifyAddress(addr Address) bool
	EqualsTo(right PublicKey) bool
}
